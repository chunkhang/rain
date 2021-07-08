package main

import (
	"log"
)

// RainDropDensity detemines number of raindrops created for a rain
// For a rain with 50w and 20h, the area is 50 * 20 = 1000
// Density of 0.5 means 0.5 * 1000 = 500 raindrops
const RainDropDensity float32 = 0.1

// RainDrop holds the state of a single raindrop
type RainDrop struct {
	char rune
	x    int
	y    int
}

// Rain holds the state of all raindrops
type Rain struct {
	drops []*RainDrop
	w     int
	h     int
}

// NewRain creates a new rain
func NewRain(w, h int) *Rain {
	log.Printf("create new rain...")

	// Create rain
	rain := &Rain{w: w, h: h}

	// Create raindrops
	area := rain.w * rain.h
	totalDrops := int(RainDropDensity * float32(area))
	drops := make([]*RainDrop, totalDrops)
	for i := range drops {
		drop := &RainDrop{
			char: '|',
			x:    RandInt(0, rain.w-1),
			y:    RandInt(-rain.h, -1),
		}
		drops[i] = drop
	}
	rain.drops = drops

	log.Printf("rain = {drops:[%d], w:%d, h:%d}", len(rain.drops), rain.w, rain.h)

	return rain
}

// Fall updates the state of rainfall by one tick
func (r *Rain) Fall() {
	for _, drop := range r.drops {
		if drop.y+1 < r.h {
			drop.y++
		} else {
			// If a raindrop has reached the ground, reset its position
			drop.x = RandInt(0, r.w-1)
			drop.y = 0
		}
	}
}
