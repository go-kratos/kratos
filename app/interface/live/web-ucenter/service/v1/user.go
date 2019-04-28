package v1

import (
	"context"
	"github.com/pkg/errors"
	pb "go-common/app/interface/live/web-ucenter/api/http"
	"go-common/app/interface/live/web-ucenter/conf"
	"go-common/app/interface/live/web-ucenter/dao/user"
	"go-common/app/interface/live/web-ucenter/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
	"strings"
)

// UserService user service
type UserService struct {
	c   *conf.Config
	dao *user.Dao
}

// NewUserService new user service
func NewUserService(c *conf.Config) (s *UserService) {
	s = &UserService{
		c:   c,
		dao: user.New(c),
	}
	return s
}

// GetUserInfo implementation
// 根据uid查询用户信息
// `midware:"auth"`，需要登录态
func (s *UserService) GetUserInfo(ctx context.Context, req *pb.GetInfoReq) (resp *pb.GetInfoResp, err error) {
	var (
		group, errCtx = errgroup.WithContext(ctx)
		userExp       int64
		userRank      string
	)
	// check login
	uid := metadata.Int64(ctx, metadata.Mid)
	if uid == 0 {
		err = errors.Wrap(err, "请先登录")
		return
	}
	platform := checkPlatform(req.Platform)

	resp = &pb.GetInfoResp{
		Uid:         uid,
		UserCharged: 0,
	}

	// 并行获取account / xuser.vip / xuser.exp / wallet / rc / rankdb
	func() {
		// account.ProfileWithStat3
		group.Go(func() (err error) {
			profile, err := s.dao.GetAccountProfile(errCtx, uid)
			if err != nil {
				log.Error("[service.v1.user|GetUserInfo] GetAccountProfile error(%v), uid(%d)", err, uid)
				return nil
			}
			resp.Uname = profile.Name
			resp.Face = strings.Replace(profile.Face, "http://", "https://", 1)
			resp.Coin = profile.Coins
			return
		})
		// wallet
		group.Go(func() (err error) {
			silver, gold, err := s.dao.GetWallet(errCtx, uid, platform)
			if err != nil {
				log.Error("[service.v1.user|GetUserInfo] GetWallet error(%v), uid(%d)", err, uid)
				return nil
			}
			resp.Silver = silver
			resp.Gold = gold
			return
		})
		// xuser.vip
		group.Go(func() (err error) {
			vipInfo, err := s.dao.GetLiveVip(errCtx, uid)
			if err != nil || vipInfo == nil || vipInfo.Info == nil {
				log.Error("[service.v1.user|GetUserInfo] GetLiveVip error(%v), uid(%d)", err, uid)
				return nil
			}
			resp.Vip = vipInfo.Info.Vip
			resp.Svip = vipInfo.Info.Svip
			return
		})
		// xuser.exp
		group.Go(func() (err error) {
			expInfo, err := s.dao.GetLiveExp(errCtx, uid)
			if err != nil || expInfo == nil || expInfo.UserLevel == nil {
				log.Error("[service.v1.user|GetUserInfo] GetLiveExp error(%v), uid(%d)", err, uid)
				return nil
			}
			userExp = expInfo.UserLevel.UserExp
			resp.UserLevel = expInfo.UserLevel.Level
			resp.UserNextLevel = expInfo.UserLevel.NextLevel
			resp.UserIntimacy = expInfo.UserLevel.UserExp - expInfo.UserLevel.UserExpLeft
			resp.UserNextIntimacy = expInfo.UserLevel.UserExpNextLevel
			resp.IsLevelTop = expInfo.UserLevel.IsLevelTop
			return
		})
		// rc
		group.Go(func() (err error) {
			achieve, err := s.dao.GetLiveAchieve(errCtx, uid)
			if err != nil {
				log.Error("[service.v1.user|GetUserInfo] GetLiveAchieve error(%v), uid(%d)", err, uid)
				return nil
			}
			resp.Achieve = achieve
			return
		})
		// rankdb
		group.Go(func() (err error) {
			if userRank, err = s.dao.GetLiveRank(errCtx, uid); err != nil {
				log.Error("[service.v1.user|GetUserInfo] GetLiveRank error(%v), uid(%d)", err, uid)
				return nil
			}
			return
		})
	}()
	group.Wait()

	// 根据exp & rankdb 判断组装返回的user_level_rank字段
	if userExp < 120000000 {
		resp.UserLevelRank = ">50000"
	} else {
		resp.UserLevelRank = userRank
	}

	log.Info("GetUserInfo.resp(%v)", resp)
	return
}

func checkPlatform(p string) string {
	if p == "" || (p != model.PlatformIos && p != model.PlatformAndroid) {
		return model.PlatformPc
	}
	return p
}
