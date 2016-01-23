package main

import (
	"encoding/xml"
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
	var d struct {
		Items []struct {
			Title         string `xml:"title"`
			ContentLength int    `xml:"contentLength"`
			Hash          string `xml:"infoHash"`
			Seeds         int    `xml:"seeds"`
			Peers         int    `xml:"peers"`
			Enclosure     struct {
				URL string `xml:"url,attr"`
			} `xml:"enclosure"`
		} `xml:"channel>item"`
	}

	if err := xml.NewDecoder(r).Decode(&d); err != nil {
		return nil, err
	}

	ts := make([]magopie.Torrent, len(d.Items))
	for i, item := range d.Items {
		ts[i] = magopie.Torrent{
			ID:       item.Hash,
			Title:    item.Title,
			Seeders:  item.Seeds,
			Leechers: item.Peers - item.Seeds,
			Size:     item.ContentLength,
			FileURL:  item.Enclosure.URL,
			SiteID:   "kat",
		}
	}

	return ts, nil
}
