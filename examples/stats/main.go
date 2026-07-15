package main

import (
	"fmt"

	"github.com/shanehull/dozen/engine"
)

func main() {
	c := engine.New()
	for i, r := range []float64{10, 15, 22, 18, 25} {
		c.X, c.Y = r, float64(i+1)
		c.Step("Σ+", 0, "Σ+")
	}
	c.Step("x̄", 0, "x̄")
	fmt.Printf("Mean: %.1f\n", c.X)
	c.Step("s", 0, "s")
	fmt.Printf("StdDev: %.2f\n", c.X)
}
