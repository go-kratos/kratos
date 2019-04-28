package archive

import (
	"context"
	"encoding/json"
	"fmt"

	tmod "go-common/app/job/main/videoup-report/model/task"
	"go-common/app/job/main/videoup-report/model/utils"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_assignConfigsSQL = "SELECT id,uids,pool,config_mids,config_tids,config_time,adminid,state,stime,etime FROM task_config WHERE state=0"
	_delAConfsSQL     = "UPDATE task_config SET state=1 WHERE id IN (%s)"
	_weightConfSQL    = "SELECT id,description,mtime FROM task_weight_config WHERE state=0" // 查
	_delWConfsSQL     = "UPDATE task_weight_config SET state=1 WHERE id IN (%s)"
)

//AssignConfigs take config
func (d *Dao) AssignConfigs(c context.Context) (tasks map[int64]*tmod.AssignConfig, err error) {
	rows, err := d.db.Query(c, _assignConfigsSQL)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	tasks = make(map[int64]*tmod.AssignConfig)
	defer rows.Close()
	for rows.Next() {
		var (
			midStr      string
			uidsStr     string
			tidStr      string
			durationStr string
			mids        []int64
			tids        []int64
			durations   []int64
		)
		t := &tmod.AssignConfig{}
		if err = rows.Scan(&t.ID, &uidsStr, &t.Pool, &midStr, &tidStr, &durationStr, &t.AdminID, &t.State, &t.STime, &t.ETime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			continue
		}
		if uidsStr != "" {
			if t.UIDs, err = xstr.SplitInts(uidsStr); err != nil {
				log.Error("xstr.SplitInts(%s) errror(%v)", uidsStr, err)
				err = nil
				continue
			}
		}
		if midStr != "" {
			if mids, err = xstr.SplitInts(midStr); err != nil {
				log.Error("xstr.SplitInts(%s) error(%v)", midStr, err)
				err = nil
				continue
			}
			t.MIDs = make(map[int64]struct{}, len(mids))
			for _, mid := range mids {
				t.MIDs[mid] = struct{}{}
			}
		}
		if tidStr != "" {
			if tids, err = xstr.SplitInts(tidStr); err != nil {
				log.Error("xstr.SplitInts(%s) error(%v)", tidStr, err)
				err = nil
				continue
			}
			t.TIDs = make(map[int16]struct{}, len(tids))
			for _, tid := range tids {
				t.TIDs[int16(tid)] = struct{}{}
			}
		}
		if durationStr != "" {
			if durations, err = xstr.SplitInts(durationStr); err != nil || len(durations) != 2 {
				log.Error("xstr.SplitInts(%s) error(%v)", durationStr, err)
				err = nil
				continue
			}
			t.MinDuration = durations[0]
			t.MaxDuration = durations[1]
		}
		if len(t.UIDs) > 0 {
			tasks[t.ID] = t
		}
	}
	return
}

// DelAssignConfs 删除指派配置
func (d *Dao) DelAssignConfs(c context.Context, ids []int64) (rows int64, err error) {
	sqlstring := fmt.Sprintf(_delAConfsSQL, xstr.JoinInts(ids))
	res, err := d.db.Exec(c, sqlstring)
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", sqlstring, err)
		return
	}
	return res.RowsAffected()
}

// WeightConf 所有有效的配置(用于检测是否和以及有的配置冲突)
func (d *Dao) WeightConf(c context.Context) (items []*tmod.ConfigItem, err error) {
	var (
		id    int64
		descb []byte
		rows  *sql.Rows
		wci   *tmod.ConfigItem
		mtime utils.FormatTime
	)
	if rows, err = d.db.Query(c, _weightConfSQL); err != nil {
		log.Error("d.db.Query(%s) error(%v)", _weightConfSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		wci = new(tmod.ConfigItem)
		if err = rows.Scan(&id, &descb, &mtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if err = json.Unmarshal(descb, wci); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(descb), err)
			continue
		}
		wci.Mtime = mtime
		wci.ID = id
		items = append(items, wci)
	}
	return
}

// DelWeightConfs 删除权重配置
func (d *Dao) DelWeightConfs(c context.Context, ids []int64) (rows int64, err error) {
	sqlstring := fmt.Sprintf(_delWConfsSQL, xstr.JoinInts(ids))
	res, err := d.db.Exec(c, sqlstring)
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", sqlstring, err)
		return
	}
	return res.RowsAffected()
}
