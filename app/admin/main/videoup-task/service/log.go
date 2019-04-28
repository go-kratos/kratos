package service

import (
	"context"
	"strconv"
	"time"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

// send to log service
func (s *Service) sendVideoLog(c context.Context, vp *model.VideoParam, others string) (err error) {
	var (
		v *model.ArcVideo
		a *model.Archive
	)
	if vp.Cid != 0 {
		v, err = s.dao.ArcVideoByCID(c, vp.Cid)
	} else if vp.ID != 0 {
		v, err = s.dao.NewVideoByID(c, vp.ID)
	}
	if err != nil || v == nil {
		v = &model.ArcVideo{} // ignore err
	}

	a, err = s.dao.Archive(c, vp.Aid)
	if err != nil || a == nil {
		a = &model.Archive{} // ignore err
	}
	// send
	logData := &report.ManagerInfo{
		Uname:    vp.Oname,
		UID:      vp.UID,
		Business: model.LogClientVideo,
		Type:     model.LogClientTypeVideo,
		Oid:      vp.Cid,
		Action:   strconv.Itoa(int(vp.Status)),
		Ctime:    time.Now(),
		Index:    []interface{}{int64(vp.Attribute), v.CTime.Unix(), vp.TagID, a.Title, vp.Note},
		Content: map[string]interface{}{
			"content": vp,
			"others":  others,
		},
	}
	report.Manager(logData)
	return
}

// sendConsumerLog send consumer log
func (s *Service) sendConsumerLog(c context.Context, cl *model.ConsumerLog) (err error) {
	logData := &report.ManagerInfo{
		Uname:    cl.Uname,
		UID:      cl.UID,
		Business: model.LogClientConsumer,
		Type:     model.LogClientTypeConsumer,
		Oid:      cl.UID,
		Action:   strconv.Itoa(int(cl.Action)),
		Ctime:    time.Now(),
		Index:    []interface{}{cl.UID, cl.Action, cl.Ctime},
		Content: map[string]interface{}{
			"content": cl,
		},
	}
	report.Manager(logData)
	log.Info("sendConsumerLog logData(%+v)", cl)
	return
}
