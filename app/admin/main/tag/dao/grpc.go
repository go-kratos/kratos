package dao

import (
	"context"

	accwarden "go-common/app/service/main/account/api"
	"go-common/library/log"
)

// UserInfo user info.
func (d *Dao) UserInfo(c context.Context, mid int64) (res *accwarden.InfoReply, err error) {
	if res, err = d.accClient.Info3(c, &accwarden.MidReq{Mid: mid}); err != nil {
		log.Error("d.UserInfo(%d) error(%v)", mid, err)
	}
	return
}

// UserInfos user infos.
func (d *Dao) UserInfos(c context.Context, mids []int64) (res *accwarden.InfosReply, err error) {
	if res, err = d.accClient.Infos3(c, &accwarden.MidsReq{Mids: mids}); err != nil {
		log.Error("d.UserInfos(%v) error(%v)", mids, err)
	}
	return
}
