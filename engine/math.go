package engine

import (
	"errors"
	"math"
)

// TVM and cash-flow math, shared by the Engine methods.  All functions in
// this file are package‑private — the Engine is the public RPN interface.

// Sign convention (same as the HP-12C): cash inflows are positive and outflows
// negative, so a loan taken out has a positive PV and negative PMT. Rates are
// per period, expressed as a fraction (0.05 == 5%). n is the number of periods.

type timing bool

var (
	errNoConvergence = errors.New("engine: solution did not converge")
	errNoSolution    = errors.New("engine: no solution for the given inputs")
)

func fv(rate, n, pv, pmt float64, when timing) float64 {
	return tvmFV(rate, n, pv, pmt, bool(when))
}

func pv(rate, n, pmt, fv float64, when timing) float64 {
	return tvmPV(rate, n, pmt, fv, bool(when))
}

func pmt(rate, n, pv, fv float64, when timing) float64 {
	return tvmPMT(rate, n, pv, fv, bool(when))
}

func nper(rate, pmt, pv, fv float64, when timing) (float64, error) {
	if isZero(rate) {
		if isZero(pmt) {
			return 0, errNoSolution
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
		return 0, errNoSolution
	}
	num := adj - fv
	if num/den <= 0 {
		return 0, errNoSolution
	}
	return clamp(math.Log(num/den) / math.Log(1+rate)), nil
}

func rate(n, pmt, pv, fv float64, when timing) (float64, error) {
	beg := bool(when)
	guess := 0.001
	for iter := 0; iter < 100; iter++ {
		f := tvmEq(guess, n, pv, pmt, fv, beg)
		fp := tvmEqDeriv(guess, n, pv, pmt, fv, beg)
		if isZero(fp) {
			return 0, errNoConvergence
		}
		next := guess - f/fp
		if math.Abs(next-guess) < 1e-10 {
			return clamp(next), nil
		}
		guess = next
	}
	return 0, errNoConvergence
}

func solveTVM(N, I, PV, PMT, FV float64, when timing) (float64, error) {
	var missing int
	if math.IsNaN(N) {
		missing++
	}
	if math.IsNaN(I) {
		missing++
	}
	if math.IsNaN(PV) {
		missing++
	}
	if math.IsNaN(PMT) {
		missing++
	}
	if math.IsNaN(FV) {
		missing++
	}
	if missing != 1 {
		return 0, errNoSolution
	}
	switch {
	case math.IsNaN(N):
		return nper(I, PMT, PV, FV, when)
	case math.IsNaN(I):
		return rate(N, PMT, PV, FV, when)
	case math.IsNaN(PV):
		return pv(I, N, PMT, FV, when), nil
	case math.IsNaN(PMT):
		return pmt(I, N, PV, FV, when), nil
	case math.IsNaN(FV):
		return fv(I, N, PV, PMT, when), nil
	}
	return 0, errNoSolution
}

func npv(rate float64, cashflows ...float64) float64 {
	if len(cashflows) == 0 {
		return 0
	}
	v := cashflows[0]
	for k := 1; k < len(cashflows); k++ {
		v += cashflows[k] / math.Pow(1+rate, float64(k))
	}
	return clamp(v)
}

func irr(cashflows ...float64) (float64, error) {
	guess := 0.1
	for iter := 0; iter < 100; iter++ {
		f := npv(guess, cashflows...)
		d := npvDeriv(guess, cashflows)
		if isZero(d) {
			return 0, errNoConvergence
		}
		next := guess - f/d
		if math.Abs(next-guess) < 1e-10 {
			return clamp(next), nil
		}
		guess = next
	}
	return 0, errNoConvergence
}

func npvDeriv(rate float64, cashflows []float64) float64 {
	d := 0.0
	for k := 1; k < len(cashflows); k++ {
		d -= float64(k) * cashflows[k] / math.Pow(1+rate, float64(k+1))
	}
	return d
}

func depSL(cost, salvage, life, period float64) (dep, remaining float64) {
	if isZero(life) {
		return math.NaN(), math.NaN()
	}
	dep = (cost - salvage) / life
	remaining = (cost - salvage) - dep*period
	return clamp(dep), clamp(remaining)
}

func depSOYD(cost, salvage, life, period float64) (dep, remaining float64) {
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

func depDB(cost, salvage, life, period, factorPct float64) (dep, remaining float64) {
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

func amort(i, n, pv, pmt float64, begin bool) (totalInt, totalPrin float64) {
	for k := 1; k <= int(n); k++ {
		interest := pv * i
		if begin {
			interest = 0
			pv += pmt
			if pv < 0 {
				pv = 0
			}
		} else {
			pv += pmt + interest
		}
		totalInt += interest
		totalPrin += -pmt - interest
	}
	return clamp(totalInt), clamp(totalPrin)
}

func bondPrice(ytm, coupon float64, n float64) float64 {
	p := 0.0
	for k := 1.0; k <= n; k++ {
		p += coupon / 2 / math.Pow(1+ytm/2, k)
	}
	p += 100 / math.Pow(1+ytm/2, n)
	return clamp(p)
}

func bondYield(price, coupon float64, n float64) (float64, error) {
	guess := coupon / price
	for iter := 0; iter < 100; iter++ {
		p := bondPrice(guess, coupon, n)
		f := p - price
		d := 0.0
		for k := 1.0; k <= n; k++ {
			d -= float64(k) * coupon / (4 * math.Pow(1+guess/2, k+1))
		}
		d -= n * 100 / (4 * math.Pow(1+guess/2, n+1))
		if isZero(d) {
			break
		}
		next := guess - f/d
		if math.Abs(next-guess) < 1e-10 {
			return clamp(next), nil
		}
		guess = next
	}
	return 0, errNoConvergence
}
