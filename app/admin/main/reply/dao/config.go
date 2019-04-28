package dao

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"go-common/app/admin/main/reply/model"
)

const (
	_addConfigSQL             = "INSERT IGNORE INTO reply_config (type, oid, adminid, operator, category, config, ctime, mtime) VALUES(?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE adminid = ?, operator = ?, config=?, mtime = ?"
	_paginationConfigSQL      = "SELECT id, type, oid, adminid, operator, category, config, ctime, mtime FROM reply_config"
	_paginationConfigCountSQL = "SELECT count(id) FROM reply_config"
	_loadConfigSQL            = "SELECT id, type, oid, adminid, operator, category, config, ctime, mtime FROM reply_config WHERE type=? AND oid=? AND category = ?"
	_loadConfigByIDSQL        = "SELECT id, type, oid, adminid, operator, category, config, ctime, mtime FROM reply_config WHERE id = ?"
	_deleteConfigSQL          = "DELETE FROM reply_config where id = ?"
)

// AddConfig create a new config.
func (d *Dao) AddConfig(c context.Context, typ, category int32, oid, adminid int64, operator, config string, now time.Time) (id int64, err error) {
	res, err := d.db.Exec(c, _addConfigSQL, typ, oid, adminid, operator, category, config, now, now, adminid, operator, config, now)
	if err != nil {
		return
	}
	return res.LastInsertId()
}

// LoadConfig load a config record.
func (d *Dao) LoadConfig(c context.Context, typ, category int32, oid int64) (m *model.Config, err error) {
	m = new(model.Config)
	row := d.db.QueryRow(c, _loadConfigSQL, typ, oid, category)
	if err = row.Scan(&m.ID, &m.Type, &m.Oid, &m.AdminID, &m.Operator, &m.Category, &m.Config, &m.CTime, &m.MTime); err != nil {
		if err == sql.ErrNoRows {
			m = nil
			err = nil
			return
		}
	}
	if m.ID > 0 && len(m.Config) > 0 {
		dat := new(model.Config)
		if err = json.Unmarshal([]byte(m.Config), dat); err == nil {
			m.ShowEntry = dat.ShowEntry
			m.ShowAdmin = dat.ShowAdmin
		}
	}
	return
}

// LoadConfigByID load a config record by id.
func (d *Dao) LoadConfigByID(c context.Context, id int64) (m *model.Config, err error) {
	m = new(model.Config)
	row := d.db.QueryRow(c, _loadConfigByIDSQL, id)
	if err = row.Scan(&m.ID, &m.Type, &m.Oid, &m.AdminID, &m.Operator, &m.Category, &m.Config, &m.CTime, &m.MTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
	}
	if m.ID > 0 {
		dat := new(model.Config)
		if err = json.Unmarshal([]byte(m.Config), dat); err == nil {
			m.ShowEntry = dat.ShowEntry
			m.ShowAdmin = dat.ShowAdmin
		}
	}
	return
}

// PaginateConfig paginate config list of records indexing from start(offset) to end(offset+count) by conditions.
func (d *Dao) PaginateConfig(c context.Context, typ, category int32, oid int64, operator string, offset, count int) (configs []*model.Config, err error) {
	var paginationSQLBuffer bytes.Buffer
	paginationSQLBuffer.WriteString(_paginationConfigSQL)
	sqlwhere, queryParams := d.constructPaginationSQLWhere(typ, category, oid, operator)
	paginationSQLBuffer.WriteString(sqlwhere)
	paginationSQLBuffer.WriteString(" limit ?, ?")
	queryParams = append(queryParams, offset, count-1)
	rows, err := d.db.Query(c, paginationSQLBuffer.String(), queryParams...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := new(model.Config)
		if err = rows.Scan(&m.ID, &m.Type, &m.Oid, &m.AdminID, &m.Operator, &m.Category, &m.Config, &m.CTime, &m.MTime); err != nil {
			return
		}
		if m.ID > 0 {
			dat := new(model.Config)
			if err = json.Unmarshal([]byte(m.Config), dat); err == nil {
				m.ShowEntry = dat.ShowEntry
				m.ShowAdmin = dat.ShowAdmin
			}
		}
		configs = append(configs, m)
	}
	return
}

// PaginateConfigCount returns a total count of records by conditions.
func (d *Dao) PaginateConfigCount(c context.Context, typ, category int32, oid int64, operator string) (totalCount int64, err error) {
	var paginationCountSQLBuffer bytes.Buffer
	paginationCountSQLBuffer.WriteString(_paginationConfigCountSQL)
	sqlwhere, countParams := d.constructPaginationSQLWhere(typ, category, oid, operator)
	paginationCountSQLBuffer.WriteString(sqlwhere)
	row := d.db.QueryRow(c, paginationCountSQLBuffer.String(), countParams...)
	if err = row.Scan(&totalCount); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

func (d *Dao) constructPaginationSQLWhere(tp, category int32, oid int64, operator string) (sqlWhere string, queryParams []interface{}) {
	var sqlBuffer bytes.Buffer
	sqlBuffer.WriteString(" where 1 = 1 ")
	if tp > 0 && oid > 0 && category > 0 {
		sqlBuffer.WriteString(" and type = ? and oid = ? and category = ?")
		queryParams = append(queryParams, tp, oid, category)
	}
	if len(operator) > 0 {
		sqlBuffer.WriteString(" and operator = ?")
		queryParams = append(queryParams, operator)
	}
	return sqlBuffer.String(), queryParams
}

// DeleteConfig delete a reply config record by id.
func (d *Dao) DeleteConfig(c context.Context, id int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _deleteConfigSQL, id)
	if err != nil {
		return
	}
	return res.RowsAffected()
}
