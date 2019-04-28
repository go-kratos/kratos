package thirdp

import (
	"context"
	"fmt"

	model "go-common/app/interface/main/tv/model/thirdp"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_mangoRecomKey = "mango_cms_recom"
	_mangoRecom    = "SELECT id, rid, rtype, title, cover , category, playcount, jid, content, staff , rorder FROM mango_recom WHERE deleted = 0 AND id IN (%s)"
)

// MangoOrder gets mango recom data.
func (d *Dao) MangoOrder(c context.Context) (s []int64, err error) {
	var (
		conn = d.mc.Get(c)
		item *memcache.Item
		res  = model.MangoOrder{}
	)
	defer conn.Close()
	if item, err = conn.Get(_mangoRecomKey); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Get(%s) error (%v)", _mangoRecomKey, err)
		return
	}
	if err = conn.Scan(item, &res); err != nil {
		log.Error("conn.Scan(%s) error(%v)", _mangoRecomKey, err)
	}
	s = res.RIDs
	return
}

// MangoRecom picks the mango recom data from DB
func (d *Dao) MangoRecom(c context.Context, ids []int64) (data []*model.MangoRecom, err error) {
	var (
		rows    *sql.Rows
		query   = fmt.Sprintf(_mangoRecom, xstr.JoinInts(ids))
		dataSet = make(map[int64]*model.MangoRecom, len(ids))
	)
	if rows, err = d.db.Query(c, query); err != nil {
		log.Error("mangoRecom, Err %v", err)
		return
	}
	// SELECT id, rid, rtype, title, cover , category, playcount, jid, content, staff , rorder
	for rows.Next() {
		var r = &model.MangoRecom{}
		if err = rows.Scan(&r.ID, &r.RID, &r.Rtype, &r.Title, &r.Cover, &r.Category, &r.Playcount, &r.JID, &r.Content, &r.Staff, &r.Rorder); err != nil {
			return
		}
		dataSet[r.ID] = r
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
		return
	}
	for _, v := range ids {
		if value, ok := dataSet[v]; ok {
			data = append(data, value)
			continue
		}
		log.Warn("MangoRecom RID %d Missing", v)
	}
	return
}
