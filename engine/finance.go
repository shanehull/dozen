package engine

import (
	"errors"
	"math"
)

// This file is the library API: pure financial functions in HP-12C
// terminology, with no dependency on the calculator's stack, prefixes, or
// keystroke model. The Engine (Step, X/Y/Z/T, ...) is a calculator front-end
// that delegates to these functions; use these directly when you just want the
// math.
//
// Sign convention (same as the HP-12C): cash inflows are positive and outflows
// negative, so a loan taken out has a positive PV and negative PMT. Rates are
// per period, expressed as a fraction (0.05 == 5%). n is the number of periods.

// Timing selects whether payments occur at the end of each period (an ordinary
// annuity) or the beginning (an annuity due, "BEGIN" mode on the HP-12C).
type Timing bool

const (
	End   Timing = false
	Begin Timing = true
)

var (
	// ErrNoConvergence is returned by Rate/IRR when the iterative solver fails
	// to converge (e.g. a cash-flow stream with no valid internal rate).
	ErrNoConvergence = errors.New("engine: solution did not converge")
	// ErrNoSolution is returned when the inputs admit no finite answer.
	ErrNoSolution = errors.New("engine: no solution for the given inputs")
)

// FV returns the future value of a series of cash flows.
func FV(rate, n, pv, pmt float64, when Timing) float64 {
	return tvmFV(rate, n, pv, pmt, bool(when))
}

// PV returns the present value of a series of cash flows.
func PV(rate, n, pmt, fv float64, when Timing) float64 {
	return tvmPV(rate, n, pmt, fv, bool(when))
}

// PMT returns the periodic payment for a loan or annuity.
func PMT(rate, n, pv, fv float64, when Timing) float64 {
	return tvmPMT(rate, n, pv, fv, bool(when))
}

// NPer solves for the number of periods. It returns ErrNoSolution when the
// cash flows never reach the target future value.
func NPer(rate, pmt, pv, fv float64, when Timing) (float64, error) {
	if isZero(rate) {
		if isZero(pmt) {
			return 0, ErrNoSolution
		}
		return clamp(-(fv + pv) / pmt), nil
	}
	b := 0.0
	if when {
		b = 1
	}
	adj := pmt * (1 + rate*b) / rate
	den := pv + adj
	if isZero(den) {
		return 0, ErrNoSolution
	}
	num := adj - fv
	if num/den <= 0 {
		return 0, ErrNoSolution
	}
	return clamp(math.Log(num/den) / math.Log(1+rate)), nil
}

// Rate solves for the periodic interest rate (as a fraction) using
// Newton-Raphson iteration. It returns ErrNoConvergence if no rate is found.
func Rate(n, pmt, pv, fv float64, when Timing) (float64, error) {
	begin := bool(when)
	guess := 0.001
	for iter := 0; iter < 100; iter++ {
		f := tvmEq(guess, n, pv, pmt, fv, begin)
		fp := tvmEqDeriv(guess, n, pv, pmt, fv, begin)
		if isZero(fp) {
			return 0, ErrNoConvergence
		}
		next := guess - f/fp
		if math.Abs(next-guess) < 1e-10 {
			return clamp(next), nil
		}
		guess = next
	}
	return 0, ErrNoConvergence
}

// SolveTVM solves for the missing time-value-of-money variable. Pass
// math.NaN() for exactly one of n, i, pv, pmt, or fv. i is the periodic
// rate as a fraction (0.05 = 5%). The result is the solved value or an
// error if the solver fails to converge.
func SolveTVM(n, i, pv, pmt, fv float64, when Timing) (float64, error) {
	var missing int
	if math.IsNaN(n) {
		missing++
	}
	if math.IsNaN(i) {
		missing++
	}
	if math.IsNaN(pv) {
		missing++
	}
	if math.IsNaN(pmt) {
		missing++
	}
	if math.IsNaN(fv) {
		missing++
	}
	if missing != 1 {
		return 0, ErrNoSolution
	}
	switch {
	case math.IsNaN(n):
		return NPer(i, pmt, pv, fv, when)
	case math.IsNaN(i):
		return Rate(n, pmt, pv, fv, when)
	case math.IsNaN(pv):
		return PV(i, n, pmt, fv, when), nil
	case math.IsNaN(pmt):
		return PMT(i, n, pv, fv, when), nil
	case math.IsNaN(fv):
		return FV(i, n, pv, pmt, when), nil
	}
	return 0, ErrNoSolution
}

// NPV returns the net present value of a cash-flow stream discounted at rate
// (a fraction). cashflows[0] is the period-0 flow (typically the initial
// outlay); cashflows[k] occurs at the end of period k.
func NPV(rate float64, cashflows ...float64) float64 {
	if len(cashflows) == 0 {
		return 0
	}
	npv := cashflows[0]
	for k := 1; k < len(cashflows); k++ {
		npv += cashflows[k] / math.Pow(1+rate, float64(k))
	}
	return clamp(npv)
}

// IRR returns the internal rate of return (as a fraction) of a cash-flow
// stream, i.e. the rate at which NPV is zero. cashflows[0] is the period-0
// flow. It returns ErrNoConvergence if the solver does not converge.
func IRR(cashflows ...float64) (float64, error) {
	guess := 0.1
	for iter := 0; iter < 100; iter++ {
		f := NPV(guess, cashflows...)
		d := npvDeriv(guess, cashflows)
		if isZero(d) {
			return 0, ErrNoConvergence
		}
		next := guess - f/d
		if math.Abs(next-guess) < 1e-10 {
			return clamp(next), nil
		}
		guess = next
	}
	return 0, ErrNoConvergence
}

func npvDeriv(rate float64, cashflows []float64) float64 {
	d := 0.0
	for k := 1; k < len(cashflows); k++ {
		d -= float64(k) * cashflows[k] / math.Pow(1+rate, float64(k+1))
	}
	return d
}

// DepSL returns straight-line depreciation for the given period and the
// remaining depreciable value after that period.
func DepSL(cost, salvage, life, period float64) (dep, remaining float64) {
	if isZero(life) {
		return math.NaN(), math.NaN()
	}
	dep = (cost - salvage) / life
	remaining = (cost - salvage) - dep*period
	return clamp(dep), clamp(remaining)
}

// DepSOYD returns sum-of-the-years'-digits depreciation for the given period
// and the remaining depreciable value after that period.
func DepSOYD(cost, salvage, life, period float64) (dep, remaining float64) {
	if isZero(life) {
		return math.NaN(), math.NaN()
	}
	soyd := life * (life + 1) / 2
	dep = (cost - salvage) * (life - period + 1) / soyd
	acc := 0.0
	for k := 1.0; k <= period; k++ {
		acc += (cost - salvage) * (life - k + 1) / soyd
	}
	return clamp(dep), clamp((cost - salvage) - acc)
}

// DepDB returns declining-balance depreciation for the given period and the
// remaining book value (above salvage) after that period. factorPct is the
// declining-balance factor as a percent (200 == double-declining balance);
// zero defaults to 200%.
func DepDB(cost, salvage, life, period, factorPct float64) (dep, remaining float64) {
	if isZero(life) {
		return math.NaN(), math.NaN()
	}
	factor := factorPct / 100
	if isZero(factor) {
		factor = 2
	}
	rate := factor / life
	book := cost
	for k := 1.0; k <= period; k++ {
		dep = book * rate
		if book-dep < salvage {
			dep = book - salvage
		}
		if dep < 0 {
			dep = 0
		}
		book -= dep
	}
	return clamp(dep), clamp(book - salvage)
}
