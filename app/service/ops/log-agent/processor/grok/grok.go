package grok

import (
	"context"

	"go-common/app/service/ops/log-agent/event"
	"go-common/app/service/ops/log-agent/processor"
	"go-common/app/service/ops/log-agent/pkg/flowmonitor"

	"github.com/vjeantet/grok"
)

type Grok struct {
	c *Config
	g *grok.Grok
}

func init() {
	err := processor.Register("grok", Process)
	if err != nil {
		panic(err)
	}
}

func Process(ctx context.Context, config interface{}, input <-chan *event.ProcessorEvent) (output chan *event.ProcessorEvent, err error) {
	g := new(Grok)

	if c, ok := config.(*Config); !ok {
		panic("Error config for Grok Processor")
	} else {
		if err = c.ConfigValidate(); err != nil {
			return nil, err
		}
		g.c = c
	}

	if g.g, err = grok.New(); err != nil {
		return nil, err
	}

	output = make(chan *event.ProcessorEvent)
	go func() {
		for {
			select {
			case e := <-input:
				values, err := g.g.Parse(g.c.Pattern, e.String())
				if err != nil || len(values) == 0 {
					flowmonitor.Fm.AddEvent(e, "log-agent.processor.grok", "WARN", "grok error")
					e.Tags = append(e.Tags, "grok_error")
					output <- e
					continue
				}
				e.ParsedFields = values
				output <- e
			case <-ctx.Done():
				return
			}
		}
	}()

	return output, nil
}
