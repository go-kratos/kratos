package service

import (
	"context"
	"fmt"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/log"
)

// RecordChangeLog 记录操作日志
func (s *Service) RecordChangeLog(c context.Context, streamLog *model.StreamChangeLog) {
	err := s.dao.InsertChangeLog(c, streamLog)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("insert change log faild = %v", err)))
	}
}

// GetChangeLogByRoomID 得到切换cdn 记录
func (s *Service) GetChangeLogByRoomID(c context.Context, rid int64, limit int64) (infos []*model.StreamChangeLog, err error) {
	return s.dao.GetChangeLogByRoomID(c, rid, limit)
}
