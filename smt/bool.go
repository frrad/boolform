package smt

import (
	"fmt"

	"github.com/frrad/boolform/bf"
)

type Bool struct {
	wrapped bf.Formula

	// name is a unique name for this var. It should correspond to the name of
	// the underlying bf.Formula.
	name string

	// prob is a pointer to the problem instance that this var belongs to.
	prob *Problem
}

func (p *Problem) NewBool() *Bool {
	n := p.nextName()

	underlying := bf.Var(n)

	ans := &Bool{
		name:    n,
		wrapped: underlying,
		prob:    p,
	}

	if ans.wrapped.String() != ans.name {
		panic("skew")
	}
	return ans
}

func (p *Problem) NewBoolConst(val bool) *Bool {
	var underlying bf.Formula
	if val {
		underlying = bf.True
	} else {
		underlying = bf.False
	}

	ans := p.NewBool()
	p.assertFormula(bf.Eq(ans.wrapped, underlying))
	/////
	return ans
}

func (p *Problem) wrap(val bf.Formula) *Bool {
	x := p.NewBool()
	p.assertFormula(bf.Eq(x.wrapped, val))
	return x
}

func (a *Bool) SolVal() bool {
	if !a.prob.solved {
		panic("tried to access solval of unsolved problem")
	}
	if _, ok := a.prob.sol[a.name]; !ok {
		panic(fmt.Sprintf("can't find solval for var %+v", a))
	}

	return a.prob.sol[a.name]
}

func (a *Bool) Unique(rest ...*Bool) *Bool {
	unwrap := make([]bf.Formula, len(rest)+1)
	unwrap[0] = a.wrapped
	for i := 0; i < len(rest); i++ {
		unwrap[i+1] = rest[i].wrapped
	}

	underlying := bf.Unique(unwrap...)
	return a.prob.wrap(underlying)
}

func (a *Bool) Eq(b *Bool) *Bool {
	underlying := bf.Eq(a.wrapped, b.wrapped)
	return a.prob.wrap(underlying)
}

func (a *Bool) And(rest ...*Bool) *Bool {
	unwrap := make([]bf.Formula, len(rest)+1)
	unwrap[0] = a.wrapped
	for i := 0; i < len(rest); i++ {
		unwrap[i+1] = rest[i].wrapped
	}

	underlying := bf.And(unwrap...)
	return a.prob.wrap(underlying)
}

func (a *Bool) Or(rest ...*Bool) *Bool {
	fmt.Println('.')
	unwrap := make([]bf.Formula, len(rest)+1)
	fmt.Println(len(unwrap))
	unwrap[0] = a.wrapped
	for i := 0; i < len(rest); i++ {
		unwrap[i+1] = rest[i].wrapped
	}

	underlying := bf.Or(unwrap...)
	return a.prob.wrap(underlying)
}
