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

func (d *Dao) GetDaHangHai(ctx context.Context, req *liveUserV1.NoteGetReq) (rep *liveUserV1.NoteGetResp_Data, err error) {
	// TODO 添加DaHangHai超时配置
	getTagTimeout := time.Duration(conf.GetTimeout("DaHangHai", 50)) * time.Millisecond
	dahanghai, err := cDao.LiveUserApi.V1Note.Get(rpcCtx.WithTimeout(ctx, getTagTimeout), req)
	rep = &liveUserV1.NoteGetResp_Data{}
	if err != nil {
		log.Error("[GetDaHangHai]live_user.v1.note.get rpc error:%+v", err)
		err = errors.WithMessage(ecode.UserDHHRPCError, fmt.Sprintf("live_user.v1.note.get rpc error:%+v", err))
		return
	}
	if dahanghai.Code != 0 {
		log.Error("[GetDaHangHai]live_user.v1.note response error:%+v,code:%d,msg:%s", err, dahanghai.Code, dahanghai.Msg)
		err = errors.WithMessage(ecode.UserDHHReturnError, fmt.Sprintf("live_user.v1.note response error,code:%d,msg:%s", dahanghai.Code, dahanghai.Msg))
		return
	}

	if dahanghai.Data == nil {
		log.Error("[GetUserTagList]live_user.v1.note empty error")
		err = errors.WithMessage(ecode.UserDHHDataNil, "live_user.v1.note empty error")
		return
	}
	rep = dahanghai.Data

	return
}
