package archive

import (
	"context"
	"fmt"

	"go-common/app/service/main/archive/api"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_statSharding = 100
	// stat
	_statSQL  = "SELECT aid,fav,share,reply,coin,dm,click,now_rank,his_rank,likes FROM archive_stat_%02d WHERE aid=?"
	_statsSQL = "SELECT aid,fav,share,reply,coin,dm,click,now_rank,his_rank,likes FROM archive_stat_%02d WHERE aid in (%s)"
	// click
	_clickSQL = "SELECT aid,web,h5,outside,ios,android FROM archive_click_%02d WHERE aid=?"
)

func statTbl(aid int64) int64 {
	return aid % _statSharding
}

// stat3 archive stat.
func (d *Dao) stat3(c context.Context, aid int64) (st *api.Stat, err error) {
	d.infoProm.Incr("stat3")
	row := d.statDB.QueryRow(c, fmt.Sprintf(_statSQL, statTbl(aid)), aid)
	st = &api.Stat{}
	if err = row.Scan(&st.Aid, &st.Fav, &st.Share, &st.Reply, &st.Coin, &st.Danmaku, &st.View, &st.NowRank, &st.HisRank, &st.Like); err != nil {
		if err == sql.ErrNoRows {
			st = nil
			err = nil
		} else {
			d.errProm.Incr("stat_db")
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	return
}

// stats3 archive stats.
func (d *Dao) stats3(c context.Context, aids []int64) (sts map[int64]*api.Stat, err error) {
	d.infoProm.Incr("stats3")
	tbls := make(map[int64][]int64)
	for _, aid := range aids {
		tbls[statTbl(aid)] = append(tbls[statTbl(aid)], aid)
	}
	sts = make(map[int64]*api.Stat, len(aids))
	for tbl, ids := range tbls {
		var rows *sql.Rows
		if rows, err = d.statDB.Query(c, fmt.Sprintf(_statsSQL, tbl, xstr.JoinInts(ids))); err != nil {
			log.Error("d.statDB.Query(%s) error(%v)", fmt.Sprintf(_statsSQL, tbl, xstr.JoinInts(ids)), err)
			d.errProm.Incr("stat_db")
			return
		}
		for rows.Next() {
			st := &api.Stat{}
			if err = rows.Scan(&st.Aid, &st.Fav, &st.Share, &st.Reply, &st.Coin, &st.Danmaku, &st.View, &st.NowRank, &st.HisRank, &st.Like); err != nil {
				log.Error("rows.Scan error(%v)", err)
				d.errProm.Incr("stat_db")
				rows.Close()
				return
			}
			sts[st.Aid] = st
		}
		rows.Close()
	}
	return
}

// click3 archive click.
func (d *Dao) click3(c context.Context, aid int64) (cl *api.Click, err error) {
	d.infoProm.Incr("click3")
	row := d.clickDB.QueryRow(c, fmt.Sprintf(_clickSQL, statTbl(aid)), aid)
	cl = &api.Click{}
	if err = row.Scan(&cl.Aid, &cl.Web, &cl.H5, &cl.Outter, &cl.Ios, &cl.Android); err != nil {
		if err == sql.ErrNoRows {
			cl = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
			d.errProm.Incr("click_db")
		}
	}
	return
}
