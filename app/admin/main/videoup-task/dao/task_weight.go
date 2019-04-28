package dao

import (
	"context"
	xsql "database/sql"
	"encoding/json"
	"fmt"
	"go-common/library/xstr"
	"strings"
	"time"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/cache/redis"
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

	// archive_config
	_confSQL   = "SELECT value FROM archive_config WHERE state=0 AND name=?"
	_upconfSQL = "UPDATE archive_config SET value=?,remark=? WHERE name=?"
	_inconfSQL = "INSERT archive_config(value,remark,name,state) VALUE (?,?,?,0)"

	_twexpire = 24 * 60 * 60 // 1 day
)

// GetMaxWeight 获取当前最大权重数值
func (d *Dao) GetMaxWeight(c context.Context) (max int64, err error) {
	if err = d.arcDB.QueryRow(c, _getMaxWeightSQL).Scan(&max); err != nil {
		log.Error("d.arcDB.QueryRow error(%v)", err)
		err = nil
	}
	return
}

// UpCwAfterAdd update config weight after add config
func (d *Dao) UpCwAfterAdd(c context.Context, id int64, desc string) (rows int64, err error) {
	row, err := d.arcDB.Exec(c, _upCwAfterAddSQL, id, desc, desc)
	if err != nil {
		log.Error("arcDB.Exec(%s,%d,%s,%s) error(%v)", _upCwAfterAddSQL, id, desc, desc, err)
		return
	}
	return row.RowsAffected()
}

// InWeightConf 写入权重配置表
func (d *Dao) InWeightConf(c context.Context, mcases map[int64]*model.WCItem) (err error) {
	tx, err := d.arcDB.Begin(c)
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
	res, err := d.arcDB.Exec(c, _delWeightConfSQL, id)
	if err != nil {
		log.Error("tx.Exec(%s %d) error(%v)", _delWeightConfSQL, id, err)
		return
	}
	return res.RowsAffected()
}

// ListWeightConf 查看权重配置表列表
func (d *Dao) ListWeightConf(c context.Context, cf *model.Confs) (citems []*model.WCItem, err error) {
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
	rows, err = d.arcDB.Query(c, sqlstring)
	if err != nil {
		log.Error("d.arcDB.Query(%s) error(%v)", sqlstring, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		wci := &model.WCItem{}
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
func (d *Dao) WeightConf(c context.Context) (items []*model.WCItem, err error) {
	var (
		id    int64
		descb []byte
		rows  *sql.Rows
		wci   *model.WCItem
	)
	if rows, err = d.arcDB.Query(c, _WeightConfSQL); err != nil {
		log.Error("d.arcDB.Query(%s) error(%v)", _WeightConfSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		wci = new(model.WCItem)
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
	rows, err = d.arcDB.Query(c, fmt.Sprintf(_lwconfigHelpSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("d.arcDB.Query(%v) error(%v)", ids, err)
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

func key(id int64) string {
	return fmt.Sprintf("tw_%d", id)
}

//SetWeightRedis 设置权重配置
func (d *Dao) SetWeightRedis(c context.Context, mcases map[int64]*model.TaskPriority) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	for tid, mcase := range mcases {
		var bs []byte
		key := key(tid)
		if bs, err = json.Marshal(mcase); err != nil {
			log.Error("json.Marshal(%+v) error(%v)", mcase, err)
			continue
		}

		if err = conn.Send("SET", key, bs); err != nil {
			log.Error("SET error(%v)", err)
			continue
		}
		if err = conn.Send("EXPIRE", key, _twexpire); err != nil {
			log.Error("EXPIRE error(%v)", err)
			continue
		}
	}

	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2*len(mcases); i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

//GetWeightRedis 获取实时任务的权重配置
func (d *Dao) GetWeightRedis(c context.Context, ids []int64) (mcases map[int64]*model.TaskPriority, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	mcases = make(map[int64]*model.TaskPriority)
	for _, id := range ids {
		var bs []byte
		key := key(int64(id))
		if bs, err = redis.Bytes(conn.Do("GET", key)); err != nil {
			if err == redis.ErrNil {
				err = nil
			} else {
				log.Error("conn.Do(GET, %v) error(%v)", key, err)
			}
			continue
		}
		p := &model.TaskPriority{}
		if err = json.Unmarshal(bs, p); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(bs), err)
			err = nil
			continue
		}
		mcases[int64(id)] = p
	}
	return
}

// WeightVC 获取权重分值
func (d *Dao) WeightVC(c context.Context) (wvc *model.WeightVC, err error) {
	var value []byte
	row := d.arcDB.QueryRow(c, _confSQL, model.ConfForWeightVC)
	if err = row.Scan(&value); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	wvc = new(model.WeightVC)
	if err = json.Unmarshal(value, wvc); err != nil {
		log.Error("json.Unmarshal error(%v)", err)
		wvc = nil
	}
	return
}

// SetWeightVC 设置权重分值
func (d *Dao) SetWeightVC(c context.Context, wvc *model.WeightVC, desc string) (rows int64, err error) {
	var (
		valueb []byte
		res    xsql.Result
	)
	if valueb, err = json.Marshal(wvc); err != nil {
		log.Error("json.Marshal(%+v) error(%v)", wvc, err)
		return
	}

	if res, err = d.arcDB.Exec(c, _upconfSQL, string(valueb), desc, model.ConfForWeightVC); err != nil {
		log.Error("d.arcDB.Exec(%s, %s, %s, %s) error(%v)", _upconfSQL, string(valueb), desc, model.ConfForWeightVC, err)
		return
	}
	return res.RowsAffected()
}

// InWeightVC 插入
func (d *Dao) InWeightVC(c context.Context, wvc *model.WeightVC, desc string) (rows int64, err error) {
	var (
		valueb []byte
		res    xsql.Result
	)
	if valueb, err = json.Marshal(wvc); err != nil {
		log.Error("json.Marshal(%+v) error(%v)", wvc, err)
		return
	}

	if res, err = d.arcDB.Exec(c, _inconfSQL, string(valueb), desc, model.ConfForWeightVC); err != nil {
		log.Error("d.arcDB.Exec(%s, %s, %s, %s) error(%v)", _inconfSQL, string(valueb), desc, model.ConfForWeightVC, err)
		return
	}
	return res.LastInsertId()
}
