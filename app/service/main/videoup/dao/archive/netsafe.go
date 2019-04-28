package archive

import (
	"context"
	"go-common/library/log"
)

const (
	_inNetsafeSQL = "INSERT INTO netsafe (nid,md5) VALUES (?,?)"
)

// AddNetSafeMd5 fn
func (d *Dao) AddNetSafeMd5(c context.Context, nid int64, md5 string) (rows int64, err error) {
	res, err := d.db.Exec(c, _inNetsafeSQL, nid, md5)
	if err != nil {
		log.Error("_inNetsafeSQL.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}
