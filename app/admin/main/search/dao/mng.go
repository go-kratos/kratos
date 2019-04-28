package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go-common/app/admin/main/search/model"
	sqlx "go-common/library/database/sql"
)

const (
	_mngBusinessListSQL       = `select id,business,description,app_ids from digger_business where business like '%%%s%%' limit ?,?`
	_mngBusinessListTotalSQL  = `select count(*) from digger_business where business like '%%%s%%'`
	_mngBusinessAllSQL        = `select id,business,description,app_ids from digger_business`
	_mngAddBusinessSQL        = `insert into digger_business (business,description,app_ids) values (?,?,?)`
	_mngUpdateBusinessSQL     = `update digger_business set business=?,description=?,app_ids=? where id=?`
	_mngBusinessInfoSQL       = `select id,business,description,app_ids from digger_business where id=?`
	_mngBusinessInfoByNameSQL = `select id,business,description,app_ids from digger_business where business=?`

	_mngAssetListSQL       = `select id,name,type,src,description from digger_asset %s limit ?,?`
	_mngAssetTotalSQL      = `select count(*) from digger_asset %s`
	_mngAssetAllSQL        = `select id,name,type,src,description from digger_asset`
	_mngAssetInfoSQL       = `select id,name,type,src,description from digger_asset where id=?`
	_mngAssetInfoByNameSQL = `select id,name,type,src,description from digger_asset where name=?`
	_mngAddAssetSQL        = `insert into digger_asset (name,type,src,description) values (?,?,?,?)`
	_mngUpdateAssetSQL     = `update digger_asset set name=?,type=?,src=?,description=? where id=?`

	_mngApplistSQL = `select id,business,appid,description,db_name,es_name,table_name,databus_name,table_prefix,table_format,index_prefix,
index_version,index_format,index_type,index_id,data_index_suffix,index_mapping,data_fields,data_extra,review_num,review_time,
sleep,size,sql_by_id,sql_by_mtime,sql_by_idmtime,databus_info,databus_index_id,query_max_indexes from digger_app where business=?`
	_mngAppInfoSQL = `select id,business,appid,description,db_name,es_name,table_name,databus_name,table_prefix,table_format,index_prefix,
index_version,index_format,index_type,index_id,data_index_suffix,index_mapping,data_fields,data_extra,review_num,review_time,
sleep,size,sql_by_id,sql_by_mtime,sql_by_idmtime,databus_info,databus_index_id,query_max_indexes from digger_app where id=?`
	_mngAppInfoByAppidSQL = `select id,business,appid,description,db_name,es_name,table_name,databus_name,table_prefix,table_format,index_prefix,
index_version,index_format,index_type,index_id,data_index_suffix,index_mapping,data_fields,data_extra,review_num,review_time,
sleep,size,sql_by_id,sql_by_mtime,sql_by_idmtime,databus_info,databus_index_id,query_max_indexes from digger_app where appid=?`
	_mngAddAppSQL    = `insert into digger_app (business,appid,description) values (?,?,?)`
	_mngUpdateAppSQL = `update digger_app set business=?,appid=?,description=?,db_name=?,es_name=?,table_name=?,databus_name=?,table_prefix=?,table_format=?,index_prefix=?,
index_version=?,index_format=?,index_type=?,index_id=?,data_index_suffix=?,index_mapping=?,data_fields=?,data_extra=?,review_num=?,review_time=?,
sleep=?,size=?,sql_by_id=?,sql_by_mtime=?,sql_by_idmtime=?,databus_info=?,databus_index_id=?,query_max_indexes=? where id=?`
	_mngUpdateAppAssetTableSQL   = `update digger_app set table_prefix=?,table_format=? where table_name=?`
	_mngUpdateAppAssetDatabusSQL = `update digger_app set databus_info=?,databus_index_id=? where databus_name=?`

	_mngCountSQL   = `select time,count from digger_count where business=? and type=? and time >= ?`
	_mngPercentSQL = `select name,count from digger_count where business=? and type=? and time = ?`
)

// BusinessList .
func (d *Dao) BusinessList(ctx context.Context, name string, offset, limit int) (list []*model.MngBusiness, err error) {
	sqlStr := fmt.Sprintf(_mngBusinessListSQL, name)
	rows, err := d.db.Query(ctx, sqlStr, offset, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.MngBusiness{}
		if err = rows.Scan(&b.ID, &b.Name, &b.Desc, &b.AppsJSON); err != nil {
			return
		}
		b.Apps = make([]*model.MngBusinessApp, 0)
		if b.AppsJSON != "" {
			if err = json.Unmarshal([]byte(b.AppsJSON), &b.Apps); err != nil {
				return
			}
		}
		list = append(list, b)
	}
	err = rows.Err()
	return
}

// BusinessTotal .
func (d *Dao) BusinessTotal(ctx context.Context, name string) (total int64, err error) {
	sqlStr := fmt.Sprintf(_mngBusinessListTotalSQL, name)
	err = d.db.QueryRow(ctx, sqlStr).Scan(&total)
	return
}

// BusinessAll .
func (d *Dao) BusinessAll(ctx context.Context) (list []*model.MngBusiness, err error) {
	rows, err := d.db.Query(ctx, _mngBusinessAllSQL)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.MngBusiness{}
		if err = rows.Scan(&b.ID, &b.Name, &b.Desc, &b.AppsJSON); err != nil {
			return
		}
		b.Apps = make([]*model.MngBusinessApp, 0)
		if b.AppsJSON != "" {
			if err = json.Unmarshal([]byte(b.AppsJSON), &b.Apps); err != nil {
				return
			}
		}
		list = append(list, b)
	}
	err = rows.Err()
	return
}

// AddBusiness .
func (d *Dao) AddBusiness(ctx context.Context, b *model.MngBusiness) (id int64, err error) {
	res, err := d.db.Exec(ctx, _mngAddBusinessSQL, b.Name, b.Desc, b.AppsJSON)
	if err != nil {
		return
	}
	id, err = res.LastInsertId()
	return
}

// UpdateBusiness .
func (d *Dao) UpdateBusiness(ctx context.Context, b *model.MngBusiness) (err error) {
	_, err = d.db.Exec(ctx, _mngUpdateBusinessSQL, b.Name, b.Desc, b.AppsJSON, b.ID)
	return
}

// BusinessInfo .
func (d *Dao) BusinessInfo(ctx context.Context, id int64) (info *model.MngBusiness, err error) {
	info = new(model.MngBusiness)
	if err = d.db.QueryRow(ctx, _mngBusinessInfoSQL, id).Scan(&info.ID, &info.Name, &info.Desc, &info.AppsJSON); err != nil {
		if err == sqlx.ErrNoRows {
			info = nil
			err = nil
		}
		return
	}
	info.Apps = make([]*model.MngBusinessApp, 0)
	if info.AppsJSON != "" {
		err = json.Unmarshal([]byte(info.AppsJSON), &info.Apps)
	}
	return
}

// BusinessInfoByName .
func (d *Dao) BusinessInfoByName(ctx context.Context, name string) (info *model.MngBusiness, err error) {
	info = new(model.MngBusiness)
	if err = d.db.QueryRow(ctx, _mngBusinessInfoByNameSQL, name).Scan(&info.ID, &info.Name, &info.Desc, &info.AppsJSON); err != nil {
		if err == sqlx.ErrNoRows {
			info = nil
			err = nil
		}
		return
	}
	info.Apps = make([]*model.MngBusinessApp, 0)
	if info.AppsJSON != "" {
		err = json.Unmarshal([]byte(info.AppsJSON), &info.Apps)
	}
	return
}

// AssetList .
func (d *Dao) AssetList(ctx context.Context, typ int, name string, offset, limit int) (list []*model.MngAsset, err error) {
	where := " where 1 "
	if typ > 0 {
		where += fmt.Sprintf(" and type=%d ", typ)
	}
	if name != "" {
		where += fmt.Sprintf(" and name like '%%%s%%'", name)
	}
	sqlStr := fmt.Sprintf(_mngAssetListSQL, where)
	rows, err := d.db.Query(ctx, sqlStr, offset, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.MngAsset{}
		if err = rows.Scan(&a.ID, &a.Name, &a.Type, &a.Config, &a.Desc); err != nil {
			return
		}
		list = append(list, a)
	}
	err = rows.Err()
	return
}

// AssetTotal .
func (d *Dao) AssetTotal(ctx context.Context, typ int, name string) (total int64, err error) {
	where := " where 1 "
	if typ > 0 {
		where += fmt.Sprintf(" and type=%d ", typ)
	}
	if name != "" {
		where += fmt.Sprintf(" and name like '%%%s%%'", name)
	}
	sqlStr := fmt.Sprintf(_mngAssetTotalSQL, where)
	err = d.db.QueryRow(ctx, sqlStr).Scan(&total)
	return
}

// AssetAll .
func (d *Dao) AssetAll(ctx context.Context) (list []*model.MngAsset, err error) {
	rows, err := d.db.Query(ctx, _mngAssetAllSQL)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.MngAsset{}
		if err = rows.Scan(&a.ID, &a.Name, &a.Type, &a.Config, &a.Desc); err != nil {
			return
		}
		list = append(list, a)
	}
	err = rows.Err()
	return
}

// AssetInfo .
func (d *Dao) AssetInfo(ctx context.Context, id int64) (info *model.MngAsset, err error) {
	info = new(model.MngAsset)
	if err = d.db.QueryRow(ctx, _mngAssetInfoSQL, id).Scan(&info.ID, &info.Name, &info.Type, &info.Config, &info.Desc); err != nil {
		if err == sqlx.ErrNoRows {
			info = nil
			err = nil
		}
		return
	}
	return
}

// AssetInfoByName .
func (d *Dao) AssetInfoByName(ctx context.Context, name string) (info *model.MngAsset, err error) {
	info = new(model.MngAsset)
	if err = d.db.QueryRow(ctx, _mngAssetInfoByNameSQL, name).Scan(&info.ID, &info.Name, &info.Type, &info.Config, &info.Desc); err != nil {
		if err == sqlx.ErrNoRows {
			info = nil
			err = nil
		}
		return
	}
	return
}

// AddAsset .
func (d *Dao) AddAsset(ctx context.Context, b *model.MngAsset) (id int64, err error) {
	res, err := d.db.Exec(ctx, _mngAddAssetSQL, b.Name, b.Type, b.Config, b.Desc)
	if err != nil {
		return
	}
	id, err = res.LastInsertId()
	return
}

// UpdateAsset .
func (d *Dao) UpdateAsset(ctx context.Context, b *model.MngAsset) (err error) {
	_, err = d.db.Exec(ctx, _mngUpdateAssetSQL, b.Name, b.Type, b.Config, b.Desc, b.ID)
	return
}

// AppList .
func (d *Dao) AppList(ctx context.Context, business string) (list []*model.MngApp, err error) {
	rows, err := d.db.Query(ctx, _mngApplistSQL, business)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.MngApp{}
		if err = rows.Scan(&a.ID, &a.Business, &a.AppID, &a.Desc, &a.DBName, &a.ESName, &a.TableName, &a.DatabusName, &a.TablePrefix, &a.TableFormat,
			&a.IndexPrefix, &a.IndexVersion, &a.IndexFormat, &a.IndexType, &a.IndexID, &a.DataIndexSuffix, &a.IndexMapping,
			&a.DataFields, &a.DataExtra, &a.ReviewNum, &a.ReviewTime, &a.Sleep, &a.Size, &a.SQLByID, &a.SQLByMtime,
			&a.SQLByIDMtime, &a.DatabusInfo, &a.DatabusIndexID, &a.QueryMaxIndexes); err != nil {
			return
		}
		list = append(list, a)
	}
	err = rows.Err()
	return
}

// AppInfo .
func (d *Dao) AppInfo(ctx context.Context, id int64) (a *model.MngApp, err error) {
	a = new(model.MngApp)
	if err = d.db.QueryRow(ctx, _mngAppInfoSQL, id).Scan(&a.ID, &a.Business, &a.AppID, &a.Desc, &a.DBName, &a.ESName, &a.TableName, &a.DatabusName,
		&a.TablePrefix, &a.TableFormat, &a.IndexPrefix, &a.IndexVersion, &a.IndexFormat, &a.IndexType, &a.IndexID, &a.DataIndexSuffix, &a.IndexMapping,
		&a.DataFields, &a.DataExtra, &a.ReviewNum, &a.ReviewTime, &a.Sleep, &a.Size, &a.SQLByID, &a.SQLByMtime,
		&a.SQLByIDMtime, &a.DatabusInfo, &a.DatabusIndexID, &a.QueryMaxIndexes); err != nil {
		if err == sqlx.ErrNoRows {
			a = nil
			err = nil
		}
		return
	}
	return
}

// AppInfoByAppid .
func (d *Dao) AppInfoByAppid(ctx context.Context, appid string) (a *model.MngApp, err error) {
	a = new(model.MngApp)
	if err = d.db.QueryRow(ctx, _mngAppInfoByAppidSQL, appid).Scan(&a.ID, &a.Business, &a.AppID, &a.Desc, &a.DBName, &a.ESName, &a.TableName, &a.DatabusName,
		&a.TablePrefix, &a.TableFormat, &a.IndexPrefix, &a.IndexVersion, &a.IndexFormat, &a.IndexType, &a.IndexID, &a.DataIndexSuffix, &a.IndexMapping,
		&a.DataFields, &a.DataExtra, &a.ReviewNum, &a.ReviewTime, &a.Sleep, &a.Size, &a.SQLByID, &a.SQLByMtime,
		&a.SQLByIDMtime, &a.DatabusInfo, &a.DatabusIndexID, &a.QueryMaxIndexes); err != nil {
		if err == sqlx.ErrNoRows {
			a = nil
			err = nil
		}
		return
	}
	return
}

// AddApp .
func (d *Dao) AddApp(ctx context.Context, a *model.MngApp) (id int64, err error) {
	res, err := d.db.Exec(ctx, _mngAddAppSQL, a.Business, a.AppID, a.Desc)
	if err != nil {
		return
	}
	id, err = res.LastInsertId()
	return
}

// UpdateApp .
func (d *Dao) UpdateApp(ctx context.Context, a *model.MngApp) (err error) {
	_, err = d.db.Exec(ctx, _mngUpdateAppSQL, a.Business, a.AppID, a.Desc, a.DBName, a.ESName, a.TableName, a.DatabusName, a.TablePrefix, a.TableFormat,
		a.IndexPrefix, a.IndexVersion, a.IndexFormat, a.IndexType, a.IndexID, a.DataIndexSuffix, a.IndexMapping,
		a.DataFields, a.DataExtra, a.ReviewNum, a.ReviewTime, a.Sleep, a.Size, a.SQLByID, a.SQLByMtime,
		a.SQLByIDMtime, a.DatabusInfo, a.DatabusIndexID, a.QueryMaxIndexes, a.ID)
	return
}

// UpdateAppAssetTable .
func (d *Dao) UpdateAppAssetTable(ctx context.Context, name string, t *model.MngAssetTable) (err error) {
	_, err = d.db.Exec(ctx, _mngUpdateAppAssetTableSQL, t.TablePrefix, t.TableFormat, name)
	return
}

// UpdateAppAssetDatabus .
func (d *Dao) UpdateAppAssetDatabus(ctx context.Context, name string, v *model.MngAssetDatabus) (err error) {
	_, err = d.db.Exec(ctx, _mngUpdateAppAssetDatabusSQL, v.DatabusInfo, v.DatabusIndexID, name)
	return
}

// MngCount .
func (d *Dao) MngCount(ctx context.Context, c *model.MngCount) (list []*model.MngCountRes, err error) {
	list = []*model.MngCountRes{}
	sTime := time.Now().AddDate(0, 0, -365).Format("2006-01-02")
	rows, err := d.db.Query(ctx, _mngCountSQL, c.Business, c.Type, sTime)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.MngCountRes{}
		if err = rows.Scan(&a.Time, &a.Count); err != nil {
			return
		}
		a.Time = a.Time[:10]
		list = append(list, a)
	}
	err = rows.Err()
	return
}

// MngPercent .
func (d *Dao) MngPercent(ctx context.Context, c *model.MngCount) (list []*model.MngPercentRes, err error) {
	list = []*model.MngPercentRes{}
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	rows, err := d.db.Query(ctx, _mngPercentSQL, c.Business, c.Type, yesterday)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.MngPercentRes{}
		if err = rows.Scan(&a.Name, &a.Count); err != nil {
			return
		}
		list = append(list, a)
	}
	err = rows.Err()
	return
}

// Unames .
func (d *Dao) Unames(c context.Context, uids []string) (res *model.UnamesData, err error) {
	params := url.Values{}
	params.Set("uids", strings.Join(uids, ","))
	if err = d.client.Get(c, d.managerUnames, "", params, &res); err != nil {
		return
	}
	if res.Code != 0 {
		return
	}
	return
}
