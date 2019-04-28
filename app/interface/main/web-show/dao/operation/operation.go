package operation

import (
	"context"
	"time"

	"go-common/app/interface/main/web-show/model/operation"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_selNoticeSQL = "SELECT id,type,content,link,ads,pic,rank FROM operations where stime <= ? AND etime > ? ORDER BY stime DESC"
)

//Operation dao
func (dao *Dao) Operation(c context.Context) (ns []*operation.Operation, err error) {
	now := xtime.Time(time.Now().Unix())
	rows, err := dao.db.Query(c, _selNoticeSQL, now, now)
	if err != nil {
		log.Error("notice.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		n := &operation.Operation{}
		if err = rows.Scan(&n.ID, &n.Type, &n.Message, &n.Link, &n.Ads, &n.Pic, &n.Rank); err != nil {
			PromError("Operation", "rows.scan err(%v)", err)
			return
		}
		ns = append(ns, n)
	}
	return
}
