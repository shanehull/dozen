package engine

import "math"

func (e *Engine) sin() {
	e.X = clamp(math.Sin(e.angleInRads()))
}

func (e *Engine) cos() {
	e.X = clamp(math.Cos(e.angleInRads()))
}

func (e *Engine) tan() {
	e.X = clamp(math.Tan(e.angleInRads()))
}

func (e *Engine) asin() {
	e.X = clamp(math.Asin(e.X))
	e.toCurAngle()
}

func (e *Engine) acos() {
	e.X = clamp(math.Acos(e.X))
	e.toCurAngle()
}

func (e *Engine) atan() {
	e.X = clamp(math.Atan(e.X))
	e.toCurAngle()
}

func (e *Engine) toRad() {
	e.X = clamp(e.X * math.Pi / 180)
}

func (e *Engine) toDeg() {
	e.X = clamp(e.X * 180 / math.Pi)
}

func (e *Engine) toH() {
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

func (e *Engine) toHMS() {
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

func (e *Engine) toRect() {
	r := e.X
	theta := e.Y
	e.LastX = r
	e.X = clamp(r * math.Cos(theta))
	e.Y = clamp(r * math.Sin(theta))
}

func (e *Engine) toPolar() {
	x := e.X
	y := e.Y
	e.LastX = x
	r := math.Hypot(x, y)
	theta := math.Atan2(y, x)
	e.X = clamp(r)
	e.Y = clamp(theta)
}

func (e *Engine) pi() {
	e.result(math.Pi)
}

func (e *Engine) ln() {
	e.X = clamp(math.Log(e.X))
}

func (e *Engine) log() {
	e.X = clamp(math.Log10(e.X))
}

func (e *Engine) exp() {
	e.X = clamp(math.Exp(e.X))
}

func (e *Engine) exp10() {
	e.X = clamp(math.Pow(10, e.X))
}

func (e *Engine) recip() {
	if isZero(e.X) {
		e.X = math.Inf(1)
	} else {
		e.X = clamp(1 / e.X)
	}
}

func (e *Engine) sqrt() {
	if e.X < 0 {
		e.X = math.NaN()
	} else {
		e.X = clamp(math.Sqrt(e.X))
	}
}

func (e *Engine) sqr() {
	e.X = clamp(e.X * e.X)
}

func (e *Engine) abs() {
	e.X = math.Abs(e.X)
}

func (e *Engine) intg() {
	e.X = math.Trunc(e.X)
}

func (e *Engine) frac() {
	e.X = e.X - math.Trunc(e.X)
}

func (e *Engine) fact() {
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

func (e *Engine) pct() {
	e.LastX = e.X
	e.X = clamp(e.Y * e.X / 100)
}

func (e *Engine) pctChg() {
	e.LastX = e.X
	if isZero(e.Y) {
		e.X = math.NaN()
	} else {
		e.X = clamp((e.X - e.Y) / e.Y * 100)
	}
	e.tuck()
	e.Flags.StackLift = true
}

func (e *Engine) pctTotal() {
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
