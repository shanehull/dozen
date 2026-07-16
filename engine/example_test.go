package engine_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/shanehull/dozen/engine"
)

func Example_mortgage() {
	c := engine.New()

	c.FinN = 30 * 12
	c.FinI = 6.0 / 12.0
	c.FinPV = 300000
	c.FinFV = 0
	c.SolvePMT()

	fmt.Printf("Monthly payment: $%.2f\n", -c.X)
	// Output: Monthly payment: $1798.65
}

func Example_futureValue() {
	c := engine.New()

	c.FinN = 20
	c.FinI = 8
	c.FinPV = -10000
	c.FinPMT = 0
	c.SolveFV()

	fmt.Printf("Future value: $%.2f\n", c.X)
	// Output: Future value: $46609.57
}

func Example_npv() {
	c := engine.New()

	c.FinCF0 = -100000
	c.FinCFj[0] = 30000
	c.FinCFj[1] = 40000
	c.FinCFj[2] = 50000
	c.FinCFj[3] = 60000
	c.FinCfCnt = 4
	c.X = 10
	c.Flags.StackLift = true
	c.ComputeNPV()

	fmt.Printf("NPV: $%.2f\n", c.X)
	// Output: NPV: $38877.13
}

func Example_irr() {
	c := engine.New()

	c.FinCF0 = -100000
	c.FinCFj[0] = 40000
	c.FinCFj[1] = 50000
	c.FinCFj[2] = 30000
	c.FinCfCnt = 3
	c.Flags.StackLift = true
	c.ComputeIRR()

	fmt.Printf("IRR: %.2f%%\n", c.X)
	// Output: IRR: 10.13%
}

func Example_bondPrice() {
	c := engine.New()

	c.FinN = 5
	c.Y = 7
	c.X = 6
	c.BondPrice()

	fmt.Printf("Bond price: $%.2f\n", c.X)
	// Output: Bond price: $97.74
}

func Example_depreciationSL() {
	c := engine.New()

	c.FinN = 10
	c.FinPV = 50000
	c.FinFV = 5000
	c.X = 1
	c.DepreciationSL()

	fmt.Printf("Annual depreciation: $%.2f\n", c.X)
	// Output: Annual depreciation: $4500.00
}

func Example_statistics() {
	c := engine.New()

	c.X, c.Y = 10, 1
	c.StatAdd()
	c.X, c.Y = 15, 2
	c.StatAdd()
	c.X, c.Y = 22, 3
	c.StatAdd()
	c.X, c.Y = 18, 4
	c.StatAdd()
	c.X, c.Y = 25, 5
	c.StatAdd()

	c.MeanX()
	meanX := c.X

	c.SDev()
	sdev := c.X

	fmt.Printf("Mean: %.1f, StdDev: %.1f\n", meanX, sdev)
	// Output: Mean: 3.0, StdDev: 1.6
}

func Example_dateDays() {
	c := engine.New()

	c.Flags.Dmy = false
	c.Y = 1.012025
	c.X = 6.152025
	c.DaysBetween()

	fmt.Printf("Days between: %.0f\n", c.X)
	// Output: Days between: 165
}

func Example_dateAdd() {
	c := engine.New()

	c.Flags.Dmy = false
	c.Y = 3.012025
	c.X = 90
	c.DateAdd()

	fmt.Printf("Date after 90 days: %.6f\n", c.X)
	// Output: Date after 90 days: 5.302025
}

func TestExamples(t *testing.T) {
	t.Run("mortgage", func(t *testing.T) {
		c := engine.New()
		c.FinN = 360
		c.FinI = 0.5
		c.FinPV = 300000
		c.FinFV = 0
		c.SolvePMT()
		if math.Abs(c.X+1798.65) > 0.1 {
			t.Fatalf("PMT: want -1798.65, got %v", c.X)
		}
	})

	t.Run("future_value", func(t *testing.T) {
		c := engine.New()
		c.FinN = 20
		c.FinI = 8
		c.FinPV = -10000
		c.FinPMT = 0
		c.SolveFV()
		if math.Abs(c.X-46609.57) > 0.5 {
			t.Fatalf("FV: want 46609.57, got %v", c.X)
		}
	})

	t.Run("npv", func(t *testing.T) {
		c := engine.New()
		c.FinCF0 = -100000
		c.FinCFj[0] = 30000
		c.FinCFj[1] = 40000
		c.FinCFj[2] = 50000
		c.FinCFj[3] = 60000
		c.FinCfCnt = 4
		c.X = 10
		c.Flags.StackLift = true
		c.ComputeNPV()
		if math.Abs(c.X-38877.13) > 1 {
			t.Fatalf("NPV: want ~38877, got %v", c.X)
		}
	})

	t.Run("irr", func(t *testing.T) {
		c := engine.New()
		c.FinCF0 = -100000
		c.FinCFj[0] = 40000
		c.FinCFj[1] = 50000
		c.FinCFj[2] = 30000
		c.FinCfCnt = 3
		c.Flags.StackLift = true
		c.ComputeIRR()
		if math.Abs(c.X-10.13) > 0.5 {
			t.Fatalf("IRR: want ~10.13, got %v", c.X)
		}
	})

	t.Run("bond_price", func(t *testing.T) {
		c := engine.New()
		c.FinN = 5
		c.Y = 7
		c.X = 6
		c.BondPrice()
		if math.Abs(c.X-97.74) > 0.1 {
			t.Fatalf("bond price: want ~97.74, got %v", c.X)
		}
	})

	t.Run("depreciation_sl", func(t *testing.T) {
		c := engine.New()
		c.FinN = 10
		c.FinPV = 50000
		c.FinFV = 5000
		c.X = 1
		c.DepreciationSL()
		if math.Abs(c.X-4500) > 0.01 {
			t.Fatalf("dep SL: want 4500, got %v", c.X)
		}
	})

	t.Run("depreciation_soyd", func(t *testing.T) {
		c := engine.New()
		c.FinN = 5
		c.FinPV = 1000
		c.FinFV = 100
		c.X = 1
		c.DepreciationSOYD()
		if math.Abs(c.X-300) > 0.01 {
			t.Fatalf("dep SOYD yr1: want 300, got %v", c.X)
		}
	})

	t.Run("depreciation_db", func(t *testing.T) {
		c := engine.New()
		c.FinN = 5
		c.FinI = 200
		c.FinPV = 1000
		c.FinFV = 100
		c.X = 1
		c.DepreciationDB()
		if math.Abs(c.X-400) > 0.01 {
			t.Fatalf("dep DB yr1: want 400, got %v", c.X)
		}
	})

	t.Run("statistics", func(t *testing.T) {
		c := engine.New()
		c.X, c.Y = 10, 1
		c.StatAdd()
		c.X, c.Y = 15, 2
		c.StatAdd()
		c.X, c.Y = 22, 3
		c.StatAdd()
		c.X, c.Y = 18, 4
		c.StatAdd()
		c.X, c.Y = 25, 5
		c.StatAdd()
		c.MeanX()
		if math.Abs(c.X-3.0) > 0.01 {
			t.Fatalf("x: want 3.0, got %v", c.X)
		}
	})

	t.Run("date_days", func(t *testing.T) {
		c := engine.New()
		c.Flags.Dmy = false
		c.Y = 1.012025
		c.X = 6.152025
		c.DaysBetween()
		if c.X != 165 {
			t.Fatalf("days: want 165, got %v", c.X)
		}
	})

	t.Run("date_add", func(t *testing.T) {
		c := engine.New()
		c.Flags.Dmy = false
		c.Y = 3.012025
		c.X = 90
		c.DateAdd()
		if math.Abs(c.X-5.302025) > 1e-6 {
			t.Fatalf("date+90: want 5.302025, got %v", c.X)
		}
	})
}
