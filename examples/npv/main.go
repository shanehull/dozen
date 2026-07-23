package main

import (
	"fmt"

	"github.com/shanehull/dozen/engine"
)

func main() {
	e := engine.New()
	e.FinCF0 = -100000
	e.FinCFj[0] = 30000
	e.FinCFj[1] = 40000
	e.FinCFj[2] = 50000
	e.FinCFj[3] = 60000
	e.FinCfCnt = 4

	e.X = 10
	e.NPV()
	fmt.Printf("NPV @ 10%%: $%.2f\n", e.X)

	e.IRR()
	fmt.Printf("IRR: %.2f%%\n", e.X)
}
