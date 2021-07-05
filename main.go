package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

// RainFallDelay is the delay between rain falls in milliseconds
const RainFallDelay = 20

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
			curses.Resize()
		}
	}()

	rain := NewRain()

	for {
		rain.Fall()
		time.Sleep(time.Duration(RainFallDelay) * time.Millisecond)
	}
}
