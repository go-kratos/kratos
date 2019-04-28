package dao

import (
	"context"
	"fmt"
	"strconv"

	"go-common/app/admin/main/tv/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

const (
	_setUploader = "REPLACE INTO ugc_uploader (mid, state) VALUES (?,?)"
	_auditMids   = "UPDATE ugc_uploader SET valid = ? WHERE mid IN (?) AND deleted = 0"
	_intervUp    = "UPDATE ugc_uploader SET cms_face = ?, cms_name = ? WHERE mid = ? AND deleted = 0"
	_stateNormal = 1
	_deleted     = 1
	// up list oder type
	_ctimeNew = 1
	_ctimeOld = 2
	_mtimeNew = 3
	_mtimeOld = 4
	size      = 20
)

// UpAdd def.
func (d *Dao) UpAdd(mid int64) (err error) {
	if err = d.DB.Exec(_setUploader, mid, _stateNormal).Error; err != nil {
		log.Error("UpAdd Error %v", err)
	}
	return
}

// UpList def.
func (d *Dao) UpList(order int, page int, ids []int64) (ups []*model.Upper, pager *model.Page, err error) {
	var db = d.DB.Where("deleted != ?", _deleted)
	if len(ids) != 0 {
		db = db.Where("mid IN (?)", ids)
	}
	pager = &model.Page{
		Num:  page,
		Size: size,
	}
	if err = db.Model(&model.Upper{}).Count(&pager.Total).Error; err != nil {
		log.Error("Uplist Count Error %v", err)
		return
	}
	db = treatOrder(order, db)
	if err = db.Offset((page - 1) * size).Limit(size).Find(&ups).Error; err != nil {
		log.Error("UpList Error %v, Order: %d", err, order)
	}
	return
}

func treatOrder(order int, db *gorm.DB) *gorm.DB {
	switch order {
	case _ctimeNew:
		db = db.Order("ctime DESC")
	case _ctimeOld:
		db = db.Order("ctime ASC")
	case _mtimeNew:
		db = db.Order("mtime DESC")
	case _mtimeOld:
		db = db.Order("mtime ASC")
	default:
		db = db.Order("ctime DESC")
	}
	return db
}

// UpCmsList def.
func (d *Dao) UpCmsList(req *model.ReqUpCms) (ups []*model.CmsUpper, pager *model.Page, err error) {
	var db = d.DB.Where("deleted = 0")
	if req.MID != 0 {
		db = db.Where("mid = ?", req.MID)
	}
	if req.Name != "" {
		db = db.Where("ori_name LIKE ?", "%"+req.Name+"%")
	}
	if req.Valid != "" {
		valid, _ := strconv.Atoi(req.Valid)
		db = db.Where("valid = ?", valid)
	}
	pager = &model.Page{
		Num:  req.Pn,
		Size: size,
	}
	if err = db.Model(&model.Upper{}).Count(&pager.Total).Error; err != nil {
		log.Error("Uplist Count Error %v", err)
		return
	}
	db = treatOrder(req.Order, db)
	if err = db.Offset((req.Pn - 1) * size).Limit(size).Find(&ups).Error; err != nil {
		log.Error("UpList Error %v, Order: %d", err, req.Order)
	}
	return
}

// VerifyIds verifies whether all the mids exist
func (d *Dao) VerifyIds(mids []int64) (okMids map[int64]*model.UpMC, err error) {
	if len(mids) == 0 {
		return
	}
	okMids = make(map[int64]*model.UpMC)
	var ups []*model.UpMC
	db := d.DB.Where("deleted != ?", _deleted).Where("mid IN (?)", mids)
	if err = db.Find(&ups).Error; err != nil {
		log.Error("VerifyIds Error %v, Mids %v", err, mids)
	}
	for _, v := range ups {
		okMids[v.MID] = v
	}
	return
}

// AuditIds carries out the action to the given mids
func (d *Dao) AuditIds(mids []int64, validAct int) (err error) {
	if err = d.DB.Exec(_auditMids, validAct, mids).Error; err != nil {
		log.Error("AuditIds Error %v, Mids %v", err, mids)
	}
	return
}

// SetUpper updates the cms info of an upper in DB
func (d *Dao) SetUpper(req *model.ReqUpEdit) (err error) {
	if err = d.DB.Exec(_intervUp, req.Face, req.Name, req.MID).Error; err != nil {
		log.Error("SetUpper Error %v, Mid %v", err, req)
	}
	return
}

func upperMetaKey(MID int64) string {
	return fmt.Sprintf("up_cms_%d", MID)
}

// SetUpMetaCache updates upinfo in MC
func (d *Dao) SetUpMetaCache(c context.Context, upper *model.UpMC) (err error) {
	var (
		key  = upperMetaKey(upper.MID)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: key, Object: upper, Flags: memcache.FlagJSON, Expiration: d.cmsExpire}); err != nil {
		log.Error("conn.Set error(%v)", err)
	}
	return
}

// DelCache deletes the cache from CM
func (d *Dao) DelCache(c context.Context, mid int64) (err error) {
	var (
		key  = upperMetaKey(mid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if err = conn.Delete(key); err != nil {
		log.Error("conn.Set error(%v)", err)
	}
	return
}
