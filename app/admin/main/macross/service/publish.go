package service

import (
	"context"

	"go-common/app/admin/main/macross/model/publish"
	"go-common/library/log"
)

// Dashborad insert dashboard info and logs.
func (s *Service) Dashborad(c context.Context, d *publish.Dashboard) (err error) {
	var id int64
	if id, err = s.dao.Dashborad(c, d); err != nil {
		log.Error("Dashborad() error(%v)", err)
		return
	}
	if len(d.Logs) > 0 && id > 0 {
		if _, err = s.dao.DashboradLogs(c, id, d.Logs); err != nil {
			log.Error("DashboradLogs() error(%v)", err)
			return
		}
	}
	return
}
