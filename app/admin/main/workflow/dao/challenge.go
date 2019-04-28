package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"go-common/app/admin/main/workflow/model"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Chall will retrive challenge by cid
func (d *Dao) Chall(c context.Context, cid int64) (chall *model.Chall, err error) {
	chall = &model.Chall{}
	if db := d.ReadORM.Table("workflow_chall").Where("id=?", cid).First(chall); db.Error != nil {
		err = db.Error
		if db.RecordNotFound() {
			chall = nil
			err = nil
		} else {
			err = errors.Wrapf(err, "chall(%d)", cid)
		}
	}
	return
}

// Challs will select challenges by ids
func (d *Dao) Challs(c context.Context, cids []int64) (challs map[int64]*model.Chall, err error) {
	challs = make(map[int64]*model.Chall, len(cids))
	if len(cids) <= 0 {
		return
	}

	chlist := make([]*model.Chall, 0)
	err = d.ReadORM.Table("workflow_chall").Where("id IN (?)", cids).Find(&chlist).Error
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	for _, c := range chlist {
		challs[c.Cid] = c
	}
	return
}

// StateChalls will select a set of groups by challenge ids and state
func (d *Dao) StateChalls(c context.Context, cids []int64, state int8) (challs map[int64]*model.Chall, err error) {
	challs = make(map[int64]*model.Chall, len(cids))
	challSlice := make([]*model.Chall, 0, len(cids))
	if len(cids) <= 0 {
		return
	}

	if err = d.ORM.Table("workflow_chall").Where("id IN (?)", cids).Find(&challSlice).Error; err != nil {
		err = errors.WithStack(err)
		return
	}

	for _, chall := range challSlice {
		chall.FromState()
		if chall.State != state {
			continue
		}
		challs[chall.Cid] = chall
	}

	return
}

// LastChallIDsByGids will select last chall ids by given gids
func (d *Dao) LastChallIDsByGids(c context.Context, gids []int64) (cids []int64, err error) {
	if len(gids) <= 0 {
		return
	}
	var rows *sql.Rows
	if rows, err = d.ReadORM.Table("workflow_chall").Select("max(id)").Where("gid IN (?)", gids).Group("gid").Rows(); err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var maxID int64
		if err = rows.Scan(&maxID); err != nil {
			return
		}
		cids = append(cids, maxID)
	}

	return
}

// TxUpChall will update state of a challenge
// Deprecated
func (d *Dao) TxUpChall(tx *gorm.DB, chall *model.Chall) (rows int64, err error) {
	// write old field
	chall.FromState()
	db := tx.Table("workflow_chall").Where("id=?", chall.Cid).
		Update("dispatch_state", chall.DispatchState)
	if db.Error != nil {
		err = errors.WithStack(db.Error)
		return
	}

	rows = db.RowsAffected
	return
}

// TxBatchUpChallByIDs will update state of challenges by cids
func (d *Dao) TxBatchUpChallByIDs(tx *gorm.DB, cids []int64, state int8) (err error) {
	challSlice := make([]*model.Chall, 0, len(cids))
	if len(cids) <= 0 {
		return
	}
	if err = tx.Table("workflow_chall").Where("id IN (?)", cids).Find(&challSlice).Error; err != nil {
		err = errors.WithStack(err)
		return
	}

	for _, chall := range challSlice {
		chall.SetState(uint32(state), uint8(0))
		if err = tx.Table("workflow_chall").Where("id=?", chall.Cid).
			Update("dispatch_state", chall.DispatchState).Error; err != nil {
			err = errors.WithStack(err)
			return
		}
	}

	return
}

// AttPathsByCids will select a set of attachments paths by challenge ids
func (d *Dao) AttPathsByCids(c context.Context, cids []int64) (paths map[int64][]string, err error) {
	var pathSlice []*struct {
		Cid  int64
		Path string
	}
	paths = make(map[int64][]string, len(cids))
	if len(cids) <= 0 {
		return
	}
	if err = d.ReadORM.Table("workflow_attachment").Where("cid IN (?)", cids).Select("cid,path").Find(&pathSlice).Error; err != nil {
		return
	}
	for _, cp := range pathSlice {
		if _, ok := paths[cp.Cid]; !ok {
			paths[cp.Cid] = make([]string, 0, 1)
		}
		paths[cp.Cid] = append(paths[cp.Cid], cp.Path)
	}
	return
}

// AttPathsByCid will select a set of attachments paths by challenge id
// Deprecated
func (d *Dao) AttPathsByCid(c context.Context, cid int64) (paths []string, err error) {
	paths = make([]string, 0)

	rows, err := d.ReadORM.Table("workflow_attachment").Select("cid,path").Where("cid=?", cid).Rows()
	if err != nil {
		err = errors.Wrapf(err, "cid(%d)", cid)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cp struct {
			cid  int32
			path string
		}
		if err = rows.Scan(&cp.cid, &cp.path); err != nil {
			err = errors.WithStack(err)
			return
		}
		paths = append(paths, cp.path)
	}

	return
}

// UpChallBusState will update specified business_state by conditions
// Deprecated
func (d *Dao) UpChallBusState(c context.Context, cid int64, busState int8, assigneeAdminid int64) (err error) {
	var chall *model.Chall
	if chall, err = d.Chall(c, cid); err != nil {
		return
	}
	if chall == nil {
		err = errors.Wrapf(err, "can not find challenge cid(%d)", cid)
		return
	}
	// double write to new field
	chall.SetState(uint32(busState), uint8(1))

	if err = d.ORM.Table("workflow_chall").Where("id=?", cid).
		Update("dispatch_state", chall.DispatchState).
		Update("assignee_adminid", assigneeAdminid).
		Error; err != nil {
		err = errors.WithStack(err)
	}
	return
}

// BatchUpChallBusState will update specified business_state by conditions
func (d *Dao) BatchUpChallBusState(c context.Context, cids []int64, busState int8, assigneeAdminid int64) (err error) {
	var challs map[int64]*model.Chall

	challs, err = d.Challs(c, cids)
	if err != nil {
		return
	}
	for cid := range challs {
		challs[cid].SetState(uint32(busState), uint8(1))
		if err = d.ORM.Table("workflow_chall").Where("id=?", cid).
			Update("dispatch_state", challs[cid].DispatchState).
			Update("assignee_adminid", assigneeAdminid).Error; err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	return
}

// TxChallsByBusStates select cids by business and oid
func (d *Dao) TxChallsByBusStates(tx *gorm.DB, business int8, oid int64, busStates []int8) (cids []int64, err error) {
	cids = make([]int64, 0)
	rows, err := tx.Table("workflow_chall").Where("business=? AND oid=? ",
		business, oid).Select("id,dispatch_state").Rows()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		c := &model.Chall{}
		if err = rows.Scan(&c.Cid, &c.DispatchState); err != nil {
			return
		}
		for _, busState := range busStates {
			c.FromState()
			if c.BusinessState == busState {
				cids = append(cids, int64(c.Cid))
			}
		}
	}
	return
}

// TxUpChallsBusStateByIDs will update specified business_state by conditions
func (d *Dao) TxUpChallsBusStateByIDs(tx *gorm.DB, cids []int64, busState int8, assigneeAdminid int64) (err error) {
	challSlice := make([]*model.Chall, 0)
	if err = tx.Table("workflow_chall").Where("id IN (?)", cids).
		Select("id,gid,mid,dispatch_state").Find(&challSlice).Error; err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, chall := range challSlice {
		chall.SetState(uint32(busState), uint8(1))
		if err = tx.Table("workflow_chall").Where("id=?", chall.Cid).
			Update("dispatch_state", chall.DispatchState).Update("assignee_adminid", assigneeAdminid).Error; err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	return
}

// TxUpChallExtraV2 will update Extra data by business oid
func (d *Dao) TxUpChallExtraV2(tx *gorm.DB, business int8, oid, adminid int64, extra map[string]interface{}) (rows int64, err error) {
	exData, err := json.Marshal(extra)
	if err != nil {
		err = errors.Wrapf(err, "business(%d) oid(%d), extra(%s)", business, oid, extra)
		return
	}
	if err = tx.Table("workflow_business").Where("business=? AND oid=?", business, oid).Update("extra", exData).Error; err != nil {
		err = errors.Wrapf(err, "business(%d), oid(%d), extra(%s)", business, oid, exData)
	}
	return
}

// UpExtraV3 will update Extra data by gids
func (d *Dao) UpExtraV3(gids []int64, adminid int64, extra string) error {
	return d.ORM.Table("workflow_business").Where("gid IN (?)", gids).Update("extra", extra).Error
}

// TxUpChallTag will update tid by cid
func (d *Dao) TxUpChallTag(tx *gorm.DB, cid int64, tid int64) (err error) {
	if err = tx.Table("workflow_chall").Where("id=?", cid).Update("tid", tid).Error; err != nil {
		err = errors.Wrapf(err, "cid(%d), tid(%d)", cid, tid)
	}
	return
}

// BatchUpChallByIDs update dispatch_state field of cids
func (d *Dao) BatchUpChallByIDs(cids []int64, dispatchState uint32, adminid int64) (err error) {
	if len(cids) <= 0 {
		return
	}
	if err = d.ORM.Table("workflow_chall").Where("id IN (?)", cids).
		Update("dispatch_state", dispatchState).Update("adminid", adminid).Error; err != nil {
		err = errors.WithStack(err)
	}
	return
}

// BatchResetAssigneeAdminID reset assignee_adminid by cids
func (d *Dao) BatchResetAssigneeAdminID(cids []int64) (err error) {
	if len(cids) <= 0 {
		return
	}
	if err = d.ORM.Table("workflow_chall").Where("id IN (?)", cids).
		Update("assignee_adminid", 0).Error; err != nil {
		err = errors.WithStack(err)
	}
	return
}

// TxUpChallAssignee update assignee_adminid and dispatch_time when admin start a mission
func (d *Dao) TxUpChallAssignee(tx *gorm.DB, cids []int64) error {
	return tx.Table("workflow_chall").Where("id IN (?)", cids).
		Update("dispatch_time", time.Now().Format("2006-01-02 15:04:05")).Error
}
