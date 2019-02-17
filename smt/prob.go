package smt

import (
	"github.com/frrad/boolform/bf"
)

// Problem is a thin wrapper around Formula to allow for gradually building up a
// complex formula out of many separate assertions.
type Problem struct {
	form bf.Formula
}

// NewProb creates a new Problem instance to start building from.
func NewProb() *Problem {
	return &Problem{
		form: bf.True,
	}
}

// Assert adds the requirement that the given assertion be true to the Problem
// instance.
func (p *Problem) Assert(ass bf.Formula) {
	p.form = bf.And(p.form, ass)
}

// AsFormula can be used to retrieve the Formula representation of the problem
// we've built up.
func (p *Problem) AsFormula() bf.Formula {
	return p.form
}
