package dao

import (
	"context"
	"time"

	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// RecommendByCategory find recommend by category
func (d *Dao) RecommendByCategory(c context.Context, categoryID int64) (res []*artmdl.Recommend, err error) {
	ts := time.Now().Unix()
	rows, err := d.recommendCategoryStmt.Query(c, ts, ts, categoryID)
	if err != nil {
		PromError("db:推荐列表")
		log.Error("dao.recommendCategoryStmt.Query error(%+v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			r = &artmdl.Recommend{Rec: true}
		)
		if err = rows.Scan(&r.ArticleID, &r.RecImageURL, &r.RecFlag, &r.Position, &r.EndTime, &r.RecImageStartTime, &r.RecImageEndTime); err != nil {
			PromError("db:推荐列表scan")
			log.Error("dao.RecommendByCategory.rows.Scan error(%+v)", err)
			return
		}
		r.RecImageURL = artmdl.CompleteURL(r.RecImageURL)
		res = append(res, r)
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// DelRecommend delete recommend
func (d *Dao) DelRecommend(c context.Context, aid int64) (err error) {
	if _, err := d.delRecommendStmt.Exec(c, time.Now().Unix(), aid); err != nil {
		PromError("db:删除推荐")
		log.Error("dao.delRecommendStmt.Exec(%v) error(%+v)", aid, err)
	}
	return
}

// AllRecommends .
func (d *Dao) AllRecommends(c context.Context, t time.Time, pn, ps int) (res []int64, err error) {
	ts := t.Unix()
	offset := (pn - 1) * ps
	rows, err := d.allRecommendStmt.Query(c, ts, ts, offset, ps)
	if err != nil {
		PromError("db:全部推荐列表")
		log.Error("dao.AllRecommends(pn: %v ps: %v)error(%+v)", pn, ps, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			aid int64
		)
		if err = rows.Scan(&aid); err != nil {
			PromError("db:全部推荐列表")
			log.Error("dao.AllRecommends.rows.Scan error(%+v)", err)
			return
		}
		res = append(res, aid)
	}
	err = rows.Err()
	promErrorCheck(err)
	return
}

// AllRecommendCount .
func (d *Dao) AllRecommendCount(c context.Context, t time.Time) (res int64, err error) {
	ts := t.Unix()
	if err = d.allRecommendCountStmt.QueryRow(c, ts, ts).Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		PromError("db:推荐列表计数")
		log.Error("dao.AllRecommendCount() error(%+v)", err)
	}
	return
}
