package main

import (
	"testing"
)

/*
func genVelocityMap(part part, notes matrix) matrix {

	type part struct {
    Name  string
    Set   noteMap
    Step  int
    Bpm   int
    Fx    []map[string]string
    Lanes []string
}
*/

type genVelocityMapTestInput struct {
	part  part
	notes matrix
}

var genVelocityMapInputList = []genVelocityMapTestInput{
	{
		part: part{
			fx: []map[string]string{
				{
					"rampv": "20-30",
				},
			},
			lanes: []string{
				"ab ab ab ab -- ab",
				"cd cd -- -- cd --",
			},
		},
		notes: matrix{
			row{10, 10, 10, 10, 0, 10},
			row{11, 11, 0, 0, 11, 0},
		},
	},
	{
		part: part{
			fx: []map[string]string{
				{
					"rampv": "127-40",
				},
			},
			lanes: []string{
				"ab ab ab ab -- ab ab -- ab",
				"cd cd -- -- cd -- cd -- cd",
			},
		},
		notes: matrix{
			row{20, 20, 20, 20, 0, 20, 20, 0, 20},
			row{21, 21, 0, 0, 21, 0, 21, 0, 21},
		},
	},
	{
		part: part{
			fx: []map[string]string{
				{
					"rampv": "90-90",
				},
			},
			lanes: []string{
				"ab ab ab ab -- ab ab -- ab",
				"cd cd -- -- cd -- cd -- cd",
			},
		},
		notes: matrix{
			row{20, 20, 20, 20, 0, 20, 20, 0, 20},
			row{21, 21, 0, 0, 21, 0, 21, 0, 21},
		},
	},
}

var genVelocityMapOutputList = []row{
	row{20, 21, 22, 23, 24, 25},
	row{127, 118, 109, 100, 91, 82, 73, 64, 55},
	row{90, 90, 90, 90, 90, 90, 90, 90, 90},
}

func TestRampv(t *testing.T) {
	for i := range genVelocityMapOutputList {
		got := fxRampv(genVelocityMapInputList[i].part, genVelocityMapInputList[i].notes)
		want := genVelocityMapOutputList[i]
		if !got.eq(want) {
			t.Errorf("got %v, wanted %v", got, want)
		}
	}
}
