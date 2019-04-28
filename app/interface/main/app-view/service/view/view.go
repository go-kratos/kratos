package view

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"strconv"
	"strings"
	"time"

	cdm "go-common/app/interface/main/app-card/model"
	"go-common/app/interface/main/app-card/model/card"
	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-view/model"
	"go-common/app/interface/main/app-view/model/creative"
	"go-common/app/interface/main/app-view/model/view"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	location "go-common/app/service/main/location/model"
	relation "go-common/app/service/main/relation/model"
	resource "go-common/app/service/main/resource/model"
	sharerpc "go-common/app/service/main/share/api"
	thumbuppb "go-common/app/service/main/thumbup/api"
	thumbup "go-common/app/service/main/thumbup/model"
	"go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"

	egV2 "go-common/library/sync/errgroup.v2"

	"github.com/pkg/errors"
)

const (
	_descLen      = 250
	_promptCoin   = 1
	_promptFav    = 2
	_avTypeAv     = 1
	_businessLike = "archive"
	_coinAv       = 1
)

// View  all view data.
func (s *Service) View(c context.Context, mid, aid, movieID int64, plat int8, build, qn, fnver, fnval, forceHost, parentMode int, ak, mobiApp, device, buvid, cdnIP, network, adExtra, from string, now time.Time) (v *view.View, err error) {
	if v, err = s.ViewPage(c, mid, aid, movieID, plat, build, ak, mobiApp, device, cdnIP, true, now); err != nil {
		ip := metadata.String(c, metadata.RemoteIP)
		if err == ecode.AccessDenied || err == ecode.NothingFound {
			log.Warn("s.ViewPage() mid(%d) aid(%d) movieID(%d) plat(%d) ak(%s) ip(%s) cdn_ip(%s) error(%v)", mid, aid, movieID, plat, ak, ip, cdnIP, err)
		} else {
			log.Error("s.ViewPage() mid(%d) aid(%d) movieID(%d) plat(%d) ak(%s) ip(%s) cdn_ip(%s) error(%v)", mid, aid, movieID, plat, ak, ip, cdnIP, err)
		}
		return
	}
	g, ctx := errgroup.WithContext(c)
	g.Go(func() (err error) {
		v.VIPActive = s.vipActiveCache[view.VIPActiveView]
		return
	})
	g.Go(func() (err error) {
		s.initReqUser(ctx, v, mid, plat, build)
		if v.AttrVal(archive.AttrBitIsPGC) != archive.AttrYes {
			s.initContributions(ctx, v)
		}
		return
	})
	g.Go(func() (err error) {
		s.initRelateCMTag(ctx, v, plat, build, qn, fnver, fnval, forceHost, parentMode, mid, buvid, mobiApp, device, network, adExtra, from, now)
		return
	})
	if v.AttrVal(archive.AttrBitIsPGC) != archive.AttrYes {
		g.Go(func() (err error) {
			s.initDM(ctx, v)
			return
		})
		g.Go(func() (err error) {
			s.initAudios(ctx, v)
			return
		})
		g.Go(func() (err error) {
			if model.IsIPhoneB(plat) || (model.IsIPhone(plat) && (build >= 7000 && build <= 8000)) {
				return
			}
			s.initElec(ctx, v, mid)
			return
		})
		if len([]rune(v.Desc)) > _descLen {
			g.Go(func() (err error) {
				if desc, _ := s.arcDao.Description(ctx, v.Aid); desc != "" {
					v.Desc = desc
				}
				return
			})
		}
	}
	if v.AttrVal(archive.AttrBitHasBGM) == archive.AttrYes {
		g.Go(func() (err error) {
			if v.Bgm, err = s.creativeDao.Bgm(ctx, v.Aid, v.FirstCid); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
	}
	return
}

// ViewPage view page data.
func (s *Service) ViewPage(c context.Context, mid, aid, movieID int64, plat int8, build int, ak, mobiApp, device, cdnIP string, nMovie bool, now time.Time) (v *view.View, err error) {
	if aid == 0 && movieID == 0 {
		err = ecode.NothingFound
		return
	}
	const (
		_androidMovie = 5220000
		_iPhoneMovie  = 6500
		_iPadMovie    = 6720
		_iPadHDMovie  = 12020
	)
	var (
		vs         *view.ViewStatic
		vp         *archive.View3
		seasoninfo map[int64]int64
		ok         bool
	)
	if movieID != 0 {
		if seasoninfo, err = s.banDao.SeasonidAid(c, movieID, now); err != nil {
			log.Error("%+v", err)
			err = ecode.NothingFound
			return
		}
		if aid, ok = seasoninfo[movieID]; !ok || aid == 0 {
			err = ecode.NothingFound
			return
		}
		var a *api.Arc
		if a, err = s.arcDao.Archive3(c, aid); err != nil {
			log.Error("%+v", err)
			err = ecode.NothingFound
			return
		}
		if a == nil {
			err = ecode.NothingFound
			return
		}
		vs = &view.ViewStatic{Archive3: archive.BuildArchive3(a)}
		s.prom.Incr("from_movieID")
	} else {
		if vp, err = s.arcDao.ViewCache(c, aid); err != nil {
			log.Error("%+v", err)
		}
		if vp == nil || vp.Archive3 == nil || len(vp.Pages) == 0 || vp.AttrVal(archive.AttrBitIsMovie) == archive.AttrYes {
			if vp, err = s.arcDao.View3(c, aid); err != nil {
				log.Error("%+v", err)
				err = ecode.NothingFound
				return
			}
		}
		if vp == nil || vp.Archive3 == nil || len(vp.Pages) == 0 {
			err = ecode.NothingFound
			return
		}
		vs = &view.ViewStatic{Archive3: vp.Archive3}
		s.initPages(c, vs, vp.Pages)
		s.prom.Incr("from_aid")
	}
	if _, ok := s.specialMids[vs.Author.Mid]; ok && env.DeployEnv == env.DeployEnvProd {
		err = ecode.NothingFound
		log.Error("aid(%d) mid(%d) can not view on prod", vs.Aid, vs.Author.Mid)
		return
	}
	// TODO 产品最帅了！
	vs.Stat.DisLike = 0
	if s.overseaCheck(vs.Archive3, plat) {
		err = ecode.AreaLimit
		return
	}
	// check region area limit
	if err = s.areaLimit(c, plat, int(vs.TypeID)); err != nil {
		return
	}
	v = &view.View{ViewStatic: vs, DMSeg: 1, PlayerIcon: s.playerIcon}
	if v.AttrVal(archive.AttrBitIsPGC) != archive.AttrYes {
		// check access
		if err = s.checkAceess(c, mid, v.Aid, int(v.State), int(v.Access), ak); err != nil {
			// archive is ForbitFixed and Transcoding and StateForbitDistributing need analysis history body .
			if v.State != archive.StateForbidFixed {
				return
			}
			err = nil
		}
		if v.Access > 0 {
			v.Stat.View = 0
		}
	}
	g, ctx := errgroup.WithContext(c)
	if mid != 0 {
		g.Go(func() (err error) {
			v.History, _ = s.arcDao.Progress(ctx, v.Aid, mid)
			return
		})
	}
	if v.AttrVal(archive.AttrBitIsPGC) == archive.AttrYes {
		if (v.AttrVal(archive.AttrBitIsMovie) != archive.AttrYes) || (plat == model.PlatAndroid && build >= _androidMovie) || (plat == model.PlatIPhone && build >= _iPhoneMovie) || (plat == model.PlatIPad && build >= _iPadMovie) ||
			(plat == model.PlatIpadHD && build > _iPadHDMovie) || plat == model.PlatAndroidTVYST || plat == model.PlatAndroidTV || plat == model.PlatAndroidI || plat == model.PlatIPhoneB {
			g.Go(func() error {
				return s.initPGC(ctx, v, mid, build, mobiApp, device)
			})
		} else {
			g.Go(func() error {
				return s.initMovie(ctx, v, mid, build, mobiApp, device, nMovie)
			})
		}
	} else {
		g.Go(func() (err error) {
			if err = s.initDownload(ctx, v, mid, cdnIP); err != nil {
				ip := metadata.String(ctx, metadata.RemoteIP)
				log.Error("aid(%d) mid(%d) ip(%s) cdn_ip(%s) error(%+v)", v.Aid, mid, ip, cdnIP, err)
			}
			return
		})
		if v.Rights.UGCPay == 1 && mid != v.Author.Mid {
			g.Go(func() (err error) {
				if err = s.initUGCPay(ctx, v, plat, mid); err != nil {
					log.Error("%+v", err)
					err = nil
					return
				}
				return nil
			})
		}
	}
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
	}
	if v.Rights.UGCPay == 1 {
		if (v.Asset == nil || v.Asset.Paid == 0) && mid != v.Author.Mid {
			v.Rights.Download = int32(location.ForbiddenDown)
		}
	}
	return
}

// AddShare add a share
func (s *Service) AddShare(c context.Context, aid, mid int64, ip string) (share int, isReport bool, upID int64, err error) {
	var a *api.Arc
	if a, err = s.arcDao.Archive(c, aid); err != nil {
		if errors.Cause(err) == ecode.NothingFound {
			err = ecode.ArchiveNotExist
		}
		return
	}
	if !a.IsNormal() {
		err = ecode.ArchiveNotExist
		return
	}
	upID = a.Author.Mid
	shareReply, err := s.shareClient.AddShare(context.Background(), &sharerpc.AddShareRequest{
		Oid:  aid,
		Mid:  mid,
		Type: 3,
		Ip:   ip,
	})
	if err != nil {
		if ecode.Cause(err).Equal(ecode.ShareAlreadyAdd) {
			err = nil
			return
		}
		log.Error("s.shareClient.AddShare(%d, %d, 3) error(%v)", aid, mid, err)
		return
	}
	if shareReply != nil && shareReply.Shares > int64(a.Stat.Share) {
		isReport = true
	}
	return
}

// Shot shot service
func (s *Service) Shot(c context.Context, aid, cid int64) (shot *view.Videoshot, err error) {
	var (
		arcShot *archive.Videoshot
		points  []*creative.Points
	)
	shot = new(view.Videoshot)
	if arcShot, err = s.arcDao.Shot(c, aid, cid); err != nil {
		log.Error("%+v", err)
		return
	}
	if arcShot == nil {
		return
	}
	shot.Videoshot = arcShot
	a := &api.Arc{Attribute: arcShot.Attr}
	if a.AttrVal(archive.AttrBitHasViewpoint) == archive.AttrYes {
		if points, err = s.creativeDao.Points(c, aid, cid); err != nil {
			log.Error("%+v", err)
			err = nil
		}
		shot.Points = points
	}
	return
}

// Like add a like.
func (s *Service) Like(c context.Context, aid, mid int64, status int8) (upperID int64, toast string, err error) {
	var (
		a    *api.Arc
		typ  int8
		stat *thumbuppb.LikeReply
	)
	if a, err = s.arcDao.Archive(c, aid); err != nil {
		if errors.Cause(err) == ecode.NothingFound {
			err = ecode.ArchiveNotExist
		}
		return
	}
	if !a.IsNormal() {
		err = ecode.ArchiveNotExist
		return
	}
	upperID = a.Author.Mid
	if status == 0 {
		typ = thumbup.TypeLike
	} else if status == 1 {
		typ = thumbup.TypeCancelLike
	}
	if stat, err = s.thumbupDao.Like(c, mid, upperID, _businessLike, a.Aid, typ); err != nil {
		if ecode.EqualError(ecode.ThumbupDupLikeErr, err) {
			log.Error("%+v", err)
			err = nil
			toast = "点赞收到！视频可能推荐哦"
		}
		return
	}
	if typ == thumbup.TypeLike {
		if stat.LikeNumber < 100 {
			toast = "点赞收到！视频可能推荐哦"
		} else if stat.LikeNumber >= 100 && stat.LikeNumber < 1000 {
			toast = "感谢点赞，推荐已收到啦"
		} else if stat.LikeNumber >= 1000 && stat.LikeNumber < 10000 {
			toast = "get！视频也许更多人能看见！"
		} else {
			toast = "点赞爆棚，感谢推荐！"
		}
	}
	return
}

// Dislike add a dislike.
func (s *Service) Dislike(c context.Context, aid, mid int64, status int8) (upperID int64, err error) {
	var (
		a   *api.Arc
		typ int8
	)
	if a, err = s.arcDao.Archive(c, aid); err != nil {
		if errors.Cause(err) == ecode.NothingFound {
			err = ecode.ArchiveNotExist
		}
		return
	}
	if !a.IsNormal() {
		err = ecode.ArchiveNotExist
		return
	}
	upperID = a.Author.Mid
	if status == 0 {
		typ = thumbup.TypeDislike
	} else if status == 1 {
		typ = thumbup.TypeCancelDislike
	}
	_, err = s.thumbupDao.Like(c, mid, upperID, _businessLike, a.Aid, typ)
	return
}

// AddCoin add a coin
func (s *Service) AddCoin(c context.Context, aid, mid, upID, avtype, multiply int64, ak string, selectLike int) (prompt, like bool, err error) {
	var maxCoin int64 = 2
	var typeID int16
	var pubTime int64
	if avtype == _avTypeAv {
		var a *api.Arc
		if a, err = s.arcDao.Archive(c, aid); err != nil {
			if errors.Cause(err) == ecode.NothingFound {
				err = ecode.ArchiveNotExist
			}
			return
		}
		if !a.IsNormal() {
			err = ecode.ArchiveNotExist
			return
		}
		if a.Copyright == int32(archive.CopyrightCopy) {
			maxCoin = 1
		}
		upID = a.Author.Mid
		typeID = int16(a.TypeID)
		pubTime = int64(a.PubDate)
	}
	if err = s.coinDao.AddCoins(c, aid, mid, upID, maxCoin, avtype, multiply, typeID, pubTime); err != nil {
		return
	}
	eg, ctx := errgroup.WithContext(c)
	eg.Go(func() (err error) {
		if avtype == _avTypeAv && selectLike == 1 {
			if _, err = s.thumbupDao.Like(ctx, mid, upID, _businessLike, aid, thumbup.TypeLike); err != nil {
				log.Error("%+v", err)
				err = nil
			} else {
				like = true
			}
		}
		return
	})
	eg.Go(func() (err error) {
		if prompt, err = s.relDao.Prompt(ctx, mid, upID, _promptCoin); err != nil {
			log.Error("%+v", err)
			err = nil
		}
		return
	})
	eg.Wait()
	return
}

// AddFav add a favorite
func (s *Service) AddFav(c context.Context, mid, vmid int64, fids []int64, aid int64, ak string) (prompt bool, err error) {
	var a *api.Arc
	if a, err = s.arcDao.Archive(c, aid); err != nil {
		if errors.Cause(err) == ecode.NothingFound {
			err = ecode.ArchiveNotExist
		}
		return
	}
	if !a.IsNormal() {
		err = ecode.ArchiveNotExist
		return
	}
	if err = s.favDao.AddVideo(c, mid, fids, aid, ak); err != nil {
		return
	}
	if prompt, err = s.relDao.Prompt(c, mid, vmid, _promptFav); err != nil {
		log.Error("%+v", err)
		err = nil
	}
	return
}

// Paster get paster if nologin.
func (s *Service) Paster(c context.Context, plat, adType int8, aid, typeID, buvid string) (p *resource.Paster, err error) {
	if p, err = s.rscDao.Paster(c, plat, adType, aid, typeID, buvid); err != nil {
		log.Error("%+v", err)
	}
	return
}

// VipPlayURL get playurl token.
func (s *Service) VipPlayURL(c context.Context, aid, cid, mid int64) (res *view.VipPlayURL, err error) {
	var (
		a    *api.Arc
		card *account.Card
	)
	res = &view.VipPlayURL{
		From: "app",
		Ts:   time.Now().Unix(),
		Aid:  aid,
		Cid:  cid,
		Mid:  mid,
	}
	if card, err = s.accDao.Card3(c, mid); err != nil {
		log.Error("%+v", err)
		err = ecode.AccessDenied
		return
	}
	if res.VIP = int(card.Level); res.VIP > 6 {
		res.VIP = 6
	}
	if card.Vip.Type != 0 && card.Vip.Status == 1 {
		res.SVIP = 1
	}
	if a, err = s.arcDao.Archive(c, aid); err != nil {
		log.Error("%+v", err)
		err = ecode.NothingFound
		return
	}
	if mid == a.Author.Mid {
		res.Owner = 1
	}
	params := url.Values{}
	params.Set("from", res.From)
	params.Set("ts", strconv.FormatInt(res.Ts, 10))
	params.Set("aid", strconv.FormatInt(res.Aid, 10))
	params.Set("cid", strconv.FormatInt(res.Cid, 10))
	params.Set("mid", strconv.FormatInt(res.Mid, 10))
	params.Set("vip", strconv.Itoa(res.VIP))
	params.Set("svip", strconv.Itoa(res.SVIP))
	params.Set("owner", strconv.Itoa(res.Owner))
	tmp := params.Encode()
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	mh := md5.Sum([]byte(strings.ToLower(tmp) + s.c.PlayURL.Secret))
	res.Fcs = hex.EncodeToString(mh[:])
	return
}

// Follow get auto follow switch from creative and acc.
func (s *Service) Follow(c context.Context, vmid, mid int64) (res *creative.PlayerFollow, err error) {
	var (
		fl      bool
		profile *relation.Stat
		fs      *creative.FollowSwitch
	)
	g, ctx := errgroup.WithContext(c)
	g.Go(func() (err error) {
		profile, err = s.relDao.Stat(ctx, vmid)
		if err != nil {
			log.Error("%+v", err)
		}
		return
	})
	if mid > 0 {
		g.Go(func() (err error) {
			fl, err = s.accDao.Following3(ctx, mid, vmid)
			if err != nil {
				log.Error("%+v", err)
			}
			return
		})
	}
	g.Go(func() (err error) {
		fs, err = s.creativeDao.FollowSwitch(ctx, vmid)
		if err != nil {
			log.Error("%+v", err)
		}
		return
	})
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
		return
	}
	res = &creative.PlayerFollow{}
	if profile.Follower >= s.c.AutoLimit && fs.State == 1 && !fl {
		res.Show = true
	}
	return
}

// UpperRecmd is
func (s *Service) UpperRecmd(c context.Context, plat int8, platform, mobiApp, device, buvid string, build int, mid, vimd int64) (res card.Handler, err error) {
	var (
		upIDs   []int64
		follow  *operate.Card
		cardm   map[int64]*account.Card
		statm   map[int64]*relation.Stat
		isAtten map[int64]int8
	)
	if follow, err = s.searchFollow(c, platform, mobiApp, device, buvid, build, mid, vimd); err != nil {
		log.Error("%+v", err)
		return
	}
	if follow == nil {
		err = ecode.AppNotData
		log.Error("follow is nil")
		return
	}
	for _, item := range follow.Items {
		upIDs = append(upIDs, item.ID)
	}
	g, ctx := errgroup.WithContext(c)
	if len(upIDs) != 0 {
		g.Go(func() (err error) {
			if cardm, err = s.accDao.Cards3(ctx, upIDs); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
		g.Go(func() (err error) {
			if statm, err = s.relDao.Stats(ctx, upIDs); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
		if mid != 0 {
			g.Go(func() error {
				isAtten = s.accDao.IsAttention(ctx, upIDs, mid)
				return nil
			})
		}
	}
	g.Wait()
	op := &operate.Card{}
	op.From(cdm.CardGt(model.GotoSearchUpper), 0, 0, plat, build)
	h := card.Handle(plat, cdm.CardGt(model.GotoSearchUpper), "", cdm.ColumnSvrSingle, nil, nil, isAtten, statm, cardm)
	if h == nil {
		err = ecode.AppNotData
		return
	}
	op = follow
	h.From(nil, op)
	if h.Get().Right {
		res = h
	} else {
		err = ecode.AppNotData
	}
	return
}

// LikeTriple like & coin & fav
func (s *Service) LikeTriple(c context.Context, aid, mid int64, ak string) (res *view.TripleRes, err error) {
	res = new(view.TripleRes)
	maxCoin := int64(1)
	multiply := int64(1)
	var a *api.Arc
	if a, err = s.arcDao.Archive(c, aid); err != nil {
		if errors.Cause(err) == ecode.NothingFound {
			err = ecode.ArchiveNotExist
		}
		return
	}
	if !a.IsNormal() {
		err = ecode.ArchiveNotExist
		return
	}
	if a.Copyright == int32(archive.CopyrightOriginal) {
		maxCoin = 2
		multiply = 2
	}
	res.UpID = a.Author.Mid
	eg := egV2.WithContext(c)
	eg.Go(func(ctx context.Context) (err error) {
		if multiply == 2 {
			userCoins, _ := s.coinDao.UserCoins(ctx, mid)
			if userCoins < 1 {
				return
			}
			if userCoins < 2 {
				multiply = 1
			}
		}
		err = s.coinDao.AddCoins(ctx, aid, mid, a.Author.Mid, maxCoin, _coinAv, multiply, int16(a.TypeID), int64(a.PubDate))
		if err == nil || ecode.EqualError(ecode.CoinOverMax, err) {
			res.Multiply = multiply
			res.Anticheat = true
			res.Coin = true
			err = nil
		} else {
			log.Error("s.coinDao.AddCoins err(%+v) aid(%d) mid(%d)", err, aid, mid)
			err = nil
			arcUserCoins, _ := s.coinDao.ArchiveUserCoins(ctx, aid, mid, _coinAv)
			if arcUserCoins != nil && arcUserCoins.Multiply > 0 {
				res.Coin = true
			}
		}
		return
	})
	eg.Go(func(ctx context.Context) (err error) {
		var isFav bool
		if isFav, err = s.favDao.IsFavVideo(ctx, mid, aid); err != nil {
			log.Error("s.favDao.IsFavVideo err(%+v) aid(%d) mid(%d)", err, aid, mid)
			err = nil
		} else if isFav {
			res.Fav = true
			return
		}
		if err = s.favDao.AddFav(ctx, mid, aid); err != nil {
			log.Error("s.favDao.AddFav err(%+v) aid(%d) mid(%d)", err, aid, mid)
			if ecode.EqualError(ecode.FavVideoExist, err) || ecode.EqualError(ecode.FavResourceExist, err) {
				res.Fav = true
			}
			err = nil
		} else {
			res.Fav = true
			res.Anticheat = true
		}
		return
	})
	eg.Go(func(ctx context.Context) (err error) {
		if _, err = s.thumbupDao.Like(ctx, mid, res.UpID, _businessLike, aid, thumbup.TypeLike); err != nil {
			log.Error("s.thumbupDao.Like err(%+v) aid(%d) mid(%d)", err, aid, mid)
			if ecode.EqualError(ecode.ThumbupDupLikeErr, err) {
				res.Like = true
			}
			err = nil
		} else {
			res.Like = true
			res.Anticheat = true
		}
		return
	})
	eg.Go(func(ctx context.Context) (err error) {
		if res.Prompt, err = s.relDao.Prompt(ctx, mid, a.Author.Mid, _promptFav); err != nil {
			log.Error("s.relDao.Prompt err(%+v)", err)
			err = nil
		}
		return
	})
	eg.Wait()
	if !res.Coin && !res.Fav {
		res.Prompt = false
	}
	return
}
