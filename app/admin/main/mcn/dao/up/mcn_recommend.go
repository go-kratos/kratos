package up

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"sync"

	"go-common/app/admin/main/mcn/model"
	xsql "go-common/library/database/sql"
	"go-common/library/xstr"
)

const (
	// private condition
	_inMcnUpRecommendSQL = `INSERT mcn_up_recommend_pool(up_mid,fans_count,archive_count,play_count_accumulate,play_count_average,active_tid,source) VALUES (?,?,?,?,?,?,2) 
ON DUPLICATE KEY UPDATE fans_count=?, archive_count=?, play_count_accumulate=?, play_count_average=?, active_tid=?, state=1, source=2, fans_count_increase_month=0,
last_archive_time='1970-01-01 08:00:00', generate_time='1970-01-01 08:00:00'`
	_upMcnUpRecommendOPSQL = "UPDATE mcn_up_recommend_pool SET state = ?  WHERE up_mid IN (%s)"
	_mcnUpRecommendsSQL    = `SELECT id,up_mid,fans_count,fans_count_increase_month,archive_count,
play_count_accumulate,play_count_average,active_tid,last_archive_time,state,source,generate_time,ctime,mtime FROM mcn_up_recommend_pool %s`
	_mcnUpRecommendTotalSQL = "SELECT COUNT(1) FROM mcn_up_recommend_pool %s"
	_mcnUpBindMidsSQL       = "SELECT up_mid FROM mcn_up WHERE up_mid IN (%s) AND state IN (2,10,11,15)"
	_mcnUpRecommendMidSQL   = `SELECT id,up_mid,fans_count,fans_count_increase_month,archive_count,
play_count_accumulate,play_count_average,active_tid,last_archive_time,state,source,generate_time,ctime,mtime FROM mcn_up_recommend_pool WHERE up_mid = ?`
	_mcnUpRecommendMidsSQL = `SELECT id,up_mid,fans_count,fans_count_increase_month,archive_count,
play_count_accumulate,play_count_average,active_tid,last_archive_time,state,source,generate_time,ctime,mtime FROM mcn_up_recommend_pool WHERE up_mid IN (%s)`

	// common condition
	_orderByConditionSQL           = " %s ORDER BY %s %s LIMIT ?,?"
	_orderByConditionNotLimitSQL   = " %s ORDER BY %s %s"
	_orderByNoConditionSQL         = " ORDER BY %s %s LIMIT ?,?"
	_orderByNoConditionNotLimitSQL = " ORDER BY %s %s"
)

// AddMcnUpRecommend .
func (d *Dao) AddMcnUpRecommend(c context.Context, arg *model.McnUpRecommendPool) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _inMcnUpRecommendSQL, arg.UpMid, arg.FansCount, arg.ArchiveCount, arg.PlayCountAccumulate, arg.PlayCountAverage, arg.ActiveTid,
		arg.FansCount, arg.ArchiveCount, arg.PlayCountAccumulate, arg.PlayCountAverage, arg.ActiveTid); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// UpMcnUpsRecommendOP .
func (d *Dao) UpMcnUpsRecommendOP(c context.Context, upMids []int64, state model.MCNUPRecommendState) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_upMcnUpRecommendOPSQL, xstr.JoinInts(upMids)), state); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// McnUpRecommends .
func (d *Dao) McnUpRecommends(c context.Context, arg *model.MCNUPRecommendReq) (res []*model.McnUpRecommendPool, err error) {
	sql, values := d.buildUpRecommendSQL("list", arg)
	rows, err := d.db.Query(c, sql, values...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.McnUpRecommendPool{}
		err = rows.Scan(&m.ID, &m.UpMid, &m.FansCount, &m.FansCountIncreaseMonth, &m.ArchiveCount, &m.PlayCountAccumulate, &m.PlayCountAverage,
			&m.ActiveTid, &m.LastArchiveTime, &m.State, &m.Source, &m.GenerateTime, &m.Ctime, &m.Mtime)
		if err != nil {
			if err == xsql.ErrNoRows {
				err = nil
			}
			res = nil
			return
		}
		res = append(res, m)
	}
	return
}

// McnUpRecommendTotal .
func (d *Dao) McnUpRecommendTotal(c context.Context, arg *model.MCNUPRecommendReq) (count int, err error) {
	sql, values := d.buildUpRecommendSQL("count", arg)
	row := d.db.QueryRow(c, sql, values...)
	if err = row.Scan(&count); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			return
		}
	}
	return
}

// McnUpBindMids .
func (d *Dao) McnUpBindMids(c context.Context, mids []int64) (bmids []int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_mcnUpBindMidsSQL, xstr.JoinInts(mids)))
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var upMids int64
		err = rows.Scan(&upMids)
		if err != nil {
			if err == xsql.ErrNoRows {
				err = nil
			}
			return
		}
		bmids = append(bmids, upMids)
	}
	return
}

// McnUpRecommendMid .
func (d *Dao) McnUpRecommendMid(c context.Context, mid int64) (m *model.McnUpRecommendPool, err error) {
	row := d.db.QueryRow(c, _mcnUpRecommendMidSQL, mid)
	m = &model.McnUpRecommendPool{}
	err = row.Scan(&m.ID, &m.UpMid, &m.FansCount, &m.FansCountIncreaseMonth, &m.ArchiveCount, &m.PlayCountAccumulate, &m.PlayCountAverage,
		&m.ActiveTid, &m.LastArchiveTime, &m.State, &m.Source, &m.GenerateTime, &m.Ctime, &m.Mtime)
	if err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		}
		m = nil
		return
	}
	return
}

// McnUpRecommendMids .
func (d *Dao) McnUpRecommendMids(c context.Context, mids []int64) (mrp map[int64]*model.McnUpRecommendPool, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_mcnUpRecommendMidsSQL, xstr.JoinInts(mids)))
	if err != nil {
		return
	}
	defer rows.Close()
	mrp = make(map[int64]*model.McnUpRecommendPool, len(mids))
	for rows.Next() {
		m := &model.McnUpRecommendPool{}
		err = rows.Scan(&m.ID, &m.UpMid, &m.FansCount, &m.FansCountIncreaseMonth, &m.ArchiveCount, &m.PlayCountAccumulate, &m.PlayCountAverage,
			&m.ActiveTid, &m.LastArchiveTime, &m.State, &m.Source, &m.GenerateTime, &m.Ctime, &m.Mtime)
		if err != nil {
			if err == xsql.ErrNoRows {
				err = nil
			}
			mrp = nil
			return
		}
		mrp[m.UpMid] = m
	}
	return
}

// buildUpRecommendSQL build a up recommend sql string.
func (d *Dao) buildUpRecommendSQL(tp string, arg *model.MCNUPRecommendReq) (sql string, values []interface{}) {
	values = make([]interface{}, 0, 11)
	var (
		cond    []string
		condStr string
	)
	if arg.TID != 0 {
		cond = append(cond, "active_tid = ?")
		values = append(values, arg.TID)
	}
	if arg.UpMid != 0 {
		cond = append(cond, "up_mid = ?")
		values = append(values, arg.UpMid)
	}
	if arg.FansMin != 0 {
		cond = append(cond, "fans_count >= ?")
		values = append(values, arg.FansMin)
	}
	if arg.FansMax != 0 {
		cond = append(cond, "fans_count <= ?")
		values = append(values, arg.FansMax)
	}
	if arg.PlayMin != 0 {
		cond = append(cond, "play_count_accumulate >= ?")
		values = append(values, arg.PlayMin)
	}
	if arg.PlayMax != 0 {
		cond = append(cond, "play_count_accumulate <= ?")
		values = append(values, arg.PlayMax)
	}
	if arg.PlayAverageMin != 0 {
		cond = append(cond, "play_count_average >= ?")
		values = append(values, arg.PlayAverageMin)
	}
	if arg.PlayAverageMax != 0 {
		cond = append(cond, "play_count_average <= ?")
		values = append(values, arg.PlayAverageMax)
	}
	if arg.State != model.MCNUPRecommendStateUnKnown {
		cond = append(cond, "state = ?")
		values = append(values, arg.State)
	} else {
		cond = append(cond, "state IN (1,2,3)")
	}
	if arg.Source != model.MCNUPRecommendSourceUnKnown {
		cond = append(cond, "source = ?")
		values = append(values, arg.Source)
	}
	condStr = d.joinStringSQL(cond)
	switch tp {
	case "count":
		if condStr != "" {
			sql = fmt.Sprintf(_mcnUpRecommendTotalSQL+" %s", "WHERE", condStr)
			return
		}
		sql = fmt.Sprintf(_mcnUpRecommendTotalSQL, condStr)
	case "list":
		// 导出
		if arg.Export == model.ResponeModelCSV {
			if condStr != "" {
				sql = fmt.Sprintf(_mcnUpRecommendsSQL+_orderByConditionNotLimitSQL, "WHERE", condStr, arg.Order, arg.Sort)
				return
			}
			sql = fmt.Sprintf(_mcnUpRecommendsSQL+_orderByNoConditionNotLimitSQL, condStr, arg.Order, arg.Sort)
			return
		}
		// 非导出
		if condStr != "" {
			sql = fmt.Sprintf(_mcnUpRecommendsSQL+_orderByConditionSQL, "WHERE", condStr, arg.Order, arg.Sort)
		} else {
			sql = fmt.Sprintf(_mcnUpRecommendsSQL+_orderByNoConditionSQL, condStr, arg.Order, arg.Sort)
		}
		limit, offset := arg.PageArg.CheckPageValidation()
		values = append(values, offset, limit)
	}
	return
}

func (d *Dao) joinStringSQL(is []string) string {
	if len(is) == 0 {
		return ""
	}
	if len(is) == 1 {
		return is[0]
	}
	var bfPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer([]byte{})
		},
	}
	buf := bfPool.Get().(*bytes.Buffer)
	for _, i := range is {
		buf.WriteString(i)
		buf.WriteString(" AND ")
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 4)
	}
	s := buf.String()
	buf.Reset()
	bfPool.Put(buf)
	return s
}
