package account

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/appeal"
	accapi "go-common/app/service/main/account/api"
	relaMdl "go-common/app/service/main/relation/model"
	relation "go-common/app/service/main/relation/rpc/client"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/sync/errgroup"
	"net/http"
	"net/url"
	"time"
)

const (
	_pass         = "/web/site/userInfo"
	_richRelation = "/x/account/relation/rich"
	_relation     = "/x/internal/relation"
)

// Dao is account dao.
type Dao struct {
	c *conf.Config
	// rpc
	acc  accapi.AccountClient
	rela *relation.Service
	// http client
	client     *httpx.Client
	fastClient *httpx.Client
	mc         *memcache.Pool
	mcExpire   int32
	// user
	passURI         string
	relationURI     string
	richRelationURI string
	picUpInfoURL    string
	blinkUpInfoURL  string
	upInfoURL       string
}

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:    c,
		rela: relation.New(c.RelationRPC),
		// http client
		client:          httpx.NewClient(c.HTTPClient.Normal),
		fastClient:      httpx.NewClient(c.HTTPClient.Fast),
		passURI:         c.Host.Passport + _pass,
		relationURI:     c.Host.API + _relation,
		richRelationURI: c.Host.API + _richRelation,
		picUpInfoURL:    c.Host.Live + _picUpInfoURL,
		blinkUpInfoURL:  c.Host.Live + _blinkUpInfoURL,
		upInfoURL:       c.Host.API + _upInfoURL,
		mc:              memcache.NewPool(c.Memcache.Archive.Config),
		mcExpire:        int32(time.Duration(c.Memcache.Archive.TplExpire) / time.Second),
	}
	var err error
	if d.acc, err = accapi.NewClient(c.AccClient); err != nil {
		panic(err)
	}
	return
}

// Profile get account.
func (d *Dao) Profile(c context.Context, mid int64, ip string) (res *accapi.Profile, err error) {
	var (
		arg = &accapi.MidReq{
			Mid: mid,
		}
		rpcRes *accapi.ProfileReply
	)
	if rpcRes, err = d.acc.Profile3(c, arg); err != nil {
		log.Error("d.acc.Profile3 error(%v)", err)
		err = ecode.CreativeAccServiceErr
	}
	if rpcRes != nil {
		res = rpcRes.Profile
	}
	return
}

// ProfileWithStat get account.
func (d *Dao) ProfileWithStat(c context.Context, mid int64) (res *accapi.ProfileStatReply, err error) {
	var (
		arg = &accapi.MidReq{
			Mid: mid,
		}
	)
	if res, err = d.acc.ProfileWithStat3(c, arg); err != nil {
		log.Error("d.acc.ProfileWithStat3() error(%v)", err)
		err = ecode.CreativeAccServiceErr
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
		log.Error("d.acc.Card3() error(%v)", err)
		err = ecode.CreativeAccServiceErr
	}
	if rpcRes != nil {
		res = rpcRes.Card
	}
	return
}

// Cards get cards from rpc
func (d *Dao) Cards(c context.Context, mids []int64, ip string) (cards map[int64]*accapi.Card, err error) {
	if len(mids) == 0 {
		return
	}
	arg := &accapi.MidsReq{
		Mids: mids,
	}
	var reply *accapi.CardsReply
	if reply, err = d.acc.Cards3(c, arg); err != nil {
		log.Error("d.acc.Cards3 error(%v) | mids(%v) ip(%s) arg(%v)", err, mids, ip, arg)
		err = ecode.CreativeAccServiceErr
		return
	}
	cards = reply.Cards
	return
}

// PhoneEmail get user email & phone
func (d *Dao) PhoneEmail(c context.Context, ck, ip string) (ct *appeal.Contact, err error) {
	params := url.Values{}
	params.Set("Cookie", ck)
	// init req set cookie
	req, err := http.NewRequest("GET", d.passURI, nil)
	if err != nil {
		log.Error("passport url(%s) error(%v)", d.passURI, err)
		return
	}
	req.Header.Set("Cookie", ck)
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	var res struct {
		Code int             `json:"code"`
		Data *appeal.Contact `json:"data"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("passport url(%s) response(%+v) error(%v)", d.passURI+"?"+params.Encode(), res, err)
		err = ecode.CreativeAccServiceErr
		return
	}
	if res.Code != 0 {
		log.Error("passport url(%s) res(%v)", d.passURI, res)
		err = ecode.CreativeAccServiceErr
		return
	}
	ct = res.Data
	return
}

// RichRelation get multi user relations
func (d *Dao) RichRelation(c context.Context, owner int64, mids []int64, ip string) (richRel map[int64]int32, err error) {
	var rpcRes *accapi.RichRelationsReply
	if rpcRes, err = d.acc.RichRelations3(c, &accapi.RichRelationReq{Owner: owner, Mids: mids}); err != nil {
		log.Error("d.acc.RichRelations3(%d, %v) error(%v)", owner, mids, err)
		err = ecode.CreativeAccServiceErr
	}
	if rpcRes != nil {
		richRel = rpcRes.RichRelations
	}
	return
}

// Followers get users only follower relation which attr eq 2 with errgroup
func (d *Dao) Followers(c context.Context, owner int64, mids []int64, ip string) (relations map[int64]int, err error) {
	relations = make(map[int64]int, len(mids))
	type midRel struct {
		mid int64
		rel int
	}
	rechan := make(chan midRel, len(mids)) // avoid concurrent write map
	g, ctx := errgroup.WithContext(c)
	for _, mid := range mids {
		var mid2 = mid // for closure, copy to avoid same mid
		g.Go(func() error {
			rl, e := d.acc.Relation3(ctx, &accapi.RelationReq{Mid: mid2, Owner: owner})
			if e != nil || rl == nil || !rl.Following {
				if e != nil {
					log.Error("d.acc.Relation3(mid:%d,owner:%d,relation:%v) error(%v)", mid2, owner, rl, e)
				}
				rechan <- midRel{mid2, 1}
			} else {
				rechan <- midRel{mid2, 2}
			}
			return nil
		})
	}
	g.Wait()
	for i := 0; i < len(mids); i++ {
		wc := <-rechan
		relations[wc.mid] = wc.rel
	}
	return
}

// Infos get user info by mids.
func (d *Dao) Infos(c context.Context, mids []int64, ip string) (res map[int64]*accapi.Info, err error) {
	res = make(map[int64]*accapi.Info)
	if len(mids) == 0 {
		return
	}
	var arg = &accapi.MidsReq{
		Mids: mids,
	}
	var rpcRes *accapi.InfosReply
	if rpcRes, err = d.acc.Infos3(c, arg); err != nil {
		log.Error("d.acc.Infos3() error(%v)|ip(%s)", err, ip)
		err = ecode.CreativeAccServiceErr
	}
	if rpcRes != nil {
		res = rpcRes.Infos
	}
	return
}

// IdentifyInfo 获取用户实名认证状态
// tel_status	int	0未绑定，1已绑定有效手机号 2绑定虚拟号段170/171
// identification	int	身份证绑定状态，0:未绑定 1:已绑定
func (d *Dao) IdentifyInfo(c context.Context, mid int64, phoneOnly int8, ip string) (ret int, err error) {
	var (
		rpcRes *accapi.ProfileReply
		arg    = &accapi.MidReq{
			Mid: mid,
		}
		mf *accapi.Profile
	)
	if rpcRes, err = d.acc.Profile3(c, arg); err != nil {
		log.Error("d.acc.Profile3 error(%v) | mid(%d) ip(%s) arg(%v)", err, mid, ip, arg)
		err = ecode.CreativeAccServiceErr
		return
	}
	if rpcRes != nil {
		mf = rpcRes.Profile
	}
	//switch for FrontEnd return json format
	ret = d.switchPhoneRet(int(mf.TelStatus))
	if phoneOnly == 1 {
		return
	}
	if mf.TelStatus == 1 || mf.Identification == 1 {
		return 0, err
	}
	return
}

// MidByName 获取mid
func (d *Dao) MidByName(c context.Context, name string) (mid int64, err error) {
	var (
		rpcRes *accapi.InfosReply
		arg    = &accapi.NamesReq{
			Names: []string{name},
		}
		infos map[int64]*accapi.Info
	)
	if rpcRes, err = d.acc.InfosByName3(c, arg); err != nil {
		log.Error("d.acc.InfosByName3 error(%v)", err)
		err = ecode.CreativeAccServiceErr
		return
	}
	if rpcRes != nil {
		infos = rpcRes.Infos
	}
	for _, v := range infos {
		if v != nil && v.Name == name {
			return v.Mid, nil
		}
	}
	err = ecode.AccountInexistence
	return
}

// 0: "已实名认证",
// 1: "根据国家实名制认证的相关要求，您需要换绑一个非170/171的手机号，才能继续进行操作。",
// 2: "根据国家实名制认证的相关要求，您需要绑定手机号，才能继续进行操作。",
func (d *Dao) switchPhoneRet(newV int) (oldV int) {
	switch newV {
	case 0:
		oldV = 2
	case 1:
		oldV = 0
	case 2:
		oldV = 1
	}
	return
}

// CheckIdentify fn
func (d *Dao) CheckIdentify(identify int) (err error) {
	switch identify {
	case 0:
		err = nil
	case 1:
		err = ecode.UserCheckInvalidPhone
	case 2:
		err = ecode.UserCheckNoPhone
	}
	return
}

// RelationFollowers get all relation state.
func (d *Dao) RelationFollowers(c context.Context, mid int64, ip string) (res map[int64]int32, err error) {
	var fls []*relaMdl.Following
	if fls, err = d.rela.Followers(c, &relaMdl.ArgMid{Mid: mid, RealIP: ip}); err != nil {
		log.Error("d.rela.Followers mid(%d)|ip(%s)|error(%v)", mid, ip, err)
		return
	}
	if len(fls) == 0 {
		log.Info("d.rela.Followers mid(%d)|ip(%s)", mid, ip)
		return
	}
	res = make(map[int64]int32, len(fls))
	for _, v := range fls {
		res[v.Mid] = int32(v.Attribute)
	}
	log.Info("d.rela.Followers mid(%d)|res(%+v)|ip(%s)", mid, res, ip)
	return
}

// Relations get all relation state.
func (d *Dao) Relations(c context.Context, mid int64, fids []int64, ip string) (res map[int64]int, err error) {
	var rls map[int64]*relaMdl.Following
	if rls, err = d.rela.Relations(c, &relaMdl.ArgRelations{Mid: mid, Fids: fids, RealIP: ip}); err != nil {
		log.Error("d.rela.Relations mid(%d)|ip(%s)|error(%v)", mid, ip, err)
		err = ecode.CreativeAccServiceErr
		return
	}
	if len(rls) == 0 {
		log.Info("d.rela.Relations mid(%d)|ip(%s)", mid, ip)
		return
	}
	res = make(map[int64]int, len(rls))
	for _, v := range rls {
		res[v.Mid] = int(v.Attribute)
	}
	log.Info("d.rela.Relations mid(%d)|res(%+v)|rls(%+v)|ip(%s)", mid, res, rls, ip)
	return
}

// Relations2
func (d *Dao) Relations2(c context.Context, mid int64, fids []int64, ip string) (res map[int64]int, err error) {
	if res, err = d.Relations(c, mid, fids, ip); err != nil {
		return
	}
	for k, v := range res {
		if v == 2 || v == 6 { //2表示我关注他,6表示双向关注,为了兼容客户端统一吐出6作为已关注状态.
			res[k] = 6
		} else {
			delete(res, k)
		}
	}
	return
}

// ShouldFollow  fn
func (d *Dao) ShouldFollow(c context.Context, mid int64, fids []int64, ip string) (shouldMids []int64, err error) {
	var rls map[int64]*relaMdl.Following
	if rls, err = d.rela.Relations(c, &relaMdl.ArgRelations{Mid: mid, Fids: fids, RealIP: ip}); err != nil {
		log.Error("d.rela.Relations mid(%d)|fids(%+v)|ip(%s)|error(%v)", mid, fids, ip, err)
		return
	}
	if len(rls) == 0 {
		shouldMids = fids
		return
	}
	shouldMids = make([]int64, 0)
	for _, v := range rls {
		if v.Attribute == 0 {
			shouldMids = append(shouldMids, v.Mid)
		}
	}
	log.Info("d.rela.Relations mid(%d)|shouldMids(%+v)|ip(%s)", mid, shouldMids, ip)
	return
}
