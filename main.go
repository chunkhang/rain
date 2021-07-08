package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// RainEngine is the engine that controls rain animation
type RainEngine struct {
	delay    int
	density  float64
	stopping chan bool
	stopped  chan bool
}

// NewRainEngine creates a new RainEngine
func NewRainEngine(delay int, density float64) *RainEngine {
	log.Println("create new rain engine...")

	engine := &RainEngine{
		delay:    delay,
		density:  density,
		stopping: make(chan bool),
		stopped:  make(chan bool),
	}

	log.Printf("rain engine = {delay:%d, density:%f}", engine.delay, engine.density)

	return engine
}

// Start starts the rain animation
// A new rain is created before animation begins
func (e *RainEngine) Start() {
	log.Println("rain engine starting...")

	// Create new rain based on screen size
	rain := NewRain(screen.w, screen.h, e.density)

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
			case <-time.After(time.Duration(e.delay) * time.Millisecond):
				continue
			}
		}
	}()
}

// Stop stops the rain animation
// We wait until the animation has stopped completely
func (e *RainEngine) Stop() {
	log.Println("rain engine stopping...")
	e.stopping <- true
	<-e.stopped
	log.Println("rain engine stopped...")
}

// Application flags
var debugMode = flag.Bool("debug", false, "enable debug logs")
var debugFilename = flag.String("logfile", "rain.log", "filename for debug logs to output to")
var rainSpeed = flag.Int("speed", 3, "speed of raindrops [1 - 5]")
var rainDensity = flag.Int("density", 3, "density of raindrops [1 - 5]")

// Parse application flags
func init() {
	flag.Parse()

	if *rainSpeed <= 0 || *rainSpeed > 5 {
		fmt.Println("rain speed must be between 1 - 5")
		os.Exit(1)
	}

	if *rainDensity <= 0 || *rainDensity > 5 {
		fmt.Println("rain density must be between 1 - 5")
		os.Exit(1)
	}
}

func main() {
	logger := &Logger{enabled: *debugMode, filename: *debugFilename}
	logger.Setup()
	defer logger.Teardown()

	// Print flag values
	log.Println("print flags...")
	flag.VisitAll(func(f *flag.Flag) {
		log.Printf("%s = %s\n", f.Name, f.Value)
	})

	screen.Setup()
	defer screen.Teardown()

	// Convert rain speed and rain density from flags
	// We want default delay to be around 20 ms
	// We want default density to be around 0.15
	delay := int(60.0 / float64(*rainSpeed))
	density := float64(*rainDensity) / 20.0
	engine := NewRainEngine(delay, density)
	engine.Start()

	quitting := make(chan bool)

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGWINCH, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			sig := <-s
			log.Println("receive signal...")
			log.Printf("signal = %s\n", sig)
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
