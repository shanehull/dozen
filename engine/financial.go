package engine

import "math"

func (e *Engine) SetN(n float64)     { e.FinN = n }
func (e *Engine) SetI(i float64)     { e.FinI = i }
func (e *Engine) SetPV(pv float64)   { e.FinPV = pv }
func (e *Engine) SetPMT(pmt float64) { e.FinPMT = pmt }
func (e *Engine) SetFV(fv float64)   { e.FinFV = fv }

func (e *Engine) SolveN() { e.solve(math.NaN(), e.FinI/100, e.FinPV, e.FinPMT, e.FinFV) }
func (e *Engine) SolveI() {
	v, err := solveTVM(e.FinN, math.NaN(), e.FinPV, e.FinPMT, e.FinFV, timing(e.Flags.Begin))
	if err != nil {
		e.X = math.NaN()
		return
	}
	e.Push()
	e.X = clamp(v * 100)
	e.Flags.StackLift = true
}
func (e *Engine) SolvePV()  { e.solve(e.FinN, e.FinI/100, math.NaN(), e.FinPMT, e.FinFV) }
func (e *Engine) SolvePMT() { e.solve(e.FinN, e.FinI/100, e.FinPV, math.NaN(), e.FinFV) }
func (e *Engine) SolveFV()  { e.solve(e.FinN, e.FinI/100, e.FinPV, e.FinPMT, math.NaN()) }

func (e *Engine) solve(n, i, pv, pmt, fv float64) {
	v, err := solveTVM(n, i, pv, pmt, fv, timing(e.Flags.Begin))
	if err != nil {
		e.X = math.NaN()
		return
	}
	e.Push()
	e.X = clamp(v)
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
	return clamp(-fv/math.Pow(1+i, n) - pmt*(1+i*b)*(1-1/math.Pow(1+i, n))/i)
}

func tvmPMT(i, n, pv, fv float64, begin bool) float64 {
	if isZero(i) {
		return clamp(-(fv + pv) / n)
	}
	b := 0.0
	if begin {
		b = 1.0
	}
	return clamp(-(pv + fv/math.Pow(1+i, n)) * i / ((1 + i*b) * (1 - 1/math.Pow(1+i, n))))
}

func tvmFV(i, n, pv, pmt float64, begin bool) float64 {
	if isZero(i) {
		return clamp(-(pv + pmt*n))
	}
	b := 0.0
	if begin {
		b = 1.0
	}
	return clamp(-pv*math.Pow(1+i, n) - pmt*(1+i*b)*(math.Pow(1+i, n)-1)/i)
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
	return pv*n*p + pmt*(1+i*b)*(n*p*(1+i)-(math.Pow(1+i, n)-1)/i)/(1+i)
}

func (e *Engine) cashflows() []float64 {
	cfs := make([]float64, 0, e.FinCfCnt+1)
	cfs = append(cfs, e.FinCF0)
	cfs = append(cfs, e.FinCFj[:e.FinCfCnt]...)
	return cfs
}

func (e *Engine) NPV() {
	e.X = clamp(npv(e.X/100, e.cashflows()...))
	e.Flags.StackLift = true
}

func (e *Engine) IRR() {
	r, _ := irr(e.cashflows()...)
	e.result(r * 100)
}

func (e *Engine) Amortize() {
	totalInt, totalPrin := amort(e.FinI/100, e.X, e.FinPV, e.FinPMT, bool(e.Flags.Begin))
	e.AmortInt = totalInt
	e.AmortPrin = totalPrin
	e.AmortN = e.X
	e.X = e.AmortInt
	e.Push()
	e.X = e.AmortPrin
	e.Y = e.AmortInt
	e.Flags.StackLift = true
}

func (e *Engine) BondPrice() {
	price := bondPrice(e.Y/100, e.X, e.FinN)
	e.LastX = e.X
	e.X = price
	e.Tuck()
	e.Flags.StackLift = true
}

func (e *Engine) BondYield() {
	ytm, err := bondYield(e.X, e.Y, e.FinN)
	if err != nil {
		e.X = math.NaN()
		return
	}
	e.LastX = e.X
	e.X = clamp(ytm * 100)
	e.Tuck()
	e.Flags.StackLift = true
}

func (e *Engine) depResult(dep, remaining float64) {
	e.LastX = e.X
	e.X = clamp(dep)
	e.Y = clamp(remaining)
	e.Flags.StackLift = true
}

func (e *Engine) DepreciationSL() {
	dep, rem := depSL(e.FinPV, e.FinFV, e.FinN, e.X)
	if math.IsNaN(dep) {
		e.X = math.NaN()
		return
	}
	e.depResult(dep, rem)
}

func (e *Engine) DepreciationSOYD() {
	dep, rem := depSOYD(e.FinPV, e.FinFV, e.FinN, e.X)
	if math.IsNaN(dep) {
		e.X = math.NaN()
		return
	}
	e.depResult(dep, rem)
}

func (e *Engine) DepreciationDB() {
	dep, rem := depDB(e.FinPV, e.FinFV, e.FinN, e.X, e.FinI)
	if math.IsNaN(dep) {
		e.X = math.NaN()
		return
	}
	e.depResult(dep, rem)
}

func (e *Engine) result(v float64) {
	if e.Flags.StackLift {
		e.Push()
	}
	e.X = clamp(v)
	e.Flags.StackLift = true
}
