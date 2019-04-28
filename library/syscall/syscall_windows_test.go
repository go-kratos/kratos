package syscall

import (
	"log"
	"testing"
)

func TestSignal(t *testing.T) {
	if int(SIGSTOP) != 0x19 {
		t.FailNow()
	}
	if int(SIGXFSZ) != 0x25 {
		t.FailNow()
	}
}
