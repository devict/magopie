package magopie

import (
	"fmt"
)

// Greetings ...
func Greetings(n string) string {
	return fmt.Sprintf("Hello, %s!", n)
}

// Search ...
func Search(s string) *TorrentCollection {
	ret := &TorrentCollection{}
	t := &Torrent{}
	t.ID = "Banana"
	ret.list = append(ret.list, *t)
	return ret
}
