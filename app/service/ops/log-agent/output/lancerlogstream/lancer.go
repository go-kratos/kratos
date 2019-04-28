package lancerlogstream

import (
	"context"
	"fmt"
	"bytes"
	"sync"
	"encoding/binary"
	"strconv"
	"time"

	"go-common/app/service/ops/log-agent/event"
	"go-common/app/service/ops/log-agent/output"
	"go-common/app/service/ops/log-agent/pkg/flowmonitor"
	"go-common/app/service/ops/log-agent/pkg/common"
	"go-common/app/service/ops/log-agent/output/cache/file"
	"go-common/library/log"
	"go-common/app/service/ops/log-agent/pkg/lancermonitor"
)

const (
	_logLenStart       = 2
	_logLenEnd         = 6
	_tokenHeaderFormat = "logId=%s&timestamp=%s&version=1.1"
	_protocolLen       = 6
	_appIdKey          = `"app_id":`
	_levelKey          = `"level":`
	_logTime           = `"time":`
)

var (
	logMagic    = []byte{0xAC, 0xBE}
	logMagicBuf = []byte{0xAC, 0xBE}
	_logType    = []byte{0, 1}
	_logLength  = []byte{0, 0, 0, 0}
	local, _    = time.LoadLocation("Local")
)

type logDoc struct {
	b     []byte
	logId string
}

func init() {
	err := output.Register("lancer", NewLancer)
	if err != nil {
		panic(err)
	}
}

type Lancer struct {
	c              *Config
	next           chan string
	i              chan *event.ProcessorEvent
	cache          *file.FileCache
	logAggrBuf     map[string]*bytes.Buffer
	logAggrBufLock sync.Mutex
	sendChan       chan *logDoc
	connPool       *connPool
	ctx            context.Context
	cancel         context.CancelFunc
}

func NewLancer(ctx context.Context, config interface{}) (output.Output, error) {
	var err error

	lancer := new(Lancer)
	if c, ok := config.(*Config); !ok {
		return nil, fmt.Errorf("Error config for Lancer output")
	} else {
		if err = c.ConfigValidate(); err != nil {
			return nil, err
		}
		lancer.c = c
	}
	if output.OutputRunning(lancer.c.Name) {
		return nil, fmt.Errorf("Output %s already running", lancer.c.Name)
	}

	lancer.i = make(chan *event.ProcessorEvent)
	lancer.next = make(chan string, 1)
	lancer.logAggrBuf = make(map[string]*bytes.Buffer)
	lancer.sendChan = make(chan *logDoc)
	cache, err := file.NewFileCache(lancer.c.CacheConfig)
	if err != nil {
		return nil, err
	}
	lancer.cache = cache
	lancer.c.PoolConfig.Name = lancer.c.Name
	lancer.connPool, err = initConnPool(lancer.c.PoolConfig)
	if err != nil {
		return nil, err
	}
	lancer.ctx, lancer.cancel = context.WithCancel(ctx)
	return lancer, nil
}

func (l *Lancer) InputChan() (chan *event.ProcessorEvent) {
	return l.i
}

func (l *Lancer) Run() (err error) {
	go l.writeToCache()
	go l.readFromCache()
	go l.flushLogAggrPeriodically()
	for i := 0; i < l.c.SendConcurrency; i++ {
		go l.sendToLancer()
	}
	output.RegisterOutput(l.c.Name, l)
	return nil
}

func (l *Lancer) Stop() {
	l.cancel()
}

// writeToCache write the log to cache
func (l *Lancer) writeToCache() {
	for e := range l.i {
		if e.Length < _logLancerHeaderLen {
			event.PutEvent(e)
			continue
		}
		l.cache.WriteToCache(e)
	}
}

func (l *Lancer) readFromCache() {
	for {
		e := l.cache.ReadFromCache()
		if e.Length < _logLancerHeaderLen {
			event.PutEvent(e)
			continue
		}
		// monitor should be called before event recycle
		l.parseOpslog(e)
		flowmonitor.Fm.AddEvent(e, "log-agent.output.lancer", "OK", "write to lancer")
		lancermonitor.IncreaseLogCount("agent.send.success.count", e.LogId)
		if l.c.Name == "lancer-ops-log" {
			l.logAggr(e)
		} else {
			l.sendLogDirectToLancer(e)
		}

	}
}

func (l *Lancer) parseOpslog(e *event.ProcessorEvent) {
	if l.c.Name == "lancer-ops-log" && e.Length > _logLancerHeaderLen {
		logBody := e.Body[(_logLancerHeaderLen):(e.Length)]
		e.AppId, _ = common.SeekValue([]byte(_appIdKey), logBody)

		if timeValue, err := common.SeekValue([]byte(_logTime), logBody); err == nil {
			if len(timeValue) >= 19 {
				// parse time
				var t time.Time
				if t, err = time.Parse(time.RFC3339Nano, string(timeValue)); err != nil {
					if t, err = time.ParseInLocation("2006-01-02T15:04:05", string(timeValue), local); err != nil {
						if t, err = time.ParseInLocation("2006-01-02T15:04:05", string(timeValue[0:19]), local); err != nil {
						}
					}
				}
				if !t.IsZero() {
					e.TimeRangeKey = strconv.FormatInt(t.Unix()/100*100, 10)
				}
			}
		}
	}
}


// sendLogDirectToLancer send log direct to lancer without aggr
func (l *Lancer) sendLogDirectToLancer(e *event.ProcessorEvent) {
	logDoc := new(logDoc)
	logDoc.b = make([]byte, e.Length)
	copy(logDoc.b, e.Bytes())
	logDoc.logId = e.LogId
	event.PutEvent(e)
	l.sendChan <- logDoc
}

// sendproc send the proc to lancer
func (l *Lancer) sendToLancer() {
	logSend := new(bytes.Buffer)
	tokenHeaderLen := []byte{0, 0}
	for {
		select {
		case logDoc := <-l.sendChan:
			var err error
			if len(logDoc.b) == 0 {
				continue
			}
			// header
			logSend.Reset()
			logSend.Write(logMagicBuf)
			logSend.Write(_logLength) // placeholder
			logSend.Write(_logType)
			// token header
			tokenheader := []byte(fmt.Sprintf(_tokenHeaderFormat, logDoc.logId, strconv.FormatInt(time.Now().Unix()/100*100, 10)))
			binary.BigEndian.PutUint16(tokenHeaderLen, uint16(len(tokenheader)))
			logSend.Write(tokenHeaderLen)
			logSend.Write(tokenheader)
			// log body
			logSend.Write(logDoc.b)

			// set log length
			bs := logSend.Bytes()
			binary.BigEndian.PutUint32(bs[_logLenStart:_logLenEnd], uint32(logSend.Len()-_protocolLen))

			// write
			connBuf, err := l.connPool.getBufConn()
			if err != nil {
				flowmonitor.Fm.Add("log-agent", "log-agent.output.lancer", "", "ERROR", "get conn failed")
				log.Error("get conn error: %v", err)
				continue
			}
			if _, err = connBuf.write(bs); err != nil {
				log.Error("wr.Write(log) error(%v)", err)
				connBuf.enabled = false
				l.connPool.putBufConn(connBuf)
				flowmonitor.Fm.Add("log-agent", "log-agent.output.lancer", "", "ERROR", "write to lancer failed")
				continue
			}
			l.connPool.putBufConn(connBuf)
			// TODO: flowmonitor for specific appId
		}
	}
}
