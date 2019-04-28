package service

import "context"

// ClearRoomCacheByRID 清除一个房间缓存
func (s *Service) ClearRoomCacheByRID(c context.Context, rid int64) error {
	return s.dao.DeleteStreamByRIDFromCache(c, rid)
}
