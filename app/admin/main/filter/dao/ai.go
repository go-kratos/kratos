package dao

import (
	"context"
	"database/sql"
	"go-common/app/admin/main/filter/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_selWhite      = "SELECT id,mid,state,ctime,mtime FROM filter_ai_white ORDER BY mtime DESC LIMIT ?,?"
	_selWhiteByMid = "SELECT id FROM filter_ai_white WHERE mid=?"
	_selWhiteCount = "SELECT count(*) FROM filter_ai_white"
	_insWhite      = "INSERT INTO filter_ai_white (mid) VALUES (?)"
	_editWhite     = "UPDATE filter_ai_white SET state=? WHERE mid=?"
	_insCase       = "INSERT INTO filter_ai_badcase (source,content,type) VALUES (?,?,?)"
)

// AiWhite get AI white mids.
func (d *Dao) AiWhite(c context.Context, pn, ps int) (res []*model.AiWhite, err error) {
	var rows *xsql.Rows
	res = make([]*model.AiWhite, 0)
	var (
		start  = (pn - 1) * ps
		offset = ps
	)
	if rows, err = d.mysql.Query(c, _selWhite, start, offset); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			re = &model.AiWhite{}
		)
		if err = rows.Scan(&re.ID, &re.MID, &re.State, &re.Ctime, &re.Mtime); err != nil {
			return
		}
		res = append(res, re)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

// AiWhiteCount AI white count.
func (d *Dao) AiWhiteCount(c context.Context) (num int64, err error) {
	row := d.mysql.QueryRow(c, _selWhiteCount)
	if err != nil {
		log.Error("AiWhiteCount: d.mysql.QueryRow() error(%v)", err)
		return
	}
	if err = row.Scan(&num); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// AiWhiteByMid AI white by mid.
func (d *Dao) AiWhiteByMid(c context.Context, mid int64) (id int64, err error) {
	row := d.mysql.QueryRow(c, _selWhiteByMid, mid)
	if err != nil {
		log.Error("AiWhiteByMid: d.mysql.QueryRow(%d) error(%v)", mid, err)
		return
	}
	if err = row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// InsertAiWhite add white.
func (d *Dao) InsertAiWhite(c context.Context, mid int64) (afftected int64, err error) {
	var res sql.Result
	if res, err = d.mysql.Exec(c, _insWhite, mid); err != nil {
		log.Error("InsertAiWhite, d.mysql.Exec(%d) error(%v)", mid, err)
		return
	}
	return res.RowsAffected()
}

// EditAiWhite edit white.
func (d *Dao) EditAiWhite(c context.Context, mid int64, state int8) (afftected int64, err error) {
	var res sql.Result
	if res, err = d.mysql.Exec(c, _editWhite, state, mid); err != nil {
		log.Error("EditAiWhite, d.mysql.Exec(%d) error(%v)", mid, err)
		return
	}
	return res.RowsAffected()
}

// InsertAiCase add case.
func (d *Dao) InsertAiCase(c context.Context, aiCase *model.AiCase) (afftected int64, err error) {
	var res sql.Result
	if res, err = d.mysql.Exec(c, _insCase, aiCase.Source, aiCase.Content, aiCase.Type); err != nil {
		log.Error("InsertAiCase, d.mysql.Exec(%+v) error(%v)", aiCase, err)
		return
	}
	return res.RowsAffected()
}
