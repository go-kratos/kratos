package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"

	"go-common/app/admin/main/growup/model"
)

const (
	// select
	_upSpySQL = "SELECT mid,signed_at,nickname,fans,cheat_fans,play_count,cheat_play_count,account_state FROM up_spy_statistics %s LIMIT ?, ?"
	_avSpySQL = "SELECT archive_id,mid,nickname,upload_time,total_income,cheat_play_count,cheat_favorite,cheat_coin,deducted FROM archive_spy_statistics %s LIMIT ?,?"
	// select count(*)
	_upSpyCountSQL = "SELECT count(*) FROM up_spy_statistics"
	_avSpyCountSQL = "SELECT count(*) FROM archive_spy_statistics %s"

	// update
	_updateUpState = "UPDATE up_spy_statistics SET account_state=? WHERE mid=?"
	_updateAvState = "UPDATE archive_spy_statistics SET deducted=? WHERE archive_id IN (%s)"

	// cheat fans
	_cheatFansSQL      = "SELECT mid,nickname,real_fans,cheat_fans,signed_at,deduct_at FROM cheat_fans_info WHERE is_deleted = 0 LIMIT ?,?"
	_cheatFansCountSQL = "SELECT count(*) FROM cheat_fans_info WHERE is_deleted = 0"
	_delCheatUpSQL     = "UPDATE cheat_fans_info SET is_deleted = 1 WHERE mid = ?"

	_insertCheatFansSQL = "INSERT INTO cheat_fans_info(mid, nickname, signed_at, real_fans, cheat_fans, deduct_at) VALUES %s ON DUPLICATE KEY UPDATE mid = values(mid), nickname = values(nickname), signed_at = values(signed_at), real_fans = values(real_fans), cheat_fans = values(cheat_fans), deduct_at = values(deduct_at), is_deleted = 0"
)

// TxUpdateUpSpyState tx update up_spy_state
func (d *Dao) TxUpdateUpSpyState(tx *sql.Tx, state int, mid int64) (rows int64, err error) {
	res, err := tx.Exec(_updateUpState, state, mid)
	if err != nil {
		log.Error("d.db.TxUpdateUpSpyState error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateAvSpyState tx update av_spy_state
func (d *Dao) TxUpdateAvSpyState(tx *sql.Tx, state int, archives []int64) (rows int64, err error) {
	if len(archives) == 0 {
		return
	}
	res, err := tx.Exec(fmt.Sprintf(_updateAvState, xstr.JoinInts(archives)), state)
	if err != nil {
		log.Error("d.db.TxUpdateAvSpyState error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// UpSpyCount get up spy count
func (d *Dao) UpSpyCount(c context.Context) (count int, err error) {
	row := d.rddb.QueryRow(c, _upSpyCountSQL)
	if err = row.Scan(&count); err != nil {
		log.Error("d.rddb.UpSpyCount error(%v)", err)
	}
	return
}

// UpSpies get up spy.
func (d *Dao) UpSpies(c context.Context, query string, offset, limit int) (spies []*model.UpSpy, err error) {
	rows, err := d.rddb.Query(c, fmt.Sprintf(_upSpySQL, query), offset, limit)
	if err != nil {
		log.Error("d.db.Query UpSpies error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		spy := &model.UpSpy{}
		if err = rows.Scan(&spy.MID, &spy.SignedAt, &spy.Nickname, &spy.Fans, &spy.CheatFans, &spy.PlayCount, &spy.CheatPlayCount, &spy.AccountState); err != nil {
			log.Error("dao.UpSpies scan error(%v)", err)
			return
		}
		spies = append(spies, spy)
	}
	return
}

// ArchiveSpyCount get archive count
func (d *Dao) ArchiveSpyCount(c context.Context, query string) (count int, err error) {
	row := d.rddb.QueryRow(c, fmt.Sprintf(_avSpyCountSQL, query))
	if err = row.Scan(&count); err != nil {
		log.Error("dao.GetArchiveSpy count error(%v)", err)
	}
	return
}

// ArchiveSpies get av spy.
func (d *Dao) ArchiveSpies(c context.Context, query string, offset, limit int) (spies []*model.ArchiveSpy, err error) {
	rows, err := d.rddb.Query(c, fmt.Sprintf(_avSpySQL, query), offset, limit)
	if err != nil {
		log.Error("dao.GetArchiveSpy query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		spy := &model.ArchiveSpy{}
		if err = rows.Scan(&spy.ArchiveID, &spy.MID, &spy.Nickname, &spy.UploadTime, &spy.TotalIncome, &spy.CheatPlayCount, &spy.CheatFavorite, &spy.CheatCoin, &spy.Deducted); err != nil {
			log.Error("dao.GetArchiveSpy scan error(%v)", err)
			return
		}
		spies = append(spies, spy)
	}
	return
}

// CheatFansCount get cheat fans count.
func (d *Dao) CheatFansCount(c context.Context) (count int64, err error) {
	row := d.rddb.QueryRow(c, _cheatFansCountSQL)
	if err = row.Scan(&count); err != nil {
		log.Error("dao.CheatFansInfo count error(%v)", err)
	}
	return
}

// CheatFans get cheat fans info.
func (d *Dao) CheatFans(c context.Context, from, limit int64) (fans []*model.CheatFans, err error) {
	rows, err := d.rddb.Query(c, _cheatFansSQL, from, limit)
	if err != nil {
		log.Error("dao.CheatFansInfo query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		fan := &model.CheatFans{}
		if err = rows.Scan(&fan.MID, &fan.Nickname, &fan.RealFans, &fan.CheatFans, &fan.SignedAt, &fan.DeductAt); err != nil {
			log.Error("dao.CheatFansInfo scan error(%v)", err)
			return
		}
		fans = append(fans, fan)
	}
	err = rows.Err()
	return
}

// DelCheatUp update cheat_fans_info.
func (d *Dao) DelCheatUp(c context.Context, mid int64) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _delCheatUpSQL, mid)
	if err != nil {
		log.Error("dao.UpdateCheatFans error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// InsertCheatFansInfo insert into cheat_fans_info.
func (d *Dao) InsertCheatFansInfo(c context.Context, values string) (rows int64, err error) {
	res, err := d.rddb.Exec(c, fmt.Sprintf(_insertCheatFansSQL, values))
	if err != nil {
		log.Error("dao.InsertCheatFansInfo error(%v)", err)
		return
	}
	return res.RowsAffected()
}

const (
	_realFansCount  = "/x/internal/relation/stat"
	_cheatFansCount = "/x/internal/v1/spy/stat"
)

// GetUpRealFansCount get up real fans count
func (d *Dao) GetUpRealFansCount(c context.Context, host string, mid int64) (count int, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			FansCount int `json:"follower"`
		} `json:"data"`
	}
	if err = d.client.Get(c, host+_realFansCount, "", params, &res); err != nil {
		log.Error("dao.GetUpRealFansCount get cheat count error(%v)", err)
		return
	}

	if res.Code != 0 {
		err = fmt.Errorf("get real fans count error code: %d", res.Code)
		return
	}
	count = res.Data.FansCount
	return
}

// GetUpCheatFansCount get up cheat fan count.
func (d *Dao) GetUpCheatFansCount(c context.Context, host string, mid int64) (count int, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int                 `json:"code"`
		Data []*model.CheatCount `json:"data"`
	}
	if err = d.client.Get(c, host+_cheatFansCount, "", params, &res); err != nil {
		log.Error("dao.GetUpCheatFansCount get cheat count error(%v)", err)
		return
	}
	if res.Code != 0 {
		log.Error("growup-job GetUpCheatFansCount code != 0. res.Code(%d) | params(%s) error(%v)", res.Code, params.Encode(), err)
		err = fmt.Errorf("get cheat fans count error code: %d", res.Code)
		return
	}
	for _, v := range res.Data {
		if v.EventID == "异常粉丝量" {
			count = v.Quantity
		}
	}
	return
}
