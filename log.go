package main

import (
	"io/ioutil"
	"log"
	"os"
)

// Logging is a module to handle logging
type Logging struct {
	enabled  bool
	filename string
	file     *os.File
}

// Setup performs setup for logging
// Debug mode determines whether we are logging to a file
func (l *Logging) Setup() {
	if l.enabled {
		file, err := os.OpenFile(l.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		l.file = file
		log.SetOutput(l.file)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

// Teardown performs teardown for logging
func (l *Logging) Teardown() {
	if l.enabled {
		l.file.Close()
	}
}
