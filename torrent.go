package magopie

// A Torrent is an individual result from a search operation representing a
// single torrent file.
type Torrent struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	FileURL string `json:"fileURL"`
	Site    Site   `json:"site"`

	// Fields we hopefully can populate
	Description string `json:"description"`
	Seeders     int    `json:"seeders"`
	Leechers    int    `json:"leechers"`
	Size        int    `json:"size"`
}
