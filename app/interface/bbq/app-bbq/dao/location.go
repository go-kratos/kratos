package dao

import (
	"context"
	"go-common/app/interface/bbq/app-bbq/model"
)

var (
	locationQueryChild = "select `loc_id`, `pid`, `name` from `bbq_location` where `pid` = ?;"
	locationQueryAll   = "select `loc_id`, `pid`, `name` from `bbq_location`;"
)

// GetLocationAll .
func (d *Dao) GetLocationAll(c context.Context) (*map[int32][]*model.Location, error) {
	rows, err := d.db.Query(c, locationQueryAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[int32][]*model.Location)
	var id, pid int32
	var name string
	for rows.Next() {
		rows.Scan(&id, &pid, &name)
		m[pid] = append(m[pid], &model.Location{
			ID:   id,
			PID:  pid,
			Name: name,
		})
	}

	return &m, err
}

// GetLocationChild .
func (d *Dao) GetLocationChild(c context.Context, locID int32) (*map[int32][]*model.Location, error) {
	rows, err := d.db.Query(c, locationQueryChild, locID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[int32][]*model.Location)
	var id, pid int32
	var name string
	for rows.Next() {
		rows.Scan(&id, &pid, &name)
		m[pid] = append(m[pid], &model.Location{
			ID:   id,
			PID:  pid,
			Name: name,
		})
	}

	return &m, err
}
