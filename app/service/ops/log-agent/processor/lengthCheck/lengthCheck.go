package lengthCheck

import (
	"context"

	"go-common/app/service/ops/log-agent/event"
	"go-common/app/service/ops/log-agent/processor"
	"go-common/app/service/ops/log-agent/pkg/flowmonitor"
)

type LengthCheck struct {
	c *Config
}

func init() {
	err := processor.Register("lengthCheck", Process)
	if err != nil {
		panic(err)
	}
}

func Process(ctx context.Context, config interface{}, input <-chan *event.ProcessorEvent) (output chan *event.ProcessorEvent, err error) {
	lcheck := new(LengthCheck)

	if c, ok := config.(*Config); !ok {
		panic("Error config for lengthCheck Processor")
	} else {
		if err = c.ConfigValidate(); err != nil {
			return nil, err
		}
		lcheck.c = c
	}

	output = make(chan *event.ProcessorEvent)
	go func() {
		for {
			select {
			case e := <-input:
				// log length check
				if e.Length > lcheck.c.MaxLength {
					flowmonitor.Fm.AddEvent(e, "log-agent.processor.lengthCheck", "ERROR", "too long")
					event.PutEvent(e)
					continue
				}
				output <- e
			case <-ctx.Done():
				return
			}
		}
	}()
	return output, nil
}
