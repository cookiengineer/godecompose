package main

import (
	"fmt"
	"math/big"
)

func bigExercise() {
	x := new(big.Int)
	x.SetString("12345678901234567890", 10)
	_ = x.String()

	y := new(big.Int)
	y.SetString("98765432109876543210", 10)

	z := new(big.Int)
	z.Add(x, y)
	z.Mul(x, y)
	z.Div(y, x)
	z.Mod(y, x)

	r := new(big.Rat)
	r.SetString("355/113")

	f := new(big.Float)
	f.SetString("3.141592653589793")
}

func main() {
	bigExercise()
	fmt.Println("math/big e2e test done")
}
