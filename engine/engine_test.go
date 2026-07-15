package engine

import (
	"math"
	"testing"
)

func TestStack(t *testing.T) {
	e := New()

	e.X = 1
	e.push()
	if e.X != 1 || e.Y != 1 || e.Z != 0 || e.T != 0 {
		t.Fatalf("push1: want X=1 Y=1 Z=0 T=0, got X=%v Y=%v Z=%v T=%v", e.X, e.Y, e.Z, e.T)
	}

	e.X = 2
	e.push()
	if e.X != 2 || e.Y != 2 || e.Z != 1 || e.T != 0 {
		t.Fatalf("push2: want X=2 Y=2 Z=1 T=0, got X=%v Y=%v Z=%v T=%v", e.X, e.Y, e.Z, e.T)
	}

	e.X = 3
	e.push()
	e.X = 4
	e.push()
	e.tuck()
	// After 4 pushes (1,1,2,2,3,3,4,4) and tuck:
	// X=4, Y=3 (was Z), Z=2 (was T), T=2 (replicated)
	if e.X != 4 || e.Y != 3 || e.Z != 2 || e.T != 2 {
		t.Fatalf("tuck: want X=4 Y=3 Z=2 T=2, got X=%v Y=%v Z=%v T=%v", e.X, e.Y, e.Z, e.T)
	}
}

func TestArithmetic(t *testing.T) {
	e := New()
	e.X, e.Y = 3, 4
	e.add()
	if e.X != 7 {
		t.Fatalf("3+4: want 7, got %v", e.X)
	}

	e.X, e.Y = 3, 10
	e.sub()
	if e.X != 7 {
		t.Fatalf("10-3: want 7, got %v", e.X)
	}

	e.X, e.Y = 6, 7
	e.mul()
	if e.X != 42 {
		t.Fatalf("6*7: want 42, got %v", e.X)
	}

	e.X, e.Y = 4, 12
	e.div()
	if e.X != 3 {
		t.Fatalf("12/4: want 3, got %v", e.X)
	}
}

func TestDigitEntry(t *testing.T) {
	e := New()
	e.Flags.StackLift = true

	e.Step("4", 4, "4")
	if e.X != 4 {
		t.Fatalf("digit 4: want 4, got %v", e.X)
	}

	e.Step("2", 2, "2")
	if e.X != 42 {
		t.Fatalf("digit 42: want 42, got %v", e.X)
	}

	e.Step(".", 0, ".")
	e.Step("5", 5, "5")
	if e.X != 42.5 {
		t.Fatalf("decimal: want 42.5, got %v", e.X)
	}
}

func TestBasicOps(t *testing.T) {
	e := New()

	// CHS
	e.X = 5
	e.chs()
	if e.X != -5 {
		t.Fatalf("chs(5): want -5, got %v", e.X)
	}

	// CLx
	e.X = 42
	e.clx()
	if e.X != 0 {
		t.Fatal("clx: want 0")
	}

	// x↔y
	e.X, e.Y = 3, 7
	e.xy()
	if e.X != 7 || e.Y != 3 {
		t.Fatalf("xy: want X=7 Y=3, got X=%v Y=%v", e.X, e.Y)
	}
}

func TestEnterAndLastX(t *testing.T) {
	e := New()
	e.Flags.StackLift = true

	e.Step("5", 5, "5")
	e.Step("ENTER", 0, "")
	if e.X != 5 || e.Y != 5 {
		t.Fatalf("ENTER: want X=5 Y=5, got X=%v Y=%v", e.X, e.Y)
	}

	e.Step("3", 3, "3")
	e.Step("+", 0, "")
	if e.X != 8 {
		t.Fatalf("5+3: want 8, got %v", e.X)
	}

	if e.LastX != 3 {
		t.Fatalf("lastX: want 3, got %v", e.LastX)
	}

	e.Step("LSTx", 0, "")
	if e.X != 3 {
		t.Fatalf("LSTx op: want 3, got %v", e.X)
	}
}

func TestStoreRecall(t *testing.T) {
	e := New()
	e.X = 99
	e.store(0, "")
	if e.Mem[0] != 99 {
		t.Fatalf("STO 0: want 99, got %v", e.Mem[0])
	}

	e.X = 0
	e.recall(0, "")
	if e.X != 99 {
		t.Fatalf("RCL 0: want 99, got %v", e.X)
	}
}

func TestPercent(t *testing.T) {
	e := New()
	e.X, e.Y = 15, 200
	e.pct()
	if e.X != 30 {
		t.Fatalf("200%% d15: want 30, got %v", e.X)
	}
}

func TestPctChg(t *testing.T) {
	e := New()
	e.X, e.Y = 125, 100
	e.pctChg()
	if e.X != 25 {
		t.Fatalf("pctChg: want 25, got %v", e.X)
	}
}

func TestScientific(t *testing.T) {
	e := New()

	e.X = 5
	e.fact()
	if e.X != 120 {
		t.Fatalf("5!: want 120, got %v", e.X)
	}

	e.X = 16
	e.sqrt()
	if e.X != 4 {
		t.Fatalf("sqrt(16): want 4, got %v", e.X)
	}

	e.X = 100
	e.log()
	if math.Abs(e.X-2) > 0.0001 {
		t.Fatalf("log(100): want 2, got %v", e.X)
	}

	e.X = math.E
	e.ln()
	if math.Abs(e.X-1) > 0.0001 {
		t.Fatalf("ln(e): want 1, got %v", e.X)
	}

	e.X = 4
	e.recip()
	if e.X != 0.25 {
		t.Fatalf("1/4: want 0.25, got %v", e.X)
	}

	e.X = -7
	e.abs()
	if e.X != 7 {
		t.Fatalf("|-7|: want 7, got %v", e.X)
	}

	e.X = 3.14159
	e.intg()
	if e.X != 3 {
		t.Fatalf("INTG: want 3, got %v", e.X)
	}

	e.X = 3.75
	e.frac()
	if math.Abs(e.X-0.75) > 1e-10 {
		t.Fatalf("FRAC: want 0.75, got %v", e.X)
	}
}

func TestTrig(t *testing.T) {
	e := New()
	e.Flags.Angle = Deg

	e.X = 90
	e.sin()
	if math.Abs(e.X-1) > 0.0001 {
		t.Fatalf("sin(90): want 1, got %v", e.X)
	}

	e.X = 0
	e.cos()
	if math.Abs(e.X-1) > 0.0001 {
		t.Fatalf("cos(0): want 1, got %v", e.X)
	}

	e.X = 45
	e.tan()
	if math.Abs(e.X-1) > 0.001 {
		t.Fatalf("tan(45): want ≈1, got %v", e.X)
	}
}

func TestPi(t *testing.T) {
	e := New()
	e.pi()
	if math.Abs(e.X-math.Pi) > 0.0001 {
		t.Fatalf("pi: want %v, got %v", math.Pi, e.X)
	}
}

func TestPolarRect(t *testing.T) {
	e := New()
	e.X, e.Y = 3, 4
	e.toPolar()
	if math.Abs(e.X-5) > 0.0001 {
		t.Fatalf("toPolar r: want 5, got %v", e.X)
	}

	e.Flags.Angle = Rad
	e.X, e.Y = 5, 0.927295
	e.toRect()
	if math.Abs(e.X-3) > 0.01 || math.Abs(e.Y-4) > 0.01 {
		t.Fatalf("toRect: want X≈3 Y≈4, got X=%v Y=%v", e.X, e.Y)
	}
}

func TestYPowX(t *testing.T) {
	e := New()
	e.X, e.Y = 3, 2
	e.yPowX()
	if e.X != 8 {
		t.Fatalf("2^3: want 8, got %v", e.X)
	}
}

func TestTVM(t *testing.T) {
	e := New()

	e.FinN = 12
	e.FinI = 5
	e.FinPV = -1000
	e.FinPMT = 0

	e.push()
	e.tvmFV()
	if math.Abs(e.X-1795.86) > 0.1 {
		t.Fatalf("FV: want ≈1795.86, got %v", e.X)
	}

	e.FinFV = 2000
	e.push()
	e.tvmI()
	if math.Abs(e.X-5.946) > 0.1 {
		t.Fatalf("IRR: want ≈5.95, got %v", e.X)
	}

	e.FinPV = 1000
	e.FinFV = -2000
	e.push()
	e.tvmN()
	if math.Abs(e.X-14.2) > 0.2 {
		t.Fatalf("n: want ≈14.2, got %v", e.X)
	}

	e.FinI = 5
	e.FinN = 12
	e.FinPV = 10000
	e.FinFV = 0
	e.push()
	e.tvmPMT()
	if math.Abs(e.X+1128.25) > 0.5 {
		t.Fatalf("PMT: want ≈-1128.25, got %v", e.X)
	}
}

func TestNPVAndIRR(t *testing.T) {
	e := New()

	e.FinCF0 = -100
	e.FinCFj[0] = 50
	e.FinCFj[1] = 60
	e.FinCFj[2] = 70
	e.FinCfCnt = 3
	e.X = 10
	e.finNPV()
	if math.Abs(e.X-47.63) > 0.1 {
		t.Fatalf("NPV: want ≈47.63, got %v", e.X)
	}

	e.FinCF0 = -100
	e.FinCFj[0] = 60
	e.FinCFj[1] = 60
	e.FinCfCnt = 2
	e.Flags.StackLift = true
	e.finIRR()
	if math.Abs(e.X-13.07) > 0.5 {
		t.Fatalf("IRR: want ≈13.07, got %v", e.X)
	}
}

func TestStats(t *testing.T) {
	e := New()
	// HP-12C: Σ+ takes x from Y, y from X
	e.X, e.Y = 20, 10 // y=20, x=10
	e.statAdd()
	e.X, e.Y = 30, 20 // y=30, x=20
	e.statAdd()

	e.statMeanX()
	if e.X != 15 {
		t.Fatalf("x̄: want 15, got %v", e.X)
	}

	e.statMeanY()
	if e.X != 25 {
		t.Fatalf("ȳ: want 25, got %v", e.X)
	}

	e.clearStats()
	e.X, e.Y = 0, 4
	e.statAdd()
	e.X, e.Y = 0, 8
	e.statAdd()
	e.statSDev()
	if math.Abs(e.X-2.828) > 0.01 {
		t.Fatalf("s: want ≈2.828, got %v", e.X)
	}
}

func TestLinEst(t *testing.T) {
	e := New()
	// y = 2x + 1: (x=1,y=3), (x=2,y=5), (x=3,y=7)
	e.X, e.Y = 3, 1
	e.statAdd()
	e.X, e.Y = 5, 2
	e.statAdd()
	e.X, e.Y = 7, 3
	e.statAdd()

	e.statLinEst()
	// r=1, intercept=1 (y-intercept)
	if math.Abs(e.X-1) > 0.001 {
		t.Fatalf("r: want 1, got %v", e.X)
	}
	if math.Abs(e.Y-1) > 0.001 {
		t.Fatalf("intercept: want 1, got %v", e.Y)
	}
}

func TestDateDays(t *testing.T) {
	e := New()
	e.Flags.Dmy = false
	e.Y = 1.012025
	e.X = 1.152025
	e.dateDys()
	if e.X != 14 {
		t.Fatalf("daysBetween: want 14, got %v", e.X)
	}
}

func TestDateFuture(t *testing.T) {
	e := New()
	e.Flags.Dmy = false
	e.Y = 1.012025
	e.X = 30
	e.dateDate()
	if math.Abs(e.X-1.312025) > 1e-6 {
		t.Fatalf("date+30: want 1.312025, got %v", e.X)
	}
}

func TestToHMS(t *testing.T) {
	e := New()
	e.X = 2.5
	e.toHMS()
	if math.Abs(e.X-2.3) > 0.001 {
		t.Fatalf("2.5h→HMS: want 2.3, got %v", e.X)
	}
}

func TestToH(t *testing.T) {
	e := New()
	e.X = 2.3
	e.toH()
	if math.Abs(e.X-2.5) > 0.001 {
		t.Fatalf("2.3 HMS→h: want 2.5, got %v", e.X)
	}
}

func TestClearRegs(t *testing.T) {
	e := New()
	e.Mem[0] = 99
	e.StatsN = 5
	e.FinN = 10
	e.clearReg()
	if e.Mem[0] != 0 || e.StatsN != 0 || e.FinN != 0 {
		t.Fatal("clearReg did not reset")
	}
}

func TestPrefixF(t *testing.T) {
	e := New()
	e.X = 42
	e.Step("f", 0, "")
	if e.Flags.Prefix != "f" {
		t.Fatal("prefix not set")
	}
	e.Step("n", 0, "")
	if e.FinN != 42 {
		t.Fatalf("f+n: want FinN=42, got %v", e.FinN)
	}
	if e.Flags.Prefix != "" {
		t.Fatal("prefix not cleared")
	}
}

func TestDisplay(t *testing.T) {
	e := New()
	e.Flags.DispMode = Fix
	e.Flags.DispDigits = 2

	e.X = 1234.5678
	d := e.Format()
	if d.Mantissa == "" {
		t.Fatal("display empty")
	}

	e.X = math.NaN()
	d = e.Format()
	if d.Mantissa != " Error    " {
		t.Fatalf("NaN display: got '%v'", d.Mantissa)
	}
}

func TestIntegration(t *testing.T) {
	e := New()
	e.Flags.StackLift = true

	// 5 ENTER 3 + → 8
	e.Step("5", 5, "5")
	e.Step("ENTER", 0, "")
	e.Step("3", 3, "3")
	e.Step("+", 0, "")
	if e.X != 8 {
		t.Fatalf("5 ENTER 3 +: want 8, got %v", e.X)
	}

	// 4 × → 32
	e.Step("4", 4, "4")
	e.Step("×", 0, "")
	if e.X != 32 {
		t.Fatalf("4 ×: want 32, got %v", e.X)
	}
}
