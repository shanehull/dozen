package engine

import "math"

// TVM solver: 0 = PV*(1+i)^n + PMT*((1+i*B)*(1+i)^n-1)/i + FV
// where B=1 for BEGIN, B=0 for END
// Signed convention: inflows positive, outflows negative.

// Store TVM register values from X.
func (e *Engine) finN()  { e.FinN = e.X }
func (e *Engine) finI()  { e.FinI = e.X }
func (e *Engine) finPV() { e.FinPV = e.X }
func (e *Engine) finPMT() { e.FinPMT = e.X }
func (e *Engine) finFV() { e.FinFV = e.X }

// Solve for n (number of periods).
func (e *Engine) tvmN() {
	n, err := NPer(e.FinI/100, e.FinPMT, e.FinPV, e.FinFV, Timing(e.Flags.Begin))
	if err != nil {
		e.X = math.NaN()
		return
	}
	e.push()
	e.X = clamp(n)
	e.Flags.StackLift = true
}

func (e *Engine) tvmI() {
	r, err := Rate(e.FinN, e.FinPMT, e.FinPV, e.FinFV, Timing(e.Flags.Begin))
	if err != nil {
		e.X = math.NaN()
		return
	}
	e.push()
	e.X = clamp(r * 100)
	e.Flags.StackLift = true
}

func (e *Engine) tvmPV() {
	e.push()
	e.X = clamp(tvmPV(e.FinI/100, e.FinN, e.FinPMT, e.FinFV, e.Flags.Begin))
	e.Flags.StackLift = true
}

func (e *Engine) tvmPMT() {
	e.push()
	e.X = clamp(tvmPMT(e.FinI/100, e.FinN, e.FinPV, e.FinFV, e.Flags.Begin))
	e.Flags.StackLift = true
}

func (e *Engine) tvmFV() {
	e.push()
	e.X = clamp(tvmFV(e.FinI/100, e.FinN, e.FinPV, e.FinPMT, e.Flags.Begin))
	e.Flags.StackLift = true
}

func tvmPV(i, n, pmt, fv float64, begin bool) float64 {
	if isZero(i) {
		return clamp(-(fv + pmt*n))
	}
	b := 0.0
	if begin {
		b = 1.0
	}
	pv := -fv/math.Pow(1+i, n) - pmt*(1+i*b)*(1-1/math.Pow(1+i, n))/i
	return clamp(pv)
}

func tvmPMT(i, n, pv, fv float64, begin bool) float64 {
	if isZero(i) {
		if isZero(n) {
			return clamp(-(fv + pv) / n)
		}
		return clamp(-(fv + pv) / n)
	}
	b := 0.0
	if begin {
		b = 1.0
	}
	pmt := -(pv + fv/math.Pow(1+i, n)) * i / ((1 + i*b) * (1 - 1/math.Pow(1+i, n)))
	return clamp(pmt)
}

func tvmFV(i, n, pv, pmt float64, begin bool) float64 {
	if isZero(i) {
		return clamp(-(pv + pmt*n))
	}
	b := 0.0
	if begin {
		b = 1.0
	}
	fv := -pv*math.Pow(1+i, n) - pmt*(1+i*b)*(math.Pow(1+i, n)-1)/i
	return clamp(fv)
}

func tvmEq(i, n, pv, pmt, fv float64, begin bool) float64 {
	b := 0.0
	if begin {
		b = 1.0
	}
	if isZero(i) {
		return pv + fv + pmt*n*(1+i*b)
	}
	return pv*math.Pow(1+i, n) + pmt*(1+i*b)*(math.Pow(1+i, n)-1)/i + fv
}

func tvmEqDeriv(i, n, pv, pmt, fv float64, begin bool) float64 {
	b := 0.0
	if begin {
		b = 1.0
	}
	if isZero(i) {
		return pv*n + pmt*n*n/2 + pmt*n + fv*n
	}
	p := math.Pow(1+i, n-1)
	return pv*n*p + pmt*(1+i*b)*(n*p*(1+i) - (math.Pow(1+i, n)-1)/i) / (1 + i)
}

// cashflows assembles CF0 and the CFj entries into a single slice for the
// functional NPV/IRR helpers.
func (e *Engine) cashflows() []float64 {
	cfs := make([]float64, 0, e.FinCfCnt+1)
	cfs = append(cfs, e.FinCF0)
	cfs = append(cfs, e.FinCFj[:e.FinCfCnt]...)
	return cfs
}

// NPV using stored CF0 and CFj, discounted at the rate keyed into X (percent).
func (e *Engine) finNPV() {
	// NPV consumes the discount rate that was in X and replaces it with the
	// result; the rest of the stack is untouched.
	e.X = clamp(NPV(e.X/100, e.cashflows()...))
	e.Flags.StackLift = true
}

func (e *Engine) finIRR() {
	r, _ := IRR(e.cashflows()...)
	e.result(r * 100)
}

// Amortization: computes interest and principal for n periods.
func (e *Engine) finAmort() {
	n := e.X
	i := e.FinI / 100
	pv := e.FinPV
	pmt := e.FinPMT
	begin := e.Flags.Begin

	totalInt := 0.0
	totalPrin := 0.0

	for k := 1; k <= int(n); k++ {
		interest := pv * i
		if begin {
			interest = 0
			pv -= pmt
			if pv < 0 {
				pv = 0
			}
		} else {
			principal := pmt - interest
			pv -= principal
		}
		totalInt += interest
		totalPrin += pmt - interest
	}

	e.AmortInt = clamp(totalInt)
	e.AmortPrin = clamp(totalPrin)
	e.AmortN = n

	e.X = e.AmortInt
	e.push()
	e.X = e.AmortPrin
	e.Y = e.AmortInt
	e.Flags.StackLift = true
}

func (e *Engine) finInt() {
	e.X = e.AmortInt
}

// Bond calculations.
func (e *Engine) bondPrice() {
	// Simplified: price as % of par using yield
	ytm := e.Y / 100
	coupon := e.X
	n := e.FinN
	price := 0.0
	for k := 1.0; k <= n; k++ {
		price += coupon / 2 / math.Pow(1+ytm/2, k)
	}
	price += 100 / math.Pow(1+ytm/2, n)
	e.LastX = e.X
	e.X = clamp(price)
	e.tuck()
	e.Flags.StackLift = true
}

func (e *Engine) bondYield() {
	// Simplified yield-to-maturity
	price := e.X
	coupon := e.Y
	n := e.FinN

	guess := coupon / price
	for iter := 0; iter < 100; iter++ {
		p := 0.0
		for k := 1.0; k <= n; k++ {
			p += coupon / 2 / math.Pow(1+guess/2, k)
		}
		p += 100 / math.Pow(1+guess/2, n)
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
			e.LastX = e.X
			e.X = clamp(next * 100)
			e.tuck()
			e.Flags.StackLift = true
			return
		}
		guess = next
	}
	e.LastX = e.X
	e.X = clamp(guess * 100)
	e.tuck()
	e.Flags.StackLift = true
}

// Depreciation (SL / SOYD / DB) follows the HP-12C register convention:
//
//	FinN  = asset life (years)
//	FinPV = depreciable cost basis
//	FinFV = salvage value
//	FinI  = declining-balance factor as a percent (DB only, e.g. 200 = 200%)
//	X     = period (year) number for which to compute depreciation
//
// Each returns the depreciation for that period in X and the remaining
// depreciable/book value in Y.

// depResult stores a depreciation amount in X and remaining value in Y.
func (e *Engine) depResult(dep, remaining float64) {
	e.LastX = e.X
	e.X = clamp(dep)
	e.Y = clamp(remaining)
	e.Flags.StackLift = true
}

func (e *Engine) depSL() {
	dep, rem := DepSL(e.FinPV, e.FinFV, e.FinN, e.X)
	if math.IsNaN(dep) {
		e.X = math.NaN()
		return
	}
	e.depResult(dep, rem)
}

func (e *Engine) depSOYD() {
	dep, rem := DepSOYD(e.FinPV, e.FinFV, e.FinN, e.X)
	if math.IsNaN(dep) {
		e.X = math.NaN()
		return
	}
	e.depResult(dep, rem)
}

func (e *Engine) depDB() {
	dep, rem := DepDB(e.FinPV, e.FinFV, e.FinN, e.X, e.FinI)
	if math.IsNaN(dep) {
		e.X = math.NaN()
		return
	}
	e.depResult(dep, rem)
}
