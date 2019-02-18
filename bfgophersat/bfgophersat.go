package bfgophersat

import (
	"github.com/crillab/gophersat/solver"
	"github.com/frrad/boolform/bf"
)

type Problem struct {
	GSProb *solver.Problem
	Lookup map[int]string
}

func Export(f bf.Formula) Problem {
	cnf := bf.AsCNF(f)

	return Problem{
		GSProb: solver.ParseSlice(cnf.Clauses),
		Lookup: cnf.Lookup(),
	}
}

// Solve solves the given formula. CNF is given to gophersat. If it is
// satisfiable, the function returns a model, associating each variable name
// with its binding.  Else, the function returns nil.
func Solve(f bf.Formula) map[string]bool {
	pb := Export(f)

	s := solver.New(pb.GSProb)
	if s.Solve() != solver.Sat {
		return nil
	}
	m := s.Model()
	vars := make(map[string]bool)

	for idx, name := range pb.Lookup {
		vars[name] = m[idx-1]
	}

	return vars
}
