package sample

import (
	"math/rand"
	"context"

	"go-common/app/service/ops/log-agent/event"
	"go-common/app/service/ops/log-agent/processor"
	"go-common/app/service/ops/log-agent/pkg/common"
	"go-common/app/service/ops/log-agent/pkg/flowmonitor"
)

type Sample struct {
	c *Config
}

func init() {
	err := processor.Register("sample", Process)
	if err != nil {
		panic(err)
	}
}

func Process(ctx context.Context, config interface{}, input <-chan *event.ProcessorEvent) (output chan *event.ProcessorEvent, err error) {
	sample := new(Sample)

	if c, ok := config.(*Config); !ok {
		panic("Error config for sample Processor")
	} else {
		if err = c.ConfigValidate(); err != nil {
			return nil, err
		}
		sample.c = c
	}

	output = make(chan *event.ProcessorEvent)
	go func() {
		for {
			select {
			case e := <-input:
				// only do sample for ops-log
				if e.Destination != "lancer-ops-log" {
					output <- e
					continue
				}
				if !sample.sample(e) {
					output <- e
				} else {
					flowmonitor.Fm.AddEvent(e, "log-agent.processor.sample", "WARN", "sampled")
					event.PutEvent(e)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return output, nil
}

//sample log, if return ture, the log should be discard
func (s *Sample) sample(e *event.ProcessorEvent) bool {
	if common.CriticalLog(e.Level) {
		return false // keep log if level isn't INFO or DEBUG
	}

	if e.Priority == "high" {
		return false
	}

	if val, ok := s.c.SampleConfig[string(e.AppId)]; ok {

		if rand.Intn(100) < 100-int(val) {
			return true // discard
		}
	}
	return false
}
