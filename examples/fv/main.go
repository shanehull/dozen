package main

import (
	"fmt"

	"github.com/shanehull/dozen/engine"
)

func main() {
	// $10K invested for 20 years at 8% (outflow now → negative PV).
	fv := engine.FV(0.08, 20, -10000, 0, engine.End)
	fmt.Printf("$10K @ 8%% / 20yr → $%.2f\n", fv)
}
