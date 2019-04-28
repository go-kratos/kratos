package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/main/growup/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/time"
)

const (
	// insert
	_inUpsSQL = "INSERT INTO %s (mid,nickname,account_type,category_id,fans,account_state,sign_type,apply_at,%s,is_deleted) VALUES (?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE nickname=?,account_type=?,category_id=?,fans=?,account_state=?,sign_type=?,apply_at=?,%s=?,is_deleted=?"

	_inUpBgmSQL        = "INSERT INTO up_info_bgm(mid,nickname,bgms,account_type,fans,account_state,sign_type,signed_at,is_deleted) VALUES(?,?,?,?,?,3,?,?,0) ON DUPLICATE KEY UPDATE nickname=VALUES(nickname),bgms=VALUES(bgms),account_type=VALUES(account_type),fans=VALUES(fans),account_state=3,signed_at=VALUES(signed_at),sign_type=VALUES(sign_type),is_deleted=0"
	_inCreditRecordSQL = "INSERT INTO credit_score_record (mid,operate_at,operator,reason,deducted,remaining) VALUES(?,?,?,?,?,?)"
	_inCreditScoreSQL  = "INSERT INTO credit_score(mid) VALUES(?) ON DUPLICATE KEY UPDATE mid=VALUES(mid)"

	// select
	_upsCateInfoSQL   = "SELECT nick_name,main_category_id FROM up_category_info WHERE mid=?"
	_upsStatInfoSQL   = "SELECT fans FROM up_base_statistics WHERE mid=?"
	_blockMIDSQL      = "SELECT mid FROM up_blocked WHERE mid=? AND is_deleted=0"
	_upNicknameSQL    = "SELECT nickname FROM up_info_video WHERE mid = ?"
	_upCreditScoreSQL = "SELECT score FROM credit_score WHERE mid = ?"
	_accountStateSQL  = "SELECT account_state FROM %s WHERE mid = ? AND is_deleted = 0"
	_upSignedAtSQL    = "SELECT signed_at FROM %s WHERE mid = ? AND account_state = 3 AND is_deleted = 0"

	// will delete next version
	_upWhiteListSQL = "SELECT type FROM up_white_list WHERE mid=? AND is_deleted=0"

	//	_upsArchiveSQL = "SELECT account_type,account_state,reason,expired_in,quit_at,ctime FROM up_info_video WHERE mid=? AND is_deleted=0"
	_avAccountStateSQL     = "SELECT account_type,account_state,reason,expired_in,quit_at,ctime FROM up_info_video WHERE mid=? AND is_deleted=0"
	_bgmAccountStateSQL    = "SELECT account_type,account_state,reason,expired_in,quit_at,ctime FROM up_info_bgm WHERE mid=? AND is_deleted=0"
	_columnAccountStateSQL = "SELECT account_type,account_state,reason,expired_in,quit_at,ctime FROM up_info_column WHERE mid=? AND is_deleted=0"
	_bgmUpCountSQL         = "SELECT count(*) FROM background_music WHERE mid=?"
	// update
	_upQuitSQL            = "UPDATE %s SET account_state=5,quit_at=?,expired_in=?,reason=? WHERE mid=? AND account_state=3 AND is_deleted=0"
	_deductCreditScoreSQL = "UPDATE credit_score SET score=score-%d WHERE mid=?"
)

// GetAccountState account state
func (d *Dao) GetAccountState(c context.Context, table string, mid int64) (state int, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_accountStateSQL, table), mid)
	if err = row.Scan(&state); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row scan error(%v)", err)
		}
	}
	return
}

// GetUpSignedAt get up signed_at
func (d *Dao) GetUpSignedAt(c context.Context, table string, mid int64) (signedAt time.Time, err error) {
	if err = d.db.QueryRow(c, fmt.Sprintf(_upSignedAtSQL, table), mid).Scan(&signedAt); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			signedAt = 0
		} else {
			log.Error("row scan error(%v)", err)
		}
	}
	return
}

// InsertUpInfo add upinfo
func (d *Dao) InsertUpInfo(c context.Context, table string, totalCountField string, v *model.UpInfo) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inUpsSQL, table, totalCountField, totalCountField), v.MID, v.Nickname, v.AccountType, v.MainCategory, v.Fans, v.AccountState, v.SignType, v.ApplyAt, v.TotalPlayCount, 0, v.Nickname, v.AccountType, v.MainCategory, v.Fans, v.AccountState, v.SignType, v.ApplyAt, v.TotalPlayCount, 0)
	if err != nil {
		log.Error("db.inUpsStmt.Exec(%s) error(%v)", _inUpsSQL, err)
		return
	}
	return res.RowsAffected()
}

// TxInsertBgmUpInfo insert bgm up info
func (d *Dao) TxInsertBgmUpInfo(tx *sql.Tx, v *model.UpInfo) (rows int64, err error) {
	res, err := tx.Exec(_inUpBgmSQL, v.MID, v.Nickname, v.Bgms, v.AccountType, v.Fans, v.SignType, v.SignedAt)
	if err != nil {
		log.Error("db.inBgmUpStmt.Exec(%s) error(%v)", _inUpsSQL, err)
		return
	}
	return res.RowsAffected()
}

// TxInsertCreditScore insert credit score
func (d *Dao) TxInsertCreditScore(tx *sql.Tx, mid int64) (rows int64, err error) {
	res, err := tx.Exec(_inCreditScoreSQL, mid)
	if err != nil {
		log.Error("db.Exec(%s) error(%v)", _inCreditScoreSQL, err)
		return
	}
	return res.RowsAffected()
}

// Blocked query mid in blacklist
func (d *Dao) Blocked(c context.Context, mid int64) (id int64, err error) {
	row := d.db.QueryRow(c, _blockMIDSQL, mid)
	if err = row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// White map: key: type value: bool
func (d *Dao) White(c context.Context, mid int64) (m map[int]bool, err error) {
	m = make(map[int]bool)
	rows, err := d.db.Query(c, _upWhiteListSQL, mid)
	if err != nil {
		log.Error("row.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var typ int
		err = rows.Scan(&typ)
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
		}
		m[typ] = true
	}
	return
}

// AvUpStatus return av up status
func (d *Dao) AvUpStatus(c context.Context, mid int64) (status *model.BusinessStatus, err error) {
	status = &model.BusinessStatus{}
	row := d.db.QueryRow(c, _avAccountStateSQL, mid)
	if err = row.Scan(&status.AccountType, &status.AccountState, &status.Reason, &status.ExpiredIn, &status.QuitAt, &status.CTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// BgmUpStatus return bgm up status
func (d *Dao) BgmUpStatus(c context.Context, mid int64) (status *model.BusinessStatus, err error) {
	status = &model.BusinessStatus{}
	row := d.db.QueryRow(c, _bgmAccountStateSQL, mid)
	if err = row.Scan(&status.AccountType, &status.AccountState, &status.Reason, &status.ExpiredIn, &status.QuitAt, &status.CTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// ColumnUpStatus return bgm up status
func (d *Dao) ColumnUpStatus(c context.Context, mid int64) (status *model.BusinessStatus, err error) {
	status = &model.BusinessStatus{}
	row := d.db.QueryRow(c, _columnAccountStateSQL, mid)
	if err = row.Scan(&status.AccountType, &status.AccountState, &status.Reason, &status.ExpiredIn, &status.QuitAt, &status.CTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// CategoryInfo return nickname & categoryID
func (d *Dao) CategoryInfo(c context.Context, mid int64) (nickname string, categoryID int, err error) {
	row := d.rddb.QueryRow(c, _upsCateInfoSQL, mid)
	if err = row.Scan(&nickname, &categoryID); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// Fans return fans count
func (d *Dao) Fans(c context.Context, mid int64) (fans int, err error) {
	row := d.rddb.QueryRow(c, _upsStatInfoSQL, mid)
	if err = row.Scan(&fans); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// TxQuit update up status
func (d *Dao) TxQuit(tx *sql.Tx, table string, mid int64, quitAt, expiredIn time.Time, reason string) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_upQuitSQL, table), quitAt, expiredIn, reason, mid)
	if err != nil {
		log.Error("db.TxQuit.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxInsertCreditRecord tx insert credit deduct record
func (d *Dao) TxInsertCreditRecord(tx *sql.Tx, cr *model.CreditRecord) (rows int64, err error) {
	res, err := tx.Exec(_inCreditRecordSQL, cr.MID, cr.OperateAt, cr.Operator, cr.Reason, cr.Deducted, cr.Remaining)
	if err != nil {
		log.Error("db.TxInsertCreditRecord.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// Nickname get nickname from up_info_video
func (d *Dao) Nickname(c context.Context, mid int64) (nickname string, err error) {
	row := d.rddb.QueryRow(c, _upNicknameSQL, mid)
	if err = row.Scan(&nickname); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// CreditScore get current credit score from up_info_video
func (d *Dao) CreditScore(c context.Context, mid int64) (score int, err error) {
	row := d.rddb.QueryRow(c, _upCreditScoreSQL, mid)
	if err = row.Scan(&score); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// BgmUpCount bgm up count in table background_music
func (d *Dao) BgmUpCount(c context.Context, mid int64) (count int, err error) {
	row := d.db.QueryRow(c, _bgmUpCountSQL, mid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			count = 0
		} else {
			log.Error("db.QueryRow(%s) error(%v)", _bgmUpCountSQL, err)
		}
	}
	return
}

// TxDeductCreditScore tx update credit score
func (d *Dao) TxDeductCreditScore(tx *sql.Tx, score int, mid int64) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_deductCreditScoreSQL, score), mid)
	if err != nil {
		log.Error("tx.Exec(%s) error(%v)", _deductCreditScoreSQL, err)
		return
	}
	return res.RowsAffected()
}
