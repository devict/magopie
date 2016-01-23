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

func TestGetSites(t *testing.T) {
	req := mustNewRequest(t, "GET", "/sites", nil)
	res := httptest.NewRecorder()

	a := &app{
		sites: sitedb{
			sites: []site{
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
			},
		},
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
