package gorm

import (
	"context"
	"strings"

	"github.com/jinzhu/gorm"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"
	"go-common/library/log"
)

// TokenByID .
func (d *Dao) TokenByID(c context.Context, id int64) (t *net.Token, err error) {
	t = &net.Token{}
	err = d.orm.Where("id=?", id).First(t).Error
	if err == gorm.ErrRecordNotFound {
		err = ecode.AegisTokenNotFound
		return
	}
	if err != nil {
		log.Error("TokenByID(%+v) error(%v)", id, err)
	}
	return
}

// Tokens .
func (d *Dao) Tokens(c context.Context, ids []int64) (list []*net.Token, err error) {
	list = []*net.Token{}
	err = d.orm.Where("id in (?)", ids).Find(&list).Error
	if err != nil {
		log.Error("Tokens(%v) error(%v)", ids, err)
	}

	return
}

func (d *Dao) tokenListDB(netID []int64, id []int64, name string, onlyAssign bool) (db *gorm.DB) {
	db = d.orm.Table(net.TableToken).Where("net_id in (?)", netID)
	if len(id) > 0 {
		db = db.Where("id in (?)", id)
	}
	name = strings.TrimSpace(name)
	if name != "" {
		db = db.Where("name=?", name)
	}
	if onlyAssign {
		db = db.Where("compare=?", net.TokenCompareAssign)
	}
	return
}

// TokenListWithPager .
func (d *Dao) TokenListWithPager(c context.Context, pm *net.ListTokenParam) (result *net.ListTokenRes, err error) {
	result = &net.ListTokenRes{
		Pager: net.Pager{
			Num:  pm.Pn,
			Size: pm.Ps,
		},
	}
	db := d.tokenListDB([]int64{pm.NetID}, pm.ID, pm.Name, pm.Assign)
	err = db.Count(&result.Pager.Total).Scopes(pager(pm.Ps, pm.Pn, pm.Sort)).Find(&result.Result).Error
	if err != nil {
		log.Error("TokenListWithPager find error(%v) params(%+v)", err, pm)
	}
	return
}

// TokenList .
func (d *Dao) TokenList(c context.Context, netID []int64, id []int64, name string, onlyAssign bool) (list []*net.Token, err error) {
	err = d.tokenListDB(netID, id, name, onlyAssign).Find(&list).Error
	if err != nil {
		log.Error("TokenList find error(%v)", err)
	}
	return
}

// TokenByUnique .
func (d *Dao) TokenByUnique(c context.Context, netID int64, name string, compare int8, value string) (t *net.Token, err error) {
	t = &net.Token{}
	err = d.orm.Where("net_id=? AND name=? AND compare=? AND value=?", netID, name, compare, value).First(t).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		t = nil
		return
	}
	if err != nil {
		log.Error("TokenByUnique(%d,%s,%d,%s) error(%v)", netID, name, compare, value, err)
	}
	return
}

// TokenBinds .
func (d *Dao) TokenBinds(c context.Context, id []int64) (t []*net.TokenBind, err error) {
	t = []*net.TokenBind{}
	err = d.orm.Where("id in (?)", id).Find(&t).Error
	if err != nil {
		log.Error("TokenBinds(%+v) error(%v)", id, err)
	}
	return
}

// TokenBindByElement .
func (d *Dao) TokenBindByElement(c context.Context, elementID []int64, tp []int8, onlyAvailable bool) (binds map[int64][]*net.TokenBind, err error) {
	binds = map[int64][]*net.TokenBind{}
	list := []*net.TokenBind{}
	db := d.orm.Where("element_id in (?) AND type in (?)", elementID, tp)
	if onlyAvailable {
		db = db.Scopes(Available)
	}

	if err = db.Find(&list).Error; err != nil {
		log.Error("TokenBindByElement error(%v) elementid(%d) type(%v) onlyavailable(%v)", err, elementID, tp, onlyAvailable)
		return
	}

	for _, item := range list {
		if _, exist := binds[item.ElementID]; !exist {
			binds[item.ElementID] = []*net.TokenBind{item}
			continue
		}

		binds[item.ElementID] = append(binds[item.ElementID], item)
	}
	return
}
