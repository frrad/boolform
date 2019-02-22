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

func (p *Problem) NewBitVectConst(val []bool) *BitVect {
	z := make([]*Bool, len(val))
	for i := 0; i < len(val); i++ {
		z[i] = p.NewBoolConst(val[i])
	}
	ans := BitVect(z)
	return &ans
}

func (a *BitVect) termwise(f func(s, t *Bool) *Bool, b *BitVect) *BitVect {
	x, y := []*Bool(*a), []*Bool(*b)

	if len(x) != len(y) {
		panic("")
	}

	n := len(x)
	z := make([]*Bool, n)
	for i := 0; i < n; i++ {
		z[i] = f(x[i], y[i])
	}
	ans := BitVect(z)
	return &ans
}

func (a *BitVect) Eq(b *BitVect) *Bool {
	z := a.termwise((*Bool).Eq, b)

	zz := []*Bool(*z)
	if len(zz) == 0 {
		return zz[0]
	}

	return zz[0].And(zz[1:]...)
}

func (a *BitVect) Concat(b *BitVect) *BitVect {
	x, y := []*Bool(*a), []*Bool(*b)
	z := BitVect(append(x, y...))
	return &z
}

func (a *BitVect) And(b *BitVect) *BitVect {
	g := func(a, b *Bool) *Bool { return a.And(b) }
	return a.termwise(g, b)
}

func (a *BitVect) Or(b *BitVect) *BitVect {
	g := func(a, b *Bool) *Bool { return a.Or(b) }
	return a.termwise(g, b)
}

func (a *BitVect) SolVal() []bool {
	x := []*Bool(*a)

	ans := make([]bool, len(x))
	for i := 0; i < len(x); i++ {
		ans[i] = x[i].SolVal()
	}

	return ans
}
