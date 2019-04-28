package global

import (
	"context"

	accgrpc "go-common/app/service/main/account/api"
	accmdl "go-common/app/service/main/account/model"
	"go-common/library/log"
)

//GetInfo get info
func GetInfo(c context.Context, mid int64) (res *accmdl.Info, err error) {
	if mid == 0 {
		return
	}
	var infoReply *accgrpc.InfoReply
	if infoReply, err = accGRPC.Info3(c, &accgrpc.MidReq{Mid: mid}); err != nil {
		return
	}
	res = infoReply.Info
	return
}

//GetInfos get many infos
func GetInfos(c context.Context, mids []int64) (res map[int64]*accmdl.Info, err error) {
	if len(mids) == 0 {
		return
	}
	var infosReply *accgrpc.InfosReply
	if infosReply, err = accGRPC.Infos3(c, &accgrpc.MidsReq{Mids: mids}); err != nil {
		return
	}
	res = infosReply.Infos
	return
}

//GetName get user name
func GetName(c context.Context, mid int64) (nickname string) {
	accInfos, e := GetInfos(c, []int64{mid})
	if e == nil && accInfos != nil {
		var info, ok = accInfos[mid]
		if ok {
			nickname = info.Name
		}
	} else {
		log.Warn("get up info fail, err=%s", e)
	}
	return
}
