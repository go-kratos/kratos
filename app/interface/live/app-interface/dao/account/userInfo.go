package account

import (
    "context"
    "github.com/pkg/errors"
	"go-common/app/service/main/account/api"
	accountM "go-common/app/service/main/account/model"
    actmdl "go-common/app/service/main/account/model"
    "go-common/library/ecode"
    "go-common/library/log"
)

// GetUserInfoData ...
// 调用account grpc接口cards获取用户信息
func (d *Dao) GetUserInfoData(c context.Context, UIDs []int64) (userResult map[int64]*accountM.Card, err error) {
	userResult = make(map[int64]*accountM.Card)
	lens := len(UIDs)
	if lens <= 0 {
		return
	}
	ret, err := d.accountRPC.Cards3(c, &actmdl.ArgMids{Mids: UIDs})
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

func (d *Dao) GetUserInfos(c context.Context, uids []int64) (userResult map[int64]*api.Info, err error) {
	userResult = make(map[int64]*api.Info)
	if len(uids) <= 0 {
		return
	}
	ret, err := d.accountRPC.Infos3(c, &actmdl.ArgMids{Mids: uids})
	if err != nil {
		err = errors.WithMessage(ecode.AccountGRPCError, "GET USER INFO FAIL")
		log.Error("Call main.Account.Info3 Error.Infos(%+v) error(%+v)", uids, err)
	}
	userResult = ret
	return
}
