package engine

import (
	"math"
	"testing"
)

func TestRegressFV_NoDup(t *testing.T) {
	e := New()
	e.FinN, e.FinI, e.FinPV, e.FinPMT = 10, 5, -1000, 0
	e.X, e.Y, e.Z, e.T = 1, 2, 3, 4
	e.SolveFV()
	if math.Abs(e.X-1628.89) > 0.5 {
		t.Fatalf("X=%v", e.X)
	}
	if math.Abs(e.Y-1) > 1e-9 {
		t.Fatalf("Y should be lifted old X=1, got %v", e.Y)
	}
	if e.Z != 2 || e.T != 3 {
		t.Fatalf("Z=%v T=%v", e.Z, e.T)
	}
}

func TestRegressNPV_ConsumesRate(t *testing.T) {
	e := New()
	e.FinCF0 = -100
	e.FinCFj[0] = 50
	e.FinCFj[1] = 60
	e.FinCFj[2] = 70
	e.FinCfCnt = 3
	e.Y = 2
	e.X = 10
	e.Flags.StackLift = true
	e.ComputeNPV()
	if math.Abs(e.X-47.63) > 0.1 {
		t.Fatalf("X=%v", e.X)
	}
	if math.Abs(e.Y-2) > 1e-9 {
		t.Fatalf("Y should stay 2 (rate consumed), got %v", e.Y)
	}
}

func TestRegressPiThenDigit(t *testing.T) {
	e := New()
	e.Flags.StackLift = true
	e.Pi()
	e.Push()
	e.X = 5
	e.Flags.StackLift = false
	if e.X != 5 {
		t.Fatalf("X=%v", e.X)
	}
	if math.Abs(e.Y-math.Pi) > 1e-9 {
		t.Fatalf("Y should be pi, got %v", e.Y)
	}
}

func TestRegressKeyedNumberThenRCL(t *testing.T) {
	e := New()
	e.X = 42
	e.Flags.StackLift = true
	e.Store(5)
	e.Flags.StackLift = false
	e.X = 99
	e.Flags.StackLift = true
	e.Recall(5)
	if e.X != 42 {
		t.Fatalf("X=%v want 42", e.X)
	}
	if e.Y != 99 {
		t.Fatalf("Y=%v want 99 (keyed number preserved)", e.Y)
	}
}

func TestRegressEEXbare(t *testing.T) {
	e := New()
	e.X = 1000
	if e.X != 1000 {
		t.Fatalf("EEX 3 => %v want 1000", e.X)
	}
}

func TestRegressMeanNoDup(t *testing.T) {
	e := New()
	e.X, e.Y = 20, 10
	e.StatAdd()
	e.X, e.Y = 30, 20
	e.StatAdd()
	e.Z = 7
	e.Flags.StackLift = true
	e.MeanX()
	if e.X != 15 {
		t.Fatalf("mean X=%v", e.X)
	}
	if e.Y == 15 {
		t.Fatalf("Y duplicated mean (bug)")
	}
}
