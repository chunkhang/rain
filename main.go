package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tncardoso/gocurses"
)

// Term holds the state of the terminal
type Term struct {
	w int
	h int
}

// Drop holds the state of a single raindrop
type Drop struct {
	x int
	y int
}

// Global variables
var term Term
var win *gocurses.Window
var file os.File

func setup() {
	file, err := os.OpenFile("rain.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	gocurses.Initscr()

	log.SetOutput(file)

	gocurses.Cbreak()
	gocurses.Noecho()
	gocurses.CursSet(0)
	gocurses.Stdscr.Keypad(true)

	resize()
	win = gocurses.NewWindow(term.h, term.w, 0, 0)
}

func teardown() {
	file.Close()
	gocurses.End()
}

func resize() {
	gocurses.End()
	gocurses.Refresh()

	row, col := gocurses.Getmaxyx()

	term.w = col
	term.h = row
}

func main() {
	setup()

	// Handle interruption and termination
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		teardown()
		os.Exit(0)
	}()

	// Handle window resize
	d := make(chan os.Signal, 1)
	signal.Notify(d, syscall.SIGWINCH)
	go func() {
		for {
			<-d
			resize()
		}
	}()

	drop := Drop{}

	for {
		gocurses.Clear()

		gocurses.Mvaddstr(0, 5, fmt.Sprintf("%d, %d", term.w, term.h))

		gocurses.Mvaddch(drop.y, drop.x, '*')

		gocurses.Refresh()

		time.Sleep(1 * time.Second)

		if drop.x >= term.w || drop.y >= term.h {
			teardown()
			os.Exit(0)
		}

		drop.x++
		drop.y++
	}
}
