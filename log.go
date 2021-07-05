package main

import (
	"log"
	"os"
)

// Logging is a module to handle logging
type Logging struct{}

// Filename is the log filename
const Filename = "rain.log"

var logging Logging
var file *os.File

// Setup performs setup for logging
func (l *Logging) Setup() {
	file, err := os.OpenFile(Filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)
}

// Teardown performs teardown for logging
func (l *Logging) Teardown() {
	file.Close()
}
