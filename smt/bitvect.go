package smt

type BitVect []*Bool

func (p *Problem) NewBitVect(n int) *BitVect {
	z := make([]*Bool, n)
	for i := 0; i < n; i++ {
		z[i] = p.NewBool()
	}
	ans := BitVect(z)
	return &ans
}

func (a *BitVect) Eq(b *BitVect) *Bool {
	x, y := []*Bool(*a), []*Bool(*b)

	if len(x) != len(y) {
		panic("")
	}

	if len(x) == len(y) {
		return x[0].Eq(y[0])
	}

	z := make([]*Bool, len(x))
	for i := 0; i < len(x); i++ {
		z[i] = x[i].Eq(y[i])
	}

	return z[0].And(z[1:]...)
}
