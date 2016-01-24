package main

import "errors"

var (
	// ErrFailedRequest is returned when we get a non-200 status from one of our
	// upstream sites.
	ErrFailedRequest = errors.New("http request had non-200 status")

	// ErrHashlessMagnet is returned when we fail to parse the BTIH hash from a
	// Magnet URI
	ErrHashlessMagnet = errors.New("could not find hash in magnet URI")

	// ErrCannotFindFileSize is used when we fail to parse the file size for a
	// torrent returned from TPB.
	ErrCannotFindFileSize = errors.New("could not find tpb filesize in html")
)
