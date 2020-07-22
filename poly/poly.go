package poly

import (
	cr "crypto/rand"
	"math/rand"
	"time"

	"github.com/seemenkina/go-ntskem/ff"
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
	for i := 0; i < tau; i++ {
		buf := make([]byte, 2)
		_, _ = cr.Read(buf)
		g.Pol[i] = uint16((buf[0] << (m - 8)) | (buf[1] >> (m - 8)))
	}
	g.Pol = append(g.Pol, 1)
	pl.degree = tau
	pl.size = size
	pl.Pol = g.Pol
}

func (pl *Polynomial) ModuloReduce(mod *Polynomial, ff2 *ff.FF) *Polynomial {

	for pl.degree >= mod.degree {
		a := ff2.Mul(pl.Pol[pl.degree], ff2.Inv(mod.Pol[mod.degree]))
		j := pl.degree - mod.degree
		for i := 0; i < mod.degree; i++ {
			if mod.Pol[i] != 0 {
				pl.Pol[j] = ff2.Add(pl.Pol[j], ff2.Mul(mod.Pol[i], a))
			}
			j++
		}
		pl.Pol[j] = 0
		for pl.degree >= 0 && pl.Pol[pl.degree] != 0 {
			pl.degree--
		}
	}

	return pl
}

// Generate a length n permutation vector p
func GeneratePermutVector() []uint16 {
	p := make([]uint16, n)
	for i := 0; i < n-1; i++ {
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

// Partition the vectors a and h. Return a = (ab|ac), h = (hb|hc). Parts aa and ha of length k - l
func PartVectors(a, h []uint16) ([]uint16, []uint16) {
	return a[k-l:], h[k-l:]
}
