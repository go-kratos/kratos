package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/tag/model"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

const (
	// newest arcs of tag of region
	_regionArcKey = "ra_%d_%d"
	// origin newest arcs of tag of region
	_regionOriArcKey = "ro_%d_%d"
)

func regionArcKey(rid, tid int64) string {
	return fmt.Sprintf(_regionArcKey, rid, tid)
}

func regionOriArcKey(rid, tid int64) string {
	return fmt.Sprintf(_regionOriArcKey, rid, tid)
}

// AddTagArcs add region archives.
func (d *Dao) AddTagArcs(c context.Context, tid int64, arcMap map[int64]*model.SearchRes) (err error) {
	conn := d.redisRank.Get(c)
	defer conn.Close()
	var count int
	for _, arc := range arcMap {
		if arc.ID == 0 || arc.PubDate == "" || !arc.IsNormal() {
			log.Warn("d.AddTagArcs(%d) archive(%v) not normal.", tid, arc)
			continue
		}
		var pubTime time.Time
		if pubTime, err = time.ParseInLocation(_timeLayoutFormat, arc.PubDate, time.Local); err != nil {
			log.Error("d.AddTagArcs(%d) time.ParseInLocation(%v) error(%v)", tid, arc, err)
			continue
		}
		count++
		if err = conn.Send("ZADD", regionArcKey(arc.TypeID, tid), pubTime.Unix(), arc.ID); err != nil {
			log.Error("d.AddRegionArcs(ZADD, %d, %d, %v) error(%v)", tid, arc.ID, pubTime, err)
			return
		}
		if arc.Copyright == int32(archive.CopyrightOriginal) {
			count++
			if err = conn.Send("ZADD", regionOriArcKey(arc.TypeID, tid), pubTime.Unix(), arc.ID); err != nil {
				log.Error("d.AddRegionArcs Original Archives(ZADD, %d, %d, %v) error(%v)", tid, arc.ID, pubTime, err)
				return
			}
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("d.AddTagArcs(%d,%+v) Flush() error(%v)", tid, arcMap, err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("d.AddTagArcs(%d,%+v) Receive() error(%v)", tid, arcMap, err)
			return
		}
	}
	return
}
