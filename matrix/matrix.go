package matrix

import (
	"reflect"

	"github.com/seemenkina/go-ntskem/ff"
	"github.com/seemenkina/go-ntskem/poly"
)

type MatrixFF struct {
	nRows    uint32
	nColumns uint32
	m        [][]uint16
}

func (mff *MatrixFF) New(nr, nc uint32) {
	mff.nRows = nr
	mff.nColumns = nc
	mff.ZeroMatrix()
}

func (mff *MatrixFF) ZeroMatrix() {
	buf := make([][]uint16, mff.nRows)

	for i := range buf {
		buf[i] = make([]uint16, mff.nColumns)
	}
	mff.m = buf
}

func (mff *MatrixFF) Copy() *MatrixFF {
	duplicate := make([][]uint16, len(mff.m))
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
		mff.m[r][i], mff.m[r][j] = mff.m[r][j], mff.m[r][i]
	}
}

func (mff *MatrixFF) ReduceRowEchelon(ff2 *ff.FF) int {
	lead := uint32(0)
	for r := uint32(0); r < mff.nRows; r++ {
		if lead >= mff.nColumns {
			return mff.GetRank()
		}
		i := r
		for mff.m[i][lead] == 0 {
			i++
			if mff.nRows == i {
				i = r
				lead++
				if mff.nColumns == lead {
					return mff.GetRank()
				}
			}
		}
		mff.m[i], mff.m[r] = mff.m[r], mff.m[i]
		f := ff2.Inv(mff.m[r][lead])
		for j := range mff.m[r] {
			mff.m[r][j] = ff2.Mul(mff.m[r][j], f)
		}
		for i = 0; i < mff.nRows; i++ {
			if i != r {
				f = mff.m[i][lead]
				for j, e := range mff.m[r] {
					mff.m[i][j] = ff2.Add(mff.m[i][j], ff2.Mul(e, f))
				}
			}
		}
		lead++
	}
	return mff.GetRank()
}

func (mff *MatrixFF) GetRank() int {
	rank := 0
	for i := 0; i < int(mff.nRows); i++ {
		for j := 0; j < int(mff.nColumns); j++ {
			if mff.m[i][j] != 0 {
				break
			}
		}
		rank = i + 1
	}
	return rank
}

func (mff *MatrixFF) CreateMatrixG(pol *poly.Polynomial, p []uint16, ff2 *ff.FF, degree int) ([]uint16, []uint16) {
	n := 1 << ff2.M
	k := n - degree*ff2.M
	a := make([]uint16, n)
	h := make([]uint16, n)
	aPr := make([]uint16, n)
	hPr := make([]uint16, n)

	aPr[0] = 0
	for i := 0; i < n; i++ {
		aPr[i] = 0
		for j := 0; j < ff2.M; j++ {
			aPr[i] ^= ((i & (1 << j)) >> j) * ff2.Basis[j]
		}
	}
	for i := 0; i < n; i++ {
		a[i] = aPr[p[i]]
	}

	hPr = ff2.AdaptiveFft(pol)
	if hPr == nil {
		return nil, nil
	}

	for i := 0; i < n; i++ {
		h[i] = ff2.Sqr(ff2.Inv(hPr[p[i]]))
	}

	H := MatrixFF{}
	H.New(uint32(degree*ff2.M), n)

	for i := 0; i < n; i++ {
		hPr[i] = 1
	}

	for i := 0; i < degree; i++ {
		for j := 0; j < n; j++ {
			e := uint16(ff2.M) - 1
			for e > 0 {
				f := ff2.Mul(hPr[j], h[j])
				if f&(1<<e) != 0 {
					H.m[i*ff2.M+(ff2.M-int(e)-1)][j] = uint16(j)
				}
				e--
			}
			hPr[j] = ff2.Mul(hPr[j], a[j])
		}
	}

	rank := H.ReduceRowEchelon(ff2)
	if n-degree*ff2.M != n-rank {
		return nil, nil
	}

	for j, i := 0, int(H.nRows-1); i >= 0; i-- {
		for H.m[i][int(H.nColumns)-j-1] == 0 {
			j++
		}
		H.ColumnSwap(k+i, int(H.nColumns)-j-1)

		p[k+i], p[int(H.nColumns)-j-1] = p[int(H.nColumns)-j-1], p[k+i]
		a[k+i], a[int(H.nColumns)-j-1] = a[int(H.nColumns)-j-1], a[k+i]
		h[k+i], h[int(H.nColumns)-j-1] = h[int(H.nColumns)-j-1], h[k+i]
	}

	Q := MatrixFF{}

	Q.New(k, n-k)
	for i := 0; i < n-k; i++ {
		for j := 0; j < k; j++ {
			Q.m[j][i] = H.m[j][j]
		}
	}

	mff.nRows = Q.nRows
	mff.nColumns = Q.nColumns
	copy(mff.m, Q.m)

	return a, h
}
