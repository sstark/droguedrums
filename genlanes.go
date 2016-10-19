package main

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func sineFunc(period, xshift, yshift float64) func(float64) float64 {
	return func(x float64) float64 {
		return math.Sin((x*period)-xshift) + yshift
	}
}

func genPlace(gl map[string]string) (out string, err error) {
	//{note: hc, pos: 1 3 5 8}
	// -> hc -- hc -- hc -- -- hc
	var buffer bytes.Buffer
	inpNote, ok := gl["note"]
	if !ok {
		err = errors.New("note value missing")
		return
	}
	inpPos, ok := gl["pos"]
	if !ok {
		inpPos = "1"
	}
	positions := strings.Split(inpPos, " ")
	if len(positions) == 0 {
		return
	}
	lastI := 0
	for _, pos := range positions {
		posI, err := strconv.Atoi(pos)
		// just ignore non-integers
		if err != nil {
			continue
		}
		for ii := lastI; ii < posI-1; ii++ {
			buffer.WriteString("-- ")
		}
		lastI = posI
		buffer.WriteString(inpNote)
		buffer.WriteString(" ")
	}
	out = strings.TrimSpace(buffer.String())
	debugf("genPlace(): %v", out)
	return
}

func genSinez(gl map[string]string) (out string, err error) {
	//{note: hc, length: 13, period: 1.0, xshift: 0.4, yshift: -0.37}
	var buffer bytes.Buffer
	inpNote, ok := gl["note"]
	if !ok {
		err = errors.New("note value missing")
		return
	}
	inpLength, ok := gl["length"]
	if !ok {
		err = errors.New("length value missing")
		return
	}
	inpPeriod, ok := gl["period"]
	if !ok {
		inpPeriod = "1.0"
	}
	inpXshift, ok := gl["xshift"]
	if !ok {
		inpXshift = "0.0"
	}
	inpYshift, ok := gl["yshift"]
	if !ok {
		inpYshift = "0.0"
	}
	length, err := strconv.Atoi(inpLength)
	if err != nil {
		return
	}
	period, err := strconv.ParseFloat(inpPeriod, 32)
	if err != nil {
		return
	}
	xshift, err := strconv.ParseFloat(inpXshift, 32)
	if err != nil {
		return
	}
	yshift, err := strconv.ParseFloat(inpYshift, 32)
	if err != nil {
		return
	}
	step := (2 * math.Pi) / float64(length)
	debugf("genSinez(): step: %f, xshift: %f, yshift: %f", step, xshift, yshift)
	sine := sineFunc(1/period, ((xshift)*step)/period, yshift)
	// the i-1th step
	lastsign := math.Signbit(sine(-step))
	for i := 1; i <= length; i++ {
		x := float64(i) * step
		// find zero crossings
		if math.Signbit(sine(x)) != lastsign {
			buffer.WriteString(inpNote)
			lastsign = !lastsign
		} else {
			buffer.WriteString("--")
		}
		buffer.WriteString(" ")
	}
	out = strings.TrimSpace(buffer.String())
	debugf("genSinez(): %v", out)
	return
}

func genEquid(gl map[string]string) (out string, err error) {
	//{note: hc, length: 13, dist: 2, start: 1}
	var buffer bytes.Buffer
	inpNote, ok := gl["note"]
	if !ok {
		err = errors.New("note value missing")
		return
	}
	inpLength, ok := gl["length"]
	if !ok {
		err = errors.New("length value missing")
		return
	}
	inpDist, ok := gl["dist"]
	if !ok {
		inpDist = "1"
	}
	inpStart, ok := gl["start"]
	if !ok {
		inpStart = "1"
	}
	length, err := strconv.Atoi(inpLength)
	if err != nil {
		return
	}
	dist, err := strconv.Atoi(inpDist)
	if err != nil {
		return
	}
	start, err := strconv.Atoi(inpStart)
	if err != nil {
		return
	}
	for i := 0; i < start-1; i++ {
		buffer.WriteString("-- ")
	}
	for i := 0; i < length-start+1; i++ {
		if i%dist == 0 {
			buffer.WriteString(inpNote)
		} else {
			buffer.WriteString("--")
		}
		buffer.WriteString(" ")
	}
	out = strings.TrimSpace(buffer.String())
	debugf("genEquid(): %v", out)
	return
}

func renderGenlanes(lanes []map[string]map[string]string) (genlanes []string, outerr error) {
	for i, inLane := range lanes {
		if gen, ok := inLane["equid"]; ok {
			debugf("renderGenlanes(): found equid gen %v", gen)
			outLane, err := genEquid(gen)
			if err != nil {
				debugf("renderGenlanes(): gen_equid() failed")
				outerr = fmt.Errorf("error in genlane#%d: %s", i, err)
				return
			}
			genlanes = append(genlanes, outLane)
		}
		if gen, ok := inLane["sinez"]; ok {
			debugf("renderGenlanes(): found sinez gen %v", gen)
			outLane, err := genSinez(gen)
			if err != nil {
				debugf("renderGenlanes(): gen_sinez() failed")
				outerr = fmt.Errorf("error in genlane#%d: %s", i, err)
				return
			}
			genlanes = append(genlanes, outLane)
		}
		if gen, ok := inLane["place"]; ok {
			debugf("renderGenlanes(): found place gen %v", gen)
			outLane, err := genPlace(gen)
			if err != nil {
				debugf("renderGenlanes(): gen_place() failed")
				outerr = fmt.Errorf("error in genlane#%d: %s", i, err)
				return
			}
			genlanes = append(genlanes, outLane)
		}
	}
	return
}
