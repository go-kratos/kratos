package gorm

import (
	"context"

	"github.com/jinzhu/gorm"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/log"
)

// FRByFlow .
func (d *Dao) FRByFlow(c context.Context, flowID []int64) (fr *net.FlowResource, err error) {
	fr = &net.FlowResource{}
	if err = d.orm.Where("flow_id in (?)", flowID).First(fr).Error; err == gorm.ErrRecordNotFound {
		fr = nil
		err = nil
		return
	}

	if err != nil {
		log.Error("FRByFlow error(%v) flowid(%v)", err, flowID)
	}
	return
}

// FRByNetRID .
func (d *Dao) FRByNetRID(c context.Context, netID []int64, rids []int64, onlyRunning bool) (frs []*net.FlowResource, err error) {
	frs = []*net.FlowResource{}
	db := d.orm.Where("rid in (?) AND net_id in (?)", rids, netID)
	if onlyRunning {
		db = db.Scopes(running)
	}
	if err = db.Find(&frs).Error; err != nil {
		log.Error("FRByNetRID error(%v) netid(%v) rids(%v)", err, netID, rids)
	}
	return
}

// FRByUniques .
func (d *Dao) FRByUniques(c context.Context, rids []int64, flowID []int64, onlyRunning bool) (frs []*net.FlowResource, err error) {
	frs = []*net.FlowResource{}
	db := d.orm
	if len(rids) > 0 {
		db = db.Where("rid in (?)", rids)
	}
	if len(flowID) > 0 {
		db = db.Where("flow_id in (?)", flowID)
	}
	if onlyRunning {
		db = db.Scopes(running)
	}
	if err = db.Find(&frs).Error; err != nil {
		log.Error("FRByUniques error(%v) rids(%+v) flowid(%d)", err, rids, flowID)
	}

	return
}

// CancelFlowResource .
func (d *Dao) CancelFlowResource(c context.Context, tx *gorm.DB, rids []int64) (err error) {
	fields := map[string]interface{}{"state": net.FRStateDeleted}
	if err = tx.Table(net.TableFlowResource).Where("rid in (?)", rids).Updates(fields).Error; err != nil {
		log.Error("CancelFlowResource error(%v) rids(%+v)", err, rids)
	}
	return
}

func running(db *gorm.DB) *gorm.DB {
	return db.Where("state!=?", net.FRStateDeleted)
}
