package log

import (
	"context"
	"sync"
	"testing"
)

func TestValue(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	_ = logger.Log(LevelInfo, "msg", "helloworld")

	logger = DefaultLogger
	logger = With(logger)
	_ = logger.Log(LevelDebug, "msg", "helloworld")

	var v1 interface{}
	got := Value(context.Background(), v1)
	if got != v1 {
		t.Errorf("Value() = %v, want %v", got, v1)
	}
	var v2 Valuer = func(ctx context.Context) interface{} {
		return 3
	}
	got = Value(context.Background(), v2)
	res := got.(int)
	if res != 3 {
		t.Errorf("Value() = %v, want %v", res, 3)
	}
}

func TestCallerDepth(t *testing.T) {
	logger := With(DefaultLogger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	filter := NewFilter(logger, FilterLevel(LevelDebug))
	helper := NewHelper(logger)
	mLog := MultiLogger(logger, filter, helper)
	logs := []Logger{logger, filter, helper, mLog}
	for i := 0; i < 2; i++ {
		for _, lgr := range logs {
			logs = append(logs, With(lgr))
			logs = append(logs, NewFilter(lgr, FilterLevel(LevelDebug)))
			logs = append(logs, NewHelper(lgr))
		}
	}
	logs = append(logs, MultiLogger(logs...))
	filter = NewFilter(DefaultLogger, FilterLevel(LevelDebug))
	helper = NewHelper(DefaultLogger)
	logs = append(logs, filter, helper)
	for _, lgr := range logs {
		_ = lgr.Log(LevelDebug, "msg", "51")
		if h, ok := lgr.(*Helper); ok {
			h.Debug("53")
		}
	}
}

func TestCancel(t *testing.T) {
	logger := With(DefaultLogger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	filter := NewFilter(logger, FilterLevel(LevelDebug))
	helper := NewHelper(logger)
	mLog := MultiLogger(logger, filter, helper)
	logs := []Logger{logger, filter, helper, mLog}
	for i := 0; i < 2; i++ {
		for _, lgr := range logs {
			logs = append(logs, With(lgr))
			logs = append(logs, NewFilter(lgr, FilterLevel(LevelDebug)))
			logs = append(logs, NewHelper(lgr))
		}
	}
	logs = append(logs, MultiLogger(logs...))
	filter = NewFilter(DefaultLogger, FilterLevel(LevelDebug))
	helper = NewHelper(DefaultLogger)
	logs = append(logs, filter, helper)

	var wg sync.WaitGroup
	for _, lgr := range logs {
		lgr := lgr
		wg.Add(1)
		go func() {
			ctx, cancel := context.WithCancel(context.Background())
			lg := WithContext(ctx, lgr)
			_ = lg.Log(LevelDebug, "msg", "83")
			select {
			case <-ctx.Done():
				t.Error("Canceled")
			default:
			}
			cancel()
			select {
			case <-ctx.Done():
			default:
				t.Error("Not cancelled")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
