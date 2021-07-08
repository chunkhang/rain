package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tncardoso/gocurses"
)

// RainFallDelay is the delay between rainfall ticks in milliseconds
const RainFallDelay = 20

func startRaining(stopping, stopped chan bool) {
	log.Printf("animate rain starting...")

	// Create new rain
	rain := NewRain(term.w, term.h)

	for {
		// Draw current state of rain
		gocurses.Clear()
		for _, drop := range rain.drops {
			gocurses.Mvaddch(drop.y, drop.x, drop.char)
		}
		gocurses.Refresh()

		// Update state of rain
		rain.Fall()

		select {
		// Time to stop, so clear the screen and we are done
		case <-stopping:
			log.Printf("animate rain stopping...")
			gocurses.Clear()
			gocurses.Refresh()
			stopped <- true
			log.Printf("animate rain stopped...")
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

	stopping := make(chan bool)
	stopped := make(chan bool)
	quitting := make(chan bool)

	go startRaining(stopping, stopped)

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGWINCH, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			sig := <-s
			log.Println("receive signal...")
			log.Printf("signal = %s", sig)
			switch sig {
			// Handle terminal resize
			case syscall.SIGWINCH:
				stopping <- true
				<-stopped
				curses.Resize()
				go startRaining(stopping, stopped)
			// Handle interruption / termination
			case syscall.SIGINT, syscall.SIGTERM:
				stopping <- true
				<-stopped
				quitting <- true
			}
		}
	}()

	<-quitting
}
