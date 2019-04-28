package archive

import (
	"context"
)

// AddNetSafeMd5 fn
func (s *Service) AddNetSafeMd5(c context.Context, nid int64, md5 string) (err error) {
	err = s.arc.AddNetSafeMd5(c, nid, md5)
	go s.NotifyNetSafe(c, nid)
	return
}

// NotifyNetSafe fn
func (s *Service) NotifyNetSafe(c context.Context, nid int64) (err error) {
	_ = s.arc.NotifyNetSafe(c, nid)
	return
}
