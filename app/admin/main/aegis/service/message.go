package service

import (
	"context"
	"strconv"

	"go-common/app/admin/main/aegis/model/databus"
	"go-common/library/log"
)

func (s *Service) sendCreateTaskMsg(c context.Context, rid, flowID, dispatchLimit, bizid int64) (err error) {
	msg := &databus.CreateTaskMsg{
		BizID:         bizid,
		RID:           rid,
		FlowID:        flowID,
		DispatchLimit: dispatchLimit,
	}

	return s.async.Do(c, func(c context.Context) {
		log.Info("start to send msg(%+v)", msg)
		for retry := 0; retry < 3; retry++ {
			if err = s.aegisPub.Send(c, strconv.Itoa(int(msg.RID)), msg); err == nil {
				break
			}
		}
		if err != nil {
			log.Error("s.aegisPub.Send error(%v) msg(%+v) ", err, msg)
		}
	})
}
