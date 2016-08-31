package main

import (
	"testing"
)

type t2mInput struct {
	set noteMap
	txt []string
}

type t2mOutput struct {
	channels matrix
	notes    matrix
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
				"ab ab -- --",
				"-- ab -- ab",
			},
		},
		out: t2mOutput{
			channels: matrix{
				row{1, 1, 0, 0},
				row{0, 1, 0, 1},
			},
			notes: matrix{
				row{60, 60, 0, 0},
				row{0, 60, 0, 60},
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
				". ab . ab",
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
		},
	},
}

func TestT2m(t *testing.T) {
	for _, pair := range t2mTestPairs {
		gotC, gotN := text2matrix(pair.in.set, pair.in.txt)
		if !gotC.eq(pair.out.channels) {
			t.Errorf("got %v, wanted %v", gotC, pair.out.channels)
		}
		if !gotN.eq(pair.out.notes) {
			t.Errorf("got %v, wanted %v", gotN, pair.out.notes)
		}
	}
}
