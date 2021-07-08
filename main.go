package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tncardoso/gocurses"
)

// RainFallDelay is the delay between rainfall ticks in milliseconds
const RainFallDelay = 20

func animate(rain *Rain, stop chan bool) {
	for {
		// Update state of rain
		rain.Fall()

		// Draw the new state of rain
		gocurses.Clear()
		for _, drop := range rain.drops {
			gocurses.Mvaddch(drop.y, drop.x, drop.char)
		}
		gocurses.Refresh()

		select {
		// Time to stop, so clear the screen and we are done
		case <-stop:
			gocurses.Clear()
			gocurses.Refresh()
			return
		// Otherwise, wait a bit before the next tick
		case <-time.After(time.Duration(RainFallDelay) * time.Millisecond):
			continue
		}
	}
}

func main() {
	logging.Setup()
	defer logging.Teardown()

	curses.Setup()
	defer curses.Teardown()

	stop := make(chan bool)
	quit := make(chan bool)

	rain := NewRain(term.w, term.h)
	go animate(rain, stop)

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGWINCH, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			sig := <-s
			if sig == syscall.SIGWINCH {
				// Handle terminal resize
				stop <- true
				close(stop)
				// rain.Stop()
				// curses.Resize()
				// rain.Setup()
				// go rain.Start()
			} else {
				// Handle interruption / termination
				stop <- true
				close(stop)
				quit <- true
				return
			}
		}
	}()

	<-quit
}
