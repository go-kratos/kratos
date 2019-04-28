package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/dm2/model"
	"go-common/library/log"
)

const (
	_selectTransfer = "SELECT id,from_cid,to_cid,mid,offset,state,ctime FROM dm_transfer_job WHERE state=? limit 1"
	_updateTransfer = "UPDATE dm_transfer_job SET state=?,dmid=? WHERE id=?"
	_idxsSQL        = "SELECT id,type,oid,mid,progress,state,pool,attr,ctime,mtime FROM dm_index_%03d FORCE INDEX(ix_oid_state) WHERE type=? AND oid=? AND id >? ORDER BY id limit ?"
)

// Transfers get all transfer job
func (d *Dao) Transfers(c context.Context, state int8) (trans []*model.Transfer, err error) {
	rows, err := d.biliDMWriter.Query(c, _selectTransfer, model.StatInit)
	if err != nil {
		log.Error("d.biliDMWriter.Query(sql:%s) error(%v)", _selectTransfer, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := &model.Transfer{}
		if err = rows.Scan(&t.ID, &t.FromCid, &t.ToCid, &t.Mid, &t.Offset, &t.State, &t.Ctime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		trans = append(trans, t)
	}
	return
}

// UpdateTransfer change transfer job state.
func (d *Dao) UpdateTransfer(c context.Context, t *model.Transfer) (affect int64, err error) {
	row, err := d.biliDMWriter.Exec(c, _updateTransfer, t.State, t.Dmid, t.ID)
	if err != nil {
		log.Error("d.biliDMWriter.Exec(%+v) error(%v)", t, err)
		return
	}
	return row.RowsAffected()
}

// DMIndexs get dm indexs info
func (d *Dao) DMIndexs(c context.Context, tp int32, oid, minID, limit int64) (idxMap map[int64]*model.DM, dmids, special []int64, err error) {
	query := fmt.Sprintf(_idxsSQL, d.hitIndex(oid))
	rows, err := d.dmReader.Query(c, query, tp, oid, minID, limit)
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	idxMap = make(map[int64]*model.DM)
	for rows.Next() {
		idx := new(model.DM)
		if err = rows.Scan(&idx.ID, &idx.Type, &idx.Oid, &idx.Mid, &idx.Progress, &idx.State, &idx.Pool, &idx.Attr, &idx.Ctime, &idx.Mtime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		idxMap[idx.ID] = idx
		dmids = append(dmids, idx.ID)
		if idx.Pool == model.PoolSpecial {
			special = append(special, idx.ID)
		}
	}
	return
}
