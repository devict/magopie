package magopie

import "fmt"

// Site defines a site that serves torrent files.
type Site struct {
	ID      string
	Name    string
	URL     string
	Enabled bool
}

// SiteCollection is a collection of sites because gomobile can't
// handle slices
type SiteCollection struct {
	list []Site
}

// Length returns how many sites are in the collection
func (sc *SiteCollection) Length() int {
	return len(sc.list)
}

// Get returns the site at idx or a nil site
func (sc *SiteCollection) Get(idx int) *Site {
	if idx <= sc.Length() {
		return &sc.list[idx]
	}
	return nil
}

// Clear empties the list of sites
func (sc *SiteCollection) Clear() {
	sc.list = sc.list[:0]
}

// Index finds the index of a site or -1 if not found
func (sc *SiteCollection) Index(s *Site) int {
	for i, tst := range sc.list {
		if tst == *s {
			return i
		}
	}
	return -1
}

// Insert inserts a site into the collection at i
func (sc *SiteCollection) Insert(i int, s *Site) {
	if i < 0 || i > sc.Length() {
		fmt.Printf("Magopie-go:: Attempted to insert a site at an invalid index")
		return
	}
	sc.list = append(sc.list, Site{})
	copy(sc.list[i+1:], sc.list[i:])
	sc.list[i] = *s
}

// Remove a site from the collection at i
func (sc *SiteCollection) Remove(i int) {
	if i < 0 || i > sc.Length() {
		fmt.Printf("Magopie-go:: Attempted to remove a site from an invalid index")
		return
	}
	copy(sc.list[i:], sc.list[i+1:])
	sc.list[len(sc.list)-1] = Site{}
	sc.list = sc.list[:len(sc.list)-1]
}
