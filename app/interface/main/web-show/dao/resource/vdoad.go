package resource

import (
	"context"
	"database/sql"
	"time"

	"go-common/app/interface/main/web-show/model/resource"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_selAdVdoActSQL   = `SELECT id,name,aid,ad_cid,ad_url,skipable,ad_strategy,mtime FROM video_ads WHERE platform=0 and type=0 and state=0 AND verified=1 AND starttime<? AND endtime>?`
	_selAdMtCntVdoSQL = `SELECT FROM_UNIXTIME(ROUND(AVG(UNIX_TIMESTAMP(mtime)))) FROM video_ads WHERE platform=0 and type=0 and state=0 AND verified=1 AND starttime<? AND endtime>?`
)

func (dao *Dao) initAd() {
	dao.selAdVdoActStmt = dao.videodb.Prepared(_selAdVdoActSQL)
	dao.selAdMtCntVdoStmt = dao.videodb.Prepared(_selAdMtCntVdoSQL)
}

// AllVideoActive dao
func (dao *Dao) AllVideoActive(c context.Context, now time.Time) (ads []resource.VideoAD, err error) {
	var rows *xsql.Rows
	if rows, err = dao.selAdVdoActStmt.Query(c, now, now); err != nil {
		log.Error("dao..Exec(%v, %v), err (%v)", now, now, err)
		return
	}
	defer rows.Close()
	ads = make([]resource.VideoAD, 0, 100)
	for rows.Next() {
		ad := resource.VideoAD{}
		if err = rows.Scan(&ad.ID, &ad.Name, &ad.AidS, &ad.Cid, &ad.URL, &ad.Skipable, &ad.Strategy, &ad.MTime); err != nil {
			log.Error("rows.Scan(), err (%v)", err)
			ads = nil
			return
		}
		ads = append(ads, ad)
	}
	return
}

// AdVideoMTimeCount dao
func (dao *Dao) AdVideoMTimeCount(c context.Context, now time.Time) (mtime xtime.Time, err error) {
	row := dao.selAdMtCntVdoStmt.QueryRow(c, now, now)
	if err = row.Scan(&mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("mysqlDB error(%v)", err)
		}
	}
	return
}
