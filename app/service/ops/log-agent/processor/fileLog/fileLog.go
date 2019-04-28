package fileLog

import (
	"time"
	"context"
	"encoding/json"
	"os"

	"go-common/app/service/ops/log-agent/event"
	"go-common/app/service/ops/log-agent/processor"
)

const (
	_logIdLen           = 6
	_logLancerHeaderLen = 19
	_appIdKey           = `"app_id":`
	_levelKey           = `"level":`
	_logTime            = `"time":`
)

var (
	local, _    = time.LoadLocation("Local")
	hostname, _ = os.Hostname()
)

type FileLog struct {
	c *Config
}

func init() {
	err := processor.Register("fileLog", Process)
	if err != nil {
		panic(err)
	}
}

func Process(ctx context.Context, config interface{}, input <-chan *event.ProcessorEvent) (output chan *event.ProcessorEvent, err error) {
	fileLog := new(FileLog)

	if c, ok := config.(*Config); !ok {
		panic("Error config for jsonLog Processor")
	} else {
		if err = c.ConfigValidate(); err != nil {
			return nil, err
		}
		fileLog.c = c
	}

	output = make(chan *event.ProcessorEvent)
	go func() {
		for {
			select {
			case e := <-input:
				// only do jsonLog for ops-log
				if e.Destination != "lancer-ops-log" {
					output <- e
					continue
				}

				// format message
				message := make(map[string]interface{})
				if len(e.ParsedFields) != 0 {
					for k, v := range e.ParsedFields {
						message[k] = v
					}
				}

				message["log"] = e.String()

				e.Fields["hostname"] = hostname

				if len(e.Tags) != 0 {
					e.Fields["tag"] = e.Tags
				}

				if len(e.Fields) != 0 {
					message["fields"] = e.Fields
				}

				message["app_id"] = string(e.AppId)
				message["time"] = e.Time.UTC().Format(time.RFC3339Nano)

				if body, err := json.Marshal(message); err == nil {
					e.Write(body)
					output <- e
				}

			case <-ctx.Done():
				return
			}
		}
	}()
	return output, nil
}
