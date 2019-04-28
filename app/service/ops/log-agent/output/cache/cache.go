package cache

import (
	"go-common/app/service/ops/log-agent/event"
)

type Cahce interface {
	WriteToCache(e *event.ProcessorEvent)
	ReadFromCache() (e *event.ProcessorEvent)
}
