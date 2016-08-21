package main

import (
	"math/rand"
)

func playChord(c row) {
	debugf("playChord(): %v", c)
	dev := 60
	for _, i := range c {
		v := (rand.Int() % dev) - (dev / 2)
		v = 127 - (dev / 2) + v
		sendMidiNote(5, i, v)
	}
}
