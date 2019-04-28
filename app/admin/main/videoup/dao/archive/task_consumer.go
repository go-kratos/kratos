package archive

import (
	"context"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_taskUserCheckInSQL  = "INSERT INTO task_consumer (uid,state) VALUES (?,1) ON DUPLICATE KEY UPDATE state = 1"
	_taskUserCheckOffSQL = "UPDATE task_consumer SET state = 0 WHERE uid=?"
	_consumersSQL        = "SELECT id,uid,state,ctime,mtime FROM task_consumer where state=1"
	_isConsumerOnSQL     = "SELECT state FROM task_consumer WHERE uid=?"
)

// TaskUserCheckIn insert or update task consumer check state
func (d *Dao) TaskUserCheckIn(c context.Context, uid int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _taskUserCheckInSQL, uid)
	if err != nil {
		log.Error("tx.Exec(%s, %d) error(%v)", _taskUserCheckInSQL, uid, err)
		return
	}
	return res.RowsAffected()
}

// TaskUserCheckOff update task consumer check state
func (d *Dao) TaskUserCheckOff(c context.Context, uid int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _taskUserCheckOffSQL, uid)
	if err != nil {
		log.Error("tx.Exec(%s, %d) error(%v)", _taskUserCheckOffSQL, uid, err)
		return
	}
	return res.RowsAffected()
}

// Consumers 用户列表
func (d *Dao) Consumers(c context.Context) (cms []*archive.Consumers, err error) {
	rows, err := d.rddb.Query(c, _consumersSQL)
	if err != nil {
		log.Error("d.rddb.Query(%s) error(%v)", _consumersSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		cm := new(archive.Consumers)
		err = rows.Scan(&cm.ID, &cm.UID, &cm.State, &cm.Ctime, &cm.Mtime)
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		cms = append(cms, cm)
	}
	return
}

// IsConsumerOn 判断是否登入
func (d *Dao) IsConsumerOn(c context.Context, uid int64) (state int8) {
	err := d.rddb.QueryRow(c, _isConsumerOnSQL, uid).Scan(&state)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error("d.rddb.QueryRow error(%v)", err)
		}
	}
	return
}
