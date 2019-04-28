package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

const (
	_sharding  = 100
	_statSQL   = "SELECT aid,click,fav,share,reply,coin,dm,now_rank,his_rank,likes,dislike FROM archive_stat_%s WHERE aid=%d"
	_upStatSQL = `INSERT INTO archive_stat_%s (aid,click,fav,share,reply,coin,dm,now_rank,his_rank,ctime,mtime,likes,dislike) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)
				ON DUPLICATE KEY UPDATE click=?,fav=?,share=?,reply=?,coin=?,dm=?,now_rank=?,his_rank=?,mtime=?,likes=?,dislike=?`
	_upMStatSQL = `REPLACE INTO archive_stat_%02d (aid,click,fav,share,reply,coin,dm,now_rank,his_rank,mtime,likes,dislike) VALUES %s`
	_clickSQL   = "SELECT aid,web,h5,outside,ios,android FROM archive_click_%02d WHERE aid=?"
)

func (d *Dao) hit(id int64) string {
	return fmt.Sprintf("%02d", id%_sharding)
}

func statTbl(aid int64) int64 {
	return aid % _sharding
}

// Stat returns stat info
func (d *Dao) Stat(c context.Context, aid int64) (stat *api.Stat, err error) {
	stat = &api.Stat{}
	err = d.db.QueryRow(c, fmt.Sprintf(_statSQL, d.hit(aid), aid)).Scan(&stat.Aid, &stat.View, &stat.Fav, &stat.Share, &stat.Reply, &stat.Coin, &stat.Danmaku, &stat.NowRank, &stat.HisRank, &stat.Like, &stat.DisLike)
	if err == sql.ErrNoRows {
		err = nil
		stat = nil
	} else if err != nil {
		log.Error("Stat(%v) error(%v)", aid, err)
	}
	return
}

// Update update stat's fields
func (d *Dao) Update(c context.Context, stat *api.Stat) (rows int64, err error) {
	now := time.Now()
	res, err := d.db.Exec(c, fmt.Sprintf(_upStatSQL, d.hit(stat.Aid)), stat.Aid, stat.View, stat.Fav, stat.Share, stat.Reply, stat.Coin, stat.Danmaku, stat.NowRank, stat.HisRank, now, now, stat.Like, stat.DisLike,
		stat.View, stat.Fav, stat.Share, stat.Reply, stat.Coin, stat.Danmaku, stat.NowRank, stat.HisRank, now, stat.Like, stat.DisLike)
	if err != nil {
		log.Error("UpdateStat(%d,%v) error(%v)", stat.Aid, stat, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// MultiUpdate update some stat's fields
func (d *Dao) MultiUpdate(c context.Context, yu int64, stats ...*api.Stat) (rows int64, err error) {
	if len(stats) == 0 {
		return
	}
	const field = `(%d,%d,%d,%d,%d,%d,%d,%d,%d,'%s',%d,%d)`
	var (
		fsqls = make([]string, 0, len(stats))
		now   = time.Now().Format("2006-01-02 15:04:05")
	)
	for _, stat := range stats {
		fsqls = append(fsqls, fmt.Sprintf(field, stat.Aid, stat.View, stat.Fav, stat.Share, stat.Reply, stat.Coin, stat.Danmaku, stat.NowRank, stat.HisRank, now, stat.Like, stat.DisLike))
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_upMStatSQL, yu, strings.Join(fsqls, ",")))
	if err != nil {
		log.Error("upMstat error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// Click archive click.
func (d *Dao) Click(c context.Context, aid int64) (cl *archive.Click3, err error) {
	row := d.clickDB.QueryRow(c, fmt.Sprintf(_clickSQL, statTbl(aid)), aid)
	cl = &archive.Click3{}
	if err = row.Scan(&cl.Aid, &cl.Web, &cl.H5, &cl.Outter, &cl.Ios, &cl.Android); err != nil {
		if err == sql.ErrNoRows {
			cl = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}
