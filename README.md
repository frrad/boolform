boolform is a fork of `crillab/gophersat/bf`

The main package is `bf` which allows for the specification of boolean formula
and the translation of specified formulas into CNF.

There is also a helper package `bfgophersat` to take a formula created by bf and
change it into a gophersat problem instance. 

In the future, it should be easy to write similar adapters for other SAT solvers.


``` golang
package main

import (
	"fmt"

	"github.com/frrad/boolform/bf"
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
}
```

