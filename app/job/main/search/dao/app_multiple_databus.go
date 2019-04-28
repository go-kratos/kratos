package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go-common/app/job/main/search/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// AppMultipleDatabus .
type AppMultipleDatabus struct {
	d               *Dao
	appid           string
	attrs           *model.Attrs
	db              *xsql.DB
	dtb             *databus.Databus
	offsets         model.LoopOffsets
	mapData         []model.MapData
	tableName       []string
	indexNameSuffix []string
	commits         map[int32]*databus.Message
}

// IndexNameSuffix .
func (amd *AppMultipleDatabus) IndexNameSuffix(format string, startDate string) (res []string, err error) {
	var (
		sTime time.Time
		eTime = time.Now()
	)
	sTime, err = time.Parse(format, startDate)
	if err != nil {
		log.Error("d.LogAuditIndexName(%v)", startDate)
		return
	}
	resDict := map[string]bool{}
	if strings.Contains(format, "02") {
		for {
			resDict[amd.getIndexName(format, eTime)] = true
			eTime = eTime.AddDate(0, 0, -1)
			if sTime.After(eTime) {
				break
			}
		}
	} else if strings.Contains(format, "week") {
		for {
			resDict[amd.getIndexName(format, eTime)] = true
			eTime = eTime.AddDate(0, 0, -7)
			if sTime.After(eTime) {
				break
			}
		}
	} else if strings.Contains(format, "01") {
		// 1月31日时AddDate(0, -1, 0)会出现错误
		year, month, _ := eTime.Date()
		hour, min, sec := eTime.Clock()
		eTime = time.Date(year, month, 1, hour, min, sec, 0, eTime.Location())
		for {
			resDict[amd.getIndexName(format, eTime)] = true
			eTime = eTime.AddDate(0, -1, 0)
			if sTime.After(eTime) {
				break
			}
		}
	} else if strings.Contains(format, "2006") {
		// 2月29日时AddDate(-1, 0, 0)会出现错误
		year, _, _ := eTime.Date()
		hour, min, sec := eTime.Clock()
		eTime = time.Date(year, 1, 1, hour, min, sec, 0, eTime.Location())
		for {
			resDict[amd.getIndexName(format, eTime)] = true
			eTime = eTime.AddDate(-1, 0, 0)
			if sTime.After(eTime) {
				break
			}
		}
	}
	for k := range resDict {
		res = append(res, k)
	}
	return
}

func (amd *AppMultipleDatabus) getIndexName(format string, time time.Time) (index string) {
	var (
		week = map[int]string{
			0: "0108",
			1: "0916",
			2: "1724",
			3: "2531",
		}
	)
	return strings.Replace(time.Format(format), "week", week[time.Day()/9], -1)
}

// NewAppMultipleDatabus .
func NewAppMultipleDatabus(d *Dao, appid string) (amd *AppMultipleDatabus) {
	var err error
	amd = &AppMultipleDatabus{
		d:               d,
		appid:           appid,
		attrs:           d.AttrPool[appid],
		offsets:         make(map[int]*model.LoopOffset),
		tableName:       []string{},
		indexNameSuffix: []string{},
		commits:         make(map[int32]*databus.Message),
	}
	amd.db = d.DBPool[amd.attrs.DBName]
	amd.dtb = d.DatabusPool[amd.attrs.Databus.Databus]
	if amd.attrs.Table.TableSplit == "int" || amd.attrs.Table.TableSplit == "single" {
		for i := amd.attrs.Table.TableFrom; i <= amd.attrs.Table.TableTo; i++ {
			tableName := fmt.Sprintf("%s%0"+amd.attrs.Table.TableZero+"d", amd.attrs.Table.TablePrefix, i)
			amd.tableName = append(amd.tableName, tableName)
			amd.offsets[i] = &model.LoopOffset{}
		}
	} else {
		var tableNameSuffix []string
		tableFormat := strings.Split(amd.attrs.Table.TableFormat, ",")
		if tableNameSuffix, err = amd.IndexNameSuffix(tableFormat[0], tableFormat[1]); err != nil {
			log.Error("amd.IndexNameSuffix(%v)", err)
			return
		}
		for _, v := range tableNameSuffix {
			amd.tableName = append(amd.tableName, amd.attrs.Table.TablePrefix+v)
		}
		for i := range amd.tableName {
			amd.offsets[i] = &model.LoopOffset{}
		}
	}
	return
}

// Business return business.
func (amd *AppMultipleDatabus) Business() string {
	return amd.attrs.Business
}

// InitIndex .
func (amd *AppMultipleDatabus) InitIndex(c context.Context) {
	var (
		err             error
		indexAliasName  string
		indexEntityName string
	)
	indexFormat := strings.Split(amd.attrs.Index.IndexFormat, ",")
	aliases, aliasErr := amd.d.GetAliases(amd.attrs.ESName, amd.attrs.Index.IndexAliasPrefix)
	if indexFormat[0] == "int" || indexFormat[0] == "single" {
		for i := amd.attrs.Index.IndexFrom; i <= amd.attrs.Index.IndexTo; i++ {
			// == "0" 有问题，不通用
			if amd.attrs.Index.IndexZero == "0" {
				indexAliasName = amd.attrs.Index.IndexAliasPrefix
				indexEntityName = amd.attrs.Index.IndexEntityPrefix
			} else {
				indexAliasName = fmt.Sprintf("%s%0"+amd.attrs.Index.IndexZero+"d", amd.attrs.Index.IndexAliasPrefix, i)
				indexEntityName = fmt.Sprintf("%s%0"+amd.attrs.Index.IndexZero+"d", amd.attrs.Index.IndexEntityPrefix, i)
			}
			if aliasErr != nil {
				amd.d.InitIndex(c, nil, amd.attrs.ESName, indexAliasName, indexEntityName, amd.attrs.Index.IndexMapping)
			} else {
				amd.d.InitIndex(c, aliases, amd.attrs.ESName, indexAliasName, indexEntityName, amd.attrs.Index.IndexMapping)
			}
		}
	} else {
		if amd.indexNameSuffix, err = amd.IndexNameSuffix(indexFormat[0], indexFormat[1]); err != nil {
			log.Error("amd.IndexNameSuffix(%v)", err)
			return
		}
		for _, v := range amd.indexNameSuffix {
			if aliasErr != nil {
				amd.d.InitIndex(c, nil, amd.attrs.ESName, amd.attrs.Index.IndexAliasPrefix+v, amd.attrs.Index.IndexEntityPrefix+v, amd.attrs.Index.IndexMapping)
			} else {
				amd.d.InitIndex(c, aliases, amd.attrs.ESName, amd.attrs.Index.IndexAliasPrefix+v, amd.attrs.Index.IndexEntityPrefix+v, amd.attrs.Index.IndexMapping)
			}
		}
	}
}

// InitOffset insert init value to offset.
func (amd *AppMultipleDatabus) InitOffset(c context.Context) {
	amd.d.InitOffset(c, amd.offsets[0], amd.attrs, amd.tableName)
}

// Offset .
func (amd *AppMultipleDatabus) Offset(c context.Context) {
	for i, v := range amd.tableName {
		offset, err := amd.d.Offset(c, amd.attrs.AppID, v)
		if err != nil {
			log.Error("amd.d.offset error(%v)", err)
			time.Sleep(time.Second * 3)
		}
		amd.offsets[i].SetReview(offset.ReviewID, offset.ReviewTime)
		amd.offsets[i].SetOffset(offset.OffsetID(), offset.OffsetTime())
	}
}

// SetRecover set recover
func (amd *AppMultipleDatabus) SetRecover(c context.Context, recoverID int64, recoverTime string, i int) {
	amd.offsets.SetRecoverOffsets(i, recoverID, recoverTime)
}

// IncrMessages .
func (amd *AppMultipleDatabus) IncrMessages(c context.Context) (length int, err error) {
	ticker := time.NewTicker(time.Duration(time.Millisecond * time.Duration(amd.attrs.Databus.Ticker)))
	defer ticker.Stop()
	for {
		select {
		case msg, ok := <-amd.dtb.Messages():
			if !ok {
				log.Error("databus: %s binlog consumer exit!!!", amd.attrs.Databus)
				break
			}
			m := &model.Message{}
			amd.commits[msg.Partition] = msg
			if err = json.Unmarshal(msg.Value, m); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
				continue
			}
			if amd.attrs.Business == "creative_reply" {
				r, _ := regexp.Compile("reply_\\d+")
				if !r.MatchString(m.Table) {
					continue
				}
			}
			if (amd.attrs.Table.TableSplit == "string" && m.Table == amd.attrs.Table.TablePrefix) ||
				(amd.attrs.Table.TableSplit != "string" && strings.HasPrefix(m.Table, amd.attrs.Table.TablePrefix)) {
				if m.Action == "insert" || m.Action == "update" {
					var parseMap map[string]interface{}
					parseMap, err = amd.d.JSON2map(m.New)
					if err != nil {
						log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
						continue
					}
					// esports fav type filter
					if amd.attrs.AppID == "esports_fav" {
						if t, ok := parseMap["type"]; ok && t.(int64) != 10 {
							continue
						}
					}
					// playlist fav type and attr filter
					if amd.attrs.AppID == "fav_playlist" {
						if t, ok := parseMap["type"]; ok && t.(int64) != 2 {
							continue
						}
						if t, ok := parseMap["attr"]; ok {
							if t.(int64)>>0&1 == 0 || (m.Action == "insert" && t.(int64)>>1&1 == 1) {
								continue
							}
						}
					}
					var newParseMap map[string]interface{}
					newParseMap, err = amd.newParseMap(c, m.Table, parseMap)
					if err != nil {
						if amd.attrs.AppID == "creative_reply" {
							continue
						}
						log.Error("amd.newParseMap error(%v)", err)
						continue
					}
					amd.mapData = append(amd.mapData, newParseMap)
				}
			}
			if len(amd.mapData) < amd.attrs.Databus.AggCount {
				continue
			}
		case <-ticker.C:
		}
		break
	}
	if len(amd.mapData) > 0 {
		amd.mapData, err = amd.d.ExtraData(c, amd.mapData, amd.attrs, "dtb", []string{})
	}
	length = len(amd.mapData)
	//amd.d.extraData(c, amd, "dtb")
	return
}

// AllMessages .
func (amd *AppMultipleDatabus) AllMessages(c context.Context) (length int, err error) {
	amd.mapData = []model.MapData{}
	for i, v := range amd.tableName {
		var (
			rows *xsql.Rows
			sql  string
		)
		tableFormat := strings.Split(amd.attrs.Table.TableFormat, ",")
		if amd.attrs.AppID == "dm_search" || amd.attrs.AppID == "dm" {
			sql = fmt.Sprintf(amd.attrs.DataSQL.SQLByID, amd.attrs.DataSQL.SQLFields, i, i)
		} else if tableFormat[0] == "int" || tableFormat[0] == "single" { // 兼容只传后缀，不传表名
			sql = fmt.Sprintf(amd.attrs.DataSQL.SQLByID, amd.attrs.DataSQL.SQLFields, i)
			log.Info(sql, amd.offsets[i].OffsetID, amd.attrs.Other.Size)
		} else {
			sql = fmt.Sprintf(amd.attrs.DataSQL.SQLByID, amd.attrs.DataSQL.SQLFields, v)
		}
		if rows, err = amd.db.Query(c, sql, amd.offsets[i].OffsetID, amd.attrs.Other.Size); err != nil {
			log.Error("AllMessages db.Query error(%v)", err)
			return
		}
		tempList := []model.MapData{}
		for rows.Next() {
			item, row := InitMapData(amd.attrs.DataSQL.DataIndexFields)
			if err = rows.Scan(row...); err != nil {
				log.Error("AppMultipleDatabus.AllMessages rows.Scan() error(%v)", err)
				continue
			}
			var newParseMap map[string]interface{}
			newParseMap, err = amd.newParseMap(c, v, item)
			if err != nil {
				log.Error("amd.newParseMap error(%v)", err)
				continue
			}
			tempList = append(tempList, newParseMap)
			amd.mapData = append(amd.mapData, newParseMap)
		}
		rows.Close()
		tmpLength := len(tempList)
		if tmpLength > 0 {
			amd.offsets[i].SetTempOffset(tempList[tmpLength-1].PrimaryID(), tempList[tmpLength-1].StrMTime())
		}
	}
	if len(amd.mapData) > 0 {
		amd.mapData, err = amd.d.ExtraData(c, amd.mapData, amd.attrs, "db", []string{})
	}
	length = len(amd.mapData)
	//amd.d.extraData(c, amd, "db")
	return
}

// BulkIndex .
func (amd *AppMultipleDatabus) BulkIndex(c context.Context, start int, end int, writeEntityIndex bool) (err error) {
	partData := amd.mapData[start:end]
	if amd.d.c.Business.Index {
		err = amd.d.BulkDBData(c, amd.attrs, writeEntityIndex, partData...)
	} else {
		err = amd.d.BulkDatabusData(c, amd.attrs, writeEntityIndex, partData...)
	}
	return
}

// Commit .
func (amd *AppMultipleDatabus) Commit(c context.Context) (err error) {
	if amd.d.c.Business.Index {
		if amd.attrs.Table.TableSplit == "int" || amd.attrs.Table.TableSplit == "single" { // 兼容只传后缀，不传表名
			for i := amd.attrs.Table.TableFrom; i <= amd.attrs.Table.TableTo; i++ {
				tableName := fmt.Sprintf("%s%0"+amd.attrs.Table.TableZero+"d", amd.attrs.Table.TablePrefix, i)
				if err = amd.d.CommitOffset(c, amd.offsets[i], amd.attrs.AppID, tableName); err != nil {
					log.Error("AppMultipleDatabus.Commit error(%v)", err)
					continue
				}
			}
		} else {
			for i, v := range amd.indexNameSuffix {
				if err = amd.d.CommitOffset(c, amd.offsets[i], amd.attrs.AppID, v); err != nil {
					log.Error("Commit error(%v)", err)
					continue
				}
			}
		}
	} else {
		for k, c := range amd.commits {
			if err = c.Commit(); err != nil {
				log.Error("AppMultipleDatabus.Commit error(%v)", err)
				continue
			}
			delete(amd.commits, k)
		}
	}
	amd.mapData = []model.MapData{}
	return
}

// Sleep .
func (amd *AppMultipleDatabus) Sleep(c context.Context) {
	time.Sleep(time.Second * time.Duration(amd.attrs.Other.Sleep))
}

// Size .
func (amd *AppMultipleDatabus) Size(c context.Context) (size int) {
	return amd.attrs.Other.Size
}

// indexField .
// func (amd *AppMultipleDatabus) indexField(c context.Context, tableName string) (fieldName string, fieldValue int) {
// 	suffix, _ := strconv.Atoi(strings.Split(tableName, "_")[2])
// 	s := strings.Split(amd.attrs.DataSQL.DataIndexSuffix, ";")
// 	v := strings.Split(s[1], ":")
// 	fieldName = v[0]
// 	indexNum, _ := strconv.Atoi(v[2])
// 	fieldValue = suffix + indexNum
// 	return
// }

// newParseMap .
func (amd *AppMultipleDatabus) newParseMap(c context.Context, table string, parseMap map[string]interface{}) (res map[string]interface{}, err error) {
	res = parseMap
	//TODO 实体索引写不进去
	if (amd.attrs.AppID == "dm_search" || amd.attrs.AppID == "dm") && !amd.d.c.Business.Index {
		indexSuffix := strings.Split(table, "_")[2]
		res["index_name"] = amd.attrs.Index.IndexAliasPrefix + indexSuffix
		if _, ok := res["msg"]; ok {
			// dm_content_
			res["index_field"] = true // 删除ctime
			res["index_id"] = fmt.Sprintf("%v", res["dmid"])
		} else {
			// dm_index_
			res["index_id"] = fmt.Sprintf("%v", res["id"])
		}
	} else if amd.attrs.AppID == "dmreport" {
		if ztime, ok := res["ctime"].(*interface{}); ok { // 数据库
			if ctime, cok := (*ztime).(time.Time); cok {
				res["index_name"] = amd.attrs.Index.IndexAliasPrefix + ctime.Format("2006")
			}
		} else if ztime, ok := res["ctime"].(string); ok { // databus
			var ctime time.Time
			if ctime, err = time.Parse("2006-01-02 15:04:05", ztime); err == nil {
				res["index_name"] = amd.attrs.Index.IndexAliasPrefix + ctime.Format("2006")
			}
		}
	} else if amd.attrs.AppID == "creative_reply" && !amd.d.c.Business.Index {
		if replyType, ok := res["type"].(int64); ok {
			if replyType == 1 || replyType == 12 || replyType == 14 {
			} else {
				err = fmt.Errorf("多余数据")
			}
		} else {
			err = fmt.Errorf("错误数据")
		}
	} else if amd.attrs.Index.IndexSplit == "single" {
		res["index_name"] = amd.attrs.Index.IndexAliasPrefix
	} else {
		indexSuffix := string([]rune(table)[strings.Count(amd.attrs.Table.TablePrefix, "")-1:])
		res["index_name"] = amd.attrs.Index.IndexAliasPrefix + indexSuffix
	}
	//dtb index_id
	if amd.attrs.AppID == "favorite" && !amd.d.c.Business.Index {
		if fid, ok := res["fid"].(int64); ok {
			if oid, ok := res["oid"].(int64); ok {
				res["index_id"] = fmt.Sprintf("%d_%d", fid, oid)
				return
			}
		}
		res["index_id"] = "err"
		res["indexName"] = ""
	}
	return
}
