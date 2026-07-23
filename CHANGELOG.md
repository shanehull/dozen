# Changelog

## [1.0.0](https://github.com/shanehull/dozen/compare/v0.3.0...v1.0.0) (2026-07-23)


### ⚠ BREAKING CHANGES

* The engine package no longer exports pure functions. All math (FV, PV, PMT, NPV, IRR, DepSL, DepSOYD, DepDB, Amort, BondPrice, BondYield, SolveTVM, Rate, NPer, Timing, End, Begin, NaNf64, ErrNoConvergence, ErrNoSolution, MaxDigits) is now package-private. SetN/SetI/SetPV/SetPMT/SetFV now take a value argument instead of reading from the X register.

### Features

* add HP 12c Platinum g-shift key labels and backspace ([#16](https://github.com/shanehull/dozen/issues/16)) ([7ba3ce4](https://github.com/shanehull/dozen/commit/7ba3ce4f88200797a9f374aeeb73da9975c86a13))
* rename ComputeNPV/IRR to NPV/IRR, add solve-tvm example ([#11](https://github.com/shanehull/dozen/issues/11)) ([ed4f114](https://github.com/shanehull/dozen/commit/ed4f1145da35c6dc341580a3f738219bb464fd78))


### Code Refactoring

* remove pure-function API, engine is RPN-only ([#15](https://github.com/shanehull/dozen/issues/15)) ([b21acba](https://github.com/shanehull/dozen/commit/b21acbac0d4b64fe6eb2f41fa41927d08e8cdf54))

## [0.3.0](https://github.com/shanehull/dozen/compare/v0.2.0...v0.3.0) (2026-07-16)


### Features

* add SolveTVM to library API, delegate engine to it ([#10](https://github.com/shanehull/dozen/issues/10)) ([a8b0441](https://github.com/shanehull/dozen/commit/a8b0441b82d9435e955420999760f24a3eaedc81))


### Bug Fixes

* tighten spacing between ENTER and LSTx label ([#8](https://github.com/shanehull/dozen/issues/8)) ([3cf13ca](https://github.com/shanehull/dozen/commit/3cf13cae2a31e40c2dc0dcf1ead28b3ca4e971b6))

## [0.2.0](https://github.com/shanehull/dozen/compare/v0.1.0...v0.2.0) (2026-07-16)


### Features

* generate app icon for macOS .app bundle ([#3](https://github.com/shanehull/dozen/issues/3)) ([7ad3f03](https://github.com/shanehull/dozen/commit/7ad3f0301504cf03abc1e43bdb1bd3afbbbd1abd))


### Bug Fixes

* resolve all golangci-lint warnings and add lint CI job ([#5](https://github.com/shanehull/dozen/issues/5)) ([a4e8cba](https://github.com/shanehull/dozen/commit/a4e8cba7eb04fea4c7e55a03254662f9a3aa27e7))

## 0.1.0 (2026-07-16)


### Features

* initial release ([4ad3d77](https://github.com/shanehull/dozen/commit/4ad3d779b5ae2fea6d9d27169da888997a7ccac7))
