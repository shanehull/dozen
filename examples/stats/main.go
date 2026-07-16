package main

import (
	"fmt"

	"github.com/shanehull/dozen/engine"
)

func main() {
	c := engine.New()
	for i, r := range []float64{10, 15, 22, 18, 25} {
		c.X, c.Y = r, float64(i+1)
		c.StatAdd()
	}
	c.MeanX()
	fmt.Printf("Mean: %.1f\n", c.X)
	c.SDev()
	fmt.Printf("StdDev: %.2f\n", c.X)
}
