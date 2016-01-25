package main

import (
	"io"
	"net/http"
	"testing"

	"github.com/devict/magopie"
	"github.com/pborman/uuid"
)

var testKey = "Margo the magpie"

func mustNewRequest(t *testing.T, method, urlStr string, body io.Reader) *http.Request {
	r, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		t.Fatal(err)
	}

	id := uuid.NewRandom().String()
	r.Header.Set("X-Request-ID", id)
	r.Header.Set("X-HMAC", magopie.HMAC(id, testKey))

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
