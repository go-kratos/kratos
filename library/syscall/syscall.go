// +build !windows

package syscall

import (
	gosyscall "syscall"
)

// Signals
const (
	SIGABRT   = gosyscall.Signal(0x6)
	SIGALRM   = gosyscall.Signal(0xe)
	SIGBUS    = gosyscall.Signal(0x7)
	SIGCHLD   = gosyscall.Signal(0x11)
	SIGCLD    = gosyscall.Signal(0x11)
	SIGCONT   = gosyscall.Signal(0x12)
	SIGFPE    = gosyscall.Signal(0x8)
	SIGHUP    = gosyscall.Signal(0x1)
	SIGILL    = gosyscall.Signal(0x4)
	SIGINT    = gosyscall.Signal(0x2)
	SIGIO     = gosyscall.Signal(0x1d)
	SIGIOT    = gosyscall.Signal(0x6)
	SIGKILL   = gosyscall.Signal(0x9)
	SIGPIPE   = gosyscall.Signal(0xd)
	SIGPOLL   = gosyscall.Signal(0x1d)
	SIGPROF   = gosyscall.Signal(0x1b)
	SIGPWR    = gosyscall.Signal(0x1e)
	SIGQUIT   = gosyscall.Signal(0x3)
	SIGSEGV   = gosyscall.Signal(0xb)
	SIGSTKFLT = gosyscall.Signal(0x10)
	SIGSTOP   = gosyscall.Signal(0x13)
	SIGSYS    = gosyscall.Signal(0x1f)
	SIGTERM   = gosyscall.Signal(0xf)
	SIGTRAP   = gosyscall.Signal(0x5)
	SIGTSTP   = gosyscall.Signal(0x14)
	SIGTTIN   = gosyscall.Signal(0x15)
	SIGTTOU   = gosyscall.Signal(0x16)
	SIGUNUSED = gosyscall.Signal(0x1f)
	SIGURG    = gosyscall.Signal(0x17)
	SIGUSR1   = gosyscall.Signal(0xa)
	SIGUSR2   = gosyscall.Signal(0xc)
	SIGVTALRM = gosyscall.Signal(0x1a)
	SIGWINCH  = gosyscall.Signal(0x1c)
	SIGXCPU   = gosyscall.Signal(0x18)
	SIGXFSZ   = gosyscall.Signal(0x19)
)
