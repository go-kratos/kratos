package dao

import (
	"context"
	"time"

	"go-common/app/interface/main/kvo/model"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getUserConf = "SELECT mid,module_key,check_sum,timestamp FROM user_conf WHERE mid=? AND module_key=?"
	_upUserConf  = "INSERT INTO user_conf(mid,module_key,check_sum,timestamp,ctime,mtime) VALUES(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE check_sum=?, timestamp=?"
)

// UserConf get userconf
func (d *Dao) UserConf(ctx context.Context, mid int64, moduleKey int) (userConf *model.UserConf, err error) {
	row := d.getUserConf.QueryRow(ctx, mid, moduleKey)
	userConf = &model.UserConf{}
	err = row.Scan(&userConf.Mid, &userConf.ModuleKey, &userConf.CheckSum, &userConf.Timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			userConf = nil
			err = nil
			return
		}
		log.Error("row.Scan err:%v", err)
	}
	return
}

// TxUpUserConf add or update user conf
func (d *Dao) TxUpUserConf(ctx context.Context, tx *sql.Tx, mid int64, moduleKey int, checkSum int64, now time.Time) (err error) {
	_, err = tx.Exec(_upUserConf, mid, moduleKey, checkSum, now.Unix(), now, now, checkSum, now.Unix())
	if err != nil {
		log.Error("db.exec err:%v", err)
	}
	return
}
