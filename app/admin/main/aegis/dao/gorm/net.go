package gorm

import (
	"context"
	"time"

	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

// NetByID .
func (d *Dao) NetByID(c context.Context, id int64) (n *net.Net, err error) {
	n = &net.Net{}
	err = d.orm.Where("id=?", id).First(n).Error
	if err == gorm.ErrRecordNotFound {
		err = ecode.NothingFound
		return
	}
	if err != nil {
		log.Error("NetByID(%+v) error(%v)", id, err)
	}
	return
}

func (d *Dao) Nets(c context.Context, ids []int64) (n []*net.Net, err error) {
	n = []*net.Net{}
	if err = d.orm.Where("id in (?)", ids).Find(&n).Error; err != nil {
		log.Error("Nets error(%v) ids(%v)", err, ids)
	}
	return
}

// NetByUnique .
func (d *Dao) NetByUnique(c context.Context, name string) (n *net.Net, err error) {
	n = &net.Net{}
	err = d.orm.Where("ch_name=?", name).First(n).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		n = nil
		return
	}
	if err != nil {
		log.Error("NetByUnique(%+v) error(%v)", name, err)
	}
	return
}

// NetList .
func (d *Dao) NetList(c context.Context, pm *net.ListNetParam) (result *net.ListNetRes, err error) {
	result = &net.ListNetRes{
		Pager: net.Pager{
			Num:  pm.Pn,
			Size: pm.Ps,
		},
	}
	db := d.orm.Table(net.TableNet)
	if len(pm.ID) > 0 {
		db = db.Where("id in (?)", pm.ID)
	}
	if pm.BusinessID > 0 {
		db = db.Where("business_id=?", pm.BusinessID)
	}
	err = db.Scopes(state(pm.State)).Count(&result.Pager.Total).Scopes(pager(pm.Ps, pm.Pn, pm.Sort)).Find(&result.Result).Error
	if err != nil {
		log.Error("NetList find error(%v) params(%+v)", err, pm)
	}
	return
}

// NetBindStartFlow .
func (d *Dao) NetBindStartFlow(c context.Context, tx *gorm.DB, id int64, flowID int64) (err error) {
	if err = d.UpdateFields(c, tx, net.TableNet, id, map[string]interface{}{"start_flow_id": flowID}); err != nil {
		log.Error("NetBindStartFlow d.UpdateFields error(%v) id(%d) flowid(%d)", err, id, flowID)
	}
	return
}

// NetIDByBusiness .
func (d *Dao) NetIDByBusiness(c context.Context, businessID []int64) (bizmap map[int64][]int64, err error) {
	res := []struct {
		ID         int64 `gorm:"column:id"`
		BusinessID int64 `gorm:"column:business_id"`
	}{}
	list := []*net.Net{}
	bizmap = map[int64][]int64{}
	if err = d.orm.Select("id, business_id").Where("business_id in (?)", businessID).
		Scopes(Available).Find(&list).Scan(&res).Error; err != nil {
		return
	}

	for _, item := range res {
		bizmap[item.BusinessID] = append(bizmap[item.BusinessID], item.ID)
	}
	return
}

// NetsByBusiness .
func (d *Dao) NetsByBusiness(c context.Context, businessID []int64, onlyAvailable bool) (list []*net.Net, err error) {
	list = []*net.Net{}
	db := d.orm
	if len(businessID) > 0 {
		db = db.Where("business_id in (?)", businessID)
	}
	if onlyAvailable {
		db = db.Scopes(Available)
	}
	if err = db.Find(&list).Error; err != nil {
		log.Error("NetsByBusiness(%v) error(%v) onlyavailable(%v)", businessID, err, onlyAvailable)
	}
	return
}

// DisableNet .
func (d *Dao) DisableNet(c context.Context, tx *gorm.DB, id int64) (err error) {
	err = d.UpdateFields(c, tx, net.TableNet, id, map[string]interface{}{"disable_time": time.Now()})
	return
}
