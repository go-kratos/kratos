package archive

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"go-common/app/job/main/videoup-report/model/task"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getTaskWeight = "SELECT t.id,t.state,a.mid,t.ctime,t.upspecial,t.ptime,e.description FROM `task_dispatch` AS t " +
		"LEFT JOIN `task_dispatch_extend` AS e ON t.id=e.task_id INNER JOIN archive as a ON a.id=t.aid WHERE t.state=0 AND t.id>? LIMIT 1000"
	_inDispatchExtendSQL = "INSERT INTO task_dispatch_extend(task_id,description) VALUE (?,?)"
	_delTaskExtendSQL    = "DELETE FROM task_dispatch_extend WHERE mtime < ? LIMIT ?"
)

// GetTaskWeight 从数据库读取权重配置
func (d *Dao) GetTaskWeight(c context.Context, lastid int64) (mcases map[int64]*task.WeightParams, err error) {
	var (
		rows *xsql.Rows
		desc sql.NullString
	)
	if rows, err = d.db.Query(c, _getTaskWeight, lastid); err != nil {
		log.Error("d.db.Query(%s, %d) error(%v)", _getTaskWeight, lastid, err)
		return
	}
	defer rows.Close()
	mcases = make(map[int64]*task.WeightParams)
	for rows.Next() {
		tp := new(task.WeightParams)
		if err = rows.Scan(&tp.TaskID, &tp.State, &tp.Mid, &tp.Ctime, &tp.Special, &tp.Ptime, &desc); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if desc.Valid && len(desc.String) > 0 {
			arr := []*task.ConfigItem{}
			if err = json.Unmarshal([]byte(desc.String), &arr); err != nil {
				arr = nil
				log.Error("json.Unmarshal error(%v)", err)
				continue
			}
			tp.CfItems = arr
		}
		mcases[tp.TaskID] = tp
	}
	return
}

// InDispatchExtend 扩展表,记录权重配置信息
func (d *Dao) InDispatchExtend(c context.Context, taskid int64, desc string) (lastid int64, err error) {
	res, err := d.db.Exec(c, _inDispatchExtendSQL, taskid, desc)
	if err != nil {
		log.Error("tx.Exec(%s, %d, %v) error(%v)", _inDispatchExtendSQL, taskid, desc, err)
		return
	}
	return res.LastInsertId()
}

// DelTaskExtend del task_dispatch_extend
func (d *Dao) DelTaskExtend(c context.Context, before time.Time, limit int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delTaskExtendSQL, before.Format("2006-01-02 15:04:05"), limit)
	if err != nil {
		log.Error("d.db.Exec(%s, %s, %d) error(%v)", _delTaskExtendSQL, before, limit, err)
		return
	}
	return res.RowsAffected()
}
