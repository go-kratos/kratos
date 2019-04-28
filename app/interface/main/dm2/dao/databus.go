package dao

import (
	"context"
	"encoding/json"
	"strconv"

	"go-common/app/interface/main/dm2/model"

	"go-common/library/log"
)

// PubDatabus pub cache update message to databus.
func (d *Dao) PubDatabus(c context.Context, tp int32, pid, oid, cnt, n, duration int64) (err error) {
	var (
		jobParams = &model.JobParam{
			Type:     tp,
			Oid:      oid,
			Pid:      pid,
			Cnt:      cnt,
			Num:      n,
			Duration: duration,
		}
	)
	value, err := json.Marshal(jobParams)
	if err != nil {
		log.Error("json.Marshal(%v) error(%v)", jobParams, err)
		return
	}
	msg := model.Action{Action: model.ActionIdx, Data: value}
	if err = d.databus.Send(c, strconv.FormatInt(oid, 10), msg); err != nil {
		log.Error("databus.Send(%v) error(%v)", msg, err)
	}
	return
}

// SendAction send action to job.
func (d *Dao) SendAction(c context.Context, k string, act *model.Action) (err error) {
	if err = d.actionPub.Send(c, k, act); err != nil {
		log.Error("actionPub.Send(action:%s,data:%s) error(%v)", act.Action, act.Data, err)
	} else {
		log.Info("actionPub.Send(action:%s,data:%s) success", act.Action, act.Data)
	}
	return
}

// SendSubtitleCheck .
func (d *Dao) SendSubtitleCheck(c context.Context, key string, msg *model.SubtitleCheckMsg) (err error) {
	if err = d.subtitleCheckPub.Send(c, key, msg); err != nil {
		log.Error("actionPub.Send(key:%s,msg:%+v) error(%v)", key, msg, err)
	} else {
		log.Error("actionPub.Send(key:%s,msg:%+v) success", key, msg)
	}
	return
}
