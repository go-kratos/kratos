package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/growup/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/time"
	"go-common/library/xstr"
)

const (
	// insert
	_inUpsSQL      = "INSERT INTO up_info_video (mid,nickname,account_type,original_archive_count,category_id,fans,account_state,sign_type,reason,is_deleted) VALUES (?,?,?,?,?,?,?,?,?,0) ON DUPLICATE KEY UPDATE nickname=?,account_type=?,original_archive_count=?,category_id=?,fans=?,account_state=?,sign_type=?,reason=?,is_deleted=0"
	_inUpColumnSQL = "INSERT INTO up_info_column (mid,nickname,category_id,fans,account_type,account_state,sign_type,is_deleted) VALUES (?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE nickname=VALUES(nickname),account_type=VALUES(account_type),account_state=VALUES(account_state),category_id=VALUES(category_id),sign_type=VALUES(sign_type),is_deleted=VALUES(is_deleted)"

	_inUpBgmSQL = "INSERT INTO up_info_bgm(mid,nickname,bgms,play_count,apply_count,fans,account_state,account_type,sign_type,is_deleted) VALUES(?,?,?,?,?,?,?,?,?,0) ON DUPLICATE KEY UPDATE nickname=VALUES(nickname),bgms=VALUES(bgms),play_count=VALUES(play_count),apply_count=VALUES(apply_count),fans=VALUES(fans),account_state=VALUES(account_state),account_type=VALUES(account_type),sign_type=VALUES(sign_type),is_deleted=VALUES(is_deleted)"

	_inWhitelistSQL = "INSERT INTO up_white_list(mid,type) VALUES(?,?) ON DUPLICATE KEY UPDATE type=VALUES(type)"

	// select
	_upsCateInfoSQL   = "SELECT nick_name,main_category_id FROM up_category_info WHERE mid=?"
	_upsStatInfoSQL   = "SELECT fans,avs FROM up_base_statistics WHERE mid=?"
	_upsCountSQL      = "SELECT count(*) FROM %s WHERE %s "
	_upsInfoSQL       = "SELECT mid,nickname,account_type,original_archive_count,category_id,fans,account_state,sign_type,reason,apply_at,signed_at,reject_at,forbid_at,quit_at,dismiss_at,expired_in,ctime,mtime,is_deleted,credit_score,total_play_count,avs FROM up_info_video WHERE %s"
	_upsColumnInfoSQL = "SELECT mid,nickname,account_type,article_count,category_id,fans,account_state,sign_type,total_view_count,apply_at,signed_at,reject_at,forbid_at,quit_at,dismiss_at,expired_in FROM up_info_column WHERE %s"
	_upsBgmInfoSQL    = "SELECT mid,nickname,bgms,play_count,apply_count,fans,account_state,signed_at,forbid_at,quit_at,dismiss_at,expired_in FROM up_info_bgm WHERE %s"
	_upInfoSQL        = "SELECT mid,nickname,fans,signed_at FROM up_info_video WHERE mid=? AND account_state=? AND is_deleted=0"
	_upInfoStateSQL   = "SELECT mid FROM %s WHERE account_state = ? AND mid in (%s) AND is_deleted = 0"
	_upStateSQL       = "SELECT account_state FROM %s WHERE mid = ? AND is_deleted = 0 LIMIT 1"

	_pendingsSQL = "SELECT mid FROM %s WHERE mid IN (%s) AND account_state=2 AND is_deleted=0"
	_unusualSQL  = "SELECT mid FROM %s WHERE mid IN (%s) AND account_state IN (5, 6, 7) AND is_deleted = 0"
	_bgmCountSQL = "SELECT count(distinct sid) FROM background_music WHERE mid=?"

	// update
	_rejectUpsSQL       = "UPDATE %s SET account_state=?,reason=?,reject_at=?,expired_in=? WHERE mid IN (%s) AND is_deleted=0"
	_passUpsSQL         = "UPDATE %s SET account_state=?,signed_at=? WHERE mid IN (%s) AND is_deleted = 0"
	_dismissUpSQL       = "UPDATE %s SET account_state=?,reason=?,dismiss_at=?,quit_at=? WHERE mid=? AND account_state=? AND is_deleted=0"
	_forbidUpSQL        = "UPDATE %s SET account_state=?,reason=?,forbid_at=?,expired_in=? WHERE mid=? AND account_state=? AND is_deleted=0"
	_updateAccStateSQL  = "UPDATE %s SET account_state=? WHERE mid = ?"
	_updateUpInfoDelSQL = "UPDATE %s SET is_deleted=? WHERE mid=?"
	_delUpAccountSQL    = "UPDATE up_account SET is_deleted=1 WHERE mid=?"
	_updateUpAccountSQL = "UPDATE up_account SET is_deleted = ?, withdraw_date_version = '%s' WHERE mid = ?"
	_delCreditRecordSQL = "UPDATE credit_score_record SET is_deleted=1 WHERE id=?"
)

// InsertWhitelist insert white mid
func (d *Dao) InsertWhitelist(c context.Context, mid int64, typ int) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _inWhitelistSQL, mid, typ)
	if err != nil {
		log.Error("db.inWhitelist.Exec(%s) error(%v)", _inWhitelistSQL, err)
		return
	}
	return res.RowsAffected()
}

// Pendings get mids for account_state=2
func (d *Dao) Pendings(c context.Context, mids []int64, table string) (ms []int64, err error) {
	rows, err := d.rddb.Query(c, fmt.Sprintf(_pendingsSQL, table, xstr.JoinInts(mids)))
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid int64
		err = rows.Scan(&mid)
		if err != nil {
			return
		}
		ms = append(ms, mid)
	}
	return
}

// UnusualUps get mids for account_state=5,6,7
func (d *Dao) UnusualUps(c context.Context, mids []int64, table string) (ms []int64, err error) {
	ms = make([]int64, 0)
	rows, err := d.rddb.Query(c, fmt.Sprintf(_unusualSQL, table, xstr.JoinInts(mids)))
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid int64
		err = rows.Scan(&mid)
		if err != nil {
			return
		}
		ms = append(ms, mid)
	}
	return
}

// InsertUpVideo add upinfo video
func (d *Dao) InsertUpVideo(c context.Context, v *model.UpInfo) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _inUpsSQL, v.MID, v.Nickname, v.AccountType, v.OriginalArchiveCount, v.MainCategory, v.Fans, v.AccountState, v.SignType, v.Reason, v.Nickname, v.AccountType, v.OriginalArchiveCount, v.MainCategory, v.Fans, v.AccountState, v.SignType, v.Reason)
	if err != nil {
		log.Error("db.inUpsStmt.Exec(%s) error(%v)", _inUpsSQL, err)
		return
	}
	return res.RowsAffected()
}

// InsertUpColumn insert up column
func (d *Dao) InsertUpColumn(c context.Context, up *model.UpInfo) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _inUpColumnSQL, up.MID, up.Nickname, up.MainCategory, up.Fans, up.AccountType, up.AccountState, up.SignType, 0)
	if err != nil {
		log.Error("db.inUpsStmt.Exec(%s) error(%v)", _inUpColumnSQL, err)
		return
	}
	return res.RowsAffected()
}

// InsertBgmUpInfo insert up bgm
func (d *Dao) InsertBgmUpInfo(c context.Context, m *model.UpInfo) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _inUpBgmSQL, m.MID, m.Nickname, m.BGMs, m.BgmPlayCount, m.BgmApplyCount, m.Fans, m.AccountState, m.AccountType, m.SignType)
	if err != nil {
		log.Error("db.inUpsStmt.Exec(%s) error(%v)", _inUpBgmSQL, err)
		return
	}
	return res.RowsAffected()
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

// Stat return fans and avs count
func (d *Dao) Stat(c context.Context, mid int64) (fans int, avs int, err error) {
	row := d.rddb.QueryRow(c, _upsStatInfoSQL, mid)
	if err = row.Scan(&fans, &avs); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// UpsCount get count by query
func (d *Dao) UpsCount(c context.Context, table, query string) (count int, err error) {
	row := d.rddb.QueryRow(c, fmt.Sprintf(_upsCountSQL, table, query))
	err = row.Scan(&count)
	return
}

// UpsVideoInfo get up infos by query
func (d *Dao) UpsVideoInfo(c context.Context, query string) (ups []*model.UpInfo, err error) {
	rows, err := d.rddb.Query(c, fmt.Sprintf(_upsInfoSQL, query))
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := &model.UpInfo{}
		err = rows.Scan(&up.MID, &up.Nickname, &up.AccountType, &up.OriginalArchiveCount, &up.MainCategory, &up.Fans, &up.AccountState, &up.SignType, &up.Reason, &up.ApplyAt, &up.SignedAt, &up.RejectAt, &up.ForbidAt, &up.QuitAt, &up.DismissAt, &up.ExpiredIn, &up.CTime, &up.MTime, &up.IsDeleted, &up.CreditScore, &up.TotalPlayCount, &up.Avs)
		if err != nil {
			return
		}
		ups = append(ups, up)
	}
	return
}

// UpsColumnInfo get up column infos by query
func (d *Dao) UpsColumnInfo(c context.Context, query string) (ups []*model.UpInfo, err error) {
	rows, err := d.rddb.Query(c, fmt.Sprintf(_upsColumnInfoSQL, query))
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := &model.UpInfo{}
		err = rows.Scan(&up.MID, &up.Nickname, &up.AccountType, &up.ArticleCount, &up.MainCategory, &up.Fans, &up.AccountState, &up.SignType, &up.TotalViewCount, &up.ApplyAt, &up.SignedAt, &up.RejectAt, &up.ForbidAt, &up.QuitAt, &up.DismissAt, &up.ExpiredIn)
		if err != nil {
			return
		}
		ups = append(ups, up)
	}
	return
}

// UpsBgmInfo get ups bgm infos by query
func (d *Dao) UpsBgmInfo(c context.Context, query string) (ups []*model.UpInfo, err error) {
	rows, err := d.rddb.Query(c, fmt.Sprintf(_upsBgmInfoSQL, query))
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := &model.UpInfo{}
		err = rows.Scan(&up.MID, &up.Nickname, &up.BGMs, &up.BgmPlayCount, &up.BgmApplyCount, &up.Fans, &up.AccountState, &up.SignedAt, &up.ForbidAt, &up.QuitAt, &up.DismissAt, &up.ExpiredIn)
		if err != nil {
			return
		}
		ups = append(ups, up)
	}
	return
}

// Reject batch update reject
func (d *Dao) Reject(c context.Context, table string, state int, reason string, rejectAt, expiredIn time.Time, mids []int64) (rows int64, err error) {
	res, err := d.rddb.Exec(c, fmt.Sprintf(_rejectUpsSQL, table, xstr.JoinInts(mids)), state, reason, rejectAt, expiredIn)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// Pass pass the apply
func (d *Dao) Pass(c context.Context, table string, state int, signedAt time.Time, mids []int64) (rows int64, err error) {
	res, err := d.rddb.Exec(c, fmt.Sprintf(_passUpsSQL, table, xstr.JoinInts(mids)), state, signedAt)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// Dismiss dismiss up
func (d *Dao) Dismiss(c context.Context, table string, newState, oldState int, reason string, dismissAt, quitAt time.Time, mid int64) (rows int64, err error) {
	res, err := d.rddb.Exec(c, fmt.Sprintf(_dismissUpSQL, table), newState, reason, dismissAt, quitAt, mid, oldState)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// TxDismiss tx dismiss up
func (d *Dao) TxDismiss(tx *sql.Tx, table string, newState, oldState int, reason string, dismissAt, quitAt time.Time, mid int64) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_dismissUpSQL, table), newState, reason, dismissAt, quitAt, mid, oldState)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// Forbid forbid up
func (d *Dao) Forbid(c context.Context, table string, newState, oldState int, reason string, forbidAt, expiredIn time.Time, mid int64) (rows int64, err error) {
	res, err := d.rddb.Exec(c, fmt.Sprintf(_forbidUpSQL, table), newState, reason, forbidAt, expiredIn, mid, oldState)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// TxForbid tx forbid up
func (d *Dao) TxForbid(tx *sql.Tx, table string, newState, oldState int, reason string, forbidAt, expiredIn time.Time, mid int64) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_forbidUpSQL, table), newState, reason, forbidAt, expiredIn, mid, oldState)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// UpdateAccountState update up account state
func (d *Dao) UpdateAccountState(c context.Context, table string, state int, mid int64) (rows int64, err error) {
	res, err := d.rddb.Exec(c, fmt.Sprintf(_updateAccStateSQL, table), state, mid)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// DelUpInfo soft delete up info
func (d *Dao) DelUpInfo(c context.Context, table string, mid int64) (rows int64, err error) {
	return d.updateUpInfoDel(c, table, mid, 1)
}

// RecUpInfo recover up info from soft delete
func (d *Dao) RecUpInfo(c context.Context, table string, mid int64) (rows int64, err error) {
	return d.updateUpInfoDel(c, table, mid, 0)
}

func (d *Dao) updateUpInfoDel(c context.Context, table string, mid int64, isDeleted int) (rows int64, err error) {
	res, err := d.rddb.Exec(c, fmt.Sprintf(_updateUpInfoDelSQL, table), isDeleted, mid)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// DelUpAccount update mid is_deleted = 1 in up_account
func (d *Dao) DelUpAccount(c context.Context, mid int64) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _delUpAccountSQL, mid)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// UpdateUpAccount update up_account
func (d *Dao) UpdateUpAccount(c context.Context, mid int64, isDeleted int, withdrawDate string) (rows int64, err error) {
	res, err := d.rddb.Exec(c, fmt.Sprintf(_updateUpAccountSQL, withdrawDate), isDeleted, mid)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// DelCreditRecord soft del credit record by id
func (d *Dao) DelCreditRecord(c context.Context, id int64) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _delCreditRecordSQL, id)
	if err != nil {
		log.Error("db.delCreditRecordSQL.Exec(%s) error(%v)", _delCreditRecordSQL, err)
		return
	}
	return res.RowsAffected()
}

// TxDelCreditRecord tx soft del credit record by id
func (d *Dao) TxDelCreditRecord(tx *sql.Tx, id int64) (rows int64, err error) {
	res, err := tx.Exec(_delCreditRecordSQL, id)
	if err != nil {
		log.Error("tx.delCreditRecordSQL.Exec(%s) error(%v)", _delCreditRecordSQL, err)
		return
	}
	return res.RowsAffected()
}

// UpInfo get up info by account_state and mid
func (d *Dao) UpInfo(c context.Context, mid, state int64) (info *model.UpInfo, err error) {
	row := d.rddb.QueryRow(c, _upInfoSQL, mid, state)
	info = &model.UpInfo{}
	if err = row.Scan(&info.MID, &info.Nickname, &info.Fans, &info.SignedAt); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// GetUpInfoByState get up_info by state
func (d *Dao) GetUpInfoByState(c context.Context, table string, mids []int64, state int) (info map[int64]struct{}, err error) {
	info = make(map[int64]struct{})
	rows, err := d.rddb.Query(c, fmt.Sprintf(_upInfoStateSQL, table, xstr.JoinInts(mids)), state)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid int64
		err = rows.Scan(&mid)
		if err != nil {
			log.Error("GetUpInfoByState rows.Scan error(%v)", err)
			return
		}
		info[mid] = struct{}{}
	}
	return
}

// GetUpState get up state
func (d *Dao) GetUpState(c context.Context, table string, mid int64) (state int, err error) {
	err = d.rddb.QueryRow(c, fmt.Sprintf(_upStateSQL, table), mid).Scan(&state)
	if err == sql.ErrNoRows {
		state = 0
		err = nil
	}
	return
}

// BGMCount bgm count by mid
func (d *Dao) BGMCount(c context.Context, mid int64) (count int, err error) {
	row := d.rddb.QueryRow(c, _bgmCountSQL, mid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			count = 0
		} else {
			log.Error("db.QueryRow(%s) error(%v)", _bgmCountSQL, err)
		}
	}
	return
}
