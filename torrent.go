package magopie

import "fmt"

// A Torrent is an individual result from a search operation representing a
// single torrent file.
type Torrent struct {
	ID        string
	Title     string
	MagnetURI string
	SiteID    string
	Seeders   int
	Leechers  int
	Size      int
}

// BySeeders implements sort.Interface for []Torrent based on Seeeders.
type BySeeders []Torrent

func (s BySeeders) Len() int           { return len(s) }
func (s BySeeders) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s BySeeders) Less(i, j int) bool { return s[i].Seeders > s[j].Seeders }

// TorrentCollection is a collection of torrents because gomobile can't
// handle slices
type TorrentCollection struct {
	list []Torrent
}

// Length returns how many torrents are in the collection
func (tc *TorrentCollection) Length() int {
	return len(tc.list)
}

// Get returns the torrent at idx or a nil torrent
func (tc *TorrentCollection) Get(idx int) *Torrent {
	if idx <= tc.Length() {
		return &tc.list[idx]
	}
	return nil
}

// Clear empties the list of torrents
func (tc *TorrentCollection) Clear() {
	tc.list = tc.list[:0]
}

// Index finds the index of a torrent or -1 if not found
func (tc *TorrentCollection) Index(t *Torrent) int {
	for i, tst := range tc.list {
		if tst == *t {
			return i
		}
	}
	return -1
}

// Insert inserts a torrent into the collection at i
func (tc *TorrentCollection) Insert(i int, t *Torrent) {
	if i < 0 || i > tc.Length() {
		fmt.Printf("Magopie-go:: Attempted to insert a torrent at an invalid index")
		return
	}
	tc.list = append(tc.list, Torrent{})
	copy(tc.list[i+1:], tc.list[i:])
	tc.list[i] = *t
}

// Remove a torrent from the collection at i
func (tc *TorrentCollection) Remove(i int) {
	if i < 0 || i >= tc.Length() {
		fmt.Printf("Magopie-go:: Attempted to remove a torrent from an invalid index")
		return
	}
	copy(tc.list[i:], tc.list[i+1:])
	tc.list[len(tc.list)-1] = Torrent{}
	tc.list = tc.list[:len(tc.list)-1]
}

// Push adds an element to the end of the collection
func (tc *TorrentCollection) Push(t *Torrent) {
	tc.Insert(tc.Length(), t)
}

// Pop removes the last element from the collection
func (tc *TorrentCollection) Pop(t *Torrent) {
	tc.Remove(tc.Length() - 1)
}

// Unshift adds an element to the front of the collection
func (tc *TorrentCollection) Unshift(t *Torrent) {
	tc.Insert(0, t)
}

// Shift removes an element from the front of the collection
func (tc *TorrentCollection) Shift(t *Torrent) {
	tc.Remove(0)
}
