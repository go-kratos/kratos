package gorm

import (
	"context"
	"database/sql"
	"encoding/json"

	"go-common/app/admin/main/aegis/model/resource"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
	"go-common/app/admin/main/aegis/model"
)

const (
	_resourceResSQL = "SELECT r.id, r.business_id, r.oid, r.mid,r.content,r.extra1,r.extra2,r.extra3,r.extra4,r.extra1s,r.extra2s,r.metadata, rr.attribute, rr.note, rr.reject_reason, rr.reason_id, rr.state, rr.pubtime, rr.deltime " +
		"FROM resource r LEFT JOIN resource_result rr ON r.id = rr.rid WHERE r.id = ?"
	_resByOIDSQL = "SELECT r.id, r.business_id, r.oid, r.mid,r.content,r.extra1,r.extra2,r.extra3,r.extra4,r.extra1s,r.extra2s,r.metadata, rr.attribute, rr.note, rr.reject_reason, rr.reason_id, rr.state, rr.pubtime, rr.deltime " +
		"FROM resource r LEFT JOIN resource_result rr ON r.id = rr.rid WHERE r.business_id = ? AND r.oid = ?"
)

var _changeableFields = map[string]struct{}{
	"extra1":     {},
	"extra2":     {},
	"extra3":     {},
	"extra4":     {},
	"extra5":     {},
	"extra6":     {},
	"extra1s":    {},
	"extra2s":    {},
	"extra3s":    {},
	"extra4s":    {},
	"extratime1": {},
}

// ListHelperForTask 补充任务列表里面的oid和content
func (d *Dao) ListHelperForTask(c context.Context, rids []int64) (res map[int64][]interface{}, err error) {
	var (
		rows                   *sql.Rows
		id                     int64
		oid, content, metadata string
	)
	if rows, err = d.orm.Table("resource").Select("id,oid,content,metadata").Where("id IN (?)", rids).Rows(); err != nil {
		log.Error("listHelperForTask err(%v)", err)
		return
	}

	res = make(map[int64][]interface{})
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&id, &oid, &content, &metadata); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		metas := make(map[string]interface{})
		if len(metadata) > 0 {
			if err = json.Unmarshal([]byte(metadata), &metas); err != nil {
				log.Error("ListHelperForTask err(%+v)\n", err)
			}
			err = nil
		}
		res[id] = []interface{}{oid, content, metas}
	}
	return
}

// TxAddResource .
func (d *Dao) TxAddResource(tx *gorm.DB, r *resource.Resource, rr *resource.Result) (rid int64, err error) {
	err = tx.Table("resource").Where("business_id=? AND oid=?", r.BusinessID, r.OID).
		Assign(map[string]interface{}{
			"oid":        r.OID,
			"content":    r.Content,
			"extra1":     r.Extra1,
			"extra2":     r.Extra2,
			"extra3":     r.Extra3,
			"extra4":     r.Extra4,
			"extra1s":    r.Extra1s,
			"extra2s":    r.Extra2s,
			"metadata":   r.MetaData,
			"extra5":     r.Extra5,
			"extra6":     r.Extra6,
			"extra3s":    r.Extra3s,
			"extra4s":    r.Extra4s,
			"extratime1": r.ExtraTime1,
			"octime":     r.OCtime,
			"ptime":      r.Ptime,
		}).FirstOrCreate(&r).Error
	if err != nil {
		return
	}
	rid = r.ID
	rr.RID = r.ID
	err = tx.Table("resource_result").Where("rid=?", rr.RID).Assign(map[string]interface{}{
		"attribute":     rr.Attribute,
		"note":          rr.Note,
		"reject_reason": rr.RejectReason,
		"reason_id":     rr.ReasonID,
		"state":         rr.State,
		"pubtime":       rr.PubTime,
		"deltime":       rr.DelTime,
	}).FirstOrCreate(&rr).Error
	return
}

// ResourceByOID .
func (d *Dao) ResourceByOID(c context.Context, OID string, bizID int64) (res *resource.Resource, err error) {
	res = &resource.Resource{}
	if err = d.orm.Where("oid = ? AND business_id = ?", OID, bizID).First(res).Error; err == gorm.ErrRecordNotFound {
		res = nil
		err = nil
	}
	return
}

// RidsByOids .
func (d *Dao) RidsByOids(c context.Context, bizID int64, oids []string) (rids string, err error) {
	err = d.orm.Table("resource").Select("GROUP_CONCAT(id)").Where("business_id=? AND oid IN (?)", bizID, oids).Row().Scan(&rids)
	return
}

// OidByRID 通过rid查到oid
func (d *Dao) OidByRID(c context.Context, rid int64) (oid string, err error) {
	if err = d.orm.Table("resource").Select("oid").Where("id=?", rid).Row().Scan(&oid); err == sql.ErrNoRows {
		err = nil
	}
	return
}

// TxUpdateResult .
func (d *Dao) TxUpdateResult(ormTx *gorm.DB, rid int64, res map[string]interface{}, rscRes *resource.Result) (err error) {
	params := make(map[string]interface{})
	for k, v := range res {
		if k != "state" && k != "attribute" {
			continue
		}
		params[k] = v
	}

	db := ormTx.Table("resource_result").
		Where("rid = ?", rid)
		//Update(params)
	if rscRes != nil {
		if rscRes.Attribute != -1 {
			params["attribute"] = rscRes.Attribute
		}
		if rscRes.Note != "" {
			params["note"] = rscRes.Note
		}
		if rscRes.RejectReason != "" {
			params["reject_reason"] = rscRes.RejectReason
		}
		if rscRes.ReasonID != 0 {
			params["reason_id"] = rscRes.ReasonID
		}
	}
	return db.Update(params).Error
}

// TxUpdateResource . TODO Resource表的字段要不要直接审核提交变更 还是等待业务方同步更新？(例如单话的上线下线)
func (d *Dao) TxUpdateResource(ormTx *gorm.DB, rid int64, res map[string]interface{}) (err error) {
	params := make(map[string]interface{})
	for k, v := range res {
		if _, ok := _changeableFields[k]; !ok {
			continue
		}
		params[k] = v
	}

	return ormTx.Table("resource").Where("id = ?", rid).Update(params).Error
}

// TxUpdateState 更新状态
func (d *Dao) TxUpdateState(tx *gorm.DB, rids []int64, state int) (err error) {
	return tx.Table("resource_result").Where("rid IN (?)", rids).Update("state", state).Error
}

// UpdateResource 更新资源
func (d *Dao) UpdateResource(c context.Context, bizid int64, oid string, update map[string]interface{}) (rows int64, err error) {
	db := d.orm.Table("resource").Where("business_id=? AND oid=?", bizid, oid).Update(update)
	return db.RowsAffected, db.Error
}

// ResourceRes .
func (d *Dao) ResourceRes(c context.Context, rid int64) (res *resource.Res, err error) {
	res = &resource.Res{}
	if err = d.orm.Raw(_resourceResSQL, rid).Scan(res).Error; err == gorm.ErrRecordNotFound {
		res = nil
		err = nil
	}
	return
}

// ResByOID .
func (d *Dao) ResByOID(c context.Context, bizID int64, OID string) (res *resource.Res, err error) {
	res = &resource.Res{}
	if err = d.orm.Raw(_resByOIDSQL, bizID, OID).Scan(res).Error; err == gorm.ErrRecordNotFound {
		res = nil
		err = nil
	}
	return
}

//ResOIDByID 根据id获取资源oid
func (d *Dao) ResOIDByID(c context.Context, rids []int64) (res map[int64]string, err error) {
	list := []struct {
		ID  int64  `gorm:"column:id"`
		OID string `gorm:"column:oid"`
	}{}
	res = map[int64]string{}
	if err = d.orm.Table("resource").Select("id,oid").Where("id in (?)", rids).Find(&list).Error; err != nil {
		log.Error("ResOIDByID error(%v) rids(%v)", err, rids)
		return
	}

	for _, item := range list {
		res[item.ID] = item.OID
	}
	return
}

//ResIDByOID 根据oid获取rid
func (d *Dao) ResIDByOID(C context.Context, bizID int64, oids []string) (res map[string]int64, err error) {
	list := []struct {
		ID  int64  `gorm:"column:id"`
		OID string `gorm:"column:oid"`
	}{}
	res = map[string]int64{}
	if err = d.orm.Table("resource").Select("id,oid").Where("business_id=? and oid in (?)", bizID, oids).Find(&list).Error; err != nil {
		log.Error("ResIDByOID error(%v) oids(%v)", err, oids)
		return
	}

	for _, item := range list {
		res[item.OID] = item.ID
	}
	return
}

//ResourceHit 根据状态筛选资源
func (d *Dao) ResourceHit(c context.Context, rids []int64) (hitids map[int64]int64, err error) {
	hitids = make(map[int64]int64)
	rows, err := d.orm.Table("resource_result").Select("rid,state").Where("rid IN (?)", rids).Rows()
	if err != nil {
		log.Error("ResourceHitState error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var rid, state int64
		if err = rows.Scan(&rid, &state); err != nil {
			log.Error("ResourceHitState error(%v)", err)
		}
		hitids[rid] = state
	}
	return
}

//ResultByRID result by rid
func (d *Dao) ResultByRID(c context.Context, rid int64) (res *resource.Result, err error) {
	res = &resource.Result{}
	if err = d.orm.Where("rid=?", rid).First(&res).Error; err == gorm.ErrRecordNotFound {
		log.Error("ResultByRID rid(%d) error(%v)", rid, err)
	}
	return
}

//MetaByRID .
func (d *Dao) MetaByRID(c context.Context, rids []int64) (metas map[int64]string, err error) {
	var rows *sql.Rows
	metas = make(map[int64]string)
	rows, err = d.orm.Table("resource").Select("id,metadata").Where("id IN (?)", rids).Rows()
	if err != nil {
		log.Error("MetaByRID Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id   int64
			meta string
		)
		if err = rows.Scan(&id, &meta); err != nil {
			log.Error("MetaByRID rows.Scan error(%v)", err)
			return
		}
		metas[id] = meta
	}
	return
}

//UpsertByOIDs 获取需要更新到搜索部分
func (d *Dao) UpsertByOIDs(c context.Context, businessID int64, oids []string) (res []*model.UpsertItem, err error) {
	res = []*model.UpsertItem{}
	sqlstr := `SELECT r.id,r.extra1,r.extra2,r.extra3,r.extra4,rr.state
FROM resource r
LEFT JOIN resource_result rr ON r.id=rr.rid
WHERE r.business_id=? AND r.oid IN (?)`
	if err = d.orm.Raw(sqlstr, businessID, oids).Find(&res).Error; err != nil {
		log.Error("UpsertByOIDs error(%+v) businessid(%d) oids(%+v)", err, businessID, oids)
	}

	return
}
