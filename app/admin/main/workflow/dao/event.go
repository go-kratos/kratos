package dao

import (
	"context"

	"go-common/app/admin/main/workflow/model"

	"github.com/pkg/errors"
)

// EventsByCid will select events by cid
func (d *Dao) EventsByCid(c context.Context, cid int64) (events map[int64]*model.Event, err error) {
	events = make(map[int64]*model.Event)

	elist := make([]*model.Event, 0)
	if err = d.ReadORM.Table("workflow_event").Where("cid=?", cid).Find(&elist).Error; err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, e := range elist {
		e.FixAttachments()
		events[e.Eid] = e
	}
	return
}

// EventsByIDs will select events by eids
func (d *Dao) EventsByIDs(c context.Context, eids []int64) (events map[int64]*model.Event, err error) {
	if len(eids) == 0 {
		return
	}
	events = make(map[int64]*model.Event, len(eids))
	elist := make([]*model.Event, 0)
	if err = d.ReadORM.Table("workflow_event").Where("id IN (?)", eids).Find(&elist).Error; err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, e := range elist {
		e.FixAttachments()
		events[e.Eid] = e
	}
	return
}

// LastEventByCid will retrive last event by cid
func (d *Dao) LastEventByCid(c context.Context, cid int64) (event *model.Event, err error) {
	event = new(model.Event)
	err = d.ReadORM.Table("workflow_event").Where("cid=?", cid).Order("id").Last(&event).Error
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// BatchLastEventIDs will retrive the last event ids by serveral conditions
func (d *Dao) BatchLastEventIDs(c context.Context, cids []int64) (eids []int64, err error) {
	eids = make([]int64, 0, len(cids))
	if len(cids) <= 0 {
		return
	}

	rows, err := d.ReadORM.Table("workflow_event").Select("max(id)").Where("cid IN (?)", cids).Group("cid").Rows()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var eid int64
		if err = rows.Scan(&eid); err != nil {
			err = errors.WithStack(err)
			return
		}
		eids = append(eids, eid)
	}

	return
}
