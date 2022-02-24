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
		if posI <= lastI {
			return "", errors.New("positions must be sorted and unique")
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

func genEuclid(gl map[string]string) (out string, err error) {
	//{note: hc, length: 15, accents: 4, rotation: 1}
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
	inpAccents, ok := gl["accents"]
	if !ok {
		err = errors.New("accents value missing")
		return
	}
	inpRotation, ok := gl["rotation"]
	if !ok {
		inpRotation = "0"
	}
	length, err := strconv.Atoi(inpLength)
	if err != nil {
		return
	}
	accents, err := strconv.Atoi(inpAccents)
	if err != nil {
		return
	}
	if accents >= length {
		err = errors.New("accents must be smaller than length")
		return
	}
	rotation, err := strconv.Atoi(inpRotation)
	if err != nil {
		return
	}

	debugf("genEuclid(): len %v, acc %v, rot %v", length, accents, rotation)

	var l int
	var c, rem []int
	var divisor = length - accents

	rem = append(rem, accents)
	for {
		c = append(c, divisor/rem[l])
		rem = append(rem, divisor%rem[l])
		divisor = rem[l]
		l += 1
		if rem[l] < 2 {
			break
		}
	}
	c = append(c, divisor)

	var gen func(int)
	gen = func(l int) {
		if l == -1 {
			buffer.WriteString("-- ")
		} else if l == -2 {
			buffer.WriteString(inpNote + " ")
		} else {
			for i := 0; i < c[l]; i++ {
				gen(l - 1)
			}
			if rem[l] != 0 {
				gen(l - 2)
			}
		}
	}

	gen(l)

	// i = self.pattern.index(1)
	// self.pattern = self.pattern[i:] + self.pattern[0:i]

	out = strings.TrimSpace(buffer.String())
	debugf("genEuclid(): %v", out)
	return
}

func renderGenlanes(lanes []map[string]map[string]string) (genlanes []string, outerr error) {
	genFuncs := map[string]func(map[string]string) (string, error){
		"equid":  genEquid,
		"sinez":  genSinez,
		"place":  genPlace,
		"euclid": genEuclid,
	}
	for i, inLane := range lanes {
		for k, v := range genFuncs {
			if gen, ok := inLane[k]; ok {
				debugf("renderGenlanes(): found %s gen %v", k, gen)
				outLane, err := v(gen)
				if err != nil {
					debugf("renderGenlanes(): gen%s() failed", k)
					outerr = fmt.Errorf("error in genlane#%d: %s", i, err)
					return
				}
				genlanes = append(genlanes, outLane)
			}
		}
	}
	return
}
