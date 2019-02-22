package smt

// Add returns the sum of two Bitvectors. The inputs are assumed to have the
// same length.
func (a *BitVect) Add(b *BitVect) *BitVect {
	x := []*Bool(*a)
	zero := x[0].prob.NewBoolConst(false)

	return add(a, b, zero)
}

// add returns a+b+c
func add(a, b *BitVect, c *Bool) *BitVect {
	x, y := []*Bool(*a), []*Bool(*b)

	if len(x) != len(y) {
		panic("don't do this")
	}

	sum, carry := fullAdder(x[len(x)-1], y[len(y)-1], c)

	lsb := BitVect([]*Bool{sum})

	if len(x) == 1 {
		return &lsb
	}

	xx := BitVect(x[:len(x)-1])
	yy := BitVect(y[:len(y)-1])
	front := add(&xx, &yy, carry)
	return front.Concat(&lsb)
}

// A fullAdder circuit:
// https://en.wikipedia.org/wiki/Adder_(electronics)#Full_adder
func fullAdder(a, b, c *Bool) (*Bool, *Bool) {
	one := a.Xor(b)
	sum := one.Xor(c)
	three := one.And(c)
	four := a.And(b)
	carry := three.Or(four)
	return sum, carry
}
