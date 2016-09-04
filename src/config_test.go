package main

import (
	"reflect"
	"testing"
)

type t2mInput struct {
	set     noteMap
	txt     []string
	figures map[string]figure
}

type t2mOutput struct {
	channels matrix
	notes    matrix
	mfigures map[int][]midiFigure
}

type t2mTestPair struct {
	in  t2mInput
	out t2mOutput
}

var t2mTestPairs = []t2mTestPair{
	{
		in: t2mInput{
			set: noteMap{
				"ab": midiNote{
					Channel: 1,
					Note:    60,
				},
			},
			txt: []string{
				"ab +f1 -- --",
				"--  ab -- ab",
			},
			figures: map[string]figure{
				"f1": figure{"ab", 111, "x.x"},
			},
		},
		out: t2mOutput{
			channels: matrix{
				row{1, 0, 0, 0},
				row{0, 1, 0, 1},
			},
			notes: matrix{
				row{60, 0, 0, 0},
				row{0, 60, 0, 60},
			},
			mfigures: map[int][]midiFigure{
				1: []midiFigure{
					midiFigure{
						midiEvent{1, 60, 111},
						midiEvent{},
						midiEvent{1, 60, 111},
					},
				},
			},
		},
	},
	{
		in: t2mInput{
			set: noteMap{
				"ab": midiNote{
					Channel: 2,
					Note:    70,
				},
			},
			txt: []string{
				"ab ab . . .",
				". ab +f1 ab",
			},
			figures: map[string]figure{
				"f1": figure{"ab", 109, ".x."},
			},
		},
		out: t2mOutput{
			channels: matrix{
				row{2, 2, 0, 0, 0},
				row{0, 2, 0, 2},
			},
			notes: matrix{
				row{70, 70, 0, 0, 0},
				row{0, 70, 0, 70},
			},
			mfigures: map[int][]midiFigure{
				2: []midiFigure{
					midiFigure{
						midiEvent{},
						midiEvent{2, 70, 109},
						midiEvent{},
					},
				},
			},
		},
	},
}

func TestT2m(t *testing.T) {
	for _, pair := range t2mTestPairs {
		gotC, gotN, gotF := text2matrix(pair.in.set, pair.in.figures, pair.in.txt)
		if !gotC.eq(pair.out.channels) {
			t.Errorf("got %v, wanted %v", gotC, pair.out.channels)
		}
		if !gotN.eq(pair.out.notes) {
			t.Errorf("got %v, wanted %v", gotN, pair.out.notes)
		}
		if !reflect.DeepEqual(gotF, pair.out.mfigures) {
			t.Errorf("got %v, wanted %v", gotF, pair.out.mfigures)
		}
	}
}
