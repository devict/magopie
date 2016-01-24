package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gophergala2016/magopie"
	"github.com/spf13/afero"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	a := &app{
		fs: afero.OsFs{},
		sites: sitedb{
			sites: []site{
				kickAssTorrents,
			},
		},
		torcacheURL: "http://torcache.net",
		downloadDir: ".",
	}

	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf("%s:%s", "localhost", port)
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router(a)))
}

func router(a *app) http.Handler {
	r := http.NewServeMux()

	// TODO switch to a muxer where I define the method as part of the route
	r.HandleFunc("/sites", a.handleSites)
	r.HandleFunc("/sites/", a.handleSites)
	r.HandleFunc("/torrents", a.handleTorrents)
	r.HandleFunc("/download/", a.handleDownload)

	return r
}

type app struct {
	// sites is the collection of upstream sites registered for this app
	sites sitedb

	// torcacheURL is the base URL for the service that serves torrent files
	torcacheURL string

	// fs is the filesystem where we'll store .torrent files
	fs afero.Fs

	// downloadDir is the path to the the dir in fs where we'll store files
	downloadDir string
}

func (a *app) handleSites(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	pfx := len("/sites/")

	if len(r.URL.Path) <= pfx {
		if err := enc.Encode(a.sites.GetAllSites()); err != nil {
			log.Println(err)
		}
		return
	}

	id := r.URL.Path[pfx:]
	if err := enc.Encode(a.sites.GetSite(id)); err != nil {
		log.Println(err)
	}
}

func (a *app) handleTorrents(w http.ResponseWriter, r *http.Request) {
	// TODO better error response?
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

	if err := json.NewEncoder(w).Encode(torrents); err != nil {
		log.Println(err)
	}
}

func (a *app) handleDownload(w http.ResponseWriter, r *http.Request) {
	pfx := len("/download/")
	if len(r.URL.Path) <= pfx {
		// TODO what to do here?
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := r.URL.Path[pfx:]
	url := fmt.Sprintf(
		"%s/torrent/%s.torrent",
		a.torcacheURL,
		url.QueryEscape(strings.ToUpper(hash)),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
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

	err = a.fs.MkdirAll(a.downloadDir, 0755)
	if err != nil {
		log.Print("Error making download dir", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	path := filepath.Join(a.downloadDir, hash+".torrent")
	file, err := a.fs.Create(path)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(file, res.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
