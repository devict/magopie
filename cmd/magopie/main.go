package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gophergala2016/magopie/entities"
)

func main() {
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf("%s:%s", "localhost", port)
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router()))
}

func router() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/sites", handleSearch)

	return r
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	// TODO get list of sites from somewhere
	s := []entities.Site{
		{},
	}

	if err := json.NewEncoder(w).Encode(s); err != nil {
		log.Println(err)
		return
	}
}
