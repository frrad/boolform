boolform is a fork of `crillab/gophersat/bf`

The main package is `bf` which allows for the specification of boolean formula
and the translation of specified formulas into CNF.

There are also a helper packages `bfgophersat`, `bfgini`, `bfgosat` to take a
formula created by bf and change it into problem instances for these solvers.

As a rule each of these helper packages has a `Solve` function if you are just
interested in a solution, and a `Export` function that translates a given
problem into a formula / solver instance which you can interact with in more
interesting ways.


Solving a simple problem with all three solvers:
``` golang
package main

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

