package engine

import "encoding/json"

type EngineState struct {
	X    float64      `json:"x"`
	Y    float64      `json:"y"`
	Z    float64      `json:"z"`
	T    float64      `json:"t"`
	Last float64      `json:"last"`
	Mem  [20]float64  `json:"mem"`

	FinN     float64     `json:"finN"`
	FinI     float64     `json:"finI"`
	FinPV    float64     `json:"finPV"`
	FinPMT   float64     `json:"finPMT"`
	FinFV    float64     `json:"finFV"`
	FinCF0   float64     `json:"finCF0"`
	FinCFj   [10]float64 `json:"finCFj"`
	FinNj    [10]int     `json:"finNj"`
	FinCfCnt int         `json:"finCfCnt"`
	AmortN   float64     `json:"amortN"`
	AmortInt float64     `json:"amortInt"`
	AmortPrin float64    `json:"amortPrin"`

	Begin bool      `json:"begin"`
	Dmy   bool      `json:"dmy"`
	Angle AngleMode `json:"angle"`
}

func (e *Engine) Snapshot() EngineState {
	return EngineState{
		X: e.X, Y: e.Y, Z: e.Z, T: e.T,
		Last: e.LastX,
		Mem:  e.Mem,
		FinN: e.FinN, FinI: e.FinI, FinPV: e.FinPV, FinPMT: e.FinPMT, FinFV: e.FinFV,
		FinCF0: e.FinCF0, FinCFj: e.FinCFj, FinNj: e.FinNj, FinCfCnt: e.FinCfCnt,
		AmortN: e.AmortN, AmortInt: e.AmortInt, AmortPrin: e.AmortPrin,
		Begin: e.Flags.Begin, Dmy: e.Flags.Dmy, Angle: e.Flags.Angle,
	}
}

func (e *Engine) Restore(s EngineState) {
	e.X = s.X; e.Y = s.Y; e.Z = s.Z; e.T = s.T
	e.LastX = s.Last
	e.Mem = s.Mem
	e.FinN = s.FinN; e.FinI = s.FinI; e.FinPV = s.FinPV; e.FinPMT = s.FinPMT; e.FinFV = s.FinFV
	e.FinCF0 = s.FinCF0; e.FinCFj = s.FinCFj; e.FinNj = s.FinNj; e.FinCfCnt = s.FinCfCnt
	e.AmortN = s.AmortN; e.AmortInt = s.AmortInt; e.AmortPrin = s.AmortPrin
	e.Flags.Begin = s.Begin; e.Flags.Dmy = s.Dmy; e.Flags.Angle = s.Angle
}

func (s EngineState) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func UnmarshalState(data []byte) (EngineState, error) {
	var s EngineState
	err := json.Unmarshal(data, &s)
	return s, err
}
