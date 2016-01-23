package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gophergala2016/magopie/entities"
)

func mustNewRequest(t *testing.T, method, urlStr string, body io.Reader) *http.Request {
	r, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		t.Fatal(err)
	}

	return r
}

func describeReq(req *http.Request) string {
	return req.Method + " " + req.URL.Path
}

func fooBarSites() []site {
	return []site{
		{
			Site: entities.Site{
				ID:      "foo",
				Name:    "Foo",
				URL:     "http://foo.foo",
				Enabled: true,
			},
		},
		{
			Site: entities.Site{
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
		Site: entities.Site{
			Enabled: true,
		},
	}
	siteA.search = func(term string) []entities.Torrent {
		termA = term
		return []entities.Torrent{
			{
				ID:      "torrentA",
				Title:   "ubuntu 1",
				FileURL: "http://sitea/torrentA",
			},
		}
	}

	var termB string
	siteB := site{
		Site: entities.Site{
			Enabled: true,
		},
	}
	siteB.search = func(term string) []entities.Torrent {
		termB = term
		return []entities.Torrent{
			{
				ID:   "b",
				Site: siteB.Site,
			},
		}
	}

	var termC string
	siteC := site{
		Site: entities.Site{
			Enabled: false,
		},
		search: func(term string) []entities.Torrent {
			termC = term
			return nil
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
			"id":   "a",
			"site": siteA.Site,
		},
		{
			"id":   "b",
			"site": siteB.Site,
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("JSON actual: %v\nExpected: %v", actual, expected)
	}
}
