package service

import (
	"context"
	"time"

	"go-common/app/job/main/search/dao"
)

// incr increment data
func (s *Service) incr(c context.Context, app dao.App) {
	app.InitIndex(c)
	app.Sleep(c)
	app.Offset(c)
Loop:
	for {
		length, err := app.IncrMessages(c)
		if err != nil {
			s.base.D.PromError("IncrMessages", "IncrMessages error(%v)", err)
			app.Sleep(c)
			continue
		}
		for start := 0; start < length; start += _bulkSize {
			diff := length - start
			if diff > _bulkSize {
				diff = _bulkSize
			}
			if err := app.BulkIndex(c, start, start+diff, false); err != nil {
				s.base.D.PromError("BulkIndex", "BulkIndex error(%v)", err)
				time.Sleep(120 * time.Second) // 使databus readTimeout重新消费
				continue Loop
			}
		}
		if err := app.Commit(c); err != nil {
			s.base.D.PromError("UpdateOffsetID", "UpdateOffsetID error(%v)", err)
		}
		app.Sleep(c)
	}
}
