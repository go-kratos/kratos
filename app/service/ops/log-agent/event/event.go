package event

import (
	"sync"
	"time"
	"go-common/library/log"
)

var pool sync.Pool

// event between input and processor
type ProcessorEvent struct {
	Source       string
	Destination  string
	LogId        string
	AppId        []byte
	Level        []byte
	Time         time.Time
	Body         []byte
	Priority     string
	Length       int
	TimeRangeKey string
	Fields       map[string]interface{}
	ParsedFields map[string]string
	Tags         []string
}

//  GetEvent get event from pool.
func GetEvent() (e *ProcessorEvent) {
	var (
		ok  bool
		tmp = pool.Get()
	)
	if e, ok = tmp.(*ProcessorEvent); !ok {
		e = &ProcessorEvent{Body: make([]byte, 1024*64), Tags: make([]string, 0, 1)} // max 64K, should be longer than max log lentth
	}
	e.LogId = ""
	e.Length = 0
	e.AppId = nil
	e.Level = nil
	e.Time = time.Time{}
	e.TimeRangeKey = ""
	e.Source = ""
	e.Priority = ""
	e.Destination = ""
	e.Tags = e.Tags[:0]
	e.Fields = make(map[string]interface{})
	e.ParsedFields = make(map[string]string)
	return e
}

// PutEvent put event back to pool
func PutEvent(e *ProcessorEvent) {
	pool.Put(e)
}

func (e *ProcessorEvent) Bytes() []byte {
	return e.Body[:e.Length]
}

func (e *ProcessorEvent) String() string {
	return string(e.Body[:e.Length])
}

func (e *ProcessorEvent) Write(b []byte) {
	if len(b) > cap(e.Body) {
		log.Error("bytes write beyond e.Body capacity")
	}

	e.Length = copy(e.Body, b)
}
