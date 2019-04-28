package realname

import (
	"context"

	"go-common/app/interface/main/account/conf"
	dao "go-common/app/interface/main/account/dao/realname"
	"go-common/library/log"
)

func (s *Service) alipayAntispamCheck(c context.Context, mid int64) bool {
	var (
		value *dao.AlipayAntispamValue
		err   error
	)
	if value, err = s.realnameDao.AlipayAntispam(c, mid); err != nil {
		log.Error("%+v", err)
		return true
	}
	if value == nil {
		return true
	}
	if value.Count() > conf.Conf.Realname.AlipayAntispamThreshold && !value.Pass() {
		return false
	}
	return true
}

// alipayAntispamIncrease 增加用户申请计数
func (s *Service) alipayAntispamIncrease(c context.Context, mid int64) (err error) {
	var (
		value *dao.AlipayAntispamValue
	)
	if value, err = s.realnameDao.AlipayAntispam(c, mid); err != nil {
		return
	}
	if value == nil {
		value = new(dao.AlipayAntispamValue)
	}
	value.IncreaseCount()
	if err = s.realnameDao.SetAlipayAntispam(c, mid, value); err != nil {
		return
	}
	return
}

// setAlipayAntispamPassFlag 更新用户通过标识位
func (s *Service) setAlipayAntispamPassFlag(c context.Context, mid int64, flag bool) (err error) {
	var (
		value *dao.AlipayAntispamValue
	)
	if value, err = s.realnameDao.AlipayAntispam(c, mid); err != nil {
		return
	}
	if value == nil {
		value = new(dao.AlipayAntispamValue)
	}
	value.SetPass(flag)
	if err = s.realnameDao.SetAlipayAntispam(c, mid, value); err != nil {
		return
	}
	return
}
