package service

import (
	"context"
	"time"

	"go-common/library/log"
)

const (
	_hbaseRetryCnt      = 3
	_hbaseRetryInterval = 100 * time.Millisecond
)

func (s *Service) replyToHbase(c context.Context, l int8, id int64, area, msg string) {
	if l > 0 && id != 0 {
		s.hbaseCh.Save(func() {
			for i := 0; i < _hbaseRetryCnt; i++ {
				if err1 := s.dao.SetContent(context.Background(), id, area, msg); err1 != nil {
					log.Error("s.dao.SetContent(%d, %s, %s) retry:%d error(%v)", id, area, msg, i, err1)
					time.Sleep(_hbaseRetryInterval)
				} else {
					break
				}
			}
		})
	}
}
