package magopie

import (
	"fmt"
	"log"
	"net/http"
)

// Client provides methods for talking to a magopie server.
type Client struct {
	ServerAddr string
}

// NewClient creates a Client
func NewClient(server string) Client {
	return Client{ServerAddr: server}
}

// Search asks the Magopie server for a list of results
func (c Client) Search(s string) *TorrentCollection {
	ret := &TorrentCollection{}

	req, err := http.NewRequest("GET", c.ServerAddr+"/torrents", nil)
	if err != nil {
		// TODO can we return an error?
		log.Print(err)
		return ret
	}
	vals := req.URL.Query()
	vals.Add("q", s)
	req.URL.RawQuery = vals.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// TODO can we return an error?
		log.Print(err)
		return ret
	}
	defer res.Body.Close()

	if err := ret.UnmarshalJSONReader(res.Body); err != nil {
		// TODO can we return an error?
		log.Print(err)
		return ret
	}

	return ret
}

// Download triggers the Magopie server to download
func (c Client) Download(t *Torrent) bool {
	fmt.Println("Magopie-go: Triggering Download")
	return true
}
