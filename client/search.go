package client

import (
	"fmt"

	"github.com/gophergala2016/magopie/entities"
)

// TorrentCollection ...
type TorrentCollection struct {
	list []entities.Torrent
}

// GetResultsLength returns the number of results in the TorrentCollection
func (tc *TorrentCollection) GetResultsLength() int {
	return len(tc.list)
}

// GetResultID ...
func (tc *TorrentCollection) GetResultID(idx int) string {
	if idx <= len(tc.list) {
		return tc.list[idx].ID
	}
	return "Index not found"
}

// Greetings ...
func Greetings(n string) string {
	return fmt.Sprintf("Hello, %s!", n)
}

// Search ...
func Search(s string) *TorrentCollection {
	ret := &TorrentCollection{}
	t := &entities.Torrent{}
	t.ID = "Banana"
	ret.list = append(ret.list, *t)
	return ret
}