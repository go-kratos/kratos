package report

import (
	"context"

	mdlpgc "go-common/app/job/main/tv/model/pgc"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_styleSQL = "SELECT id,style,category FROM tv_ep_season WHERE is_deleted=0 AND `check`=1 AND valid=1"
	_labelSQL = `SELECT name,value,category FROM tv_label WHERE deleted=0 AND param="style_id" AND valid=1`
)

// FindStyle style all .
func (d *Dao) FindStyle(ctx context.Context) (res []*mdlpgc.StyleRes, err error) {
	var (
		rows *sql.Rows
	)
	if rows, err = d.DB.Query(ctx, _styleSQL); err != nil {
		log.Error("d.DB.Query sql(%s) error(%v)", _styleSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &mdlpgc.StyleRes{}
		if err = rows.Scan(&r.ID, &r.Style, &r.Category); err != nil {
			log.Error("d.DB.QueryRow error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// FindLabelID label all .
func (d *Dao) FindLabelID(ctx context.Context) (res map[int]map[string]int, err error) {
	var (
		rows *sql.Rows
		m    map[string]int
	)
	res = make(map[int]map[string]int)
	if rows, err = d.DB.Query(ctx, _labelSQL); err != nil {
		log.Error("d.DB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &mdlpgc.LabelRes{}
		if err = rows.Scan(&r.Name, &r.Value, &r.Category); err != nil {
			log.Error("d.DB.Query Scan error(%v)", err)
			return
		}
		if _, ok := res[r.Category]; ok {
			m[r.Name] = r.Value
		} else {
			m = make(map[string]int)
		}
		res[r.Category] = m
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}
