package main

import (
	"strconv"
	"strings"
	"time"
)

func makeTiming(bpm string, step int, numEvents int) (timing, timingIncrement time.Duration, err error) {
	var bpmA, bpmB int
	var timingB time.Duration
	bpmList := strings.Split(bpm, "-")
	debugf("makeTiming(): bpmList: %v, numEvents: %v", bpmList, numEvents)
	switch len(bpmList) {
	case 1:
		bpmA, err = strconv.Atoi(bpmList[0])
		if err != nil {
			return
		}
	case 2:
		bpmA, err = strconv.Atoi(bpmList[0])
		if err != nil {
			return
		}
		bpmB, err = strconv.Atoi(bpmList[1])
		if err != nil {
			return
		}
	default:
		return
	}
	timing = time.Minute / time.Duration(bpmA) / time.Duration(step/4)
	debugf("makeTiming(): timing: %v, bpmA: %v, bpmB: %v", timing, bpmA, bpmB)
	if (bpmA == bpmB) || bpmB == 0 {
		timingIncrement = 0
	} else {
		timingB = time.Minute / time.Duration(bpmB) / time.Duration(step/4)
		timingIncrement = (timingB - timing) / time.Duration(numEvents)
		debugf("makeTiming(): timingB: %v, timingIncrement: %v", timingB, timingIncrement)
	}
	return
}

type varTicker struct {
	C    <-chan time.Time
	ch   chan<- time.Time
	t    *time.Ticker
	done chan bool
}

func (t *varTicker) SetDuration(d time.Duration) {
	if t.t != nil {
		t.t.Stop()
		close(t.done)
	} else {
		var ticker = make(chan time.Time)
		t.C = ticker
		t.ch = ticker
	}
	t.done = make(chan bool)
	t.t = time.NewTicker(d)
	go func(out chan<- time.Time, in <-chan time.Time, done <-chan bool) {
		for {
			select {
			case tick := <-in:
				out <- tick
			case <-done:
				return
			}
		}
	}(t.ch, t.t.C, t.done)
}
