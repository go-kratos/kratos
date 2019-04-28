package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"go-common/app/job/main/coin/model"
	"go-common/library/log"
)

const (
	// SHARDING table shard.
	SHARDING = 50
	// coin_settle_x
	_getSettle       = "SELECT id,mid,aid,type,coin_count,exp_sub,state FROM coin_settle_%d WHERE id> ? LIMIT 10000"
	_hitSettlePeriod = "SELECT id,from_year,from_month,from_day,to_year,to_month,to_day FROM coin_settle_period WHERE from_year*10000+from_month*100+from_day<=? AND to_year*10000+to_month*100+to_day>?"
	_getSettlePeriod = "SELECT id,from_year,from_month,from_day,to_year,to_month,to_day FROM coin_settle_period WHERE id=?"
	_clearCoinCount  = "UPDATE coin_settle_%d SET coin_count=0, mtime=?"
	_updateCoinCount = "UPDATE coin_settle_%d SET coin_count=?, mtime=? WHERE aid=? AND type=?"
	_updateSettle    = "UPDATE coin_settle_%d SET state=1, exp_total=?, mtime=? WHERE id=?"
	_upsertSettle    = "INSERT INTO coin_settle_%d (mid,aid,type,coin_count,ctime,mtime) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE coin_count=coin_count+?,state=0, mtime=?"
	// coin_archive_x
	_getTotalCoins = "SELECT aid,type,SUM(multiply) FROM coin_archive_%d WHERE timestamp>=? AND timestamp<? GROUP BY aid,type"
)

func hashField(aid, tp int64) int64 {
	return aid*1000 + tp
}

// SettlePeriod settle coin by tableid.
func (dao *Dao) SettlePeriod(c context.Context, id int64) (period *model.CoinSettlePeriod, err error) {
	row := dao.getSettlePeriodStmt.QueryRow(c, id)
	period = &model.CoinSettlePeriod{}
	if err = row.Scan(&period.ID,
		&period.FromYear,
		&period.FromMonth,
		&period.FromDay,
		&period.ToYear,
		&period.ToMonth,
		&period.ToDay); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// HitSettlePeriod get table id.
func (dao *Dao) HitSettlePeriod(c context.Context, now time.Time) (period *model.CoinSettlePeriod, err error) {
	ymd, _ := strconv.Atoi(now.Format("20060102"))
	row := dao.hitSettlePeriodStmt.QueryRow(c, ymd, ymd)
	period = &model.CoinSettlePeriod{}
	if err = row.Scan(&period.ID,
		&period.FromYear,
		&period.FromMonth,
		&period.FromDay,
		&period.ToYear,
		&period.ToMonth,
		&period.ToDay); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	if period.ID == 0 {
		err = fmt.Errorf("zero id at(%d)", ymd)
	}
	return
}

// UpdateSettle update settle info.
func (dao *Dao) UpdateSettle(c context.Context, tableID, id, expTotal int64, now time.Time) (err error) {
	sqlStr := fmt.Sprintf(_updateSettle, tableID)
	if _, err = dao.coinDB.Exec(c, sqlStr, expTotal, now, id); err != nil {
		log.Error("dao.coinDB.Exec(%s, %d, %v, %d) error(%v)", sqlStr, expTotal, now, id, err)
		PromError("db:UpdateSettle")
	}
	return
}

// UpdateCoinCount update coin.
func (dao *Dao) UpdateCoinCount(c context.Context, tableID, aid, tp, count int64, now time.Time) (err error) {
	sqlStr := fmt.Sprintf(_updateCoinCount, tableID)
	if _, err = dao.coinDB.Exec(c, sqlStr, count, now, aid, tp); err != nil {
		log.Error("dao.coinDB.Exec(%s, %d, %v, %d) error(%v)", sqlStr, count, now, aid, err)
	}
	return
}

// Every10000 get 10000 coin record.
func (dao *Dao) Every10000(c context.Context, tableID int64, idx int64) (settles []*model.CoinSettle, maxid int64, err error) {
	sqlStr := fmt.Sprintf(_getSettle, tableID)
	rows, err := dao.coinDB.Query(c, sqlStr, idx)
	if err != nil {
		log.Error("dao.coinDB.Query(%s, %d) error(%v)", sqlStr, idx, err)
		return
	}
	defer rows.Close()
	maxid = idx
	settles = make([]*model.CoinSettle, 0, 10000)
	for rows.Next() {
		settle := &model.CoinSettle{}
		if err = rows.Scan(&settle.ID,
			&settle.Mid,
			&settle.Aid,
			&settle.AvType,
			&settle.CoinCount,
			&settle.ExpSub,
			&settle.State); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if settle.ID > maxid {
			maxid = settle.ID
		}
		settles = append(settles, settle)
	}
	return
}

// UpsertSettle update coin settle.
func (dao *Dao) UpsertSettle(c context.Context, tableID, mid, aid, tp, coinCount int64, now time.Time) (err error) {
	sqlStr := fmt.Sprintf(_upsertSettle, tableID)
	if _, err = dao.coinDB.Exec(c, sqlStr, mid, aid, tp, coinCount, now, now, coinCount, now); err != nil {
		log.Error("dao.coinDB.Exec(%s,%d,%d,%d,%v,%v,%d,%v) error(%v)", sqlStr, mid, aid, coinCount, now, now, coinCount, now, err)
	}
	return
}

// ClearCoinCount clear settle.
func (dao *Dao) ClearCoinCount(c context.Context, tableID int64, now time.Time) (err error) {
	sqlStr := fmt.Sprintf(_clearCoinCount, tableID)
	if _, err = dao.coinDB.Exec(c, sqlStr, now); err != nil {
		log.Error("dao.coinDB.Exec(%s, %v) error(%v)", sqlStr, now, err)
	}
	return
}

// TotalCoins get total coins.
func (dao *Dao) TotalCoins(c context.Context, id int, start, end time.Time) (coins map[int64]int64, err error) {
	rows, err := dao.getTotalCoinsStmt[id].Query(c, start.Unix(), end.Unix())
	if err != nil {
		log.Error("dao.getTotalCoinsStmt[%d].Query(%v, %v) error(%v)", id, start, end, err)
		return
	}
	coins = make(map[int64]int64)
	defer rows.Close()
	for rows.Next() {
		var aid, tp, count int64
		if err = rows.Scan(&aid, &tp, &count); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		coins[hashField(aid, tp)] = count
	}
	return
}
