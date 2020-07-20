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

// Randomly generate a Goppa polynomial of degree tau.
func GenerateGoppaPol(tau int) []uint16 {
	g := make([]uint16, tau)
	for i := 0; i < tau; i++ {
		buf := make([]byte, 2)
		_, _ = cr.Read(buf)
		g[i] = uint16((buf[0] << (m - 8)) | (buf[1] >> (m - 8)))
	}
	g = append(g, 1)
	return g
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

// Check the validity of Goppa polynomial
// 1 - g0 != 0
// 2 - Goppa polynomial  has no roots in F{2^m}. Check by Additive FFT
// 3 - Goppa polynomial  has no repeated roots in any extension 􏰀field.  Check by GCD(G)
func CheckGoppaPoly(g []uint16) bool {

	if g[0] == 0 {
		return false
	}
	if !checkFft(g) {
		return false
	}

	// if (formal_derivative_poly(Gz, Dz)) {
	// 	if (gcd_poly(ff2m, Gz, Dz, Fz)) {
	// 		status = (Fz->degree < 1);
	// 	}
	// }

	return true
}

// Return true, if Goppa polynomial  has roots in  F{2^m}
func checkFft(pol []uint16) bool {
	return false
}

// Partition the vectors a and h. Return a = (ab|ac), h = (hb|hc). Parts aa and ha of length k - l
func PartVectors(a, h []uint16) ([]uint16, []uint16) {
	return a[k-l:], h[k-l:]
}
