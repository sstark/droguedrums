package main

import (
	"testing"
)

type eqTestPair struct {
	a, b row
}

type transposeTestPair struct {
	normal, transposed matrix
}

var eqPairs = []eqTestPair{
	{
		a: row{1, 2, 3},
		b: row{1, 2, 3},
	},
}

var neqPairs = []eqTestPair{
	{
		a: row{1, 2, 3},
		b: row{1, 2, 4},
	},
	{
		a: row{1, 2, 3},
		b: row{1, 2},
	},
	{
		a: row{1, 2},
		b: row{1, 2, 3},
	},
}

var transposePairs = []transposeTestPair{
	{
		normal: matrix{
			row{1, 2, 3},
			row{4, 5, 6},
		},
		transposed: matrix{
			row{1, 4},
			row{2, 5},
			row{3, 6},
		},
	},
	{
		normal: matrix{
			row{1, 2, 3},
		},
		transposed: matrix{
			row{1},
			row{2},
			row{3},
		},
	},
	{
		normal: matrix{
			row{1, 2, 3},
			row{1, 2},
			row{1, 2, 3},
		},
		transposed: matrix{
			row{1, 1, 1},
			row{2, 2, 2},
			row{3, 0, 3},
		},
	},
}

func TestEq(t *testing.T) {
	for _, r := range eqPairs {
		if !r.a.eq(r.b) {
			t.Errorf("equality test failed: %v == %v", r.a, r.b)
		}
	}
	for _, r := range neqPairs {
		if r.a.eq(r.b) {
			t.Errorf("non-equality test failed: %v != %v", r.a, r.b)
		}
	}
}

func assertMatrixEq(t *testing.T, got matrix, want matrix) {
	var eq bool = true
	for i := range got {
		if got[i].eq(want[i]) == false {
			eq = false
			break
		}
	}
	if eq == false {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestTranspose(t *testing.T) {
	for _, pair := range transposePairs {
		assertMatrixEq(t, pair.normal.transpose(), pair.transposed)
	}
}
