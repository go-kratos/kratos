package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/dm/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_upFilterSharding     = 10
	_countUpFiltersSQL    = "SELECT count(*) FROM dm_filter_up_%02d WHERE mid=? AND type=? AND active=1"
	_countUpFiltersAllSQL = "SELECT count(*) FROM dm_filter_up_%02d WHERE mid=? AND active=1"
	_getUpFiltersSQL      = "SELECT id,oid,type,filter,ctime FROM dm_filter_up_%02d WHERE mid=? AND type=? AND active=1 ORDER BY id DESC limit ?,?"
	_getUpFiltersAllSQL   = "SELECT id,oid,type,filter,ctime FROM dm_filter_up_%02d WHERE mid=? AND active=1 ORDER BY id DESC limit ?,?"
	_uptUpFilterSQL       = "UPDATE dm_filter_up_%02d SET active=? WHERE id=?"
	_updateUpFilterCntSQL = "UPDATE dm_filter_up_count SET count=count+? WHERE mid=? AND type=? AND count<?"
)

func (d *Dao) hitUpFilter(mid int64) int64 {
	return mid % _upFilterSharding
}

// UpFilters return filter rules according type
func (d *Dao) UpFilters(c context.Context, mid, tp, pn, ps int64) (res []*model.UpFilter, total int64, err error) {
	res = make([]*model.UpFilter, 0)
	countRow := d.biliDM.QueryRow(c, fmt.Sprintf(_countUpFiltersSQL, d.hitUpFilter(mid)), mid, tp)
	if err = countRow.Scan(&total); err != nil {
		log.Error("row.ScanCount error(%v)", err)
		return
	}
	rows, err := d.biliDM.Query(c, fmt.Sprintf(_getUpFiltersSQL, d.hitUpFilter(mid)), mid, tp, (pn-1)*ps, ps)
	if err != nil {
		log.Error("dbDM.Query(mid:%d, oid:%d) error(%v)", mid, tp, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		f := &model.UpFilter{}
		if err = rows.Scan(&f.ID, &f.Oid, &f.Type, &f.Filter, &f.Ctime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, f)
	}
	return
}

// UpFiltersAll return all filter rules
func (d *Dao) UpFiltersAll(c context.Context, mid, pn, ps int64) (res []*model.UpFilter, total int64, err error) {
	res = make([]*model.UpFilter, 0)
	countRow := d.biliDM.QueryRow(c, fmt.Sprintf(_countUpFiltersAllSQL, d.hitUpFilter(mid)), mid)
	if err = countRow.Scan(&total); err != nil {
		log.Error("row.ScanCount error(%v)", err)
		return
	}
	rows, err := d.biliDM.Query(c, fmt.Sprintf(_getUpFiltersAllSQL, d.hitUpFilter(mid)), mid, (pn-1)*ps, ps)
	if err != nil {
		log.Error("dbDM.Query(mid:%d) error(%v)", mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		f := &model.UpFilter{}
		if err = rows.Scan(&f.ID, &f.Oid, &f.Type, &f.Filter, &f.Ctime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, f)
	}
	return
}

// UpdateUpFilter batch edit filter.
func (d *Dao) UpdateUpFilter(tx *sql.Tx, mid, id int64, active int8) (affect int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_uptUpFilterSQL, d.hitUpFilter(mid)), active, id)
	if err != nil {
		log.Error("tx.Exec(mid:%d,id:%d,active:%d) error(%v)", mid, id, active, err)
		return
	}
	return res.RowsAffected()
}

// UpdateUpFilterCnt set count
func (d *Dao) UpdateUpFilterCnt(tx *sql.Tx, mid int64, ftype int8, count, limit int) (affect int64, err error) {
	res, err := tx.Exec(_updateUpFilterCntSQL, count, mid, ftype, limit)
	if err != nil {
		log.Error("d.UpdateUpFilterCnt(mid:%d, type:%d, count:%d) error(%v)", mid, ftype, count, err)
		return
	}
	return res.RowsAffected()
}
