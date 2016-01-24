package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gophergala2016/magopie"
	"github.com/rs/xmux"
	"github.com/spf13/afero"
	"golang.org/x/net/context"
)

type server struct {
	// sites is the collection of upstream sites registered for this app
	sites sitedb

	// torcacheURL is the base URL for the service that serves torrent files
	torcacheURL string

	// fs is the filesystem where we'll store .torrent files
	fs afero.Fs

	// downloadDir is the path to the the dir in fs where we'll store files
	downloadDir string
}

func (a *server) handleAllSites(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(a.sites.GetAllSites())
	if err != nil {
		log.Println(err)
	}
}

func (a *server) handleSingleSite(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := xmux.Param(ctx, "id")
	err := json.NewEncoder(w).Encode(a.sites.GetSite(id))
	if err != nil {
		log.Println(err)
	}
}

func (a *server) handleTorrents(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	term := r.FormValue("q")
	if term == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO do this concurrently
	torrents := []magopie.Torrent{}
	for _, s := range a.sites.GetEnabledSites() {
		results, err := s.search(term)
		if err != nil {
			log.Printf("Error searching %q on %v: %v", term, s.ID, err)
			continue
		}
		for _, t := range results {
			torrents = append(torrents, t)
		}
	}

	sort.Sort(magopie.BySeeders(torrents))

	if err := json.NewEncoder(w).Encode(torrents); err != nil {
		log.Println(err)
	}
}

func (a *server) handleDownload(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	hash := strings.ToUpper(xmux.Param(ctx, "hash"))
	url := fmt.Sprintf(
		"%s/torrent/%s.torrent",
		a.torcacheURL,
		url.QueryEscape(hash),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Header.Set("User-Agent", "Magopie Server")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// Ensure download directory exists
	// TODO should this be done every request?
	err = a.fs.MkdirAll(a.downloadDir, 0755)
	if err != nil {
		log.Print("Error making download dir", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create file to write torrent
	// TODO create with .part or something then move when done
	path := filepath.Join(a.downloadDir, hash+".torrent")
	file, err := a.fs.Create(path)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Copy file contents from torcache response to os file
	_, err = io.Copy(file, res.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
