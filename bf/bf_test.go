package bf

import (
	"fmt"
	"os"
	"testing"
)

func TestString(t *testing.T) {
	f := And(Or(Var("a"), Not(Var("b"))), Not(Var("c")))
	const expected = "and(or(a, not(b)), not(c))"
	if f.String() != expected {
		t.Errorf("string representation of formula not as expected: wanted %q, got %q", expected, f.String())
	}
}

func ExampleDimacs() {
	f := Eq(And(Or(Var("a"), Not(Var("b"))), Not(Var("a"))), Var("b"))
	if err := Dimacs(f, os.Stdout); err != nil {
		fmt.Printf("Could not generate DIMACS file: %v", err)
	}
	// Output:
	// p cnf 4 6
	// c a=2
	// c b=3
	// -2 -1 0
	// 3 -1 0
	// 1 2 3 0
	// 2 -3 -4 0
	// -2 -4 0
	// 4 -3 0
}
