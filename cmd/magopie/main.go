package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
