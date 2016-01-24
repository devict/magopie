// The magopie server aggregates search results from various torrent sites
// and downloads specified torrent files.
//
// Available routes are:
//
// GET    /sites
// GET    /sites/{id}
//*POST   /sites/{id}/enabled
//*DELETE /sites/{id}/enabled

// GET    /torrents?q=ubuntu
// POST   /download/{id}
package main
