package dao

import (
	"context"
	"fmt"
	"time"

	pb "go-common/app/service/main/coin/api"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_dayOfEachMonth             = 25
	_sharding                   = 50
	_getTotalCoins              = "SELECT aid,type,SUM(multiply) FROM coin_member_%d WHERE mid = ? GROUP BY aid,type ORDER BY cid desc LIMIT 20"
	_getTotalCoinsByMid         = "SELECT IFNULL(SUM(multiply),0) FROM coin_member_%d WHERE mid=? AND aid=? AND type=?"
	_getTotalCoinsByMidAndUpMid = "SELECT IFNULL(coin_count,0) FROM coin_user_count_%02d WHERE mid=? AND up_mid=?"
	_getCoinCount               = "SELECT IFNULL(count,0) FROM coin_count WHERE aid=? AND type=?"
	_updateCoinCount            = "INSERT INTO coin_count(aid, type, `count`, ctime, mtime) VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE count=count+?,mtime=?"
	_insertCoinArchive          = "INSERT LOW_PRIORITY INTO coin_archive_%d (aid,type,mid,timestamp,multiply) VALUES (?,?,?,?,?)"
	_insertCoinMember           = "INSERT LOW_PRIORITY INTO coin_member_%d (aid,type,mid,timestamp,multiply,up_mid) VALUES (?,?,?,?,?,?)"
	_updateCoinMemberCount      = "INSERT INTO coin_user_count_%02d(mid, up_mid, coin_count, ctime) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE coin_count=coin_count+?"
	_getMemberVideoList         = "SELECT aid,ip,timestamp,multiply FROM coin_member_%d WHERE mid = ? AND type=? AND timestamp > ? ORDER BY cid DESC LIMIT ?"
	// coin settle
	_hitSettlePeriod      = "SELECT id FROM coin_settle_period WHERE to_year=? AND to_month=?"
	_updateSettleBD       = "UPDATE coin_settle_%d SET exp_sub=?,description=?,itime=?,mtime=? WHERE aid=? and type=?" // update for bigdata
	_updateArchiveCoinsBD = "UPDATE coin_count SET count=?, mtime=? WHERE aid=? AND type=?"                            // update for bigdata
)

func (dao *Dao) midHit(mid int64) int64 {
	return mid % _sharding
}

func (dao *Dao) aidHit(aid int64) int64 {
	return aid % _sharding
}

func (dao *Dao) hitCoinPeriod(c context.Context, now time.Time) (id int64, err error) {
	year, month, day := now.Year(), now.Month(), now.Day()
	if now.Day() < _dayOfEachMonth {
		err = fmt.Errorf("opreation from bigdata must after the %dth of each month, today is (%d)", _dayOfEachMonth, day)
		log.Error("%v", err)
		return
	}
	row := dao.coin.QueryRow(c, _hitSettlePeriod, year, month)
	if err != nil {
		log.Error("dao.hitCoinPeriodStmt.QueryRow(%d, %d) error(%v)", year, month, err)
		return
	}
	if err = row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	if id == 0 {
		err = fmt.Errorf("zero id at(%s)", now.String())
		log.Error("%v", err)
	}
	return
}

// UpdateCoinSettleBD update table coin_settle_%d.
func (dao *Dao) UpdateCoinSettleBD(c context.Context, aid, tp, expSub int64, describe string, now time.Time) (affect int64, err error) {
	id, err := dao.hitCoinPeriod(c, now)
	if err != nil {
		return
	}
	sqlStr := fmt.Sprintf(_updateSettleBD, id)
	result, err := dao.coin.Exec(c, sqlStr, expSub, describe, now, now, aid, tp)
	if err != nil {
		log.Error("dao.coin.Exec(%s, %d, %s, %v, %v, %d) error(%v)", sqlStr, expSub, describe, now, now, aid, err)
		return
	}
	affect, err = result.LastInsertId()
	return
}

// CoinList return video list of coin added in one month
func (dao *Dao) CoinList(c context.Context, mid, tp, ts, size int64) (rcs []*pb.ModelList, err error) {
	var rows *sql.Rows
	if rows, err = dao.coin.Query(c, fmt.Sprintf(_getMemberVideoList, dao.midHit(mid)), mid, tp, ts, size); err != nil {
		log.Error("dao.getMemberVideoList.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var rc = &pb.ModelList{}
		if err = rows.Scan(&rc.Aid, &rc.IP, &rc.Ts, &rc.Number); err != nil {
			log.Error("row.Scan(),error(%v)", err)
			rcs = nil
			return
		}
		rcs = append(rcs, rc)
	}
	return
}

// CoinsAddedByMid get coin added by mid of aid&tp.
func (dao *Dao) CoinsAddedByMid(c context.Context, mid, aid, tp int64) (added int64, err error) {
	row := dao.coin.QueryRow(c, fmt.Sprintf(_getTotalCoinsByMid, dao.midHit(mid)), mid, aid, tp)
	if err = row.Scan(&added); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			PromError("mysql:CoinsAddedByMid")
			log.Errorv(c,
				log.KV("log", "row.Scan"),
				log.KV("err", err),
				log.KV("mid", mid),
				log.KV("aid", aid),
			)
		}
	}
	return
}

// AddedCoins get coins added to up_mid.
func (dao *Dao) AddedCoins(c context.Context, mid, upMid int64) (added int64, err error) {
	row := dao.coin.QueryRow(c, fmt.Sprintf(_getTotalCoinsByMidAndUpMid, dao.midHit(mid)), mid, upMid)
	if err = row.Scan(&added); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// UserCoinsAdded get user coin added.
func (dao *Dao) UserCoinsAdded(c context.Context, mid int64) (addeds map[int64]int64, err error) {
	var rows *sql.Rows
	addeds = make(map[int64]int64)
	if rows, err = dao.coin.Query(c, fmt.Sprintf(_getTotalCoins, dao.midHit(mid)), mid); err != nil {
		log.Errorv(c,
			log.KV("log", "dao.getTotalCoins"),
			log.KV("err", err),
			log.KV("mid", mid),
		)
		PromError("db:UserCoinsAdded")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			added, aid, tp int64
		)
		if err = rows.Scan(&aid, &tp, &added); err != nil {
			log.Errorv(c,
				log.KV("log", "row.Scan"),
				log.KV("err", err),
				log.KV("mid", mid),
			)
			PromError("db:UserCoinsAdded")
			addeds = nil
			return
		}
		addeds[hashField(aid, tp)] = added
	}
	return
}

// InsertCoinArchive .
func (dao *Dao) InsertCoinArchive(c context.Context, aid, tp, mid, timestamp, multiply int64) (err error) {
	if _, err = dao.coin.Exec(c, fmt.Sprintf(_insertCoinArchive, dao.aidHit(aid)), aid, tp, mid, timestamp, multiply); err != nil {
		log.Errorv(c,
			log.KV("log", "coin.Exec"),
			log.KV("err", err),
			log.KV("mid", mid),
			log.KV("aid", aid),
		)
		PromError("db:InsertCoinArchive")
	}
	return
}

// InsertCoinMember .
func (dao *Dao) InsertCoinMember(c context.Context, aid, tp, mid, timestamp, multiply int64, upMid int64) (err error) {
	if _, err = dao.coin.Exec(c, fmt.Sprintf(_insertCoinMember, dao.midHit(mid)), aid, tp, mid, timestamp, multiply, upMid); err != nil {
		log.Errorv(c,
			log.KV("log", "coin.Exec"),
			log.KV("err", err),
			log.KV("mid", mid),
			log.KV("aid", aid),
		)
		PromError("db:InsertCoinMember")
	}
	return
}

// UpdateItemCoinCount update coin_count.
func (dao *Dao) UpdateItemCoinCount(c context.Context, aid, tp, count int64) (err error) {
	now := time.Now()
	if _, err = dao.coin.Exec(c, _updateCoinCount, aid, tp, count, now, now, count, now); err != nil {
		log.Errorv(c,
			log.KV("log", "UpdateItemCoinCount("),
			log.KV("err", err),
			log.KV("count", count),
			log.KV("aid", aid),
		)
		PromError("db:UpdateItemCoinCount")
	}
	return
}

// RawItemCoin get count by aid.
func (dao *Dao) RawItemCoin(c context.Context, aid, tp int64) (count int64, err error) {
	row := dao.coin.QueryRow(c, _getCoinCount, aid, tp)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Errorv(c,
				log.KV("log", "row.Scan"),
				log.KV("err", err),
			)
			PromError("db:RawItemCoin")
		}
	}
	return
}

// UpdateCoinMemberCount for archive
func (dao *Dao) UpdateCoinMemberCount(c context.Context, mid, upMid, count int64) (err error) {
	if _, err = dao.coin.Exec(c, fmt.Sprintf(_updateCoinMemberCount, dao.midHit(mid)), mid, upMid, count, time.Now(), count); err != nil {
		log.Errorv(c,
			log.KV("log", "updateCoinMemberCount"),
			log.KV("mid", mid),
			log.KV("upmid", upMid),
			log.KV("count", count),
			log.KV("err", err),
		)
		PromError("db:UpdateCoinMemberCount")
	}
	return
}

// BeginTran begin a tx.
func (dao *Dao) BeginTran(c context.Context) (t *sql.Tx, err error) {
	t, err = dao.coin.Begin(c)
	if err != nil {
		PromError("db:BeginTran")
		log.Errorv(c,
			log.KV("log", "d.BeginTran"),
			log.KV("err", err),
		)
	}
	return
}

// UpdateItemCoins update table coin_archive
func (dao *Dao) UpdateItemCoins(c context.Context, aid, tp, coins int64, now time.Time) (affect int64, err error) {
	if err != nil {
		return
	}
	result, err := dao.coin.Exec(c, _updateArchiveCoinsBD, coins, now, aid, tp)
	if err != nil {
		log.Error("dao.coin.Exec(%s, %d, %s, %d) error(%v)", _updateArchiveCoinsBD, coins, now, aid, err)
		return
	}
	affect, err = result.RowsAffected()
	return
}
