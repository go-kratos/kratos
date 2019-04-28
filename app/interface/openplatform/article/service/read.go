package service

import (
	"context"

	"go-common/library/ecode"
)

// ReadPing 处理用户阅读心跳
func (s *Service) ReadPing(c context.Context, buvid string, aid int64, mid int64, ip string, cur int64, source string) (err error) {
	var last int64
	if last, err = s.dao.GetsetReadPing(c, buvid, aid, cur); err != nil {
		err = ecode.RequestErr
		return
	}
	if last != 0 {
		return
	}
	if s.dao.AddReadPingSet(c, buvid, aid, mid, ip, cur, source); err != nil {
		err = ecode.RequestErr
		return
	}
	return
}
