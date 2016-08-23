package main

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func gen_equid(gl map[string]string) (out string, err error) {
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
	debugf("gen_equid(): %v", out)
	return
}

func renderGenlanes(lanes []map[string]map[string]string) (genlanes []string, outerr error) {
	for i, inLane := range lanes {
		if gen, ok := inLane["equid"]; ok {
			debugf("renderGenlanes(): found equid gen %v", gen)
			outLane, err := gen_equid(gen)
			if err != nil {
				debugf("renderGenlanes(): gen_equid() failed")
				outerr = fmt.Errorf("error in genlane#%d: %s", i, err)
				return
			}
			genlanes = append(genlanes, outLane)
		}
	}
	return
}
