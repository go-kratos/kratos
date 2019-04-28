package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/dm/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_selAllIdxSQL   = "SELECT id,type,oid,mid,progress,state,pool,attr,ctime,mtime FROM dm_index_%03d WHERE type=? AND oid=? AND state IN(0,6)"
	_selIdxHidesSQL = "SELECT id,type,oid,mid,progress,state,pool,attr,ctime,mtime FROM dm_index_%03d FORCE INDEX(ix_oid_state) WHERE type=? AND oid=? AND state=2 ORDER BY id DESC limit ?" // NOTE slow query
	_upIdxSQL       = "UPDATE dm_index_%03d SET mid=?,progress=?,state=?,pool=?,attr=? WHERE id=?"
	_upIdxStatesSQL = "UPDATE dm_index_%03d SET state=? WHERE id IN(%s)"
)

// DMInfos get indexs of oid.
func (d *Dao) DMInfos(c context.Context, tp int32, oid int64) (dms []*model.DM, err error) {
	rows, err := d.dmReader.Query(c, fmt.Sprintf(_selAllIdxSQL, d.hitIndex(oid)), tp, oid)
	if err != nil {
		log.Error("db.Query(%d %d) error(%v)", tp, oid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		dm := &model.DM{}
		if err = rows.Scan(&dm.ID, &dm.Type, &dm.Oid, &dm.Mid, &dm.Progress, &dm.State, &dm.Pool, &dm.Attr, &dm.Ctime, &dm.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		dms = append(dms, dm)
	}
	return
}

// DMHides get hide index info from db by oid and state.
func (d *Dao) DMHides(c context.Context, typ int32, oid, limit int64) (res []*model.DM, err error) {
	rows, err := d.dmReader.Query(c, fmt.Sprintf(_selIdxHidesSQL, d.hitIndex(oid)), typ, oid, limit)
	if err != nil {
		log.Error("db.Query(%d %d) error(%v)", typ, oid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		dm := &model.DM{}
		if err = rows.Scan(&dm.ID, &dm.Type, &dm.Oid, &dm.Mid, &dm.Progress, &dm.State, &dm.Pool, &dm.Attr, &dm.Ctime, &dm.Mtime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		res = append(res, dm)
	}
	return
}

// UpdateDM update index of dm.
func (d *Dao) UpdateDM(c context.Context, m *model.DM) (affect int64, err error) {
	res, err := d.dmWriter.Exec(c, fmt.Sprintf(_upIdxSQL, d.hitIndex(m.Oid)), m.Mid, m.Progress, m.State, m.Pool, m.Attr, m.ID)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// UpdateDMStates multi update index state of dm.
func (d *Dao) UpdateDMStates(c context.Context, oid int64, dmids []int64, state int32) (affect int64, err error) {
	upSQL := fmt.Sprintf(_upIdxStatesSQL, d.hitIndex(oid), xstr.JoinInts(dmids))
	res, err := d.dmWriter.Exec(c, upSQL, state)
	if err != nil {
		log.Error("db.Exec(%s) error(%v)", upSQL, err)
		return
	}
	return res.RowsAffected()
}
