package dao

import (
	"context"
	"database/sql"
	"encoding/json"

	"go-common/library/log"
)

const (
	_settingSQL    = "SELECT value FROM push_user_settings WHERE mid=? AND dtime=0"
	_setSettingSQL = "INSERT INTO push_user_settings (mid,value) VALUES (?,?) ON DUPLICATE KEY UPDATE value=?"
)

// SetSetting saves user notify settings.
func (d *Dao) SetSetting(c context.Context, mid int64, st map[int]int) (err error) {
	bs, err := json.Marshal(st)
	if err != nil {
		log.Error("SetSetting(%d) json.Marshal(%v) error(%v)", mid, st, err)
		return
	}
	if _, err = d.setSettingStmt.Exec(c, mid, string(bs), string(bs)); err != nil {
		PromError("mysql:保存用户通知开关")
		log.Error("d.SetSetting(%d,%s) error(%v)", mid, bs, err)
	}
	return
}

// Setting gets user push setting.
func (d *Dao) Setting(c context.Context, mid int64) (st map[int]int, err error) {
	var v string
	if err = d.settingStmt.QueryRow(c, mid).Scan(&v); err != nil {
		if err == sql.ErrNoRows {
			st = nil
			err = nil
			return
		}
		log.Error("d.settingStmt.Query() error(%v)", err)
		PromError("mysql:获取用户通知开关配置")
		return
	}
	if v == "" {
		return
	}
	if err = json.Unmarshal([]byte(v), &st); err != nil {
		log.Error("d.Setting(%d) json.Unmarshal(%s) error(%v)", mid, v, err)
	}
	return
}
