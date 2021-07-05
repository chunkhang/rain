package main

import (
	"time"

	"github.com/tncardoso/gocurses"
)

// DropDensity detemines density of raindrops
// Higher density creates more raindrops in total
const DropDensity = 10

// RainFallDelay is the delay between rain falls in milliseconds
const RainFallDelay = 20

// Rain holds the state of all raindrops
type Rain struct {
	drops []*RainDrop
	stop  bool
}

// Setup performs the setup for rainfall
func (r *Rain) Setup() {
	r.drops = make([]*RainDrop, DropDensity*term.w)
	for i := range r.drops {
		drop := &RainDrop{char: '|'}
		drop.Reset()
		r.drops[i] = drop
	}
}

// Start starts the rainfall
func (r *Rain) Start() {
	for {
		if r.stop {
			return
		}
		gocurses.Clear()
		for _, drop := range r.drops {
			drop.Fall()
			gocurses.Mvaddch(drop.y, drop.x, drop.char)
		}
		gocurses.Refresh()
		time.Sleep(time.Duration(RainFallDelay) * time.Millisecond)
	}
}

// Stop stops the rainfall
func (r *Rain) Stop() {
	r.stop = true
}

// RainDrop holds the state of a single raindrop
type RainDrop struct {
	x    int
	y    int
	char rune
}

// Reset positions the raindrop in a random spot in the sky out of view
func (r *RainDrop) Reset() {
	r.x = RandInt(0, term.w-1)
	r.y = RandInt(-term.h, -1)
}

// Fall updates the raindrop after it falls
// If the raindrop has reached the ground, it is reset
func (r *RainDrop) Fall() {
	if r.y+1 < term.h {
		r.y++
	} else {
		r.Reset()
		r.y = 0
	}
}
