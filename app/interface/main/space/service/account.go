package service

import (
	"context"
	"encoding/json"

	"go-common/app/interface/main/space/model"
	tagmdl "go-common/app/interface/main/tag/model"
	artmdl "go-common/app/interface/openplatform/article/model"
	accwar "go-common/app/service/main/account/api"
	accmdl "go-common/app/service/main/account/model"
	favmdl "go-common/app/service/main/favorite/model"
	memmdl "go-common/app/service/main/member/model"
	relmdl "go-common/app/service/main/relation/model"
	upmdl "go-common/app/service/main/up/api/v1"
	"go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

const (
	_samplePn        = 1
	_samplePs        = 1
	_silenceForbid   = 1
	_accBlockDefault = 0
	_accBlockDue     = 1
	_officialNoType  = -1
	_audioCardOn     = 1
	_noticeForbid    = 1
)

var (
	_emptyThemeList = make([]*model.ThemeDetail, 0)
	_emptyArcItem   = make([]*model.ArcItem, 0)
)

// NavNum get space nav num by mid.
func (s *Service) NavNum(c context.Context, mid, vmid int64) (res *model.NavNum) {
	ip := metadata.String(c, metadata.RemoteIP)
	res = new(model.NavNum)
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		if reply, err := s.upClient.UpCount(errCtx, &upmdl.UpCountReq{Mid: vmid}); err != nil {
			log.Error("s.upClient.UpCount(%d) error(%v)", vmid, err)
		} else if reply.Count > 0 {
			res.Video = reply.Count
		}
		return nil
	})
	group.Go(func() error {
		res.Channel = new(model.Num)
		if chs, err := s.ChannelList(errCtx, vmid, false); err != nil {
			log.Error("s.ChannelList(%d) error(%v)", vmid, err)
		} else if chCnt := len(chs); chCnt > 0 {
			res.Channel.Master = chCnt
			for _, v := range chs {
				if v.Count > 0 {
					res.Channel.Guest++
				}
			}
		}
		return nil
	})
	group.Go(func() error {
		res.Favourite = new(model.Num)
		if favs, err := s.dao.FavFolder(errCtx, mid, vmid); err != nil {
			log.Error("s.dao.FavFolder(%d) error(%v)", vmid, err)
		} else if favCnt := len(favs); favCnt > 0 {
			res.Favourite.Master = favCnt
			for _, v := range favs {
				if v.IsPublic() {
					res.Favourite.Guest++
				}
			}
		}
		return nil
	})
	group.Go(func() error {
		if _, cnt, err := s.dao.BangumiList(errCtx, vmid, _samplePn, _samplePs); err != nil {
			log.Error("s.dao.BangumiList(%d) error(%v)", vmid, err)
		} else if cnt > 0 {
			res.Bangumi = cnt
		}
		return nil
	})
	group.Go(func() error {
		if tag, err := s.tag.SubTags(errCtx, &tagmdl.ArgSub{Mid: vmid, Pn: _samplePn, Ps: _samplePs, RealIP: ip}); err != nil {
			log.Error("s.tag.SubTags(%d) error(%v)", vmid, err)
		} else if tag != nil {
			res.Tag = tag.Total
		}
		return nil
	})
	group.Go(func() error {
		if art, err := s.art.UpArtMetas(errCtx, &artmdl.ArgUpArts{Mid: vmid, Pn: 1, Ps: 10, RealIP: ip}); err != nil {
			log.Error("s.art.UpArtMetas(%d) error(%v)", vmid, err)
		} else if art != nil {
			res.Article = art.Count
		}
		return nil
	})
	group.Go(func() error {
		if cnt, err := s.fav.CntUserFolders(errCtx, &favmdl.ArgCntUserFolders{Type: favmdl.TypePlayVideo, Mid: vmid, RealIP: ip}); err != nil {
			log.Error("s.dao.Playlist(%d) error(%v)", vmid, err)
		} else if cnt > 0 {
			res.Playlist = cnt
		}
		return nil
	})
	group.Go(func() error {
		if cnt, err := s.dao.AlbumCount(errCtx, vmid); err == nil && cnt > 0 {
			res.Album = cnt
		}
		return nil
	})
	group.Go(func() error {
		if cnt, err := s.dao.AudioCnt(errCtx, vmid); err != nil {
			log.Error("s.dao.AudioCnt(%d) error(%v)", vmid, err)
		} else if cnt > 0 {
			res.Audio = cnt
		}
		return nil
	})
	group.Wait()
	return
}

// UpStat get up stat.
func (s *Service) UpStat(c context.Context, mid int64) (res *model.UpStat, err error) {
	var (
		info           *accwar.InfoReply
		arcStat        *model.UpArcStat
		artStat        *model.UpArtStat
		arcErr, artErr error
	)
	if info, err = s.accClient.Info3(c, &accwar.MidReq{Mid: mid}); err != nil || info == nil {
		log.Error("s.accClient.Info3(%d) error(%v)", mid, err)
		return
	}
	res = new(model.UpStat)
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		if arcStat, arcErr = s.UpArcStat(errCtx, mid); arcErr != nil {
			log.Error("s.UpArcStat(%d) error(%v)", mid, arcErr)
		} else if arcStat != nil {
			res.Archive.View = arcStat.View
		}
		return nil
	})
	group.Go(func() error {
		if artStat, artErr = s.UpArtStat(errCtx, mid); artErr != nil {
			log.Error("s.UpArtStat(%d) error(%v)", mid, artErr)
		} else if artStat != nil {
			res.Article.View = artStat.View
		}
		return nil
	})
	group.Wait()
	return
}

// MyInfo get my info.
func (s *Service) MyInfo(c context.Context, mid int64) (res *accmdl.ProfileStat, err error) {
	var reply *accwar.ProfileStatReply
	if reply, err = s.accClient.ProfileWithStat3(c, &accwar.MidReq{Mid: mid}); err != nil {
		log.Error("s.accClient.ProfileWithStat3(%d) error(%v)", mid, err)
		return
	}
	level := memmdl.LevelInfo{
		Cur:     reply.LevelInfo.Cur,
		Min:     reply.LevelInfo.Min,
		NowExp:  reply.LevelInfo.NowExp,
		NextExp: reply.LevelInfo.NextExp,
	}
	res = &accmdl.ProfileStat{
		Profile:   reply.Profile,
		LevelExp:  level,
		Coins:     reply.Coins,
		Following: reply.Follower,
		Follower:  reply.Follower,
	}
	return
}

// AccTags get account tags.
func (s *Service) AccTags(c context.Context, mid int64) (res json.RawMessage, err error) {
	return s.dao.AccTags(c, mid)
}

// SetAccTags set account tags.
func (s *Service) SetAccTags(c context.Context, tags, ck string) (err error) {
	return s.dao.SetAccTags(c, tags, ck)
}

// AccInfo web acc info.
func (s *Service) AccInfo(c context.Context, mid, vmid int64) (res *model.AccInfo, err error) {
	if env.DeployEnv == env.DeployEnvProd {
		if _, ok := s.BlacklistValue[vmid]; ok {
			err = ecode.NothingFound
			return
		}
	}
	var (
		reply    *accwar.ProfileStatReply
		topPhoto *model.TopPhoto
		topErr   error
	)
	if reply, err = s.accClient.ProfileWithStat3(c, &accwar.MidReq{Mid: vmid}); err != nil || reply == nil {
		log.Error("s.accClient.ProfileWithStat3(%d) error(%v)", vmid, err)
		if ecode.Cause(err) != ecode.UserNotExist {
			return
		}
		reply = model.DefaultProfileStat
	}
	res = new(model.AccInfo)
	res.FromCard(reply)
	if res.Mid == 0 {
		res.Mid = vmid
	}
	group, errCtx := errgroup.WithContext(c)
	//check privacy
	if mid != vmid {
		group.Go(func() error {
			if privacyErr := s.privacyCheck(errCtx, vmid, model.PcyUserInfo); privacyErr != nil {
				res.JoinTime = 0
				res.Sex = "保密"
				res.Birthday = ""
			}
			return nil
		})
		if mid > 0 {
			group.Go(func() error {
				if relation, err := s.accClient.Relation3(errCtx, &accwar.RelationReq{Mid: mid, Owner: vmid}); err != nil {
					log.Error("s.accClient.Relation3(%d,%d) error (%v)", mid, vmid, err)
				} else if relation != nil {
					res.IsFollowed = relation.Following
				}
				return nil
			})
		}
	}
	//get top photo
	group.Go(func() error {
		topPhoto, topErr = s.dao.WebTopPhoto(errCtx, vmid)
		return nil
	})
	//get all theme
	group.Go(func() error {
		res.Theme = struct{}{}
		if theme, err := s.dao.Theme(errCtx, vmid); err == nil && theme != nil && len(theme.List) > 0 {
			for _, v := range theme.List {
				if v.IsActivated == 1 {
					res.Theme = v
					break
				}
			}
		}
		return nil
	})
	//get live metal
	group.Go(func() error {
		res.FansBadge, _ = s.dao.LiveMetal(errCtx, vmid)
		return nil
	})
	group.Wait()
	if topErr != nil || topPhoto == nil || topPhoto.LImg == "" {
		res.TopPhoto = s.c.Rule.TopPhoto
	} else {
		res.TopPhoto = topPhoto.LImg
	}
	return
}

// ThemeList get theme list.
func (s *Service) ThemeList(c context.Context, mid int64) (data []*model.ThemeDetail, err error) {
	var theme *model.ThemeDetails
	if theme, err = s.dao.Theme(c, mid); err != nil {
		return
	}
	if theme == nil || len(theme.List) == 0 {
		data = _emptyThemeList
		return
	}
	data = theme.List
	return
}

// ThemeActive theme active.
func (s *Service) ThemeActive(c context.Context, mid, themeID int64) (err error) {
	var (
		theme *model.ThemeDetails
		check bool
	)
	if theme, err = s.dao.Theme(c, mid); err != nil {
		return
	}
	if theme == nil || len(theme.List) == 0 {
		err = ecode.RequestErr
		return
	}
	for _, v := range theme.List {
		if v.ID == themeID {
			if v.IsActivated == 1 {
				err = ecode.NotModified
				return
			}
			check = true
		}
	}
	if !check {
		err = ecode.RequestErr
		return
	}
	if err = s.dao.ThemeActive(c, mid, themeID); err == nil {
		s.dao.DelCacheTheme(c, mid)
	}
	return
}

// Relation .
func (s *Service) Relation(c context.Context, mid, vmid int64) (data *model.Relation) {
	data = &model.Relation{Relation: struct{}{}, BeRelation: struct{}{}}
	ip := metadata.String(c, metadata.RemoteIP)
	if mid == vmid {
		return
	}
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		if relation, err := s.relation.Relation(errCtx, &relmdl.ArgRelation{Mid: mid, Fid: vmid, RealIP: ip}); err != nil {
			log.Error("Relation s.relation.Relation(Mid:%d,Fid:%d,%s) error %v", mid, vmid, ip, err)
		} else if relation != nil {
			data.Relation = relation
		}
		return nil
	})
	group.Go(func() error {
		if beRelation, err := s.relation.Relation(errCtx, &relmdl.ArgRelation{Mid: vmid, Fid: mid, RealIP: ip}); err != nil {
			log.Error("Relation s.relation.Relation(Mid:%d,Fid:%d,%s) error %v", vmid, mid, ip, err)
		} else if beRelation != nil {
			data.BeRelation = beRelation
		}
		return nil
	})
	group.Wait()
	return
}

// WebIndex web index.
func (s *Service) WebIndex(c context.Context, mid, vmid int64, pn, ps int32) (data *model.WebIndex, err error) {
	data = new(model.WebIndex)
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		info, infoErr := s.AccInfo(errCtx, mid, vmid)
		if infoErr != nil {
			return infoErr
		}
		data.Account = info
		return nil
	})
	group.Go(func() error {
		if setting, e := s.SettingInfo(errCtx, vmid); e == nil {
			data.Setting = setting
		}
		return nil
	})
	group.Go(func() error {
		if upArc, e := s.UpArcs(errCtx, vmid, pn, ps); e != nil {
			data.Archive = &model.WebArc{Archives: _emptyArcItem}
		} else {
			arc := &model.WebArc{
				Page:     model.WebPage{Pn: pn, Ps: ps, Count: upArc.Count},
				Archives: upArc.List,
			}
			data.Archive = arc
		}
		return nil
	})
	err = group.Wait()
	return
}
