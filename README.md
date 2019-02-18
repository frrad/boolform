[![](https://godoc.org/github.com/frrad/boolform?status.svg)](http://godoc.org/github.com/frrad/boolform)



The package `bf` (forked from `crillab/gophersat/bf`) allows for the
specification of boolean formula and the translation of specified formulas into
CNF.

The package `smt` lets the user specify constraints in a friendly form. There is
also support for some higher-level theories built out of boolean variables.

There are also a helper packages `bfgophersat`, `bfgini`, `bfgosat` to take a
formula created by `bf` and change it into problem instances for these solvers.

As a rule each of these helper packages has a `Solve` function if you are just
interested in a solution, and a `Export` function that translates a given
problem into a formula / solver instance which you can interact with in more
interesting ways.


Solving a simple problem with all three solvers:
``` golang
package bf

import (
	"fmt"

	"github.com/frrad/boolform/bf"
	"github.com/frrad/boolform/bfgini"
	"github.com/frrad/boolform/bfgophersat"
	"github.com/frrad/boolform/bfgosat"
)

func main() {
	x := bf.Var("x")
	y := bf.Var("y")
	z := bf.Var("z")
	f := bf.And(bf.And(x, y), bf.Not(z))

	fmt.Println(bfgophersat.Solve(f))
	fmt.Println(bfgosat.Solve(f))
	fmt.Println(bfgini.Solve(f))
}
```

Solving a bitvector problem with gosat:
``` golang
package main

import (
	"fmt"

	solver "github.com/frrad/boolform/bfgosat"
	"github.com/frrad/boolform/smt"
)

func main() {
	prob := smt.NewProb()
	a := prob.NewBitVectConst([]bool{true, false, true, false})
	b := prob.NewBitVectConst([]bool{true, true, false, false})

	x := a.Or(b)
	y := a.And(b)
	prob.Solve(solver.Solve)

	fmt.Println(x.SolVal())
	fmt.Println(y.SolVal())
}
```
