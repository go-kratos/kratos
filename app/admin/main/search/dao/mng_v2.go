package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/admin/main/search/model"
	"go-common/library/log"
)

const (
	_typeDatabus = `databus`
	_typeDB      = `db`
	_typeTable   = `table`

	_businessAllV2SQL   = `SELECT id,pid,name,description,state FROM gf_business`
	_businessInfoV2SQL  = `SELECT id,pid,name,description,data_conf,index_conf,business_conf,state,mtime FROM gf_business WHERE name=?`
	_bussinessInsSQL    = `INSERT INTO gf_business (pid,name,description) VALUES(?,?,?)`
	_bussinessUpdateSQL = `UPDATE gf_business SET %s=? WHERE name=?`

	_assetDBTablesV2SQL  = `SELECT id,type,db,name,regex,fields,description,state FROM gf_asset WHERE type=? OR type=?`
	_assetDBInsSQL       = `INSERT INTO gf_asset (type,name,description,dsn) VALUES(?,?,?,?)`
	_assetTableInsSQL    = `INSERT INTO gf_asset (type,name,db,regex,fields,description) VALUES(?,?,?,?,?,?)`
	_assetTableUpdateSQL = `UPDATE gf_asset set fields=? WHERE name=?`
	_assetSQL            = `SELECT id,type,name,dsn,db,regex,fields,description,state FROM gf_asset WHERE name=?`
)

// BusinessAllV2 .
func (d *Dao) BusinessAllV2(c context.Context) (list []*model.GFBusiness, err error) {
	rows, err := d.db.Query(c, _businessAllV2SQL)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := new(model.GFBusiness)
		if err = rows.Scan(&t.ID, &t.PID, &t.Name, &t.Description, &t.State); err != nil {
			return
		}
		list = append(list, t)
	}
	err = rows.Err()
	return
}

// BusinessInfoV2 .
func (d *Dao) BusinessInfoV2(c context.Context, name string) (b *model.GFBusiness, err error) {
	row := d.db.QueryRow(c, _businessInfoV2SQL, name)
	if err != nil {
		return
	}
	b = new(model.GFBusiness)
	if err = row.Scan(&b.ID, &b.PID, &b.Name, &b.Description, &b.DataConf, &b.IndexConf, &b.BusinessConf, &b.State, &b.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			b = nil
			return
		}
	}
	tm, _ := time.Parse(time.RFC3339, b.Mtime)
	b.Mtime = tm.Format("2006-01-02 15:04:05")
	return
}

// BusinessIns insert business.
func (d *Dao) BusinessIns(c context.Context, pid int64, name, description string) (rows int64, err error) {
	res, err := d.db.Exec(c, _bussinessInsSQL, pid, name, description)
	if err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// BusinessUpdate update business.
func (d *Dao) BusinessUpdate(c context.Context, name, field, value string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_bussinessUpdateSQL, field), value, name)
	if err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// AssetDBTables .
func (d *Dao) AssetDBTables(c context.Context) (list []*model.GFAsset, err error) {
	rows, err := d.db.Query(c, _assetDBTablesV2SQL, _typeDB, _typeTable)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := new(model.GFAsset)
		if err = rows.Scan(&t.ID, &t.Type, &t.DB, &t.Name, &t.Regex, &t.Fields, &t.Description, &t.State); err != nil {
			return
		}
		list = append(list, t)
	}
	err = rows.Err()
	return
}

// AssetDBIns insert db asset.
func (d *Dao) AssetDBIns(c context.Context, name, description, dsn string) (rows int64, err error) {
	res, err := d.db.Exec(c, _assetDBInsSQL, _typeDB, name, description, dsn)
	if err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// AssetTableIns insert table asset.
func (d *Dao) AssetTableIns(c context.Context, name, db, regex, fields, description string) (rows int64, err error) {
	res, err := d.db.Exec(c, _assetTableInsSQL, _typeTable, name, db, regex, fields, description)
	if err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// UpdateAssetTable update table asset.
func (d *Dao) UpdateAssetTable(c context.Context, name, fields string) (rows int64, err error) {
	res, err := d.db.Exec(c, _assetTableUpdateSQL, fields, name)
	if err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// Asset .
func (d *Dao) Asset(c context.Context, name string) (r *model.GFAsset, err error) {
	row := d.db.QueryRow(c, _assetSQL, name)
	r = new(model.GFAsset)
	if err = row.Scan(&r.ID, &r.Type, &r.Name, &r.DSN, &r.DB, &r.Regex, &r.Fields, &r.Description, &r.State); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			r = nil
			return
		}
	}
	return
}
