// The magopie server aggregates search results from various torrent sites
// and downloads specified torrent files.
//
// Available routes are:
//
// GET    /sites
// GET    /sites/{id}
// POST   /sites/{id}/enabled
// DELETE /sites/{id}/enabled

// GET    /torrents?q=ubuntu
// GET    /torrents/{id}
// POST   /torrents/{id}/download
// DELETE /torrents/{id}/download

// GET    /downloads
// GET    /downloads/{id}
// POST   /downloads/{id}
package main
