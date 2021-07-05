package main

import (
	"github.com/tncardoso/gocurses"
)

// Curses is a module to handle ncurses
type Curses struct{}

var curses Curses

// Setup performs setup for curses
func (c *Curses) Setup() {
	gocurses.Initscr()

	gocurses.Cbreak()
	gocurses.Noecho()
	gocurses.CursSet(0)
	// gocurses.Stdscr.Keypad(true)

	curses.Resize()
}

// Resize updates the terminal size to the latest
func (c *Curses) Resize() {
	gocurses.End()
	gocurses.Refresh()
	row, col := gocurses.Getmaxyx()

	term.w = col
	term.h = row
}

// Teardown performs teardown for curses
func (c *Curses) Teardown() {
	gocurses.End()
}
