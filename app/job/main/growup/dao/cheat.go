package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/library/log"
	"go-common/library/xstr"

	"go-common/app/job/main/growup/model"
)

const (
	// select
	_upsSQL       = "SELECT mid,nickname,fans,signed_at,is_deleted FROM up_info_video WHERE account_state = 3 AND mid IN (%s)"
	_avIncomeSQL  = "SELECT mid,av_id,upload_time,total_income FROM av_income_statis WHERE av_id IN (%s) AND mtime > ?"
	_playCountSQL = "SELECT mid,play FROM up_base_statistics WHERE mid IN (%s)"
	_avBreachSQL  = "SELECT id,av_id FROM av_breach_record WHERE id > ? ORDER BY id LIMIT ?"

	_deleteArchiveSpySQL = "DELETE FROM archive_spy_statistics LIMIT ?"
	_deleteUpSpySQL      = "DELETE FROM up_spy_statistics LIMIT ?"

	// insert
	_inCheatUpsSQL     = "INSERT INTO up_spy_statistics(mid,signed_at,nickname,fans,cheat_fans,play_count,cheat_play_count,account_state) VALUES %s ON DUPLICATE KEY UPDATE signed_at=VALUES(signed_at),nickname=VALUES(nickname),fans=VALUES(fans),cheat_fans=VALUES(cheat_fans),play_count=VALUES(play_count),cheat_play_count=VALUES(cheat_play_count),account_state=VALUES(account_state)"
	_inCheatArchiveSQL = "INSERT INTO archive_spy_statistics(archive_id,mid,nickname,upload_time,total_income,cheat_play_count,cheat_favorite,cheat_coin,deducted) VALUES %s ON DUPLICATE KEY UPDATE nickname=VALUES(nickname),total_income=VALUES(total_income),cheat_play_count=VALUES(cheat_play_count),cheat_favorite=VALUES(cheat_favorite),cheat_coin=VALUES(cheat_coin),deducted=VALUES(deducted)"
)

// DelArchiveSpy del archive spy
func (d *Dao) DelArchiveSpy(c context.Context, limit int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _deleteArchiveSpySQL, limit)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// DelUpSpy del up spy
func (d *Dao) DelUpSpy(c context.Context, limit int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _deleteUpSpySQL, limit)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// AvBreachRecord get av_ids deducted
func (d *Dao) AvBreachRecord(c context.Context, id int64, limit int64) (last int64, ds map[int64]bool, err error) {
	rows, err := d.db.Query(c, _avBreachSQL, id, limit)
	if err != nil {
		return
	}
	ds = make(map[int64]bool)
	defer rows.Close()
	for rows.Next() {
		var avID int64
		err = rows.Scan(&last, &avID)
		if err != nil {
			return
		}
		ds[avID] = true
	}
	return
}

// Ups get ups in up_info_video
func (d *Dao) Ups(c context.Context, mids []int64) (cs map[int64]*model.Cheating, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_upsSQL, xstr.JoinInts(mids)))
	if err != nil {
		return
	}
	defer rows.Close()
	cs = make(map[int64]*model.Cheating)
	for rows.Next() {
		c := &model.Cheating{}
		err = rows.Scan(&c.MID, &c.Nickname, &c.Fans, &c.SignedAt, &c.IsDeleted)
		if err != nil {
			return
		}
		cs[c.MID] = c
	}
	return
}

// Avs get avs in av_income_statis
func (d *Dao) Avs(c context.Context, date time.Time, aids []int64) (cs map[int64]*model.Cheating, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_avIncomeSQL, xstr.JoinInts(aids)), date)
	if err != nil {
		return
	}
	defer rows.Close()
	cs = make(map[int64]*model.Cheating)
	for rows.Next() {
		c := &model.Cheating{}
		err = rows.Scan(&c.MID, &c.AvID, &c.UploadTime, &c.TotalIncome)
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		cs[c.AvID] = c
	}
	return
}

// PlayCount get play count in up_base_statistics
func (d *Dao) PlayCount(c context.Context, mids []int64) (cs map[int64]int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_playCountSQL, xstr.JoinInts(mids)))
	if err != nil {
		return
	}
	cs = make(map[int64]int64)
	defer rows.Close()
	for rows.Next() {
		var mid, count int64
		err = rows.Scan(&mid, &count)
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		cs[mid] = count
	}
	return
}

// InsertCheatUps insert cheat ups
func (d *Dao) InsertCheatUps(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inCheatUpsSQL, values))
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// InsertCheatArchives insert cheat archives
func (d *Dao) InsertCheatArchives(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inCheatArchiveSQL, values))
	if err != nil {
		return
	}
	return res.RowsAffected()
}
