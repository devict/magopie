package magopie

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearch(t *testing.T) {
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
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srvMethod = r.Method
		srvPath = r.URL.Path
		srvQuery = r.FormValue("q")
		json.NewEncoder(w).Encode([]Torrent{torrent1, torrent2})
	}))
	defer srv.Close()

	var (
		expectedMethod = "GET"
		expectedPath   = "/torrents"
		expectedQuery  = "ubuntu 15.10"
	)

	ret := NewClient(srv.URL).Search(expectedQuery)

	if srvMethod != expectedMethod {
		t.Errorf("Search server method %q, expected %q", srvMethod, expectedMethod)
	}
	if srvPath != expectedPath {
		t.Errorf("Search server path %q, expected %q", srvPath, expectedPath)
	}
	if srvQuery != expectedQuery {
		t.Errorf("Search server query %q, expected %q", srvQuery, expectedQuery)
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
