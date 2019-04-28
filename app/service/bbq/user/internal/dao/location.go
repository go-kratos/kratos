package dao

import (
	"context"
	"go-common/app/service/bbq/user/api"
	"go-common/library/database/sql"
	"go-common/library/log"
)

var (
	locationQuery = "select `loc_id`, `pid`, `name` from `bbq_location` where `loc_id` = ?;"
)

// GetLocation return the location info
func (d *Dao) GetLocation(c context.Context, locId int32) (*api.LocationItem, error) {
	row := d.db.QueryRow(c, locationQuery, locId)
	var location api.LocationItem
	err := row.Scan(&location.Id, &location.Pid, &location.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Errorv(c, log.KV("event", "mysql_select"), log.KV("table", "bbq_location"), log.KV("loc_id", locId))
		return nil, err
	}
	return &location, nil
}
