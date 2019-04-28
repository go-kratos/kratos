package service

import (
	"context"

	"go-common/app/admin/main/videoup/model/manager"
	"go-common/library/log"
)

func (s *Service) isLeader(c context.Context, uid int64) bool {
	role, e := s.mng.GetUserRole(c, uid)
	if e != nil {
		log.Error("s.mng.GetUserRole(%d) error(%v)", uid, e)
		return false
	}
	if role == manager.TaskLeader {
		return true
	}
	return false
}
