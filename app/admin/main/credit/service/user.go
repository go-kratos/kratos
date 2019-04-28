package service

import (
	"context"

	account "go-common/app/service/main/account/api"
	"go-common/library/log"
)

// UserInfo get account info.
func (s *Service) UserInfo(c context.Context, mid int64) (res *account.InfoReply, err error) {
	if res, err = s.accDao.RPCInfo(c, mid); err != nil {
		log.Error("s.accDao.RPCInfo error (%v)", err)
	}
	return
}
