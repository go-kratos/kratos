// +build !windows

package syscall

import (
	"testing"
)

func TestSignal(t *testing.T) {
	if int(SIGSTOP) != 0x13 {
		t.FailNow()
	}
	if int(SIGXFSZ) != 0x19 {
		t.FailNow()
	}
}
