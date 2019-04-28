package gorm

import (
	"context"
	"strings"

	"github.com/jinzhu/gorm"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"
	"go-common/library/log"
)

// FlowByID .
func (d *Dao) FlowByID(c context.Context, id int64) (f *net.Flow, err error) {
	f = &net.Flow{}
	err = d.orm.Where("id=?", id).First(f).Error
	if err == gorm.ErrRecordNotFound {
		err = ecode.AegisFlowNotFound
		return
	}
	if err != nil {
		log.Error("FlowByID(%+v) error(%v)", id, err)
	}
	return
}

// FlowList .
func (d *Dao) FlowList(c context.Context, pm *net.ListNetElementParam) (result *net.ListFlowRes, err error) {
	result = &net.ListFlowRes{
		Pager: net.Pager{
			Num:  pm.Pn,
			Size: pm.Ps,
		},
	}

	db := d.orm.Table(net.TableFlow).Where("net_id=?", pm.NetID)
	if len(pm.ID) > 0 {
		db = db.Where("id in (?)", pm.ID)
	}
	pm.Name = strings.TrimSpace(pm.Name)
	if pm.Name != "" {
		db = db.Where("name=?", pm.Name)
	}

	err = db.Scopes(state(pm.State)).Count(&result.Pager.Total).Scopes(pager(pm.Ps, pm.Pn, pm.Sort)).Find(&result.Result).Error
	if err != nil {
		log.Error("FlowList find error(%v) params(%+v)", err, pm)
	}
	return
}

// FlowByUnique .
func (d *Dao) FlowByUnique(c context.Context, netID int64, name string) (f *net.Flow, err error) {
	f = &net.Flow{}
	err = d.orm.Where("net_id=? AND name=?", netID, name).First(f).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		f = nil
		return
	}
	if err != nil {
		log.Error("FlowByUnique(%d,%s) error(%v)", netID, name, err)
	}
	return
}

// Flows .
func (d *Dao) Flows(c context.Context, ids []int64) (fs []*net.Flow, err error) {
	fs = []*net.Flow{}
	err = d.orm.Where("id in (?)", ids).Find(&fs).Error
	if err != nil {
		log.Error("Flows(%v) error(%v)", ids, err)
	}

	return
}

// FlowsByNet .
func (d *Dao) FlowsByNet(c context.Context, netID []int64) (fs []*net.Flow, err error) {
	fs = []*net.Flow{}
	err = d.orm.Scopes(Available).Where("net_id in (?)", netID).Order("id ASC").Find(&fs).Error
	if err != nil {
		log.Error("Flows(%d) error(%v)", netID, err)
	}
	return
}

func (d *Dao) FlowIDByNet(c context.Context, nid []int64) (res map[int64][]int64, err error) {
	res = map[int64][]int64{}
	listi := []struct {
		ID    int64 `gorm:"column:id"`
		NetID int64 `gorm:"column:net_id"`
	}{}

	if err = d.orm.Table(net.TableFlow).Select("id,net_id").
		Where("net_id in (?)", nid).Scopes(Available).Scan(&listi).Error; err != nil {
		log.Error("FlowIDByNet error(%v) nid(%v)", err, nid)
		return
	}

	for _, item := range listi {
		res[item.NetID] = append(res[item.NetID], item.ID)
	}
	return
}
