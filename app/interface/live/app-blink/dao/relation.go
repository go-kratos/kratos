package dao

import (
	"context"

	"github.com/pkg/errors"

	"go-common/app/service/live/relation/api/liverpc/v1"
	"go-common/library/log"
)

//GetUserFc 获取用户关注数
func (d *Dao) GetUserFc(c context.Context, uid int64) (res *v1.FeedGetUserFcResp_Data, err error) {
	reply, err := d.RelationApi.V1Feed.GetUserFc(c, &v1.FeedGetUserFcReq{Follow: uid})
	if err != nil {
		log.Error("relation_GetUserFc_error:%v", err)
		return
	}
	if reply.Code != 0 {
		err = errors.New("code error")
		log.Error("relation_gGetUserFc_error:%d,%s,%v", reply.Code, reply.Msg, reply.Data)
		return
	}
	if reply.Data == nil {
		err = errors.New("data error")
		log.Error("relation_gGetUserFc_error:%d,%s,%v", reply.Code, reply.Msg, reply.Data)
		return
	}
	log.Info("relation_gGetUserFc:%d,%s,$v", reply.Code, reply.Msg, reply.Data)
	res = reply.Data
	return
}
