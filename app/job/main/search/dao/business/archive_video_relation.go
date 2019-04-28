package business

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/search/conf"
	"go-common/app/job/main/search/dao"
	"go-common/app/job/main/search/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// Avr single table consume databus.
type Avr struct {
	c       *conf.Config
	d       *dao.Dao
	appid   string
	attrs   *model.Attrs
	db      *xsql.DB
	dtb     *databus.Databus
	offset  *model.LoopOffset
	mapData []model.MapData
	commits map[int32]*databus.Message
}

// NewAvr .
func NewAvr(d *dao.Dao, appid string, c *conf.Config) (a *Avr) {
	a = &Avr{
		c:       c,
		d:       d,
		appid:   appid,
		attrs:   d.AttrPool[appid],
		offset:  &model.LoopOffset{},
		mapData: []model.MapData{},
		db:      d.DBPool[d.AttrPool[appid].DBName],
		dtb:     d.DatabusPool[d.AttrPool[appid].Databus.Databus],
		commits: make(map[int32]*databus.Message),
	}
	return
}

// Business return business.
func (a *Avr) Business() string {
	return a.attrs.Business
}

// InitIndex init index.
func (a *Avr) InitIndex(c context.Context) {
	if aliases, err := a.d.GetAliases(a.attrs.ESName, a.attrs.Index.IndexAliasPrefix); err != nil {
		a.d.InitIndex(c, nil, a.attrs.ESName, a.attrs.Index.IndexAliasPrefix, a.attrs.Index.IndexEntityPrefix, a.attrs.Index.IndexMapping)
	} else {
		a.d.InitIndex(c, aliases, a.attrs.ESName, a.attrs.Index.IndexAliasPrefix, a.attrs.Index.IndexEntityPrefix, a.attrs.Index.IndexMapping)
	}
}

// InitOffset insert init value to offset.
func (a *Avr) InitOffset(c context.Context) {
	a.d.InitOffset(c, a.offset, a.attrs, []string{})
	nowFormat := time.Now().Format("2006-01-02 15:04:05")
	a.offset.SetOffset(0, nowFormat)
}

// Offset get offset.
func (a *Avr) Offset(c context.Context) {
	for {
		offset, err := a.d.Offset(c, a.appid, a.attrs.Table.TablePrefix)
		if err != nil {
			log.Error("a.d.Offset error(%v)", err)
			time.Sleep(time.Second * 3)
			continue
		}
		a.offset.SetReview(offset.ReviewID, offset.ReviewTime)
		a.offset.SetOffset(offset.OffsetID(), offset.OffsetTime())
		break
	}
}

// SetRecover set recover
func (a *Avr) SetRecover(c context.Context, recoverID int64, recoverTime string, i int) {
	a.offset.SetRecoverOffset(recoverID, recoverTime)
}

// IncrMessages .
func (a *Avr) IncrMessages(c context.Context) (length int, err error) {
	ticker := time.NewTicker(time.Duration(time.Millisecond * time.Duration(a.attrs.Databus.Ticker)))
	defer ticker.Stop()
	for {
		select {
		case msg, ok := <-a.dtb.Messages():
			if !ok {
				log.Error("databus: %s binlog consumer exit!!!", a.attrs.Databus)
				break
			}
			m := &model.Message{}
			a.commits[msg.Partition] = msg
			if err = json.Unmarshal(msg.Value, m); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
				continue
			}
			//fmt.Println("origin msg", m)
			//log.Info("origin msg: (%v)", m.Table, a.mapData)
			if m.Table == a.attrs.Table.TablePrefix {
				if m.Action == "insert" || m.Action == "update" {
					var parseMap map[string]interface{}
					parseMap, err = a.d.JSON2map(m.New)
					if err != nil {
						log.Error("a.JSON2map error(%v)", err)
						continue
					}
					a.mapData = append(a.mapData, parseMap)
				}
			}
			if len(a.mapData) < a.attrs.Databus.AggCount {
				continue
			}
		case <-ticker.C:
		}
		break
	}
	//log.Info("origin msg: (%v)", a.mapData)
	if len(a.mapData) > 0 {
		//fmt.Println("before", a.mapData)
		a.mapData, err = a.d.ExtraData(c, a.mapData, a.attrs, "dtb", []string{"archive", "video", "audit", "ups"})
		//fmt.Println("after", a.mapData)
		log.Info("dtb msg: (%v)", a.mapData)
	}
	length = len(a.mapData)
	return
}

// AllMessages .
func (a *Avr) AllMessages(c context.Context) (length int, err error) {
	rows, err := a.db.Query(c, a.attrs.DataSQL.SQLByID, a.offset.RecoverID, a.attrs.Other.Size)
	log.Info("appid: %s allMessages Current RecoverID: %d", a.appid, a.offset.RecoverID)
	if err != nil {
		log.Error("AllMessages db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		item, row := dao.InitMapData(a.attrs.DataSQL.DataIndexFields)
		if err = rows.Scan(row...); err != nil {
			log.Error("AllMessages rows.Scan() error(%v)", err)
			return
		}
		a.mapData = append(a.mapData, item)
	}
	if len(a.mapData) > 0 {
		//fmt.Println("before", a.mapData)
		a.mapData, err = a.d.ExtraData(c, a.mapData, a.attrs, "db", []string{"audit", "ups"})
		if v, ok := a.mapData[len(a.mapData)-1]["_id"]; ok {
			a.offset.SetTempOffset(v.(int64), "")
			a.offset.SetRecoverTempOffset(v.(int64), "")
		}
	}
	length = len(a.mapData)
	return
}

// BulkIndex .
func (a *Avr) BulkIndex(c context.Context, start, end int, writeEntityIndex bool) (err error) {
	if len(a.mapData) > 0 {
		partData := a.mapData[start:end]
		err = a.d.BulkDBData(c, a.attrs, writeEntityIndex, partData...)
	}
	return
}

// Commit commit offset.
func (a *Avr) Commit(c context.Context) (err error) {
	if a.c.Business.Index {
		err = a.d.CommitOffset(c, a.offset, a.attrs.AppID, a.attrs.Table.TablePrefix)
	} else {
		for k, cos := range a.commits {
			if err = cos.Commit(); err != nil {
				log.Error("appid(%s) commit error(%v)", a.attrs.AppID, err)
				continue
			}
			delete(a.commits, k)
		}
	}
	a.mapData = []model.MapData{}
	return
}

// Sleep interval duration.
func (a *Avr) Sleep(c context.Context) {
	time.Sleep(time.Second * time.Duration(a.attrs.Other.Sleep))
}

// Size return size.
func (a *Avr) Size(c context.Context) int {
	return a.attrs.Other.Size
}
