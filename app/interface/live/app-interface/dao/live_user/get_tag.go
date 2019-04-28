package live_user

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go-common/app/interface/live/app-interface/conf"
	cDao "go-common/app/interface/live/app-interface/dao"
	liveUserV1 "go-common/app/service/live/live_user/api/liverpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	rpcCtx "go-common/library/net/rpc/liverpc/context"
	"time"
)

func (d *Dao) GetUserTagList(ctx context.Context, req *liveUserV1.UserSettingGetTagReq) (rep *liveUserV1.UserSettingGetTagResp_Data, err error) {
	getTagTimeout := time.Duration(conf.GetTimeout("getMyTag", 50)) * time.Millisecond
	var tagListRep *liveUserV1.UserSettingGetTagResp
	tagListRep, err = cDao.LiveUserApi.V1UserSetting.GetTag(rpcCtx.WithTimeout(ctx, getTagTimeout), req)
	rep = &liveUserV1.UserSettingGetTagResp_Data{}
	if err != nil {
		log.Error("[GetUserTagList]live_user.v1.getTag rpc error:%+v", err)
		// 这个是推荐房间列表的基础信息，如果失败需要cancel，不然返回值会很奇怪
		err = errors.WithMessage(ecode.UserTagRPCError, fmt.Sprintf("live_user.v1.getTag rpc error:%+v", err))
		return
	}
	if tagListRep.Code != 0 {
		log.Error("[GetUserTagList]live_user.v1.getTag response error:%+v,code:%d,msg:%s", err, tagListRep.Code, tagListRep.Msg)
		// 这个是推荐房间列表的基础信息，如果失败需要cancel，不然返回值会很奇怪
		err = errors.WithMessage(ecode.UserTagReturnError, fmt.Sprintf("live_user.v1.getTag response error,code:%d,msg:%s", tagListRep.Code, tagListRep.Msg))
		return
	}

	if tagListRep.Data == nil {
		log.Error("[GetUserTagList]live_user.v1.getTag empty error")
		// 这个是推荐房间列表的基础信息，如果失败需要cancel，不然返回值会很奇怪
		err = errors.WithMessage(ecode.UserTagReturnError, "live_user.v1.getTag empty error")
		return
	}
	rep = tagListRep.Data

	return
}

