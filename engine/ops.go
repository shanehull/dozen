package engine

import "math"

func (e *Engine) Enter() {
	e.LastX = e.X
	e.Push()
	e.Flags.StackLift = false
}

func (e *Engine) Clx() {
	e.X = 0
}

func (e *Engine) Chs() {
	e.X = -e.X
}

func (e *Engine) Add() {
	e.LastX = e.X
	e.X = clamp(e.Y + e.X)
	e.Tuck()
	e.Flags.StackLift = true
}

func (e *Engine) Sub() {
	e.LastX = e.X
	e.X = clamp(e.Y - e.X)
	e.Tuck()
	e.Flags.StackLift = true
}

func (e *Engine) Mul() {
	e.LastX = e.X
	e.X = clamp(e.Y * e.X)
	e.Tuck()
	e.Flags.StackLift = true
}

func (e *Engine) Div() {
	e.LastX = e.X
	if isZero(e.X) {
		e.X = math.Inf(1)
	} else {
		e.X = clamp(e.Y / e.X)
	}
	e.Tuck()
	e.Flags.StackLift = true
}

func (e *Engine) XY() {
	e.LastX = e.X
	e.X, e.Y = e.Y, e.X
}

func (e *Engine) LastXRecall() {
	if e.Flags.StackLift {
		e.Push()
	}
	e.X = e.LastX
	e.Flags.StackLift = false
}

func (e *Engine) Store(idx int) {
	if idx >= 0 && idx < 20 {
		e.Mem[idx] = e.X
	}
	e.Flags.StackLift = false
}

func (e *Engine) Recall(idx int) {
	if idx >= 0 && idx < 20 {
		if e.Flags.StackLift {
			e.Push()
		}
		e.X = e.Mem[idx]
	}
	e.Flags.StackLift = false
}

func (e *Engine) RollDown() {
	e.LastX = e.X
	e.Stack.rollDown()
}

func (e *Engine) RollUp() {
	e.LastX = e.X
	e.Stack.rollUp()
}

func (e *Engine) YPowX() {
	e.LastX = e.X
	if isZero(e.Y) && e.X < 0 {
		e.X = math.Inf(1)
	} else if e.Y < 0 && e.X != math.Trunc(e.X) {
		e.X = math.NaN()
	} else {
		e.X = clamp(math.Pow(e.Y, e.X))
	}
	e.Tuck()
	e.Flags.StackLift = true
}
