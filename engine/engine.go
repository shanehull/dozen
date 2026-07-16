package engine

import "math"

const MaxDigits = 10

type AngleMode int

const (
	Deg  AngleMode = iota
	Rad
	Grad
)

type Engine struct {
	LastX float64
	Mem   [20]float64

	StatsN, StatsSx, StatsSy, StatsSxx, StatsSyy, StatsSxy float64
	StatsLx, StatsLy, StatsWx, StatsWxx                     float64

	PgmLen int
	PgmPC  int
	Program [200]Instruction

	Flags Flags

	Stack
	Registers
}

type Instruction struct {
	Op    string
	Arg   float64
	ArgS  string
	IsPfx bool
}

type Flags struct {
	StackLift bool
	Begin     bool
	Dmy       bool
	Angle     AngleMode
}

func New() *Engine {
	return &Engine{}
}

func (e *Engine) Push() {
	e.Stack.push()
}

func (e *Engine) Tuck() {
	e.Stack.tuck()
}

func isZero(f float64) bool {
	return math.Abs(f) < 1e-12
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
