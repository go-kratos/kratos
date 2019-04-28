package main

import (
	"os"
	"os/signal"
	"syscall"

	"go-common/app/service/ep/saga-agent/conf"
	"go-common/app/service/ep/saga-agent/service/agent"
	"go-common/library/log"
)

func listenSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1,
		syscall.SIGUSR2, syscall.SIGTSTP)
	select {
	case <-sigs:
		log.Info("get sig=%v\n, RunnerUnRegisterAll", sigs)
		agent.RunnerUnRegisterAll()
	default:
		log.Info("get sig=%v\n", sigs)
	}
}

func main() {
	log.Info("agent start......")
	err := conf.Init()
	if err != nil {
		panic(err)
	}
	go listenSignal()
	go agent.UpdateRegister()
	agent.ExecRegister()
	agent.RunnerStart()
	agent.RunnerUnRegisterAll()
}
