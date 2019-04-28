package business

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/search/dao"
	"go-common/app/job/main/search/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	minIDSQL = "SELECT id FROM dm_index_%03d WHERE ctime > ? ORDER BY id ASC LIMIT 1"
)

// DmDate .
type DmDate struct {
	d                    *dao.Dao
	appid                string
	attrs                *model.Attrs
	db                   *xsql.DB
	dtb                  *databus.Databus
	offsets              model.LoopOffsets
	mapData              []model.MapData
	commits              map[int32]*databus.Message
	frontTwelveMonthDate string
	tableName            []string
	oidDayMap            map[string]string
}

// NewDmDate .
func NewDmDate(d *dao.Dao, appid string) (dd *DmDate) {
	dd = &DmDate{
		d:                    d,
		appid:                appid,
		attrs:                d.AttrPool[appid],
		offsets:              make(map[int]*model.LoopOffset),
		commits:              make(map[int32]*databus.Message),
		frontTwelveMonthDate: "2017-08-01",
		oidDayMap:            make(map[string]string),
	}
	for i := dd.attrs.Table.TableFrom; i <= dd.attrs.Table.TableTo; i++ {
		dd.offsets[i] = &model.LoopOffset{}
	}
	dd.db = d.DBPool[dd.attrs.DBName]
	dd.dtb = d.DatabusPool[dd.attrs.Databus.Databus]
	return
}

// Business return business.
func (dd *DmDate) Business() string {
	return dd.attrs.Business
}

// InitIndex init index.
func (dd *DmDate) InitIndex(c context.Context) {
	var (
		indexAliasName  string
		indexEntityName string
	)
	aliases, err := dd.d.GetAliases(dd.attrs.ESName, dd.attrs.Index.IndexAliasPrefix)
	now := time.Now()
	for i := -12; i < 18; i++ {
		newDate := now.AddDate(0, i, 0).Format("2006-01")
		indexAliasName = dd.attrs.Index.IndexAliasPrefix + strings.Replace(newDate, "-", "_", -1)
		indexEntityName = dd.attrs.Index.IndexEntityPrefix + strings.Replace(newDate, "-", "_", -1)
		if err != nil {
			dd.d.InitIndex(c, nil, dd.attrs.ESName, indexAliasName, indexEntityName, dd.attrs.Index.IndexMapping)
		} else {
			dd.d.InitIndex(c, aliases, dd.attrs.ESName, indexAliasName, indexEntityName, dd.attrs.Index.IndexMapping)
		}
	}
}

// InitOffset .
func (dd *DmDate) InitOffset(c context.Context) {
	dd.d.InitOffset(c, dd.offsets[0], dd.attrs, dd.tableName)
	log.Info("in InitOffset")
	for i := dd.attrs.Table.TableFrom; i <= dd.attrs.Table.TableTo; i++ {
		var (
			id  int64
			err error
			row *xsql.Row
		)
		row = dd.db.QueryRow(c, fmt.Sprintf(minIDSQL, i), dd.frontTwelveMonthDate)
		if err = row.Scan(&id); err != nil {
			if err == xsql.ErrNoRows {
				log.Info("in ErrNoRows")
				err = nil
			} else {
				log.Info("row.Scan error(%v)", err)
				log.Error("row.Scan error(%v)", err)
				time.Sleep(time.Second * 3)
				continue
			}
		}
		log.Info("here i am %d", i)
		dd.offsets[i] = &model.LoopOffset{}
		dd.offsets[i].OffsetID = id
	}
	log.Info("InitOffset over")
}

// Offset get offset.
func (dd *DmDate) Offset(c context.Context) {
	for i := dd.attrs.Table.TableFrom; i <= dd.attrs.Table.TableTo; i++ {
		tableName := fmt.Sprintf("%s%0"+dd.attrs.Table.TableZero+"d", dd.attrs.Table.TablePrefix, i)
		offset, err := dd.d.Offset(c, dd.attrs.AppID, tableName)
		if err != nil {
			log.Error("dd.d.Offset error(%v)", err)
			time.Sleep(time.Second * 3)
		}
		dd.offsets[i].SetReview(offset.ReviewID, offset.ReviewTime)
		dd.offsets[i].SetOffset(offset.OffsetID(), offset.OffsetTime())
	}
}

// SetRecover set recover
func (dd *DmDate) SetRecover(c context.Context, recoverID int64, recoverTime string, i int) {
}

// IncrMessages .
func (dd *DmDate) IncrMessages(c context.Context) (length int, err error) {
	ticker := time.NewTicker(time.Duration(time.Millisecond * time.Duration(dd.attrs.Databus.Ticker)))
	defer ticker.Stop()
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	tomorrowZeroTimestamp := t.AddDate(0, 0, 1).Unix()
	nowTimestamp := time.Now().Unix()
	if tomorrowZeroTimestamp-nowTimestamp < 180 {
		dd.oidDayMap = nil
		dd.oidDayMap = make(map[string]string)
	}
	for {
		select {
		case msg, ok := <-dd.dtb.Messages():
			if !ok {
				log.Error("databus: %s binlog consumer exit!!!", dd.attrs.Databus)
				break
			}
			m := &model.Message{}
			dd.commits[msg.Partition] = msg
			if err = json.Unmarshal(msg.Value, m); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
				continue
			}
			if m.Action == "insert" && strings.HasPrefix(m.Table, "dm_index") {
				var parseMap map[string]interface{}
				parseMap, err = dd.d.JSON2map(m.New)
				if err != nil {
					log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
					continue
				}
				newParseMap := dd.newDtbParseMap(c, parseMap)
				indexID := newParseMap["index_id"].(string)
				indexName := newParseMap["index_name"].(string)
				if _, exists := dd.oidDayMap[indexID]; exists {
					continue
				}
				dd.oidDayMap[indexID] = indexName
				dd.mapData = append(dd.mapData, newParseMap)
			}
			if len(dd.mapData) < dd.attrs.Databus.AggCount {
				continue
			}
		case <-ticker.C:
		}
		break
	}
	if len(dd.mapData) > 0 {
		dd.mapData, err = dd.d.ExtraData(c, dd.mapData, dd.attrs, "dtb", []string{})
	}
	length = len(dd.mapData)
	//amd.d.extraData(c, amd, "dtb")
	return
}

// AllMessages .
func (dd *DmDate) AllMessages(c context.Context) (length int, err error) {
	dd.mapData = []model.MapData{}
	for i := dd.attrs.Table.TableFrom; i <= dd.attrs.Table.TableTo; i++ {
		var rows *xsql.Rows
		if dd.offsets[i].OffsetID == 0 {
			continue
		}
		if rows, err = dd.db.Query(c, fmt.Sprintf(dd.attrs.DataSQL.SQLByID, dd.attrs.DataSQL.SQLFields, i), dd.offsets[i].OffsetID, dd.attrs.Other.Size); err != nil {
			log.Error("AllMessages db.Query error(%v)", err)
			return
		}
		tempList := []model.MapData{}
		for rows.Next() {
			item, row := dao.InitMapData(dd.attrs.DataSQL.DataIndexFields)
			if err = rows.Scan(row...); err != nil {
				log.Error("appMultipleDatabus.AllMessages rows.Scan() error(%v)", err)
				continue
			}
			newParseMap := dd.newParseMap(c, item)
			ctime, ok := newParseMap["ctime"].(*interface{})
			if ok {
				dbTime := (*ctime).(time.Time)
				dbTimeStr := dbTime.Format("2006-01-02")
				t1, err1 := time.Parse("2006-01-02", dd.frontTwelveMonthDate)
				t2, err2 := time.Parse("2006-01-02", dbTimeStr)
				if err1 != nil || err2 != nil || t1.After(t2) {
					continue
				}
			} else {
				continue
			}
			tempList = append(tempList, newParseMap)
			dd.mapData = append(dd.mapData, newParseMap)
		}
		rows.Close()
		tmpLength := len(tempList)
		if tmpLength > 0 {
			dd.offsets[i].SetTempOffset(tempList[tmpLength-1].PrimaryID(), tempList[tmpLength-1].StrCTime())
		}
	}
	length = len(dd.mapData)
	if length > 0 {
		dd.mapData, err = dd.d.ExtraData(c, dd.mapData, dd.attrs, "db", []string{})
	}
	log.Info("length is %d", length)
	return
}

// BulkIndex .
func (dd *DmDate) BulkIndex(c context.Context, start int, end int, writeEntityIndex bool) (err error) {
	partData := dd.mapData[start:end]
	// if dd.d.GetConfig(c).Business.Index {
	// 	err = dd.d.BulkDBData(c, dd.attrs, partData...)
	// } else {
	// 	err = dd.d.BulkDatabusData(c, dd.attrs, partData...)
	// }
	err = dd.d.BulkDBData(c, dd.attrs, writeEntityIndex, partData...)
	return
}

// Commit commit offset.
func (dd *DmDate) Commit(c context.Context) (err error) {
	if dd.d.GetConfig(c).Business.Index {
		for i := dd.attrs.Table.TableFrom; i <= dd.attrs.Table.TableTo; i++ {
			tOffset := dd.offsets[i]
			if tOffset.TempOffsetID != 0 {
				tOffset.OffsetID = tOffset.TempOffsetID
			}
			if tOffset.TempOffsetTime != "" {
				tOffset.OffsetTime = tOffset.TempOffsetTime
			}
			tableName := fmt.Sprintf("%s%0"+dd.attrs.Table.TableZero+"d", dd.attrs.Table.TablePrefix, i)
			if err = dd.d.CommitOffset(c, tOffset, dd.attrs.AppID, tableName); err != nil {
				log.Error("appMultipleDatabus.Commit error(%v)", err)
				continue
			}
		}
	} else {
		for k, c := range dd.commits {
			if err = c.Commit(); err != nil {
				log.Error("appMultipleDatabus.Commit error(%v)", err)
				continue
			}
			delete(dd.commits, k)
		}
	}
	dd.mapData = []model.MapData{}
	return

}

// Sleep interval duration.
func (dd *DmDate) Sleep(c context.Context) {
}

// Size return size.
func (dd *DmDate) Size(c context.Context) int {
	return 0
}

// newParseMap .
func (dd *DmDate) newParseMap(c context.Context, parseMap map[string]interface{}) (res map[string]interface{}) {
	res = parseMap
	indexName, strID := "", ""
	if res["month"] != nil {
		if month, ok := res["month"].(*interface{}); ok {
			mth := strings.Replace(dd.b2s((*month).([]uint8)), "-", "_", -1)
			indexName = "dm_date_" + mth
		}
	}
	if res["date"] != nil {
		if date, ok := res["date"].(*interface{}); ok {
			dte := strings.Replace(dd.b2s((*date).([]uint8)), "-", "_", -1)
			if oid, ok := res["oid"].(*interface{}); ok {
				strID = strconv.FormatInt((*oid).(int64), 10) + "_" + dte
			}
		}
	}
	res["index_name"] = indexName
	res["index_id"] = strID
	return
}

// newDtbParseMap .
func (dd *DmDate) newDtbParseMap(c context.Context, parseMap map[string]interface{}) (res map[string]interface{}) {
	res = parseMap
	indexName, strID, mth, dte, id := "", "", "", "", ""
	if res["ctime"] != nil {
		if ctime, ok := res["ctime"].(string); ok {
			t, _ := time.Parse("2006-01-02 15:04:05", ctime)
			mth = t.Format("2006-01")
			dte = t.Format("2006-01-02")
			indexName = "dm_date_" + strings.Replace(mth, "-", "_", -1)
		}
	}
	if res["oid"] != nil {
		if oid, ok := res["oid"].(int64); ok {
			strOid := strconv.FormatInt(oid, 10)
			strID = strOid + "_" + strings.Replace(dte, "-", "_", -1)
		}
	}
	if res["id"] != nil {
		if newID, ok := res["id"].(int64); ok {
			id = strconv.Itoa(int(newID))
		}
	}
	for k := range res {
		if k == "id" || k == "oid" {
			continue
		}
		delete(res, k)
	}
	res["index_name"] = indexName
	res["index_id"] = strID
	res["month"] = mth
	res["date"] = dte
	res["id"] = id
	return
}

// bs2 []uint8 to string.
func (dd *DmDate) b2s(bs []uint8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		b[i] = byte(v)
	}
	return string(b)
}
