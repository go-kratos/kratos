package newbiedao

import (
	"context"

	accApi "go-common/app/service/main/account/api"

	"go-common/library/ecode"
	"go-common/library/log"
)

// GetInfo get info
func (d *Dao) GetInfo(c context.Context, mid int64) (res *accApi.InfoReply, err error) {
	res, err = d.accGRPC.Info3(c, &accApi.MidReq{Mid: mid})
	if err != nil {
		return
	}
	if res == nil || res.Info == nil {
		err = ecode.GrowupUpInfoNotExist
		log.Error("s.dao.GetInfo get up info is empty")
		return
	}
	return
}

// GetInfos get infos
func (d *Dao) GetInfos(c context.Context, mids []int64) (res *accApi.InfosReply, err error) {
	res, err = d.accGRPC.Infos3(c, &accApi.MidsReq{Mids: mids})
	if err != nil {
		return
	}
	if res == nil {
		err = ecode.GrowupUpInfoNotExist
		log.Error("s.dao.GetInfos get ups info are empty")
		return
	}
	return
}
