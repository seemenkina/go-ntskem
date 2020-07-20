package go-ntskem


func Decapsulation (a, h, c []uint16, l,r,t uint16)  {
	/*ab :=a[:l]
	ac:=a[l:]

	hb := h[:l]
	hc:=h[l:]*/

	//Construct matrice H
	var H = make ([][]uint16, 2*t)
	var i uint16 = 0
	var j uint16 = 0

	//OPERATIONS IN FIELDS!
	for i=0 ; i < t*2; i++ {
		H[i] = make([]uint16, l+r)
		for j= 0; j < l+r; j++ {
			H[i][j] = IntPow(a[j], i)*h[j] //in Fields!
		}
	}
}
