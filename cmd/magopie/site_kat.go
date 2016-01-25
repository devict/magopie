package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/devict/magopie"
)

var kickAssTorrents = site{
	Site: magopie.Site{
		ID:      "kat",
		Name:    "Kick Ass Torrents",
		URL:     "https://kat.ph",
		Enabled: true,
	},
	search: func(term string) ([]magopie.Torrent, error) {
		url := fmt.Sprintf(
			"https://kat.cr/usearch/%s/?rss=1",
			url.QueryEscape(term),
		)

		res, err := http.DefaultClient.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return nil, ErrFailedRequest
		}

		return katParse(res.Body)
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
			MagnetURI     string `xml:"magnetURI"`
		} `xml:"channel>item"`
	}

	if err := xml.NewDecoder(r).Decode(&d); err != nil {
		return nil, err
	}

	ts := make([]magopie.Torrent, len(d.Items))
	for i, item := range d.Items {
		ts[i] = magopie.Torrent{
			ID:        item.Hash,
			Title:     item.Title,
			Seeders:   item.Seeds,
			Leechers:  item.Peers - item.Seeds,
			Size:      item.ContentLength,
			MagnetURI: item.MagnetURI,
			SiteID:    "kat",
		}
	}

	return ts, nil
}
