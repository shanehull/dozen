package engine

type Registers struct {
	FinN   float64
	FinI   float64
	FinPV  float64
	FinPMT float64
	FinFV  float64
	FinCF0 float64
	FinCFj [10]float64
	FinNj  [10]int
	FinCfCnt int
	AmortN   float64
	AmortInt float64
	AmortPrin float64
}
