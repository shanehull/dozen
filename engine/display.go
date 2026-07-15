package engine

import (
	"fmt"
	"math"
)

type Display struct {
	Mantissa string
	Exponent string
	Sign     string
	Flags    []string
}

func (e *Engine) Format() Display {
	d := Display{}

	if e.entry.hasDigits {
		d.Mantissa = trimLeadingZeros(e.entry.buf)
		s := e.entry.buf
		if len(s) > 0 && s[0] == '-' {
			d.Sign = "-"
		}
		d.Flags = e.annunciators()
		return d
	}

	x := e.X

	if math.IsInf(x, 0) {
		d.Mantissa = "9.999999 99"
		d.Sign = " "
		return d
	}
	if math.IsNaN(x) {
		d.Mantissa = " Error    "
		d.Sign = " "
		return d
	}

	switch {
	case x < -9.999999999e99:
		d.Mantissa = "-9.999999 99"
		return d
	case x > 9.999999999e99:
		d.Mantissa = " 9.999999 99"
		return d
	}

	if x < 0 {
		d.Sign = "-"
		x = -x
	}

	switch e.Flags.DispMode {
	case Fix:
		fmtStr := fmt.Sprintf("%%.%df", e.Flags.DispDigits)
		d.Mantissa = fmt.Sprintf(fmtStr, x)
	case Sci:
		d.Mantissa = formatSci(x, e.Flags.DispDigits)
	case Eng:
		d.Mantissa = formatEng(x, e.Flags.DispDigits)
	}

	d.Flags = e.annunciators()
	return d
}

func formatSci(x float64, digits int) string {
	if x == 0 {
		s := fmt.Sprintf("%%.%df", digits)
		return fmt.Sprintf(s, 0.0) + "   00"
	}
	exp := 0
	for x >= 10 {
		x /= 10
		exp++
	}
	for x < 1 {
		x *= 10
		exp--
	}
	fmtStr := fmt.Sprintf("%%.%df", digits)
	mant := fmt.Sprintf(fmtStr, x)
	return mant + formatExp(exp)
}

func formatEng(x float64, digits int) string {
	if x == 0 {
		s := fmt.Sprintf("%%.%df", digits)
		return fmt.Sprintf(s, 0.0) + "   00"
	}
	exp := 0
	for x >= 1000 {
		x /= 1000
		exp += 3
	}
	for x < 1 {
		x *= 1000
		exp -= 3
	}
	fmtStr := fmt.Sprintf("%%.%df", digits)
	mant := fmt.Sprintf(fmtStr, x)
	return mant + formatExp(exp)
}

func formatExp(e int) string {
	if e >= 0 {
		return fmt.Sprintf("  %02d", e)
	}
	return fmt.Sprintf(" -%02d", -e)
}

func trimLeadingZeros(s string) string {
	if s == "" || s == "0" {
		return s
	}
	start := 0
	if s[0] == '-' {
		start = 1
	}
	i := start
	for i < len(s)-1 && s[i] == '0' && s[i+1] != '.' {
		i++
	}
	if i > start {
		return s[:start] + s[i:]
	}
	return s
}

func (e *Engine) annunciators() []string {
	var a []string
	if e.Flags.Prefix == "f" {
		a = append(a, "f")
	} else if e.Flags.Prefix == "g" {
		a = append(a, "g")
	}
	if e.Flags.Begin {
		a = append(a, "BEGIN")
	}
	if e.Flags.Dmy {
		a = append(a, "D.MY")
	}
	if e.Flags.ProgramMode {
		a = append(a, "PRGM")
	}
	if e.Flags.Angle == Rad {
		a = append(a, "RAD")
	} else if e.Flags.Angle == Grad {
		a = append(a, "GRAD")
	}
	return a
}
