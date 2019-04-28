package gorm

import (
	"context"

	"github.com/jinzhu/gorm"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"
	"go-common/library/log"
)

// DirectionByFlowID .
func (d *Dao) DirectionByFlowID(c context.Context, flowID []int64, direction int8) (dirs []*net.Direction, err error) {
	dirs = []*net.Direction{}
	db := d.orm.Where("flow_id in (?)", flowID)
	if direction == net.DirInput || direction == net.DirOutput {
		db = db.Where("direction=?", direction)
	}
	err = db.Scopes(Available).Find(&dirs).Error
	if err != nil {
		log.Error("DirectionByFlowID find error(%v) flowid(%v) direction(%d)", err, flowID, direction)
	}
	return
}

// DirectionByTransitionID .
func (d *Dao) DirectionByTransitionID(c context.Context, transitionID []int64, direction int8, onlyAvailable bool) (dirs []*net.Direction, err error) {
	dirs = []*net.Direction{}
	db := d.orm.Where("transition_id in (?)", transitionID)
	if direction == net.DirInput || direction == net.DirOutput {
		db = db.Where("direction=?", direction)
	}
	if onlyAvailable {
		db = db.Scopes(Available)
	}
	err = db.Find(&dirs).Error
	if err != nil {
		log.Error("DirectionByTransitionID find error(%v) transitionid(%v) direction(%d)", err, transitionID, direction)
	}
	return
}

// DirectionByID .
func (d *Dao) DirectionByID(c context.Context, id int64) (n *net.Direction, err error) {
	n = &net.Direction{}
	err = d.orm.Where("id=?", id).First(n).Error
	if err == gorm.ErrRecordNotFound {
		err = ecode.NothingFound
		return
	}
	if err != nil {
		log.Error("DirectionByID(%+v) error(%v)", id, err)
	}

	return
}

func (d *Dao) Directions(c context.Context, ids []int64) (n []*net.Direction, err error) {
	n = []*net.Direction{}
	if err = d.orm.Where("id in (?)", ids).Find(&n).Error; err != nil {
		log.Error("Directions error(%v) ids(%v)", err, ids)
	}

	return
}

// DirectionList .
func (d *Dao) DirectionList(c context.Context, pm *net.ListDirectionParam) (result *net.ListDirectionRes, err error) {
	result = &net.ListDirectionRes{
		Pager: net.Pager{
			Num:  pm.Pn,
			Size: pm.Ps,
		},
	}
	db := d.orm.Table(net.TableDirection).Where("net_id=?", pm.NetID)
	if len(pm.ID) > 0 {
		db = db.Where("id in (?)", pm.ID)
	}
	if pm.FlowID > 0 {
		db = db.Where("flow_id=?", pm.FlowID)
	}
	if pm.TransitionID > 0 {
		db = db.Where("transition_id=?", pm.TransitionID)
	}
	if pm.Direction == net.DirInput || pm.Direction == net.DirOutput {
		db = db.Where("direction=?", pm.Direction)
	}

	err = db.Scopes(state(pm.State)).Count(&result.Pager.Total).Scopes(pager(pm.Ps, pm.Pn, pm.Sort)).Find(&result.Result).Error
	if err != nil {
		log.Error("DirectionList find error(%v) params(%+v)", err, pm)
	}
	return
}

// DirectionByUnique .
func (d *Dao) DirectionByUnique(c context.Context, netID int64, flowID int64, transitionID int64, direction int8) (t *net.Direction, err error) {
	t = &net.Direction{}
	err = d.orm.Where("net_id=? AND flow_id=? AND transition_id=? AND direction=?", netID, flowID, transitionID, direction).
		First(t).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		t = nil
		return
	}
	if err != nil {
		log.Error("DirectionByUnique(%d,%d,%d,%d) error(%v)", netID, flowID, transitionID, direction, err)
	}

	return
}

// DirectionByNet .
func (d *Dao) DirectionByNet(c context.Context, netID int64) (n []*net.Direction, err error) {
	n = []*net.Direction{}
	err = d.orm.Where("net_id=?", netID).Scopes(Available).Find(&n).Error
	if err != nil {
		log.Error("DirectionByNet(%d) error(%v)", netID, err)
	}

	return
}
