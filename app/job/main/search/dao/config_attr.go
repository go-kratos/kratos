package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/search/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getAttrsSQL = "SELECT appid,db_name,es_name,table_prefix,table_format,index_prefix,index_version,index_format,index_type,index_id,index_mapping, " +
		"data_index_suffix,review_num,review_time,sleep,size,business,data_fields,data_extra,sql_by_id,sql_by_mtime,sql_by_idmtime,databus_info,databus_index_id FROM digger_app WHERE appid=?"
)

type attr struct {
	d     *Dao
	appID string
	attrs *model.Attrs
}

func newAttr(d *Dao, appID string) (ar *attr) {
	ar = &attr{
		d:     d,
		appID: appID,
		attrs: new(model.Attrs),
	}
	if err := ar.initAttrs(); err != nil {
		//fmt.Println("strace:init>", err)
		log.Error("d.initAttrs error (%v)", err)
	}
	return
}

func (ar *attr) initAttrs() (err error) {
	var sqlAttrs *model.SQLAttrs
	for {
		if sqlAttrs, err = ar.getSQLAttrs(context.TODO()); err != nil || sqlAttrs == nil {
			log.Error("d.Attrs error (%v)", err)
			time.Sleep(time.Second * 3)
			continue
		}
		break
	}
	// attr-src
	ar.attrs.Business = sqlAttrs.Business
	ar.attrs.AppID = sqlAttrs.AppID
	ar.attrs.DBName = sqlAttrs.DBName
	ar.attrs.ESName = sqlAttrs.ESName
	ar.attrs.DtbName = sqlAttrs.DtbName
	// attr-table
	if err = ar.parseTable(sqlAttrs); err != nil {
		err = fmt.Errorf("parseTable appid(%s) err(%v)", ar.appID, err)
		return
	}
	// attr-index
	if err = ar.parseIndex(sqlAttrs); err != nil {
		err = fmt.Errorf("parseIndex appid(%s) err(%v)", ar.appID, err)
		return
	}
	// attr-datasql
	if err = ar.parseDataSQL(sqlAttrs); err != nil {
		err = fmt.Errorf("parseDataSQL appid(%s) err(%v)", ar.appID, err)
		return
	}
	// attr-sql

	// attr-data_extra
	if err = ar.parseExtraData(sqlAttrs); err != nil {
		err = fmt.Errorf("parseExtraData appid(%s) err(%v)", ar.appID, err)
		return
	}
	// attr-databus
	if err = ar.parseDatabus(sqlAttrs); err != nil {
		err = fmt.Errorf("parseDatabus appid(%s) err(%v)", ar.appID, err)
		return
	}
	// attr-other
	ar.attrs.Other = &model.AttrOther{
		ReviewNum:  sqlAttrs.ReviewNum,
		ReviewTime: sqlAttrs.ReviewTime,
		Sleep:      sqlAttrs.Sleep,
		Size:       sqlAttrs.Size,
	}
	return
}

func (ar *attr) getSQLAttrs(c context.Context) (res *model.SQLAttrs, err error) {
	res = new(model.SQLAttrs)
	row := ar.d.SearchDB.QueryRow(c, _getAttrsSQL, ar.appID)
	//fmt.Println("appID", ar.appID)
	if err = row.Scan(&res.AppID, &res.DBName, &res.ESName, &res.TablePrefix, &res.TableFormat, &res.IndexAliasPrefix, &res.IndexVersion, &res.IndexFormat, &res.IndexType, &res.IndexID, &res.IndexMapping,
		&res.DataIndexSuffix, &res.ReviewNum, &res.ReviewTime, &res.Sleep, &res.Size, &res.Business, &res.DataFields, &res.DataExtraInfo, &res.SQLByID, &res.SQLByMTime, &res.SQLByIDMTime, &res.DatabusInfo, &res.DatabusIndexID); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res = nil
		}
	}
	return
}

func (ar *attr) parseTable(sqlAttrs *model.SQLAttrs) (err error) {
	table := new(model.AttrTable)
	table.TablePrefix = sqlAttrs.TablePrefix
	table.TableFormat = sqlAttrs.TableFormat
	tableFormat := strings.Split(table.TableFormat, ",")
	if len(tableFormat) != 5 {
		err = fmt.Errorf("wrong tableForamt(%s)", tableFormat)
		return
	}
	if table.TableSplit = tableFormat[0]; table.TableSplit != "single" {
		if table.TableFrom, err = strconv.Atoi(tableFormat[1]); err != nil {
			return
		}
		if table.TableTo, err = strconv.Atoi(tableFormat[2]); err != nil {
			return
		}
	}
	table.TableZero = tableFormat[3]
	table.TableFixed = (tableFormat[4] == "fixed")
	ar.attrs.Table = table
	return
}

func (ar *attr) parseIndex(sqlAttrs *model.SQLAttrs) (err error) {
	index := new(model.AttrIndex)
	index.IndexAliasPrefix = sqlAttrs.IndexAliasPrefix
	index.IndexEntityPrefix = sqlAttrs.IndexAliasPrefix + sqlAttrs.IndexVersion
	index.IndexFormat = sqlAttrs.IndexFormat
	index.IndexType = sqlAttrs.IndexType
	index.IndexID = sqlAttrs.IndexID
	index.IndexMapping = sqlAttrs.IndexMapping
	indexFormat := strings.Split(index.IndexFormat, ",")
	if len(indexFormat) != 5 {
		err = fmt.Errorf("wrong indexFormat(%s)", indexFormat)
		return
	}
	if index.IndexID == "base" {
		err = fmt.Errorf("indexID Prohibition 'base' (%s)", indexFormat)
		return
	}
	if index.IndexSplit = indexFormat[0]; index.IndexSplit != "single" {
		if index.IndexFrom, err = strconv.Atoi(indexFormat[1]); err != nil {
			return
		}
		if index.IndexTo, err = strconv.Atoi(indexFormat[2]); err != nil {
			return
		}
	}
	index.IndexZero = indexFormat[3]
	index.IndexFixed = (indexFormat[4] == "fixed")
	ar.attrs.Index = index
	return
}

func (ar *attr) parseDataSQL(sqlAttrs *model.SQLAttrs) (err error) {
	dataSQL := new(model.AttrDataSQL)
	dataSQL.DataIndexFormatFields = make(map[string]string)
	dataSQL.DataDtbFields = make(map[string][]string)
	dataSQL.DataFieldsV2 = make(map[string]model.AttrDataFields)
	dataSQL.DataIndexSuffix = sqlAttrs.DataIndexSuffix
	dataSQL.DataFields = sqlAttrs.DataFields
	dataSQL.DataExtraInfo = sqlAttrs.DataExtraInfo
	if dataSQL.DataFields == "" {
		return
	}
	p := []model.AttrDataFields{} //DataFieldsV2
	sqlFields := []string{}
	if e := json.Unmarshal([]byte(dataSQL.DataFields), &p); e != nil {
		fields := strings.Split(dataSQL.DataFields, ",")
		for _, v := range fields {
			exp := strings.Split(v, ":")
			indexFieldName := exp[0]
			dataSQL.DataIndexFields = append(dataSQL.DataIndexFields, indexFieldName)
			sqlFields = append(sqlFields, exp[1])
			dataSQL.DataIndexFormatFields[indexFieldName] = exp[2]
			if exp[3] == "n" {
				dataSQL.DataIndexRemoveFields = append(dataSQL.DataIndexRemoveFields, indexFieldName)
			}
		}
	} else {
		// json方式
		for _, v := range p {
			dataSQL.DataFieldsV2[v.ESField] = v
			dataSQL.DataIndexFields = append(dataSQL.DataIndexFields, v.ESField)
			sqlFields = append(sqlFields, v.SQL)
			dataSQL.DataIndexFormatFields[v.ESField] = v.Expect
			if v.Stored == "n" {
				dataSQL.DataIndexRemoveFields = append(dataSQL.DataIndexRemoveFields, v.ESField)
			}
			if v.InDtb == "y" {
				dataSQL.DataDtbFields[v.Field] = append(dataSQL.DataDtbFields[v.Field], v.ESField)
			}
		}
	}
	//fmt.Println(dataSQL.DataDtbFields)
	//sqlFields顺序和attr.DataIndexFields要一致
	if (len(sqlFields) != len(dataSQL.DataIndexFields)) && (len(sqlFields) == 0 || len(dataSQL.DataIndexFields) == 0) {
		log.Error("sqlFields and attr.DataIndexFields are different")
		return
	}
	dataSQL.SQLFields = strings.Join(sqlFields, ",")
	if ar.attrs.Table.TableSplit == "single" {
		dataSQL.SQLByID = fmt.Sprintf(sqlAttrs.SQLByID, dataSQL.SQLFields)
		dataSQL.SQLByMTime = fmt.Sprintf(sqlAttrs.SQLByMTime, dataSQL.SQLFields)
		dataSQL.SQLByIDMTime = fmt.Sprintf(sqlAttrs.SQLByIDMTime, dataSQL.SQLFields)
	} else {
		dataSQL.SQLByID = sqlAttrs.SQLByID
		dataSQL.SQLByMTime = sqlAttrs.SQLByMTime
		dataSQL.SQLByIDMTime = sqlAttrs.SQLByIDMTime
	}
	ar.attrs.DataSQL = dataSQL
	return
}

func (ar *attr) parseExtraData(sqlAttrs *model.SQLAttrs) (err error) {
	if sqlAttrs.DataExtraInfo != "" {
		err = json.Unmarshal([]byte(sqlAttrs.DataExtraInfo), &ar.attrs.DataExtras)
	}
	// append all format field from extra data
	for _, v := range ar.attrs.DataExtras {
		if v.FieldsStr == "" {
			continue
		}
		fields := strings.Split(v.FieldsStr, ",")
		for _, v := range fields {
			exp := strings.Split(v, ":")
			ar.attrs.DataSQL.DataIndexFormatFields[exp[0]] = exp[2]
		}
	}
	return
}

func (ar *attr) parseDatabus(sqlAttrs *model.SQLAttrs) (err error) {
	dtb := new(model.AttrDatabus)
	if sqlAttrs.DatabusInfo != "" {
		databusInfo := strings.Split(sqlAttrs.DatabusInfo, ",")
		if len(databusInfo) != 3 {
			err = fmt.Errorf("wrong databusInfo(%s)", databusInfo)
			return
		}
		dtb.Databus = databusInfo[0]
		if dtb.AggCount, err = strconv.Atoi(databusInfo[1]); err != nil {
			return
		}
		if dtb.Ticker, err = strconv.Atoi(databusInfo[2]); err != nil {
			return
		}
	}
	if sqlAttrs.DatabusIndexID != "" {
		databusIndexID := strings.Split(sqlAttrs.DatabusIndexID, ":")
		if len(databusIndexID) != 2 {
			err = fmt.Errorf("wrong databusIndexID(%s)", databusIndexID)
			return
		}
		dtb.PrimaryID = databusIndexID[0]
		dtb.RelatedID = databusIndexID[1]
	}
	ar.attrs.Databus = dtb
	return
}
