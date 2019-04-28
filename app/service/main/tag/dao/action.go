package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/tag/model"
	"go-common/library/log"
)

var actionsSQL = "SELECT oid,tid,action FROM resource_tag_action_%s WHERE mid=? AND type=? LIMIT 500"

// Actions return user's acionts .
func (d *Dao) Actions(c context.Context, mid int64, typ int32) (rs []*model.ResourceAction, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(actionsSQL, d.hit(mid)), mid, typ)
	if err != nil {
		log.Error("db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.ResourceAction{}
		if err = rows.Scan(&r.Oid, &r.Tid, &r.Action); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		rs = append(rs, r)
	}
	return
}

var addActionSQL = "INSERT IGNORE INTO resource_tag_action_%s (oid,type,tid,mid,action) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE action=?"

// AddAction update a resource action.
func (d *Dao) AddAction(c context.Context, a *model.ResourceAction) (rows int64, err error) {
	row, err := d.db.Exec(c, fmt.Sprintf(addActionSQL, d.hit(a.Mid)), a.Oid, a.Type, a.Tid, a.Mid, a.Action, a.Action)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return row.RowsAffected()
}
