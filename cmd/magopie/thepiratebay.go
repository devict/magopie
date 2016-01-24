package main

import (
	"io"
	"log"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/dustin/go-humanize"
	"github.com/gophergala2016/magopie"
)

var thePirateBay = site{
	Site: magopie.Site{
		ID:      "tpb",
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

func hashFromMagnet(m string) (string, error) {
	re := regexp.MustCompile(`magnet:\?xt=urn:btih:([^&]+)`)
	match := re.FindStringSubmatch(m)

	if len(match) <= 0 {
		return "", ErrHashlessMagnet
	}

	return match[1], nil
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
