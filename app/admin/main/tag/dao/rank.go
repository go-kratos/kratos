package dao

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/admin/main/tag/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_rankTopSQL          = "SELECT id,tag_type,tid,tname,highlight,rank,is_business,ctime,mtime FROM rank_top WHERE prid=? AND rid=? AND type=? ORDER BY rank ASC"
	_rankFilterSQL       = "SELECT id,tag_type,tid,tname,rank,ctime,mtime FROM rank_filter WHERE prid=? AND rid=? AND type=? ORDER BY rank ASC"
	_rankCountSQL        = "SELECT id,prid,rid,type,top_count,view_count,filter_count,ctime,mtime FROM rank_count WHERE prid=? AND rid=? AND type=?"
	_rankResultSQL       = "SELECT id,tag_type,tid,tname,highlight,rank,is_business,ctime,mtime FROM rank_result WHERE prid=? AND rid=? AND type=? ORDER BY rank ASC"
	_updateRankCountSQL  = "UPDATE rank_count SET top_count=? ,view_count=?,filter_count=? WHERE id=?"
	_insertRankCountSQL  = "INSERT INTO rank_count(prid,rid,type,top_count,view_count,filter_count) VALUES (?,?,?,?,?,?)"
	_deleteRankTopSQL    = "DELETE FROM rank_top WHERE prid=? AND rid=? AND type=?"
	_deleteRankFilterSQL = "DELETE FROM rank_filter WHERE prid=? AND rid=? AND type=?"
	_deleteRankResultSQL = "DELETE FROM rank_result WHERE prid=? AND rid=? AND type=?"
	_insertRankTopsSQL   = "INSERT INTO rank_top(prid,rid,type,tag_type,tid,tname,highlight,rank,is_business) VALUES %s"
	_insertRankFilterSQL = "INSERT INTO rank_filter(prid,rid,type,tag_type,tid,tname,rank) VALUES %s"
	_insertRankResultSQL = "INSERT INTO rank_result(prid,rid,type,tag_type,tid,tname,highlight,rank,is_business) VALUES %s"
)

// RankCount rank count.
func (d *Dao) RankCount(c context.Context, prid, rid int64, tp int32) (res *model.RankCount, err error) {
	res = new(model.RankCount)
	row := d.db.QueryRow(c, _rankCountSQL, prid, rid, tp)
	if err = row.Scan(&res.ID, &res.Prid, &res.Rid, &res.Type, &res.TopCount, &res.ViewCount, &res.FilterCount, &res.CTime, &res.MTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error("rank count scan(%d,%d,%d) error(%v)", prid, rid, tp, err)
	}
	return
}

// RankTop get rankTop info by prid,rid,type.
func (d *Dao) RankTop(c context.Context, prid, rid int64, tp int32) (top []*model.RankTop, topMap map[int64]*model.RankTop, tids []int64, err error) {
	topMap = make(map[int64]*model.RankTop)
	rows, err := d.db.Query(c, _rankTopSQL, prid, rid, tp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil, nil
		}
		log.Error("query rank_top(%d,%d,%d) error(%v)", prid, rid, tp, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := &model.RankTop{}
		if err = rows.Scan(&t.ID, &t.TagType, &t.Tid, &t.TName, &t.HighLight, &t.Rank, &t.Business, &t.CTime, &t.MTime); err != nil {
			log.Error("rank top scan(%d,%d,%d) error(%v)", prid, rid, tp, err)
			return
		}
		top = append(top, t)
		tids = append(tids, t.Tid)
		topMap[t.Tid] = t
	}
	err = rows.Err()
	return
}

// RankResult get rank result info by prid,rid,type.
func (d *Dao) RankResult(c context.Context, prid, rid int64, tp int32) (resultMap map[int64]*model.RankResult, err error) {
	resultMap = make(map[int64]*model.RankResult)
	rows, err := d.db.Query(c, _rankResultSQL, prid, rid, tp)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("d.dao.RankResult(%d,%d,%d) error(%v)", prid, rid, tp, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := &model.RankResult{}
		if err = rows.Scan(&t.ID, &t.TagType, &t.Tid, &t.TName, &t.HighLight, &t.Rank, &t.Business, &t.CTime, &t.MTime); err != nil {
			log.Error("d.dao.RankResult() scan(%d,%d,%d) error(%v)", prid, rid, tp, err)
			return
		}
		resultMap[t.Tid] = t
	}
	err = rows.Err()
	return
}

// RankFilter rank filter info.
func (d *Dao) RankFilter(c context.Context, prid, rid int64, tp int32) (filter []*model.RankFilter, filterMap map[int64]*model.RankFilter, tids []int64, err error) {
	filterMap = make(map[int64]*model.RankFilter)
	rows, err := d.db.Query(c, _rankFilterSQL, prid, rid, tp)
	if err != nil {
		log.Error("query rank_filter(%d,%d,%d) error(%v)", prid, rid, tp, err)
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	defer rows.Close()
	for rows.Next() {
		f := &model.RankFilter{}
		if err = rows.Scan(&f.ID, &f.TagType, &f.Tid, &f.TName, &f.Rank, &f.CTime, &f.MTime); err != nil {
			log.Error("rank filter scan(%d,%d,%d) error(%v)", prid, rid, tp, err)
			return
		}
		filter = append(filter, f)
		tids = append(tids, f.Tid)
		filterMap[f.Tid] = f
	}
	return
}

// TxUpdateRankCount update rank_count info.
func (d *Dao) TxUpdateRankCount(tx *sql.Tx, id int64, topCount, filterCount, viewCount int) (affect int64, err error) {
	res, err := tx.Exec(_updateRankCountSQL, topCount, viewCount, filterCount, id)
	if err != nil {
		log.Error("update rank count(%d,%d,%d,%d) error(%v)", topCount, viewCount, filterCount, id, err)
		return
	}
	return res.RowsAffected()
}

// TxAddRankCount add rank count.
func (d *Dao) TxAddRankCount(tx *sql.Tx, prid, rid int64, tp int32, topCount, filterCount, viewCount int) (id int64, err error) {
	res, err := tx.Exec(_insertRankCountSQL, prid, rid, tp, topCount, viewCount, filterCount)
	if err != nil {
		log.Error("add rank count(%d,%d,%d,%d,%d,%d) error(%v)", prid, rid, tp, topCount, viewCount, filterCount, err)
		return
	}
	return res.LastInsertId()
}

// TxRemoveRankTop remove rank_top records.
func (d *Dao) TxRemoveRankTop(tx *sql.Tx, prid, rid int64, tp int32) (affect int64, err error) {
	res, err := tx.Exec(_deleteRankTopSQL, prid, rid, tp)
	if err != nil {
		log.Error("remove rank top(%d,%d,%d) error(%v)", prid, rid, tp, err)
		return
	}
	return res.RowsAffected()
}

// TxRemoveRankFilter remove rank_filter records.
func (d *Dao) TxRemoveRankFilter(tx *sql.Tx, prid, rid int64, tp int32) (affect int64, err error) {
	res, err := tx.Exec(_deleteRankFilterSQL, prid, rid, tp)
	if err != nil {
		log.Error("remove rank filter(%d,%d,%d) error(%v)", prid, rid, tp, err)
		return
	}
	return res.RowsAffected()
}

//TxRemoveRankResult remove rank_result records.
func (d *Dao) TxRemoveRankResult(tx *sql.Tx, prid, rid int64, tp int32) (affect int64, err error) {
	res, err := tx.Exec(_deleteRankResultSQL, prid, rid, tp)
	if err != nil {
		log.Error("remove rank result(%d,%d,%d) error(%v)", prid, rid, tp, err)
		return
	}
	return res.RowsAffected()
}

// TxInsertRankTop tnsert records into rank_top table.
func (d *Dao) TxInsertRankTop(tx *sql.Tx, prid, rid int64, tp int32, rankTop []*model.RankTop) (id int64, err error) {
	var (
		sql    []string
		sqlStr = " (%d,%d,%d,%d,%d,%q,%d,%d,%d) "
	)
	for _, v := range rankTop {
		s := fmt.Sprintf(sqlStr, prid, rid, tp, v.TagType, v.Tid, v.TName, v.HighLight, v.Rank, v.Business)
		sql = append(sql, s)
	}
	insertSQL := strings.Join(sql, " , ")
	res, err := tx.Exec(fmt.Sprintf(_insertRankTopsSQL, insertSQL))
	if err != nil {
		log.Error("insert rank top(%s) error(%v)", insertSQL, err)
		return
	}
	return res.LastInsertId()
}

// TxInsertRankFilter insert records into rank_filter table.
func (d *Dao) TxInsertRankFilter(tx *sql.Tx, prid, rid int64, tp int32, rankFilter []*model.RankFilter) (id int64, err error) {
	var (
		sql    []string
		sqlStr = " (%d,%d,%d,%d,%d,%q,%d) "
	)
	for _, v := range rankFilter {
		s := fmt.Sprintf(sqlStr, prid, rid, tp, v.TagType, v.Tid, v.TName, v.Rank)
		sql = append(sql, s)
	}
	insertSQL := strings.Join(sql, " , ")
	res, err := tx.Exec(fmt.Sprintf(_insertRankFilterSQL, insertSQL))
	if err != nil {
		log.Error("insert rank filter(%s) error(%v)", insertSQL, err)
		return
	}
	return res.LastInsertId()
}

// TxInsertRankResult insert records into rank_result table.
func (d *Dao) TxInsertRankResult(tx *sql.Tx, prid, rid int64, tp int32, rankLists []*model.RankResult) (id int64, err error) {
	var (
		sql    []string
		sqlStr = " (%d,%d,%d,%d,%d,%q,%d,%d,%d) "
	)
	for _, v := range rankLists {
		s := fmt.Sprintf(sqlStr, prid, rid, tp, v.TagType, v.Tid, v.TName, v.HighLight, v.Rank, v.Business)
		sql = append(sql, s)
	}
	insertSQL := strings.Join(sql, " , ")
	res, err := tx.Exec(fmt.Sprintf(_insertRankResultSQL, insertSQL))
	if err != nil {
		log.Error("insert rank result(%s) error(%v)", insertSQL, err)
		return
	}
	return res.LastInsertId()
}
