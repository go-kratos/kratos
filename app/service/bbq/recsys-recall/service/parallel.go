package service

import (
	"context"
	"time"

	"go-common/library/log"

	"github.com/Dai0522/workpool"
)

func (s *Service) parallel(ctx *context.Context, tasks []*RecallTask) []*workpool.FutureTask {
	ftasks := make([]*workpool.FutureTask, len(tasks))
	for i := range tasks {
		ft := workpool.NewFutureTask(tasks[i])
		err := s.wp.Submit(ft)

		retry := 0
		for err != nil && retry < 3 {
			log.Errorv(*ctx, log.KV("workpool", err))
			err = s.wp.Submit(ft)
			retry++
		}
		ftasks[i] = ft
	}

	return ftasks
}

func (s *Service) wait(ctx context.Context, ftasks []*workpool.FutureTask) []*RecallResult {
	result := make([]*RecallResult, len(ftasks))

	for i, t := range ftasks {
		pt := t.T.(*RecallTask)
		result[i] = &RecallResult{
			Tag:      (*pt).info.Tag,
			Name:     (*pt).info.Name,
			Priority: (*pt).info.Priority,
		}

		raw, err := t.Wait(100 * time.Millisecond)
		if err != nil || raw == nil || len(*raw) <= 0 {
			log.Errorv(ctx, log.KV("future task wait", err))
			continue
		}
		res, err := parseResult(raw)
		if err != nil {
			log.Errorv(ctx, log.KV("parse recall result", err))
			continue
		}
		result[i].Result = *res
	}

	return result
}
