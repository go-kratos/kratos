package service

import (
	"context"
	"time"

	"go-common/app/interface/main/growup/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// avStat get stat from hbase.
func (s *Service) avStat(c context.Context, mid int64, ip string) (up *model.UpBaseStat, err error) {
	hbaseDate := time.Now().AddDate(0, 0, -1).Add(-12 * time.Hour).Format("20060102")
	up, err = s.dao.UpStat(c, mid, hbaseDate)
	if err != nil || up == nil {
		log.Error("s.data.UpStat error(%v) mid(%d) up(%v) ip(%s)", err, mid, up, ip)
		err = ecode.CreativeDataErr
		return
	}
	pfl, err := s.dao.ProfileWithStat(c, mid)
	if err != nil {
		return
	}
	up.Fans = int64(pfl.Follower)
	log.Info("s.data.UpStat hbaseDate(%+v) mid(%d)", up, mid)
	return
}
