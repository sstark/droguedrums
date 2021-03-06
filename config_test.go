package main

import (
	"math/rand"
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

type tlfTestPair struct {
	in  string
	out midiFigure
}

var tlfTestPairs = []tlfTestPair{
	{
		in:  "f1",
		out: midiFigure{{0, 0, 0}, {7, 65, 65}, {7, 65, 65}},
	},
	{
		in:  "f2",
		out: midiFigure{{0, 0, 0}, {0, 0, 0}, {5, 50, 100}},
	},
	{
		in:  "h1",
		out: midiFigure{{0, 0, 0}, {0, 0, 0}, {5, 70, 80}},
	},
	{
		in:  "qu",
		out: midiFigure{{5, 58, 95}, {5, 58, 95}, {5, 58, 95}},
	},
}

func TestFigures(t *testing.T) {
	drums := new(drums)
	drums.loadFromFile("testfiles/beat8.yml")
	sets := drums.getSets()
	figures := drums.getFigures()
	for _, fig := range tlfTestPairs {
		got := translateFigure(sets[defaultSet], figures, fig.in)
		if !reflect.DeepEqual(got, fig.out) {
			t.Errorf("%s: got %v, wanted %v", fig.in, got, fig.out)
		}
	}
}

type seqFlattenTestPair struct {
	in    seqMap
	start string
	out   []string
}

var seqFlattenTestPairs = []seqFlattenTestPair{
	{
		in: seqMap{
			"bla":   {"b"},
			"blu":   {":bla"},
			"start": {"a", ":bla", ":blu", ":bla"},
		},
		start: "start",
		out:   []string{"a", "b", "b", "b"},
	},
	{
		in: seqMap{
			"a":     {"b"},
			"b":     {"c", ":a", "x"},
			"start": {"d", ":a", ":b", ":a"},
		},
		start: "start",
		out:   []string{"d", "b", "c", "b", "x", "b"},
	},
}

func TestSeqFlatten(t *testing.T) {
	for _, pair := range seqFlattenTestPairs {
		got := pair.in.flatten(pair.start)
		if !reflect.DeepEqual(got, pair.out) {
			t.Errorf("got %v, wanted %v", got, pair.out)
		}
	}
}

type translateKitSeries []struct {
	channel int
	note    int
}

type translateKitPair struct {
	inset noteMap
	instr string
	out   translateKitSeries
}

var translateKitPairs = []translateKitPair{
	{
		inset: noteMap{
			"ab": midiNote{
				Channel: 2,
				Note:    70,
			},
			"cd": midiNote{
				Channel: 3,
				Note:    80,
			},
		},
		instr: "ab",
		out:   translateKitSeries{{2, 70}, {2, 70}, {2, 70}},
	},
	{
		inset: noteMap{
			"ab": midiNote{
				Channel: 2,
				Note:    70,
			},
			"cd": midiNote{
				Channel: 3,
				Note:    80,
			},
		},
		instr: "ab|cd",
		out:   translateKitSeries{{2, 70}, {3, 80}, {2, 70}, {3, 80}, {3, 80}},
	},
	{
		inset: noteMap{
			"ab": midiNote{
				Channel: 2,
				Note:    70,
			},
			"cd": midiNote{
				Channel: 3,
				Note:    80,
			},
			"ef": midiNote{
				Channel: 4,
				Note:    90,
			},
		},
		instr: "ab|cd|ef",
		out: translateKitSeries{{2, 70}, {2, 70}, {2, 70}, {4, 90},
			{4, 90}, {2, 70}, {2, 70}, {2, 70}, {4, 90}, {4, 90}, {3, 80}},
	},
}

func TestTranslateKit(t *testing.T) {
	// for each iteration a different random value is chosen, which
	// is reflected in the translateKitSeries slice. Because we use
	// pseudo random numbers it's always the same ones.
	rand.Seed(0)
	for _, pair := range translateKitPairs {
		for i := range pair.out {
			outch, outn := translateKit(pair.inset, pair.instr)
			if outch != pair.out[i].channel {
				t.Errorf("%d: got %v, wanted %v", i, outch, pair.out[i].channel)
			}
			if outn != pair.out[i].note {
				t.Errorf("%d: got %v, wanted %v", i, outn, pair.out[i].note)
			}
		}
	}
}
