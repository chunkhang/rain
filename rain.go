package main

import (
	"log"
)

// RainDrop holds the state of a single raindrop
type RainDrop struct {
	char rune
	x    int
	y    int
}

// Rain holds the state of all raindrops
// Density detemines number of raindrops created for a rain
// For a rain with 50w and 20h, the area is 50 * 20 = 1000
// Density of 0.5 means 0.5 * 1000 = 500 raindrops
type Rain struct {
	w       int
	h       int
	density float64
	drops   []*RainDrop
}

// NewRain creates a new rain
func NewRain(w, h int, density float64) *Rain {
	log.Println("create new rain...")

	// Create rain
	rain := &Rain{w: w, h: h, density: density}

	// Create raindrops
	area := rain.w * rain.h
	totalDrops := int(rain.density * float64(area))
	drops := make([]*RainDrop, totalDrops)
	for i := range drops {
		// We want more heavy drops than light drops
		// To create the illusion of depth
		var char rune
		switch RandInt(0, 5) {
		case 0, 1, 2:
			char = '|'
		case 3, 4:
			char = ':'
		case 5:
			char = '.'
		}
		drop := &RainDrop{
			char: char,
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
