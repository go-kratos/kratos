package service

import (
	"context"

	"go-common/library/log"
)

// Blacklist space blacklist
func (s *Service) Blacklist(c context.Context) {
	var (
		blTmp []int64
		err   error
	)
	if blTmp, err = s.dao.Blacklist(c); err != nil {
		log.Error("Service.Blacklist error(%v)", err)
		return
	}
	blacklist := make(map[int64]struct{}, len(blTmp))
	for _, mid := range blTmp {
		blacklist[mid] = struct{}{}
	}
	s.BlacklistValue = blacklist
}
