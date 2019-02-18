package bfgosat

import (
	"github.com/frrad/boolform/bf"
	gosat "github.com/mitchellh/go-sat"
	"github.com/mitchellh/go-sat/cnf"
)

type Problem struct {
	GSForm *cnf.Formula
	Lookup map[int]string
}

// Export takes a bf formula and returns a lightly wrapped go-sat formula.
func Export(f bf.Formula) Problem {
	c := bf.AsCNF(f)

	form := cnf.NewFormulaFromInts(c.Clauses)

	return Problem{
		GSForm: &form,
		Lookup: c.Lookup(),
	}
}

// Solve solves the given formula. If it is satisfiable, the function returns a
// model, associating each variable name with its binding. Else, the function
// returns nil.
func Solve(f bf.Formula) map[string]bool {
	pb := Export(f)

	s := gosat.New()
	s.AddFormula(*pb.GSForm)
	sat := s.Solve()
	if !sat {
		return nil
	}

	solution := s.Assignments()

	vars := map[string]bool{}
	for idx, name := range pb.Lookup {
		vars[name] = solution[idx]
	}

	return vars
}
