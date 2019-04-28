package svg

import (
	"fmt"
	"os"
	"testing"
)

func Test_dot(t *testing.T) {

	dot := NewDot()
	dotcode := dot.StartDot().AddTokenBinds(tk1, tk2).
		AddFlow(flow1, flow2, flow3, flow4, flow5, flow6, flow7).
		AddTransitions(tran1, tran2, tran3).
		AddDirections(dir1, dir2, dir3, dir4, dir5, dir6, dir7, dir8, dir9).
		End()
	nv := NewNetView()
	nv.SetDot(dot)
	t.Log(nv.Execute(os.Stdout, nv.Data))

	fmt.Println(dotcode)
	t.Fail()
}
