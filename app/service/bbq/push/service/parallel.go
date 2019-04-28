package service

import (
	"context"
	"time"

	"go-common/app/service/bbq/push/api/grpc/v1"
	"go-common/library/log"

	"github.com/Dai0522/workpool"
	"github.com/json-iterator/go"
)

func (s *Service) parallel(ctx *context.Context, tasks []workpool.Task) (ftasks []*workpool.FutureTask, err error) {
	ftasks = make([]*workpool.FutureTask, len(tasks))
	for i := range tasks {
		ft := workpool.NewFutureTask(tasks[i])
		err = s.wp.Submit(ft)

		retry := 0
		for err != nil && retry < 3 {
			log.Errorv(*ctx, log.KV("workpool", err))
			err = s.wp.Submit(ft)
			retry++
		}
		if err != nil {
			return
		}
		ftasks[i] = ft
	}

	return
}

func (s *Service) wait(ctx context.Context, ftasks []*workpool.FutureTask) []*v1.PushResult {
	result := make([]*v1.PushResult, 0)

	for _, t := range ftasks {
		raw, err := t.Wait(800 * time.Millisecond)
		if err != nil || raw == nil || len(*raw) <= 0 {
			log.Errorv(ctx, log.KV("future task wait", err))
			continue
		}
		res := &v1.PushResult{}
		err = jsoniter.Unmarshal(*raw, res)
		if err != nil {
			log.Errorv(ctx, log.KV("parse push result", err))
			continue
		}
		result = append(result, res)
	}

	return result
}
