package main

import (
	"errors"
)

type row []int
type matrix []row

// check returns an error if matrix is irregular
// and cannot be transposed.
func (m matrix) check() (err error) {
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
