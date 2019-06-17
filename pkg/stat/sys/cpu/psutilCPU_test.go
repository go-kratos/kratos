package cpu

import (
	"fmt"
	"testing"
	"time"
)

func Test_PsutilCPU(t *testing.T) {
	cpu, err := newPsutilCPU(time.Millisecond * 500)
	if err != nil {
		t.Fatalf("newPsutilCPU failed!err:=%v", err)
	}
	for i := 0; i < 6; i++ {
		time.Sleep(time.Millisecond * 500)
		u, err := cpu.Usage()
		if u == 0 {
			t.Fatalf("get cpu from psutil failed!cpu usage is zero!err:=%v", err)
		}
		fmt.Println(u)
	}
}
