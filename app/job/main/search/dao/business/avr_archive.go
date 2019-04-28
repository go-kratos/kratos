package business

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/search/dao"
	"go-common/app/job/main/search/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

// AvrArchive .
type AvrArchive struct {
	d       *dao.Dao
	appid   string
	attrs   *model.Attrs
	db      *xsql.DB
	offset  *model.LoopOffset
	mapData []model.MapData
}

// NewAvrArchive .
func NewAvrArchive(d *dao.Dao, appid string) (av *AvrArchive) {
	av = &AvrArchive{
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
func (av *AvrArchive) Business() string {
	return av.attrs.Business
}

// InitIndex init index.
func (av *AvrArchive) InitIndex(c context.Context) {
	if aliases, err := av.d.GetAliases(av.attrs.ESName, av.attrs.Index.IndexAliasPrefix); err != nil {
		av.d.InitIndex(c, nil, av.attrs.ESName, av.attrs.Index.IndexAliasPrefix, av.attrs.Index.IndexEntityPrefix, av.attrs.Index.IndexMapping)
	} else {
		av.d.InitIndex(c, aliases, av.attrs.ESName, av.attrs.Index.IndexAliasPrefix, av.attrs.Index.IndexEntityPrefix, av.attrs.Index.IndexMapping)
	}
}

// InitOffset insert init value to offset.
func (av *AvrArchive) InitOffset(c context.Context) {
	av.d.InitOffset(c, av.offset, av.attrs, []string{})
	nowFormat := time.Now().Format("2006-01-02 15:04:05")
	av.offset.SetOffset(0, nowFormat)
}

// Offset get offset.
func (av *AvrArchive) Offset(c context.Context) {
	for {
		offset, err := av.d.Offset(c, av.appid, av.attrs.Table.TablePrefix)
		if err != nil {
			log.Error("ac.d.Offset error(%v)", err)
			time.Sleep(time.Second * 3)
			continue
		}
		av.offset.SetReview(offset.ReviewID, offset.ReviewTime)
		av.offset.SetOffset(offset.OffsetID(), offset.OffsetTime())
		break
	}
}

// SetRecover set recover
func (av *AvrArchive) SetRecover(c context.Context, recoverID int64, recoverTime string, i int) {
	av.offset.SetRecoverOffset(recoverID, recoverTime)
}

// IncrMessages .
func (av *AvrArchive) IncrMessages(c context.Context) (length int, err error) {
	var rows *xsql.Rows
	log.Info("appid: %s IncrMessages Current OffsetTime: %s, OffsetID: %d", av.appid, av.offset.OffsetTime, av.offset.OffsetID)
	if !av.offset.IsLoop {
		rows, err = av.db.Query(c, av.attrs.DataSQL.SQLByMTime, av.offset.OffsetTime, av.attrs.Other.Size)
	} else {
		rows, err = av.db.Query(c, av.attrs.DataSQL.SQLByIDMTime, av.offset.OffsetID, av.offset.OffsetTime, av.attrs.Other.Size)
	}
	if err != nil {
		log.Error("db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		item, row := dao.InitMapData(av.attrs.DataSQL.DataIndexFields)
		if err = rows.Scan(row...); err != nil {
			log.Error("IncrMessages rows.Scan() error(%v)", err)
			return
		}
		av.mapData = append(av.mapData, item)
	}
	length = len(av.mapData)
	if length > 0 {
		// offset
		dao.UpdateOffsetByMap(av.offset, av.mapData...)
		// extra relevant data
		length, err = av.extraData(c, "db", map[string]bool{"Avr": true})
	}
	return
}

// AllMessages .
func (av *AvrArchive) AllMessages(c context.Context) (length int, err error) {
	var rows *xsql.Rows
	log.Info("appid: %s allMessages Current RecoverID: %d", av.appid, av.offset.RecoverID)
	if rows, err = av.db.Query(c, av.attrs.DataSQL.SQLByID, av.offset.RecoverID, av.attrs.Other.Size); err != nil {
		log.Error("AllMessages db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		item, row := dao.InitMapData(av.attrs.DataSQL.DataIndexFields)
		if err = rows.Scan(row...); err != nil {
			log.Error("AllMessages rows.Scan() error(%v)", err)
			continue
		}
		av.mapData = append(av.mapData, item)
	}
	length = len(av.mapData)
	if length > 0 {
		// offset
		if av.mapData[length-1]["_id"] != nil {
			v := av.mapData[length-1]["_id"]
			if v2, ok := v.(*interface{}); ok {
				av.offset.SetTempOffset((*v2).(int64), "")
				av.offset.SetRecoverTempOffset((*v2).(int64), "")
			}
		}
		// extra relevant data
		length, err = av.extraData(c, "db", map[string]bool{"Avr": true})
	}
	return
}

// extraData extra data for appid
func (av *AvrArchive) extraData(c context.Context, way string, tags map[string]bool) (length int, err error) {
	switch way {
	case "db":
		for i, item := range av.mapData {
			item.TransData(av.attrs)
			for k, v := range item {
				av.mapData[i][k] = v
			}
		}
	case "dtb":
		for i, item := range av.mapData {
			item.TransDtb(av.attrs)
			av.mapData[i] = model.MapData{}
			for k, v := range item {
				av.mapData[i][k] = v
			}
		}
	}
	for _, ex := range av.attrs.DataExtras {
		if _, ok := tags[ex.Tag]; !ok {
			continue
		}
		switch ex.Type {
		case "slice":
			continue
			//av.extraDataSlice(c, ex)
		default:
			length, _ = av.extraDataDefault(c, ex)
		}
	}
	return
}

// extraData-default
func (av *AvrArchive) extraDataDefault(c context.Context, ex model.AttrDataExtra) (length int, err error) {
	// filter ids from in_fields
	var (
		ids   []int64
		items map[int64]model.MapData
		temp  map[int64]model.MapData
	)
	cdtInField := ex.Condition["in_field"]
	items = make(map[int64]model.MapData)
	temp = make(map[int64]model.MapData)
	for _, md := range av.mapData {
		if v, ok := md[cdtInField]; ok {
			ids = append(ids, v.(int64)) // 加去重
			temp[v.(int64)] = md
		}
	}
	// query extra data
	if len(ids) > 0 {
		var rows *xsql.Rows
		rows, err = av.d.DBPool[ex.DBName].Query(c, fmt.Sprintf(ex.SQL, xstr.JoinInts(ids))+" and 1 = ? ", 1)
		if err != nil {
			log.Error("extraDataDefault db.Query error(%v)", err)
			return
		}
		for rows.Next() {
			item, row := dao.InitMapData(ex.Fields)
			if err = rows.Scan(row...); err != nil {
				log.Error("extraDataDefault rows.Scan() error(%v)", err)
				continue
			}
			if v, ok := item[ex.InField]; ok {
				if v2, ok := v.(*interface{}); ok {
					item.TransData(av.attrs)
					items[(*v2).(int64)] = item
				}
			}
		}
		rows.Close()
	}
	//fmt.Println("a.mapData", av.mapData, "ids", ids, "items", items)
	// merge data
	fds := []string{"_id", "cid", "vid", "aid", "v_ctime"}
	av.mapData = []model.MapData{}
	for k, item := range items {
		if v, ok := temp[k]; ok {
			for _, fd := range fds {
				if f, ok := item[fd]; ok {
					v[fd] = f
				}
			}
			av.mapData = append(av.mapData, v)
		}
	}
	length = len(av.mapData)
	//fmt.Println("a.mapData:after", av.mapData)
	return
}

// BulkIndex .
func (av *AvrArchive) BulkIndex(c context.Context, start int, end int, writeEntityIndex bool) (err error) {
	partData := av.mapData[start:end]
	err = av.d.BulkDBData(c, av.attrs, writeEntityIndex, partData...)
	return
}

// Commit commit offset.
func (av *AvrArchive) Commit(c context.Context) (err error) {
	err = av.d.CommitOffset(c, av.offset, av.attrs.AppID, av.attrs.Table.TablePrefix)
	av.mapData = []model.MapData{}
	return
}

// Sleep interval duration.
func (av *AvrArchive) Sleep(c context.Context) {
	time.Sleep(time.Second * time.Duration(av.attrs.Other.Sleep))
}

// Size return size.
func (av *AvrArchive) Size(c context.Context) int {
	return av.attrs.Other.Size
}
