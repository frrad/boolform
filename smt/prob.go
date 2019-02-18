package smt

import (
	"fmt"
	"sync"

	"github.com/frrad/boolform/bf"
)

// Problem is an SMT problem instance.
type Problem struct {
	assertions bf.Formula

	varNames map[string]bool

	nameIx int
	nameMx sync.RWMutex

	solved bool
	sol    map[string]bool
}

// NewProb creates a new Problem instance.
func NewProb() *Problem {
	return &Problem{
		assertions: bf.True,

		varNames: map[string]bool{},

		nameIx: 0,
		nameMx: sync.RWMutex{},

		solved: false,
	}
}

// Assert adds the given assertion to the Problem instance.
func (p *Problem) Assert(ass *Bool) {
	p.assertFormula(ass.wrapped)
}

// assertFormula adds an assertion of a bf.Formula.
func (p *Problem) assertFormula(ass bf.Formula) {
	p.assertions = bf.And(p.assertions, ass)
}

// AsFormula can be used to retrieve the Formula representation of the problem
// we've built up.
func (p *Problem) AsFormula() bf.Formula {
	return p.assertions
}

// Solve uses the provided backend to solve the problem instance. The results
// can be retrieved by calling SolVar on the variables of interest.
func (p *Problem) Solve(backend func(bf.Formula) map[string]bool) bool {
	sol := backend(p.assertions)
	if sol == nil {
		return false
	}

	p.solved = true
	p.sol = sol

	return true
}

func (p *Problem) nextName() string {
	p.nameMx.Lock()
	name := fmt.Sprintf("%d", p.nameIx)
	for p.varNames[name] {
		p.nameIx++
		name = fmt.Sprintf("%d", p.nameIx)
	}
	p.varNames[name] = true
	p.nameMx.Unlock()
	return name
}
