package main

import (
	"io"

	"github.com/gophergala2016/magopie"
)

var thePirateBay = site{
	Site: magopie.Site{
		ID:      "the_pirate_bay",
		Name:    "The Pirate Bay",
		URL:     "https://thepiratebay.se",
		Enabled: true,
	},
	search: func(term string) ([]magopie.Torrent, error) {
		// TODO make http request the pass body to tpbParse()
		return nil, nil
	},
}

func tpbParse(r io.Reader) ([]magopie.Torrent, error) {
	// TODO implement
	return nil, nil
}
