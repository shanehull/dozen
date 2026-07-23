package main

import (
	"fmt"

	"github.com/shanehull/dozen/engine"
)

func main() {
	e := engine.New()

	// $300K mortgage, 30 years, $1,798.65/mo → what APR?
	e.SetN(360)
	e.SetPV(300000)
	e.SetPMT(-1798.65)
	e.SolveI()
	fmt.Printf("Rate: %.2f%% APR\n", e.X*12)

	// $300K at 6% APR with $2,000/mo → how many months?
	e.SetI(0.5) // 6% APR ÷ 12
	e.SetPV(300000)
	e.SetPMT(-2000)
	e.SolveN()
	fmt.Printf("Term: %.0f months at $2,000/mo\n", e.X)

	// 6% APR over 30 years, $1,500/mo → what loan?
	e.SetN(360)
	e.SetI(0.5) // 6% APR ÷ 12
	e.SetPMT(-1500)
	e.SolvePV()
	fmt.Printf("Loan: $%.0f\n", e.X)
}
