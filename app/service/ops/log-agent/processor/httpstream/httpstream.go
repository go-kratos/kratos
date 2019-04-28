package httpstream

import (
	"context"

	"go-common/app/service/ops/log-agent/event"
	"go-common/app/service/ops/log-agent/processor"
	phttpstream "go-common/app/service/ops/log-agent/pkg/httpstream"
)

type HttpStream struct {
	c *Config
}

func init() {
	err := processor.Register("httpStream", Process)
	if err != nil {
		panic(err)
	}
}

func Process(ctx context.Context, config interface{}, input <-chan *event.ProcessorEvent) (output chan *event.ProcessorEvent, err error) {
	h := new(HttpStream)

	if c, ok := config.(*Config); !ok {
		panic("Error config for lengthCheck Processor")
	} else {
		if err = c.ConfigValidate(); err != nil {
			return nil, err
		}
		h.c = c
	}

	output = make(chan *event.ProcessorEvent)
	go func() {
		for {
			select {
			case e := <-input:
				select {
				case phttpstream.LogSourceChan <- e:
				default:
				}
				output <- e
			case <-ctx.Done():
				return
			}
		}
	}()

	return output, nil
}