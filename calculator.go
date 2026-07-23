package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/shanehull/dozen/engine"
)

// ---- display ---------------------------------------------------------------

type Display struct {
	Mantissa string   `json:"mantissa"`
	Exponent string   `json:"exponent"`
	Sign     string   `json:"sign"`
	Flags    []string `json:"flags"`
}

// ---- service ---------------------------------------------------------------

type CalcService struct {
	e         *engine.Engine
	hasEntry  bool
	buf       string
	expBuf    string
	hasDot    bool
	hasSign   bool
	inExp     bool
	armed     string // "", "f", "g", "STO", "RCL"
	undoState *engine.EngineState
}

type KeyResult struct {
	Display Display  `json:"display"`
	StackX  float64  `json:"stackX"`
	StackY  float64  `json:"stackY"`
	StackZ  float64  `json:"stackZ"`
	StackT  float64  `json:"stackT"`
	LastX   float64  `json:"lastX"`
	Flags   []string `json:"flags"`
}

type KeyInput struct {
	Op   string  `json:"op"`
	Arg  float64 `json:"arg"`
	ArgS string  `json:"argS"`
}

func NewCalcService() *CalcService {
	return &CalcService{e: engine.New()}
}

// ---- key dispatch ----------------------------------------------------------

var tvms = map[string]bool{
	"n": true, "i": true, "PV": true, "PMT": true, "FV": true,
}

func (c *CalcService) PressKey(input KeyInput) KeyResult {
	op, arg := input.Op, input.Arg
	dig, isDigit := 0, false
	if len(op) == 1 && op[0] >= '0' && op[0] <= '9' {
		dig = int(op[0] - '0')
		isDigit = true
	}

	// snapshot before every state-changing operation for single-level undo
	if op != "f" && op != "g" && op != "STO" && op != "RCL" && !isDigit && op != "." && op != "CHS" && op != "EEX" && op != "ENTER" && op != "CLx" {
		s := c.e.Snapshot()
		c.undoState = &s
	}

	switch {
	case c.armed == "f":
		c.fPrefixed(op, arg, input.ArgS)
	case c.armed == "g":
		c.gPrefixed(op, arg, input.ArgS)
	case c.armed == "STO":
		c.finishEntry()
		c.e.Store(int(arg))
		c.armed = ""
	case c.armed == "RCL":
		c.finishEntry()
		c.e.Recall(int(arg))
		c.armed = ""
	case op == "f" || op == "g":
		c.armed = op
		return c.state()
	case op == "STO" || op == "RCL":
		c.armed = op
		return c.state()
	case isDigit:
		c.enterDigit(dig)
	case op == ".":
		c.enterDecimal()
	case op == "CHS":
		c.enterChs()
	case op == "EEX":
		c.enterEex()
	case op == "ENTER":
		c.finishEntry()
		c.e.Enter()
	case op == "CLx":
		c.e.Clx()
		c.clearEntry()
	case op == "+":
		c.finishEntry()
		c.e.Add()
	case op == "−":
		c.finishEntry()
		c.e.Sub()
	case op == "×":
		c.finishEntry()
		c.e.Mul()
	case op == "÷":
		c.finishEntry()
		c.e.Div()
	case op == "x↔y":
		c.finishEntry()
		c.e.XY()
	case op == "R↓":
		c.finishEntry()
		c.e.RollDown()
	case op == "R↑":
		c.finishEntry()
		c.e.RollUp()
	case op == "LSTx":
		c.finishEntry()
		c.e.LastXRecall()
	case op == "f" || op == "g":
		c.armed = op
		return c.state()
	case op == "STO" || op == "RCL":
		c.armed = op
		return c.state()
	case c.armed == "f":
		c.fPrefixed(op, arg, input.ArgS)
	case c.armed == "g":
		c.gPrefixed(op, arg, input.ArgS)
	case c.armed == "STO":
		c.finishEntry()
		c.e.Store(int(arg))
		c.armed = ""
	case c.armed == "RCL":
		c.e.Recall(int(arg))
		c.armed = ""
	case tvms[op]:
		if c.hasEntry || c.e.Flags.StackLift {
			c.finishEntry()
			c.storeTVM(op)
		} else {
			c.solveTVM(op)
		}
	default:
		c.finishEntry()
		c.unprefixed(op, arg)
	}
	return c.state()
}

func (c *CalcService) storeTVM(op string) {
	switch op {
	case "n":
		c.e.SetN(c.e.X)
	case "i":
		c.e.SetI(c.e.X)
	case "PV":
		c.e.SetPV(c.e.X)
	case "PMT":
		c.e.SetPMT(c.e.X)
	case "FV":
		c.e.SetFV(c.e.X)
	}
	c.e.Flags.StackLift = false
}

func (c *CalcService) solveTVM(op string) {
	switch op {
	case "n":
		c.e.SolveN()
	case "i":
		c.e.SolveI()
	case "PV":
		c.e.SolvePV()
	case "PMT":
		c.e.SolvePMT()
	case "FV":
		c.e.SolveFV()
	}
}

func (c *CalcService) unprefixed(op string, arg float64) {
	switch op {
	case "NPV":
		c.e.NPV()
	case "IRR":
		c.e.IRR()
	case "PRICE":
		c.e.BondPrice()
	case "YTM":
		c.e.BondYield()
	case "SL":
		c.e.DepreciationSL()
	case "SOYD":
		c.e.DepreciationSOYD()
	case "DB":
		c.e.DepreciationDB()
	case "AMORT":
		c.e.Amortize()
	case "INT":
		// AmortInt is a field — just set X
		c.e.X = c.e.AmortInt
	case "Σ+":
		x := c.e.Y
		y := c.e.X
		c.e.Y = x
		c.e.X = y
		c.e.StatAdd()
	case "x̄":
		c.e.MeanX()
	case "ŷ":
		c.e.MeanY()
	case "s":
		c.e.SDev()
	case "x̄w":
		c.e.WeightedMean()
	case "ŷ,r":
		c.e.LinEst()
	case "SIN":
		c.e.Sin()
	case "COS":
		c.e.Cos()
	case "TAN":
		c.e.Tan()
	case "LN":
		c.e.Ln()
	case "LOG":
		c.e.Log()
	case "DYS":
		c.e.DaysBetween()
	case "DATE":
		c.e.DateAdd()
	case "π":
		c.e.Pi()
	case "n!":
		c.e.Fact()
	case "%CHG":
		c.e.PctChg()
	case "%T":
		c.e.PctTotal()
	case "%":
		c.e.Pct()
	case "√x":
		c.e.Sqrt()
	case "yˣ":
		c.e.YPowX()
	case "1/x":
		c.e.Recip()
	case "|x|":
		c.e.Abs()
	case "INTG":
		c.e.Intg()
	case "FRAC":
		c.e.Frac()
	case "ex":
		c.e.Exp()
	case "10x":
		c.e.Exp10()
	}
}

func (c *CalcService) fPrefixed(op string, arg float64, argS string) {
	defer func() { c.armed = "" }()
	switch op {
	case "CLEAR FIN":
		c.e.ClearFin()
	case "CLEAR REG":
		c.e.ClearReg()
	case "CLEAR Σ":
		c.e.ClearStats()
	case "CLEAR PRGM":
		c.e.ClearPgm()
	case "CLEAR PREFIX":
		c.armed = ""
	case "P/R":
		// program mode toggle — handled in gPrefixed for now
	case "FIX":
		dispDigits = int(arg)
		dispMode = 0
	case "SCI":
		dispDigits = int(arg)
		dispMode = 1
	case "ENG":
		dispDigits = int(arg)
		dispMode = 2
	case "DEG":
		c.e.Flags.Angle = engine.Deg
	case "RAD":
		c.e.Flags.Angle = engine.Rad
	case "GRAD":
		c.e.Flags.Angle = engine.Grad
	case "AMORT":
		c.e.Amortize()
	case "INT":
		c.e.X = c.e.AmortInt
	case "NPV":
		c.e.NPV()
	case "IRR":
		c.e.IRR()
	case "PRICE":
		c.e.BondPrice()
	case "YTM":
		c.e.BondYield()
	case "SL":
		c.e.DepreciationSL()
	case "SOYD":
		c.e.DepreciationSOYD()
	case "DB":
		c.e.DepreciationDB()
	case "R↓":
		c.e.RollDown()
	case "R↑":
		c.e.RollUp()
	case "LSTx":
		c.e.LastXRecall()
	case "SIN":
		c.e.Sin()
	case "COS":
		c.e.Cos()
	case "TAN":
		c.e.Tan()
	case "DYS":
		c.e.DaysBetween()
	case "DATE":
		c.e.DateAdd()
	case "%":
		c.e.Pct()
	case "%CHG":
		c.e.PctChg()
	case "%T":
		c.e.PctTotal()
	case "LN":
		c.e.Ln()
	case "LOG":
		c.e.Log()
	case "ex":
		c.e.Exp()
	case "10x":
		c.e.Exp10()
	case "√x":
		c.e.Sqrt()
	case "x²":
		c.e.Sqr()
	case "1/x":
		c.e.Recip()
	case "|x|":
		c.e.Abs()
	case "INTG":
		c.e.Intg()
	case "FRAC":
		c.e.Frac()
	case "n!":
		c.e.Fact()
	case "π":
		c.e.Pi()
	case "yˣ":
		c.e.YPowX()
	case "→R":
		c.e.ToRect()
	case "→P":
		c.e.ToPolar()
	case "→H.MS":
		c.e.ToHMS()
	case "→H":
		c.e.ToH()
	case "x≤y":
		// handled in calculator logic
	case "x=0":
		// handled in calculator logic
	case "MEM":
		// memory status — handled by engine
	case "RND":
		// round — nop for now
	default:
		_ = argS
	}
}

func (c *CalcService) gPrefixed(op string, arg float64, argS string) {
	defer func() { c.armed = "" }()
	switch op {
	case "DEG":
		c.e.Flags.Angle = engine.Deg
	case "RAD":
		c.e.Flags.Angle = engine.Rad
	case "GRAD":
		c.e.Flags.Angle = engine.Grad
	case "√x":
		c.e.Sqrt()
	case "x²":
		c.e.Sqr()
	case "1/x":
		c.e.Recip()
	case "LN":
		c.e.Ln()
	case "%":
		c.e.Pct()
	case "%T":
		c.e.PctTotal()
	case "Δ%":
		c.e.PctChg()
	case "INTG":
		c.e.Intg()
	case "FRAC":
		c.e.Frac()
	case "x≤y":
	case "x=0":
	case "LSTx":
		c.e.LastXRecall()
	case "x↔y":
		c.e.XY()
	case "R↓":
		c.e.RollDown()
	case "R↑":
		c.e.RollUp()
	case "PSE":
	case "NOP":
	case "→H.MS":
		c.e.ToHMS()
	case "→H":
		c.e.ToH()
	case "→R":
		c.e.ToRect()
	case "→P":
		c.e.ToPolar()
	case "D.MY":
		c.e.Flags.Dmy = !c.e.Flags.Dmy
	case "DYS":
		c.e.DaysBetween()
	case "DATE":
		c.e.DateAdd()
	case "BEG":
		c.e.Flags.Begin = true
	case "END":
		c.e.Flags.Begin = false
	case "CLEAR Σ":
		c.e.ClearStats()
	case "CLEAR FIN":
		c.e.ClearFin()
	case "CLEAR REG":
		c.e.ClearReg()
	case "CLEAR PRGM":
		c.e.ClearPgm()
	case "CLEAR PREFIX":
		c.armed = ""
	case "FIX":
		dispDigits = int(arg)
		dispMode = 0
	case "SCI":
		dispDigits = int(arg)
		dispMode = 1
	case "ENG":
		dispDigits = int(arg)
		dispMode = 2
	case "MEM":
		c.e.Push()
		c.e.X = float64(200 - c.e.PgmLen)
	case "P/R":
		// program mode toggle
	case "CF0":
		c.e.FinCF0 = c.e.X
	case "CFj":
		if c.e.FinCfCnt < 10 {
			c.e.FinCFj[c.e.FinCfCnt] = c.e.X
			c.e.FinCfCnt++
		}
	case "Nj":
		n := int(c.e.X)
		if n > 0 && c.e.FinCfCnt > 0 {
			c.e.FinNj[c.e.FinCfCnt-1] = n
		}
	case "SST":
	case "BST":
	case "GTO":
	case "R/S":
	case "+":
		c.finishEntry()
		c.e.LastXRecall()
	case "×":
		c.finishEntry()
		c.e.Sqr()
	case "÷":
		if c.undoState != nil {
			c.e.Restore(*c.undoState)
			c.undoState = nil
		}
	case "−":
		c.backspace()
	case "PMT":
		c.finishEntry()
		if c.e.FinCfCnt < 10 {
			c.e.FinCFj[c.e.FinCfCnt] = c.e.X
			c.e.FinCfCnt++
		}
	case "FV":
		c.finishEntry()
		n := int(c.e.X)
		if n > 0 && c.e.FinCfCnt > 0 {
			c.e.FinNj[c.e.FinCfCnt-1] = n
		}
	default:
		_ = argS
	}
}

// ---- entry management ------------------------------------------------------

func (c *CalcService) enterDigit(d int) {
	if !c.hasEntry && c.e.Flags.StackLift {
		c.e.Push()
		c.e.Flags.StackLift = false
	}
	if !c.hasEntry {
		c.buf = ""
		c.hasDot = false
		c.hasSign = false
	}
	if c.inExp {
		c.expBuf += string(rune('0' + d))
	} else {
		c.buf += string(rune('0' + d))
	}
	c.hasEntry = true
	c.e.X = c.parseBuf()
}

func (c *CalcService) enterDecimal() {
	if !c.hasEntry {
		if c.e.Flags.StackLift {
			c.e.Push()
			c.e.Flags.StackLift = false
		}
		c.buf = "0"
		c.hasEntry = true
	}
	if !c.hasDot {
		c.hasDot = true
		c.buf += "."
	}
	c.e.X = c.parseBuf()
}

func (c *CalcService) enterChs() {
	if c.hasEntry {
		if strings.HasPrefix(c.buf, "-") {
			c.buf = c.buf[1:]
		} else {
			c.buf = "-" + c.buf
		}
		c.hasSign = !c.hasSign
		c.e.X = c.parseBuf()
	} else {
		c.e.Chs()
	}
}

func (c *CalcService) enterEex() {
	if !c.hasEntry {
		if c.e.Flags.StackLift {
			c.e.Push()
			c.e.Flags.StackLift = false
		}
		c.buf = "1"
		c.hasEntry = true
	}
	c.inExp = true
}

func (c *CalcService) finishEntry() {
	if c.hasEntry {
		c.e.X = c.parseBuf()
		c.clearEntry()
		c.e.Flags.StackLift = true
	}
}

func (c *CalcService) clearEntry() {
	c.hasEntry = false
	c.buf = ""
	c.expBuf = ""
	c.hasDot = false
	c.hasSign = false
	c.inExp = false
}

func (c *CalcService) backspace() {
	if !c.hasEntry {
		// not in an active entry; start editing the current X value
		c.buf = strconv.FormatFloat(c.e.X, 'f', -1, 64)
		if len(c.buf) > 0 {
			c.buf = c.buf[:len(c.buf)-1]
		}
		c.hasEntry = true
	} else {
		if c.inExp && len(c.expBuf) > 0 {
			c.expBuf = c.expBuf[:len(c.expBuf)-1]
		} else if !c.inExp && len(c.buf) > 0 {
			c.buf = c.buf[:len(c.buf)-1]
		}
	}
	if c.buf == "" || c.buf == "-" || c.buf == "." || c.buf == "-." {
		c.clearEntry()
		c.e.X = 0
	} else {
		c.e.X = c.parseBuf()
	}
}

func (c *CalcService) parseBuf() float64 {
	s := c.buf
	if c.inExp && c.expBuf != "" {
		s += "e" + c.expBuf
	}
	if s == "" || s == "-" || s == "." || s == "-." {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

// ---- display formatting ----------------------------------------------------

var dispMode = 0 // 0=Fix, 1=Sci, 2=Eng
var dispDigits = 2

func (c *CalcService) format() Display {
	d := Display{}

	if c.hasEntry {
		d.Mantissa = trimZero(c.buf)
		if strings.HasPrefix(c.buf, "-") {
			d.Sign = "-"
		}
		d.Flags = c.annunciators()
		return d
	}

	x := c.e.X
	switch {
	case math.IsInf(x, 0):
		d.Mantissa = "9.999999 99"
	case math.IsNaN(x):
		d.Mantissa = " Error    "
	case x < -9.999999999e99:
		d.Mantissa = "-9.999999 99"
	case x > 9.999999999e99:
		d.Mantissa = " 9.999999 99"
	default:
		if x < 0 {
			d.Sign = "-"
			x = -x
		}
		switch dispMode {
		case 0:
			d.Mantissa = fmt.Sprintf("%.*f", dispDigits, x)
		case 1:
			d.Mantissa = formatSci(x, dispDigits)
		case 2:
			d.Mantissa = formatEng(x, dispDigits)
		}
	}
	d.Flags = c.annunciators()
	return d
}

func formatSci(x float64, d int) string {
	if x == 0 {
		return fmt.Sprintf("%.*f", d, 0.0) + "   00"
	}
	e := 0
	for x >= 10 {
		x /= 10
		e++
	}
	for x < 1 {
		x *= 10
		e--
	}
	return fmt.Sprintf("%.*f", d, x) + expStr(e)
}

func formatEng(x float64, d int) string {
	if x == 0 {
		return fmt.Sprintf("%.*f", d, 0.0) + "   00"
	}
	e := 0
	for x >= 1000 {
		x /= 1000
		e += 3
	}
	for x < 1 {
		x *= 1000
		e -= 3
	}
	return fmt.Sprintf("%.*f", d, x) + expStr(e)
}

func expStr(e int) string {
	if e >= 0 {
		return fmt.Sprintf("  %02d", e)
	}
	return fmt.Sprintf(" -%02d", -e)
}

func trimZero(s string) string {
	if len(s) <= 1 || s == "0" {
		return s
	}
	i := 0
	neg := s[0] == '-'
	if neg {
		i = 1
	}
	for i < len(s)-1 && s[i] == '0' && s[i+1] != '.' {
		i++
	}
	return s[i:]
}

func (c *CalcService) annunciators() []string {
	var a []string
	switch c.armed {
	case "f":
		a = append(a, "f")
	case "g":
		a = append(a, "g")
	}
	if c.e.Flags.Begin {
		a = append(a, "BEGIN")
	}
	if c.e.Flags.Dmy {
		a = append(a, "D.MY")
	}
	switch c.e.Flags.Angle {
	case engine.Rad:
		a = append(a, "RAD")
	case engine.Grad:
		a = append(a, "GRAD")
	}
	return a
}

// ---- state -----------------------------------------------------------------

func (c *CalcService) GetState() KeyResult { return c.state() }

func (c *CalcService) Save() string {
	b, _ := c.e.Snapshot().Marshal()
	return string(b)
}

func (c *CalcService) Load(j string) {
	s, err := engine.UnmarshalState([]byte(j))
	if err != nil {
		return
	}
	c.e.Restore(s)
}

func (c *CalcService) state() KeyResult {
	f := c.format()
	return KeyResult{
		Display: f,
		StackX:  sz(c.e.X),
		StackY:  sz(c.e.Y),
		StackZ:  sz(c.e.Z),
		StackT:  sz(c.e.T),
		LastX:   sz(c.e.LastX),
		Flags:   f.Flags,
	}
}

func sz(f float64) float64 {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0
	}
	return f
}
