package main

import (
	"math"

	"github.com/shanehull/dozen/engine"
)

type CalcService struct {
	e *engine.Engine
}

type KeyResult struct {
	Display engine.Display `json:"display"`
	StackX  float64        `json:"stackX"`
	StackY  float64        `json:"stackY"`
	StackZ  float64        `json:"stackZ"`
	StackT  float64        `json:"stackT"`
	LastX   float64        `json:"lastX"`
	Flags   []string       `json:"flags"`
}

type KeyInput struct {
	Op   string  `json:"op"`
	Arg  float64 `json:"arg"`
	ArgS string  `json:"argS"`
}

func NewCalcService() *CalcService {
	return &CalcService{e: engine.New()}
}

func (c *CalcService) PressKey(input KeyInput) KeyResult {
	c.e.Step(input.Op, input.Arg, input.ArgS)
	return c.state()
}

func (c *CalcService) GetState() KeyResult {
	return c.state()
}

func (c *CalcService) Save() string {
	b, _ := c.e.Snapshot().Marshal()
	return string(b)
}

func (c *CalcService) Load(j string) {
	s, err := engine.UnmarshalState([]byte(j))
	if err != nil {
		return
	}
	c.e.Restore(s)
}

func sanitize(f float64) float64 {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0
	}
	return f
}

func (c *CalcService) state() KeyResult {
	return KeyResult{
		Display: c.e.Format(),
		StackX:  sanitize(c.e.X),
		StackY:  sanitize(c.e.Y),
		StackZ:  sanitize(c.e.Z),
		StackT:  sanitize(c.e.T),
		LastX:   sanitize(c.e.LastX),
		Flags:   c.e.Format().Flags,
	}
}
