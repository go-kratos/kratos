package dao

import (
	"context"

	"go-common/app/interface/main/dm/model"

	"go-common/library/log"
)

// SendAction send action to job.
func (d *Dao) SendAction(c context.Context, k string, act *model.ReportAction) (err error) {
	if err = d.databus.Send(c, k, act); err != nil {
		log.Error("actionPub.Send(data:%v) error(%v)", act, err)
	} else {
		log.Info("actionPub.Send(action:%v) success", act)
	}
	return
}
