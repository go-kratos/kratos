package income

import (
	"context"
	"fmt"

	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// select
	_blackListByMIDSQL  = "SELECT av_id FROM av_black_list WHERE mid = ? AND ctype = ? AND is_delete = 0"
	_blackListByAvIDSQL = "SELECT av_id FROM av_black_list WHERE av_id in (%s) AND ctype = ? AND is_delete = 0"
	// insert
	_inBlackListSQL = "INSERT INTO av_black_list(av_id,mid,ctype,reason,nickname,has_signed,is_delete) VALUES %s ON DUPLICATE KEY UPDATE reason=VALUES(reason),nickname=VALUES(nickname),has_signed=VALUES(has_signed),is_delete=VALUES(is_delete)"
)

// ListAvBlackList list av_blakc_list by av_id
func (d *Dao) ListAvBlackList(c context.Context, avID []int64, ctype int) (avb map[int64]struct{}, err error) {
	avb = make(map[int64]struct{})
	rows, err := d.db.Query(c, fmt.Sprintf(_blackListByAvIDSQL, xstr.JoinInts(avID)), ctype)
	if err != nil {
		log.Error("ListAvBlackList d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			log.Error("ListAvBlackList rows scan error(%v)", err)
			return
		}
		avb[id] = struct{}{}
	}
	err = rows.Err()
	return
}

// GetAvBlackListByMID list av_blakc_list by av_id
func (d *Dao) GetAvBlackListByMID(c context.Context, mid int64, typ int) (avb map[int64]struct{}, err error) {
	avb = make(map[int64]struct{})
	rows, err := d.db.Query(c, _blackListByMIDSQL, mid, typ)
	if err != nil {
		log.Error("GetAvBlackListByMID d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var avID int64
		err = rows.Scan(&avID)
		if err != nil {
			log.Error("GetAvBlackListByMID rows scan error(%v)", err)
			return
		}
		avb[avID] = struct{}{}
	}
	err = rows.Err()
	return
}

// TxInsertAvBlackList insert val into av_black_list
func (d *Dao) TxInsertAvBlackList(tx *sql.Tx, val string) (rows int64, err error) {
	if val == "" {
		return
	}
	res, err := tx.Exec(fmt.Sprintf(_inBlackListSQL, val))
	if err != nil {
		log.Error("TxInsertAvBlackList tx.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}
