package engine

import "math"

const MaxDigits = 10

type AngleMode int

const (
	Deg AngleMode = iota
	Rad
	Grad
)

type DisplayMode int

const (
	Fix DisplayMode = iota
	Sci
	Eng
)

type Engine struct {
	Stack
	Registers
	Mem     [20]float64
	Flags   Flags
	LastX   float64
	Display Display
	entry   numberEntry

	StatsN   int
	StatsSx  float64
	StatsSy  float64
	StatsSxx float64
	StatsSyy float64
	StatsSxy float64
	StatsLx  float64
	StatsLy  float64
	StatsWx  float64
	StatsWxx float64

	Program [200]Instruction
	PgmLen  int
	PgmPC   int
}

type numberEntry struct {
	buf       string
	hasDigits bool
	hasDot    bool
	hasEex    bool
	expSign   float64
	expDigits string
	inExp     bool
}

type Instruction struct {
	Op    string
	Arg   float64
	ArgS  string
	IsPfx bool
}

type Flags struct {
	StackLift   bool
	Begin       bool
	Dmy         bool
	Angle       AngleMode
	DispMode    DisplayMode
	DispDigits  int
	Prefix      string // "f", "g", "STO", "RCL", ""
	Running     bool
	ProgramMode bool
}

func New() *Engine {
	e := &Engine{}
	e.Flags.DispDigits = 2
	e.Flags.Angle = Deg
	return e
}

func (e *Engine) Step(op string, arg float64, argS string) {
	isPfx := op == "f" || op == "g" || op == "STO" || op == "RCL"
	if e.Flags.ProgramMode && !e.Flags.Running && e.PgmLen < 200 && op != "P/R" && !isPfx {
		e.Program[e.PgmLen] = Instruction{Op: op, Arg: arg, ArgS: argS, IsPfx: false}
		e.PgmLen++
		return
	}

	entryOps := map[string]bool{"0": true, "1": true, "2": true, "3": true, "4": true, "5": true, "6": true, "7": true, "8": true, "9": true, ".": true, "CHS": true, "EEX": true}
	if !entryOps[op] && e.entry.hasDigits {
		e.X = e.parseEntry()
		e.commitEntry()
		e.Flags.StackLift = true
	}
	switch {
	case e.Flags.Prefix == "f":
		e.prefixedF(op, arg, argS)
	case e.Flags.Prefix == "g":
		e.prefixedG(op, arg, argS)
	case e.Flags.Prefix == "STO":
		e.store(arg, argS)
	case e.Flags.Prefix == "RCL":
		e.recall(arg, argS)
	default:
		e.unprefixed(op, arg, argS)
	}
}

func (e *Engine) unprefixed(op string, arg float64, argS string) {
	switch op {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		e.digit(int(arg))
	case ".":
		e.decimal()
	case "CHS":
		e.chs()
	case "EEX":
		e.eex()
	case "ENTER":
		e.enter()
	case "CLx":
		e.clx()
	case "+":
		e.add()
	case "−":
		e.sub()
	case "×":
		e.mul()
	case "÷":
		e.div()
	case "NPV":
		e.finNPV()
	case "IRR":
		e.finIRR()
	case "PRICE":
		e.bondPrice()
	case "YTM":
		e.bondYield()
	case "SL":
		e.depSL()
	case "SOYD":
		e.depSOYD()
	case "DB":
		e.depDB()
	case "AMORT":
		e.finAmort()
	case "INT":
		e.finInt()
	case "Σ+":
		e.statAdd()
	case "x̄":
		e.statMeanX()
	case "ŷ":
		e.statMeanY()
	case "s":
		e.statSDev()
	case "x̄w":
		e.statMeanW()
	case "ŷ,r":
		e.statLinEst()
	case "SIN":
		e.sin()
	case "COS":
		e.cos()
	case "TAN":
		e.tan()
	case "LN":
		e.ln()
	case "LOG":
		e.log()
	case "DYS":
		e.dateDys()
	case "DATE":
		e.dateDate()
	case "π":
		e.pi()
	case "n!":
		e.fact()
	case "%CHG":
		e.pctChg()
	case "%T":
		e.pctTotal()
	case "%":
		e.pct()
	case "√x":
		e.sqrt()
	case "yˣ":
		e.yPowX()
	case "1/x":
		e.recip()
	case "|x|":
		e.abs()
	case "INTG":
		e.intg()
	case "FRAC":
		e.frac()
	case "ex":
		e.exp()
	case "10x":
		e.exp10()
	case "R↓":
		e.rollDown()
	case "R↑":
		e.rollUp()
	case "x↔y":
		e.xy()
	case "LSTx":
		e.lastX()
	case "n":
		if e.Flags.StackLift {
			e.FinN = e.X
			e.Flags.StackLift = false
		} else {
			e.tvmN()
		}
	case "i":
		if e.Flags.StackLift {
			e.FinI = e.X
			e.Flags.StackLift = false
		} else {
			e.tvmI()
		}
	case "PV":
		if e.Flags.StackLift {
			e.FinPV = e.X
			e.Flags.StackLift = false
		} else {
			e.tvmPV()
		}
	case "PMT":
		if e.Flags.StackLift {
			e.FinPMT = e.X
			e.Flags.StackLift = false
		} else {
			e.tvmPMT()
		}
	case "FV":
		if e.Flags.StackLift {
			e.FinFV = e.X
			e.Flags.StackLift = false
		} else {
			e.tvmFV()
		}
	case "R/S":
		e.runStop()
	case "PSE":
		// pause - noop
	case "SST":
		e.sst()
	case "BST":
		e.bst()
	case "f":
		e.Flags.Prefix = "f"
	case "g":
		e.Flags.Prefix = "g"
	case "STO":
		e.Flags.Prefix = "STO"
	case "RCL":
		e.Flags.Prefix = "RCL"
	}
}

func (e *Engine) prefixedF(op string, arg float64, argS string) {
	switch op {
	case "n":
		e.finN()
	case "i":
		e.finI()
	case "PV":
		e.finPV()
	case "PMT":
		e.finPMT()
	case "FV":
		e.finFV()
	case "NPV":
		e.finNPV()
	case "IRR":
		e.finIRR()
	case "AMORT":
		e.finAmort()
	case "INT":
		e.finInt()
	case "PRICE":
		e.bondPrice()
	case "YTM":
		e.bondYield()
	case "SL":
		e.depSL()
	case "SOYD":
		e.depSOYD()
	case "DB":
		e.depDB()
	case "Σ+":
		e.statAdd()
	case "x̄":
		e.statMeanX()
	case "ŷ":
		e.statMeanY()
	case "s":
		e.statSDev()
	case "x̄w":
		e.statMeanW()
	case "ŷ,r":
		e.statLinEst()
	case "R↓":
		e.rollDown()
	case "x↔y":
		e.xy()
	case "LSTx":
		e.lastX()
	case "CLEAR Σ":
		e.clearStats()
	case "CLEAR FIN":
		e.clearFin()
	case "CLEAR REG":
		e.clearReg()
	case "CLEAR PRGM":
		e.clearPgm()
	case "CLEAR PREFIX":
		e.clearPrefix()
	case "P/R":
		e.Flags.ProgramMode = !e.Flags.ProgramMode
	case "SIN":
		e.sin()
	case "COS":
		e.cos()
	case "TAN":
		e.tan()
	case "SIN⁻¹":
		e.asin()
	case "COS⁻¹":
		e.acos()
	case "TAN⁻¹":
		e.atan()
	case "RAD":
		e.toRad()
	case "DEG":
		e.toDeg()
	case "R↑":
		e.rollUp()
	case "→H.MS":
		e.toHMS()
	case "→H":
		e.toH()
	case "π":
		e.pi()
	case "%":
		e.pct()
	case "%CHG":
		e.pctChg()
	case "%T":
		e.pctTotal()
	case "LN":
		e.ln()
	case "LOG":
		e.log()
	case "eˣ":
		e.exp()
	case "10ˣ":
		e.exp10()
	case "yˣ":
		e.yPowX()
	case "1/x":
		e.recip()
	case "√x":
		e.sqrt()
	case "x²":
		e.sqr()
	case "|x|":
		e.abs()
	case "INTG":
		e.intg()
	case "FRAC":
		e.frac()
	case "n!":
		e.fact()
	case "→R":
		e.toRect()
	case "→P":
		e.toPolar()
	case "D.MY":
		e.Flags.Dmy = !e.Flags.Dmy
	case "DYS":
		e.dateDys()
	case "DATE":
		e.dateDate()
	case "BEG":
		e.Flags.Begin = !e.Flags.Begin
	case "END":
		e.Flags.Begin = false
	case "FIX":
		e.flagsDisp(Fix, int(arg))
	case "SCI":
		e.flagsDisp(Sci, int(arg))
	case "ENG":
		e.flagsDisp(Eng, int(arg))
	case "GRAD":
		e.Flags.Angle = Grad
	default:
		e.Flags.Prefix = ""
		return
	}
	e.Flags.Prefix = ""
}

func (e *Engine) prefixedG(op string, arg float64, argS string) {
	switch op {
	case "DEG":
		e.Flags.Angle = Deg
	case "RAD":
		e.Flags.Angle = Rad
	case "GRAD":
		e.Flags.Angle = Grad
	case "√x":
		e.sqrt()
	case "x²":
		e.sqr()
	case "1/x":
		e.recip()
	case "LN":
		e.ln()
	case "%":
		e.pct()
	case "%T":
		e.pctTotal()
	case "Δ%":
		e.pctChg()
	case "INTG":
		e.intg()
	case "FRAC":
		e.frac()
	case "x≤y":
		e.cmpLE()
	case "x=0":
		e.cmpEQ()
	case "LSTx":
		e.lastX()
	case "x↔y":
		e.xy()
	case "R↓":
		e.rollDown()
	case "R↑":
		e.rollUp()
	case "PSE": // pause - noop in emulator
	case "NOP": // noop
	case "→H.MS":
		e.toHMS()
	case "→H":
		e.toH()
	case "→R":
		e.toRect()
	case "→P":
		e.toPolar()
	case "D.MY":
		e.Flags.Dmy = !e.Flags.Dmy
	case "DYS":
		e.dateDys()
	case "DATE":
		e.dateDate()
	case "BEG":
		e.Flags.Begin = true
	case "END":
		e.Flags.Begin = false
	case "CLEAR Σ":
		e.clearStats()
	case "CLEAR FIN":
		e.clearFin()
	case "CLEAR REG":
		e.clearReg()
	case "CLEAR PRGM":
		e.clearPgm()
	case "CLEAR PREFIX":
		e.clearPrefix()
	case "FIX":
		e.flagsDisp(Fix, int(arg))
	case "SCI":
		e.flagsDisp(Sci, int(arg))
	case "ENG":
		e.flagsDisp(Eng, int(arg))
	case "MEM":
		e.memStatus()
	case "P/R":
		e.Flags.ProgramMode = !e.Flags.ProgramMode
	case "CF0":
		e.FinCF0 = e.X
	case "CFj":
		if e.FinCfCnt < 10 {
			e.FinCFj[e.FinCfCnt] = e.X
			e.FinCfCnt++
		}
	case "Nj":
		n := int(e.X)
		if n > 0 && e.FinCfCnt > 0 {
			e.FinNj[e.FinCfCnt-1] = n
		}
	case "SST":
		e.sst()
	case "BST":
		e.bst()
	case "GTO":
		e.gto(int(arg))
	case "R/S":
		e.runStop()
	default:
		e.Flags.Prefix = ""
		return
	}
	e.Flags.Prefix = ""
}

func (e *Engine) flagsDisp(m DisplayMode, digits int) {
	e.Flags.DispMode = m
	e.Flags.DispDigits = digits
}

func (e *Engine) memStatus() {
	// Returns used program steps - sets X to available/unused
	used := e.PgmLen
	e.push()
	e.X = float64(200 - used)
}

// result places a computed value into X following HP-12C stack semantics: the
// previous X is lifted into Y (when stack lift is enabled) and the result
// overwrites X. Stack lift is then enabled so a subsequently keyed number
// lifts the result up. Use this for functions that PRODUCE a value from the
// financial/statistics registers (TVM, IRR, means, etc.) rather than from X.
func (e *Engine) result(v float64) {
	if e.Flags.StackLift {
		e.push()
	}
	e.X = clamp(v)
	e.Flags.StackLift = true
}

func clamp(val float64) float64 {
	if math.IsInf(val, 0) || math.IsNaN(val) {
		return val
	}
	if math.Abs(val) > 9.999999999e99 {
		return math.Copysign(9.999999999e99, val)
	}
	return val
}

func isZero(f float64) bool {
	return math.Abs(f) < 1e-12
}
