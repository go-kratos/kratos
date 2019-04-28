package dao

import (
	"context"
	"go-common/app/service/live/userexp/model"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
	"go-common/library/time"
)

// consts
const (
	LiveUserExpID = 104 //http://info.bilibili.co/pages/viewpage.action?pageId=8731603
)

// AddUserExpLog 加用户等级经验行为日志 is
func (d *Dao) AddUserExpLog(ctx context.Context, expInfo *model.ExpLog) {
	d.addLog(ctx, LiveUserExpID, "exp_change", expInfo, "增加用户经验")
}

// AddAnchorExpLog 加主播等级经验行为日志 is
func (d *Dao) AddAnchorExpLog(ctx context.Context, expInfo *model.ExpLog) {
	d.addLog(ctx, LiveUserExpID, "exp_change", expInfo, "增加主播经验")
}

func (d *Dao) addLog(ctx context.Context, business int, action string, expInfo *model.ExpLog, desc string) {
	t := time.Time(expInfo.Ts)
	content := make(map[string]interface{}, len(expInfo.Content))
	for k, v := range expInfo.Content {
		content[k] = v
	}
	ui := &report.UserInfo{
		Mid:      expInfo.Mid,
		Platform: desc,
		Build:    0,
		Buvid:    expInfo.Buvid,
		Business: business,
		Type:     0,
		Action:   action,
		Ctime:    t.Time(),
		IP:       expInfo.Ip,
		// extra
		Index:   []interface{}{int64(expInfo.Mid), 0, "", "", ""},
		Content: content,
	}
	report.User(ui)
	log.Info("add log to report: userexplog: %+v userinfo: %+v,error(%v)", expInfo, ui)
}
