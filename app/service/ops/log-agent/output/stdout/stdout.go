package stdout

import (
	"context"
	"fmt"

	"go-common/app/service/ops/log-agent/output"
	"go-common/app/service/ops/log-agent/event"
)

type Stdout struct {
	c      *Config
	ctx    context.Context
	cancel context.CancelFunc
	i      chan *event.ProcessorEvent
}

func init() {
	err := output.Register("stdout", NewStdout)
	if err != nil {
		panic(err)
	}
}

func NewStdout(ctx context.Context, config interface{}) (output.Output, error) {
	var err error

	stdout := new(Stdout)
	if c, ok := config.(*Config); !ok {
		return nil, fmt.Errorf("Error config for Lancer output")
	} else {
		if err = c.ConfigValidate(); err != nil {
			return nil, err
		}
		stdout.c = c
	}

	stdout.i = make(chan *event.ProcessorEvent)
	stdout.ctx, stdout.cancel = context.WithCancel(ctx)
	return stdout, nil
}

func (s *Stdout) Run() (err error) {
	go func() {
		for {
			select {
			case e := <-s.i:
				fmt.Println(string(e.Body))
			case <-s.ctx.Done():
				return
			}
		}
	}()
	if s.c.Name != "" {
		output.RegisterOutput(s.c.Name, s)
	}
	return nil
}

func (s *Stdout) Stop() {
	s.cancel()
}

func (s *Stdout) InputChan() (chan *event.ProcessorEvent) {
	return s.i
}