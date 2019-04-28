package dao

import (
	"context"

	"go-common/app/job/main/tag/model"
	"go-common/library/log"
)

const (
	_bussinessSQL = "SELECT type, name, appkey, remark, alias FROM business WHERE state=?"
)

// Business Gets gets all business records
func (d *Dao) Business(c context.Context, state int32) (business map[string]*model.Business, err error) {
	business = make(map[string]*model.Business)
	rows, err := d.platform.Query(c, _bussinessSQL, state)
	if err != nil {
		log.Error("d.Business(%d) error(%v)", state, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.Business{}
		if err = rows.Scan(&b.Type, &b.Name, &b.Appkey, &b.Remark, &b.Alias); err != nil {
			log.Error("d.Business(%d) rows.Scan() error(%v)", state, err)
			return
		}
		business[b.Alias] = b
	}
	err = rows.Err()
	return
}
