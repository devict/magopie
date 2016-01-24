package magopie

import (
	"net/http"
	"testing"
)

func TestHMAC(t *testing.T) {
	var (
		message = "hello"
		key     = "world"
		mac     = "3cfa76ef14937c1c0ea519f8fc057a80fcd04a7420f8e8bcd0a7567c272e007b"
		badMAC  = "3cfa76ef14937c1c0ea519f8fc057a80fcd04a7420f8e8bcd0a7567c272e007c"
	)

	actual := HMAC(message, key)

	if actual != mac {
		t.Errorf("HMAC(%q, %q) = %q, expected %q", message, key, actual, mac)
	}
	if same := CheckMAC(message, mac, key); !same {
		t.Error("CheckMAC failed: returned false for known good MAC")
	}
	if same := CheckMAC(message, badMAC, key); same {
		t.Error("CheckMAC should have failed: returned true for bad MAC")
	}
}

func TestSignRequest(t *testing.T) {
	var (
		req, _ = http.NewRequest("GET", "/foo", nil)
		key    = "banana"
	)

	signRequest(req, key)

	id := req.Header.Get("X-Request-ID")
	if id == "" {
		t.Fatal("Request header X-Request-ID was blank")
	}

	hmac := req.Header.Get("X-HMAC")
	if hmac == "" {
		t.Fatal("Request header X-HMAC was blank")
	}

	if !CheckMAC(id, hmac, key) {
		t.Fatal("Request X-HMAC was invalid")
	}
}

func TestRequestIsSigned(t *testing.T) {
	var (
		req, _ = http.NewRequest("GET", "/foo", nil)
		id     = "hello"
		key    = "world"
		mac    = "3cfa76ef14937c1c0ea519f8fc057a80fcd04a7420f8e8bcd0a7567c272e007b"
		badMAC = "4cfa76ef14937c1c0ea519f8fc057a80fcd04a7420f8e8bcd0a7567c272e007b"
	)

	req.Header.Set("X-Request-ID", id)
	req.Header.Set("X-HMAC", mac)

	if !requestIsSigned(req, key) {
		t.Fatal("RequestIsSigned rejected our message, should have accepted it.")
	}

	req.Header.Set("X-HMAC", badMAC)
	if requestIsSigned(req, key) {
		t.Fatal("RequestIsSigned accepted our message, should have rejected it.")
	}
}
