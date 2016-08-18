package main

type row []int
type matrix []row

func (m matrix) transpose() matrix {
	r := make(matrix, len(m[0]))
	for x, _ := range r {
		r[x] = make(row, len(m))
	}
	for y, s := range m {
		for x, e := range s {
			r[x][y] = e
		}
	}
	return r
}
