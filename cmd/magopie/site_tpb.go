package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/dustin/go-humanize"
	"github.com/devict/magopie"
)

var thePirateBay = site{
	Site: magopie.Site{
		ID:      "tpb",
		Name:    "The Pirate Bay",
		URL:     "https://thepiratebay.se",
		Enabled: true,
	},
	search: func(term string) ([]magopie.Torrent, error) {
		url := fmt.Sprintf(
			"https://thepiratebay.se/search/%s/0/7/0",
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

		return tpbParse(res.Body)
	},
}

func tpbParse(r io.Reader) ([]magopie.Torrent, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	var torrents []magopie.Torrent
	doc.Find("table#searchResult tr").Each(func(i int, s *goquery.Selection) {
		// First child is header row
		if i <= 0 {
			return
		}

		magnet, ok := s.Find("div.detName + a").Attr("href")
		if !ok {
			return
		}

		// hash gets ripped from magnet
		hash, err := hashFromMagnet(magnet)
		if err != nil {
			log.Print(err)
			return
		}

		seed, err := strconv.Atoi(s.Children().Eq(2).Text())
		if err != nil {
			log.Print(err)
			seed = 0
		}

		leech, err := strconv.Atoi(s.Children().Eq(3).Text())
		if err != nil {
			log.Print(err)
			leech = 0
		}

		// Need to rip human-friendly filesize from description and make it ugly
		details := s.Find("font.detDesc").Text()
		size, err := bytesFromDetails(details)
		if err != nil {
			log.Print(err)
			size = 0
		}

		t := magopie.Torrent{
			ID:        hash,
			Title:     s.Find("div.detName a").Text(),
			MagnetURI: magnet,
			SiteID:    "tpb",
			Seeders:   seed,
			Leechers:  leech,
			Size:      size,
		}

		torrents = append(torrents, t)
	})

	return torrents, nil
}

// btihRE finds the BTIH hash from the xt segment of a Magnet URI
var btihRE = regexp.MustCompile(`urn:btih:(.+)`)

func hashFromMagnet(m string) (string, error) {
	parsed, err := url.Parse(m)
	if err != nil {
		return "", err
	}

	matches := btihRE.FindStringSubmatch(parsed.Query().Get("xt"))
	if len(matches) <= 0 {
		return "", ErrHashlessMagnet
	}

	return matches[1], nil
}

func bytesFromDetails(d string) (int, error) {
	re := regexp.MustCompile(`, Size ([^,]+)`)
	match := re.FindStringSubmatch(d)

	if len(match) <= 0 {
		return 0, ErrCannotFindFileSize
	}

	bytes, err := humanize.ParseBigBytes(match[1])
	if err != nil {
		log.Print(err)
		return 0, err
	}

	return int(bytes.Int64()), nil
}
