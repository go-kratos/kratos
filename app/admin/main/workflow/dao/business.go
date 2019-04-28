package dao

import (
	"context"
	"encoding/json"

	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/manager"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const _businessRoleURI = "http://manager.bilibili.co/x/admin/manager/internal/business/role"

var metas map[int8]*model.Meta

func init() {
	metas = make(map[int8]*model.Meta)
	data := `[
	{
		"business": 1,
		"name": "稿件投诉",
		"item_type": "group",
		"rounds": [
			{
				"id": 1,
				"name": "一审"
			},
			{
				"id": 2,
				"name": "回查"
			}
		]
	},
	{
		"business": 2,
		"name": "稿件申诉",
		"item_type": "challenge",
		"rounds": [
			{
				"id": 1,
				"name": "一审"
			},
			{
				"id": 2,
				"name": "回查"
			},
			{
				"id": 3,
				"name": "三查"
			},
			{
				"id": 11,
				"name": "客服"
			}
		]
	},
	{
		"business": 3,
		"name": "短点评投诉",
		"item_type": "group",
		"rounds": [
			{
				"id": 1,
				"name": "一审"
			},
			{
				"id": 2,
				"name": "回查"
			}
		]
	},
	{
		"business": 4,
		"name": "长点评投诉",
		"item_type": "group",
		"rounds": [
			{
				"id": 1,
				"name": "一审"
			},
			{
				"id": 2,
				"name": "回查"
			}
		]
	},
	{
		"business": 5,
		"name": "小黑屋",
		"item_type": "challenge",
		"rounds": [
			{
				"id": 1,
				"name": "评论"
			},
			{
				"id": 2,
				"name": "弹幕"
			},
			{
				"id": 3,
				"name": "私信"
			},
			{
				"id": 4,
				"name": "标签"
			},
			{
				"id": 5,
				"name": "个人资料"
			},
			{
				"id": 6,
				"name": "投稿"
			},
			{
				"id": 7,
				"name": "音频"
			},
			{
				"id": 8,
				"name": "专栏"
			},
			{
				"id": 9,
				"name": "空间头图"
			}
		]
	},
	{
		"business": 6,
		"name": "稿件审核",
		"item_type": "challenge",
		"rounds": [
			{
				"id": 1,
				"name": "一审"
			}
		]
	},
	{
		"business": 7,
		"name": "任务质检",
		"item_type": "challenge",
		"rounds": [
			{
				"id": 1,
				"name": "一审"
			}
		]
	},
	{
		"business": 8,
		"name": "频道举报",
		"item_type": "group",
		"rounds": [
			{
				"id": 1,
				"name": "一审"
			}
		]
	}
]`
	ml := make([]*model.Meta, 0)
	err := json.Unmarshal([]byte(data), &ml)
	if err != nil {
		panic(err)
	}

	for _, m := range ml {
		metas[m.Business] = m
	}
}

// BatchLastBusRecIDs will retrive the last business record ids by serveral conditions
func (d *Dao) BatchLastBusRecIDs(c context.Context, oids []int64, business int8) (bids []int64, err error) {
	bids = make([]int64, 0, len(oids))
	if len(oids) <= 0 {
		return
	}
	rows, err := d.ReadORM.Table("workflow_business").Select("max(id)").
		Where("oid IN (?) AND business=?", oids, business).
		Group("oid,business").Rows()
	if err != nil {
		err = errors.Wrapf(err, "Query(%v, %d)", oids, business)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var bid int64
		if err = rows.Scan(&bid); err != nil {
			err = errors.WithStack(err)
			return
		}
		bids = append(bids, bid)
	}
	return
}

// BusinessRecs will retrive the business record by ids
func (d *Dao) BusinessRecs(c context.Context, bids []int64) (bs map[int32]*model.Business, err error) {
	bs = make(map[int32]*model.Business, len(bids))
	if len(bids) <= 0 {
		return
	}

	blist := make([]*model.Business, 0, len(bids))
	err = d.ReadORM.Table("workflow_business").Where("id IN (?)", bids).Find(&blist).Error
	if err != nil {
		err = errors.Wrapf(err, "Query(%v)", bids)
		return
	}

	for _, b := range blist {
		bs[b.Bid] = b
	}
	return
}

// LastBusRec will retrive last business record by business oid
func (d *Dao) LastBusRec(c context.Context, business int8, oid int64) (bs *model.Business, err error) {
	bs = new(model.Business)
	err = d.ReadORM.Table("workflow_business").Where("oid=? AND business=?", oid, business).Last(bs).Error
	if err != nil || bs.Bid == 0 {
		err = errors.Wrapf(err, "Query(%d, %d)", business, oid)
		bs = nil
		return
	}
	return
}

// BatchBusRecByCids will retrive businesses by cids
func (d *Dao) BatchBusRecByCids(c context.Context, cids []int64) (cidToBus map[int64]*model.Business, err error) {
	cidToBus = make(map[int64]*model.Business)
	if len(cids) <= 0 {
		return
	}

	blist := make([]*model.Business, 0, len(cids))
	err = d.ReadORM.Table("workflow_business").Where("cid IN (?)", cids).Find(&blist).Error
	if err != nil {
		err = errors.Wrapf(err, "Query(%v)", cids)
		return
	}

	for _, b := range blist {
		cidToBus[b.Cid] = b
	}
	return
}

// BusObjectByGids will retrive businesses by gids
func (d *Dao) BusObjectByGids(c context.Context, gids []int64) (gidToBus map[int64]*model.Business, err error) {
	gidToBus = make(map[int64]*model.Business, len(gids))
	if len(gids) <= 0 {
		return
	}
	blist := make([]*model.Business, 0, len(gids))
	if err = d.ReadORM.Table("workflow_business").Where("gid IN (?)", gids).Find(&blist).Error; err != nil {
		err = errors.Wrapf(err, "Query(%v)", gids)
		return
	}
	for _, b := range blist {
		gidToBus[b.Gid] = b
	}
	return
}

// AllMetas will retrive business meta infomation from pre-configured
func (d *Dao) AllMetas(c context.Context) map[int8]*model.Meta {
	return metas
}

// LoadRole .
func (d *Dao) LoadRole(c context.Context) (role map[int8]map[int8]string, err error) {
	var (
		resp *manager.RoleResponse
		ok   bool
		uri  = _businessRoleURI
	)
	role = make(map[int8]map[int8]string)
	if err = d.httpRead.Get(c, uri, "", nil, &resp); err != nil {
		log.Error("failed call %s error(%v)", uri, err)
		return
	}

	if resp.Code != ecode.OK.Code() {
		err = ecode.Int(resp.Code)
		log.Error("call %s error response code(%d) message(%s)", uri, resp.Code, resp.Message)
		return
	}
	for _, r := range resp.Data {
		if _, ok = role[r.Bid]; !ok {
			role[r.Bid] = make(map[int8]string)
		}
		role[r.Bid][r.Rid] = r.Name
	}
	return
}
