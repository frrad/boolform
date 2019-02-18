package smt

import (
	"fmt"

	"github.com/frrad/boolform/bf"
)

type Bool struct {
	wrapped bf.Formula

	dummy bool

	// name is a unique name for this var. It is only guaranteed to correspond
	// to the name of the underlying bf.Formula for non-dummy variables
	name string

	// prob is a pointer to the problem instance that this var belongs to.
	prob *Problem
}

func (p *Problem) wrap(n string, dum bool, underlying bf.Formula) *Bool {
	ans := &Bool{
		name:    n,
		dummy:   dum,
		wrapped: underlying,
		prob:    p,
	}

	if !ans.dummy && ans.wrapped.String() != ans.name {
		panic("skew")
	}

	return ans
}

func (p *Problem) NewBool() *Bool {
	name := p.nextName()

	underlying := bf.Var(name)
	return p.wrap(name, false, underlying)
}

func (p *Problem) NewBoolConst(val bool) *Bool {
	name := p.nextName()
	var underlying bf.Formula

	if val {
		underlying = bf.True
	} else {
		underlying = bf.False
	}

	return p.wrap(name, true, underlying)
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
	return a.prob.wrap(a.prob.nextName(), true, underlying)
}

func (a *Bool) Eq(b *Bool) *Bool {
	underlying := bf.Eq(a.wrapped, b.wrapped)
	return a.prob.wrap(a.prob.nextName(), true, underlying)
}
