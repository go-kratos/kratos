package dao

import (
	"go-common/app/job/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

//GetAbleCode get able code
func (d *Dao) GetAbleCode(tx *sql.Tx, batchCodeID int64) (code string, err error) {
	row := tx.QueryRow(_getAbleCodeSQL, batchCodeID)
	if err = row.Scan(&code); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		err = errors.WithStack(err)
	}
	return
}

//UpdateCodeRelationID update code relationID
func (d *Dao) UpdateCodeRelationID(tx *sql.Tx, code, relationID string, bmid int64) (err error) {
	if _, err = tx.Exec(_updateCodeRelationIDSQL, relationID, bmid, code); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//SelBatchCodeByID .
func (d *Dao) SelBatchCodeByID(tx *sql.Tx, batchCodeID int64) (r *model.VipResourceBatchCode, err error) {
	row := tx.QueryRow(_selBatchCodeSQL, batchCodeID)
	r = new(model.VipResourceBatchCode)
	if err = row.Scan(&r.ID, &r.BusinessID, &r.PoolID, &r.Status, &r.Type, &r.BatchName, &r.Reason, &r.Unit, &r.Count, &r.SurplusCount, &r.Price, &r.StartTime, &r.EndTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		err = errors.WithStack(err)
	}
	return
}
