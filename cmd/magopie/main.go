package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gophergala2016/magopie"
)

func main() {
	// TODO populate with registered sites
	a := &app{}

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

	r.HandleFunc("/sites", a.handleSites)
	r.HandleFunc("/sites/", a.handleSites)
	r.HandleFunc("/torrents", a.handleTorrents)

	return r
}

type app struct {
	sites sitedb
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
