package batchwrite

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"go-common/app/service/main/dapper/model"
	"go-common/library/log"
)

var (
	_writeTimeout = time.Second
	// ErrClosed .
	ErrClosed = errors.New("batchwriter already closed")
)

// BatchWriter BatchWriter
type BatchWriter interface {
	WriteSpan(span *model.Span) error
	Close() error
	// internale queue length
	QueueLen() int
}

type rawBundle struct {
	key  string
	data map[string][]byte
}

// NewRawDataBatchWriter NewRawDataBatchWriter
func NewRawDataBatchWriter(writeFunc func(context.Context, string, map[string][]byte) error, bufSize, chanSize, workers int, interval time.Duration) BatchWriter {
	if workers <= 0 {
		workers = 1
	}
	if interval <= 0 {
		interval = 5 * time.Second
	}
	rbw := &rawDataBatchWrite{
		maxBufSize: bufSize,
		ch:         make(chan *rawBundle, chanSize),
		bufMap:     make(map[string]map[string][]byte),
		timeout:    10 * time.Second,
		writeFunc:  writeFunc,
	}
	rbw.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go rbw.worker()
	}
	rbw.flushTicker = time.NewTicker(interval)
	go rbw.daemonFlush()
	return rbw
}

type rawDataBatchWrite struct {
	mx          sync.Mutex
	closed      bool
	maxBufSize  int
	sizeCount   int
	bufMap      map[string]map[string][]byte
	ch          chan *rawBundle
	timeout     time.Duration
	writeFunc   func(context.Context, string, map[string][]byte) error
	wg          sync.WaitGroup
	flushTicker *time.Ticker
}

func (r *rawDataBatchWrite) WriteSpan(span *model.Span) error {
	data, err := span.Marshal()
	if err != nil {
		return err
	}
	traceID := span.TraceIDStr()
	spanID := span.SpanIDStr()
	kind := "_s"
	if !span.IsServer() {
		kind = "_c"
	}
	key := spanID + kind
	var bufMap map[string]map[string][]byte
	r.mx.Lock()
	if r.sizeCount > r.maxBufSize {
		bufMap = r.bufMap
		r.bufMap = make(map[string]map[string][]byte)
		r.sizeCount = 0
	}
	r.sizeCount += len(data)
	if _, ok := r.bufMap[traceID]; !ok {
		r.bufMap[traceID] = make(map[string][]byte)
	}
	r.bufMap[traceID][key] = data
	closed := r.closed
	r.mx.Unlock()
	if closed {
		return ErrClosed
	}
	if bufMap != nil {
		return r.flushBufMap(bufMap)
	}
	return nil
}

func (r *rawDataBatchWrite) QueueLen() int {
	return len(r.ch)
}

func (r *rawDataBatchWrite) daemonFlush() {
	for range r.flushTicker.C {
		if err := r.flush(); err != nil {
			log.Error("flush raw data error: %s", err)
		}
	}
}

func (r *rawDataBatchWrite) flush() error {
	var bufMap map[string]map[string][]byte
	r.mx.Lock()
	if r.sizeCount != 0 {
		bufMap = r.bufMap
		r.bufMap = make(map[string]map[string][]byte)
		r.sizeCount = 0
	}
	r.mx.Unlock()
	if bufMap != nil {
		return r.flushBufMap(bufMap)
	}
	return nil
}

func (r *rawDataBatchWrite) flushBufMap(bufMap map[string]map[string][]byte) error {
	timer := time.NewTimer(_writeTimeout)
	for traceID, data := range bufMap {
		select {
		case <-timer.C:
			return errors.New("write span timeout, raw data buffer channel is full")
		case r.ch <- &rawBundle{
			key:  traceID,
			data: data,
		}:
		}
	}
	return nil
}

func (r *rawDataBatchWrite) Close() error {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.closed = true
	r.flushTicker.Stop()
	bufMap := r.bufMap
	r.bufMap = make(map[string]map[string][]byte)
	r.sizeCount = 0
	r.flushBufMap(bufMap)
	close(r.ch)
	r.wg.Wait()
	return nil
}

func (r *rawDataBatchWrite) worker() {
	for bundle := range r.ch {
		if err := r.write(bundle); err != nil {
			log.Error("batch write raw data error: %s", err)
		}
	}
	r.wg.Done()
}

func (r *rawDataBatchWrite) write(bundle *rawBundle) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.writeFunc(ctx, bundle.key, bundle.data)
}
