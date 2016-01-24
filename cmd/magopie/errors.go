package main

import "errors"

// ErrFailedRequest is returned when we get a non-200 status from one of our
// upstream sites.
var ErrFailedRequest = errors.New("http request had non-200 status")
var ErrHashlessMagnet = errors.New("could not find hash in magnet URI")
var ErrCannotFindFileSize = errors.New("could not find tpb filesize in html")
