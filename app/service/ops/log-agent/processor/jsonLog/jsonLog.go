package jsonLog

import (
	"time"
	"strconv"
	"context"

	"go-common/app/service/ops/log-agent/event"
	"go-common/app/service/ops/log-agent/processor"
	"go-common/app/service/ops/log-agent/pkg/common"
	"go-common/app/service/ops/log-agent/pkg/flowmonitor"
)

const (
	_appIdKey = `"app_id":`
	_levelKey = `"level":`
	_logTime  = `"time":`
)

var (
	local, _ = time.LoadLocation("Local")
)

type JsonLog struct {
	c *Config
}

func init() {
	err := processor.Register("jsonLog", Process)
	if err != nil {
		panic(err)
	}
}

func Process(ctx context.Context, config interface{}, input <-chan *event.ProcessorEvent) (output chan *event.ProcessorEvent, err error) {
	jsonLog := new(JsonLog)

	if c, ok := config.(*Config); !ok {
		panic("Error config for jsonLog Processor")
	} else {
		if err = c.ConfigValidate(); err != nil {
			return nil, err
		}
		jsonLog.c = c
	}

	output = make(chan *event.ProcessorEvent)
	var (
		t time.Time
	)
	go func() {
		for {
			select {
			case e := <-input:
				// only do jsonLog for ops-log
				if e.Destination != "lancer-ops-log" {
					output <- e
					continue
				}
				if e.Length == 0 {
					event.PutEvent(e)
					continue
				}

				// seek app_id
				if appId, err := common.SeekValue([]byte(_appIdKey), e.Bytes()); err == nil {
					e.AppId = appId
				}

				// priority
				if priority, err := common.GetPriority(e.Bytes()); err == nil {
					e.Priority = string(priority)
				}

				// seek time
				if timeValue, err := common.SeekValue([]byte(_logTime), e.Bytes()); err == nil {
					if len(timeValue) >= 19 {
						// parse time
						if t, err = time.Parse(time.RFC3339Nano, string(timeValue)); err != nil {
							if t, err = time.ParseInLocation("2006-01-02T15:04:05", string(timeValue), local); err != nil {
								if t, err = time.ParseInLocation("2006-01-02T15:04:05", string(timeValue[0:19]), local); err != nil {
								}
							}
						}
						e.Time = t
					}
				}

				// TimeRangeKey for flow monitor
				if !e.Time.IsZero() {
					e.TimeRangeKey = strconv.FormatInt(e.Time.Unix()/100*100, 10)
				}
				// seek level
				if level, err := common.SeekValue([]byte(_levelKey), e.Bytes()); err == nil {
					e.Level = level
				}
				flowmonitor.Fm.AddEvent(e, "log-agent.processor.jsonLog", "OK", "received")
				output <- e
			case <-ctx.Done():
				return
			}
		}
	}()
	return output, nil
}
