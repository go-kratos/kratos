package member

import (
	"context"
	"fmt"
)

func notifyKey(mid int64) string {
	return fmt.Sprintf("AccountInterface-AccountNotify-T%d", mid)
}

// NotifyInfo notify info.
type NotifyInfo struct {
	Uname   string `json:"uname"`
	Mid     int64  `json:"mid"`
	Type    string `json:"type"`
	NewName string `json:"newName"`
	Action  string `json:"action"`
}

// NotityPurgeCache is
func (s *Service) NotityPurgeCache(ctx context.Context, mid int64, action string) error {
	msg := &NotifyInfo{
		Mid:    mid,
		Action: action,
	}
	key := notifyKey(mid)
	return s.accountNotify.Send(ctx, key, msg)
}
