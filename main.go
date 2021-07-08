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

// RainEngine is the engine that controls rain animation
type RainEngine struct {
	stopping chan bool
	stopped  chan bool
}

// NewRainEngine creates a new RainEngine
func NewRainEngine() *RainEngine {
	engine := &RainEngine{
		stopping: make(chan bool),
		stopped:  make(chan bool),
	}
	return engine
}

// Start starts the rain animation
// A new rain is created before animation begins
func (e *RainEngine) Start() {
	log.Printf("rain engine starting...")

	// Create new rain based on terminal size
	rain := NewRain(term.w, term.h)

	go func() {
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
			case <-e.stopping:
				gocurses.Clear()
				gocurses.Refresh()
				e.stopped <- true
				return
			// Otherwise, wait a bit before the next tick
			case <-time.After(time.Duration(RainFallDelay) * time.Millisecond):
				continue
			}
		}
	}()
}

// Stop stops the rain animation
// We wait until the animation has stopped completely
func (e *RainEngine) Stop() {
	log.Printf("rain engine stopping...")
	e.stopping <- true
	<-e.stopped
	log.Printf("rain engine stopped...")
}

func main() {
	logging.Setup()
	defer logging.Teardown()

	curses.Setup()
	defer curses.Teardown()

	engine := NewRainEngine()
	engine.Start()

	quitting := make(chan bool)

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
				engine.Stop()
				curses.Resize()
				engine.Start()
			// Handle interruption / termination
			case syscall.SIGINT, syscall.SIGTERM:
				engine.Stop()
				quitting <- true
			}
		}
	}()

	<-quitting
}
