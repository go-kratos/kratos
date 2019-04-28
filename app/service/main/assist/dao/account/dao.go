package account

import (
	"context"
	accapi "go-common/app/service/main/account/api"
	"go-common/app/service/main/assist/conf"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Dao is account dao.
type Dao struct {
	c   *conf.Config
	acc accapi.AccountClient
}

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
	}
	var err error
	if d.acc, err = accapi.NewClient(c.AccClient); err != nil {
		panic(err)
	}
	return
}

// IsFollow check assist follow up.
func (d *Dao) IsFollow(c context.Context, mid, assistMid int64) (follow bool, err error) {
	var arg = &accapi.RelationReq{
		Mid:   assistMid,
		Owner: mid,
	}
	res, err := d.acc.Relation3(c, arg)
	if err != nil {
		log.Error("d.acc.Relation2(%d,%d) error(%v)", mid, assistMid, err)
		return
	}
	follow = res.Following
	return
}

// IdentifyInfo 获取用户实名认证状态
func (d *Dao) IdentifyInfo(c context.Context, mid int64, ip string) (err error) {
	var (
		arg = &accapi.MidReq{
			Mid: mid,
		}
		rpcRes *accapi.ProfileReply
		mf     *accapi.Profile
	)
	if rpcRes, err = d.acc.Profile3(c, arg); err != nil {
		log.Error("d.acc.Profile3 error(%v) | mid(%d) ip(%s) arg(%v)", err, mid, ip, arg)
		err = ecode.CreativeAccServiceErr
		return
	}
	if rpcRes != nil {
		mf = rpcRes.Profile
	}
	if mf.Identification == 1 {
		return
	}
	if err = d.switchIDInfoRet(mf.TelStatus); err != nil {
		log.Error("switchIDInfoRet res(%v)", mf.TelStatus)
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

// UserBanned 获取用户封禁状态, disabled when spacesta == 2
func (d *Dao) UserBanned(c context.Context, mid int64) (err error) {
	var card *accapi.Card
	if card, err = d.Card(c, mid, ""); err != nil {
		log.Error("d.Card() error(%v)", err)
		err = nil
		return
	}
	if card.Silence == 1 {
		err = ecode.UserDisabled
		return
	}
	return
}

// Card get account.
func (d *Dao) Card(c context.Context, mid int64, ip string) (res *accapi.Card, err error) {
	var (
		rpcRes *accapi.CardReply
		arg    = &accapi.MidReq{
			Mid: mid,
		}
	)
	if rpcRes, err = d.acc.Card3(c, arg); err != nil {
		log.Error("s.acc.Card3() error(%v)", err)
		err = ecode.CreativeAccServiceErr
	}
	if rpcRes != nil {
		res = rpcRes.Card
	}
	return
}

// Cards get infos for space
func (d *Dao) Cards(c context.Context, mids []int64) (res map[int64]*accapi.Card, err error) {
	var (
		arg = &accapi.MidsReq{
			Mids: mids,
		}
		rpcRes *accapi.CardsReply
	)
	if rpcRes, err = d.acc.Cards3(c, arg); err != nil {
		log.Error("s.acc.Cards3() error(%v)", err)
		err = ecode.CreativeAccServiceErr
	}
	if rpcRes != nil {
		res = rpcRes.Cards
	}
	return
}
