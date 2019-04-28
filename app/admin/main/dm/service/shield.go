package service

import (
	"context"

	"go-common/library/log"
)

const (
	_bnjShieldFileName = "bnj_shield.csv"
	_bnjShieldBucket   = "dm"
)

// DmShield .
func (s *Service) DmShield(c context.Context, bs []byte) (err error) {
	if err = s.dao.Upload(c, _bnjShieldBucket, _bnjShieldFileName, "text/csv", bs); err != nil {
		log.Error("DmShield(err:%v)", err)
		return
	}
	return
}
