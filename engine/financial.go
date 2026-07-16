package engine

import "math"

func (e *Engine) SetN()  { e.FinN = e.X }
func (e *Engine) SetI()  { e.FinI = e.X }
func (e *Engine) SetPV() { e.FinPV = e.X }
func (e *Engine) SetPMT() { e.FinPMT = e.X }
func (e *Engine) SetFV() { e.FinFV = e.X }

func (e *Engine) SolveN() {
	n, err := NPer(e.FinI/100, e.FinPMT, e.FinPV, e.FinFV, Timing(e.Flags.Begin))
	if err != nil { e.X = math.NaN(); return }
	e.Push()
	e.X = clamp(n)
	e.Flags.StackLift = true
}

func (e *Engine) SolveI() {
	r, err := Rate(e.FinN, e.FinPMT, e.FinPV, e.FinFV, Timing(e.Flags.Begin))
	if err != nil { e.X = math.NaN(); return }
	e.Push()
	e.X = clamp(r * 100)
	e.Flags.StackLift = true
}

func (e *Engine) SolvePV() {
	e.Push()
	e.X = clamp(tvmPV(e.FinI/100, e.FinN, e.FinPMT, e.FinFV, e.Flags.Begin))
	e.Flags.StackLift = true
}

func (e *Engine) SolvePMT() {
	e.Push()
	e.X = clamp(tvmPMT(e.FinI/100, e.FinN, e.FinPV, e.FinFV, e.Flags.Begin))
	e.Flags.StackLift = true
}

func (e *Engine) SolveFV() {
	e.Push()
	e.X = clamp(tvmFV(e.FinI/100, e.FinN, e.FinPV, e.FinPMT, e.Flags.Begin))
	e.Flags.StackLift = true
}

func tvmPV(i, n, pmt, fv float64, begin bool) float64 {
	if isZero(i) { return clamp(-(fv + pmt*n)) }
	b := 0.0; if begin { b = 1.0 }
	return clamp(-fv/math.Pow(1+i, n) - pmt*(1+i*b)*(1-1/math.Pow(1+i, n))/i)
}

func tvmPMT(i, n, pv, fv float64, begin bool) float64 {
	if isZero(i) { return clamp(-(fv + pv) / n) }
	b := 0.0; if begin { b = 1.0 }
	return clamp(-(pv + fv/math.Pow(1+i, n)) * i / ((1 + i*b) * (1 - 1/math.Pow(1+i, n))))
}

func tvmFV(i, n, pv, pmt float64, begin bool) float64 {
	if isZero(i) { return clamp(-(pv + pmt*n)) }
	b := 0.0; if begin { b = 1.0 }
	return clamp(-pv*math.Pow(1+i, n) - pmt*(1+i*b)*(math.Pow(1+i, n)-1)/i)
}

func tvmEq(i, n, pv, pmt, fv float64, begin bool) float64 {
	b := 0.0; if begin { b = 1.0 }
	if isZero(i) { return pv + fv + pmt*n*(1+i*b) }
	return pv*math.Pow(1+i, n) + pmt*(1+i*b)*(math.Pow(1+i, n)-1)/i + fv
}

func tvmEqDeriv(i, n, pv, pmt, fv float64, begin bool) float64 {
	b := 0.0; if begin { b = 1.0 }
	if isZero(i) { return pv*n + pmt*n*n/2 + pmt*n + fv*n }
	p := math.Pow(1+i, n-1)
	return pv*n*p + pmt*(1+i*b)*(n*p*(1+i) - (math.Pow(1+i, n)-1)/i) / (1 + i)
}

func (e *Engine) cashflows() []float64 {
	cfs := make([]float64, 0, e.FinCfCnt+1)
	cfs = append(cfs, e.FinCF0)
	cfs = append(cfs, e.FinCFj[:e.FinCfCnt]...)
	return cfs
}

func (e *Engine) ComputeNPV() {
	e.X = clamp(NPV(e.X/100, e.cashflows()...))
	e.Flags.StackLift = true
}

func (e *Engine) ComputeIRR() {
	r, _ := IRR(e.cashflows()...)
	e.result(r * 100)
}

func (e *Engine) Amortize() {
	n := e.X; i := e.FinI/100; pv := e.FinPV; pmt := e.FinPMT; begin := e.Flags.Begin
	totalInt, totalPrin := 0.0, 0.0
	for k := 1; k <= int(n); k++ {
		interest := pv * i
		if begin { interest = 0; pv -= pmt; if pv < 0 { pv = 0 } } else { pv -= pmt - interest }
		totalInt += interest; totalPrin += pmt - interest
	}
	e.AmortInt = clamp(totalInt); e.AmortPrin = clamp(totalPrin); e.AmortN = n
	e.X = e.AmortInt; e.Push(); e.X = e.AmortPrin; e.Y = e.AmortInt
	e.Flags.StackLift = true
}

func (e *Engine) BondPrice() {
	ytm := e.Y/100; coupon := e.X; n := e.FinN; price := 0.0
	for k := 1.0; k <= n; k++ { price += coupon/2 / math.Pow(1+ytm/2, k) }
	price += 100 / math.Pow(1+ytm/2, n)
	e.LastX = e.X; e.X = clamp(price); e.Tuck(); e.Flags.StackLift = true
}

func (e *Engine) BondYield() {
	price := e.X; coupon := e.Y; n := e.FinN; guess := coupon / price
	for iter := 0; iter < 100; iter++ {
		p := 0.0
		for k := 1.0; k <= n; k++ { p += coupon/2 / math.Pow(1+guess/2, k) }
		p += 100 / math.Pow(1+guess/2, n)
		f := p - price; d := 0.0
		for k := 1.0; k <= n; k++ { d -= float64(k)*coupon / (4*math.Pow(1+guess/2, k+1)) }
		d -= n*100 / (4*math.Pow(1+guess/2, n+1))
		if isZero(d) { break }
		next := guess - f/d
		if math.Abs(next-guess) < 1e-10 {
			e.LastX = e.X; e.X = clamp(next*100); e.Tuck(); e.Flags.StackLift = true; return
		}
		guess = next
	}
	e.LastX = e.X; e.X = clamp(guess*100); e.Tuck(); e.Flags.StackLift = true
}

func (e *Engine) depResult(dep, remaining float64) {
	e.LastX = e.X; e.X = clamp(dep); e.Y = clamp(remaining); e.Flags.StackLift = true
}

func (e *Engine) DepreciationSL() {
	dep, rem := DepSL(e.FinPV, e.FinFV, e.FinN, e.X)
	if math.IsNaN(dep) { e.X = math.NaN(); return }
	e.depResult(dep, rem)
}

func (e *Engine) DepreciationSOYD() {
	dep, rem := DepSOYD(e.FinPV, e.FinFV, e.FinN, e.X)
	if math.IsNaN(dep) { e.X = math.NaN(); return }
	e.depResult(dep, rem)
}

func (e *Engine) DepreciationDB() {
	dep, rem := DepDB(e.FinPV, e.FinFV, e.FinN, e.X, e.FinI)
	if math.IsNaN(dep) { e.X = math.NaN(); return }
	e.depResult(dep, rem)
}

func (e *Engine) result(v float64) {
	if e.Flags.StackLift { e.Push() }
	e.X = clamp(v)
	e.Flags.StackLift = true
}
