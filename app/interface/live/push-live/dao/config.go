package dao

import (
	"context"
	"go-common/app/interface/live/push-live/model"
	"go-common/library/log"
)

const (
	_getPushConfig   = "SELECT `type` FROM ap_push_config WHERE value=? ORDER BY `order` ASC" // 获取推送选项配置
	_getPushInterval = "SELECT `order` FROM ap_push_config WHERE type=?"                      // 获取推送间隔时间
)

// GetPushConfig 从DB中获取推送配置
func (d *Dao) GetPushConfig(c context.Context) (types []string, err error) {
	var t string
	types = make([]string, 0)
	rows, err := d.db.Query(c, _getPushConfig, model.LivePushConfigOn)
	if err != nil {
		log.Error("[dao.config|GetPushConfig] db.Query() error(%v)", err)
		return
	}
	for rows.Next() {
		if err = rows.Scan(&t); err != nil {
			log.Error("[dao.config|GetPushConfig] rows.Scan() error(%v)", err)
			return
		}
		types = append(types, t)
	}
	return
}

// GetPushInterval 获取推送时间间隔
func (d *Dao) GetPushInterval(c context.Context) (interval int32, err error) {
	var i int32
	row := d.db.QueryRow(c, _getPushInterval, model.PushIntervalKey)
	if err = row.Scan(&i); err != nil {
		log.Error("[dao.config|GetPushInterval] row.Scan() error(%v)", err)
		return
	}
	interval = i * 60 // min to sec
	return
}
