package main

import (
	"io/ioutil"
	"log"
	"os"
)

// Logger represents the application logger
type Logger struct {
	enabled  bool
	filename string
	file     *os.File
}

// Setup performs setup for logger
// Debug mode determines whether we are logging to a file
func (l *Logger) Setup() {
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

// Teardown performs teardown for logger
func (l *Logger) Teardown() {
	if l.enabled {
		l.file.Close()
	}
}
