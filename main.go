package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tncardoso/gocurses"
)

// DropDensity detemines density of raindrops
// Higher density creates more raindrops in total
const DropDensity = 10

// FrameDelay is the delay between refreshes in milliseconds
const FrameDelay = 30

func main() {
	logging.Setup()
	curses.Setup()

	// Handle interruption and termination
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		logging.Teardown()
		curses.Teardown()
		os.Exit(0)
	}()

	// Handle window resize
	d := make(chan os.Signal, 1)
	signal.Notify(d, syscall.SIGWINCH)
	go func() {
		for {
			<-d
			logging.Setup()
			curses.Resize()
		}
	}()

	drops := []*RainDrop{}
	for i := 0; i < DropDensity*term.w; i++ {
		drop := NewRainDrop()
		drops = append(drops, drop)
	}

	for {
		gocurses.Clear()

		for _, drop := range drops {
			drop.Fall()
			gocurses.Mvaddch(drop.y, drop.x, '|')
		}

		gocurses.Refresh()

		time.Sleep(time.Duration(FrameDelay) * time.Millisecond)
	}
}
