package dao

import (
	"context"
	"strconv"

	"go-common/app/service/main/push/model"
	"go-common/library/log"
)

// PubReport add report to databus.
func (d *Dao) PubReport(c context.Context, info *model.Report) (err error) {
	if err = d.reportPub.Send(c, info.Buvid, info); err != nil {
		PromError("databus:发送上报的设备信息")
		log.Error("d.reportPub.Send(%+v) error(%v)", info, err)
		return
	}
	PromInfo("databus:发送上报的设备信息")
	log.Info("PubReport(%+v) success.", info)
	return
}

// PubCallback add push arrive/click callback to databus.
func (d *Dao) PubCallback(c context.Context, v []*model.Callback) (err error) {
	if err = d.callbackPub.Send(c, strconv.Itoa(len(v)), v); err != nil {
		PromError("databus:发送callback")
		log.Error("d.callbackPub.Send(%+v) error(%v)", v, err)
		return
	}
	PromInfo("databus:发送callback")
	log.Info("PubCallback(%+v) success.", v)
	return
}
