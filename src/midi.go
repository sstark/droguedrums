package main

import (
	"math/rand"
	"strconv"
)

type event struct {
	Notes      row
	Velocities row
}

func FXrandv(part part) int {
	for _, ef := range part.Fx {
		if v, ok := ef["randv"]; ok {
			debugf("FXrandv(): found randv value %v", v)
			randomness, err := strconv.Atoi(v)
			if err == nil {
				return randomness
			}
		}
	}
	return 0
}

func genVelocityMap(notes row, part part) row {
	randomness := FXrandv(part)
	vmap := make(row, len(notes))
	for i := range notes {
		if randomness == 0 {
			vmap[i] = 127
		} else {
			v := (rand.Int() % randomness) - (randomness / 2)
			v = 127 - (randomness / 2) + v
			vmap[i] = v
		}
	}
	return vmap
}

func playChord(e event) {
	debugf("playChord(): %v", e)
	for i, note := range e.Notes {
		sendMidiNote(5, note, e.Velocities[i])
	}
}
