package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var debugMode = flag.Bool("debug", false, "enable debug logs")
var debugFilename = flag.String("logfile", "rain.log", "filename for debug logs to output to")
var rainfallDelay = flag.Int("delay", 20, "delay between rainfall ticks")
var rainDensity = flag.Float64("density", 0.2, "density of raindrops")

// RainEngine is the engine that controls rain animation
type RainEngine struct {
	tickDelay int
	stopping  chan bool
	stopped   chan bool
}

// NewRainEngine creates a new RainEngine
func NewRainEngine(tickDelay int) *RainEngine {
	engine := &RainEngine{
		tickDelay: tickDelay,
		stopping:  make(chan bool),
		stopped:   make(chan bool),
	}
	return engine
}

// Start starts the rain animation
// A new rain is created before animation begins
func (e *RainEngine) Start() {
	log.Printf("rain engine starting...")

	// Create new rain based on screen size
	rain := NewRain(screen.w, screen.h, *rainDensity)

	go func() {
		for {
			// Draw current state of rain
			screen.Clear()
			for _, drop := range rain.drops {
				screen.AddChar(drop.char, drop.x, drop.y)
			}
			screen.Flush()

			// Update state of rain
			rain.Fall()

			select {
			// Time to stop, so clear the screen and we are done
			case <-e.stopping:
				screen.Clear()
				screen.Flush()
				e.stopped <- true
				return
			// Otherwise, wait a bit before the next tick
			case <-time.After(time.Duration(e.tickDelay) * time.Millisecond):
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

func init() {
	flag.Parse()
}

func main() {
	logger := &Logger{enabled: *debugMode, filename: *debugFilename}
	logger.Setup()
	defer logger.Teardown()

	// Print flag values
	log.Printf("print flags...")
	flag.VisitAll(func(f *flag.Flag) {
		log.Printf("%s = %s", f.Name, f.Value)
	})

	screen.Setup()
	defer screen.Teardown()

	engine := NewRainEngine(*rainfallDelay)
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
				screen.Resize()
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
