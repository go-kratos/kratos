package archive

import (
	"context"
	"go-common/library/database/sql"
	"go-common/library/log"
	"time"
)

const (
	_slByAID     = "SELECT `id` FROM `archive_first_pass` WHERE `aid`=? LIMIT 1;"
	_inFirstPass = "INSERT INTO `archive_first_pass`(`aid`, `ctime`, `mtime`) VALUES(?,?,?);"
)

//GetFirstPassByAID 根据aid获取第一次过审的记录
func (d *Dao) GetFirstPassByAID(c context.Context, aid int64) (id int64, err error) {
	row := d.db.QueryRow(c, _slByAID, aid)
	if err = row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("GetFirstPassByAID error(%v) aid(%d)", err, aid)
		}
		return
	}

	return
}

//AddFirstPass 添加一条 第一次过审的记录
func (d *Dao) AddFirstPass(tx *sql.Tx, aid int64) (err error) {
	now := time.Now()
	if _, err = tx.Exec(_inFirstPass, aid, now, now); err != nil {
		log.Error("AddFirstPass error(%v) aid(%d)", err)
	}

	return
}
