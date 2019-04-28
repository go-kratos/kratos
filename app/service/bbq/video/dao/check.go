package dao

import (
	"context"
	"database/sql"
	"fmt"
	"go-common/library/log"
)

const (
	_queryBvcResource   = "select id from %s where svid = %d"
	_queryCoverResource = "select cover_url,cover_width,cover_height from video_repository where svid = ?"
)

//CheckSVResource ...
func (d *Dao) CheckSVResource(c context.Context, svid int64) (err error) {
	var (
		ID   int64
		cURL string
		cH   int64
		cW   int64
	)
	tN := fmt.Sprintf("video_bvc_%02d", svid%100)

	if err = d.db.QueryRow(c, fmt.Sprintf(_queryBvcResource, tN, svid)).Scan(&ID); err == sql.ErrNoRows {
		log.Error("CheckSVResource bvc err,svid:%d,err:%v", svid, err)
		return
	}
	//cover,err := d.cmsdb.QueryRow(c, query, ...)
	if err = d.db.QueryRow(c, _queryCoverResource, svid).Scan(&cURL, &cW, &cH); err == sql.ErrNoRows {
		log.Error("CheckSVResource cover err,svid:%d,err:%v", svid, err)
		return
	}
	return
}
