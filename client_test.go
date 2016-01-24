package magopie

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testKey = "Margo the magpie"

func TestClientSearch(t *testing.T) {
	torrent1 := Torrent{
		ID:        "ID 1",
		Title:     "Title 1",
		MagnetURI: "Magnet 1",
		SiteID:    "test",
		Seeders:   10,
		Leechers:  100,
		Size:      123456,
	}
	torrent2 := Torrent{
		ID:        "ID 2",
		Title:     "Title 2",
		MagnetURI: "Magnet 2",
		SiteID:    "test",
		Seeders:   20,
		Leechers:  200,
		Size:      654321,
	}

	var srvMethod, srvPath, srvQuery string
	var srvValid bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srvMethod = r.Method
		srvPath = r.URL.Path
		srvQuery = r.FormValue("q")
		srvValid = requestIsSigned(r, testKey)
		json.NewEncoder(w).Encode([]Torrent{torrent1, torrent2})
	}))
	defer srv.Close()

	var (
		expectedMethod = "GET"
		expectedPath   = "/torrents"
		expectedQuery  = "ubuntu 15.10"
	)

	ret := NewClient(srv.URL, testKey).Search(expectedQuery)

	if srvMethod != expectedMethod {
		t.Errorf("Search server method %q, expected %q", srvMethod, expectedMethod)
	}
	if srvPath != expectedPath {
		t.Errorf("Search server path %q, expected %q", srvPath, expectedPath)
	}
	if srvQuery != expectedQuery {
		t.Errorf("Search server query %q, expected %q", srvQuery, expectedQuery)
	}
	if !srvValid {
		t.Error("Search server request was not signed")
	}

	if l := ret.Length(); l != 2 {
		t.Fatalf("Search result length = %d, expected 2", l)
	}

	if trnt := ret.Get(0); *trnt != torrent1 {
		t.Errorf("Search result[0] = %v, expected %v", *trnt, torrent1)
	}
	if trnt := ret.Get(1); *trnt != torrent2 {
		t.Errorf("Search result[0] = %v, expected %v", *trnt, torrent2)
	}
}

func TestClientDownload(t *testing.T) {
	var (
		torrent = Torrent{ID: "someHash"}

		srvMethod, srvPath string
		srvValid           bool

		expectedMethod = "POST"
		expectedPath   = "/download/someHash"
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srvMethod = r.Method
		srvPath = r.URL.Path
		srvValid = requestIsSigned(r, testKey)
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	ret := NewClient(srv.URL, testKey).Download(&torrent)

	if srvMethod != expectedMethod {
		t.Errorf("Download server method %q, expected %q", srvMethod, expectedMethod)
	}
	if srvPath != expectedPath {
		t.Errorf("Download server path %q, expected %q", srvPath, expectedPath)
	}
	if !srvValid {
		t.Error("Download server request was not signed")
	}

	if ret != true {
		t.Errorf("Download result %v, expected %v", ret, true)
	}
}
