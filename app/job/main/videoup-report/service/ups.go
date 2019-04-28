package service

import (
	"context"
	"encoding/json"
	"go-common/app/job/main/videoup-report/model/manager"
	"go-common/library/log"
)

// hdlManagerUpsBinlog UP主分组表变更事件
func (s *Service) hdlManagerUpsBinlog(msg *manager.BinMsg) {
	defer func() {
		if pErr := recover(); pErr != nil {
			log.Error("s.hdlManagerUpsBinlog() panic(%v)", pErr)
		}
	}()
	var (
		oldUps, newUps *manager.Ups
		err            error
		c              = context.TODO()
	)
	if msg.New != nil {
		if err = json.Unmarshal(msg.New, &newUps); err != nil {
			log.Error("s.hdlManagerUpsBinlog() json.Unmarshal error(%v) msg.new(%s)", err, string(msg.New))
			return
		}
	}
	if msg.Old != nil {
		if err = json.Unmarshal(msg.Old, &oldUps); err != nil {
			log.Error("s.hdlManagerUpsBinlog() json.Unmarshal error(%v) msg.old(%s)", err, string(msg.Old))
			return
		}
	}
	if newUps != nil && newUps.Type == manager.UpTypeExcitationWhite { //新增UP主到激励回查白名单，需要将该UP主的未回查的稿件从激励回查去除
		go s.ignoreUpsExcitation(c, newUps.MID)
	}
}
