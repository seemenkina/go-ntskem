package go_ntskem

func Decapsulation(a, h, c []uint16, l, r, t uint16) {
		//Build matrice H
	/*var H = make ([][]uint16, t*2)
	var i uint16 = 0
	var j uint16 = 0
	for i=0 ; i < t*2; i++ {
		H[i] = make([]uint16, l+r)
		for j= 0; j < l+r; j++ {
			H[i][j] = IntPow(a[j], i)*h[j] //in Fields!
			print(" ", H[i][j])
		}
		println()
	}*/

	//Build Transpose matrix.
	var H = make ([][]uint16, l+r)
	var i uint16 = 0
	var j uint16 = 0

	//OPERATIONS IN FIELDS!
	for i = 0 ; i < l+r; i++ {
		H[i] = make([]uint16, 2*t)
		for j = 0; j < 2*t; j++ {
			H[i][j] = IntPow(a[i], j)*h[i] //in Fields!
		}
	}
	
	//Compute 2t syndromes
	var s = PolyOnMatriceMult(c, H)
}
