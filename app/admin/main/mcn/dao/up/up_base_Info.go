package up

import (
	"context"
	"fmt"

	"go-common/app/admin/main/mcn/model"
	xsql "go-common/library/database/sql"
	"go-common/library/xstr"
)

const (
	_selUpBaseInfoMapSQL = "SELECT mid,fans_count,active_tid,article_count_accumulate FROM up_base_info WHERE mid IN (%s) AND business_type = 1"
	_selUpPlayInfoMapSQL = "SELECT mid,play_count_accumulate,article_count FROM up_play_info WHERE mid IN (%s) AND business_type = 1"
)

// UpBaseInfoMap .
func (d *Dao) UpBaseInfoMap(c context.Context, mids []int64) (mbi map[int64]*model.UpBaseInfo, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_selUpBaseInfoMapSQL, xstr.JoinInts(mids))); err != nil {
		return
	}
	defer rows.Close()
	mbi = make(map[int64]*model.UpBaseInfo, len(mids))
	for rows.Next() {
		bi := new(model.UpBaseInfo)
		if err = rows.Scan(&bi.Mid, &bi.FansCount, &bi.ActiveTid, &bi.ArticleCountAccumulate); err != nil {
			return
		}
		mbi[bi.Mid] = bi
	}
	err = rows.Err()
	return
}

// UpPlayInfoMap .
func (d *Dao) UpPlayInfoMap(c context.Context, mids []int64) (mpi map[int64]*model.UpPlayInfo, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_selUpPlayInfoMapSQL, xstr.JoinInts(mids))); err != nil {
		return
	}
	defer rows.Close()
	mpi = make(map[int64]*model.UpPlayInfo, len(mids))
	for rows.Next() {
		pi := new(model.UpPlayInfo)
		if err = rows.Scan(&pi.Mid, &pi.PlayCountAccumulate, &pi.ArticleCount); err != nil {
			return
		}
		if pi.ArticleCount != 0 {
			pi.PlayCountAverage = pi.PlayCountAccumulate / pi.ArticleCount
		}
		mpi[pi.Mid] = pi
	}
	err = rows.Err()
	return
}
