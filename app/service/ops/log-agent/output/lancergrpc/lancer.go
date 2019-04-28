package lancergrpc

import (
	"context"
	"fmt"
	"bytes"
	"sync"
	"strconv"
	"time"
	"math"

	"go-common/app/service/ops/log-agent/event"
	"go-common/app/service/ops/log-agent/output"
	"go-common/app/service/ops/log-agent/pkg/flowmonitor"
	"go-common/app/service/ops/log-agent/pkg/common"
	"go-common/app/service/ops/log-agent/output/cache/file"
	"go-common/library/log"
	"go-common/app/service/ops/log-agent/pkg/lancermonitor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"go-common/app/service/ops/log-agent/output/lancergrpc/lancergateway"
)

const (
	_appIdKey = `"app_id":`
	_levelKey = `"level":`
	_logTime  = `"time":`
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
	err := output.Register("lancergrpc", NewLancer)
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
	lancerClient   lancergateway.Gateway2ServerClient
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
	lancer.lancerClient, err = lancergateway.NewClient(lancer.c.LancerGateway)
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
	go l.readFromProcessor()
	go l.consumeCache()
	go l.flushLogAggrPeriodically()
	for i := 0; i < l.c.SendConcurrency; i++ {
		go l.sendToLancer()
	}
	if l.c.Name != "" {
		output.RegisterOutput(l.c.Name, l)
	}
	return nil
}

func (l *Lancer) Stop() {
	l.cancel()
}

func (l *Lancer) readFromProcessor() {
	for e := range l.i {
		// only cache for sock input
		if e.Source == "sock" {
			l.cache.WriteToCache(e)
			continue
		}
		// without cache
		l.preWriteToLancer(e)
	}
}

func (l *Lancer) preWriteToLancer(e *event.ProcessorEvent) {
	flowmonitor.Fm.AddEvent(e, "log-agent.output.lancer", "OK", "write to lancer")
	lancermonitor.IncreaseLogCount("agent.send.success.count", e.LogId)
	if l.c.Name == "lancer-ops-log" {
		l.logAggr(e)
	} else {
		l.sendLogDirectToLancer(e)
	}
}

// consumeCache consume logs from cache
func (l *Lancer) consumeCache() {
	for {
		e := l.cache.ReadFromCache()
		if e.Length < _logLancerHeaderLen {
			event.PutEvent(e)
			continue
		}
		// monitor should be called before event recycle
		l.parseOpslog(e)
		l.preWriteToLancer(e)
	}
}

func (l *Lancer) parseOpslog(e *event.ProcessorEvent) {
	if l.c.Name == "lancer-ops-log" {
		e.AppId, _ = common.SeekValue([]byte(_appIdKey), e.Bytes())

		if timeValue, err := common.SeekValue([]byte(_logTime), e.Bytes()); err == nil {
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

func (l *Lancer) nextRetry(retry int) (time.Duration) {
	// avoid d too large
	if retry > 10 {
		return time.Duration(l.c.MaxRetryDuration)
	}

	d := time.Duration(math.Pow(2, float64(retry))) * time.Duration(l.c.InitialRetryDuration)

	if d > time.Duration(l.c.MaxRetryDuration) {
		return time.Duration(l.c.MaxRetryDuration)
	}

	return d
}

func (l *Lancer) bulkSendToLancerWithRetry(in *lancergateway.EventList) {
	retry := 0
	for {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(l.c.SendBatchTimeout))
		t1 := time.Now()
		resp, err := l.lancerClient.SendList(ctx, in)
		if err == nil {
			if resp.Code == lancergateway.StatusCode_SUCCESS {
				log.Info("get 200 from lancer gateway: size %d, count %d, cost %s", in.Size(), len(in.Events), time.Since(t1).String())
				return
			}
			flowmonitor.Fm.Add("log-agent", "log-agent.output.lancer", "", "ERROR", fmt.Sprintf("write to lancer None 200: %s", resp.Code))
			log.Warn("get None 200 from lancer gateway, retry: %s", resp.Code)
		}

		if err != nil {
			switch grpc.Code(err) {
			case codes.Canceled, codes.DeadlineExceeded, codes.Unavailable, codes.ResourceExhausted:
				flowmonitor.Fm.Add("log-agent", "log-agent.output.lancer", "", "ERROR", fmt.Sprintf("write to lancer failed, retry: %s", err))
				log.Warn("get error from lancer gateway, retry: %s", err)
			default:
				flowmonitor.Fm.Add("log-agent", "log-agent.output.lancer", "", "ERROR", fmt.Sprintf("write to lancer failed, no retry: %s", err))
				log.Warn("get error from lancer gateway, no retry: %s", err)
				return
			}
		}

		time.Sleep(l.nextRetry(retry))
		retry ++
	}
}

// sendproc send the proc to lancer
func (l *Lancer) sendToLancer() {
	eventList := new(lancergateway.EventList)
	eventListLock := sync.Mutex{}
	lastSend := time.Now()
	ticker := time.Tick(time.Second * 1)
	size := 0
	for {
		select {
		case <-ticker:
			if lastSend.Add(time.Duration(l.c.SendFlushInterval)).Before(time.Now()) && len(eventList.Events) > 0 {
				eventListLock.Lock()
				l.bulkSendToLancerWithRetry(eventList)
				eventList.Reset()
				size = 0
				eventListLock.Unlock()
				lastSend = time.Now()
			}
		case logDoc := <-l.sendChan:
			event := new(lancergateway.SimpleEvent)
			event.LogId = logDoc.logId
			event.Header = map[string]string{"timestamp": strconv.FormatInt(time.Now().Unix()/100*100, 10)}
			event.Data = logDoc.b
			size += len(event.Data)
			eventListLock.Lock()
			eventList.Events = append(eventList.Events, event)
			if size > l.c.SendBatchSize || len(eventList.Events) > l.c.SendBatchNum {
				l.bulkSendToLancerWithRetry(eventList)
				eventList.Reset()
				size = 0
				lastSend = time.Now()
			}
			eventListLock.Unlock()
		}
	}
}
