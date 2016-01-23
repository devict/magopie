package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mp "github.com/gophergala2016/magopie"
)

func mustNewRequest(t *testing.T, method, urlStr string, body io.Reader) *http.Request {
	r, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		t.Fatal(err)
	}

	return r
}

func describeReq(req *http.Request) string {
	return req.Method + " " + req.URL.String()
}

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

func TestGetSites(t *testing.T) {
	req := mustNewRequest(t, "GET", "/sites", nil)
	res := httptest.NewRecorder()

	a := &app{
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
			"id":      "foo",
			"name":    "Foo",
			"url":     "http://foo.foo",
			"enabled": true,
		},
		{
			"id":      "bar",
			"name":    "Bar",
			"url":     "http://bar.bar",
			"enabled": false,
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("JSON actual: %v\nExpected: %v", actual, expected)
	}
}

func TestGetSiteSingle(t *testing.T) {
	req := mustNewRequest(t, "GET", "/sites/bar", nil)
	res := httptest.NewRecorder()

	a := &app{
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
		"id":      "bar",
		"name":    "Bar",
		"url":     "http://bar.bar",
		"enabled": false,
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
	}
	siteA.search = func(term string) ([]mp.Torrent, error) {
		termA = term
		return []mp.Torrent{
			{
				ID:      "torrentA",
				Title:   "ubuntu 1",
				FileURL: "http://sitea/torrentA",
			},
		}, nil
	}

	var termB string
	siteB := site{
		Site: mp.Site{
			Enabled: true,
		},
	}
	siteB.search = func(term string) ([]mp.Torrent, error) {
		termB = term
		return []mp.Torrent{
			{
				ID:   "b",
				Site: siteB.Site,
			},
		}, nil
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

	a := &app{
		sites: sitedb{sites: []site{siteA, siteB, siteC}},
	}

	router(a).ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("%s : status = %d, expected %d", describeReq(req), res.Code, http.StatusOK)
	}

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
			"id": "a",
		},
		{
			"id": "b",
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("JSON actual: %v\nExpected: %v", actual, expected)
	}
}

// TestGetTorrentsFail tests what happens when you don't include the
// required param q
func TestGetTorrentsFail(t *testing.T) {
	req := mustNewRequest(t, "GET", "/torrents", nil)
	res := httptest.NewRecorder()

	router(&app{}).ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("%s : status = %d, expected %d", describeReq(req), res.Code, http.StatusBadRequest)
	}
}
