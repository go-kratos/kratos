package user

import (
	"context"
	"github.com/pkg/errors"
	"go-common/app/interface/live/web-ucenter/conf"
	"go-common/app/interface/live/web-ucenter/dao"
	rankdbv1 "go-common/app/service/live/rankdb/api/liverpc/v1"
	rcv1 "go-common/app/service/live/rc/api/liverpc/v1"
	xuserv1 "go-common/app/service/live/xuser/api/grpc/v1"
	accModel "go-common/app/service/main/account/model"
	account "go-common/app/service/main/account/rpc/client"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"net/http"
	"strconv"
)

var (
	_walletApiUrl = "/x/internal/livewallet/wallet/getAll"
)

// Dao user dao, wrap clients
type Dao struct {
	c             *conf.Config
	bmClient      *bm.Client
	vipClient     xuserv1.VipClient
	expClient     xuserv1.UserExpClient
	walletUrl     string
	accountClient *account.Service3
	rankdbClient  rankdbv1.UserRank
	rcClient      rcv1.AchvRPCClient
}

// New new user dao
func New(c *conf.Config) *Dao {
	conn, err := xuserv1.NewClient(c.Warden)
	if err != nil {
		panic(err)
	}
	d := &Dao{
		c:             c,
		bmClient:      bm.NewClient(c.HTTPClient),
		walletUrl:     c.Host.LiveRpc + _walletApiUrl,
		accountClient: account.New3(c.AccountRPC),
		rankdbClient:  dao.RankdbApi.V1UserRank,
		rcClient:      dao.RcApi.V1Achv,
	}
	d.vipClient = conn.VipClient
	d.expClient = conn.UserExpClient
	return d
}

// GetAccountProfile get account profile
func (d *Dao) GetAccountProfile(ctx context.Context, uid int64) (profile *accModel.ProfileStat, err error) {
	arg := &accModel.ArgMid{Mid: uid}
	if profile, err = d.accountClient.ProfileWithStat3(ctx, arg); err != nil || profile == nil {
		log.Error("[dao.user|GetAccountProfile] get account profile3 error(%v), uid(%d), profile(%v)", err, uid, profile)
		return
	}
	return
}

// GetWallet get silver/gold from go-wallet by http request
func (d *Dao) GetWallet(ctx context.Context, uid int64, platform string) (silver, gold int64, err error) {
	m := make(map[string]string)
	m["uid"] = strconv.FormatInt(uid, 10)
	paramString := dao.EncodeHttpParams(m, d.c.HTTPClient.Key, d.c.HTTPClient.Secret)
	req, _ := http.NewRequest("GET", d.walletUrl+"?"+paramString, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("platform", platform)
	var wr struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Gold   string `json:"gold"`
			Silver string `json:"silver"`
		} `json:"data"`
	}
	if err = d.bmClient.Do(ctx, req, &wr); err != nil {
		log.Error("[dao.user|GetWallet] connect error(%v), uid(%d), platform(%s)", err, uid, platform)
		return
	}
	if wr.Code != 0 {
		err = errors.Wrap(ecode.Int(wr.Code), d.walletUrl+"?"+paramString)
		log.Error("[dao.user|GetWallet] request error(%v), uid(%d), platform(%s)", err, uid, platform)
		return
	}

	gold, _ = strconv.ParseInt(wr.Data.Gold, 10, 64)
	silver, _ = strconv.ParseInt(wr.Data.Silver, 10, 64)
	return
}

// GetLiveVip get live vip/svip from xuser.vip.Info
func (d *Dao) GetLiveVip(ctx context.Context, uid int64) (vipInfo *xuserv1.InfoReply, err error) {
	uidReq := &xuserv1.UidReq{
		Uid: uid,
	}
	if vipInfo, err = d.vipClient.Info(ctx, uidReq); err != nil || vipInfo == nil {
		log.Error("[dao.user|GetLiveVip] get vip error(%v), uid(%d)", err, uid)
		return
	}
	return
}

// GetLiveExp get live exp from xuser.exp.GetUserExp
func (d *Dao) GetLiveExp(ctx context.Context, uid int64) (expInfo *xuserv1.LevelInfo, err error) {
	req := &xuserv1.GetUserExpReq{
		Uids: []int64{uid},
	}
	resp, err := d.expClient.GetUserExp(ctx, req)
	if err != nil {
		log.Error("[dao.user|GetLiveExp] get exp error(%v), uid(%d)", err, uid)
		return
	}
	var ok bool
	if expInfo, ok = resp.Data[uid]; !ok {
		log.Error("[dao.user|GetLiveExp] get exp empty, uid(%d)", uid)
		return
	}
	return
}

// GetLiveAchieve get rc achieve by liverpc
func (d *Dao) GetLiveAchieve(ctx context.Context, uid int64) (achieve int64, err error) {
	resp, err := d.rcClient.Userstatus(ctx, &rcv1.AchvUserstatusReq{})
	if err != nil || resp == nil || resp.Data == nil {
		log.Error("[dao.user|GetLiveAchieve] get rc achieve error(%v), uid(%d), resp(%v)", err, uid, resp)
		return
	}
	achieve = resp.Data.Point
	return
}

// GetLiveRank get user rank by liverpc
func (d *Dao) GetLiveRank(ctx context.Context, uid int64) (rank string, err error) {
	rank = "1000000"
	req := &rankdbv1.UserRankGetUserRankReq{
		Uid:  uid,
		Type: "user_level",
	}
	resp, err := d.rankdbClient.GetUserRank(ctx, req)
	if err != nil || resp == nil || resp.Data == nil {
		log.Error("[dao.user|GetLiveRank] get rankdb user rank error(%v), uid(%d)", err, uid)
		return
	}
	rank = strconv.FormatInt(resp.Data.Rank, 10)
	return
}
