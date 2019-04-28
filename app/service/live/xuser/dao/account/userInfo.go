package account

import (
	"context"
	"github.com/pkg/errors"
	v1 "go-common/app/service/main/account/api"
	accountM "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// GetUserCard3 ...
// 调用account grpc接口cards获取用户信息
func (d *Dao) GetUserCard3(c context.Context, UIDs []int64) (userResult map[int64]*accountM.Card, err error) {
	userResult = make(map[int64]*accountM.Card)
	lens := len(UIDs)
	if lens <= 0 {
		return
	}
	ret, err := d.accountRPC.Cards3(c, &accountM.ArgMids{Mids: UIDs})
	if err != nil {
		err = errors.WithMessage(ecode.AccountGRPCError, "GET SEA PATROL FAIL")
		log.Error("Call main.Account.Cards Error.Infos(%+v) error(%+v)", UIDs, err)
	}
	// 整理数据
	for _, item := range ret {
		if item != nil {
			userResult[item.Mid] = item
		}
	}
	return
}

// GetUserInfo ...
// 调用account grpc接口info获取用户信息
func (d *Dao) GetUserInfo(c context.Context, UID int64) (userResult *v1.Info, err error) {
	userResult = &v1.Info{}
	ret, err := d.accountRPC.Info3(c, &accountM.ArgMid{Mid: UID})
	if err != nil {
		err = errors.WithMessage(ecode.AccountGRPCError, "GET SEA PATROL FAIL")
		log.Error("Call main.Account.Cards Error.Infos(%+v) error(%+v)", UID, err)
		return
	}
	userResult = ret
	return
}
