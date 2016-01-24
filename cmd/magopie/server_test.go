package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	mp "github.com/gophergala2016/magopie"
	"github.com/spf13/afero"
)

func fooBarSites() []site {
	return []site{
		{
			Site: mp.Site{
				ID:      "foo",
				Name:    "Foo",
				URL:     "http://foo.foo",
				Enabled: true,
			},
		},
		{
			Site: mp.Site{
				ID:      "bar",
				Name:    "Bar",
				URL:     "http://bar.bar",
				Enabled: false,
			},
		},
	}
}

// TestAuthentication checks that every route we register gets denied with a
// 401 if it is not signed.
func TestAuthentication(t *testing.T) {
	tests := []struct {
		method, path string
	}{
		{"GET", "/sites"},
		{"GET", "/sites/foo"},
		{"GET", "/torrents?q=bar"},
		{"POST", "/download/1234567890"},
	}

	r := router(&server{key: testKey})

	for i, test := range tests {
		req, err := http.NewRequest(test.method, test.path, nil)
		res := httptest.NewRecorder()
		if err != nil {
			t.Fatal("Error making request")
		}
		r.ServeHTTP(res, req)
		if res.Code != http.StatusUnauthorized {
			t.Errorf("Test %d: %s %s gave status %d, expected %d", i, test.method, test.path, res.Code, http.StatusUnauthorized)
		}
	}
}

func TestGetSites(t *testing.T) {
	req := mustNewRequest(t, "GET", "/sites", nil)
	res := httptest.NewRecorder()

	a := &server{
		key:   testKey,
		sites: sitedb{sites: fooBarSites()},
	}

	router(a).ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("%s : status = %d, expected %d", describeReq(req), res.Code, http.StatusOK)
	}

	actual := []map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&actual); err != nil {
		t.Fatal("Error parsing json response as []. Err:", err)
	}

	expected := []map[string]interface{}{
		{
			"ID":      "foo",
			"Name":    "Foo",
			"URL":     "http://foo.foo",
			"Enabled": true,
		},
		{
			"ID":      "bar",
			"Name":    "Bar",
			"URL":     "http://bar.bar",
			"Enabled": false,
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("JSON actual: %v\nExpected: %v", actual, expected)
	}
}

func TestGetSiteSingle(t *testing.T) {
	req := mustNewRequest(t, "GET", "/sites/bar", nil)
	res := httptest.NewRecorder()

	a := &server{
		key:   testKey,
		sites: sitedb{sites: fooBarSites()},
	}

	router(a).ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("%s : status = %d, expected %d", describeReq(req), res.Code, http.StatusOK)
	}

	actual := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&actual); err != nil {
		t.Fatal("Error parsing json response as []. Err:", err)
	}

	expected := map[string]interface{}{
		"ID":      "bar",
		"Name":    "Bar",
		"URL":     "http://bar.bar",
		"Enabled": false,
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("JSON actual: %v\nExpected: %v", actual, expected)
	}
}

func TestGetTorrents(t *testing.T) {
	term := "ubuntu 15.10"
	req := mustNewRequest(t, "GET", "/torrents?q=ubuntu%2015.10", nil)
	res := httptest.NewRecorder()

	var termA string
	siteA := site{
		Site: mp.Site{
			Enabled: true,
		},
		search: func(term string) ([]mp.Torrent, error) {
			termA = term
			return []mp.Torrent{
				{
					ID:        "torrent1",
					Title:     "ubuntu 1",
					MagnetURI: "http://sitea/torrent1",
					SiteID:    "a",
					Seeders:   10,
					Leechers:  50,
					Size:      1234567,
				},
			}, nil
		},
	}

	var termB string
	siteB := site{
		Site: mp.Site{
			Enabled: true,
		},
		search: func(term string) ([]mp.Torrent, error) {
			termB = term
			return []mp.Torrent{
				{
					ID:        "torrentB",
					Title:     "ubuntu 2",
					MagnetURI: "http://siteb/torrentB",
					SiteID:    "b",
					Seeders:   20,
					Leechers:  80,
					Size:      7654321,
				},
			}, nil
		},
	}

	var termC string
	siteC := site{
		Site: mp.Site{
			Enabled: false,
		},
		search: func(term string) ([]mp.Torrent, error) {
			termC = term
			return nil, nil
		},
	}

	a := &server{
		key:   testKey,
		sites: sitedb{sites: []site{siteA, siteB, siteC}},
	}

	router(a).ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("%s : status = %d, expected %d", describeReq(req), res.Code, http.StatusOK)
	}

	// Ensure our search term was passed to the appropriate site search funcs
	if termA != term {
		t.Errorf("Site A search was passed %q, expected %q", termA, term)
	}
	if termB != term {
		t.Errorf("Site B search was passed %q, expected %q", termB, term)
	}
	if termC != "" {
		t.Errorf("Site C search was passed %q, expected site to not be used", termC)
	}

	actual := []map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&actual); err != nil {
		t.Fatal("Error parsing json response as []. Err:", err)
	}

	expected := []map[string]interface{}{
		{
			"ID":        "torrentB",
			"Title":     "ubuntu 2",
			"MagnetURI": "http://siteb/torrentB",
			"SiteID":    "b",
			"Seeders":   float64(20),
			"Leechers":  float64(80),
			"Size":      float64(7654321),
		},
		{
			"ID":        "torrent1",
			"Title":     "ubuntu 1",
			"MagnetURI": "http://sitea/torrent1",
			"SiteID":    "a",
			"Seeders":   float64(10),
			"Leechers":  float64(50),
			"Size":      float64(1234567),
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		diffMaps(t, actual, expected)
		t.Errorf("JSON actual: %v\nExpected: %v", actual, expected)
	}
}

// TestGetTorrentsFail tests what happens when you don't include the
// required param q
func TestGetTorrentsFail(t *testing.T) {
	req := mustNewRequest(t, "GET", "/torrents", nil)
	res := httptest.NewRecorder()

	router(&server{key: testKey}).ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("%s : status = %d, expected %d", describeReq(req), res.Code, http.StatusBadRequest)
	}
}

func TestPostDownload(t *testing.T) {
	a := &server{
		key:         testKey,
		fs:          &afero.MemMapFs{},
		downloadDir: "/magopie/downloads",
	}

	hash := "337b6dbb824ff8acf38846d4698746df7bf2b5c9"
	file := strings.ToUpper(hash) + ".torrent"
	fullFile := filepath.Join(a.downloadDir, file)
	content := "torrents!"

	var torMethod, torPath string
	torcacheSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		torMethod = r.Method
		torPath = r.URL.Path
		content := strings.NewReader(content)
		http.ServeContent(w, r, file, time.Time{}, content)
	}))
	defer torcacheSrv.Close()
	a.torcacheURL = torcacheSrv.URL

	// Make the request to magopie to download a particular torrent by hash
	req := mustNewRequest(t, "POST", "/download/"+hash, nil)
	res := httptest.NewRecorder()

	router(a).ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Errorf("%s : status = %d, expected %d", describeReq(req), res.Code, http.StatusCreated)
	}

	// Ensure things about the post to the torache server
	if torMethod != "GET" {
		t.Errorf("torcache server method = %q, expected %q", torMethod, "GET")
	}
	expectedPath := "/torrent/337B6DBB824FF8ACF38846D4698746DF7BF2B5C9.torrent"
	if torPath != expectedPath {
		t.Errorf("torcache server path = %q, expected %q", torPath, expectedPath)
	}

	// Ensure we downloaded the file from the torcache server
	actualContent, err := afero.ReadFile(a.fs, fullFile)
	if err != nil {
		t.Errorf("err reading file in test: %v", err)
	}

	if string(actualContent) != content {
		t.Errorf("downloaded file contents %s\nexpected %s", actualContent, content)
	}
}
