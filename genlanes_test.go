package main

import (
	"testing"
)

type genEquidTestPair struct {
	in     map[string]string
	wanted string
}

var genEquidTestPairs = []genEquidTestPair{
	{
		in: map[string]string{
			"note":   "bd",
			"length": "1",
			"dist":   "1",
			"start":  "1",
		},
		wanted: "bd",
	},
	{
		in: map[string]string{
			"note":   "hc",
			"length": "13",
			"dist":   "2",
			"start":  "1",
		},
		wanted: "hc -- hc -- hc -- hc -- hc -- hc -- hc",
	},
	{
		in: map[string]string{
			"note":   "abc",
			"length": "4",
			"dist":   "1",
			"start":  "2",
		},
		wanted: "-- abc abc abc",
	},
	{
		in: map[string]string{
			"note":   "bd",
			"length": "8",
			"dist":   "3",
			"start":  "5",
		},
		wanted: "-- -- -- -- bd -- -- bd",
	},
}

func TestGenEquid(t *testing.T) {
	for _, p := range genEquidTestPairs {
		got, err := genEquid(p.in)
		if got != p.wanted {
			t.Errorf("got %#v, wanted %#v", got, p.wanted)
			if err == nil {
				t.Errorf("also, error was not detected")
			}
		}
	}
}

type genPlaceTestPair struct {
	in     map[string]string
	wanted string
}

var genPlaceTestPairs = []genPlaceTestPair{
	{
		in: map[string]string{
			"note": "hc",
			"pos":  "1 3 5 8",
		},
		wanted: "hc -- hc -- hc -- -- hc",
	},
	{
		in: map[string]string{
			"note": "bd",
			"pos":  "2 9",
		},
		wanted: "-- bd -- -- -- -- -- -- bd",
	},
	{
		in: map[string]string{
			"note": "bd",
			"pos":  "2 9 8",
		},
		wanted: "",
	},
	{
		in: map[string]string{
			"note": "xx",
			"pos":  "1 1 1",
		},
		wanted: "",
	},
}

func TestGenPlace(t *testing.T) {
	for _, p := range genPlaceTestPairs {
		got, err := genPlace(p.in)
		if got != p.wanted {
			t.Errorf("got %#v, wanted %#v", got, p.wanted)
			if err == nil {
				t.Errorf("also, error was not detected")
			}
		}
	}
}

type genEuclidTestPair struct {
	in     map[string]string
	wanted string
}

var genEuclidTestPairs = []genEuclidTestPair{
	{
		in: map[string]string{
			"note":     "bd",
			"length":   "2",
			"accents":  "1",
			"rotation": "0",
		},
		wanted: "bd --",
	},
	{
		in: map[string]string{
			"note":     "hc",
			"length":   "8",
			"accents":  "3",
			"rotation": "0",
		},
		wanted: "hc -- -- hc -- -- hc --",
	},
	{
		in: map[string]string{
			"note":     "x",
			"length":   "12",
			"accents":  "7",
			"rotation": "0",
		},
		wanted: "x -- x -- x x -- x -- x x --",
	},
	{
		in: map[string]string{
			"note":     "bd",
			"length":   "5",
			"accents":  "2",
			"rotation": "1",
		},
		wanted: "-- bd -- -- bd",
	},
}

func TestGenEuclid(t *testing.T) {
	for _, p := range genEuclidTestPairs {
		got, err := genEuclid(p.in)
		if got != p.wanted {
			t.Errorf("got %#v, wanted %#v", got, p.wanted)
			if err == nil {
				t.Errorf("also, error was not detected")
			}
		}
	}
}
