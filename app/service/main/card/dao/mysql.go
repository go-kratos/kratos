package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/card/model"
	xsql "go-common/library/database/sql"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_selUserEquip          = "SELECT mid,card_id,expire_time FROM user_card_equip WHERE mid=? AND deleted = 0;"
	_selUserEquips         = "SELECT mid,card_id,expire_time FROM user_card_equip WHERE mid IN (%s) AND deleted = 0;"
	_selEffectiveCard      = "SELECT id,name,state,is_hot,card_url,big_crad_url,card_type,order_num,group_id,ctime,mtime FROM card_info WHERE state = 0 AND deleted = 0 ORDER BY order_num DESC;"
	_selEffectiveCardGroup = "SELECT id,name,state,ctime,mtime,order_num FROM card_group WHERE state = 0 AND deleted = 0;"
	_cardEquip             = "INSERT INTO user_card_equip(mid,card_id,expire_time)VALUES(?,?,?) ON DUPLICATE KEY UPDATE card_id =?,expire_time=?,deleted=0;"
	_cardDemount           = "UPDATE user_card_equip SET deleted = 1 WHERE mid = ?;"
)

// RawEquip get user equip info.
func (d *Dao) RawEquip(c context.Context, mid int64) (r *model.UserEquip, err error) {
	r = new(model.UserEquip)
	row := d.db.QueryRow(c, _selUserEquip, mid)
	if err = row.Scan(&r.Mid, &r.CardID, &r.ExpireTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			r = nil
			return
		}
		r = nil
		err = errors.Wrapf(err, "dao equip mid(%d)", mid)
	}
	return
}

// RawEquips get user equip infos.
func (d *Dao) RawEquips(c context.Context, mids []int64) (res map[int64]*model.UserEquip, err error) {
	var rows *xsql.Rows
	res = make(map[int64]*model.UserEquip, len(mids))
	midStr := xstr.JoinInts(mids)
	if rows, err = d.db.Query(c, fmt.Sprintf(_selUserEquips, midStr)); err != nil {
		err = errors.Wrapf(err, "dao equips mids(%s)", midStr)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.UserEquip)
		if err = rows.Scan(&r.Mid, &r.CardID, &r.ExpireTime); err != nil {
			err = errors.Wrapf(err, "dao equips scan mids(%s)", midStr)
			res = nil
			return
		}
		res[r.Mid] = r
	}
	err = rows.Err()
	return
}

// EffectiveCards query all effective cards .
func (d *Dao) EffectiveCards(c context.Context) (res []*model.Card, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _selEffectiveCard); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Card)
		if err = rows.Scan(&r.ID, &r.Name, &r.State, &r.IsHot, &r.CardURL, &r.BigCradURL, &r.CardType, &r.OrderNum,
			&r.GroupID, &r.Ctime, &r.Mtime); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// EffectiveGroups query all effective groups .
func (d *Dao) EffectiveGroups(c context.Context) (res []*model.CardGroup, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _selEffectiveCardGroup); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.CardGroup)
		if err = rows.Scan(&r.ID, &r.Name, &r.State, &r.Ctime, &r.Mtime, &r.OrderNum); err != nil {
			err = errors.WithStack(err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// CardEquip card equip.
func (d *Dao) CardEquip(c context.Context, e *model.UserEquip) (err error) {
	if _, err = d.db.Exec(c, _cardEquip, e.Mid, e.CardID, e.ExpireTime, e.CardID, e.ExpireTime); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// DeleteEquip delete card equip.
func (d *Dao) DeleteEquip(c context.Context, mid int64) (err error) {
	if _, err = d.db.Exec(c, _cardDemount, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
