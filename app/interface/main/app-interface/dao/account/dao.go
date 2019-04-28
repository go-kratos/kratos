package account

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-interface/conf"
	account "go-common/app/service/main/account/model"
	accrpc "go-common/app/service/main/account/rpc/client"
	memberrpc "go-common/app/service/main/member/api/gorpc"
	blockmodel "go-common/app/service/main/member/model/block"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Dao is account dao.
type Dao struct {
	client *bm.Client
	// rpc
	accRPC    *accrpc.Service3
	memberRPC *memberrpc.Service
}

// New account dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:    bm.NewClient(c.HTTPClient),
		accRPC:    accrpc.New3(c.AccountRPC),
		memberRPC: memberrpc.New(c.MemberRPC),
	}
	return
}

// BlockTime get user blocktime
func (d *Dao) BlockTime(c context.Context, mid int64) (blockTime int64, err error) {
	info, err := d.memberRPC.BlockInfo(c, &blockmodel.RPCArgInfo{MID: mid})
	if err != nil {
		err = errors.Wrapf(err, "%v", mid)
		return
	}
	if info.EndTime > 0 {
		blockTime = info.EndTime
	}
	return
}

// Profile3 get profile
func (d *Dao) Profile3(c context.Context, mid int64) (card *account.ProfileStat, err error) {
	arg := &account.ArgMid{Mid: mid}
	if card, err = d.accRPC.ProfileWithStat3(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

// Card get card
func (d *Dao) Card(c context.Context, mid int64) (card *account.Card, err error) {
	if card, err = d.accRPC.Card3(c, &account.ArgMid{Mid: mid}); err != nil {
		err = errors.Wrapf(err, "%v", mid)
	}
	return
}

// ProfileByName3 rpc card get by name
func (d *Dao) ProfileByName3(c context.Context, name string) (card *account.ProfileStat, err error) {
	infos, err := d.accRPC.InfosByName3(c, &account.ArgNames{Names: []string{name}})
	if err != nil {
		err = errors.Wrapf(err, "%v", name)
		return
	}
	if len(infos) == 0 {
		err = ecode.NothingFound
		return
	}
	for mid := range infos {
		card, err = d.Profile3(c, mid)
		break
	}
	return
}

// Infos3 rpc info get by mids .
func (d *Dao) Infos3(c context.Context, mids []int64) (res map[int64]*account.Info, err error) {
	arg := &account.ArgMids{Mids: mids}
	if res, err = d.accRPC.Infos3(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

// Relations3 relations.
func (d *Dao) Relations3(c context.Context, owners []int64, mid int64) (follows map[int64]bool) {
	if len(owners) == 0 {
		return nil
	}
	follows = make(map[int64]bool, len(owners))
	for _, owner := range owners {
		follows[owner] = false
	}
	var (
		am  map[int64]*account.Relation
		err error
		ip  = metadata.String(c, metadata.RemoteIP)
	)
	arg := &account.ArgRelations{Owners: owners, Mid: mid, RealIP: ip}
	if am, err = d.accRPC.Relations3(c, arg); err != nil {
		log.Error("d.accRPC.Relations2(%v) error(%v)", arg, err)
		return
	}
	for i, a := range am {
		if _, ok := follows[i]; ok {
			follows[i] = a.Following
		}
	}
	return
}

// RichRelations3 rich relations.
func (d *Dao) RichRelations3(c context.Context, owner, mid int64) (rel int, err error) {
	var (
		res map[int64]int
		ip  = metadata.String(c, metadata.RemoteIP)
	)
	arg := &account.ArgRichRelation{Mids: []int64{mid}, Owner: owner, RealIP: ip}
	if res, err = d.accRPC.RichRelations3(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
		return
	}
	if r, ok := res[mid]; ok {
		rel = r
	}
	return
}

// Cards3 is
func (d *Dao) Cards3(c context.Context, mids []int64) (res map[int64]*account.Card, err error) {
	arg := &account.ArgMids{Mids: mids}
	if res, err = d.accRPC.Cards3(c, arg); err != nil {
		err = errors.Wrapf(err, "%v")
	}
	return
}

// UserCheck 各种入口白名单
// https://www.tapd.cn/20055921/prong/stories/view/1120055921001066980  动态互推TAPD在此！！
func (d *Dao) UserCheck(c context.Context, mid int64, checkURL string) (ok bool, err error) {
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Status int `json:"status"`
		} `json:"data"`
	}
	if err = d.client.Get(c, checkURL, "", params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), checkURL+"?"+params.Encode())
		return
	}
	if res.Data.Status == 1 {
		ok = true
	}
	return
}

// RedDot 我的页小红点逻辑
func (d *Dao) RedDot(c context.Context, mid int64, redDotURL string) (ok bool, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			RedDot bool `json:"red_dot"`
		} `json:"data"`
	}
	if err = d.client.Get(c, redDotURL, "", params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), redDotURL+"?"+params.Encode())
		return
	}
	log.Warn("reddot response mid(%d) url(%s) res(%t)", mid, redDotURL+"?"+params.Encode(), res.Data.RedDot)
	ok = res.Data.RedDot
	return
}
