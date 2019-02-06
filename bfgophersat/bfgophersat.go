package bfgophersat

import (
	"github.com/crillab/gophersat/solver"
	"github.com/frrad/boolform/bf"
)

type Problem struct {
	GSPRob *solver.Problem
	Lookup map[int]string
}

func Export(f bf.Formula) Problem {
	cnf := bf.AsCnf(f)

	return Problem{
		GSPRob: solver.ParseSlice(cnf.Clauses),
		Lookup: cnf.Lookup(),
	}
}

// solve solves the given formula.
// cnf is given to gophersat.
// If it is satisfiable, the function returns a model, associating each variable name with its binding.
// Else, the function returns nil.
func Solve(f bf.Formula) map[string]bool {
	pb := Export(f)

	s := solver.New(pb.GSPRob)
	if s.Solve() != solver.Sat {
		return nil
	}
	m := s.Model()
	vars := make(map[string]bool)

	for idx, name := range pb.Lookup {
		vars[name] = m[idx]
	}

	return vars
}
