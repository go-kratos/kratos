package syscall

import (
	"fmt"
	gosyscall "syscall"
)

// Signal for windows.
// These Signals won't be registered by invoking signal.Notify.
// Use these signal to avoid build error on Windows.
type Signal int

// Signal nothing here
func (s Signal) Signal() {}

func (s Signal) String() string {
	return fmt.Sprintf("windows signal %d", int(s))
}

// Signals for windows.
const (
	// More invented values for signals
	SIGHUP  = gosyscall.Signal(0x1)
	SIGINT  = gosyscall.Signal(0x2)
	SIGQUIT = gosyscall.Signal(0x3)
	SIGILL  = gosyscall.Signal(0x4)
	SIGTRAP = gosyscall.Signal(0x5)
	SIGABRT = gosyscall.Signal(0x6)
	SIGBUS  = gosyscall.Signal(0x7)
	SIGFPE  = gosyscall.Signal(0x8)
	SIGKILL = gosyscall.Signal(0x9)
	SIGSEGV = gosyscall.Signal(0xb)
	SIGPIPE = gosyscall.Signal(0xd)
	SIGALRM = gosyscall.Signal(0xe)
	SIGTERM = gosyscall.Signal(0xf)

	//SIG fake linux signal for windwos.
	SIGCHLD   = Signal(0x10)
	SIGCLD    = Signal(0x11)
	SIGCONT   = Signal(0x12)
	SIGIO     = Signal(0x13)
	SIGIOT    = Signal(0x14)
	SIGPOLL   = Signal(0x15)
	SIGPROF   = Signal(0x16)
	SIGPWR    = Signal(0x17)
	SIGSTKFLT = Signal(0x18)
	SIGSTOP   = Signal(0x19)
	SIGSYS    = Signal(0x1a)
	SIGTSTP   = Signal(0x1b)
	SIGTTIN   = Signal(0x1c)
	SIGTTOU   = Signal(0x1d)
	SIGUNUSED = Signal(0x1e)
	SIGURG    = Signal(0x1f)
	SIGUSR1   = Signal(0x20)
	SIGUSR2   = Signal(0x21)
	SIGVTALRM = Signal(0x22)
	SIGWINCH  = Signal(0x23)
	SIGXCPU   = Signal(0x24)
	SIGXFSZ   = Signal(0x25)
)
