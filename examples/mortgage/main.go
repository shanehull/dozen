package main

import (
	"fmt"

	"github.com/shanehull/dozen/engine"
)

func main() {
	// $300K at 6% APR over 30 years, paid monthly.
	pmt := engine.PMT(0.06/12, 360, 300000, 0, engine.End)
	fmt.Printf("$300K @ 6%% / 30yr → $%.2f/mo\n", -pmt)
}
