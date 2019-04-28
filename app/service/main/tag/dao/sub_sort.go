package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/tag/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

var (
	_insertSubSortSQL = "INSERT IGNORE INTO subscriber_sort_%s (mid,`type`,sort,tids) VALUE (?,?,0,?) ON DUPLICATE KEY UPDATE tids=?;"
)

// AddSubSort .
func (d *Dao) AddSubSort(c context.Context, mid int64, typ int, tids []int64) (affected int64, err error) {
	tidsB := model.Index(tids)
	res, err := d.db.Exec(c, fmt.Sprintf(_insertSubSortSQL, d.hit(mid)), mid, typ, tidsB, tidsB)
	if err != nil {
		log.Error("tx.Exec(%d,%d,%v) error(%v)", mid, typ, tids, err)
		return
	}
	return res.RowsAffected()
}

var (
	_insertSubChannelSQL = "INSERT IGNORE INTO subscriber_sort_%s (mid,`type`,sort,tids) VALUE (?,?,?,?) ON DUPLICATE KEY UPDATE tids=?;"
)

// AddSubChannel .
func (d *Dao) AddSubChannel(c context.Context, mid int64, tp int, tids []int64) (err error) {
	var tx *sql.Tx
	tx, err = d.db.Begin(c)
	if err != nil {
		log.Error("open channel sub tx error(%v)", err)
		return
	}
	var n = 100
	for i := 0; len(tids) > 0; i++ {
		if n > len(tids) {
			n = len(tids)
		}
		tidsB := model.Index(tids[:n])
		_, err = tx.Exec(fmt.Sprintf(_insertSubChannelSQL, d.hit(mid)), mid, tp, i, tidsB, tidsB)
		if err != nil {
			log.Error("add sub channel tx exec(%d,%d,%d,%v) error(%v)", mid, tp, i, tidsB, err)
			tx.Rollback()
			return
		}
		tids = tids[n:]
	}
	return tx.Commit()
}

// AddSubChannels .
func (d *Dao) AddSubChannels(c context.Context, mid int64, sortMap map[int32][]int64) (err error) {
	var tx *sql.Tx
	tx, err = d.db.Begin(c)
	if err != nil {
		log.Error("d.dao.AddSubChannels(%d) BeginTx error(%v)", mid, err)
		return
	}
	for tp, tids := range sortMap {
		var n = 100
		for i := 0; len(tids) > 0; i++ {
			if n > len(tids) {
				n = len(tids)
			}
			tidsB := model.Index(tids[:n])
			_, err = tx.Exec(fmt.Sprintf(_insertSubChannelSQL, d.hit(mid)), mid, tp, i, tidsB, tidsB)
			if err != nil {
				log.Error("d.dao.AddSubChannels Exec(%d,%d,%d,%v) error(%v)", mid, tp, i, tidsB, err)
				tx.Rollback()
				return
			}
			tids = tids[n:]
		}
	}
	return tx.Commit()
}

var (
	_selectSubSortSQL = "SELECT tids FROM platform_tag.subscriber_sort_%s where mid = ? AND `type` = ? ORDER BY sort ASC;"
)

// CustomSubSort .
// 每一行存100个tag，当发现该行不足100时表明该type数据已经读取完成，停止该type的数据
func (d *Dao) CustomSubSort(c context.Context, mid int64, typ int) (tids []int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_selectSubSortSQL, d.hit(mid)), mid, typ)
	if err != nil {
		log.Error("d.db.Query(%d,%d) error(%v)", mid, typ, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			b   []byte
			ids []int64
		)
		if err = rows.Scan(&b); err != nil {
			log.Error("custom sub sort rows scan(%d,%d,%v) error(%v)", mid, typ, b, err)
			return
		}
		if ids, err = model.SetIndex(b); err != nil {
			log.Error("model.SetIndex(%d,%d,%v) error(%v)", mid, typ, ids, err)
			return
		}
		if len(ids) < 100 {
			tids = append(tids, ids...)
			break
		}
		tids = append(tids, ids...)
	}
	err = rows.Err()
	return
}

var (
	_allSubSortSQL = "SELECT `type`,tids FROM platform_tag.subscriber_sort_%s where mid = ? ORDER BY sort ASC"
)

// AllCustomSubSort .
// 取数据逻辑描述：
// 每一行存100个tag，当发现该行不足100时表明该type数据已经读取完成，停止该type的数据
// tpMap 判断该type是否继续读取逻辑.
func (d *Dao) AllCustomSubSort(c context.Context, mid int64) (tidMap map[int32][]int64, err error) {
	tidMap = make(map[int32][]int64)
	rows, err := d.db.Query(c, fmt.Sprintf(_allSubSortSQL, d.hit(mid)), mid)
	if err != nil {
		log.Error("d.AllCustomSubSort(%d) error(%v)", mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			tp    int32
			b     []byte
			ids   []int64
			tpMap = make(map[int32]bool)
		)
		if err = rows.Scan(&tp, &b); err != nil {
			log.Error("d.AllCustomSubSort(%d) rows.scan() error(%v)", mid, err)
			return
		}
		if ids, err = model.SetIndex(b); err != nil {
			log.Error("d.AllCustomSubSort(%d)SetIndex(%v) error(%v)", mid, ids, err)
			return
		}
		k, ok := tpMap[tp]
		if !ok {
			tpMap[tp] = false
		}
		if k {
			continue
		}
		tidMap[tp] = append(tidMap[tp], ids...)
		if len(ids) < 100 {
			tpMap[tp] = true
		}
	}
	err = rows.Err()
	return
}
