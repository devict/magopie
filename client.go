package magopie

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

// Client provides methods for talking to a magopie server.
type Client struct {
	ServerAddr string
}

// NewClient creates a Client
func NewClient(server string) *Client {
	return &Client{ServerAddr: server}
}

// Search asks the Magopie server for a list of results
func (c *Client) Search(s string) *TorrentCollection {
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

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(res.Body); err != nil {
		// TODO can we return an error?
		log.Print(err)
		return ret
	}

	if err := ret.UnmarshalJSON(buf.Bytes()); err != nil {
		// TODO can we return an error?
		log.Print(err)
		return ret
	}

	return ret
}

// Download triggers the Magopie server to download
func (c *Client) Download(t *Torrent) bool {
	fmt.Println("Magopie-go: Triggering Download")

	url := fmt.Sprintf("%s/download/%s", c.ServerAddr, t.ID)
	res, err := http.Post(url, "", nil)
	if err != nil {
		log.Print(err)
		return false
	}

	if res.StatusCode != http.StatusCreated {
		log.Print("Download request responded with status", res.StatusCode)
		return false
	}

	return true
}
