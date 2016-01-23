package entities

// Site defines a site that serves torrent files.
type Site struct {
	ID      string
	Name    string
	URL     string
	Enabled bool
}
