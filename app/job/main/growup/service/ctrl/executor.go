package ctrl

import (
	"context"
	"reflect"
	"runtime"
	"sync"

	"go-common/library/log"
)

func NewUnboundedExecutor() *UnboundedExecutor {
	ctx := context.Background()
	return &UnboundedExecutor{
		ctx: ctx,
	}
}

type UnboundedExecutor struct {
	wg sync.WaitGroup
	// for future extension
	ctx context.Context
}

type Executor interface {
	Submit(bizFunc ...func(ctx context.Context))
}

func (executor *UnboundedExecutor) Submit(bizFunc ...func(c context.Context)) {
	for _, biz := range bizFunc {
		pc := reflect.ValueOf(biz).Pointer()
		funcName := runtime.FuncForPC(pc).Name()
		executor.wg.Add(1)
		go func(funcName string, biz func(ctx context.Context)) {
			defer executor.wg.Done()
			log.Info("Exec Task %s", funcName)
			biz(executor.ctx)
		}(funcName, biz)
	}
}
