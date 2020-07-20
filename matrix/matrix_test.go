package matrix

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/seemenkina/go-ntskem/ff"
)

func TestMatrixFF_ReduceRowEchelon(t *testing.T) {
	mf := MatrixFF{
		nRows:    3,
		nColumns: 5,
		m: [][]uint16{
			{1, 0, 1, 0, 1},
			{1, 0, 0, 1, 0},
			{0, 0, 0, 1, 0},
		},
	}

	ff2 := ff.FF{}
	ff2.New(12)

	r := mf.ReduceRowEchelon(&ff2)
	print(mf.m)
	spew.Dump(r)
}

func print(m [][]uint16) {
	for e := 0; e < len(m); e += 1 {
		fmt.Printf("%d ", m[e])
		fmt.Println()
	}
}
