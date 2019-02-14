package bf

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
)

// A Formula is any kind of boolean formula, not necessarily in CNF.
type Formula interface {
	nnf() Formula
	String() string
	Eval(model map[string]bool) bool
}

// Dimacs writes the DIMACS CNF version of the formula on w.
// It is useful so as to feed it to any SAT solver.
// The original names of variables is associated with their DIMACS integer counterparts
// in comments, between the prolog and the set of clauses.
// For instance, if the variable "a" is associated with the index 1, there will be a comment line
// "c a=1".
func Dimacs(f Formula, w io.Writer) error {
	cnf := AsCNF(f)
	nbVars := len(cnf.vars.all)
	nbClauses := len(cnf.Clauses)
	prefix := fmt.Sprintf("p cnf %d %d\n", nbVars, nbClauses)
	if _, err := io.WriteString(w, prefix); err != nil {
		return fmt.Errorf("could not write DIMACS output: %v", err)
	}
	var pbVars []string
	for v := range cnf.vars.Pb {
		if !v.dummy {
			pbVars = append(pbVars, v.Name)
		}
	}
	sort.Sort(sort.StringSlice(pbVars))
	for _, v := range pbVars {
		idx := cnf.vars.Pb[pbVar(v)]
		line := fmt.Sprintf("c %s=%d\n", v, idx)
		if _, err := io.WriteString(w, line); err != nil {
			return fmt.Errorf("could not write DIMACS output: %v", err)
		}
	}
	for _, clause := range cnf.Clauses {
		strClause := make([]string, len(clause))
		for i, lit := range clause {
			strClause[i] = strconv.Itoa(lit)
		}
		line := fmt.Sprintf("%s 0\n", strings.Join(strClause, " "))
		if _, err := io.WriteString(w, line); err != nil {
			return fmt.Errorf("could not write DIMACS output: %v", err)
		}
	}
	return nil
}

// The "true" constant.
type trueConst struct{}

// True is the constant denoting a tautology.
var True Formula = trueConst{}

func (t trueConst) nnf() Formula                    { return t }
func (t trueConst) String() string                  { return "⊤" }
func (t trueConst) Eval(model map[string]bool) bool { return true }

// The "false" constant.
type falseConst struct{}

// False is the constant denoting a contradiction.
var False Formula = falseConst{}

func (f falseConst) nnf() Formula                    { return f }
func (f falseConst) String() string                  { return "⊥" }
func (f falseConst) Eval(model map[string]bool) bool { return false }

// Var generates a named boolean variable in a formula.
func Var(name string) Formula {
	return pbVar(name)
}

func pbVar(name string) variable {
	return variable{Name: name, dummy: false}
}

func dummyVar(name string) variable {
	return variable{Name: name, dummy: true}
}

type variable struct {
	Name  string
	dummy bool
}

func (v variable) nnf() Formula {
	return lit{signed: false, v: v}
}

func (v variable) String() string {
	return v.Name
}

func (v variable) Eval(model map[string]bool) bool {
	b, ok := model[v.Name]
	if !ok {
		panic(fmt.Errorf("Model lacks binding for variable %s", v.Name))
	}
	return b
}

type lit struct {
	v      variable
	signed bool
}

func (l lit) nnf() Formula {
	return l
}

func (l lit) String() string {
	if l.signed {
		return "not(" + l.v.Name + ")"
	}
	return l.v.Name
}

func (l lit) Eval(model map[string]bool) bool {
	b := l.v.Eval(model)
	if l.signed {
		return !b
	}
	return b
}

// Not represents a negation. It negates the given subformula.
func Not(f Formula) Formula {
	return not{f}
}

type not [1]Formula

func (n not) nnf() Formula {
	switch f := n[0].(type) {
	case variable:
		l := f.nnf().(lit)
		l.signed = true
		return l
	case lit:
		f.signed = !f.signed
		return f
	case not:
		return f[0].nnf()
	case and:
		subs := make([]Formula, len(f))
		for i, sub := range f {
			subs[i] = not{sub}.nnf()
		}
		return or(subs).nnf()
	case or:
		subs := make([]Formula, len(f))
		for i, sub := range f {
			subs[i] = not{sub}.nnf()
		}
		return and(subs).nnf()
	case trueConst:
		return False
	case falseConst:
		return True
	default:
		panic("invalid formula type")
	}
}

func (n not) String() string {
	return "not(" + n[0].String() + ")"
}

func (n not) Eval(model map[string]bool) bool {
	return !n[0].Eval(model)
}

// And generates a conjunction of subformulas.
func And(subs ...Formula) Formula {
	return and(subs)
}

type and []Formula

func (a and) nnf() Formula {
	var res and
	for _, s := range a {
		nnf := s.nnf()
		switch nnf := nnf.(type) {
		case and: // Simplify: "and"s in the "and" get to the higher level
			res = append(res, nnf...)
		case trueConst: // True is ignored
		case falseConst:
			return False
		default:
			res = append(res, nnf)
		}
	}
	if len(res) == 1 {
		return res[0]
	}
	if len(res) == 0 {
		return False
	}
	return res
}

func (a and) String() string {
	strs := make([]string, len(a))
	for i, f := range a {
		strs[i] = f.String()
	}
	return "and(" + strings.Join(strs, ", ") + ")"
}

func (a and) Eval(model map[string]bool) (res bool) {
	for i, s := range a {
		b := s.Eval(model)
		if i == 0 {
			res = b
		} else {
			res = res && b
		}
	}
	return
}

// Or generates a disjunction of subformulas.
func Or(subs ...Formula) Formula {
	return or(subs)
}

type or []Formula

func (o or) nnf() Formula {
	var res or
	for _, s := range o {
		nnf := s.nnf()
		switch nnf := nnf.(type) {
		case or: // Simplify: "or"s in the "or" get to the higher level
			res = append(res, nnf...)
		case falseConst: // False is ignored
		case trueConst:
			return True
		default:
			res = append(res, nnf)
		}
	}
	if len(res) == 1 {
		return res[0]
	}
	if len(res) == 0 {
		return True
	}
	return res
}

func (o or) String() string {
	strs := make([]string, len(o))
	for i, f := range o {
		strs[i] = f.String()
	}
	return "or(" + strings.Join(strs, ", ") + ")"
}

func (o or) Eval(model map[string]bool) (res bool) {
	for i, s := range o {
		b := s.Eval(model)
		if i == 0 {
			res = b
		} else {
			res = res || b
		}
	}
	return
}

// Implies indicates a subformula implies another one.
func Implies(f1, f2 Formula) Formula {
	return or{not{f1}, f2}
}

// Eq indicates a subformula is equivalent to another one.
func Eq(f1, f2 Formula) Formula {
	return and{or{not{f1}, f2}, or{f1, not{f2}}}
}

// Xor indicates exactly one of the two given subformulas is true.
func Xor(f1, f2 Formula) Formula {
	return and{or{not{f1}, not{f2}}, or{f1, f2}}
}

// Unique indicates exactly one of the given variables must be true.
// It might create dummy variables to reduce the number of generated clauses.
func Unique(vars ...Formula) Formula {
	vars2 := make([]variable, len(vars))
	eq := make([]Formula, len(vars))

	// improvment: don't do this when the input is already variable
	for i, v := range vars {
		vars2[i] = dummyVar(fmt.Sprintf("u-%d", i)) // fix: breaks if called multiple times
		eq[i] = Eq(vars2[i], v)
	}

	return And(append(eq, uniqueRec(vars2...))...)
}

// uniqueSmall generates clauses indicating exactly one of the given variables is true.
// It is suitable when the number of variables is small (typically, <= 4).
func uniqueSmall(vars ...variable) Formula {
	res := make([]Formula, 1, 1+(len(vars)*len(vars)-1)/2)
	varsAsForms := make([]Formula, len(vars))
	for i, v := range vars {
		varsAsForms[i] = v
	}
	res[0] = Or(varsAsForms...)
	for i := 0; i < len(vars)-1; i++ {
		for j := i + 1; j < len(vars); j++ {
			res = append(res, Or(Not(varsAsForms[i]), Not(varsAsForms[j])))
		}
	}
	return And(res...)
}

func uniqueRec(vars ...variable) Formula {
	nbVars := len(vars)
	if nbVars <= 4 {
		return uniqueSmall(vars...)
	}
	sqrt := math.Sqrt(float64(nbVars))
	nbLines := int(sqrt + 0.5)
	lines := make([]variable, nbLines)
	linesF := make([][]Formula, nbLines)
	allNames := make([]string, len(vars))
	for i := range vars {
		allNames[i] = vars[i].Name
	}
	fullName := strings.Join(allNames, "-")
	for i := range lines {
		lines[i] = dummyVar(fmt.Sprintf("line-%d-%s", i, fullName))
		linesF[i] = []Formula{}
	}
	nbCols := int(math.Ceil(sqrt))
	cols := make([]variable, nbCols)
	colsF := make([][]Formula, nbCols)
	for i := range cols {
		cols[i] = dummyVar(fmt.Sprintf("col-%d-%s", i, fullName))
		colsF[i] = []Formula{}
	}
	res := make([]Formula, 0, 2*nbVars+1)
	for i, v := range vars {
		linesF[i/nbCols] = append(linesF[i/nbCols], v)
		colsF[i%nbCols] = append(colsF[i%nbCols], v)
	}
	for i := range lines {
		res = append(res, Eq(lines[i], Or(linesF[i]...)))
	}
	for i := range cols {
		res = append(res, Eq(cols[i], Or(colsF[i]...)))
	}

	res = append(res, uniqueRec(lines...))
	res = append(res, uniqueRec(cols...))
	return And(res...)
}

// vars associate variable names with numeric indices.
type vars struct {
	all map[variable]int // all vars, including those created when converting the formula
	Pb  map[variable]int // Only the vars that appeared orinigally in the problem
}

// litValue returns the int value associated with the given problem var.
// If the var was not referenced yet, it is created first.
func (vars *vars) litValue(l lit) int {
	val, ok := vars.all[l.v]
	if !ok {
		val = len(vars.all) + 1
		vars.all[l.v] = val
		if !l.v.dummy {
			vars.Pb[l.v] = val
		}

	}
	if l.signed {
		return -val
	}
	return val
}

// Dummy creates a dummy variable and returns its associated index.
func (vars *vars) dummy() int {
	val := len(vars.all) + 1
	vars.all[dummyVar(fmt.Sprintf("dummy-%d", val))] = val
	return val
}

// A CNF is the representation of a boolean formula as a conjunction of
// disjunction.  It can be solved by a SAT solver.
type CNF struct {
	vars    vars
	Clauses [][]int
}

// Lookup returns a map recording the correspondence between variable numbers in
// the CNF representation and original variable names.
func (c *CNF) Lookup() map[int]string {
	lookup := map[int]string{}
	for v, ix := range c.vars.Pb {
		lookup[ix-1] = v.Name
	}
	return lookup
}

// AsCNF returns a CNF representation of the given formula.
func AsCNF(f Formula) *CNF {
	vars := vars{all: make(map[variable]int), Pb: make(map[variable]int)}
	clauses := cnfRec(f.nnf(), &vars)
	return &CNF{vars: vars, Clauses: clauses}
}

// transforms the f NNF formula into a CNF formula.  nbDummies is the current
// number of dummy variables created.
//
// Note: code should be improved, there are a few useless allocs/deallocs here
// and there.
func cnfRec(f Formula, vars *vars) [][]int {
	switch f := f.(type) {
	case lit:
		return [][]int{{vars.litValue(f)}}
	case and:
		var res [][]int
		for _, sub := range f {
			res = append(res, cnfRec(sub, vars)...)
		}
		return res
	case or:
		var res [][]int
		var lits []int
		for _, sub := range f {
			switch sub := sub.(type) {
			case lit:
				lits = append(lits, vars.litValue(sub))
			case and:
				d := vars.dummy()
				lits = append(lits, d)
				for _, sub2 := range sub {
					cnf := cnfRec(sub2, vars)
					cnf[0] = append(cnf[0], -d)
					res = append(res, cnf...)
				}
			default:
				panic("unexpected or in or")
			}
		}
		res = append(res, lits)
		return res
	case trueConst: // True clauses are ignored
		return [][]int{}
	case falseConst: // TODO: improve this. This should simply be declared to make the problem UNSAT.
		return [][]int{{}}
	default:
		panic("invalid NNF formula")
	}
}
