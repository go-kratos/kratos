package archive

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/library/xstr"
	"strings"
	"time"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getMaxWeightSQL   = "SELECT MAX(weight) FROM task_dispatch WHERE state in (0,1)"
	_upCwAfterAddSQL   = "INSERT INTO `task_dispatch_extend` (`task_id`,`description`) VALUES(?,?) ON DUPLICATE KEY UPDATE description=?"
	_inWeightConfSQL   = "INSERT INTO task_weight_config(mid,rule,weight,uid,uname,radio,description) VALUES (?,?,?,?,?,?,?)" // 增
	_delWeightConfSQL  = "UPDATE task_weight_config SET state=1 WHERE id=?"                                                   // 软删
	_listWeightConfSQL = "SELECT id,uname,state,rule,weight,mtime,description FROM task_weight_config"                        // 查
	_WeightConfSQL     = "SELECT id,description FROM task_weight_config WHERE state=0"                                        // 查
	_lwconfigHelpSQL   = "SELECT t.id,t.cid,a.title,v.filename FROM task_dispatch t INNER JOIN archive a ON t.aid=a.id INNER JOIN archive_video v ON t.cid=v.cid WHERE t.id IN (%s)"
)

// GetMaxWeight 获取当前最大权重数值
func (d *Dao) GetMaxWeight(c context.Context) (max int64, err error) {
	if err = d.rddb.QueryRow(c, _getMaxWeightSQL).Scan(&max); err != nil {
		log.Error("d.rddb.QueryRow error(%v)", err)
		err = nil
	}
	return
}

// UpCwAfterAdd update config weight after add config
func (d *Dao) UpCwAfterAdd(c context.Context, id int64, desc string) (rows int64, err error) {
	row, err := d.db.Exec(c, _upCwAfterAddSQL, id, desc, desc)
	if err != nil {
		log.Error("db.Exec(%s,%d,%s,%s) error(%v)", _upCwAfterAddSQL, id, desc, desc, err)
		return
	}
	return row.RowsAffected()
}

// InWeightConf 写入权重配置表
func (d *Dao) InWeightConf(c context.Context, mcases map[int64]*archive.WCItem) (err error) {
	tx, err := d.db.Begin(c)
	if err != nil {
		log.Error("db.Begin() error(%v)", err)
		return
	}

	for _, item := range mcases {
		var descb []byte
		if descb, err = json.Marshal(item); err != nil {
			log.Error("json.Marshal(%+v) error(%v)", item, err)
			tx.Rollback()
			return
		}
		if _, err = tx.Exec(_inWeightConfSQL, item.CID, item.Rule, item.Weight, item.UID, item.Uname, item.Radio, string(descb)); err != nil {
			log.Error("db.Exec(%s) error(%v)", _inWeightConfSQL, err)
			tx.Rollback()
			return
		}
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	return
}

// DelWeightConf 删除权重配置
func (d *Dao) DelWeightConf(c context.Context, id int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delWeightConfSQL, id)
	if err != nil {
		log.Error("tx.Exec(%s %d) error(%v)", _delWeightConfSQL, id, err)
		return
	}
	return res.RowsAffected()
}

// ListWeightConf 查看权重配置表列表
func (d *Dao) ListWeightConf(c context.Context, cf *archive.Confs) (citems []*archive.WCItem, err error) {
	var (
		count     int64
		rows      *sql.Rows
		where     string
		wherecase []string
		descb     []byte
		bt        = cf.Bt.TimeValue()
		et        = cf.Et.TimeValue()
	)
	if cid := cf.Cid; cid != -1 {
		wherecase = append(wherecase, fmt.Sprintf("mid=%d", cid))
	}
	if operator := cf.Operator; len(operator) > 0 {
		wherecase = append(wherecase, fmt.Sprintf("uname='%s'", operator))
	}
	if rule := cf.Rule; rule != -1 {
		wherecase = append(wherecase, fmt.Sprintf("rule=%d", rule))
	}

	wherecase = append(wherecase, fmt.Sprintf("radio=%d AND state=%d", cf.Radio, cf.State))
	where = "WHERE " + strings.Join(wherecase, " AND ")

	sqlstring := fmt.Sprintf("%s %s LIMIT %d,%d", _listWeightConfSQL, where, (cf.Pn-1)*cf.Ps, cf.Pn*cf.Ps)
	rows, err = d.rddb.Query(c, sqlstring)
	if err != nil {
		log.Error("d.rddb.Query(%s) error(%v)", sqlstring, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		wci := &archive.WCItem{}
		if err = rows.Scan(&wci.ID, &wci.Uname, &wci.State, &wci.Rule, &wci.Weight, &wci.Mtime, &descb); err != nil {
			log.Error("rows.Scan(%s) error(%v)", sqlstring, err)
			return
		}
		if len(descb) > 0 {
			if err = json.Unmarshal(descb, wci); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", string(descb), err)
				err = nil
				continue
			}
			eti := wci.Et.TimeValue()
			// filter time
			if !et.IsZero() && !bt.IsZero() && (bt.After(wci.Mtime.TimeValue()) || et.Before(wci.Mtime.TimeValue())) {
				log.Info("config expired (%+v) parse et(%v)", wci, et)
				continue
			}
			// filter state
			if cf.State == 0 && !eti.IsZero() && eti.Before(time.Now()) {
				log.Info("config expired (%+v) parse et(%v)", wci, eti)
				continue
			}
		}

		if count > 50 {
			break
		}

		count++
		citems = append(citems, wci)
	}

	return
}

// WeightConf 所有有效的配置(用于检测是否和已有的配置冲突)
func (d *Dao) WeightConf(c context.Context) (items []*archive.WCItem, err error) {
	var (
		id    int64
		descb []byte
		rows  *sql.Rows
		wci   *archive.WCItem
	)
	if rows, err = d.rddb.Query(c, _WeightConfSQL); err != nil {
		log.Error("d.rddb.Query(%s) error(%v)", _WeightConfSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		wci = new(archive.WCItem)
		if err = rows.Scan(&id, &descb); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if err = json.Unmarshal(descb, wci); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(descb), err)
			err = nil
			continue
		}
		wci.ID = id
		items = append(items, wci)
	}
	return
}

// LWConfigHelp 补充任务对应稿件的title和filename
func (d *Dao) LWConfigHelp(c context.Context, ids []int64) (res map[int64][]interface{}, err error) {
	var (
		taskid, vid     int64
		filename, title string
		rows            *sql.Rows
	)
	rows, err = d.rddb.Query(c, fmt.Sprintf(_lwconfigHelpSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("d.db.Query(%v) error(%v)", ids, err)
		return
	}
	defer rows.Close()

	res = make(map[int64][]interface{})
	for rows.Next() {
		err = rows.Scan(&taskid, &vid, &title, &filename)
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
			continue
		}
		res[taskid] = []interface{}{filename, title, vid}
	}
	return
}
