package dao

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/job/main/search/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

// AppMultiple .
type AppMultiple struct {
	d       *Dao
	appid   string
	attrs   *model.Attrs
	db      *xsql.DB
	offsets model.LoopOffsets
	mapData []model.MapData
}

// NewAppMultiple .
func NewAppMultiple(d *Dao, appid string) (am *AppMultiple) {
	am = &AppMultiple{
		d:       d,
		appid:   appid,
		db:      d.DBPool[d.AttrPool[appid].DBName],
		attrs:   d.AttrPool[appid],
		offsets: make(map[int]*model.LoopOffset),
	}
	for i := am.attrs.Table.TableFrom; i <= am.attrs.Table.TableTo; i++ {
		am.offsets[i] = &model.LoopOffset{}
	}
	return
}

// Business return business
func (am *AppMultiple) Business() string {
	return am.attrs.Business
}

// InitIndex .
func (am *AppMultiple) InitIndex(c context.Context) {
	var (
		indexAliasName  string
		indexEntityName string
	)
	aliases, err := am.d.GetAliases(am.attrs.ESName, am.attrs.Index.IndexAliasPrefix)
	for i := am.attrs.Index.IndexFrom; i <= am.attrs.Index.IndexTo; i++ {
		indexAliasName = fmt.Sprintf("%s%0"+am.attrs.Index.IndexZero+"d", am.attrs.Index.IndexAliasPrefix, i)
		indexEntityName = fmt.Sprintf("%s%0"+am.attrs.Index.IndexZero+"d", am.attrs.Index.IndexEntityPrefix, i)
		if err != nil {
			am.d.InitIndex(c, nil, am.attrs.ESName, indexAliasName, indexEntityName, am.attrs.Index.IndexMapping)
		} else {
			am.d.InitIndex(c, aliases, am.attrs.ESName, indexAliasName, indexEntityName, am.attrs.Index.IndexMapping)
		}
	}
}

// InitOffset insert init value to offset.
func (am *AppMultiple) InitOffset(c context.Context) {
	am.d.InitOffset(c, am.offsets[0], am.attrs, []string{})
}

// Offset .
func (am *AppMultiple) Offset(c context.Context) {
	for i := am.attrs.Table.TableFrom; i < am.attrs.Table.TableTo; i++ {
		offset, err := am.d.Offset(c, am.attrs.AppID, am.attrs.Table.TablePrefix+strconv.Itoa(i))
		if err != nil {
			log.Error("as.d.Offset error(%v)", err)
			time.Sleep(time.Second * 3)
		}
		am.offsets[i].SetReview(offset.ReviewID, offset.ReviewTime)
		am.offsets[i].SetOffset(offset.OffsetID(), offset.OffsetTime())
	}
}

// SetRecover set recover
func (am *AppMultiple) SetRecover(c context.Context, recoverID int64, recoverTime string, i int) {
	am.offsets.SetRecoverOffsets(i, recoverID, recoverTime)
}

// IncrMessages .
func (am *AppMultiple) IncrMessages(c context.Context) (length int, err error) {
	var rows *xsql.Rows
	am.mapData = []model.MapData{}
	for i := am.attrs.Table.TableFrom; i <= am.attrs.Table.TableTo; i++ {
		if !am.offsets[i].IsLoop {
			rows, err = am.db.Query(c, fmt.Sprintf(am.attrs.DataSQL.SQLByMTime, am.attrs.DataSQL.SQLFields, i), am.offsets[i].OffsetTime, am.attrs.Other.Size)
		} else {
			rows, err = am.db.Query(c, fmt.Sprintf(am.attrs.DataSQL.SQLByIDMTime, am.attrs.DataSQL.SQLFields, i), am.offsets[i].OffsetID, am.offsets[i].OffsetTime, am.attrs.Other.Size)
		}
		if err != nil {
			log.Error("db.Query error(%v)", err)
			continue
		}
		tempList := []model.MapData{}
		for rows.Next() {
			item, row := InitMapData(am.attrs.DataSQL.DataIndexFields)
			if err = rows.Scan(row...); err != nil {
				log.Error("rows.Scan() error(%v)", err)
				continue
			}
			tempList = append(tempList, item)
			am.mapData = append(am.mapData, item)
		}
		rows.Close()
		if len(tempList) > 0 {
			UpdateOffsetByMap(am.offsets[i], tempList...)
		}
	}
	if len(am.mapData) > 0 {
		//fmt.Println("before", am.mapData)
		am.mapData, err = am.d.ExtraData(c, am.mapData, am.attrs, "db", []string{})
		//fmt.Println("after", am.mapData)
	}
	length = len(am.mapData)
	return
}

// AllMessages .
func (am *AppMultiple) AllMessages(c context.Context) (length int, err error) {
	am.mapData = []model.MapData{}
	for i := am.attrs.Table.TableFrom; i <= am.attrs.Table.TableTo; i++ {
		var rows *xsql.Rows
		if rows, err = am.db.Query(c, fmt.Sprintf(am.attrs.DataSQL.SQLByID, am.attrs.DataSQL.SQLFields, i), am.offsets[i].OffsetID, am.attrs.Other.Size); err != nil {
			log.Error("AllMessages db.Query error(%v)", err)
			return
		}
		tempList := []model.MapData{}
		for rows.Next() {
			item, row := InitMapData(am.attrs.DataSQL.DataIndexFields)
			if err = rows.Scan(row...); err != nil {
				log.Error("IncrMessages rows.Scan() error(%v)", err)
				continue
			}
			tempList = append(tempList, item)
			am.mapData = append(am.mapData, item)
		}
		rows.Close()
		tmpLength := len(tempList)
		if tmpLength > 0 {
			am.offsets[i].SetTempOffset(tempList[tmpLength-1].PrimaryID(), tempList[tmpLength-1].StrMTime())
		}
	}
	if len(am.mapData) > 0 {
		am.mapData, err = am.d.ExtraData(c, am.mapData, am.attrs, "db", []string{})
	}
	length = len(am.mapData)
	return
}

// BulkIndex .
func (am *AppMultiple) BulkIndex(c context.Context, start int, end int, writeEntityIndex bool) (err error) {
	if len(am.mapData) >= (start+1) && len(am.mapData) >= end {
		partData := am.mapData[start:end]
		err = am.d.BulkDBData(c, am.attrs, writeEntityIndex, partData...)
	}
	return
}

// Commit .
func (am *AppMultiple) Commit(c context.Context) (err error) {
	for i := am.attrs.Table.TableFrom; i <= am.attrs.Table.TableTo; i++ {
		if err = am.d.CommitOffset(c, am.offsets[i], am.attrs.AppID, am.attrs.Table.TablePrefix+strconv.Itoa(i)); err != nil {
			log.Error("Commit error(%v)", err)
			continue
		}
	}
	am.mapData = []model.MapData{}
	return
}

// Sleep .
func (am *AppMultiple) Sleep(c context.Context) {
	time.Sleep(time.Second * time.Duration(am.attrs.Other.Sleep))
}

// Size .
func (am *AppMultiple) Size(c context.Context) (size int) {
	size = am.attrs.Other.Size
	return
}
