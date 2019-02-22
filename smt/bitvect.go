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

// Neq or "not equal to" returns a bool which is true iff a is not equal to b.
func (a *BitVect) Neq(b *BitVect) *Bool {
	return a.Eq(b).Not()
}

// Lt returns a bool which is true iff a is less than b
func (a *BitVect) Lt(b *BitVect) *Bool {
	x, y := []*Bool(*a), []*Bool(*b)
	if len(x) != len(y) {
		panic("nope")
	}

	head := y[0].And(x[0].Not())
	if len(x) == 1 {
		return head
	}

	xtail := BitVect(x[1:])
	ytail := BitVect(y[1:])

	return head.Or((x[0].Eq(y[0])).And(xtail.Lt(&ytail)))
}

// Gt returns a bool which is true iff a greater than b
func (a *BitVect) Gt(b *BitVect) *Bool {
	return b.Lt(a)
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
