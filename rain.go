package main

// RainDrop holds the state of a single raindrop
type RainDrop struct {
	x int
	y int
}

// reset positions the raindrop in a random spot in the sky out of view
func (r *RainDrop) reset() {
	r.x = RandInt(0, term.w-1)
	r.y = RandInt(-term.h, -1)
}

// Fall updates the raindrop after it falls
// If the raindrop has reached the ground, it is reset
func (r *RainDrop) Fall() {
	if r.y+1 < term.h {
		r.y++
	} else {
		r.reset()
		r.y = 0
	}
}

// NewRainDrop creates a new RainDrop
func NewRainDrop() *RainDrop {
	rainDrop := &RainDrop{}
	rainDrop.reset()
	return rainDrop
}
