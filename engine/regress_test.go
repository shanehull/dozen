package engine
import "testing"
import "math"
func TestRegressFV_NoDup(t *testing.T){
	e:=New(); e.FinN,e.FinI,e.FinPV,e.FinPMT=10,5,-1000,0
	e.X,e.Y,e.Z,e.T=1,2,3,4
	e.Step("FV",0,"FV")
	// result lifts old X(1) into Y; must NOT duplicate result into Y
	if math.Abs(e.X-1628.89)>0.5 { t.Fatalf("X=%v",e.X) }
	if math.Abs(e.Y-1)>1e-9 { t.Fatalf("Y should be lifted old X=1, got %v",e.Y) }
	if e.Z!=2||e.T!=3 { t.Fatalf("Z=%v T=%v",e.Z,e.T) }
}
func TestRegressNPV_ConsumesRate(t *testing.T){
	e:=New(); e.FinCF0=-100; e.FinCFj[0]=50; e.FinCFj[1]=60; e.FinCFj[2]=70; e.FinCfCnt=3
	e.Y=2; e.X=10; e.Flags.StackLift=true
	e.Step("NPV",0,"NPV")
	if math.Abs(e.X-47.63)>0.1 { t.Fatalf("X=%v",e.X) }
	if math.Abs(e.Y-2)>1e-9 { t.Fatalf("Y should stay 2 (rate consumed), got %v",e.Y) }
}
func TestRegressPiThenDigit(t *testing.T){
	e:=New(); e.Flags.StackLift=true
	e.Step("π",0,"π"); e.Step("5",5,"5")
	if e.X!=5 { t.Fatalf("X=%v",e.X) }
	if math.Abs(e.Y-math.Pi)>1e-9 { t.Fatalf("Y should be pi, got %v",e.Y) }
}
func TestRegressKeyedNumberThenRCL(t *testing.T){
	e:=New()
	for _,k:=range []string{"4","2"} { e.Step(k,float64(k[0]-'0'),k) }
	e.Step("STO",0,"STO"); e.Step("5",5,"5")
	for _,k:=range []string{"9","9"} { e.Step(k,float64(k[0]-'0'),k) }
	e.Step("RCL",0,"RCL"); e.Step("5",5,"5")
	if e.X!=42 { t.Fatalf("X=%v want 42",e.X) }
	if e.Y!=99 { t.Fatalf("Y=%v want 99 (keyed number preserved)",e.Y) }
}
func TestRegressEEXbare(t *testing.T){
	e:=New(); e.Step("EEX",0,"EEX"); e.Step("3",3,"3")
	if e.X!=1000 { t.Fatalf("EEX 3 => %v want 1000",e.X) }
}
func TestRegressMeanNoDup(t *testing.T){
	e:=New()
	e.X,e.Y=20,10; e.statAdd(); e.X,e.Y=30,20; e.statAdd()
	e.Z=7; e.Flags.StackLift=true
	e.statMeanX()
	if e.X!=15 { t.Fatalf("mean X=%v",e.X) }
	if e.Y==15 { t.Fatalf("Y duplicated mean (bug)") }
}
