package gorm

import (
	"context"
	"strings"

	"github.com/jinzhu/gorm"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"
	"go-common/library/log"
)

// TransitionByID .
func (d *Dao) TransitionByID(c context.Context, id int64) (n *net.Transition, err error) {
	n = &net.Transition{}
	err = d.orm.Where("id=?", id).First(n).Error
	if err == gorm.ErrRecordNotFound {
		err = ecode.AegisTranNotFound
		return
	}
	if err != nil {
		log.Error("TransitionByID(%d) error(%v)", id, err)
	}
	return
}

// Transitions .
func (d *Dao) Transitions(c context.Context, id []int64) (n []*net.Transition, err error) {
	n = []*net.Transition{}
	err = d.orm.Where("id in (?)", id).Find(&n).Error
	if err != nil {
		log.Error("Transitions(%v) error(%v)", id, err)
	}
	return
}

// TransitionList .
func (d *Dao) TransitionList(c context.Context, pm *net.ListNetElementParam) (result *net.ListTransitionRes, err error) {
	result = &net.ListTransitionRes{
		Pager: net.Pager{
			Num:  pm.Pn,
			Size: pm.Ps,
		},
	}
	db := d.orm.Table(net.TableTransition).Where("net_id=?", pm.NetID)
	if len(pm.ID) > 0 {
		db = db.Where("id in (?)", pm.ID)
	}
	pm.Name = strings.TrimSpace(pm.Name)
	if pm.Name != "" {
		db = db.Where("name=?", pm.Name)
	}
	err = db.Scopes(state(pm.State)).Count(&result.Pager.Total).Scopes(pager(pm.Ps, pm.Pn, pm.Sort)).Find(&result.Result).Error
	if err != nil {
		log.Error("TransitionList find error(%v) params(%+v)", err, pm)
	}
	return
}

// TransitionByUnique .
func (d *Dao) TransitionByUnique(c context.Context, netID int64, name string) (t *net.Transition, err error) {
	t = &net.Transition{}
	err = d.orm.Where("net_id=? AND name=?", netID, name).First(t).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		t = nil
		return
	}
	if err != nil {
		log.Error("TransitionByUnique(%d,%s) error(%v)", netID, name, err)
	}
	return
}

// TransitionIDByNet .
func (d *Dao) TransitionIDByNet(c context.Context, netID []int64, onlyDispatch bool, onlyAvailable bool) (ids map[int64][]int64, err error) {
	ids = map[int64][]int64{}
	listi := []struct {
		ID    int64 `gorm:"column:id"`
		NetID int64 `gorm:"column:net_id"`
	}{}
	db := d.orm.Table(net.TableTransition).Where("net_id in (?)", netID)
	if onlyAvailable {
		db = db.Scopes(Available)
	}
	if onlyDispatch {
		db = db.Where("`limit`>0").Scopes(manual)
	}
	if err = db.Find(&listi).Error; err != nil {
		log.Error("TransitionIDByNet netid(%v) error(%v)", netID, err)
		return
	}

	for _, item := range listi {
		ids[item.NetID] = append(ids[item.NetID], item.ID)
	}
	return
}

func manual(db *gorm.DB) *gorm.DB {
	return db.Where("`trigger`=?", net.TriggerManual)
}

// TranByNet .
func (d *Dao) TranByNet(c context.Context, netID int64, onlyAvailable bool) (list []*net.Transition, err error) {
	list = []*net.Transition{}
	db := d.orm
	if netID > 0 {
		db = db.Where("net_id=?", netID)
	}
	if onlyAvailable {
		db = db.Scopes(Available)
	}
	if err = db.Find(&list).Error; err != nil {
		log.Error("TranByNet error(%v)", err)
	}
	return
}
