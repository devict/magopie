// Package client is for communicating with a Magopie server
package client

import "fmt"

// Greetings ...
func Greetings(n string) string {
	return fmt.Sprintf("Hello, %s!", n)
}
