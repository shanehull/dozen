package main

import (
	"fmt"

	"github.com/shanehull/dozen/engine"
)

func main() {
	e := engine.New()
	e.SetN(20)      // 20 years
	e.SetI(8)       // 8% per period
	e.SetPV(-10000) // $10,000 invested (outflow)
	e.SolveFV()
	fmt.Printf("$10K @ 8%% / 20yr → $%.2f\n", e.X)
}
