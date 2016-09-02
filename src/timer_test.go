package main

import (
	"testing"
	"time"
)

type makeTimingTestInput struct {
	bpm       string
	step      int
	numEvents int
}

type makeTimingTestOutput struct {
	timing          time.Duration
	timingIncrement string
	err             error
}

type makeTimingTestPair struct {
	in  makeTimingTestInput
	out makeTimingTestOutput
}

//func makeTiming(bpm string, step int, numEvents int) (timing, timingIncrement time.Duration, err error) {

var makeTimingTestPairs = []makeTimingTestPair{
	{
		in: makeTimingTestInput{
			bpm:       "120",
			step:      16,
			numEvents: 16,
		},
		out: makeTimingTestOutput{
			timing:          time.Millisecond * 125,
			timingIncrement: "0ms",
		},
	},
	{
		in: makeTimingTestInput{
			bpm:       "120-250",
			step:      16,
			numEvents: 16,
		},
		out: makeTimingTestOutput{
			timing:          time.Millisecond * 125,
			timingIncrement: "-4.0625ms",
		},
	},
	{
		in: makeTimingTestInput{
			bpm:       "300-200",
			step:      8,
			numEvents: 25,
		},
		out: makeTimingTestOutput{
			timing:          time.Millisecond * 100,
			timingIncrement: "2ms",
		},
	},
}

func TestMakeTiming(t *testing.T) {
	for _, pair := range makeTimingTestPairs {
		gotTiming, gotTimingIncrement, _ := makeTiming(pair.in.bpm, pair.in.step, pair.in.numEvents)
		if gotTiming != pair.out.timing {
			t.Errorf("got %v, wanted %v", gotTiming, pair.out.timing)
		}
		ti, _ := time.ParseDuration(pair.out.timingIncrement)
		if gotTimingIncrement != ti {
			t.Errorf("got %v, wanted %v", gotTimingIncrement, ti)
		}
	}
}
