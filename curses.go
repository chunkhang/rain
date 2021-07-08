package main

import (
	"log"

	"github.com/tncardoso/gocurses"
)

// Curses is a module to handle ncurses
type Curses struct{}

// Term holds the state of the terminal
type Term struct {
	w int
	h int
}

var curses Curses
var term Term

// Setup performs setup for curses
func (c *Curses) Setup() {
	log.Printf("setup curses...")

	gocurses.Initscr()

	gocurses.Cbreak()
	gocurses.Noecho()
	gocurses.CursSet(0)
	// gocurses.Stdscr.Keypad(true)

	term.h, term.w = gocurses.Getmaxyx()
}

// Resize updates the terminal size to the latest
func (c *Curses) Resize() {
	log.Printf("resize curses...")

	gocurses.End()
	gocurses.Refresh()
	term.h, term.w = gocurses.Getmaxyx()
}

// Teardown performs teardown for curses
func (c *Curses) Teardown() {
	log.Printf("teardown curses...")

	gocurses.End()
}
