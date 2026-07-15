package engine

import "math"

func (e *Engine) digit(d int) {
	if !e.entry.hasDigits && e.Flags.StackLift {
		e.push()
	}
	if !e.entry.hasDigits {
		e.entry.buf = ""
		e.entry.hasDot = false
	}
	e.entry.hasDigits = true
	e.entry.buf += string(rune('0' + d))
	e.X = e.parseEntry()
}

func (e *Engine) decimal() {
	if !e.entry.hasDigits {
		if e.Flags.StackLift {
			e.push()
		}
		e.entry.buf = "0"
		e.entry.hasDigits = true
	}
	if !e.entry.hasDot {
		e.entry.hasDot = true
		e.entry.buf += "."
	}
	e.X = e.parseEntry()
}

func (e *Engine) chs() {
	if e.entry.hasDigits {
		if len(e.entry.buf) > 0 && e.entry.buf[0] == '-' {
			e.entry.buf = e.entry.buf[1:]
		} else {
			e.entry.buf = "-" + e.entry.buf
		}
		e.X = e.parseEntry()
	} else {
		e.X = -e.X
	}
}

func (e *Engine) eex() {
	if !e.entry.hasDigits {
		if e.Flags.StackLift {
			e.push()
		}
		e.entry.buf = "1"
		e.entry.hasDigits = true
	}
	if !e.entry.inExp {
		e.entry.inExp = true
		e.entry.buf += "e"
	}
}

func (e *Engine) parseEntry() float64 {
	s := e.entry.buf
	if s == "" || s == "-" || s == "." || s == "-." {
		return 0
	}
	result := 0.0
	sign := 1.0
	i := 0
	if s[0] == '-' {
		sign = -1
		i++
	}
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		result = result*10 + float64(s[i]-'0')
		i++
	}
	if i < len(s) && s[i] == '.' {
		i++
		div := 10.0
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			result += float64(s[i]-'0') / div
			div *= 10
			i++
		}
	}
	result *= sign
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		i++
		esign := 1.0
		if i < len(s) && s[i] == '-' {
			esign = -1
			i++
		} else if i < len(s) && s[i] == '+' {
			i++
		}
		exp := 0.0
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			exp = exp*10 + float64(s[i]-'0')
			i++
		}
		result *= math.Pow(10, esign*exp)
	}
	return result
}

func (e *Engine) commitEntry() {
	e.entry = numberEntry{}
}

func (e *Engine) enter() {
	e.LastX = e.X
	e.push()
	e.Flags.StackLift = false
}

func (e *Engine) clx() {
	e.X = 0
}

func (e *Engine) add() {
	e.LastX = e.X
	e.X = clamp(e.Y + e.X)
	e.tuck()
	e.Flags.StackLift = true
}

func (e *Engine) sub() {
	e.LastX = e.X
	e.X = clamp(e.Y - e.X)
	e.tuck()
	e.Flags.StackLift = true
}

func (e *Engine) mul() {
	e.LastX = e.X
	e.X = clamp(e.Y * e.X)
	e.tuck()
	e.Flags.StackLift = true
}

func (e *Engine) div() {
	e.LastX = e.X
	if isZero(e.X) {
		e.X = math.Inf(1)
	} else {
		e.X = clamp(e.Y / e.X)
	}
	e.tuck()
	e.Flags.StackLift = true
}

func (e *Engine) xy() {
	e.LastX = e.X
	e.X, e.Y = e.Y, e.X
}

func (e *Engine) lastX() {
	if e.Flags.StackLift {
		e.push()
	}
	e.X = e.LastX
	e.Flags.StackLift = false
}

func (e *Engine) store(arg float64, argS string) {
	idx := int(arg)
	if idx >= 0 && idx < 20 {
		e.Mem[idx] = e.X
	}
	e.Flags.StackLift = false
	e.Flags.Prefix = ""
}

func (e *Engine) recall(arg float64, argS string) {
	idx := int(arg)
	if idx >= 0 && idx < 20 {
		if e.Flags.StackLift {
			e.push()
		}
		e.X = e.Mem[idx]
	}
	e.Flags.StackLift = false
	e.Flags.Prefix = ""
}

func (e *Engine) rollDown() {
	e.LastX = e.X
	e.Stack.rollDown()
}

func (e *Engine) yPowX() {
	e.LastX = e.X
	if isZero(e.Y) && e.X < 0 {
		e.X = math.Inf(1)
	} else if e.Y < 0 && e.X != math.Trunc(e.X) {
		e.X = math.NaN()
	} else {
		e.X = clamp(math.Pow(e.Y, e.X))
	}
	e.tuck()
	e.Flags.StackLift = true
}
