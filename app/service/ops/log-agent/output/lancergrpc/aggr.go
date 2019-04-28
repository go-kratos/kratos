package lancergrpc

import (
	"bytes"
	"time"
	"go-common/app/service/ops/log-agent/event"
)

const (
	_logSeparator       = byte('\u0001')
	_logLancerHeaderLen = 19
)

// logAggr aggregates multi logs to one log
func (l *Lancer) logAggr(e *event.ProcessorEvent) (err error) {
	logAddrbuf := l.getlogAggrBuf(e.LogId)
	l.logAggrBufLock.Lock()
	logAddrbuf.Write(e.Bytes())
	logAddrbuf.WriteByte(_logSeparator)
	l.logAggrBufLock.Unlock()
	if logAddrbuf.Len() > l.c.AggrSize {
		return l.flushLogAggr(e.LogId)
	}
	event.PutEvent(e)
	return nil
}

// getlogAggrBuf get logAggrBuf by logId
func (l *Lancer) getlogAggrBuf(logId string) (*bytes.Buffer) {
	if _, ok := l.logAggrBuf[logId]; !ok {
		l.logAggrBuf[logId] = new(bytes.Buffer)
	}
	return l.logAggrBuf[logId]
}

// flushLogAggr write aggregated logs to conn
func (l *Lancer) flushLogAggr(logId string) (err error) {
	l.logAggrBufLock.Lock()
	defer l.logAggrBufLock.Unlock()
	buf := l.getlogAggrBuf(logId)
	if buf.Len() > 0 {
		logDoc := new(logDoc)
		logDoc.b = make([]byte, buf.Len())
		copy(logDoc.b, buf.Bytes())
		logDoc.logId = logId
		l.sendChan <- logDoc
	}
	buf.Reset()
	return nil
}

// flushLogAggrPeriodically run flushLogAggr Periodically
func (l *Lancer) flushLogAggrPeriodically() {
	tick := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-tick.C:
			for logid, _ := range l.logAggrBuf {
				l.flushLogAggr(logid)
			}
		}
	}
}
