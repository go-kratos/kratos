package dao

import (
	"context"

	"go-common/app/service/main/relation/model"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
	"go-common/library/time"
)

// consts
const (
	RelationLogID = 13
)

// AddFollowingLog is
func (d *Dao) AddFollowingLog(ctx context.Context, rl *model.RelationLog) {
	d.addLog(ctx, RelationLogID, "log_add_following", rl)
	d.addLog(ctx, RelationLogID, "log_follower_incr", rl.Reverse())
}

// DelFollowingLog is
func (d *Dao) DelFollowingLog(ctx context.Context, rl *model.RelationLog) {
	d.addLog(ctx, RelationLogID, "log_del_following", rl)
	d.addLog(ctx, RelationLogID, "log_follower_decr", rl.Reverse())
}

// DelFollowerLog is
func (d *Dao) DelFollowerLog(ctx context.Context, rl *model.RelationLog) {
	d.addLog(ctx, RelationLogID, "log_del_follower", rl)
	d.addLog(ctx, RelationLogID, "log_following_decr", rl.Reverse())
}

// AddWhisperLog is
func (d *Dao) AddWhisperLog(ctx context.Context, rl *model.RelationLog) {
	d.addLog(ctx, RelationLogID, "log_add_whisper", rl)
	d.addLog(ctx, RelationLogID, "log_whisper_follower_incr", rl.Reverse())
}

// DelWhisperLog is
func (d *Dao) DelWhisperLog(ctx context.Context, rl *model.RelationLog) {
	d.addLog(ctx, RelationLogID, "log_del_whisper", rl)
	d.addLog(ctx, RelationLogID, "log_whisper_follower_decr", rl.Reverse())
}

// AddBlackLog is
func (d *Dao) AddBlackLog(ctx context.Context, rl *model.RelationLog) {
	d.addLog(ctx, RelationLogID, "log_add_black", rl)
	d.addLog(ctx, RelationLogID, "log_black_incr", rl.Reverse())
}

// DelBlackLog is
func (d *Dao) DelBlackLog(ctx context.Context, rl *model.RelationLog) {
	d.addLog(ctx, RelationLogID, "log_del_black", rl)
	d.addLog(ctx, RelationLogID, "log_black_decr", rl.Reverse())
}

func (d *Dao) addLog(ctx context.Context, business int, action string, rl *model.RelationLog) {
	t := time.Time(rl.Ts)
	content := make(map[string]interface{}, len(rl.Content))
	for k, v := range rl.Content {
		content[k] = v
	}
	content["from_attr"] = rl.FromAttr
	content["to_attr"] = rl.ToAttr
	content["from_rev_attr"] = rl.FromRevAttr
	content["to_rev_attr"] = rl.ToRevAttr
	content["source"] = rl.Source
	ui := &report.UserInfo{
		Mid:      rl.Mid,
		Platform: "",
		Build:    0,
		Buvid:    rl.Buvid,
		Business: business,
		Type:     0,
		Oid:      rl.Fid,
		Action:   action,
		Ctime:    t.Time(),
		IP:       rl.Ip,
		// extra
		Index:   []interface{}{int64(rl.Source), 0, "", "", ""},
		Content: content,
	}
	report.User(ui)
	log.Info("add log to report: relationlog: %+v userinfo: %+v", rl, ui)
}
