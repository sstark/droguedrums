package main

import (
	"math/rand"
	"strconv"
	"strings"
)

func fxRandv(part part) int {
	for _, ef := range part.Fx {
		if v, ok := ef["randv"]; ok {
			debugf("fxRandv(): found randv value %v", v)
			randomness, err := strconv.Atoi(v)
			if err == nil {
				return randomness
			}
		}
	}
	return 0
}

func fxRampv(part part, notes matrix) row {
	vmap := make(row, len(notes[0]))
	if len(vmap) == 0 {
		debugf("fxRampv(): empty part: %v", notes)
		return nil
	}
	for _, ef := range part.Fx {
		if v, ok := ef["rampv"]; ok {
			debugf("fxRampv(): found rampv value %v", v)
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

func genVelocityMap(part part, notes matrix) matrix {
	// if available, use ramp as base for vmap
	// otherwise initialise one with max velocity
	var vmatrix matrix
	debugf("genVelocityMap(): notes: %v", notes)
	vmap := fxRampv(part, notes)
	debugf("genVelocityMap(): part length: %v", len(vmap))
	if vmap == nil {
		vmap = make(row, len(notes[0]))
		for i := range vmap {
			vmap[i] = midiVmax
		}
	}
	randomness := fxRandv(part)
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
	for range notes {
		vmatrix = append(vmatrix, vmap)
	}
	debugf("genVelocityMap(): vmatrix: %v", vmatrix)
	return vmatrix
}
