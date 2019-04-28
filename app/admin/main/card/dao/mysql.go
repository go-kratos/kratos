package dao

import (
	"bytes"
	"context"
	"strconv"

	"go-common/app/admin/main/card/model"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

// AddGroup add group info.
func (d *Dao) AddGroup(c context.Context, arg *model.AddGroup) error {
	return d.DB.Table("card_group").Create(arg).Error
}

// UpdateGroup update group info.
func (d *Dao) UpdateGroup(c context.Context, arg *model.UpdateGroup) error {
	return d.DB.Table("card_group").Where("id=?", arg.ID).Update(map[string]interface{}{
		"name":     arg.Name,
		"state":    arg.State,
		"operator": arg.Operator,
	}).Error
}

// Cards query cards.
func (d *Dao) Cards(c context.Context) (res []*model.Card, err error) {
	err = d.DB.Table("card_info").Order("order_num desc").Where("state = 0 AND deleted = 0").Find(&res).Error
	return
}

// CardsByGid query cards by group id.
func (d *Dao) CardsByGid(c context.Context, gid int64) (res []*model.Card, err error) {
	err = d.DB.Table("card_info").Order("order_num desc").Where("group_id=?", gid).Where("deleted = 0").Find(&res).Error
	return
}

// CardsByIds query cards by ids.
func (d *Dao) CardsByIds(c context.Context, ids []int64) (res []*model.Card, err error) {
	err = d.DB.Table("card_info").Order("order_num desc").Where("id in (?)", ids).Where("deleted = 0").Find(&res).Error
	return
}

// GroupsByIds query groups by ids.
func (d *Dao) GroupsByIds(c context.Context, ids []int64) (res []*model.CardGroup, err error) {
	err = d.DB.Table("card_group").Order("order_num desc").Where("id in (?)", ids).Where("deleted = 0").Find(&res).Error
	return
}

// AddCard add card.
func (d *Dao) AddCard(arg *model.AddCard) error {
	return d.DB.Table("card_info").Create(arg).Error
}

// CardByName get card by name.
func (d *Dao) CardByName(name string) (res *model.Card, err error) {
	res = new(model.Card)
	q := d.DB.Table("card_info").Where("name=?", name).First(res)
	if q.Error != nil {
		if q.RecordNotFound() {
			err = nil
			res = nil
			return
		}
		err = errors.Wrapf(err, "card by name")
	}
	return
}

// GroupByName get group by name.
func (d *Dao) GroupByName(name string) (res *model.CardGroup, err error) {
	res = new(model.CardGroup)
	q := d.DB.Table("card_group").Where("name=?", name).First(res)
	if q.Error != nil {
		if q.RecordNotFound() {
			err = nil
			res = nil
			return
		}
		err = errors.Wrapf(err, "card_group by name")
	}
	return
}

// UpdateCard update card.
func (d *Dao) UpdateCard(req *model.UpdateCard) error {
	args := map[string]interface{}{}
	args["name"] = req.Name
	args["state"] = req.State
	args["is_hot"] = req.IsHot
	args["operator"] = req.Operator
	if req.CardURL != "" {
		args["card_url"] = req.CardURL
	}
	if req.BigCradURL != "" {
		args["big_crad_url"] = req.BigCradURL
	}
	return d.DB.Table("card_info").Where("id=?", req.ID).Update(args).Error
}

// UpdateCardState update card state.
func (d *Dao) UpdateCardState(c context.Context, id int64, state int8) error {
	return d.DB.Table("card_info").Where("id=?", id).Update("state", state).Error
}

// DeleteCard delete card.
func (d *Dao) DeleteCard(c context.Context, id int64) error {
	return d.DB.Table("card_info").Where("id=?", id).Delete(&model.Card{}).Error
}

// DeleteGroup delete group.
func (d *Dao) DeleteGroup(c context.Context, id int64) error {
	return d.DB.Table("card_group").Where("id=?", id).Delete(&model.CardGroup{}).Error
}

// UpdateGroupState update group state.
func (d *Dao) UpdateGroupState(c context.Context, id int64, state int8) error {
	return d.DB.Table("card_group").Where("id=?", id).Update("state", state).Error
}

// MaxCardOrder max card order num.
func (d *Dao) MaxCardOrder() (max int64, err error) {
	err = d.DB.Table("card_info").Select("MAX(order_num)").Row().Scan(&max)
	return
}

// MaxGroupOrder max card group order num.
func (d *Dao) MaxGroupOrder() (max int64, err error) {
	err = d.DB.Table("card_group").Select("MAX(order_num)").Row().Scan(&max)
	return
}

// BatchUpdateCardOrder update card order.
func (d *Dao) BatchUpdateCardOrder(c context.Context, cs []*model.Card) error {
	var (
		buf bytes.Buffer
		ids []int64
	)
	buf.WriteString("UPDATE card_info SET order_num = CASE id")
	for _, v := range cs {
		buf.WriteString(" WHEN ")
		buf.WriteString(strconv.FormatInt(v.ID, 10))
		buf.WriteString(" THEN ")
		buf.WriteString(strconv.FormatInt(v.OrderNum, 10))
		ids = append(ids, v.ID)
	}
	buf.WriteString(" END  WHERE id IN (")
	buf.WriteString(xstr.JoinInts(ids))
	buf.WriteString(");")
	return d.DB.Exec(buf.String()).Error
}

// BatchUpdateCardGroupOrder update card order.
func (d *Dao) BatchUpdateCardGroupOrder(c context.Context, cs []*model.CardGroup) error {
	var (
		buf bytes.Buffer
		ids []int64
	)
	buf.WriteString("UPDATE card_group SET order_num = CASE id")
	for _, v := range cs {
		buf.WriteString(" WHEN ")
		buf.WriteString(strconv.FormatInt(v.ID, 10))
		buf.WriteString(" THEN ")
		buf.WriteString(strconv.FormatInt(v.OrderNum, 10))
		ids = append(ids, v.ID)
	}
	buf.WriteString(" END  WHERE id IN (")
	buf.WriteString(xstr.JoinInts(ids))
	buf.WriteString(");")
	return d.DB.Exec(buf.String()).Error
}

// Groups query groups.
func (d *Dao) Groups(c context.Context, arg *model.ArgQueryGroup) (res []*model.CardGroup, err error) {
	q := d.DB.Table("card_group").Where("deleted = 0")
	if arg.GroupID > 0 {
		q = q.Where("id = ?", arg.GroupID)
	}
	if arg.State > -1 {
		q = q.Where("state = ?", arg.State)
	}
	if err = q.Order("order_num desc").Find(&res).Error; err != nil {
		if q.RecordNotFound() {
			err = nil
			return
		}
		err = errors.Wrapf(err, "card group list")
		return
	}
	return
}
