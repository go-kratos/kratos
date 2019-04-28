package service

import (
	"context"

	"go-common/library/log"
)

// AddNetSafeMd5 fn
func (s *Service) AddNetSafeMd5(c context.Context, nid int64, md5 string) (err error) {
	if _, err = s.arc.AddNetSafeMd5(c, nid, md5); err != nil {
		log.Error("s.arc.AddNetSafeMd5 nid:(%+v), md5:(%+v)", nid, md5)
		return
	}
	log.Info("s.arc.AddNetSafeMd5 nid:(%+v), md5:(%+v)", nid, md5)
	return
}
