package main

import (
	"math/rand"
	"strconv"
	"strings"
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

func FXrampv(part part) row {
	vmap := make(row, len(part.Lanes[0]))
	if len(vmap) == 0 {
		debugf("FXrampv(): empty part: %v", part)
		return nil
	}
	for _, ef := range part.Fx {
		if v, ok := ef["rampv"]; ok {
			debugf("FXrampv(): found rampv value %v", v)
			values := strings.Split(v, "-")
			if len(values) != 2 {
				return nil
			}
			start, err := strconv.Atoi(values[0])
			if err != nil {
				return nil
			}
			end, err := strconv.Atoi(values[1])
			if err != nil {
				return nil
			}
			step := (end - start) / len(vmap)
			for i := range vmap {
				vmap[i] = start + (i * step)
			}
			return vmap
		}
	}
	return nil
}

func genVelocityMap(part part) matrix {
	// if available, use ramp as base for vmap
	// otherwise initialise one with max velocity
	var vmatrix matrix
	vmap := FXrampv(part)
	if vmap == nil {
		vmap = make(row, len(part.Lanes[0]))
		for i := range vmap {
			vmap[i] = 127
		}
	}
	randomness := FXrandv(part)
	for i := range vmap {
		if randomness != 0 {
			v := (rand.Int() % randomness) - (randomness / 2)
			v = vmap[i] - (randomness / 2) + v
			if v < 0 {
				v = 0
			}
			vmap[i] = v
		}
	}
	for range part.Lanes {
		vmatrix = append(vmatrix, vmap)
	}
	debugf("genVelocityMap(): vmatrix: %v", vmatrix)
	return vmatrix
}

func playChord(e event) {
	debugf("playChord(): %v", e)
	for i, note := range e.Notes {
		sendMidiNote(5, note, e.Velocities[i])
	}
}
