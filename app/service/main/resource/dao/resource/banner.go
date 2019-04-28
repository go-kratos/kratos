package resource

import (
	"context"
	"fmt"
	"time"

	"go-common/app/service/main/resource/model"
	"go-common/library/log"
)

var (
	_bannerSQL = `SELECT rm.id,ra.position_id,rm.name,rm.url,rm.pic,ra.position,ra.platform,ra.rule,ra.atype,ra.stime FROM resource_assignment AS ra,resource_material AS rm 
	WHERE ra.resource_group_id>0 AND ra.category=0 AND ra.stime<? AND ra.etime>? AND ra.state=0 AND ra.audit_state IN (2,3,4) AND 
	ra.id=rm.resource_assignment_id AND rm.audit_state=2 AND rm.category=0 ORDER BY ra.position ASC,ra.weight DESC,rm.ctime DESC`
	_categorySQL = `SELECT rm.id,ra.resource_id,rm.name,rm.url,rm.pic,ra.position,ra.platform,ra.rule,ra.atype,ra.stime FROM resource_assignment AS ra,resource_material AS rm 
	WHERE ra.id=rm.resource_assignment_id AND rm.id IN (SELECT max(rm.id) FROM resource_assignment AS ra,resource_material AS rm WHERE ra.resource_group_id>0 AND ra.category=1 
	AND ra.stime<? AND ra.etime>? AND ra.state=0 AND ra.audit_state IN (2,3,4) AND ra.id=rm.resource_assignment_id AND rm.audit_state=2 AND rm.category=1 GROUP BY rm.resource_assignment_id) ORDER BY rand()`
	_bannerLimitSQL = "SELECT rule FROM default_one WHERE id=2"
)

// Banner get active banner from new db.
func (d *Dao) Banner(c context.Context) (res map[int8]map[int][]*model.Banner, err error) {
	var (
		now = time.Now()
		ok  bool
		pm  map[string]string
	)
	rows, err := d.db.Query(c, _bannerSQL, now, now)
	if err != nil {
		log.Error("query error(%v)", err)
		return
	}
	defer rows.Close()
	res = map[int8]map[int][]*model.Banner{}
	pm = make(map[string]string)
	for rows.Next() {
		b := &model.Banner{}
		if err = rows.Scan(&b.ID, &b.ParentID, &b.Title, &b.Value, &b.Image, &b.Rank, &b.Plat, &b.Rule, &b.Type, &b.Start); err != nil {
			log.Error("rows.Scan error(%v)", err)
			res = nil
			return
		}
		pindex := fmt.Sprintf("%d_%d_%d", b.Plat, b.ParentID, b.Rank)
		if _, ok = pm[pindex]; ok {
			continue
		}
		b.BannerChange()
		if plm, ok := res[b.Plat]; ok {
			plm[b.ParentID] = append(plm[b.ParentID], b)
		} else {
			res[b.Plat] = map[int][]*model.Banner{
				b.ParentID: []*model.Banner{b},
			}
		}
		pm[pindex] = pindex
	}
	err = rows.Err()
	return
}

// Category get category banner from new db.
func (d *Dao) Category(c context.Context) (res map[int8]map[int][]*model.Banner, err error) {
	var now = time.Now()
	rows, err := d.db.Query(c, _categorySQL, now, now)
	if err != nil {
		log.Error("query error(%v)", err)
		return
	}
	defer rows.Close()
	res = map[int8]map[int][]*model.Banner{}
	for rows.Next() {
		b := &model.Banner{}
		if err = rows.Scan(&b.ID, &b.ParentID, &b.Title, &b.Value, &b.Image, &b.Rank, &b.Plat, &b.Rule, &b.Type, &b.Start); err != nil {
			log.Error("rows.Scan error(%v)", err)
			res = nil
			return
		}
		b.BannerChange()
		if plm, ok := res[b.Plat]; ok {
			plm[b.ParentID] = append(plm[b.ParentID], b)
		} else {
			res[b.Plat] = map[int][]*model.Banner{
				b.ParentID: []*model.Banner{b},
			}
		}
	}
	err = rows.Err()
	return
}

// Limit get app banner limit num.
func (d *Dao) Limit(c context.Context) (res map[int]int, err error) {
	row := d.db.QueryRow(c, _bannerLimitSQL)
	b := &model.Limit{}
	if err = row.Scan(&b.Rule); err != nil {
		log.Error("Limit row.Scan error(%v)", err)
		return
	}
	res = b.LimitChange()
	return
}
