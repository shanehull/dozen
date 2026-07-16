package engine

import "math"

func (e *Engine) Sin() {
	e.X = clamp(math.Sin(e.angleInRads()))
}

func (e *Engine) Cos() {
	e.X = clamp(math.Cos(e.angleInRads()))
}

func (e *Engine) Tan() {
	e.X = clamp(math.Tan(e.angleInRads()))
}

func (e *Engine) Asin() {
	e.X = clamp(math.Asin(e.X))
	e.toCurAngle()
}

func (e *Engine) Acos() {
	e.X = clamp(math.Acos(e.X))
	e.toCurAngle()
}

func (e *Engine) Atan() {
	e.X = clamp(math.Atan(e.X))
	e.toCurAngle()
}

func (e *Engine) ToRad() {
	e.X = clamp(e.X * math.Pi / 180)
}

func (e *Engine) ToDeg() {
	e.X = clamp(e.X * 180 / math.Pi)
}

func (e *Engine) ToH() {
	d := e.X
	sign := 1.0
	if d < 0 {
		sign = -1
		d = -d
	}
	si := int(math.Round(d * 10000))
	deg := si / 10000
	min := (si % 10000) / 100
	sec := si % 100
	e.X = sign * (float64(deg) + float64(min)/60 + float64(sec)/3600)
}

func (e *Engine) ToHMS() {
	d := e.X
	sign := 1.0
	if d < 0 {
		sign = -1
		d = -d
	}
	deg := math.Trunc(d)
	frac := d - deg
	min := math.Trunc(frac * 60)
	sec := (frac*60 - min) * 60
	e.X = sign * (deg + min/100 + sec/10000)
}

func (e *Engine) ToRect() {
	r := e.X
	theta := e.Y
	e.LastX = r
	e.X = clamp(r * math.Cos(theta))
	e.Y = clamp(r * math.Sin(theta))
}

func (e *Engine) ToPolar() {
	x := e.X
	y := e.Y
	e.LastX = x
	r := math.Hypot(x, y)
	theta := math.Atan2(y, x)
	e.X = clamp(r)
	e.Y = clamp(theta)
}

func (e *Engine) Pi() {
	e.result(math.Pi)
}

func (e *Engine) Ln() {
	e.X = clamp(math.Log(e.X))
}

func (e *Engine) Log() {
	e.X = clamp(math.Log10(e.X))
}

func (e *Engine) Exp() {
	e.X = clamp(math.Exp(e.X))
}

func (e *Engine) Exp10() {
	e.X = clamp(math.Pow(10, e.X))
}

func (e *Engine) Recip() {
	if isZero(e.X) {
		e.X = math.Inf(1)
	} else {
		e.X = clamp(1 / e.X)
	}
}

func (e *Engine) Sqrt() {
	if e.X < 0 {
		e.X = math.NaN()
	} else {
		e.X = clamp(math.Sqrt(e.X))
	}
}

func (e *Engine) Sqr() {
	e.X = clamp(e.X * e.X)
}

func (e *Engine) Abs() {
	e.X = math.Abs(e.X)
}

func (e *Engine) Intg() {
	e.X = math.Trunc(e.X)
}

func (e *Engine) Frac() {
	e.X = e.X - math.Trunc(e.X)
}

func (e *Engine) Fact() {
	if e.X < 0 || e.X != math.Trunc(e.X) {
		e.X = math.NaN()
		return
	}
	n := int(e.X)
	if n > 69 {
		e.X = 9.999999999e99
		return
	}
	result := 1.0
	for i := 2; i <= n; i++ {
		result *= float64(i)
	}
	e.X = clamp(result)
}

func (e *Engine) Pct() {
	e.LastX = e.X
	e.X = clamp(e.Y * e.X / 100)
}

func (e *Engine) PctChg() {
	e.LastX = e.X
	if isZero(e.Y) {
		e.X = math.NaN()
	} else {
		e.X = clamp((e.X - e.Y) / e.Y * 100)
	}
	e.tuck()
	e.Flags.StackLift = true
}

func (e *Engine) PctTotal() {
	e.LastX = e.X
	e.X = clamp(e.X / e.Y * 100)
}

func (e *Engine) angleInRads() float64 {
	switch e.Flags.Angle {
	case Deg:
		return e.X * math.Pi / 180
	case Grad:
		return e.X * math.Pi / 200
	default:
		return e.X
	}
}

func (e *Engine) toCurAngle() {
	switch e.Flags.Angle {
	case Deg:
		e.X = e.X * 180 / math.Pi
	case Grad:
		e.X = e.X * 200 / math.Pi
	}
}
