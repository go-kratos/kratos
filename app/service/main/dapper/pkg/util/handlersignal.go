package util

import (
	"os"
	"os/signal"
	"syscall"
)

// HandlerExit handler exit signal
func HandlerExit(exitFn func(s os.Signal) int) {
	sch := make(chan os.Signal, 1)
	signal.Notify(sch, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	s := <-sch
	os.Exit(exitFn(s))
}

// HandlerReload handler Reload signal
func HandlerReload(reload func(s os.Signal)) {
	go func() {
		sch := make(chan os.Signal, 1)
		signal.Notify(sch, syscall.SIGHUP)
		for s := range sch {
			reload(s)
		}
	}()
}
