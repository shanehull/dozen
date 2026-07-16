package main

import (
	"fmt"
	"math"

	"github.com/shanehull/dozen/engine"
)

func main() {
	// $300K mortgage, 30 years, monthly payment of $1,798.65.
	// What's the APR? Pass NaN for the unknown.
	i, _ := engine.SolveTVM(360, math.NaN(), 300000, -1798.65, 0, engine.End)
	fmt.Printf("Rate: %.2f%% APR\n", i*12*100)

	// How many months to pay off with $2,000/mo?
	n, _ := engine.SolveTVM(math.NaN(), 0.06/12, 300000, -2000, 0, engine.End)
	fmt.Printf("Term: %.0f months at $2,000/mo\n", n)

	// What loan can I afford at $1,500/mo over 30 years?
	pv, _ := engine.SolveTVM(360, 0.06/12, math.NaN(), -1500, 0, engine.End)
	fmt.Printf("Loan: $%.0f\n", pv)
}
