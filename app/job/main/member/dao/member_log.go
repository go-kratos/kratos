package dao

import (
	"context"

	"go-common/app/job/main/member/model"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
	"go-common/library/time"
)

// consts
const (
	ExpLogID   = 11
	MoralLogID = 12
)

// AddExpLog is
func (d *Dao) AddExpLog(ctx context.Context, ul *model.UserLog) {
	d.addLog(ctx, ExpLogID, "log_exp_change", ul)
}

// AddMoralLog is
func (d *Dao) AddMoralLog(ctx context.Context, ul *model.UserLog) {
	d.addLog(ctx, MoralLogID, "log_moral_change", ul)
}

// AddExpLog is
func (d *Dao) addLog(ctx context.Context, business int, action string, ul *model.UserLog) {
	t := time.Time(ul.TS)
	content := make(map[string]interface{}, len(ul.Content))
	for k, v := range ul.Content {
		content[k] = v
	}
	if ul.LogID == "" {
		ul.LogID = model.UUID4()
	}
	content["log_id"] = ul.LogID
	ui := &report.UserInfo{
		Mid:      ul.Mid,
		Platform: "",
		Build:    0,
		Buvid:    "",
		Business: business,
		Type:     0,
		Oid:      0,
		Action:   action,
		Ctime:    t.Time(),
		IP:       ul.IP,
		// extra
		Index:   []interface{}{ul.LogID},
		Content: content,
	}
	report.User(ui)
	log.Info("add log to report: userlog: %+v userinfo: %+v", ul, ui)
}
