package poly

import (
	cr "crypto/rand"
	"math/rand"
	"time"
)

const (
	m   = 12 // m = log n
	tau = 64 // tau = (d - 1) / 2
	n   = 1 << m
	k   = n - tau*m
	l   = 256
)

type Polynomial struct {
	degree int
	size   int
	Pol    []uint16
}

func (pl *Polynomial) Size() int {
	return pl.size
}

func (pl *Polynomial) SetSize(size int) {
	pl.size = size
}

func (pl *Polynomial) New(size int) {
	pl.degree = -1
	pl.size = size
	pl.Pol = make([]uint16, size)
}

func (pl *Polynomial) GetDegree() int {
	return pl.degree
}
func (pl *Polynomial) SetDegree(d int) {
	pl.degree = d
}

// Randomly generate a Goppa polynomial of degree tau.
func (pl *Polynomial) GenerateGoppaPol(tau, size int) {
	g := Polynomial{}
	g.New(size)
	g.SetDegree(tau)
	for i := 0; i <= tau; i++ {
		buf := make([]byte, 2)
		_, _ = cr.Read(buf)
		g.Pol[i] = uint16((buf[0] << (m - 8)) | (buf[1] >> (m - 8)))
	}
	g.Pol = append(g.Pol, 1)
	pl.degree = tau
	pl.size = size
	pl.Pol = g.Pol
}

// Generate a length n permutation vector p
func GeneratePermuteVector() []uint16 {
	p := make([]uint16, n)
	for i := 0; i < n; i++ {
		p[i] = uint16(i)
	}
	p = FisherYatesShuffle(p)
	return p
}

// Algorithm generates unbiased permutations of n elements in linear time
func FisherYatesShuffle(slice []uint16) []uint16 {
	rand.Seed(time.Now().UnixNano())
	n := len(slice)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

// Create a random vector of length n bits with
// Hamming weight t
func RandomVector(n uint16, t uint16) []uint16 {
	e := make([]uint16, n-t)
	for i := n - t; i < n; i++ {
		e = append(e, 1)
	}
	// shuffle
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	for i := n - 1; i >= n-t; {
		r := random.Intn(int(i))
		e[i], e[r] = e[r], e[i]
		i = i - 1
	}
	return e
}

// Polynom Sum
// Использовать сумму из ff, m =1 ?
func PolySum(BitArray1, BitArray2 []uint16) []uint16 {
	if len(BitArray2) > len(BitArray1) {
		BitArray1, BitArray2 = BitArray2, BitArray1
	}
	var SumArray = make([]uint16, len(BitArray1))
	copy(SumArray, BitArray1)
	for i := 0; i < len(BitArray2); i++ {
		SumArray[i+len(BitArray1)-len(BitArray2)] ^= BitArray2[i]
	}
	return SumArray
}
