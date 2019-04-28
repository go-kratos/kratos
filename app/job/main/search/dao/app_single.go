package dao

import (
	"context"
	"time"

	"go-common/app/job/main/search/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

// AppSingle .
type AppSingle struct {
	d       *Dao
	appid   string
	attrs   *model.Attrs
	db      *xsql.DB
	offset  *model.LoopOffset
	mapData []model.MapData
}

// NewAppSingle .
func NewAppSingle(d *Dao, appid string) (as *AppSingle) {
	as = &AppSingle{
		d:       d,
		appid:   appid,
		attrs:   d.AttrPool[appid],
		offset:  &model.LoopOffset{},
		mapData: []model.MapData{},
		db:      d.DBPool[d.AttrPool[appid].DBName],
	}
	return
}

// Business return business.
func (as *AppSingle) Business() string {
	return as.attrs.Business
}

// InitIndex init index.
func (as *AppSingle) InitIndex(c context.Context) {
	if aliases, err := as.d.GetAliases(as.attrs.ESName, as.attrs.Index.IndexAliasPrefix); err != nil {
		as.d.InitIndex(c, nil, as.attrs.ESName, as.attrs.Index.IndexAliasPrefix, as.attrs.Index.IndexEntityPrefix, as.attrs.Index.IndexMapping)
	} else {
		as.d.InitIndex(c, aliases, as.attrs.ESName, as.attrs.Index.IndexAliasPrefix, as.attrs.Index.IndexEntityPrefix, as.attrs.Index.IndexMapping)
	}
}

// InitOffset insert init value to offset.
func (as *AppSingle) InitOffset(c context.Context) {
	as.d.InitOffset(c, as.offset, as.attrs, []string{})
	nowFormat := time.Now().Format("2006-01-02 15:04:05")
	as.offset.SetOffset(0, nowFormat)
}

// Offset get offset.
func (as *AppSingle) Offset(c context.Context) {
	for {
		offset, err := as.d.Offset(c, as.appid, as.attrs.Table.TablePrefix)
		if err != nil {
			log.Error("ac.d.Offset error(%v)", err)
			time.Sleep(time.Second * 3)
			continue
		}
		as.offset.SetReview(offset.ReviewID, offset.ReviewTime)
		as.offset.SetOffset(offset.OffsetID(), offset.OffsetTime())
		break
	}
}

// SetRecover set recover
func (as *AppSingle) SetRecover(c context.Context, recoverID int64, recoverTime string, i int) {
	as.offset.SetRecoverOffset(recoverID, recoverTime)
}

// IncrMessages .
func (as *AppSingle) IncrMessages(c context.Context) (length int, err error) {
	var rows *xsql.Rows
	//fmt.Println("start", as.offset.OffsetTime)
	if !as.offset.IsLoop {
		rows, err = as.db.Query(c, as.attrs.DataSQL.SQLByMTime, as.offset.OffsetTime, as.attrs.Other.Size)
	} else {
		rows, err = as.db.Query(c, as.attrs.DataSQL.SQLByIDMTime, as.offset.OffsetID, as.offset.OffsetTime, as.attrs.Other.Size)
	}
	if err != nil {
		log.Error("db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	tempPartList := []model.MapData{}
	for rows.Next() {
		item, row := InitMapData(as.attrs.DataSQL.DataIndexFields)
		if err = rows.Scan(row...); err != nil {
			log.Error("IncrMessages rows.Scan() error(%v)", err)
			return
		}
		as.mapData = append(as.mapData, item)
		tempPartList = append(tempPartList, item)
	}
	if len(as.mapData) > 0 {
		// extra relevant data
		as.mapData, err = as.d.ExtraData(c, as.mapData, as.attrs, "db", []string{})
		// offset
		UpdateOffsetByMap(as.offset, tempPartList...)
	}
	length = len(as.mapData)
	return
}

// AllMessages .
func (as *AppSingle) AllMessages(c context.Context) (length int, err error) {
	var rows *xsql.Rows
	if rows, err = as.db.Query(c, as.attrs.DataSQL.SQLByID, as.offset.RecoverID, as.attrs.Other.Size); err != nil {
		log.Error("AllMessages db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		item, row := InitMapData(as.attrs.DataSQL.DataIndexFields)
		if err = rows.Scan(row...); err != nil {
			log.Error("AllMessages rows.Scan() error(%v)", err)
			continue
		}
		as.mapData = append(as.mapData, item)
	}
	if len(as.mapData) > 0 {
		// extra relevant data
		as.mapData, err = as.d.ExtraData(c, as.mapData, as.attrs, "db", []string{})
		// offset
		if v, ok := as.mapData[len(as.mapData)-1]["_id"]; ok && v != nil {
			if v2, ok := v.(interface{}); ok {
				as.offset.SetTempOffset((v2).(int64), "")
				as.offset.SetRecoverTempOffset((v2).(int64), "")
			} else {
				log.Error("single.all._id interface error")
			}
		} else {
			log.Error("single.all._id nil error")
		}
	}
	length = len(as.mapData)
	return
}

// BulkIndex .
func (as *AppSingle) BulkIndex(c context.Context, start int, end int, writeEntityIndex bool) (err error) {
	if len(as.mapData) >= (start+1) && len(as.mapData) >= end {
		partData := as.mapData[start:end]
		err = as.d.BulkDBData(c, as.attrs, writeEntityIndex, partData...)
	}
	return
}

// Commit commit offset.
func (as *AppSingle) Commit(c context.Context) (err error) {
	err = as.d.CommitOffset(c, as.offset, as.attrs.AppID, as.attrs.Table.TablePrefix)
	as.mapData = []model.MapData{}
	return
}

// Sleep interval duration.
func (as *AppSingle) Sleep(c context.Context) {
	time.Sleep(time.Second * time.Duration(as.attrs.Other.Sleep))
}

// Size return size.
func (as *AppSingle) Size(c context.Context) int {
	return as.attrs.Other.Size
}
