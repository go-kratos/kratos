package rankdb

import (
	"context"
	"github.com/pkg/errors"
	"go-common/app/interface/live/app-interface/dao"
	"go-common/app/service/live/rankdb/api/liverpc/v1"
	accountM "go-common/app/service/main/account/model"
	actmdl "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"time"
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

// GetLastHourTop3 ...
// 获取小时榜前三
func (d *Dao) GetLastHourTop3(c context.Context) (uids []int64, err error) {
	uids = make([]int64, 0)
	lastHour := time.Now().Add(-time.Hour).Format("06010215")
	req := &v1.Rank2018GetHourRankReq{
		AreaV2ParentId: 0,
		AreaV2Id:       0,
		Hour:           lastHour,
		// rankDB begin from 0
		Top: 2,
	}
	resp, err := dao.RankdbApi.V1Rank2018.GetHourRank(c, req)
	if err != nil || resp.Data == nil {
		log.Error("[app-interface][rankDbItem] liveRpc call rankDb failed")
		return
	}
	if 0 != resp.Code || 0 == len(resp.Data) {
		log.Error("[app-interface][rankDbItem] liveRpc call rankDb return error, code:%d, msg:%s", resp.Code, resp.Data)
		return
	}
	uids = resp.Data
	return
}
