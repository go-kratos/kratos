package classify

import (
	"strings"
	"context"
	"go-common/app/service/ops/log-agent/event"
	"go-common/app/service/ops/log-agent/processor"
	"go-common/app/service/ops/log-agent/pkg/common"
)

type Classify struct {
	c *Config
}

func init() {
	err := processor.Register("classify", Process)
	if err != nil {
		panic(err)
	}
}

func Process(ctx context.Context, config interface{}, input <-chan *event.ProcessorEvent) (output chan *event.ProcessorEvent, err error) {
	classify := new(Classify)

	if c, ok := config.(*Config); !ok {
		panic("Error config for Classify Processor")
	} else {
		if err = c.ConfigValidate(); err != nil {
			return nil, err
		}
		classify.c = c
	}

	output = make(chan *event.ProcessorEvent)
	go func() {
		for {
			select {
			case e := <-input:
				// only do classify for ops-log
				if e.Destination == "lancer-ops-log" {
					if common.CriticalLog(e.Level) || e.Priority == "high" {
						e.LogId = classify.getLogIdByLevel("important")
					} else {
						e.LogId = classify.getLogIdByAppId(e.AppId)
					}
				}
				output <- e
			case <-ctx.Done():
				return
			}
		}
	}()
	return output, nil
}

// getLogLevel get logId level by appId
func (c *Classify) getLogIdByAppId(appId []byte) (logId string) {
	// get logId by setting
	if logLevel, ok := c.c.LogLevelMapConfig[string(appId)]; ok {
		return c.getLogIdByLevel(logLevel)
	}

	// appId format error, logId 1
	if len(strings.Split(string(appId), ".")) < 3 {
		return c.getLogIdByLevel("low") // low level
	}

	// set logLevel to 2 by default
	return c.getLogIdByLevel("normal") // normal level
}

// getLogIdByLevel return logid by level
func (c *Classify) getLogIdByLevel(level string) (logId string) {
	if logId, ok := c.c.LogIdMapConfig[level]; ok {
		return logId
	} else {
		// return 000161 by default
		return "000161"
	}
}
