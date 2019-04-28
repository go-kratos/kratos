package xuser

import (
	"context"
	"github.com/pkg/errors"
	xuserM "go-common/app/service/live/xuser/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
)

// GetUserExpData ...
// 调用account grpc接口cards获取用户信息
func (d *Dao) GetUserExpData(c context.Context, UIDs []int64) (userResult map[int64]*xuserM.LevelInfo, err error) {
	userResult = make(map[int64]*xuserM.LevelInfo)
	lens := len(UIDs)
	if lens <= 0 {
		return
	}
	ret, err := d.xuserGRPC.GetUserExp(c, &xuserM.GetUserExpReq{Uids: UIDs})
	if err != nil {
		err = errors.WithMessage(ecode.AccountGRPCError, "GET SEA PATROL FAIL")
		log.Error("Call main.Account.Cards Error.Infos(%+v) error(%+v)", UIDs, err)
	}
	// 整理数据
	for _, item := range ret.Data {
		if item != nil {
			userResult[item.Uid] = item
		}
	}
	return
}
