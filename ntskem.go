package go_ntskem

import (
	"crypto/rand"
	"encoding/binary"

	ff2 "github.com/seemenkina/go-ntskem/ff"
	"github.com/seemenkina/go-ntskem/matrix"
	"github.com/seemenkina/go-ntskem/poly"
)

const (
	m   = 12 // m = log n
	tau = 64 // tau = (d - 1) / 2
	n   = 1 << m
	k   = n - tau*m
	l   = 256
)

type PublicKey struct {
	Q   matrix.MatrixFF
	tau int // const
	l   int
}

type PrivateKey struct {
	a  []uint16
	h  []uint16
	p  []uint16 // vector of length n
	z  uint32   // random number of length l
	pk PublicKey
}

type NTSKEM struct {
	ff      *ff2.FF
	publKey PublicKey
	privKey PrivateKey
	n       int
}

func (nk *NTSKEM) New(m, l int) {
	nk.ff.New(m)
	nk.n = 1 << m
}

// GenerateKey returns a new public/private key pair
func (nk *NTSKEM) GenerateKey() {

	// Step 1: Generate Goppa polynomial of degree τ
	g := poly.GenerateGoppaPol(tau)
	for !poly.CheckGoppaPoly(g) {
		g = poly.GenerateGoppaPol(tau)
	}
	// Step 2: Randomly generate a permutation vector p of length n
	p := poly.GeneratePermutVector()

	// Step 3: Construct a generator matrix

	a := make([]uint16, nk.n)
	h := make([]uint16, nk.n)

	// Step 4: Generate random number of length l
	buf := make([]byte, l/8)
	_, _ = rand.Read(buf)
	z := binary.LittleEndian.Uint32(buf)

	// Step 5: Partition the vectors a and h from step 3. Finally define a* and h*
	aSt, hSt := poly.PartVectors(a, h)

	pk := PublicKey{
		Q:   nil,
		tau: tau,
		l:   l,
	}
	sk := PrivateKey{
		a:  aSt,
		h:  hSt,
		p:  p,
		z:  z,
		pk: pk,
	}
	nk.privKey = sk
	nk.publKey = pk
}

// Encapsulate uses a given public key produce random key
// and compute ciphertext encapsulating this key
func (nk *NTSKEM) Encapsulate() {}

// Decapsulate uses a private key to decrypt a ciphertext
func (nk *NTSKEM) Decapsulate() {}
