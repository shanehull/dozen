<p align="center">
  <strong>dozen</strong> &middot; RPN financial calculator library &amp; desktop app
</p>

[![Go Reference](https://pkg.go.dev/badge/github.com/shanehull/dozen.svg)](https://pkg.go.dev/github.com/shanehull/dozen)
[![CI](https://github.com/shanehull/dozen/actions/workflows/ci.yml/badge.svg)](https://github.com/shanehull/dozen/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

![Dozen calculator](docs/screenshot.png)

Two things in one repo:

- **`engine`**: a generic RPN calculator library: stack, registers, TVM, cash
  flows, statistics, depreciation, date math. No GUI dependency.
- **App**: a cross-platform financial calculator built with
  [Wails v3](https://v3.wails.io) on top of the engine. macOS, Windows, Linux,
  iOS, Android.

---

# Library

```bash
go get github.com/shanehull/dozen/engine
```

`Engine` is an RPN calculator with a four-level stack, 20 memory registers,
TVM and cash-flow slots, statistics accumulation, and program memory.

```go
c := engine.New()

c.X = 2; c.Enter(); c.X = 3
c.Add()                    // c.X == 5

c.SetN(360)                // FinN = 360
c.SetI(5)                  // FinI = 5
c.SetPV(-1000)             // FinPV = -1000
c.SolvePMT()               // c.X == monthly payment

c.X = 4; c.Sqrt()          // c.X == 2
c.X = 180; c.Sin()         // sine of 180 degrees

c.Snapshot()               // EngineState for persistence
c.Restore(state)
```

Set methods are variadic. Omit the argument to store the current X value:

```go
c.X = 8; c.SetI()          // equivalent to c.SetI(8)
```

### Methods

| Category     | Methods                                                                                        |
| ------------ | ---------------------------------------------------------------------------------------------- |
| Stack        | `Enter`, `Clx`, `Chs`, `XY`, `RollDown`, `RollUp`, `LastXRecall`                               |
| Arithmetic   | `Add`, `Sub`, `Mul`, `Div`, `YPowX`                                                            |
| TVM store    | `SetN(n?)`, `SetI(i?)`, `SetPV(pv?)`, `SetPMT(pmt?)`, `SetFV(fv?)`                             |
| TVM solve    | `SolveN`, `SolveI`, `SolvePV`, `SolvePMT`, `SolveFV`                                           |
| Cash flow    | `NPV`, `IRR`                                                                                   |
| Amortization | `Amortize`                                                                                     |
| Bonds        | `BondPrice`, `BondYield`                                                                       |
| Depreciation | `DepreciationSL`, `DepreciationSOYD`, `DepreciationDB`                                         |
| Statistics   | `StatAdd`, `MeanX`, `MeanY`, `SDev`, `WeightedMean`, `LinEst`, `ClearStats`                    |
| Scientific   | `Sin`, `Cos`, `Tan`, `Asin`, `Acos`, `Atan`, `Ln`, `Log`, `Sqrt`, `Sqr`, `Recip`, `Pi`, `Fact` |
| Trig helpers | `ToRad`, `ToDeg`, `ToRect`, `ToPolar`, `ToHMS`, `ToH`                                          |
| Percent      | `Pct`, `PctChg`, `PctTotal`                                                                    |
| Utility      | `Abs`, `Intg`, `Frac`, `Exp`, `Exp10`                                                          |
| Date         | `DaysBetween`, `DateAdd`                                                                       |
| Memory       | `Store(n)`, `Recall(n)`                                                                        |
| Program      | `SST`, `BST`, `Goto(line)`                                                                     |
| State        | `Snapshot`, `Restore`                                                                          |
| Clear        | `ClearFin`, `ClearReg`, `ClearStats`, `ClearPgm`                                               |

`SetN/SetI/SetPV/SetPMT/SetFV` are variadic: omit the argument to store the
current X register value (`e.SetN()` is shorthand for `e.SetN(e.X)`).

### Fields

```go
c.X, c.Y, c.Z, c.T                            // stack registers
c.LastX                                       // last X before operation
c.Mem[0]..c.Mem[19]                           // 20 general registers
c.FinN, c.FinI, c.FinPV, c.FinPMT, c.FinFV    // TVM registers
c.FinCF0, c.FinCFj, c.FinNj, c.FinCfCnt       // cash-flow registers
c.AmortN, c.AmortInt, c.AmortPrin             // amortization results
c.Flags.Begin, c.Flags.Dmy, c.Flags.Angle     // calculator mode
c.Flags.StackLift                             // stack lift flag
c.Program, c.PgmLen, c.PgmPC                  // program storage
```

## Examples

```bash
go run examples/mortgage/main.go    # $300K mortgage payment
go run examples/solve-tvm/main.go   # unified TVM solver
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

| Key                 | Action                  |
| ------------------- | ----------------------- |
| `0`–`9`, `.`        | digits                  |
| `+` `-` `*` `/`     | `+` `−` `×` `÷`         |
| `Enter`             | `ENTER`                 |
| `Esc` / `Backspace` | `CLx`                   |
| `%`                 | `%`                     |
| `f` / `g`           | arm gold / blue shift   |
| `n` `i` `p` `m` `v` | `n` `i` `PV` `PMT` `FV` |
| `s`                 | `STO` prefix            |
| `r`                 | `RCL` prefix            |
| `x`                 | `x↔y`                   |

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
