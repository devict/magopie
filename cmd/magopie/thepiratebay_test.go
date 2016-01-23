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
			ID:    "???",
			Title: "???",
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("tpbParse actual = %v\nexpected %v", actual, expected)
	}
}
