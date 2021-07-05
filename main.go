package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logging.Setup()
	curses.Setup()

	rain := &Rain{}
	rain.Setup()
	go rain.Start()

	quit := make(chan bool)

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGWINCH, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			sig := <-s
			if sig == syscall.SIGWINCH {
				// Handle terminal resize
				rain.Stop()
				curses.Resize()
			} else {
				// Handle interruption / termination
				quit <- true
				return
			}
		}
	}()

	<-quit
	rain.Stop()
	logging.Teardown()
	curses.Teardown()
}
