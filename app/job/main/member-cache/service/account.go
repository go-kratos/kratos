package service

import (
	"context"
	"fmt"

	"go-common/app/job/main/member-cache/model"
	"go-common/library/log"
)

func notifyKey(mid int64) string {
	return fmt.Sprintf("MemberCacheJob-AccountNotify-T%d", mid)
}

// NotifyPurgeCache is
func (s *Service) NotifyPurgeCache(ctx context.Context, mid int64, action string) {
	msg := &model.NotifyInfo{
		Mid:    mid,
		Action: action,
	}
	key := notifyKey(mid)
	if err := s.accountNotify.Send(ctx, key, msg); err != nil {
		log.Error("Failed to notify to purge cache with msg: %+v: %+v", msg, err)
	}
	log.Info("Succeed to notify to purge cache with msg: %+v", msg)
}
