package service

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/admin/main/search/model"
	"go-common/library/ecode"
)

// BusinessAllV2 .
func (s *Service) BusinessAllV2(c context.Context) (list []*model.GFBusiness, err error) {
	return s.dao.BusinessAllV2(c)
}

// BusinessInfoV2 .
func (s *Service) BusinessInfoV2(c context.Context, name string) (info *model.GFBusiness, err error) {
	return s.dao.BusinessInfoV2(c, name)
}

// BusinessAdd .
func (s *Service) BusinessAdd(c context.Context, pid int64, name, description string) (id int64, err error) {
	return s.dao.BusinessIns(c, pid, name, description)
}

// BusinessUpdate .
func (s *Service) BusinessUpdate(c context.Context, name, filed, value string) (id int64, err error) {
	allowFields := []string{"data_conf", "index_conf", "business_conf", "description", "state"}
	var allow bool
	for _, v := range allowFields {
		if v == filed {
			allow = true
		}
	}
	if !allow {
		err = ecode.AccessDenied
		return
	}
	return s.dao.BusinessUpdate(c, name, filed, value)
}

// AssetDBTables .
func (s *Service) AssetDBTables(c context.Context) (list []*model.GFAsset, err error) {
	return s.dao.AssetDBTables(c)
}

// AssetDBConnect .
func (s *Service) AssetDBConnect(c context.Context, host, port, user, password string) (dbNames []string, err error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4,utf8", user, password, host, port))
	if err != nil {
		return
	}
	defer db.Close()
	rows, err := db.Query("show databases")
	if err != nil {
		return
	}
	defer rows.Close()
	dbNames = make([]string, 0)
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return
		}
		dbNames = append(dbNames, name)
	}
	return dbNames, rows.Err()
}

// AssetDBAdd .
func (s *Service) AssetDBAdd(c context.Context, name, description, host, port, user, password string) (id int64, err error) {
	dbnames, err := s.AssetDBConnect(c, host, port, user, password)
	if err != nil {
		return
	}
	var dbExist bool
	for _, v := range dbnames {
		if v == name {
			dbExist = true
		}
	}
	if !dbExist {
		err = ecode.AccessDenied
		return
	}
	dsn := fmt.Sprintf(model.DBDsnFormat, user, password, host, port, name)
	return s.dao.AssetDBIns(c, name, description, dsn)
}

// AssetTableAdd .
func (s *Service) AssetTableAdd(c context.Context, db, regex, fields, description string) (id int64, err error) {
	name := db + "." + regex
	return s.dao.AssetTableIns(c, name, db, regex, fields, description)
}

// UpdateAssetTable .
func (s *Service) UpdateAssetTable(c context.Context, name, fields string) (id int64, err error) {
	return s.dao.UpdateAssetTable(c, name, fields)
}

// AssetInfoV2 .
func (s *Service) AssetInfoV2(c context.Context, name string) (info *model.GFAsset, err error) {
	return s.dao.Asset(c, name)
}

// AssetShowTables .
func (s *Service) AssetShowTables(c context.Context, dbName string) (tables []string, err error) {
	asset, err := s.dao.Asset(c, dbName)
	if err != nil {
		return
	}
	db, err := sql.Open("mysql", asset.DSN)
	if err != nil {
		return
	}
	defer db.Close()
	rows, err := db.Query("show tables")
	if err != nil {
		return
	}
	defer rows.Close()
	tables = make([]string, 0)
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return
		}
		tables = append(tables, name)
	}
	return tables, rows.Err()
}

// AssetTableFields .
func (s *Service) AssetTableFields(c context.Context, dbName, regex string) (fs []*model.TableField, count int, err error) {
	asset, err := s.dao.Asset(c, dbName)
	if err != nil {
		return
	}
	db, err := sql.Open("mysql", asset.DSN)
	if err != nil {
		return
	}
	defer db.Close()
	regex = fmt.Sprintf("^%s$", regex)
	rows, err := db.Query("SELECT COLUMN_NAME,DATA_TYPE,count(1) FROM information_schema.COLUMNS WHERE table_name REGEXP ? GROUP BY COLUMN_NAME,DATA_TYPE", regex)
	if err != nil {
		return
	}
	defer rows.Close()
	fs = make([]*model.TableField, 0)
	for rows.Next() {
		f := new(model.TableField)
		if err = rows.Scan(&f.Name, &f.Type, &f.Count); err != nil {
			return nil, 0, err
		}
		fs = append(fs, f)
	}
	if err = rows.Err(); err != nil {
		return
	}
	if len(fs) == 0 {
		err = ecode.NothingFound
		return
	}
	for _, f := range fs {
		if fs[0].Count != f.Count {
			err = ecode.NothingFound
			return
		}
		count = f.Count
	}
	row := db.QueryRow("SELECT COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE table_name REGEXP ? AND CONSTRAINT_NAME='PRIMARY' GROUP BY CONSTRAINT_NAME LIMIT 1", regex)
	var primaryCo string
	err = row.Scan(&primaryCo)
	for k, v := range fs {
		if v.Name == primaryCo {
			fs[k].Primary = true
		}
	}
	return fs, count, err
}

// ClusterOwners .
func (s *Service) ClusterOwners() map[string]string {
	clusters := s.c.Es
	res := make(map[string]string)
	res["default"] = "guanhuaxin,daizhichen,libingqi,zhapuyu"
	for name, es := range clusters {
		if es.Owner == "" {
			continue
		}
		if es.Cluster != "" {
			name = es.Cluster
		}
		res[name] = es.Owner
	}
	return res
}
