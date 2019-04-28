package dao

import (
	"context"
	"fmt"
	"go-common/app/service/live/xlottery/model"
	"strings"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getExtraDataByTypeMysql = "select type, item_value, item_extra, ctime, mtime from capsule_extra_data where uid = ? and type in (%v)"
	_addExtraDataMysql       = "insert into capsule_extra_data(uid,type,item_value,item_extra) values "
	_updateExtraValueMysql   = "update capsule_extra_data set item_value = ? where uid = ? and type = ?"
)

// GetExtraDataByType 添加数据
func (d *Dao) GetExtraDataByType(ctx context.Context, uid int64, types []string) (extraData []*model.ExtraData, err error) {
	var rows *sql.Rows
	var typeString []string
	for _, etype := range types {
		typeString = append(typeString, "'"+etype+"'")
	}
	sqlStr := fmt.Sprintf(_getExtraDataByTypeMysql, strings.Join(typeString, ","))
	if rows, err = d.db.Query(ctx, sqlStr, uid); err != nil {
		log.Error("[dao.extra | GetExtraDataByType] query(%s) error (%v)", sqlStr, err)
		return
	}
	defer rows.Close()
	extraData = make([]*model.ExtraData, 0)
	for rows.Next() {
		p := &model.ExtraData{}
		if err = rows.Scan(&p.Type, &p.ItemValue, &p.ItemExtra, &p.Ctime, &p.Mtime); err != nil {
			log.Error("[dao.extra | GetExtraDataByType] scan error, err %v", err)
			return
		}
		extraData = append(extraData, p)
	}
	return
}

// AddExtraData 添加数据
func (d *Dao) AddExtraData(ctx context.Context, extraData []*model.ExtraData) (status bool, err error) {
	sqlStr := _addExtraDataMysql + strings.Repeat("(?,?,?,?),", len(extraData))
	slen := len(sqlStr) - 1
	sqlStr = sqlStr[0:slen]
	values := make([]interface{}, len(extraData)*4)
	for i, extra := range extraData {
		values[i*4] = extra.Uid
		values[i*4+1] = extra.Type
		values[i*4+2] = extra.ItemValue
		values[i*4+3] = extra.ItemExtra
	}
	res, err := d.db.Exec(ctx, sqlStr, values...)
	if err != nil {
		log.Error("[dao.extra | AddExtraData] add(%s) error (%v)", sqlStr, err)
		return
	}
	rows, _ := res.RowsAffected()
	status = rows > 0
	return
}

// UpdateExtraValue 更新数据
func (d *Dao) UpdateExtraValue(ctx context.Context, uid int64, etype string, itemValue int64) (status bool, err error) {
	res, err := d.db.Exec(ctx, _updateExtraValueMysql, itemValue, uid, etype)
	if err != nil {
		log.Error("[dao.extra | UpdateExtraValue] update(%s) error(%v)", _updateExtraValueMysql, err)
		return false, err
	}
	var rows int64
	rows, err = res.RowsAffected()
	if err != nil {
		log.Error("[dao.extra | UpdateExtraValue]  err %v", err)
		return false, err
	}
	return rows > 0, nil
}
