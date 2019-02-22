package smt

// UInt8 is an 8-bit unsigned integer with LSB in index 0
type UInt8 []*Bool

func (p *Problem) NewUInt8() *UInt8 {
	z := p.NewBitVect(8)
	ans := UInt8(*z)
	return &ans
}

func (p *Problem) NewUInt8Const(val uint8) *UInt8 {
	z := make([]*Bool, 8)
	for i := 0; i < 8; i++ {
		z[7-i] = p.NewBoolConst((val & 1) == 1)
		val = val >> 1
	}
	ans := UInt8(z)
	return &ans
}

func (a *UInt8) SolVal() uint8 {
	x := []*Bool(*a)

	ans := uint8(0)
	for i := 0; i < 8; i++ {
		ans = ans << 1

		if x[i].SolVal() {
			ans += 1
		}
	}

	return ans
}

func (a *UInt8) Eq(b *UInt8) *Bool {
	x, y := BitVect(*a), BitVect(*b)
	return x.Eq(&y)
}

func (a *UInt8) Neq(b *UInt8) *Bool {
	x, y := BitVect(*a), BitVect(*b)
	return x.Neq(&y)
}

func (a *UInt8) Lt(b *UInt8) *Bool {
	x, y := BitVect(*a), BitVect(*b)
	return x.Lt(&y)
}

func (a *UInt8) Gt(b *UInt8) *Bool {
	x, y := BitVect(*a), BitVect(*b)
	return x.Gt(&y)
}

func (a *UInt8) Add(b *UInt8) *UInt8 {
	x, y := BitVect(*a), BitVect(*b)
	z := x.Add(&y)
	ans := UInt8(*z)
	return &ans
}
