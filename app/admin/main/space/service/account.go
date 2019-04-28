package service

import (
	"context"

	"go-common/app/service/main/relation/model"
	"go-common/library/log"
)

// Relation .
func (s *Service) Relation(c context.Context, mid int64) (follower int64, err error) {
	var stat *model.Stat
	if stat, err = s.relation.Stat(c, &model.ArgMid{Mid: mid}); err != nil {
		log.Error("Relation s.relation.Stat(mid:%d) error(%v)", mid, err)
		return
	}
	follower = stat.Follower
	return
}
