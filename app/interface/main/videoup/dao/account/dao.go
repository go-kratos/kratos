package account

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/videoup/conf"
	accapi "go-common/app/service/main/account/api"
	relation "go-common/app/service/main/relation/rpc/client"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_addFollowingURI = "/x/internal/relation/following/add"
)

// Dao dao is account dao.
type Dao struct {
	// config
	c *conf.Config
	// rpc
	rela *relation.Service
	acc  accapi.AccountClient
	// memcache
	mc                           *memcache.Pool
	mcSubExp, mcLimitAddBasicExp int32
	client                       *bm.Client
	addFollowingURL              string
}

// New new a account dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// rpc
		rela: relation.New(c.RelationRPC),
		// access memcache
		mc:                 memcache.NewPool(c.Memcache.Account.Config),
		mcSubExp:           int32(time.Duration(c.Memcache.Account.SubmitExpire) / time.Second),
		mcLimitAddBasicExp: int32(time.Duration(c.Limit.AddBasicExp) / time.Second),
		client:             bm.NewClient(c.HTTPClient.Write),
		addFollowingURL:    c.Host.APICo + _addFollowingURI,
	}
	var err error
	if d.acc, err = accapi.NewClient(c.AccClient); err != nil {
		panic(err)
	}
	return
}

// Close close resource.
func (d *Dao) Close() {
	if d.mc != nil {
		d.mc.Close()
	}
}

// Ping ping success.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingMemcache(c); err != nil {
		return
	}
	return
}

// IdentifyInfo 获取用户实名认证状态
func (d *Dao) IdentifyInfo(c context.Context, ip string, mid int64) (err error) {
	var mf *accapi.Profile
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

// AddFollowing 添加关注
func (d *Dao) AddFollowing(c context.Context, mid, fid int64, src int, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("fid", strconv.FormatInt(fid, 10))
	params.Set("src", strconv.Itoa(src))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.addFollowingURL, ip, params, &res); err != nil {
		log.Error("d.client.Do uri(%s) mid(%d) fid(%d) res(%+v) error(%v)", d.addFollowingURL+"?"+params.Encode(), mid, fid, res, err)
		return
	}
	log.Info("acc AddFollowing url(%s)", d.addFollowingURL+"?"+params.Encode())
	if res.Code != 0 {
		log.Error("acc AddFollowing (%+s)|(%+d)|(%+d)|(%+d)|(%s) (%+v)", d.addFollowingURL, mid, fid, src, ip, res)
		err = ecode.CreativeAccServiceErr
		return
	}
	return
}
