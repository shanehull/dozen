package main

import (
	"fmt"

	"github.com/shanehull/dozen/engine"
)

func main() {
	e := engine.New()
	e.SetN(360)     // 30 years × 12 months
	e.SetI(0.5)     // 6% APR ÷ 12 = 0.5% per month
	e.SetPV(300000) // $300,000 loan
	e.SetFV(0)      // paid off at end
	e.SolvePMT()
	fmt.Printf("$300K @ 6%% / 30yr → $%.2f/mo\n", -e.X)
}
