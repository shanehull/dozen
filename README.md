<p align="center">
  <strong>dozen</strong> &middot; RPN financial calculator library &amp; desktop app
</p>

[![Go Reference](https://pkg.go.dev/badge/github.com/shanehull/dozen.svg)](https://pkg.go.dev/github.com/shanehull/dozen)
[![CI](https://github.com/shanehull/dozen/actions/workflows/ci.yml/badge.svg)](https://github.com/shanehull/dozen/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Two things in one repo:

- **`engine`** — a generic RPN calculator library: stack, registers, TVM, cash
  flows, statistics, depreciation, date math, and a `Step` dispatcher for
  keystroke-driven use. No GUI dependency.
- **App** — a cross-platform financial calculator built with
  [Wails v3](https://v3.wails.io) on top of the engine. macOS, Windows, Linux,
  iOS, Android.

---

# Library

```bash
go get github.com/shanehull/dozen/engine
```

Two layers: pure functions for direct math, and `Engine` for stack/keystroke
behaviour that the app is built on.

## Functional API

Rates are per-period fractions (`0.05` = 5%). Inflows positive, outflows
negative. `engine.End` (ordinary annuity) or `engine.Begin` (annuity due).

```go
import "github.com/shanehull/dozen/engine"

pmt := engine.PMT(0.06/12, 360, 300000, 0, engine.End)
fmt.Printf("%.2f\n", -pmt) // 1798.65

npv := engine.NPV(0.10, -100000, 30000, 40000, 50000, 60000)
fmt.Printf("%.2f\n", npv)  // 38877.13

irr, _ := engine.IRR(-100000, 40000, 50000, 30000)
fmt.Printf("%.2f%%\n", irr*100) // 24.89%
```

### TVM

```go
FV(rate, n, pv, pmt float64, when Timing) float64
PV(rate, n, pmt, fv float64, when Timing) float64
PMT(rate, n, pv, fv float64, when Timing) float64
NPer(rate, pmt, pv, fv float64, when Timing) (float64, error)
Rate(n, pmt, pv, fv float64, when Timing) (float64, error)
```

### Cash flows

```go
NPV(rate float64, cashflows ...float64) float64
IRR(cashflows ...float64) (float64, error)
```

### Depreciation

`factorPct` = declining-balance factor as percent (200 = double-declining).

```go
DepSL(cost, salvage, life, period float64) (dep, remaining float64)
DepSOYD(cost, salvage, life, period float64) (dep, remaining float64)
DepDB(cost, salvage, life, period, factorPct float64) (dep, remaining float64)
```

## Engine

`Engine` is an RPN calculator with a four-register stack, 20 memory registers,
TVM and cash-flow slots, statistics, and `f`/`g` prefix keys. Everything goes
through `Step`:

```go
func New() *Engine
func (e *Engine) Step(op string, arg float64, argS string)
func (e *Engine) Format() Display
func (e *Engine) Snapshot() EngineState
func (e *Engine) Restore(EngineState)
```

```go
c := engine.New()
c.Step("2", 2, "2")
c.Step("ENTER", 0, "")
c.Step("3", 3, "3")
c.Step("+", 0, "+")   // c.X == 5

c.X, c.Y, c.Z, c.T                // stack
c.FinN, c.FinI, c.FinPV, c.FinPMT, c.FinFV  // TVM registers
c.FinCF0, c.FinCFj, c.FinNj       // cash-flow registers
c.Flags.Begin, c.Flags.Dmy        // mode flags
c.Mem[0]..c.Mem[19]               // 20 general registers
c.Format()                        // Display{Mantissa, Sign, Flags}
```

Ops recognised by `Step`:

| Op | Args | Notes |
| --- | --- | --- |
| `0`–`9`, `.`, `CHS`, `EEX` | digit in arg | entry keys |
| `ENTER`, `CLx`, `LSTx` | — | stack |
| `+` `−` `×` `÷` `x↔y` `R↓` `R↑` | — | arithmetic / stack |
| `n` `i` `PV` `PMT` `FV` | — | store X into register, or solve |
| `NPV` `IRR` `AMORT` `INT` | — | cash flow / amortization |
| `PRICE` `YTM` `SL` `SOYD` `DB` | — | bonds / depreciation |
| `SIN` `COS` `TAN` `LN` `LOG` `yˣ` `√x` `1/x` `n!` `π` | — | scientific |
| `%` `%CHG` `%T` `INTG` `FRAC` `\|x\|` | — | utility |
| `Σ+` `x̄` `ŷ` `s` `ŷ,r` | — | statistics |
| `DYS` `DATE` | — | date |
| `STO` + digit, `RCL` + digit | index in arg | memory |
| `FIX` `SCI` `ENG` + digit | digits in arg | g‑prefixed display mode |
| `CLEAR FIN` `CLEAR REG` `CLEAR Σ` `CLEAR PRGM` | — | f‑prefixed clears |
| `P/R` `R/S` `GTO` `SST` `BST` | — | program mode |
| `CF0` `CFj` `Nj` | — | g‑prefixed cash flows |
| `BEG` `END` `D.MY` `DEG` `RAD` `GRAD` `→H.MS` `→H` `→R` `→P` | — | g‑prefixed modes |

## Examples

```bash
go run examples/mortgage/main.go    # $300K mortgage payment
go run examples/npv/main.go         # NPV + IRR
go run examples/fv/main.go          # future value
go run examples/stats/main.go       # statistics (engine)
go run examples/date/main.go        # date calculations (engine)
```

---

# App

A financial RPN calculator: stack, memory, TVM, cash flows, statistics,
programs. Built with [Wails v3](https://v3.wails.io).

```bash
wails3 dev
```

## Keyboard

| Key | Action |
| --- | --- |
| `0`–`9`, `.` | digits |
| `+` `-` `*` `/` | `+` `−` `×` `÷` |
| `Enter` | `ENTER` |
| `Esc` / `Backspace` | `CLx` |
| `%` | `%` |
| `f` / `g` | arm gold / blue shift |
| `n` `i` `p` `m` `v` | `n` `i` `PV` `PMT` `FV` |
| `s` | `STO` prefix |
| `r` | `RCL` prefix |
| `x` | `x↔y` |

---

# Development

```bash
go test ./engine/           # unit + regression tests
go test -run TestUI .       # UI integration tests
go vet ./...                # lint
wails3 dev                  # app with hot reload
```

# License

MIT
