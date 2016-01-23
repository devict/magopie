package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
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

	router().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("%s : status = %d, expected %d", describeReq(req), res.Code, http.StatusOK)
	}

	body := []map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal("Error parsing response as json. Err:", err)
	}

	// TODO assert something about the response
}
