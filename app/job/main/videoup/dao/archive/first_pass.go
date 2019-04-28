package archive

import (
	"context"
	"fmt"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
	"time"
)

const (
	_slFirstPassByAID = "SELECT `id` FROM `archive_first_pass` WHERE `aid`=? LIMIT 1;"
	_firstPassCount   = "SELECT COUNT(id) FROM archive_first_pass WHERE aid IN (%s)"
	_inFirstPass      = "INSERT INTO `archive_first_pass`(`aid`, `ctime`, `mtime`) VALUES(?,?,?);"
)

//GetFirstPassByAID 根据aid获取第一次过审的记录
func (d *Dao) GetFirstPassByAID(c context.Context, aid int64) (id int64, err error) {
	row := d.db.QueryRow(c, _slFirstPassByAID, aid)
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

// FirstPassCount 根据aid获取第一次过审的数量
func (d *Dao) FirstPassCount(c context.Context, aids []int64) (count int, err error) {
	if len(aids) < 1 {
		return
	}
	row := d.rdb.QueryRow(c, fmt.Sprintf(_firstPassCount, xstr.JoinInts(aids)))
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan() error(%v) aid(%v)", err, aids)
		}
		return
	}
	return
}

//AddFirstPass 添加一条 第一次过审的记录
func (d *Dao) AddFirstPass(tx *sql.Tx, aid int64) (err error) {
	now := time.Now()
	if _, err = tx.Exec(_inFirstPass, aid, now, now); err != nil {
		log.Error("AddFirstPass error(%v) aid(%d)", aid, err)
	}

	return
}
