package block

import (
	"context"
	"strconv"
	"time"

	model "go-common/app/job/main/member/model/block"
	"go-common/library/log"
	manager "go-common/library/queue/databus/report"
)

// AddAuditLog .
func (s *Service) AddAuditLog(c context.Context, tp model.BlockAction, uid int64, uname string, oids []int64, duration time.Duration, source model.BlockSource, area model.BlockArea, reason, comment string, notify bool, stime time.Time) error {
	var (
		err       error
		dur       = int64(duration / time.Second)
		notifyStr = strconv.FormatBool(notify)
	)
	for _, oid := range oids {
		managerInfo := &manager.ManagerInfo{
			UID:      uid,
			Uname:    uname,
			Business: model.BlockLogBizID,
			Type:     int(tp),
			Action:   tp.String(),
			Oid:      oid,
			Ctime:    time.Now(),
			Index:    []interface{}{dur, uint8(source), uint8(area), reason, comment, notifyStr},
			Content: map[string]interface{}{
				"duration":    dur,
				"source":      source,
				"area":        area,
				"reason":      reason,
				"comment":     comment,
				"notify":      notifyStr,
				"action_time": stime.Unix(),
				"remove_time": stime.Add(time.Second * time.Duration(dur)).Unix(),
			},
		}
		if err = manager.Manager(managerInfo); err != nil {
			log.Error("manager.Manager(%+v) error(%+v)", managerInfo, err)
			continue
		}
		log.Info("s.managerSendLog(%+v)", managerInfo)
	}
	return err
}
