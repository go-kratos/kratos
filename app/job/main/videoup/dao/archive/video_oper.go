package archive

import (
	"context"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_inVideoOperSQL = "INSERT INTO archive_video_oper(aid,uid,vid,status,attribute,last_id) VALUES(?,399,?,?,?,1)"
)

//TranVideoOper insert a record
func (d *Dao) TranVideoOper(c context.Context, tx *sql.Tx, aid, vid int64, status int16, attr int32) (rows int64, err error) {
	res, err := tx.Exec(_inVideoOperSQL, aid, vid, status, attr)
	if err != nil {
		log.Error("tx.Exec(%s, %d, %d, %d, %d) error(%v)", _inVideoOperSQL, aid, vid, status, attr, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}
