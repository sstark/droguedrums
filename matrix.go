package main

import (
	"errors"
)

type row []int
type matrix []row

func (r1 row) eq(r2 row) bool {
	if len(r1) != len(r2) {
		return false
	}
	for i := range r1 {
		if r1[i] != r2[i] {
			return false
		}
	}
	return true
}

func (m matrix) eq(m2 matrix) bool {
	if len(m) != len(m2) {
		return false
	}
	for i := range m {
		if !m[i].eq(m2[i]) {
			return false
		}
	}
	return true
}

// check returns an error if matrix is irregular
// and cannot be transposed.
func (m matrix) check() (err error) {
	debugf("check(): %v", m)
	l := len(m[0])
	for _, x := range m {
		// as long as first line is the longest we're fine
		if len(x) > l {
			return errors.New("first lane is not longest")
		}
	}
	return nil
}

func (m matrix) transpose() matrix {
	r := make(matrix, len(m[0]))
	for x := range r {
		r[x] = make(row, len(m))
	}
	for y, s := range m {
		for x, e := range s {
			r[x][y] = e
		}
	}
	return r
}
