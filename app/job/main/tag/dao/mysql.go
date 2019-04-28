package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"go-common/app/job/main/tag/model"
	"go-common/library/log"
)

const (
	capacity    = 1000
	_tagInfoSQL = "SELECT t.id,t.name,t.cover,t.head_cover,t.short_content,c.sub,c.bind FROM tag t LEFT OUTER JOIN `count` c ON t.id=c.tid WHERE t.id>? AND t.state=0 ORDER BY t.id LIMIT 1000 ;"
)

func replace(name string) string {
	var (
		there bool
		rb    []byte
	)
	sb := []byte(strings.Trim(name, " "))
	for _, b := range sb {
		if b < 0x20 || b == 0x7f {
			there = true
			continue
		}
		rb = append(rb, b)
	}
	if there {
		log.Warn("There are invisible characters,tag byte(%v),name(%v)", sb, name)
	}
	return string(rb)
}

// TagInfo .
func (d *Dao) TagInfo(c context.Context, tid int64) (res []*model.PlatformTagInfo, err error) {
	var (
		sub  sql.NullInt64
		bind sql.NullInt64
	)
	rows, err := d.platform.Query(c, _tagInfoSQL, tid)
	if err != nil {
		log.Error("TagInfo d.platform.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*model.PlatformTagInfo, 0, capacity)
	for rows.Next() {
		t := &model.PlatformTagInfo{}
		if err = rows.Scan(&t.ID, &t.Name, &t.Cover, &t.HeadCover, &t.ShortContent, &sub, &bind); err != nil {
			log.Error("TagInfo rows.Scan() error(%v)", err)
			return
		}
		if !sub.Valid {
			log.Warn("TagInfo() count_table no data(%d)", t.ID)
			continue
		}
		t.Name = replace(t.Name)
		t.Sub = sub.Int64
		t.Bind = bind.Int64
		res = append(res, t)
	}
	return
}

var _rids = "SELECT DISTINCT prid,rid FROM rank_result" // rid prid

// Rids .
func (d *Dao) Rids(c context.Context) (rpMap map[int64]int64, err error) {
	rows, err := d.platform.Query(c, _rids)
	if err != nil {
		log.Error("d.mysql.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	rpMap = make(map[int64]int64)
	for rows.Next() {
		var prid, rid int64
		if err = rows.Scan(&prid, &rid); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		rpMap[rid] = prid
	}
	return
}

var _hotMap = "SELECT rid,tid FROM rank_result WHERE type=0 ORDER BY rank"

// HotMap .
func (d *Dao) HotMap(c context.Context) (res map[int16][]int64, err error) {
	rows, err := d.platform.Query(c, _hotMap)
	if err != nil {
		log.Error("Hot d.hot.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int16][]int64)
	for rows.Next() {
		var rid int16
		var tid int64
		if err = rows.Scan(&rid, &tid); err != nil {
			log.Error("Hot rows.Scan() error(%v)", err)
			res = nil
			return
		}
		res[rid] = append(res[rid], tid)
	}
	return
}

var (
	_shard         = 200
	resourceTagSQL = "SELECT tid FROM resource_tag_%s WHERE oid=? AND type=? AND state=0"
)

func (d *Dao) hit(mid int64) string {
	return fmt.Sprintf("%03d", mid%int64(_shard))
}

// Resources return resources by oid from mysql.
func (d *Dao) Resources(c context.Context, oid int64, typ int32) (res []int64, err error) {
	rows, err := d.platform.Query(c, fmt.Sprintf(resourceTagSQL, d.hit(oid)), oid, typ)
	if err != nil {
		log.Error("d.Resources(%d,%d) error(%v)", oid, typ, err)
		return
	}
	defer rows.Close()
	res = make([]int64, 0)
	for rows.Next() {
		var tid int64
		if err = rows.Scan(&tid); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, tid)
	}
	return
}

var tagResourceSQL = "SELECT oid FROM tag_resource_%s WHERE tid=? AND `type`=3 AND state=0 ORDER BY ctime DESC LIMIT 1000 ;"

// TagResources return resources by oid from mysql.
func (d *Dao) TagResources(c context.Context, tid int64) (res []int64, err error) {
	rows, err := d.platform.Query(c, fmt.Sprintf(tagResourceSQL, d.hit(tid)), tid)
	if err != nil {
		log.Error("d.TagResources(%d) error(%v)", tid, err)
		return
	}
	defer rows.Close()
	res = make([]int64, 0)
	for rows.Next() {
		var oid int64
		if err = rows.Scan(&oid); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, oid)
	}
	return
}
