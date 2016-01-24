package main

import (
	"io"
	"net/http"
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
	return req.Method + " " + req.URL.String()
}

func diffMaps(t *testing.T, a, b []map[string]interface{}) {
	for i := range a {
		for k := range a[i] {
			valA, valB := a[i][k], b[i][k]
			if valA != valB {
				t.Logf("%d %q: a: %v, b: %v", i, k, valA, valB)
			}
		}
	}
}
