package magopie

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pborman/uuid"
)

// SaveToFile saves 'data' to the file 'filePath'
func SaveToFile(filePath string, data []byte) error {
	return ioutil.WriteFile(filePath, data, 0644)
}

// ReadFromFile returns the bytes from 'filePath'
func ReadFromFile(filePath string) []byte {
	var data []byte
	var err error
	if data, err = ioutil.ReadFile(filePath); err != nil {
		fmt.Println(err)
	}
	return data
}

// HMAC gives the sha256 HMAC of message using key. It is expressed as a hex
// string.
func HMAC(message, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))

	return hex.EncodeToString(mac.Sum(nil))
}

// CheckMAC confirms that the mac provided with a message was signed with the
// appropriate key. The mac should be in a hex string.
func CheckMAC(message, messageMAC, key string) bool {
	msgMAC, err := hex.DecodeString(messageMAC)
	if err != nil {
		log.Print("message mac was not valid hex", err)
		return false
	}

	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))

	return hmac.Equal(msgMAC, mac.Sum(nil))
}

// SignRequest sets X-Request-ID and X-HMAC headers on a request signed with
// the given key.
func SignRequest(r *http.Request, key string) {
	id := uuid.NewRandom().String()
	r.Header.Set("X-Request-ID", id)
	r.Header.Set("X-HMAC", HMAC(id, key))
}

// RequestIsSigned validates the authentication headers on a request.
func RequestIsSigned(r *http.Request, key string) bool {
	hdr := r.Header
	return CheckMAC(hdr.Get("X-Request-ID"), hdr.Get("X-HMAC"), key)
}
