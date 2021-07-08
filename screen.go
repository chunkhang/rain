package main

import (
	"log"

	"github.com/tncardoso/gocurses"
)

// Screen represents the application screen
type Screen struct {
	w int
	h int
}

var screen Screen

// Setup performs setup for screen
func (s *Screen) Setup() {
	log.Printf("setup screen...")

	gocurses.Initscr()

	gocurses.Cbreak()
	gocurses.Noecho()
	gocurses.CursSet(0)

	s.h, s.w = gocurses.Getmaxyx()

	// Endlessly wait for input, so keypresses do not affect the screen
	go func() {
		for {
			gocurses.Stdscr.Getch()
		}
	}()
}

// AddChar adds a character to the screen buffer
func (s *Screen) AddChar(char rune, x, y int) {
	gocurses.Mvaddch(y, x, char)
}

// Clear clears the screen buffer
func (s *Screen) Clear() {
	gocurses.Clear()
}

// Flush flushes the screen buffer, writing everything to the screen
func (s *Screen) Flush() {
	gocurses.Refresh()
}

// Resize updates the screen size to the latest
func (s *Screen) Resize() {
	log.Printf("resize screen...")

	gocurses.End()
	gocurses.Refresh()

	s.h, s.w = gocurses.Getmaxyx()
}

// Teardown performs teardown for screen
func (s *Screen) Teardown() {
	log.Printf("teardown screen...")

	gocurses.End()
}
