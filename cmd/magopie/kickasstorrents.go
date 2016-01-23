package main

import (
	"io"

	"github.com/gophergala2016/magopie"
)

var kickAssTorrents = site{
	Site: magopie.Site{
		ID:      "kat",
		Name:    "Kick Ass Torrents",
		URL:     "https://kat.ph",
		Enabled: true,
	},
	search: func(term string) ([]magopie.Torrent, error) {
		// TODO make http request the pass body to katParse()
		return nil, nil
	},
}

func katParse(r io.Reader) ([]magopie.Torrent, error) {
	// TODO implement
	return nil, nil
}
