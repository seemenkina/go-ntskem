package go_ntskem

import (
	"crypto/rand"
	"encoding/binary"
	"strconv"
	"strings"

	ff2 "github.com/seemenkina/go-ntskem/ff"
	"github.com/seemenkina/go-ntskem/matrix"
	"github.com/seemenkina/go-ntskem/poly"
)

const (
	// m   = 12 // m = log n
	tau = 64 // tau = (d - 1) / 2
	// n   = 1 << m
	// k   = n - tau*m
	l = 256
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

func (nk *NTSKEM) New(m int) {
	nk.ff = &ff2.FF{}
	nk.ff.New(m)
	nk.n = 1 << m
}

// GenerateKey returns a new public/private key pair
func (nk *NTSKEM) GenerateKey() {

	// // Step 1: Generate Goppa polynomial of degree τ
	// g := poly.Polynomial{}
	// g.GenerateGoppaPol(tau, 1<<nk.ff.M)
	// for !nk.ff.CheckGoppaPoly(&g) {
	// 	g := poly.Polynomial{}
	// 	g.GenerateGoppaPol(tau, 1<<nk.ff.M)
	// }

	g := poly.Polynomial{}
	g.New(1 << nk.ff.M)
	g.SetDegree(tau)

	rawHex := "1EE 677 162 5EC 23B AA7 076 A65 A3B 519 000 B04 F3C E70 504 C07 B46 BC3 045 BAA 95B 807 6DD EE4 FF8 B02 362 500 077 42D 6F3 BB0 163 049 D0E D90 165 FDF 1A9 83D FCA CC9 FB4 C08 110 84B 0AB 330 9E3 985 DE7 17B 2A4 A95 9C6 BF1 DD8 8A9 2AC 652 4ED 2A2 CEA D1C 001"
	rawHex = strings.Replace(rawHex, " ", "", -1)
	// rawHex = strings.ToLower(rawHex)
	for i := 0; len(rawHex) > 0; i++ {
		b, _ := strconv.ParseUint(rawHex[:3], 16, 64)
		g.Pol[i] = uint16(b)
		rawHex = rawHex[3:]
	}
	// Step 2: Randomly generate a permutation vector p of length n
	p := poly.GeneratePermutVector()

	// Step 3: Construct a generator matrix
	Q := matrix.MatrixFF{}
	a, h := Q.CreateMatrixG(&g, p, nk.ff, tau)

	// Step 4: Generate random number of length l
	buf := make([]byte, l/8)
	_, _ = rand.Read(buf)
	z := binary.LittleEndian.Uint32(buf)

	// Step 5: Partition the vectors a and h from step 3. Finally define a* and h*
	aSt, hSt := poly.PartVectors(a, h)

	pk := PublicKey{
		Q:   Q,
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
