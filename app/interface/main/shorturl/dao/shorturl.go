package dao

import (
	"context"
	"database/sql"

	"fmt"
	"go-common/app/interface/main/shorturl/model"
	"go-common/library/log"
)

const (
	_prefix = "su_"
	//short_url
	_getSQL         = "SELECT id,mid,short_url,long_url,state,ctime,mtime FROM short_url WHERE short_url=?"
	_getAllSQL      = "SELECT id,mid,short_url,long_url,state,ctime,mtime FROM short_url"
	_getLimitSQL    = "SELECT id,mid,short_url,long_url,state,ctime,mtime FROM short_url WHERE state=? " // TODO limit page
	_shortInSQL     = "INSERT IGNORE INTO short_url(mid,short_url,long_url,state,ctime) VALUES(?,?,?,?,?)"
	_shortUpSQL     = "UPDATE short_url SET mid=?,long_url=? WHERE id=?"
	_shortCountSQL  = "SELECT COUNT(*) FROM short_url WHERE state=?"
	_shortByIDSQL   = "SELECT id,mid,short_url,long_url,state,ctime,mtime FROM short_url WHERE id=?"
	_updateStateSQL = "UPDATE short_url SET mid=?,state=? WHERE id=?"
)

// Short get short_url
func (d *Dao) Short(ctx context.Context, short string) (res *model.ShortUrl, err error) {
	rows := d.db.QueryRow(ctx, _getSQL, short)
	res = &model.ShortUrl{}
	if err = rows.Scan(&res.ID, &res.Mid, &res.Short, &res.Long, &res.State, &res.CTime, &res.MTime); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("rows.Scan error(%v)", err)
		}
		return
	}
	return
}

// ShortbyID get short_url by id
func (d *Dao) ShortbyID(ctx context.Context, id int64) (res *model.ShortUrl, err error) {
	rows := d.db.QueryRow(ctx, _shortByIDSQL, id)
	res = &model.ShortUrl{}
	if err = rows.Scan(&res.ID, &res.Mid, &res.Short, &res.Long, &res.State, &res.CTime, &res.MTime); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("rows.Scan err (%v)", err)
		}
		return
	}
	res.FormatDate()
	return
}

// AllShorts get all short_url
func (d *Dao) AllShorts(ctx context.Context) (res []*model.ShortUrl, err error) {
	rows, err := d.db.Query(ctx, _getAllSQL)
	if err != nil {
		log.Error("query error (%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		su := &model.ShortUrl{}
		if err = rows.Scan(&su.ID, &su.Mid, &su.Short, &su.Long, &su.State, &su.CTime, &su.MTime); err != nil {
			log.Error("rows.Scan err (%v)", err)
			return
		}
		su.FormatDate()
		res = append(res, su)
	}
	return
}

// ShortCount get all short_url
func (d *Dao) ShortCount(ctx context.Context, mid int64, long string) (count int, err error) {
	countSQL := _shortCountSQL
	if mid > 0 {
		countSQL = fmt.Sprintf("%s AND mid=%d", countSQL, mid)
	}
	if long != "" {
		countSQL += " AND long_url='" + long + "'"
	}
	row := d.db.QueryRow(ctx, countSQL, model.StateNormal)
	if err = row.Scan(&count); err != nil {
		log.Error("row.Scan error(%v)", err)
		return
	}
	return
}

// InShort add short_url
func (d *Dao) InShort(ctx context.Context, su *model.ShortUrl) (id int64, err error) {
	res, err := d.db.Exec(ctx, _shortInSQL, su.Mid, su.Short, su.Long, su.State, su.CTime)
	if err != nil {
		log.Error("tx.Exec() error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// ShortUp add short_url
func (d *Dao) ShortUp(ctx context.Context, id, mid int64, long string) (rows int64, err error) {
	res, err := d.db.Exec(ctx, _shortUpSQL, mid, long, id)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// UpdateState update state
func (d *Dao) UpdateState(ctx context.Context, id, mid int64, state int) (rows int64, err error) {
	_, err = d.db.Exec(ctx, _updateStateSQL, mid, state, id)
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", _updateStateSQL, err)
		return
	}
	return
}

// ShortLimit get short_url list
func (d *Dao) ShortLimit(ctx context.Context, pn, ps int, mid int64, long string) (res []*model.ShortUrl, err error) {
	limitSQL := _getLimitSQL
	if mid > 0 {
		limitSQL = fmt.Sprintf("%s AND mid=%d", limitSQL, mid)
	}
	if long != "" {
		limitSQL += " AND long_url='" + long + "'"
	}
	limitSQL += " ORDER BY id DESC LIMIT ?,? "
	rows, err := d.db.Query(ctx, limitSQL, model.StateNormal, pn, ps)
	if err != nil {
		log.Error("query error (%v)", err)
		return
	}
	defer rows.Close()
	res = []*model.ShortUrl{}
	for rows.Next() {
		su := &model.ShortUrl{}
		if err = rows.Scan(&su.ID, &su.Mid, &su.Short, &su.Long, &su.State, &su.CTime, &su.MTime); err != nil {
			log.Error("rows.Scan err (%v)", err)
			return
		}
		su.FormatDate()
		res = append(res, su)
	}
	return
}
