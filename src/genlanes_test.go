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
		got, err := gen_equid(p.in)
		if got != p.wanted {
			t.Errorf("got %#v, wanted %#v", got, p.wanted)
			if err == nil {
				t.Errorf("also, error was not detected")
			}
		}
	}
}
