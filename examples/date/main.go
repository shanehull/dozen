package main

import (
	"fmt"

	"github.com/shanehull/dozen/engine"
)

func main() {
	c := engine.New()
	c.Flags.Dmy = false
	c.Y, c.X = 1.012025, 6.152025
	c.Step("DYS", 0, "DYS")
	fmt.Printf("Jan 1 → Jun 15 2025: %.0f days\n", c.X)
}
