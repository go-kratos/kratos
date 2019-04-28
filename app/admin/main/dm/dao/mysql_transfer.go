package dao

import (
	"context"

	"go-common/app/admin/main/dm/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_InstTrJobSQL        = "INSERT INTO dm_transfer_job(from_cid,to_cid,mid,offset,state) VALUES(?,?,?,?,?)"
	_countTransList      = "SELECT count(*) FROM dm_transfer_job WHERE to_cid=?"
	_countTransListState = "SELECT count(*) FROM dm_transfer_job WHERE to_cid=? AND state=?"
	_TransList           = "SELECT id,from_cid,to_cid,state,ctime FROM dm_transfer_job WHERE to_cid=? limit ?,?"
	_TransListState      = "SELECT id,from_cid,to_cid,state,ctime FROM dm_transfer_job WHERE to_cid=? AND state=? limit ?, ?"
	_selTransferID       = "SELECT id,from_cid,to_cid,mid,offset,state,ctime,mtime from dm_transfer_job WHERE id=?"
	_uptTransferSQL      = "UPDATE dm_transfer_job SET state=? WHERE id=?"
)

// InsertTransferJob add transfer job
func (d *Dao) InsertTransferJob(c context.Context, from, to, mid int64, offset float64, state int8) (affect int64, err error) {
	row, err := d.biliDM.Exec(c, _InstTrJobSQL, from, to, mid, offset, state)
	if err != nil {
		log.Error("biliDM.Exec(%s, %d %d %d %v) error(%v)", _InstTrJobSQL, from, to, mid, offset, err)
		return
	}
	return row.LastInsertId()
}

// TransferList transfer list
func (d *Dao) TransferList(c context.Context, cid, state, pn, ps int64) (res []*model.TransList, total int64, err error) {
	var rows *sql.Rows
	res = make([]*model.TransList, 0)
	if state == int64(model.TransferJobStateAll) {
		countRow := d.biliDM.QueryRow(c, _countTransList, cid)
		if err = countRow.Scan(&total); err != nil {
			log.Error("row.ScanCount error(%v)", err)
			return
		}
		rows, err = d.biliDM.Query(c, _TransList, cid, (pn-1)*ps, ps)
		if err != nil {
			log.Error("biliDM.Query(%s, %d ) error(%v)", _TransList, cid, err)
			return
		}
	} else {
		countRow := d.biliDM.QueryRow(c, _countTransListState, cid, state)
		if err = countRow.Scan(&total); err != nil {
			log.Error("row.ScanCount error(%v)", err)
			return
		}
		rows, err = d.biliDM.Query(c, _TransListState, cid, state, (pn-1)*ps, ps)
		if err != nil {
			log.Error("biliDM.Query(%s, %d, %d) error(%v)", _TransList, cid, state, err)
			return
		}
	}
	defer rows.Close()
	for rows.Next() {
		dm := &model.TransList{}
		if err = rows.Scan(&dm.ID, &dm.From, &dm.To, &dm.State, &dm.Ctime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res = append(res, dm)
	}
	return
}

// CheckTransferID check transfer job state by id
func (d *Dao) CheckTransferID(c context.Context, id int64) (job *model.TransferJobInfo, err error) {
	job = new(model.TransferJobInfo)
	row := d.biliDM.QueryRow(c, _selTransferID, id)
	if err = row.Scan(&job.ID, &job.FromCID, &job.ToCID, &job.MID, &job.Offset, &job.State, &job.Ctime, &job.Mtime); err != nil {
		if err == sql.ErrNoRows {
			job = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// SetTransferState change transfer state
func (d *Dao) SetTransferState(c context.Context, id int64, state int8) (affect int64, err error) {
	row, err := d.biliDM.Exec(c, _uptTransferSQL, state, id)
	if err != nil {
		log.Error("d.biliDM.Exec(%s,%d) error(%v)", _uptTransferSQL, id, err)
		return
	}
	return row.RowsAffected()
}
