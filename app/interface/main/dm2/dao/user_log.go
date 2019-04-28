package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/queue/databus/report"
)

const (
	_businessUserLogDm = 161

	// dm log
	_typeUserLogDm       = 1
	_actionUserLogDmPost = "danmu_post"

	// garbage dm
	_typeUserLogGarbageDm   = 2
	_actionUserLogGarbageDm = "danmu_garbage"
)

func init() {
	report.InitUser(nil)
}

// ReportDmGarbageLog report garbage dm
func (d *Dao) ReportDmGarbageLog(c context.Context, dm *model.DM) (err error) {
	return d.reportUserLog(c, dm, _typeUserLogGarbageDm, _actionUserLogGarbageDm)
}

// ReportDmLog report dm post log
func (d *Dao) ReportDmLog(c context.Context, dm *model.DM) (err error) {
	return d.reportUserLog(c, dm, _typeUserLogDm, _actionUserLogDmPost)
}

func (d *Dao) reportUserLog(c context.Context, dm *model.DM, userLogType int, userLogAction string) (err error) {
	var (
		remoteIP, remoteAddr = metadata.String(c, metadata.RemoteIP), metadata.String(c, metadata.RemotePort)
		msg                  string
	)
	if dm.Content == nil {
		return
	}
	msg = dm.Content.Msg
	if dm.ContentSpe != nil {
		msg = dm.ContentSpe.Msg
	}
	if err = report.User(&report.UserInfo{
		Mid:      dm.Mid,
		Business: _businessUserLogDm,
		Type:     userLogType,
		Oid:      dm.ID,
		Action:   userLogAction,
		Ctime:    time.Now(),
		Platform: fmt.Sprint(dm.Content.Plat),
		IP:       remoteIP,
		Index:    []interface{}{dm.Oid},
		Content: map[string]interface{}{
			"ip":   remoteIP,
			"port": remoteAddr,
			"msg":  msg,
		},
	}); err != nil {
		log.Error("reportUserLog(dmid:%v) error(%v)", dm.ID, err)
		return
	}
	return
}
