package service

import (
	"context"
	"go-common/library/log"
)

// 在播列表

// refreshLiveStreamList 刷新在播列表缓存
func (s *Service) refreshLiveStreamList() {
	log.Warn("refreshLiveStreamList")
	s.dao.StoreLiveStreamList()
}

// checkLiveStreamList 确认流是否在
func (s *Service) CheckLiveStreamList(c context.Context, rids []int64) map[int64]bool {
	return s.dao.LoadLiveStreamList(c, rids)
}
