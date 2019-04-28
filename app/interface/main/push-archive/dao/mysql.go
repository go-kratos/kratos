package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-common/library/xstr"

	"go-common/app/interface/main/push-archive/model"
	"go-common/library/log"
	"strconv"
)

const (
	_settingSQL       = `SELECT value FROM push_settings WHERE mid=? and dtime=0 limit 1`
	_setSettingSQL    = `INSERT INTO push_settings (mid,value) VALUES (?,?) ON DUPLICATE KEY UPDATE value=?`
	_settingsSQL      = `SELECT mid,value FROM push_settings WHERE mid IN(%s) and dtime=0`
	_settingsAllSQL   = `SELECT mid,value FROM push_settings WHERE id > %s AND id <= %s`
	_settingsMaxIDSQL = `SELECT MAX(id) AS mx FROM push_settings`
)

// Setting gets the setting.
func (d *Dao) Setting(c context.Context, mid int64) (st *model.Setting, err error) {
	var v string
	if err = d.settingStmt.QueryRow(c, mid).Scan(&v); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("d.Setting(%d) error(%v)", mid, err)
		PromError("db:获取用户配置")
		return
	}
	st = new(model.Setting)
	if err = json.Unmarshal([]byte(v), &st); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", v, err)
	}
	return
}

// SetSetting saves the setting.
func (d *Dao) SetSetting(c context.Context, mid int64, st *model.Setting) (err error) {
	v, err := json.Marshal(st)
	if err != nil {
		log.Error("json.Marshal error(%v)", err)
		return
	}
	if _, err = d.setSettingStmt.Exec(c, mid, v, v); err != nil {
		log.Error("setSetting Exec mid(%d) error(%v)", mid, err)
		PromError("db:保存用户设置")
	}
	return
}

// Settings gets the settings.
func (d *Dao) Settings(c context.Context, mids []int64) (res map[int64]*model.Setting, err error) {
	res = make(map[int64]*model.Setting, len(mids))
	rows, err := d.db.Query(c, fmt.Sprintf(_settingsSQL, xstr.JoinInts(mids)))
	if err != nil {
		log.Error("d.db.Query() error(%v)", err)
		PromError("db:批量查询用户设置")
		return
	}
	for rows.Next() {
		var mid int64
		var v string
		if err = rows.Scan(&mid, &v); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			PromError("db:批量查询用户设置")
			return
		}
		st := new(model.Setting)
		if err = json.Unmarshal([]byte(v), &st); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", v, err)
			return
		}
		res[mid] = st
	}
	return
}

// SettingsAll gets all settings.
func (d *Dao) SettingsAll(c context.Context, startID int64, endID int64, res *map[int64]*model.Setting) (err error) {
	start := strconv.FormatInt(startID, 10)
	end := strconv.FormatInt(endID, 10)
	rows, err := d.db.Query(c, fmt.Sprintf(_settingsAllSQL, start, end))
	if err != nil {
		log.Error("d.db.Query() error(%v)", err)
		PromError("db:查询全部用户设置")
		return
	}

	for rows.Next() {
		var mid int64
		var v string
		if err = rows.Scan(&mid, &v); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			PromError("db:查询用户设置")
			return
		}
		st := new(model.Setting)
		if err = json.Unmarshal([]byte(v), &st); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", v, err)
			return
		}
		(*res)[mid] = st
	}

	return
}

//SettingsMaxID get settings' total number by max(id)
func (d *Dao) SettingsMaxID(c context.Context) (mx int64, err error) {
	if err = d.settingsMaxIDStmt.QueryRow(c).Scan(&mx); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("d.settingsMaxIDStmt.QueryRow.Scan error(%v)", err)
		PromError("db:查询用户最大ID")
		return
	}

	return
}
