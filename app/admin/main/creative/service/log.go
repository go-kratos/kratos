package service

import (
	"time"

	"go-common/app/admin/main/creative/model/academy"
	"go-common/app/admin/main/creative/model/logcli"
	"go-common/app/admin/main/creative/model/music"
	"go-common/app/admin/main/creative/model/task"
	"go-common/library/conf/env"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus/report"
)

// SendAcademyLog academy log to manager.
func (s *Service) SendAcademyLog(c *bm.Context, lp *academy.LogParam) (err error) {
	if env.DeployEnv == env.DeployEnvDev {
		return
	}
	logData := &report.ManagerInfo{
		Business: academy.LogClientAcademy,
		Type:     0,
		Uname:    lp.UName,
		UID:      lp.UID,
		Oid:      lp.TID,
		Action:   lp.Action,
		Ctime:    time.Now(),
		Index:    []interface{}{lp.TID},
		Content: map[string]interface{}{
			"academy_manager": lp,
		},
	}
	log.Info("s.SendAcademyLog logData(%+v) log param(%+v)", logData, lp)
	report.Manager(logData)
	return
}

// SendMusicLog send to log archive music
func (s *Service) SendMusicLog(c *bm.Context, clientType int, ap *music.LogParam) (err error) {
	logData := &report.ManagerInfo{
		Uname:    ap.UName,
		UID:      ap.UID,
		Business: logcli.LogClientArchiveMusic,
		Type:     clientType,
		Oid:      ap.ID,
		Action:   ap.Action,
		Ctime:    time.Now(),
		Index:    []interface{}{ap.ID},
		Content: map[string]interface{}{
			"object": ap,
		},
	}
	log.Info("sendMusicLog logData(%+v) ap(%+v)", logData, ap)
	report.Manager(logData)
	return
}

// SendTaskLog task log to manager.
func (s *Service) SendTaskLog(c *bm.Context, lp *task.LogParam) (err error) {
	if env.DeployEnv == env.DeployEnvDev {
		return
	}
	logData := &report.ManagerInfo{
		Business: task.LogClientTask,
		Type:     0,
		Uname:    lp.UName,
		UID:      lp.UID,
		Oid:      lp.OID,
		Action:   lp.Action,
		Ctime:    time.Now(),
		Index:    []interface{}{lp.OID},
		Content: map[string]interface{}{
			"params": lp.Content,
		},
	}
	log.Info("s.SendTaskLog logData(%+v) log param(%+v)", logData, lp)
	report.Manager(logData)
	return
}
