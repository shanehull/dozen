package main

import (
	"strings"
	"testing"
)

func mathAbs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func press(svc *CalcService, op string, arg float64) KeyResult {
	return svc.PressKey(KeyInput{Op: op, Arg: arg, ArgS: op})
}

func assertDisplay(t *testing.T, r KeyResult, want string, msg string) {
	t.Helper()
	s := (r.Display.Sign) + r.Display.Mantissa
	if s != want {
		t.Errorf("%s: display = %q, want %q", msg, s, want)
	}
}

func assertFlag(t *testing.T, r KeyResult, flag string, want bool, msg string) {
	t.Helper()
	found := false
	for _, f := range r.Flags {
		if strings.EqualFold(f, flag) {
			found = true
			break
		}
	}
	if found != want {
		t.Errorf("%s: flag %q found=%v, want %v (flags=%v)", msg, flag, found, want, r.Flags)
	}
}

func assertStackX(t *testing.T, r KeyResult, want float64, msg string) {
	t.Helper()
	if r.StackX != want {
		t.Errorf("%s: stackX = %v, want %v", msg, r.StackX, want)
	}
}

func TestUIPercent(t *testing.T) {
	svc := NewCalcService()

	// 200 ENTER 15 %
	press(svc, "2", 2)
	press(svc, "0", 0)
	press(svc, "0", 0)
	press(svc, "ENTER", 0)
	press(svc, "1", 1)
	press(svc, "5", 5)
	r := press(svc, "%", 0)
	assertStackX(t, r, 30, "15% of 200")

	// Δ%: 200 ENTER 230 %CHG
	press(svc, "2", 2)
	press(svc, "0", 0)
	press(svc, "0", 0)
	press(svc, "ENTER", 0)
	press(svc, "2", 2)
	press(svc, "3", 3)
	press(svc, "0", 0)
	r = press(svc, "%CHG", 0)
	assertStackX(t, r, 15, "Δ% 200→230")
}

func TestUIStoreRecall(t *testing.T) {
	svc := NewCalcService()

	press(svc, "4", 4)
	press(svc, "2", 2)
	press(svc, "STO", 0)
	press(svc, "5", 5)
	press(svc, "CLx", 0)
	press(svc, "7", 7)
	press(svc, "RCL", 0)
	press(svc, "5", 5)
	r := press(svc, "+", 0)
	assertStackX(t, r, 49, "7 + RCL5(42)")
}

func TestUITVMStoreSolve(t *testing.T) {
	// TVM: type value + press key = store; just press key = solve
	svc := NewCalcService()

	// Store n=360
	press(svc, "3", 3)
	press(svc, "6", 6)
	press(svc, "0", 0)
	press(svc, "n", 0)
	// Verify n=360 by solving (should compute n from PV/PMT/FV/i — but those are 0 so NaN)

	// Now store and solve for FV: i=5, n=10, PV=-1000
	press(svc, "f", 0)
	press(svc, "CLEAR FIN", 0)
	
	press(svc, "1", 1)
	press(svc, "0", 0)
	press(svc, "n", 0)          // n=10
	press(svc, "5", 5)
	press(svc, "i", 0)          // i=5
	press(svc, "1", 1)
	press(svc, "0", 0)
	press(svc, "0", 0)
	press(svc, "0", 0)
	press(svc, "CHS", 0)
	press(svc, "PV", 0)         // PV=-1000

	r := press(svc, "FV", 0)    // solve for FV
	if mathAbs(r.StackX-1628.89) > 0.5 {
		t.Errorf("FV(5%%,10,-1000): stackX = %v, want ~1628.89", r.StackX)
	}
}

func TestUIPrefixFlags(t *testing.T) {
	svc := NewCalcService()

	r := press(svc, "f", 0)
	assertFlag(t, r, "f", true, "f prefix active")

	press(svc, "CLx", 0)
	r = press(svc, "g", 0)
	assertFlag(t, r, "g", true, "g prefix active")
}

func TestUIProgram(t *testing.T) {
	// Record: 5 ENTER 3 + and run
	svc := NewCalcService()

	press(svc, "f", 0)
	press(svc, "P/R", 0)
	press(svc, "5", 5)
	press(svc, "ENTER", 0)
	press(svc, "3", 3)
	press(svc, "+", 0)
	press(svc, "f", 0)
	press(svc, "P/R", 0)

	press(svc, "R/S", 0)
	r := press(svc, "ENTER", 0)
	if mathAbs(r.StackX-8) > 1e-10 {
		t.Errorf("Program: stackX = %v, want 8", r.StackX)
	}
}

func TestUITrig(t *testing.T) {
	svc := NewCalcService()

	press(svc, "1", 1)
	press(svc, "8", 8)
	press(svc, "0", 0)
	r := press(svc, "SIN", 0)
	if mathAbs(r.StackX) > 1e-10 {
		t.Errorf("sin(180): stackX = %v, want ~0", r.StackX)
	}
}

func TestUIDate(t *testing.T) {
	svc := NewCalcService()

	// 12.252025 ENTER 1.022026 g DYS
	press(svc, "1", 1)
	press(svc, "2", 2)
	press(svc, ".", 0)
	press(svc, "2", 2)
	press(svc, "5", 5)
	press(svc, "2", 2)
	press(svc, "0", 0)
	press(svc, "2", 2)
	press(svc, "5", 5)
	press(svc, "ENTER", 0)
	press(svc, "1", 1)
	press(svc, ".", 0)
	press(svc, "0", 0)
	press(svc, "2", 2)
	press(svc, "2", 2)
	press(svc, "0", 0)
	press(svc, "2", 2)
	press(svc, "6", 6)
	press(svc, "g", 0)
	r := press(svc, "DYS", 0)
	assertStackX(t, r, 8, "days between 12/25/2025 and 01/02/2026")
}

func TestUIReciprocal(t *testing.T) {
	svc := NewCalcService()

	press(svc, "4", 4)
	r := press(svc, "1/x", 0)
	assertStackX(t, r, 0.25, "1/4")
}

func TestUIChainCalc(t *testing.T) {
	svc := NewCalcService()

	// (10 + 5) × (20 / 4) = 15 × 5 = 75
	press(svc, "1", 1)
	press(svc, "0", 0)
	press(svc, "ENTER", 0)
	press(svc, "5", 5)
	press(svc, "+", 0)
	press(svc, "2", 2)
	press(svc, "0", 0)
	press(svc, "ENTER", 0)
	press(svc, "4", 4)
	press(svc, "÷", 0)
	r := press(svc, "×", 0)
	assertStackX(t, r, 75, "(10+5)×(20/4)")
}

func TestUIStackBasics(t *testing.T) {
	svc := NewCalcService()

	// 5 ENTER 3 + 4 ×
	for _, d := range []float64{5, 0, 3, 0, 0, 4, 0} {
		var op string
		switch {
		case d == 5: op = "5"
		case d == 3: op = "3"
		case d == 4: op = "4"
		default: op = "ENTER"
		}
		if op == "ENTER" && d == 0 {
			press(svc, "ENTER", 0)
		} else if d == 5 {
			press(svc, "5", 5)
		} else if d == 3 {
			press(svc, "3", 3)
			press(svc, "+", 0)
		} else if d == 4 {
			_ = d
		}
	}
	// Manually: 5 ENTER 3 + 4 ×
	svc2 := NewCalcService()
	press(svc2, "5", 5)
	press(svc2, "ENTER", 0)
	press(svc2, "3", 3)
	press(svc2, "+", 0)
	press(svc2, "4", 4)
	r := press(svc2, "×", 0)
	assertStackX(t, r, 32, "5+3×4")
}

func TestUISqrt(t *testing.T) {
	svc := NewCalcService()

	press(svc, "1", 1)
	press(svc, "6", 6)
	r := press(svc, "√x", 0)
	assertStackX(t, r, 4, "sqrt(16)")
}

func TestUIAbs(t *testing.T) {
	svc := NewCalcService()

	press(svc, "5", 5)
	press(svc, "CHS", 0)
	r := press(svc, "|x|", 0)
	assertStackX(t, r, 5, "|-5|")
}

func TestUIIntgFrac(t *testing.T) {
	svc := NewCalcService()

	// 12.34
	press(svc, "1", 1)
	press(svc, "2", 2)
	press(svc, ".", 0)
	press(svc, "3", 3)
	press(svc, "4", 4)

	r := press(svc, "INTG", 0)
	if mathAbs(r.StackX-12) > 1e-10 {
		t.Errorf("INTG(12.34): stackX = %v, want 12", r.StackX)
	}

	press(svc, "1", 1)
	press(svc, "2", 2)
	press(svc, ".", 0)
	press(svc, "3", 3)
	press(svc, "4", 4)
	r = press(svc, "FRAC", 0)
	if mathAbs(r.StackX-0.34) > 1e-10 {
		t.Errorf("FRAC(12.34): stackX = %v, want 0.34", r.StackX)
	}
}

func TestUILogLn(t *testing.T) {
	svc := NewCalcService()

	press(svc, "1", 1)
	press(svc, "0", 0)
	press(svc, "0", 0)
	r := press(svc, "LN", 0)
	// ln(100) ≈ 4.605
	if r.StackX < 4.6 || r.StackX > 4.61 {
		t.Errorf("ln(100): stackX = %v, want ~4.605", r.StackX)
	}

	press(svc, "1", 1)
	press(svc, "0", 0)
	press(svc, "0", 0)
	press(svc, "0", 0)
	r = press(svc, "LOG", 0)
	// log10(1000) = 3
	assertStackX(t, r, 3, "log10(1000)")
}

func TestUIDisplayStates(t *testing.T) {
	svc := NewCalcService()

	// While typing, display shows raw entry buffer
	r := press(svc, "0", 0)
	assertDisplay(t, r, "0", "typing 0")

	r = press(svc, "1", 1)
	r = press(svc, "0", 0)
	r = press(svc, "0", 0)
	assertDisplay(t, r, "100", "typing 100")

	r = press(svc, "ENTER", 0)
	assertDisplay(t, r, "100.00", "after ENTER, Fix 2 format")

	r = press(svc, "CLx", 0)
	assertDisplay(t, r, "0.00", "after CLx")
}

func TestUIXYSwap(t *testing.T) {
	svc := NewCalcService()

	press(svc, "5", 5)
	press(svc, "ENTER", 0)
	press(svc, "3", 3)
	r := press(svc, "x↔y", 0)
	assertStackX(t, r, 5, "x↔y: X was 3, now 5")
}

func TestUIEEX(t *testing.T) {
	svc := NewCalcService()

	// 1.5 EEX 3 = 1500
	press(svc, "1", 1)
	press(svc, ".", 0)
	press(svc, "5", 5)
	press(svc, "EEX", 0)
	press(svc, "3", 3)
	r := press(svc, "ENTER", 0)
	assertStackX(t, r, 1500, "1.5 EEX 3")
}

func TestUIAllFlagsOff(t *testing.T) {
	svc := NewCalcService()

	r := press(svc, "CLx", 0)
	assertFlag(t, r, "f", false, "f off initially")
	assertFlag(t, r, "g", false, "g off initially")
	assertFlag(t, r, "BEGIN", false, "begin off initially")
}
