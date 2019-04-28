package dao

import (
	"context"
	"time"

	"go-common/library/log"
)

const (
	_updateJuryExpiredSQL = "UPDATE blocked_jury SET status=1, expired=? WHERE mid = ?"
	_selConfSQL           = "SELECT config_key,content FROM blocked_config"
)

// UpdateJuryExpired update jury expired.
func (d *Dao) UpdateJuryExpired(c context.Context, mid int64, expired time.Time) (err error) {
	if _, err = d.db.Exec(c, _updateJuryExpiredSQL, expired, mid); err != nil {
		log.Error("d.UpdateJuryExpired err(%v)", err)
	}
	return
}

// LoadConf load conf.
func (d *Dao) LoadConf(c context.Context) (cf map[string]string, err error) {
	cf = make(map[string]string)
	rows, err := d.db.Query(c, _selConfSQL)
	if err != nil {
		log.Error("d.loadConf err(%v)", err)
		return
	}
	defer rows.Close()
	var (
		key   string
		value string
	)
	for rows.Next() {
		if err = rows.Scan(&key, &value); err != nil {
			log.Error("rows.Scan err(%v)", err)
			continue
		}
		cf[key] = value
	}
	err = rows.Err()
	return
}
