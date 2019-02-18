package bfgini

import (
	"github.com/frrad/boolform/bf"
	"github.com/irifrance/gini"
	"github.com/irifrance/gini/z"
)

type Problem struct {
	GiniSol *gini.Gini
	Lookup  map[int]string
}

// Export takes a bf formula and returns a lightly wrapped gini instance.
func Export(f bf.Formula) Problem {
	c := bf.AsCNF(f)

	g := gini.New()
	for _, clause := range c.Clauses {
		for _, term := range clause {
			g.Add(z.Dimacs2Lit(term))
		}
		g.Add(0)
	}

	return Problem{
		GiniSol: g,
		Lookup:  c.Lookup(),
	}
}

// Solve solves the given formula with gini. If it is satisfiable, the function
// returns a model, associating each variable name with its binding. Else, the
// function returns nil.
func Solve(f bf.Formula) map[string]bool {
	pb := Export(f)

	worked := pb.GiniSol.Solve()
	if worked != 1 {
		return nil
	}

	vars := map[string]bool{}
	for idx, name := range pb.Lookup {
		vars[name] = pb.GiniSol.Value(z.Dimacs2Lit(idx))
	}

	return vars
}
