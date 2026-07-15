package engine_test

import (
	"fmt"

	"github.com/shanehull/dozen/engine"
)

// The functional API needs no calculator, stack, or keystrokes — just call the
// function with the values you already have.

func ExamplePMT() {
	// A $300,000 mortgage at 6% APR over 30 years, paid monthly.
	pmt := engine.PMT(0.06/12, 360, 300000, 0, engine.End)
	fmt.Printf("Monthly payment: $%.2f\n", -pmt)
	// Output: Monthly payment: $1798.65
}

func ExampleFV() {
	// $10,000 invested for 20 years at 8% (outflow now, so PV is negative).
	fv := engine.FV(0.08, 20, -10000, 0, engine.End)
	fmt.Printf("Future value: $%.2f\n", fv)
	// Output: Future value: $46609.57
}

func ExampleRate() {
	// What rate turns $1,000 today into $2,000 after 12 periods?
	r, _ := engine.Rate(12, 0, -1000, 2000, engine.End)
	fmt.Printf("Rate: %.3f%%\n", r*100)
	// Output: Rate: 5.946%
}

func ExampleNPV() {
	npv := engine.NPV(0.10, -100000, 30000, 40000, 50000, 60000)
	fmt.Printf("NPV: $%.2f\n", npv)
	// Output: NPV: $38877.13
}

func ExampleIRR() {
	irr, _ := engine.IRR(-100000, 40000, 50000, 30000)
	fmt.Printf("IRR: %.2f%%\n", irr*100)
	// Output: IRR: 10.13%
}

func ExampleDepSL() {
	// $50,000 asset, $5,000 salvage, 10-year life, first year.
	dep, remaining := engine.DepSL(50000, 5000, 10, 1)
	fmt.Printf("Year 1 depreciation: $%.2f, remaining: $%.2f\n", dep, remaining)
	// Output: Year 1 depreciation: $4500.00, remaining: $40500.00
}
