package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/growup/model"

	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

const (
	_upStateByMID   = "SELECT account_state FROM up_info_video WHERE mid = ? LIMIT 1"
	_upInfoVideoSQL = "SELECT id,mid,nickname,signed_at,fans,total_play_count,account_type,account_state,credit_score,is_deleted FROM up_info_video WHERE id > ? ORDER BY id LIMIT ?"

	_upsState           = "SELECT mid,expired_in FROM %s WHERE account_state = ?"
	_upsStateType       = "SELECT mid,expired_in FROM %s WHERE account_type = ? AND account_state = ?"
	_updateAccountState = "UPDATE %s SET account_state = ? WHERE mid IN (%s)"

	// select count(*)
	_signedDayUpsSQL = "SELECT COUNT(*) FROM up_info_video WHERE account_state = 3 AND is_deleted = 0 AND signed_at < ? AND signed_at >= ?"
	_signedAllUpsSQL = "SELECT COUNT(*) FROM up_info_video WHERE account_state = 3 AND signed_at < ? AND is_deleted = 0"

	_videoApplyCountSQL = "SELECT COUNT(*) FROM up_info_video WHERE apply_at >= ? AND apply_at < ?"

	_upBaseInfoSQL       = "SELECT mid,fans,play,avs_origin,avs FROM up_base_statistics WHERE mid IN (%s)"
	_insertUpInfoSQL     = "INSERT INTO up_info_video(mid,fans,total_play_count,original_archive_count,avs) VALUES %s ON DUPLICATE KEY UPDATE fans=VALUES(fans),total_play_count=VALUES(total_play_count),original_archive_count=VALUES(original_archive_count),avs=VALUES(avs)"
	_uidSQL              = "SELECT id,mid FROM up_info_video WHERE id > ? ORDER BY id LIMIT ?"
	_creditScoreByMIDSQL = "SELECT mid, score FROM credit_score WHERE mid IN (%s)"
)

// GetUpStateByMID get up account_state
func (d *Dao) GetUpStateByMID(c context.Context, mid int64) (state int, err error) {
	err = d.db.QueryRow(c, _upStateByMID, mid).Scan(&state)
	if err == sql.ErrNoRows {
		err = nil
		state = 0
	}
	return
}

// GetUpCreditScore get up credit score
func (d *Dao) GetUpCreditScore(c context.Context, mids []int64) (scores map[int64]int64, err error) {
	scores = make(map[int64]int64)
	if len(mids) == 0 {
		return
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_creditScoreByMIDSQL, xstr.JoinInts(mids)))
	if err != nil {
		log.Error("db Query GetUpCreditScore error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid, score int64
		err = rows.Scan(&mid, &score)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		scores[mid] = score
	}
	return
}

// UpInfoVideo get up_info_video
func (d *Dao) UpInfoVideo(c context.Context, offset int64, limit int64) (last int64, ups map[int64]*model.UpInfoVideo, err error) {
	ups = make(map[int64]*model.UpInfoVideo)
	rows, err := d.db.Query(c, _upInfoVideoSQL, offset, limit)
	if err != nil {
		log.Error("db Query UpInfoVideo error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := &model.UpInfoVideo{}
		err = rows.Scan(&last, &up.MID, &up.Nickname, &up.SignedAt, &up.Fans, &up.TotalPlayCount, &up.AccountType, &up.AccountState, &up.CreditScore, &up.IsDeleted)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		ups[up.MID] = up
	}
	return
}

// MIDsByState get mids and expired
func (d *Dao) MIDsByState(c context.Context, state int, table string) (result map[int64]xtime.Time, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_upsState, table), state)
	if err != nil {
		log.Error("d.MIDsByState Query error(%v)", err)
		return
	}

	result = make(map[int64]xtime.Time)
	defer rows.Close()
	for rows.Next() {
		var mid int64
		var exp xtime.Time
		err = rows.Scan(&mid, &exp)
		if err != nil {
			log.Error("rows scan error (%v)", err)
			return
		}
		result[mid] = exp
	}
	return
}

// MIDsByStateType get mids and expired by account_type and account_state
func (d *Dao) MIDsByStateType(c context.Context, typ int, state int, table string) (result map[int64]xtime.Time, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_upsStateType, table), typ, state)
	if err != nil {
		log.Error("d.MIDsByStateType Query error(%v)", err)
		return
	}

	result = make(map[int64]xtime.Time)
	defer rows.Close()
	for rows.Next() {
		var mid int64
		var exp xtime.Time
		err = rows.Scan(&mid, &exp)
		if err != nil {
			log.Error("rows scan error (%v)", err)
			return
		}
		result[mid] = exp
	}
	return
}

// UpdateAccountState update account state
func (d *Dao) UpdateAccountState(c context.Context, state int, mids []int64, table string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_updateAccountState, table, xstr.JoinInts(mids)), state)
	if err != nil {
		log.Error("d.UpdateAccountState Exec error (%v)", err)
		return
	}
	return res.RowsAffected()
}

// GetDateSignedUps get date signed ups
func (d *Dao) GetDateSignedUps(c context.Context, startAt time.Time, endAt time.Time) (count int, err error) {
	row := d.db.QueryRow(c, _signedDayUpsSQL, endAt, startAt)
	if err = row.Scan(&count); err != nil {
		log.Error("dao.GetDateSignedUps error(%v)", err)
	}
	return
}

// GetAllSignedUps get all signed ups.
func (d *Dao) GetAllSignedUps(c context.Context, data time.Time) (count int, err error) {
	row := d.db.QueryRow(c, _signedAllUpsSQL, data.Add(24*time.Hour))
	if err = row.Scan(&count); err != nil {
		log.Error("dao.GetAllSignedUps error(%v)", err)
	}
	return
}

// GetVideoApplyUpCount get up_info_video count
func (d *Dao) GetVideoApplyUpCount(c context.Context, startAt, endAt time.Time) (count int, err error) {
	row := d.db.QueryRow(c, _videoApplyCountSQL, startAt, endAt)
	err = row.Scan(&count)
	return
}

// GetUpBaseInfo get up_base_info
func (d *Dao) GetUpBaseInfo(c context.Context, mid []int64) (bs []*model.UpBaseInfo, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_upBaseInfoSQL, xstr.JoinInts(mid)))
	if err != nil {
		log.Error("dao.GetUpBaseInfo error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.UpBaseInfo{}
		err = rows.Scan(&b.MID, &b.Fans, &b.TotalPlayCount, &b.OriginalArchiveCount, &b.Avs)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		bs = append(bs, b)
	}
	return
}

// UpdateUpInfo update up_info_video
func (d *Dao) UpdateUpInfo(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_insertUpInfoSQL, values))
	if err != nil {
		log.Error("dao.UpdateUpInfo error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// MIDs get mids from up_info_viode
func (d *Dao) MIDs(c context.Context, offset, limit int64) (last int64, mids []int64, err error) {
	rows, err := d.db.Query(c, _uidSQL, offset, limit)
	if err != nil {
		log.Error("dao.GetUpBaseInfo error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid int64
		err = rows.Scan(&last, &mid)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		mids = append(mids, mid)
	}
	return
}
