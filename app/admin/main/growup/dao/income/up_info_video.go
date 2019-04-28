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
	_upInfoNicknameSQL      = "SELECT mid, nickname FROM up_info_video WHERE mid in (%s)"
	_upInfoNicknameByMIDSQL = "SELECT nickname FROM %s WHERE mid = ? AND account_state = 3 AND is_deleted = 0"
	// update
	_updateUpInfoScoreSQL = "UPDATE %s set credit_score = credit_score + ? WHERE mid = ? AND is_deleted = 0"
)

// GetUpInfoNickname get nickname
func (d *Dao) GetUpInfoNickname(c context.Context, mids []int64) (upInfo map[int64]string, err error) {
	upInfo = make(map[int64]string)
	if len(mids) == 0 {
		return
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_upInfoNicknameSQL, xstr.JoinInts(mids)))
	if err != nil {
		log.Error("GetUpInfoNickname d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid int64
		var nickname string
		err = rows.Scan(&mid, &nickname)
		if err != nil {
			log.Error("GetUpInfoNickname rows scan error(%v)", err)
			return
		}
		upInfo[mid] = nickname
	}
	err = rows.Err()
	return
}

// GetUpInfoNicknameByMID get nickname by mid
func (d *Dao) GetUpInfoNicknameByMID(c context.Context, mid int64, table string) (nickname string, err error) {
	err = d.db.QueryRow(c, fmt.Sprintf(_upInfoNicknameByMIDSQL, table), mid).Scan(&nickname)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return
}

// TxUpdateUpInfoScore update up_info_video credit score
func (d *Dao) TxUpdateUpInfoScore(tx *sql.Tx, table string, score int, mid int64) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_updateUpInfoScoreSQL, table), score, mid)
	if err != nil {
		return
	}
	return res.RowsAffected()
}
