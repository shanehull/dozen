package main

import (
	"fmt"

	"github.com/shanehull/dozen/engine"
)

func main() {
	// Initial outlay followed by four years of returns, discounted at 10%.
	npv := engine.NPV(0.10, -100000, 30000, 40000, 50000, 60000)
	irr, _ := engine.IRR(-100000, 30000, 40000, 50000, 60000)
	fmt.Printf("NPV @ 10%%: $%.2f\n", npv)
	fmt.Printf("IRR: %.2f%%\n", irr*100)
}
