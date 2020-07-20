package go_ntskem

import (
"golang.org/x/crypto/sha3"
"math/rand"
"time"
)


//TODO: BitArrays Type
//	Decapsulation: Operations in fields
//	Change math/rand to crypto/rand
func main() {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	//Q gen
	var Q = make ([][]uint16, 3328)
	for i:=0 ; i < 3328; i++ {
		Q[i] = make([]uint16, 768)
		for j:= 0; j < 768; j++ {
			if j % 2 == 0 {
				Q[i][j] = 1
			} else {
				Q[i][j] = 0
			}
		}
	}

	//a* gen
	var a = make([]uint16, 968)
	for i:=0; i < 968; i++{
		a[i] = uint16(random.Intn(3))
	}
	println()
	//h* gen
	var h = make([]uint16, 968)
	for i:=0; i < 968; i++{
		h[i] = uint16(random.Intn(3))
	}
	var kr, c = Encapsulation(3328,256,4096,64, Q)
	//Decapsulation(a,h,c,256,712, 64)

	print(kr)
	print(c)
}

func randomVector (n uint16, t uint16) []uint16 {
	var e []uint16 = make([]uint16, n-t)
	//vector, n -length, t = weight of e
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

func Encapsulation (k uint16, l uint16, n uint16, t uint16, Q [][]uint16) ([]uint16, []uint16){
	//1. Random error vector
	e:= randomVector(n,t)

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
	cc := PolySum(PolyOnMatriceMult(m, Q), ec)

	//worng matrices mult!
	var co []uint16
	for i=0; i < l;i++{
		co=append(co,cb[i])
	}
	for i=0; i < uint16(len(cc));i++{
		co=append(co,cc[i])
	}
	/*println("cb len: %v", len(cb))
	println("cc len: %v", len(cc))
	println("co len: %v", len(co))*/

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

//Converts BitArray to ByteArray
func BitArrayToByteArray (IntArray []uint16) []byte{
	var ByteArray []byte = make ([]byte, len(IntArray)/8)
	for i:=0; i < len(IntArray); i++ {
		ByteArray[i/8] = ByteArray[i/8] | byte(IntArray[i])
		if i % 8 != 7 {
			ByteArray[i/8] <<= 1
		}
	}
	return ByteArray
}

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

//Polynom Sum
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

//wrong result!
func PolyOnMatriceMult (poly []uint16, matrice [][]uint16) []uint16{
	if len(poly) != len(matrice){
		println("err handling here")
	}
	result := make ([]uint16, len(poly))
	for i:=0; i < len(poly); i++{
		for j:=0; j < len(matrice[i]); j++{
			result[i] ^= poly[j] * matrice[i][j]
		}
	}
	return result
}

//No intPow in golang?
func IntPow (x, pow uint16) uint16{
	var i uint16
	var result uint16 = 1
	for i=0; i < pow; i++ {
		result *= x
	}
	return result
}
