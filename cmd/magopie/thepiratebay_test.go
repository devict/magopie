package main

import (
	"os"
	"reflect"
	"testing"

	mp "github.com/gophergala2016/magopie"
)

func TestTPBParse(t *testing.T) {
	data, err := os.Open("testdata/tpb_ubuntu.html")
	if err != nil {
		t.Fatal("Error opening fixture", err)
	}
	defer data.Close()

	actual, err := tpbParse(data)

	// TODO add more fields to expected slice
	expected := []mp.Torrent{
		{
			ID:       "1619ecc9373c3639f4ee3e261638f29b33a6cbd6",
			Title:    "Ubuntu 14.10 i386 (Desktop ISO)",
			FileURL:  "",
			SiteID:   "tpb",
			Seeders:  66,
			Leechers: 8,
			Size:     1191853424,
		},
		{
			ID:       "4896fde14efbc0f66a274d2a69104fbb57fbd2cb",
			Title:    "Ubuntu 15.04 Desktop i386, [Iso - MultiLang] [TNTVillage]",
			FileURL:  "",
			SiteID:   "tpb",
			Seeders:  33,
			Leechers: 3,
			Size:     1191853424,
		},
		{
			ID:       "b415c913643e5ff49fe37d304bbb5e6e11ad5101",
			Title:    "Ubuntu 14.10 desktop  x64",
			FileURL:  "",
			SiteID:   "tpb",
			Seeders:  23,
			Leechers: 1,
			Size:     1159641169,
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("tpbParse actual = %v\nexpected %v", actual, expected)
	}
}
