package magopie

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

// Search asks the Magopie server for a list of results
func Search(s string) *TorrentCollection {
	ret := &TorrentCollection{}
	for i := 0; i < 100; i++ {
		ai := strconv.Itoa(i)
		ret.Push(&Torrent{
			ID:        s + "-" + ai,
			Title:     "the-" + s + "-" + ai,
			MagnetURI: s + "-" + ai,
			SiteID:    "the_pirate_bay",
		})
	}
	return ret
}

// Download triggers the Magopie server to download
func Download(t *Torrent) bool {
	fmt.Println("Magopie-go: Triggering Download")
	return true
}

// SaveToFile saves 'data' to the file 'filePath'
func SaveToFile(filePath string, data []byte) error {
	return ioutil.WriteFile(filePath, data, 0644)
}

// ReadFromFile returns the bytes from 'filePath'
func ReadFromFile(filePath string) []byte {
	var data []byte
	var err error
	if data, err = ioutil.ReadFile(filePath); err != nil {
		fmt.Println(err)
	}
	return data
}
