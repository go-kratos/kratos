package space

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/app-interface/model"
	"go-common/app/interface/main/app-interface/model/audio"
	"go-common/app/interface/main/app-interface/model/bangumi"
	"go-common/app/interface/main/app-interface/model/bplus"
	"go-common/app/interface/main/app-interface/model/community"
	"go-common/app/interface/main/app-interface/model/favorite"
	"go-common/app/interface/main/app-interface/model/shop"
	"go-common/app/interface/main/app-interface/model/space"
	article "go-common/app/interface/openplatform/article/model"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	thumbup "go-common/app/service/main/thumbup/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"

	"github.com/pkg/errors"
)

const (
	_shopName     = "的店铺"
	_businessLike = "archive"
)

// Space aggregation space data.
func (s *Service) Space(c context.Context, mid, vmid int64, plat int8, build int, pn, ps int, platform, device, mobiApp, name string, now time.Time) (sp *space.Space, err error) {
	if _, ok := s.BlackList[vmid]; ok {
		err = ecode.NothingFound
		return
	}
	sp = &space.Space{}
	// get card
	card, err := s.card(c, vmid, name)
	if err != nil {
		if ecode.String(errors.Cause(err).Error()) == ecode.UserNotExist || vmid < 1 {
			err = ecode.NothingFound
			return
		}
		sp.Card = &space.Card{Mid: strconv.FormatInt(vmid, 10)}
		log.Error("%+v", err)
	} else if card == nil {
		log.Warn("space mid(%d) vmid(%d) name(%s) plat(%s) device(%s) card is nil", mid, vmid, name, platform, device)
		err = ecode.NothingFound
		return
	} else {
		if card.Mid == "" {
			err = ecode.NothingFound
			return
		}
		if vmid > 0 && card.Mid != strconv.FormatInt(vmid, 10) {
			err = ecode.NothingFound
			return
		}
		if vmid < 1 {
			if vmid, _ = strconv.ParseInt(card.Mid, 10, 64); vmid < 1 {
				err = ecode.NothingFound
				return
			}
		}
		sp.Card = card
		// ignore some card field
		sp.Card.Rank = ""
		sp.Card.DisplayRank = ""
		sp.Card.Regtime = 0
		sp.Card.Spacesta = 0
		sp.Card.Birthday = ""
		sp.Card.Place = ""
		sp.Card.Article = 0
		sp.Card.Friend = 0
		sp.Card.Attentions = nil
		sp.Card.Pendant.Pid = 0
		sp.Card.Pendant.Image = ""
		sp.Card.Pendant.Expire = 0
		sp.Card.Nameplate.Nid = 0
		sp.Card.Nameplate.Name = ""
		sp.Card.Nameplate.Image = ""
		sp.Card.Nameplate.ImageSmall = ""
		sp.Card.Nameplate.Level = ""
		sp.Card.Nameplate.Condition = ""
	}
	// space
	g, ctx := errgroup.WithContext(c)
	g.Go(func() error {
		sp.Space, _ = s.spcDao.SpaceMob(ctx, mid, vmid, platform, device)
		return nil
	})
	g.Go(func() error {
		sp.Card.FansGroup, _ = s.bplusDao.GroupsCount(ctx, mid, vmid)
		return nil
	})
	g.Go(func() (err error) {
		cardm, err := s.audioDao.Card(ctx, vmid)
		if err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		if card, ok := cardm[vmid]; ok && card.Type == 1 && card.Status == 1 {
			sp.Card.Audio = 1
		}
		return
	})
	if vmid == mid {
		g.Go(func() (err error) {
			if sp.Card.FansUnread, err = s.relDao.FollowersUnread(ctx, vmid); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	g.Go(func() (err error) {
		cert, err := s.audioDao.UpperCert(ctx, vmid)
		if err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		if cert == nil || cert.Cert == nil || cert.Cert.Type == -1 || cert.Cert.Desc == "" {
			return
		}
		if sp.Card.OfficialVerify.Type == -1 {
			sp.Card.OfficialVerify.Type = int8(cert.Cert.Type)
		}
		if sp.Card.OfficialVerify.Desc != "" {
			sp.Card.OfficialVerify.Desc = sp.Card.OfficialVerify.Desc + "、" + cert.Cert.Desc
		} else {
			sp.Card.OfficialVerify.Desc = cert.Cert.Desc
		}
		return
	})
	// iOS蓝版 强行去除充电模块
	if model.IsIPhone(plat) && (build >= 7000 && build <= 8000) {
		sp.Elec = nil
	} else {
		// elec rank
		g.Go(func() (err error) {
			info, err := s.elecDao.Info(ctx, vmid, mid)
			if err != nil {
				log.Error("%+v", err)
				err = nil
				return
			}
			if info != nil {
				info.Show = true
				sp.Elec = info
			}
			return
		})
	}
	g.Go(func() error {
		rel, _ := s.accDao.RichRelations3(ctx, mid, vmid)
		// -999:无关系, -1:我已经拉黑该用户, 0:我悄悄关注了该用户, 1:我公开关注了该用户
		// 1- 悄悄关注 2 关注  6-好友 128-拉黑
		if rel == 1 {
			sp.Relation = 0
		} else if rel == 2 || rel == 6 {
			sp.Relation = 1
		} else if rel >= 128 {
			sp.Relation = -1
		} else {
			sp.Relation = -999
		}
		return nil
	})
	g.Go(func() (err error) {
		if sp.Medal, err = s.liveDao.MedalStatus(ctx, vmid); err != nil {
			log.Error("%+v", err)
			err = nil
		}
		return
	})
	// g.Go(func() (err error) {
	// 	if sp.Attention, err = s.relDao.Attention(ctx, mid, vmid); err != nil {
	// 		log.Error("%+v", err)
	// 		err = nil
	// 	}
	// 	return
	// })
	g.Go(func() (err error) {
		if sp.Live, err = s.liveDao.Live(ctx, vmid, platform); err != nil {
			log.Error("%+v", err)
			err = nil
		}
		return
	})
	sp.Tab = &space.Tab{}
	// up archives
	g.Go(func() error {
		sp.Archive = s.UpArcs(ctx, vmid, pn, ps, now)
		if sp.Archive != nil && len(sp.Archive.Item) > 0 {
			sp.Tab.Archive = true
		}
		return nil
	})
	g.Go(func() (err error) {
		if sp.Tab.Dynamic, err = s.bplusDao.Dynamic(ctx, vmid); err != nil {
			log.Error("%+v", err)
			err = nil
		}
		return
	})
	g.Go(func() error {
		sp.Article = s.UpArticles(ctx, vmid, 0, 3)
		if sp.Article != nil && len(sp.Article.Item) != 0 {
			sp.Tab.Article = true
		}
		return nil
	})
	g.Go(func() error {
		sp.Clip = s.upClips(ctx, vmid)
		if sp.Clip != nil && len(sp.Clip.Item) != 0 {
			sp.Tab.Clip = true
		}
		return nil
	})
	g.Go(func() error {
		sp.Album = s.upAlbums(ctx, vmid)
		if sp.Album != nil && len(sp.Album.Item) != 0 {
			sp.Tab.Album = true
		}
		return nil
	})
	g.Go(func() error {
		sp.Audios = s.audios(ctx, vmid, 1, 3)
		if sp.Audios != nil && len(sp.Audios.Item) != 0 {
			sp.Tab.Audios = true
		}
		return nil
	})
	g.Go(func() (err error) {
		var info *shop.Info
		if info, err = s.shopDao.Info(ctx, vmid, mobiApp, device, build); err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		if info != nil && info.Shop != nil && info.Shop.Status == 2 {
			sp.Tab.Shop = true
			sp.Shop = &space.Shop{ID: info.Shop.ID, Name: info.Shop.Name + _shopName}
		}
		return
	})
	if sp.Setting, err = s.spcDao.Setting(c, vmid); err != nil {
		log.Error("%+v", err)
		err = nil
	}
	if sp.Setting == nil {
		err = g.Wait()
		return
	}
	g.Go(func() error {
		sp.Favourite = s.favFolders(ctx, mid, vmid, sp.Setting, plat, build, mobiApp)
		if sp.Favourite != nil && len(sp.Favourite.Item) != 0 {
			sp.Tab.Favorite = true
		}
		return nil
	})
	g.Go(func() error {
		sp.Season = s.Bangumi(ctx, mid, vmid, sp.Setting, pn, ps)
		if sp.Season != nil && len(sp.Season.Item) != 0 {
			sp.Tab.Bangumi = true
		}
		return nil
	})
	g.Go(func() error {
		sp.CoinArc = s.CoinArcs(ctx, mid, vmid, sp.Setting, pn, ps)
		if sp.CoinArc != nil && len(sp.CoinArc.Item) != 0 {
			sp.Tab.Coin = true
		}
		return nil
	})
	g.Go(func() error {
		sp.LikeArc = s.LikeArcs(ctx, mid, vmid, sp.Setting, pn, ps)
		if sp.LikeArc != nil && len(sp.LikeArc.Item) != 0 {
			sp.Tab.Like = true
		}
		return nil
	})
	err = g.Wait()
	return
}

// UpArcs get upload archive .
func (s *Service) UpArcs(c context.Context, uid int64, pn, ps int, now time.Time) (res *space.ArcList) {
	var (
		arcs []*api.Arc
		err  error
	)
	res = &space.ArcList{Item: []*space.ArcItem{}}
	if res.Count, err = s.arcDao.UpCount2(c, uid); err != nil {
		log.Error("%+v", err)
	}
	if arcs, err = s.arcDao.UpArcs3(c, uid, pn, ps); err != nil {
		log.Error("%+v", err)
		return
	}
	if len(arcs) != 0 {
		res.Item = make([]*space.ArcItem, 0, len(arcs))
		for _, v := range arcs {
			si := &space.ArcItem{}
			si.FromArc(v)
			res.Item = append(res.Item, si)
		}
	}
	return
}

// UpArticles get article.
func (s *Service) UpArticles(c context.Context, uid int64, pn, ps int) (res *space.ArticleList) {
	res = &space.ArticleList{Item: []*space.ArticleItem{}, Lists: []*article.List{}}
	g, ctx := errgroup.WithContext(c)
	g.Go(func() (err error) {
		var ams []*article.Meta
		if ams, res.Count, err = s.artDao.UpArticles(ctx, uid, pn, ps); err != nil {
			return err
		}
		if len(ams) != 0 {
			res.Item = make([]*space.ArticleItem, 0, len(ams))
			for _, v := range ams {
				if v.AttrVal(article.AttrBitNoDistribute) {
					continue
				}
				si := &space.ArticleItem{}
				si.FromArticle(v)
				res.Item = append(res.Item, si)
			}
		}
		return err
	})
	g.Go(func() (err error) {
		var lists []*article.List
		lists, res.ListsCount, err = s.artDao.UpLists(c, uid)
		if err != nil {
			return err
		}
		if len(lists) > 0 {
			res.Lists = lists
		}
		return err
	})
	if err := g.Wait(); err != nil {
		log.Error("%+v", err)
	}
	return
}

// favFolders get favorite folders
func (s *Service) favFolders(c context.Context, mid, vmid int64, setting *space.Setting, plat int8, build int, mobiApp string) (res *space.FavList) {
	const (
		_oldAndroidBuild = 427100
		_oldIOSBuild     = 3910
	)
	var (
		fs  []*favorite.Folder
		err error
	)
	res = &space.FavList{Item: []*favorite.Folder{}}
	if mid != vmid {
		if setting == nil {
			setting, err = s.spcDao.Setting(c, vmid)
			if err != nil {
				log.Error("%+v", err)
				return
			}
		}
		if setting.FavVideo != 1 {
			return
		}
	}
	var mediaList bool
	// 双端版本号限制，符合此条件显示为“默认收藏夹”：
	// iPhone <5.36.1(8300) 或iPhone>5.36.1(8300)
	// Android <5360001或Android>5361000
	// 双端版本号限制，符合此条件显示为“默认播单”：
	// iPhone=5.36.1(8300)
	// 5360001 <=Android <=5361000
	if (plat == model.PlatIPhone && build == 8300) || (plat == model.PlatAndroid && build >= 5360001 && build <= 5361000) {
		mediaList = true
	}
	if fs, err = s.favDao.Folders(c, mid, vmid, mobiApp, build, mediaList); err != nil {
		log.Error("%+v", err)
		return
	}
	for _, v := range fs {
		if ((plat == model.PlatAndroid || plat == model.PlatAndroidG) && build <= _oldAndroidBuild) || ((plat == model.PlatIPhone || plat == model.PlatIPhoneI) && build <= _oldIOSBuild) {
			v.Videos = v.Cover
			v.Cover = nil
		}
	}
	res.Item = fs
	res.Count = len(fs)
	return
}

// Bangumi get concern season
func (s *Service) Bangumi(c context.Context, mid, vmid int64, setting *space.Setting, pn, ps int) (res *space.BangumiList) {
	var (
		seasons []*bangumi.Season
		err     error
	)
	res = &space.BangumiList{Item: []*space.BangumiItem{}}
	if mid != vmid {
		if setting == nil {
			setting, err = s.spcDao.Setting(c, vmid)
			if err != nil {
				log.Error("%+v", err)
				return
			}
		}
		if setting.Bangumi != 1 {
			return
		}
	}
	if seasons, res.Count, err = s.bgmDao.Concern(c, mid, vmid, pn, ps); err != nil {
		log.Error("%+v", err)
		return
	}
	if len(seasons) != 0 {
		res.Item = make([]*space.BangumiItem, 0, len(seasons))
		for _, v := range seasons {
			si := &space.BangumiItem{}
			si.FromSeason(v)
			res.Item = append(res.Item, si)
		}
	}
	return
}

// Community get community
func (s *Service) Community(c context.Context, uid int64, pn, ps int, ak, platform string) (res *space.CommuList) {
	var (
		comm []*community.Community
		err  error
	)
	res = &space.CommuList{Item: []*space.CommItem{}}
	if comm, res.Count, err = s.commDao.Community(c, uid, ak, platform, pn, ps); err != nil {
		log.Error("%+v", err)
		return
	}
	if len(comm) != 0 {
		res.Item = make([]*space.CommItem, 0, len(comm))
		for _, v := range comm {
			si := &space.CommItem{}
			si.FromCommunity(v)
			res.Item = append(res.Item, si)
		}
	}
	return
}

// CoinArcs get coin archives.
func (s *Service) CoinArcs(c context.Context, mid, vmid int64, setting *space.Setting, pn, ps int) (res *space.ArcList) {
	var (
		coins []*api.Arc
		err   error
	)
	res = &space.ArcList{Item: []*space.ArcItem{}}
	if mid != vmid {
		if setting == nil {
			setting, err = s.spcDao.Setting(c, vmid)
			if err != nil {
				log.Error("%+v", err)
				return
			}
		}
		if setting.CoinsVideo != 1 {
			return
		}
	}
	if coins, res.Count, err = s.coinDao.CoinList(c, vmid, pn, ps); err != nil {
		log.Error("%+v", err)
		return
	}
	if len(coins) != 0 {
		res.Item = make([]*space.ArcItem, 0, len(coins))
		for _, v := range coins {
			si := &space.ArcItem{}
			si.FromCoinArc(v)
			res.Item = append(res.Item, si)
		}
	}
	return
}

// LikeArcs get like archives.
func (s *Service) LikeArcs(c context.Context, mid, vmid int64, setting *space.Setting, pn, ps int) (res *space.ArcList) {
	var (
		likes []*thumbup.ItemLikeRecord
		err   error
	)
	res = &space.ArcList{Item: []*space.ArcItem{}}
	if mid != vmid {
		if setting == nil {
			setting, err = s.spcDao.Setting(c, vmid)
			if err != nil {
				log.Error("%+v", err)
				return
			}
		}
		if setting.LikesVideo != 1 {
			return
		}
	}
	if likes, res.Count, err = s.thumbupDao.UserTotalLike(c, vmid, _businessLike, pn, ps); err != nil {
		log.Error("%+v", err)
		return
	}
	if len(likes) != 0 {
		aids := make([]int64, 0, len(likes))
		for _, v := range likes {
			aids = append(aids, v.MessageID)
		}
		var as map[int64]*api.Arc
		if as, err = s.arcDao.Archives2(c, aids); err != nil {
			log.Error("%+v", err)
			return
		}
		if len(as) == 0 {
			return
		}
		res.Item = make([]*space.ArcItem, 0, len(as))
		for _, v := range likes {
			if a, ok := as[v.MessageID]; ok {
				si := &space.ArcItem{}
				si.FromLikeArc(a)
				res.Item = append(res.Item, si)
			}
		}
	}
	return
}

func (s *Service) upClips(c context.Context, uid int64) (res *space.ClipList) {
	var (
		clips []*bplus.Clip
		err   error
	)
	res = &space.ClipList{Item: []*space.Item{}}
	if clips, res.Count, err = s.bplusDao.AllClip(c, uid, 200); err != nil {
		log.Error("%+v", err)
		return
	}
	if len(clips) != 0 {
		res.Item = make([]*space.Item, 0, len(clips))
		for k, v := range clips {
			si := &space.Item{}
			si.FromClip(v)
			res.Item = append(res.Item, si)
			if k == 1 {
				break
			}
		}
	}
	return
}

func (s *Service) upAlbums(c context.Context, uid int64) (res *space.AlbumList) {
	var (
		album []*bplus.Album
		err   error
	)
	res = &space.AlbumList{Item: []*space.Item{}}
	if album, res.Count, err = s.bplusDao.AllAlbum(c, uid, 6); err != nil {
		log.Error("%+v", err)
		return
	}
	if len(album) != 0 {
		res.Item = make([]*space.Item, 0, len(album))
		for k, v := range album {
			si := &space.Item{}
			si.FromAlbum(v)
			res.Item = append(res.Item, si)
			if k == 5 {
				break
			}
		}
	}
	return
}

// card get card by mid, vmid or name.
func (s *Service) card(c context.Context, vmid int64, name string) (scard *space.Card, err error) {
	scard = &space.Card{}
	var profile *account.ProfileStat
	if vmid > 0 {
		profile, err = s.accDao.Profile3(c, vmid)
	} else if name != "" {
		profile, err = s.accDao.ProfileByName3(c, name)
	}
	if err != nil {
		err = errors.Wrapf(err, "%v,%v", vmid, name)
		return
	}
	scard.Mid = strconv.FormatInt(profile.Mid, 10)
	scard.Name = profile.Name
	// scard.Approve =
	scard.Sex = profile.Sex
	scard.Rank = strconv.FormatInt(int64(profile.Rank), 10)
	scard.Face = profile.Face
	scard.DisplayRank = strconv.FormatInt(int64(profile.Rank), 10)
	// scard.Regtime =
	if profile.Silence == 1 {
		scard.Spacesta = -2
	}
	// scard.Birthday =
	// scard.Place =
	scard.Description = profile.Official.Desc
	scard.Article = 0
	scard.Attentions = []int64{}
	scard.Fans = int(profile.Follower)
	// scard.Friend = profile.Following
	scard.Attention = int(profile.Following)
	scard.Sign = profile.Sign
	scard.LevelInfo.Cur = profile.Level
	scard.LevelInfo.Min = profile.LevelExp.Min
	scard.LevelInfo.NowExp = profile.LevelExp.NowExp
	scard.LevelInfo.NextExp = profile.LevelExp.NextExp
	if profile.LevelExp.NextExp == -1 {
		scard.LevelInfo.NextExp = "--"
	}
	scard.Pendant.Pid = profile.Pendant.Pid
	scard.Pendant.Name = profile.Pendant.Name
	scard.Pendant.Image = profile.Pendant.Image
	scard.Pendant.Expire = profile.Pendant.Expire
	scard.OfficialVerify.Role = profile.Official.Role
	scard.OfficialVerify.Title = profile.Official.Title
	if profile.Official.Role == 0 {
		scard.OfficialVerify.Type = -1
	} else {
		if profile.Official.Role <= 2 {
			scard.OfficialVerify.Type = 0
		} else {
			scard.OfficialVerify.Type = 1
		}
		scard.OfficialVerify.Desc = profile.Official.Title
	}
	scard.Vip.Type = int(profile.Vip.Type)
	scard.Vip.VipStatus = int(profile.Vip.Status)
	scard.Vip.DueDate = profile.Vip.DueDate
	// scard.FansGroup =
	// scard.Audio =
	// scard.FansUnread =
	return
}

// audios
func (s *Service) audios(c context.Context, mid int64, pn, ps int) (res *space.AudioList) {
	var (
		audios []*audio.Audio
		err    error
	)
	res = &space.AudioList{Item: []*space.AudioItem{}}
	if audios, res.Count, err = s.audioDao.Audios(c, mid, pn, ps); err != nil {
		log.Error("%+v", err)
		return
	}
	if len(audios) != 0 {
		res.Item = make([]*space.AudioItem, 0, len(audios))
		for _, v := range audios {
			si := &space.AudioItem{}
			si.FromAudio(v)
			res.Item = append(res.Item, si)
		}
	}
	return
}

// Report func
func (s *Service) Report(c context.Context, mid int64, reason, ak string) (err error) {
	return s.spcDao.Report(c, mid, reason, ak)
}
