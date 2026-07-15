package engine

import "math"

func (e *Engine) dateDys() {
	d1 := e.Y
	d2 := e.X
	var dp1, dp2 dateParts
	if e.Flags.Dmy {
		dp1 = dmyParse(d1)
		dp2 = dmyParse(d2)
	} else {
		dp1 = mdyParse(d1)
		dp2 = mdyParse(d2)
	}
	days := daysBetween(dp1, dp2)
	e.X = clamp(days)
	e.Y = 0
	e.Z = 0
	e.T = 0
}

func (e *Engine) dateDate() {
	date := e.Y
	days := int(e.X)
	var d dateParts
	if e.Flags.Dmy {
		d = dmyParse(date)
	} else {
		d = mdyParse(date)
	}
	jd := julianDay(d.year, d.month, d.day)
	jd += days
	nd := fromJulian(jd)
	var result float64
	if e.Flags.Dmy {
		result = float64(nd.day) + float64(nd.month)/100 + float64(nd.year)/1000000
	} else {
		result = float64(nd.month) + float64(nd.day)/100 + float64(nd.year)/1000000
	}
	dow := float64(weekday(jd))
	e.X = clamp(result)
	e.Y = dow
	e.Z = 0
	e.T = 0
}

type dateParts struct {
	year, month, day int
}

func mdyParse(d float64) dateParts {
	mm := int(d)
	frac := d - float64(mm)
	fs := int(math.Round(frac * 1000000))
	dd := fs / 10000
	yyyy := fs % 10000
	return dateParts{yyyy, mm, dd}
}

func dmyParse(d float64) dateParts {
	dd := int(d)
	frac := d - float64(dd)
	fs := int(math.Round(frac * 1000000))
	mm := fs / 10000
	yyyy := fs % 10000
	return dateParts{yyyy, mm, dd}
}

func daysBetween(d1, d2 dateParts) float64 {
	jd1 := julianDay(d1.year, d1.month, d1.day)
	jd2 := julianDay(d2.year, d2.month, d2.day)
	return float64(jd2 - jd1)
}

func julianDay(y, m, d int) int {
	if m <= 2 {
		y--
		m += 12
	}
	a := y / 100
	b := 2 - a + a/4
	return int(365.25*float64(y+4716)) + int(30.6001*float64(m+1)) + d + b - 1524
}

func fromJulian(jd int) dateParts {
	l := jd + 68569
	n := 4 * l / 146097
	l = l - (146097*n+3)/4
	i := 4000 * (l + 1) / 1461001
	l = l - 1461*i/4 + 31
	j := 80 * l / 2447
	d := l - 2447*j/80
	l = j / 11
	m := j + 2 - 12*l
	y := 100*(n-49) + i + l
	return dateParts{y, m, d}
}

func weekday(jd int) int {
	return (jd + 1) % 7
}
