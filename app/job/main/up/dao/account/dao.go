package account

import (
	"context"

	"go-common/app/job/main/up/conf"
	"go-common/app/service/main/account/model"
	account "go-common/app/service/main/account/rpc/client"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Dao dao is account dao.
type Dao struct {
	// config
	c *conf.Config
	// rpc
	acc *account.Service3
}

// New new a account dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// config
		c: c,
		// rpc
		acc: account.New3(c.AccountRPC),
	}
	return
}

// Close close resource.
func (d *Dao) Close() {
}

// Ping ping success.
func (d *Dao) Ping(c context.Context) (err error) {
	return
}

// IdentifyInfo 获取用户实名认证状态
func (d *Dao) IdentifyInfo(c context.Context, ak, ck, ip string, mid int64) (err error) {
	var mf *model.Profile
	if mf, err = d.Profile(c, mid, ip); err != nil {
		log.Error("d.Profile mid(%d),ip(%s),error(%v)", mid, ip, err)
		err = ecode.CreativeAccServiceErr
		return
	}
	if mf.Identification == 1 {
		return
	}
	//switch for FrontEnd return json format, return OldPhone, and newError
	if err = d.switchIDInfoRet(mf.TelStatus); err != nil {
		log.Error("switchIDInfoRet url(%s) res(%v)", mf.TelStatus)
		return
	}
	return
}

func (d *Dao) switchIDInfoRet(phoneRet int32) (err error) {
	switch phoneRet {
	case 0:
		err = ecode.UserCheckNoPhone
	case 1:
		err = nil
	case 2:
		err = ecode.UserCheckInvalidPhone
	}
	return
}
