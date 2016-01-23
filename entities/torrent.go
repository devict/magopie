package entities

// A Torrent is an individual result from a search operation representing a
// single torrent file.
type Torrent struct {
	ID      string
	Title   string
	FileURL string

	// Fields we hopefully can populate
	Description string
	Seeders     int
	Leachers    int
	Size        int
}
