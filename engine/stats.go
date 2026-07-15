package engine

import "math"

func (e *Engine) statAdd() {
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
	// Weighted mean accumulates separately
	e.StatsWx += 1
	e.StatsWxx += x
	e.tuck()
}

func (e *Engine) statMeanX() {
	if e.StatsN == 0 {
		e.X = math.NaN()
		return
	}
	e.result(e.StatsSx / float64(e.StatsN))
}

func (e *Engine) statMeanY() {
	if e.StatsN == 0 {
		e.X = math.NaN()
		return
	}
	e.result(e.StatsSy / float64(e.StatsN))
}

func (e *Engine) statMeanW() {
	if e.StatsWx == 0 {
		e.X = math.NaN()
		return
	}
	e.result(e.StatsWxx / e.StatsWx)
}

func (e *Engine) statSDev() {
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

func (e *Engine) statLinEst() {
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

func (e *Engine) clearStats() {
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

func (e *Engine) clearFin() {
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

func (e *Engine) clearReg() {
	for i := range e.Mem {
		e.Mem[i] = 0
	}
	e.clearStats()
	e.clearFin()
	e.Flags.StackLift = true
}

func (e *Engine) clearPgm() {
	e.Program = [200]Instruction{}
	e.PgmLen = 0
	e.PgmPC = 0
}

func (e *Engine) clearPrefix() {
	e.Flags.Prefix = ""
}

func (e *Engine) cmpLE() {
	e.LastX = e.X
	if e.X <= e.Y {
		e.PgmPC++
	}
	e.tuck()
	e.Flags.StackLift = true
}

func (e *Engine) cmpEQ() {
	e.LastX = e.X
	if isZero(e.X) {
		e.PgmPC++
	}
	e.tuck()
	e.Flags.StackLift = true
}

func (e *Engine) sst() {
	if e.PgmPC < e.PgmLen {
		e.PgmPC++
	}
}

func (e *Engine) bst() {
	if e.PgmPC > 0 {
		e.PgmPC--
	}
}

func (e *Engine) gto(line int) {
	if line >= 0 && line < e.PgmLen {
		e.PgmPC = line
	}
}

func (e *Engine) runStop() {
	if e.Flags.Running {
		e.Flags.Running = false
		return
	}
	e.Flags.Running = true
	e.PgmPC = 0
	for e.Flags.Running && e.PgmPC < e.PgmLen {
		inst := e.Program[e.PgmPC]
		e.PgmPC++
		e.Step(inst.Op, inst.Arg, inst.ArgS)
	}
	e.Flags.Running = false
}
