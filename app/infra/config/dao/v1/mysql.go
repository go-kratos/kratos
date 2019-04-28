package v1

import (
	"bytes"
	"context"

	"go-common/app/infra/config/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_getToken        = "SELECT token FROM service_name WHERE name=? AND environment=?"
	_getBuildVersion = "SELECT b.config_id FROM service_name s, build_version b WHERE b.version=? AND s.name=? And s.environment=? AND s.id=b.service_id AND b.state=2"
	_getVersions     = "SELECT c.id,c.remark FROM service_config c, service_name s WHERE s.name=? AND s.environment=? AND s.id =c.service_id AND c.state=? ORDER BY c.id DESC"
	_getNamespace    = "SELECT id,namespace FROM service_namespace WHERE config_id=?"
	_getValue        = "SELECT config_id,namespace_id,name,config FROM service_config_value WHERE config_id=?"
	_getFile         = "SELECT config FROM service_config_value WHERE config_id=? AND name =?"
	_getBuilds       = "SELECT b.version FROM service_name s ,build_version b WHERE s.name=? AND s.environment=? AND s.id=b.service_id AND b.state=2 ORDER BY b.id DESC"
	_getServiceID    = "SELECT id FROM service_name where name=? AND environment =?"
	_insertVersion   = "INSERT INTO service_config(service_id,state,operator) VALUES (?,?,?)"
	_insertConfigs   = "INSERT INTO service_config_value(config_id,name,config,operator) VALUES "
	_updateConfigs   = "UPDATE service_config_value SET config=?,operator=? WHERE config_id = ? AND name = ?"
	_insertLog       = "INSERT INTO log(username,business,info) VALUES (?,?,?)"
)

// Token return a Secret from mysql.
func (d *Dao) Token(c context.Context, svr, env string) (token string, err error) {
	row := d.db.QueryRow(c, _getToken, svr, env)
	if err = row.Scan(&token); err != nil {
		log.Error("row.Scan error(%v) svrName(%v)", err, svr)
		if err == sql.ErrNoRows {
			err = nil
			return
		}
	}
	return
}

// BuildVersion return service build version from mysql.
func (d *Dao) BuildVersion(c context.Context, svr, bver, env string) (version int64, err error) {
	row := d.db.QueryRow(c, _getBuildVersion, bver, svr, env)
	if err = row.Scan(&version); err != nil {
		if err == sql.ErrNoRows {
			version = model.UnknownVersion
			err = nil
			return
		}
		log.Error("row.Scan error(%v)", err)
	}
	return
}

// Values return values from mysql.
func (d *Dao) Values(c context.Context, ver int64) (rs []*model.NSValue, err error) {
	rows, err := d.db.Query(c, _getValue, ver)
	if err != nil {
		log.Error("db.Query(%d) error(%v)", ver, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r model.NSValue
		if err = rows.Scan(&r.ConfigID, &r.NamespaceID, &r.Name, &r.Config); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		rs = append(rs, &r)
	}
	return
}

// Namespaces return namespaces from mysql
func (d *Dao) Namespaces(c context.Context, ver int64) (rs map[int64]string, err error) {
	rows, err := d.db.Query(c, _getNamespace, ver)
	if err != nil {
		log.Error("db.Query(%d) error(%v)", err)
		return
	}
	rs = make(map[int64]string)
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		if err = rows.Scan(&id, &name); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		rs[id] = name
	}
	return
}

//Versions return versions from mysql
func (d *Dao) Versions(c context.Context, svr, env string, state int8) (rs []*model.ReVer, err error) {
	rows, err := d.db.Query(c, _getVersions, svr, env, state)
	if err != nil {
		log.Error("db.Query(%s) error(%v)", svr, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r model.ReVer
		if err = rows.Scan(&r.Version, &r.Remark); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		rs = append(rs, &r)
	}
	return
}

// Value return values from mysql.
func (d *Dao) Value(c context.Context, fname string, ver int64) (file string, err error) {
	row := d.db.QueryRow(c, _getFile, ver, fname)
	if err = row.Scan(&file); err != nil {
		log.Error("row.Scan error(%v)", err)
		if err == sql.ErrNoRows {
			err = ecode.NothingFound
			return
		}
	}
	return
}

// Builds get service builds.
func (d *Dao) Builds(c context.Context, svr, env string) (rs []string, err error) {
	rows, err := d.db.Query(c, _getBuilds, svr, env)
	if err != nil {
		log.Error("db.Query(%s) error(%v)", svr, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r string
		if err = rows.Scan(&r); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		rs = append(rs, r)
	}
	return
}

// ServiceID get ServiceID.
func (d *Dao) ServiceID(c context.Context, svr, env string) (ID int64, err error) {
	row := d.db.QueryRow(c, _getServiceID, svr, env)
	if err != nil {
		log.Error("db.Query(%s) error(%v)", svr, err)
		return
	}
	if err = row.Scan(&ID); err != nil {
		log.Error("row.Scan error(%v)", err)
		if err == sql.ErrNoRows {
			err = ecode.NothingFound
			return
		}
	}
	return
}

// TxInsertVer insert version.
func (d *Dao) TxInsertVer(tx *sql.Tx, svrID int64, user string) (verID int64, err error) {
	row, err := tx.Exec(_insertVersion, svrID, model.ConfigIng, user)
	if err != nil {
		log.Error("db.insert(%d) error(%v)", svrID, err)
		return
	}
	return row.LastInsertId()
}

// TxInsertValues insert config values.
func (d *Dao) TxInsertValues(c context.Context, tx *sql.Tx, verID int64, user string, data map[string]string) (err error) {
	var (
		buffer   bytes.Buffer
		insertTp string
		stmt     *sql.Stmt
		is       []interface{}
	)
	buffer.WriteString(_insertConfigs)
	insertTp = "(?,?,?,?),"
	for key, val := range data {
		buffer.WriteString(insertTp)
		is = append(is, verID)
		is = append(is, key)
		is = append(is, val)
		is = append(is, user)
	}
	buffer.Truncate(buffer.Len() - 1)
	if stmt, err = tx.Prepare(buffer.String()); err != nil {
		log.Error("d.insert() error(%v)", err)
		return
	}
	_, err = stmt.Exec(c, is...)
	if err != nil {
		log.Error("d.insert() error(%v)", err)
	}
	return
}

// TxUpdateValues update config values.
func (d *Dao) TxUpdateValues(tx *sql.Tx, verID int64, user string, data map[string]string) (err error) {
	for key, val := range data {
		if _, err = tx.Exec(_updateConfigs, val, user, verID, key); err != nil {
			log.Error("db.UpdateValues(%d) error(%v)", user, err)
			break
		}
	}
	if err != nil {
		log.Error("d.insert() error(%v)", err)
	}
	return
}

// InsertLog insert log.
func (d *Dao) InsertLog(c context.Context, user, business, info string) (err error) {
	_, err = d.db.Exec(c, _insertLog, user, business, info)
	if err != nil {
		log.Error("db.InsertLog(%d) error(%v)", user, err)
		return
	}
	return
}
