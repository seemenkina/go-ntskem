package matrix

import (
	"math/big"
	"reflect"
)

type MatrixFF struct {
	nRows    uint32
	nColumns uint32
	m        []*big.Int
}

func (mff *MatrixFF) New(nr, nc uint32) {
	mff.nRows = nr
	mff.nColumns = nc
	mff.ZeroMatrix()
}

func (mff *MatrixFF) ZeroMatrix() {
	buf := make([]*big.Int, mff.nRows)

	for i := range buf {
		buf[i] = big.NewInt(0)
	}
	mff.m = buf
}

func (mff *MatrixFF) Copy() *MatrixFF {
	duplicate := make([]*big.Int, len(mff.m))
	copy(duplicate, mff.m)

	return &MatrixFF{
		nRows:    mff.nRows,
		nColumns: mff.nColumns,
		m:        duplicate,
	}
}

func (mff *MatrixFF) IsEqual(second *MatrixFF) bool {
	// if mff.nRows != second.nRows || mff.nColumns != second.nColumns || mff == nil || second == nil {
	// 	return false
	// }
	// return reflect.DeepEqual(mff.m, second.m)

	return reflect.DeepEqual(mff, second)
}

func (mff *MatrixFF) ColumnSwap(i, j int) {
	if i == j {
		return
	}

	for r := uint32(0); r < mff.nRows; r++ {
		bi := mff.m[r].Bit(i)
		bj := mff.m[r].Bit(j)
		mff.m[r].SetBit(mff.m[r], i, bj)
		mff.m[r].SetBit(mff.m[r], j, bi)
	}
}

// func (mff *MatrixFF) ReduceRowEchelon () int {
// 	lead := uint32(0)
// 	for r := uint32(0); r < mff.nRows; r++ {
// 		if lead >= mff.nColumns {
// 			return
// 		}
// 		i := r
// 		for mff.m[i][lead] == 0 {
// 			i++
// 			if mff.nRows == i {
// 				i = r
// 				lead++
// 				if mff.nColumns == lead {
// 					return
// 				}
// 			}
// 		}
// 		mff.m[i], mff.m[r] = mff.m[r], mff.m[i]
// 		f := 1 / mff.m[r][lead]
// 		for j, _ := range mff.m[r] {
// 			mff.m[r][j] *= f
// 		}
// 		for i = 0; i < mff.nRows; i++ {
// 			if i != r {
// 				f = mff.m[i][lead]
// 				for j, e := range mff.m[r] {
// 					mff.m[i][j] -= e * f
// 				}
// 			}
// 		}
// 		lead++
// 	}
// }
