package ff

import "github.com/seemenkina/go-ntskem/poly"

type FF struct {
	M     int
	Basis []uint16
}

func (ff2 *FF) New(m int) {
	ff2.M = m
	basis := make([]uint16, m)
	for i := 0; i < m; i++ {
		basis[m-i-1] = 1 << i
	}
	ff2.Basis = basis
}

func (ff2 *FF) Add(a, b uint16) uint16 {
	return a ^ b
}

func (ff2 *FF) Mul(a, b uint16) uint16 {
	var buf uint32

	buf = uint32(a * (b & 1))
	buf ^= uint32(a * (b & 0x0002))
	buf ^= uint32(a * (b & 0x0004))
	buf ^= uint32(a * (b & 0x0008))
	buf ^= uint32(a * (b & 0x0010))
	buf ^= uint32(a * (b & 0x0020))
	buf ^= uint32(a * (b & 0x0040))
	buf ^= uint32(a * (b & 0x0080))
	buf ^= uint32(a * (b & 0x0100))
	buf ^= uint32(a * (b & 0x0200))
	buf ^= uint32(a * (b & 0x0400))
	buf ^= uint32(a * (b & 0x0800))

	return reduce(buf)
}

func (ff2 *FF) Sqr(a uint16) uint16 {
	buf := uint32(a)
	buf = (buf | (buf << 8)) & 0x00FF00FF
	buf = (buf | (buf << 4)) & 0x0F0F0F0F
	buf = (buf | (buf << 2)) & 0x33333333
	buf = (buf | (buf << 1)) & 0x55555555

	return reduce(buf)
}

func (ff2 *FF) Inv(a uint16) uint16 {
	var a3, a15, b uint16

	a3 = ff2.Sqr(a)     /* a^2 */
	a3 = ff2.Mul(a3, a) /* a^3 */

	a15 = ff2.Sqr(a3)      /* a^6 */
	a15 = ff2.Sqr(a15)     /* a^12 */
	a15 = ff2.Mul(a15, a3) /* a^15 */

	b = ff2.Sqr(a15)    /* a^30 */
	b = ff2.Sqr(b)      /* a^60 */
	b = ff2.Sqr(b)      /* a^120 */
	b = ff2.Sqr(b)      /* a^240 */
	b = ff2.Mul(b, a15) /* a^255 */

	b = ff2.Sqr(b)     /* a^510 */
	b = ff2.Sqr(b)     /* a^1020 */
	b = ff2.Mul(b, a3) /* a^1023 */

	b = ff2.Sqr(b)    /* a^2046 */
	b = ff2.Mul(b, a) /* a^2047 */

	return ff2.Sqr(b) /* a^4094 */
}

/* GF(2^12), generated by f(x) = x^12 + x^3 + 1 */
func reduce(a uint32) uint16 {
	var buf uint32

	buf = a & 0x7F0000
	a ^= buf >> 9
	a ^= buf >> 12

	buf = a & 0x00F000
	a ^= buf >> 9
	a ^= buf >> 12

	return uint16(a & 0xFFF)
}

// Check the validity of Goppa polynomial
// 1 - g0 != 0
// 2 - Goppa polynomial  has no roots in F{2^m}. Check by Additive FFT
// 3 - Goppa polynomial  has no repeated roots in any extension 􏰀field.  Check by GCD(G)
func (ff2 *FF) CheckGoppaPoly(g *poly.Polynomial) bool {

	if g.Pol[0] == 0 {
		return false
	}
	if !ff2.checkFft(g) {
		return false
	}

	dx := ff2.Derivative(g)
	if dx != nil {
		gcd := ff2.GCD(g, dx)
		if gcd == nil || gcd.GetDegree() < 1 {
			return false
		}
	}

	return true
}

func (ff2 *FF) Derivative(g *poly.Polynomial) *poly.Polynomial {
	dx := poly.Polynomial{}
	dx.New(1 << ff2.M)

	if dx.Size() < g.Size()-1 {
		return nil
	}

	for i := 0; i < g.GetDegree(); i++ {
		dx.Pol[i] = 0
		if (i & 1) == 0 {
			dx.Pol[i] = g.Pol[i+1]
		}
	}
	dx.SetDegree(g.GetDegree() - 1)
	for i := 0; i < g.GetDegree(); i++ {
		if dx.Pol[g.GetDegree()-i-1] == 0 {
			break
		}
		dx.SetDegree(dx.GetDegree() - 1)
	}
	return &dx
}

func (ff2 *FF) GCD(f, g *poly.Polynomial) *poly.Polynomial {
	if f == nil {
		return g
	}

	if g == nil {
		return f
	}

	if g.GetDegree() < f.GetDegree() {
		g, f = f, g
	}

	return ff2.GCD(f, g.ModuloReduce(f, ff2))
}

// Return true, if Goppa polynomial  has roots in  F{2^m}
func (ff2 *FF) checkFft(pol *poly.Polynomial) bool {
	w := ff2.AdaptiveFft(pol)

	for i := 0; i < ff2.M; i++ {
		if w[i] == 0 {
			return false
		}
	}
	return true
}

func (ff2 *FF) AdaptiveFft(pol *poly.Polynomial) []uint16 {

	return nil
}

func TaylorExpansion() {

}
