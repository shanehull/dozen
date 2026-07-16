package engine

import "math"

func (e *Engine) StatAdd() {
	y := e.X
	x := e.Y
	e.StatsN++
	e.StatsSx += x
	e.StatsSy += y
	e.StatsSxx += x * x
	e.StatsSyy += y * y
	e.StatsSxy += x * y
	e.StatsLx = x
	e.StatsLy = y
	e.StatsWx += 1
	e.StatsWxx += x
	e.tuck()
}

func (e *Engine) MeanX() {
	if e.StatsN == 0 {
		e.X = math.NaN()
		return
	}
	e.result(e.StatsSx / float64(e.StatsN))
}

func (e *Engine) MeanY() {
	if e.StatsN == 0 {
		e.X = math.NaN()
		return
	}
	e.result(e.StatsSy / float64(e.StatsN))
}

func (e *Engine) WeightedMean() {
	if e.StatsWx == 0 {
		e.X = math.NaN()
		return
	}
	e.result(e.StatsWxx / e.StatsWx)
}

func (e *Engine) SDev() {
	if e.StatsN < 2 {
		e.X = math.NaN()
		return
	}
	n := float64(e.StatsN)
	v := (e.StatsSxx - e.StatsSx*e.StatsSx/n) / (n - 1)
	if v < 0 {
		v = 0
	}
	e.result(math.Sqrt(v))
}

func (e *Engine) LinEst() {
	if e.StatsN < 2 {
		e.X = math.NaN()
		return
	}
	n := float64(e.StatsN)
	denom := n*e.StatsSxx - e.StatsSx*e.StatsSx
	if isZero(denom) {
		e.X = math.NaN()
		return
	}
	slope := (n*e.StatsSxy - e.StatsSx*e.StatsSy) / denom
	intercept := (e.StatsSy - slope*e.StatsSx) / n

	numR := n*e.StatsSxy - e.StatsSx*e.StatsSy
	denR := math.Sqrt((n*e.StatsSxx - e.StatsSx*e.StatsSx) * (n*e.StatsSyy - e.StatsSy*e.StatsSy))
	r := 0.0
	if !isZero(denR) {
		r = numR / denR
	}

	if e.Flags.StackLift {
		e.push()
	}
	e.X = clamp(intercept)
	e.Y = clamp(r)
	e.Flags.StackLift = true
}

func (e *Engine) ClearStats() {
	e.StatsN = 0
	e.StatsSx = 0
	e.StatsSy = 0
	e.StatsSxx = 0
	e.StatsSyy = 0
	e.StatsSxy = 0
	e.StatsLx = 0
	e.StatsLy = 0
	e.StatsWx = 0
	e.StatsWxx = 0
}

func (e *Engine) ClearFin() {
	e.FinN = 0
	e.FinI = 0
	e.FinPV = 0
	e.FinPMT = 0
	e.FinFV = 0
	e.FinCF0 = 0
	e.FinCfCnt = 0
	e.AmortN = 0
	e.AmortInt = 0
	e.AmortPrin = 0
	e.Flags.StackLift = false
}

func (e *Engine) ClearReg() {
	for i := range e.Mem {
		e.Mem[i] = 0
	}
	e.ClearStats()
	e.ClearFin()
	e.Flags.StackLift = true
}

func (e *Engine) ClearPgm() {
	e.Program = [200]Instruction{}
	e.PgmLen = 0
	e.PgmPC = 0
}

func (e *Engine) SST() {
	if e.PgmPC < e.PgmLen {
		e.PgmPC++
	}
}

func (e *Engine) BST() {
	if e.PgmPC > 0 {
		e.PgmPC--
	}
}

func (e *Engine) Goto(line int) {
	if line >= 0 && line < e.PgmLen {
		e.PgmPC = line
	}
}
