package dao

import (
	"context"

	"go-common/app/service/main/member/model"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
	"go-common/library/time"
)

// consts
const (
	MoralLogID = 12
)

// AddMoralLogReport is
func (d *Dao) AddMoralLogReport(ctx context.Context, ul *model.UserLog) {
	d.addLogReport(ctx, MoralLogID, "log_moral_change", ul)
}

// addLogReport is
func (d *Dao) addLogReport(ctx context.Context, business int, action string, ul *model.UserLog) {
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
