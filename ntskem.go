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
	tau = 64 // tau = (d - 1) / 2
	l   = 256
	k = 3328
	n = 4096
	r = 768 //n-k
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
func (pk PublicKey) Encapsulate() ([]uint16, []uint16){
	//1. Random error vector
	e:= randomVector(n,tau)

	//2. Partition e
	ea := e[:k-l]
	eb := e[k-l:k]
	ec := e[k:n]

	//3. Compute ke, using pseudorandom func - sha3
	ke := sha3.Sum256(BitArrayToByteArray(e))
	keInBits := ByteArrayToBitArray(ke[:32])

	//4. Construct message vector
	var m = ea
	var i uint16 = 0
	for ; i < l; i++ {
		m = append(m, keInBits[i])
	}
	//5. Encoding of m with Q:
	cb := PolySum(keInBits, eb)
	cc := PolySum(pk.Q.PolyOnMatriceMult(m), ec)
	//cc := PolySum(matrix.PolyOnMatriceMult(m, pk.Q), ec)

	var co []uint16
	for i=0; i < l;i++{
		co=append(co,cb[i])
	}
	for i=0; i < uint16(len(cc));i++{
		co=append(co,cc[i])
	}
	var seed []uint16
	for i=0; i < l;i++{
		seed=append(seed,keInBits[i])
	}
	for i=0; i < uint16(len(e));i++{
		seed=append(seed,e[i])
	}

	kr := sha3.Sum256(BitArrayToByteArray(seed))
	return ByteArrayToBitArray(kr[:32]), co
}

// Decapsulate uses a private key to decrypt a ciphertext
func (sk PrivateKey) Decapsulate (c []uint16) []uint16 {
	//field gen
	var ff ff2.FF
	ff.New(12)
	/*ab :=a[:l]
	ac:=a[l:]

	hb := h[:l]
	hc:=h[l:]*/
	
	//1.b Build Transpose matrix.
	Q := matrix.MatrixFF{}
	Q.New(l+r, 2*tau)
	Q.CreateMatrixH(sk.a,sk.h, ff)

	//1.c  Compute all 2τ syndromes of c*
	var s= Q.PolyOnMatriceMult(c)
	
	//1.d Compute the error locator polynomial σ(x) and the first coordinate error indicator ξ 
	var sigma, xi = ff.BerlekampMasseyAlgorithm(s)
	
	//1.e Evaluate the polynomial σ(x) on all elements of F2m
	var A =	ff.Roots(sigma)
	
	//1.f  obtain the error vector e`
	var e = make([]uint16, n)
	for i:=0; i <n; i++ {
		if A[i] == 0{
			e[i] = 1
		}
	}
	if xi == 1 {
		e[0] = 1
	}
	
	//2.Apply the permutation
	e = Permutation(e,sk.p)
	
	//3.Consider e = (ea | eb | ec), and compute ke = cb − eb.
	var eb = e[k-l:k]
	var cb = c[:l]
	var cc = c[l:]
	var ke = PolySum(eb,cb)

	Hl := sha3.Sum256(BitArrayToByteArray(e))
	var q = BitArrayToByteArray(ke)
	
	//4 Verify that H`(e) = ke and hw(e) = τ...
	var str1 = string(Hl[:])
	var str2 = string(q[:])
	var hw uint16
	for i:=0;i<len(e); i++ {
		if e[i] == 1 {
			hw++
		}
	}
	if (str1 == str2 && hw == tau) {
		var out = ke
		for i:=0 ; i < len(e); i++ {
			out = append(out, e[i])
		}
		var kr = sha3.Sum256(BitArrayToByteArray(out))
		return ByteArrayToBitArray(kr[:])
	} else {
		var zInBitArray []uint16
		for  sk.z > 0 {
			zInBitArray = append(zInBitArray, uint16(sk.z & 0x0001)) //Краш из-за длины z 
			sk.z = sk.z >> 1 //Не очень хорошо. Стоит завести локальную переменную
		}

		var out = zInBitArray
		//1a
		for i:=0 ; i < k-l; i++ {
			out = append(out, 1)
		}
		//cb
		for i:=0 ; i < l; i++ {
			out = append(out, cb[i])
		}
		//cc
		for i:=0 ; i < l; i++ {
			out = append(out, cc[i])
		}
		
		var kr = sha3.Sum256(BitArrayToByteArray(out))
		return ByteArrayToBitArray(kr[:])
	}
}
//added
//Использовать генерацию из poly?
func randomVector (n uint16, t uint16) []uint16 {
	var e []uint16 = make([]uint16, n-t)
	for i:=n-t; i <n;i++{
		e=append(e,1)
	}
	//shuffle
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	for i:= n-1; i >= n-t; {
		r:= random.Intn(int(i))
		e[i], e[r] = e[r], e[i]
		i=i-1
	}
	return e
}

//Converts BitArray to ByteArray
func BitArrayToByteArray (BitArray []uint16) []byte{
	println()
	print("BITARRAYLEN: ", len(BitArray))
	var ByteArray []byte = make ([]byte, len(BitArray)/8)
	for i:=0; i < len(BitArray); i++ {
		ByteArray[i/8] = ByteArray[i/8] | byte(BitArray[i])
		if i % 8 != 7 {
			ByteArray[i/8] <<= 1
		}
	}
	return ByteArray
}

//added
//Converts ByteArray to BitArray
func ByteArrayToBitArray(ByteArray []byte) []uint16{
	ByteArraylen := len(ByteArray)
	var BitArray []uint16 = make([]uint16, ByteArraylen*8)
	for i:=0; i < ByteArraylen; i++ {
		for j:=7; j >= 0; j--{
			BitArray[8*i + j] = uint16(ByteArray[i]) & 1
			ByteArray[i] >>= 1
		}
	}
	return BitArray
}

//added
//Polynom Sum
//Использовать сумму из ff, m =1 ?
func PolySum(BitArray1, BitArray2 []uint16) []uint16 {
	if len(BitArray2) > len(BitArray1) {
		BitArray1, BitArray2 = BitArray2, BitArray1
	}
	var SumArray = make ([]uint16, len(BitArray1))
	copy(SumArray, BitArray1)
	for i:=0; i < len(BitArray2); i++ {
		SumArray[i+len(BitArray1)- len(BitArray2)] ^= BitArray2[i]
	}
	return SumArray
}

//added
func Permutation(e,p []uint16) []uint16{
	var out = make([]uint16, len(e))
	for i := 0; i < len(out); i++  {
		out[i] = e[p[i]]
	}
	return out
}
