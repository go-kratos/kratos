package archive

import (
	"context"

	"go-common/app/service/main/archive/api"
	"go-common/library/log"
)

const (
	_rgnArcsSQL = "SELECT aid,attribute,copyright,pubtime FROM archive WHERE typeid=? and (state>=0 or state=-6) LIMIT ?,?"
)

// RegionArcs multi get archvies by rid.
func (d *Dao) RegionArcs(c context.Context, rid int16, start, length int) (ras []*api.RegionArc, err error) {
	d.infoProm.Incr("RegionArcs")
	rows, err := d.rgnArcsStmt.Query(c, rid, start, length)
	if err != nil {
		log.Error("rgnArcsStmt.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &api.RegionArc{}
		if err = rows.Scan(&a.Aid, &a.Attribute, &a.Copyright, &a.PubDate); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		ras = append(ras, a)
	}
	return
}
