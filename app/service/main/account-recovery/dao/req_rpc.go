package dao

import (
	"context"

	account "go-common/app/service/main/account/api"
	location "go-common/app/service/main/location/model"
	member "go-common/app/service/main/member/api"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Info3 get info by mid
func (d *Dao) Info3(c context.Context, mid int64) (info *account.Info, err error) {
	var (
		arg = &account.MidReq{
			Mid:    mid,
			RealIp: metadata.String(c, metadata.RemoteIP),
		}
		res *account.InfoReply
	)
	if res, err = d.accountClient.Info3(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
		return nil, err
	}
	return res.Info, nil
}

// Infos get the ips info.
func (d *Dao) Infos(c context.Context, ipList []string) (res map[string]*location.Info, err error) {
	if res, err = d.locRPC.Infos(c, ipList); err != nil {
		log.Error("s.locaRPC err(%v)", err)
	}
	return
}

// CheckRealnameStatus realname status
func (d *Dao) CheckRealnameStatus(c context.Context, mid int64) (status int8, err error) {
	var (
		relnameStatus *member.RealnameStatusReply
	)
	if relnameStatus, err = d.memberClient.RealnameStatus(c, &member.MemberMidReq{Mid: mid, RemoteIP: metadata.String(c, metadata.RemoteIP)}); err != nil {
		log.Error("s.memberSvr.RealnameStatus err(%v)", err)
		return
	}
	return relnameStatus.RealnameStatus, nil
}
