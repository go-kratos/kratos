package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	dm2mdl "go-common/app/interface/main/dm2/model"
	tagmdl "go-common/app/interface/main/tag/model"
	"go-common/app/interface/main/web/model"
	accmdl "go-common/app/service/main/account/api"
	infomdl "go-common/app/service/main/account/model"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	coinmdl "go-common/app/service/main/coin/api"
	favmdl "go-common/app/service/main/favorite/model"
	locmdl "go-common/app/service/main/location/model"
	sharemdl "go-common/app/service/main/share/api"
	thumbup "go-common/app/service/main/thumbup/model"
	ugcmdl "go-common/app/service/main/ugcpay/api/grpc/v1"
	"go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/log/infoc"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

var (
	_emptyTags     = make([]*tagmdl.Tag, 0)
	_emptyReplyHot = new(model.ReplyHot)
	_emptyArchive3 = make([]*arcmdl.Arc, 0)
	_emptyArc      = make([]*arcmdl.Arc, 0)
)

const (
	_japan            = "日本"
	_businessAppeal   = 1
	_notForward       = 0
	_member           = 10000
	_viewBakCacheRand = 10
	_shareArcType     = 3
	_hasUGCPay        = 1
	_ugcOtypeArc      = "archive"
	_ugcCurrencyBp    = "bp"
	_ugcAssetPaid     = "paid"
	_ugcPaidState     = 1
)

// View get a  view page by aid.
func (s *Service) View(c context.Context, aid, cid, mid int64, cdnIP, ck string) (rs *model.View, err error) {
	var (
		viewReply              *arcmdl.ViewReply
		video                  *arcmdl.Page
		longDesc               string
		pDesc, forward, report bool
		remoteIP               = metadata.String(c, metadata.RemoteIP)
	)
	if viewReply, err = s.arcClient.View(c, &arcmdl.ViewRequest{Aid: aid}); err != nil {
		log.Error("s.arcClient.View(%d) error %v", aid, err)
		return
	}
	if viewReply != nil && viewReply.Arc != nil {
		if _, ok := s.specialMids[viewReply.Arc.Author.Mid]; ok && env.DeployEnv == env.DeployEnvProd {
			err = ecode.NothingFound
			log.Warn("aid(%d) mid(%d) can not view on prod", viewReply.Arc.Aid, viewReply.Arc.Author.Mid)
			return
		}
		if viewReply.Arc.FirstCid == 0 {
			if len(viewReply.Pages) > 0 {
				viewReply.Arc.FirstCid = viewReply.Pages[0].Cid
			}
		}
		if viewReply.Arc.State == archive.StateForbidLock && viewReply.Arc.Forward > 0 {
			forward = true
		} else if viewReply.Arc.State == archive.StateForbidRecicle || viewReply.Arc.State == archive.StateForbidLock {
			if viewReply.Arc.ReportResult == "" {
				err = ecode.NothingFound
				return
			}
			report = true
		}
		check := []int32{archive.StateForbidWait, archive.StateForbidFixed, archive.StateForbidLater, archive.StateForbidAdminDelay}
		for _, v := range check {
			if viewReply.Arc.State == v {
				err = ecode.ArchiveChecking
				return
			}
		}
		if viewReply.Arc.State == archive.StateForbidUserDelay {
			err = ecode.ArchivePass
			return
		} else if !viewReply.Arc.IsNormal() && !forward && !report {
			err = ecode.ArchiveDenied
			return
		}
		rs = &model.View{Arc: viewReply.Arc, Pages: viewReply.Pages}
		group, errCtx := errgroup.WithContext(c)
		if !forward && !report {
			group.Go(func() error {
				return s.zlimit(errCtx, rs, mid, remoteIP, cdnIP)
			})
			group.Go(func() error {
				return s.checkAccess(errCtx, mid, rs, ck, remoteIP)
			})
		}
		group.Go(func() error {
			dmCid := cid
			if dmCid == 0 {
				dmCid = viewReply.Arc.FirstCid
			}
			rs.Subtitle = s.dmSubtitle(errCtx, aid, dmCid)
			return nil
		})
		if rs.Arc.Rights.UGCPay == _hasUGCPay {
			rs.Asset = s.ugcPayAsset(errCtx, aid, mid)
		}
		if err = group.Wait(); err != nil {
			return
		}
		if cid > 0 {
			if video, err = s.arc.Video3(c, &archive.ArgVideo2{Aid: aid, Cid: cid, RealIP: remoteIP}); err != nil {
				err = nil
				log.Error("s.arc.Video3(%d,%d,%s) error %v", aid, cid, remoteIP, err)
			}
			if video.Desc != "" {
				rs.Desc = video.Desc
				pDesc = true
			}
		}
		if !pDesc {
			if longDesc, err = s.arc.Description2(c, &archive.ArgAid{Aid: aid, RealIP: remoteIP}); err != nil {
				err = nil
				log.Error("view s.arc.Description2(%d) error(%v)", aid, err)
			} else if longDesc != "" {
				rs.Desc = longDesc
			}
		}
		if rs.AttrVal(archive.AttrBitIsPGC) == archive.AttrNo {
			if err = s.initDownload(c, rs, mid, remoteIP, cdnIP); err != nil {
				log.Error("s.initDownload  aid(%d) mid(%d) ip(%s) error(%+v)", rs.Aid, mid, remoteIP, err)
				err = nil
			}
		}
		if !forward {
			rs.Forward = _notForward
		}
		if !report {
			rs.ReportResult = ""
		}
		if rs.AttrVal(archive.AttrBitJumpUrl) == archive.AttrNo {
			rs.RedirectURL = ""
		}
		if rs.Access >= _member {
			rs.Stat.View = -1
		}
		if rs.AttrVal(archive.AttrBitLimitArea) == archive.AttrYes {
			rs.NoCache = true
		}
		s.cache.Do(c, func(c context.Context) {
			if s.r.Intn(_viewBakCacheRand) == 1 {
				s.dao.SetViewBakCache(c, aid, rs)
			}
		})
	} else {
		rs, err = s.dao.ViewBakCache(c, aid)
	}
	return
}

func (s *Service) zlimit(c context.Context, view *model.View, mid int64, remoteIP, cdnIP string) (err error) {
	var data *int64
	if data, err = s.loc.Archive(c, &locmdl.Archive{Aid: view.Aid, Mid: mid, IP: remoteIP, CIP: cdnIP}); err != nil {
		log.Error("s.loc.Archive(%d%d%s%s) error(%v)", view.Aid, mid, remoteIP, cdnIP, err)
		return
	}
	if *data == locmdl.Forbidden {
		log.Warn("s.loc.Archive aid(%d) zlimit.Forbidden", view.Aid)
		err = ecode.NothingFound
	} else {
		err = s.specialLimit(c, view, remoteIP)
	}
	return
}

// specialLimit spacialLimit special type id limit in japan
func (s *Service) specialLimit(c context.Context, view *model.View, remoteIP string) (err error) {
	var zone *locmdl.Info
	for _, typeID := range model.LimitTypeIDs {
		if int32(typeID) == view.TypeID {
			// TODO all no cache
			view.NoCache = true
			if zone, err = s.loc.Info(c, &locmdl.ArgIP{IP: remoteIP}); err != nil || zone == nil {
				log.Error("s.loc.Info(%s) error(%v) or zone is nil", remoteIP, err)
				err = nil
			} else if zone.Country == _japan {
				err = ecode.NothingFound
			}
			break
		}
	}
	return
}

// checkAccess check mid aid access
func (s *Service) checkAccess(c context.Context, mid int64, view *model.View, ck, ip string) (err error) {
	var (
		p *accmdl.CardReply
	)
	if view.Access == 0 {
		return
	}
	view.NoCache = true
	if mid <= 0 {
		log.Warn("user not login  aid(%d)", view.Aid)
		err = ecode.AccessDenied
		return
	}
	if p, err = s.accClient.Card3(c, &accmdl.MidReq{Mid: mid}); err != nil {
		log.Error("s.accClient.Card3(%d) error(%v)", mid, err)
		return
	}
	if p == nil {
		log.Warn("Info2 result is null aid(%d) state(%d) access(%d)", view.Aid, view.State, view.Access)
		err = ecode.AccessDenied
		return
	}
	card := p.Card
	isVip := (card.Vip.Type > 0) && (card.Vip.Status == 1)
	if view.Access > 0 && card.Rank < view.Access && (!isVip) {
		log.Warn("mid(%d) rank(%d) vip(tp:%d,status:%d) have not access(%d) view archive(%d) ", mid, card.Rank, card.Vip.Type, card.Vip.Status, view.Access, view.Aid)
		if mid > 0 {
			err = ecode.NothingFound
		} else {
			err = ecode.ArchiveNotLogin
		}
	}
	return
}

func (s *Service) initDownload(c context.Context, v *model.View, mid int64, ip, cdnIP string) (err error) {
	var download int64
	if v.AttrVal(archive.AttrBitLimitArea) == archive.AttrYes {
		if download, err = s.downLimit(c, mid, v.Aid, cdnIP); err != nil {
			return
		}
	} else {
		download = locmdl.AllowDown
	}
	if download == locmdl.ForbiddenDown {
		v.Rights.Download = int32(download)
		return
	}
	for _, p := range v.Pages {
		if p.From == "qq" {
			download = locmdl.ForbiddenDown
			break
		}
	}
	v.Rights.Download = int32(download)
	return
}

// downLimit ip limit
func (s *Service) downLimit(c context.Context, mid, aid int64, cdnIP string) (down int64, err error) {
	var (
		auth *locmdl.Auth
		ip   = metadata.String(c, metadata.RemoteIP)
	)
	if auth, err = s.loc.Archive2(c, &locmdl.Archive{Aid: aid, Mid: mid, IP: ip, CIP: cdnIP}); err != nil {
		log.Error("s.loc.Archive2(%d) error(%v)", mid, err)
		return
	}
	if auth.Play == locmdl.Forbidden {
		err = ecode.AccessDenied
	} else {
		down = auth.Down
	}
	return
}

func (s *Service) dmSubtitle(c context.Context, aid, cid int64) (subtitle *model.Subtitle) {
	var (
		dmSub      *dm2mdl.VideoSubtitles
		err        error
		mids       []int64
		infosReply *accmdl.InfosReply
		subs       []*model.SubtitleItem
	)
	subtitle = new(model.Subtitle)
	if dmSub, err = s.dm2.SubtitleGet(c, &dm2mdl.ArgSubtitleGet{Aid: aid, Oid: cid, Type: dm2mdl.SubTypeVideo}); err != nil {
		log.Warn("dmSubtitle s.dm2.SubtitleGet aid(%d) cid(%d) warn(%v)", aid, cid, err)
	} else if dmSub != nil {
		subtitle.AllowSubmit = dmSub.AllowSubmit
		if len(dmSub.Subtitles) > 0 {
			for _, v := range dmSub.Subtitles {
				if v.AuthorMid > 0 {
					mids = append(mids, v.AuthorMid)
				}
			}
			infoData := make(map[int64]*infomdl.Info)
			if len(mids) > 0 {
				if infosReply, err = s.accClient.Infos3(c, &accmdl.MidsReq{Mids: mids, RealIp: metadata.String(c, metadata.RemoteIP)}); err != nil {
					log.Error("dmSubtitle aid(%d) cid(%d) s.acc.Infos3 mids(%v) error(%v)", aid, cid, mids, err)
				} else {
					infoData = infosReply.Infos
				}
			}
			for _, v := range dmSub.Subtitles {
				sub := &model.SubtitleItem{VideoSubtitle: v, Author: &infomdl.Info{Mid: v.AuthorMid}}
				if info, ok := infoData[v.AuthorMid]; ok && info != nil {
					sub.Author = info
				}
				subs = append(subs, sub)
			}
			subtitle.List = subs
		}
	}
	if len(subtitle.List) == 0 {
		subtitle.List = make([]*model.SubtitleItem, 0)
	}
	return
}

// ArchiveStat get archive stat data by aid.
func (s *Service) ArchiveStat(c context.Context, aid int64) (stat *model.Stat, err error) {
	var (
		arcReply *arcmdl.ArcReply
		view     interface{}
	)
	if aid == s.c.Bnj2019.LiveAid && s.bnj2019LiveArc != nil {
		arcReply = s.bnj2019LiveArc
	} else {
		if arcReply, err = s.arcClient.Arc(c, &arcmdl.ArcRequest{Aid: aid}); err != nil {
			log.Error("s.arcClient.Arc(%d) error(%v)", aid, err)
			return
		}
	}
	arc := arcReply.Arc
	if !model.CheckAllowState(arc) {
		err = ecode.AccessDenied
		return
	}
	view = arc.Stat.View
	if arc.Access > 0 {
		view = "--"
	}
	stat = &model.Stat{
		Aid:       arc.Stat.Aid,
		View:      view,
		Danmaku:   arc.Stat.Danmaku,
		Reply:     arc.Stat.Reply,
		Fav:       arc.Stat.Fav,
		Coin:      arc.Stat.Coin,
		Share:     arc.Stat.Share,
		Like:      arc.Stat.Like,
		NowRank:   arc.Stat.NowRank,
		HisRank:   arc.Stat.HisRank,
		NoReprint: arc.Rights.NoReprint,
		Copyright: arc.Copyright,
	}
	return
}

// AddShare share add count
func (s *Service) AddShare(c context.Context, aid, mid int64, ua, refer, path, buvid, sid string) (shares int64, err error) {
	var (
		shareReply *sharemdl.AddShareReply
		remoteIP   = metadata.String(c, metadata.RemoteIP)
	)
	// add Anti cheat
	if s.CheatInfoc != nil {
		ac := map[string]string{
			"itemType": infoc.ItemTypeAv,
			"action":   infoc.ActionShare,
			"ip":       remoteIP,
			"mid":      strconv.FormatInt(mid, 10),
			"fid":      "",
			"aid":      strconv.FormatInt(aid, 10),
			"sid":      sid,
			"ua":       ua,
			"buvid":    buvid,
			"refer":    refer,
			"url":      path,
		}
		s.CheatInfoc.ServiceAntiCheat(ac)
	}
	if shareReply, err = s.shareClient.AddShare(c, &sharemdl.AddShareRequest{Oid: aid, Mid: mid, Type: _shareArcType, Ip: remoteIP}); err != nil {
		log.Error("AddShare s.shareClient.AddShare (oid:%d,mid:%d) warn(%v)", aid, mid, err)
		return
	}
	shares = shareReply.Shares
	return
}

// Description get archive description by aid.
func (s *Service) Description(c context.Context, aid, page int64) (res string, err error) {
	var (
		viewReply *arcmdl.ViewReply
		video     *arcmdl.Page
		cid       int64
		longDesc  string
		ip        = metadata.String(c, metadata.RemoteIP)
	)
	if viewReply, err = s.arcClient.View(c, &arcmdl.ViewRequest{Aid: aid}); err != nil {
		log.Error("s.arcClient.View(%d,%s) error %v", aid, ip, err)
		return
	}
	if viewReply != nil && viewReply.Arc != nil {
		if !viewReply.Arc.IsNormal() {
			err = ecode.ArchiveDenied
			return
		}
	}
	if page > 0 {
		if int(page-1) >= len(viewReply.Pages) || viewReply.Pages[page-1] == nil {
			err = ecode.NothingFound
			return
		}
		cid = viewReply.Pages[page-1].Cid
		if cid > 0 {
			if video, err = s.arc.Video3(c, &archive.ArgVideo2{Aid: aid, Cid: cid, RealIP: ip}); err != nil {
				log.Error("s.arc.Video2(%d,%d,%s) error %v", aid, cid, ip, err)
			}
			if video.Desc != "" {
				res = video.Desc
				return
			}
		}
	} else {
		res = viewReply.Arc.Desc
	}
	if longDesc, err = s.arc.Description2(c, &archive.ArgAid{Aid: aid, RealIP: ip}); err != nil {
		log.Error("s.arc.Description2(%d) error(%v)", aid, err)
	} else if longDesc != "" {
		res = longDesc
	}
	return
}

// ArcReport add archive report
func (s *Service) ArcReport(c context.Context, mid, aid, tp int64, reason, pics string) (err error) {
	if err = s.dao.ArcReport(c, mid, aid, tp, reason, pics); err != nil {
		log.Error("s.dao.ArcReport(%d,%d,%d,%s,%s) err (%v)", mid, aid, tp, reason, pics, err)
	}
	return
}

// AppealTags get appeal tags
func (s *Service) AppealTags(c context.Context) (rs json.RawMessage, err error) {
	if rs, err = s.dao.AppealTags(c, _businessAppeal); err != nil {
		log.Error("s.dao.AppealTags(1) error(%v)", err)
	}
	return
}

// ArcAppeal add archive appeal.
func (s *Service) ArcAppeal(c context.Context, mid int64, data map[string]string) (err error) {
	aid, _ := strconv.ParseInt(data["oid"], 10, 64)
	if err = s.dao.ArcAppealCache(c, mid, aid); err != nil {
		if err == ecode.ArcAppealLimit {
			log.Warn("s.arcAppealLimit mid(%d) aid(%d)", mid, aid)
			return
		}
		err = nil
	}
	if err = s.dao.ArcAppeal(c, mid, data, _businessAppeal); err != nil {
		log.Error("s.dao.ArcAppeal(%d,%v,1) error(%v)", mid, data, err)
		return
	}
	if err = s.dao.SetArcAppealCache(c, mid, aid); err != nil {
		log.Error("s.dao.SetArcAppealCache(%d,%d)", mid, aid)
		err = nil
	}
	return
}

// AuthorRecommend get author recommend data
func (s *Service) AuthorRecommend(c context.Context, aid int64) (res []*arcmdl.Arc, err error) {
	var (
		arcReply        *arcmdl.ArcReply
		aids            []int64
		recArcs, upArcs []*arcmdl.Arc
		arcs            *arcmdl.ArcsReply
		ip              = metadata.String(c, metadata.RemoteIP)
	)
	defer func() {
		if len(res) == 0 {
			res = _emptyArchive3
		}
	}()
	resAids := make(map[int64]int64)
	if arcReply, err = s.arcClient.Arc(c, &arcmdl.ArcRequest{Aid: aid}); err != nil {
		log.Error("s.arcClient.Arc(%d) error(%v)", aid, err)
		return
	}
	resAids[aid] = aid
	if recArcs, err = s.arc.Recommend3(c, &archive.ArgAid2{Aid: aid, RealIP: ip}); err != nil {
		log.Error("s.arc.Recommend3(%d) error(%v)", aid, err)
		err = nil
	} else {
		for _, v := range recArcs {
			res = append(res, v)
			resAids[v.Aid] = v.Aid
		}
	}
	if len(res) < s.c.Rule.AuthorRecCnt {
		if upArcs, err = s.arc.UpArcs3(c, &archive.ArgUpArcs2{Mid: arcReply.Arc.Author.Mid, Pn: _firstPn, Ps: s.c.Rule.AuthorRecCnt, RealIP: ip}); err != nil {
			log.Error("s.arc.UpArcs3(%d) error(%v)", arcReply.Arc.Author.Mid, err)
			err = nil
		} else {
			for _, v := range upArcs {
				if _, ok := resAids[v.Aid]; !ok {
					res = append(res, v)
					resAids[v.Aid] = v.Aid
					if len(res) >= s.c.Rule.AuthorRecCnt {
						return
					}
				}
			}
		}
	}
	if len(res) < s.c.Rule.AuthorRecCnt {
		if aids, err = s.dao.RelatedAids(c, aid); err != nil {
			log.Error("s.dao.RelatedArchives(%d) error(%v)", aid, err)
			err = nil
		} else if len(aids) > 0 {
			ps := s.c.Rule.AuthorRecCnt - len(res)
			if len(aids) > ps {
				aids = aids[0:ps]
			}
			archivesArgLog("AuthorRecommend", aids)
			if arcs, err = s.arcClient.Arcs(c, &arcmdl.ArcsRequest{Aids: aids}); err != nil {
				log.Error("s.arcClient.Arcs(%v) error(%v)", aids, err)
				err = nil
			} else {
				for _, aid := range aids {
					if _, ok := resAids[aid]; !ok {
						if arc, ok := arcs.Arcs[aid]; ok {
							res = append(res, model.FmtArc(arc))
							if len(res) >= s.c.Rule.AuthorRecCnt {
								return
							}
						}
					}
				}
			}
		}
	}
	return
}

// RelatedArcs get related archives
func (s *Service) RelatedArcs(c context.Context, aid int64) (res []*arcmdl.Arc, err error) {
	var (
		aids      []int64
		arcReply  *arcmdl.ArcReply
		arcsReply *arcmdl.ArcsReply
	)
	if _, ok := s.noRelAids[aid]; ok {
		res = _emptyArc
		return
	}
	if arcReply, err = s.arcClient.Arc(c, &arcmdl.ArcRequest{Aid: aid}); err != nil {
		log.Error("s.arcClient.Arc(%d) error(%v)", aid, err)
		return
	}
	// StateForbidUpDelete can have related arc
	if !arcReply.Arc.IsNormal() && arcReply.Arc.GetState() != archive.StateForbidUpDelete {
		res = _emptyArc
		return
	}
	if aids, err = s.dao.RelatedAids(c, aid); err != nil {
		log.Error("s.dao.RelatedArchives(%d) error(%v)", aid, err)
		err = nil
	} else if len(aids) > 0 {
		if len(aids) > s.c.Rule.RelatedArcCnt {
			aids = aids[:s.c.Rule.RelatedArcCnt]
		}
		archivesArgLog("RelatedArcs", aids)
		if arcsReply, err = s.arcClient.Arcs(c, &arcmdl.ArcsRequest{Aids: aids}); err != nil {
			log.Error("s.arcClient.Arcs(%v) error(%v)", aids, err)
			err = nil
		} else {
			for _, aid := range aids {
				if arc, ok := arcsReply.Arcs[aid]; ok && arc.IsNormal() {
					res = append(res, arc)
				}
			}
		}
	}
	if len(res) == 0 {
		res = _emptyArc
	}
	return
}

// Detail get merge view  card tag reply related.
func (s *Service) Detail(c context.Context, aid, mid int64, cdnIP, ck string) (rs *model.Detail, err error) {
	var (
		view                *model.View
		card                *model.Card
		tags                []*tagmdl.Tag
		reply               *model.ReplyHot
		related             []*arcmdl.Arc
		group               *errgroup.Group
		cardErr, relatedErr error
	)
	if view, err = s.View(c, aid, 0, mid, cdnIP, ck); err != nil {
		log.Error("s.View(%d) error %+v", aid, err)
		return
	}
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		if card, cardErr = s.Card(c, view.Author.Mid, mid, false, false); cardErr != nil {
			log.Error("s.Card(%d) error %+v", aid, cardErr)
		}
		return nil
	})
	group.Go(func() error {
		tags, _ = s.arcTags(errCtx, aid, mid)
		return nil
	})
	group.Go(func() error {
		reply, _ = s.replyHot(errCtx, aid)
		return nil
	})
	group.Go(func() error {
		if related, relatedErr = s.RelatedArcs(errCtx, aid); relatedErr != nil {
			log.Error("s.RelatedArcs(%d) error %+v", aid, relatedErr)
		}
		return nil
	})
	group.Wait()
	rs = &model.Detail{
		View:    view,
		Card:    card,
		Tags:    tags,
		Reply:   reply,
		Related: related,
	}
	return
}

// ArcUGCPay get arc ugc pay relation.
func (s *Service) ArcUGCPay(c context.Context, mid, aid int64) (data *model.AssetRelation, err error) {
	var relation *ugcmdl.AssetRelationResp
	data = new(model.AssetRelation)
	if arcReply, e := s.arcClient.Arc(c, &arcmdl.ArcRequest{Aid: aid}); e != nil {
		log.Error("s.arcClient.Arc(%d) error(%v)", aid, e)
	} else if arcReply.Arc.Author.Mid == mid {
		data.State = _ugcPaidState
		return
	}
	if relation, err = s.ugcPayClient.AssetRelation(c, &ugcmdl.AssetRelationReq{Mid: mid, Oid: aid, Otype: _ugcOtypeArc}); err != nil {
		log.Error("ArcUGCPay s.ugcPayClient.AssetRelation mid:%d aid:%d error(%v)", mid, aid, err)
		err = nil
		return
	}
	if relation.State == _ugcAssetPaid {
		data.State = _ugcPaidState
	}
	return
}

// ArcRelation .
func (s *Service) ArcRelation(c context.Context, mid, aid int64) (data *model.ReqUser, err error) {
	var arc *arcmdl.ArcReply
	data = new(model.ReqUser)
	if arc, err = s.arcClient.Arc(c, &arcmdl.ArcRequest{Aid: aid}); err != nil || arc.Arc == nil || !arc.Arc.IsNormal() {
		log.Error("ArcRelation s.arcClient.Arc(%d) error(%v)", aid, err)
		err = nil
		return
	}
	authorMid := arc.Arc.Author.Mid
	ip := metadata.String(c, metadata.RemoteIP)
	group, errCtx := errgroup.WithContext(c)
	// attention
	group.Go(func() error {
		if resp, e := s.accClient.Relation3(errCtx, &accmdl.RelationReq{Mid: mid, Owner: authorMid, RealIp: ip}); e != nil {
			log.Error("ArcRelation s.accClient.Relation3(%d,%d,%s) error(%v)", mid, authorMid, ip, e)
		} else if resp != nil {
			data.Attention = resp.Following
		}
		return nil
	})
	// favorite
	group.Go(func() error {
		if resp, e := s.fav.IsFav(errCtx, &favmdl.ArgIsFav{Type: favmdl.TypeVideo, Mid: mid, Oid: aid, RealIP: ip}); e != nil {
			log.Error("ArcRelation s.fav.IsFav(%d,%d,%s) error(%v)", mid, aid, ip, e)
		} else {
			data.Favorite = resp
		}
		return nil
	})
	// like
	group.Go(func() error {
		if resp, e := s.thumbup.HasLike(errCtx, &thumbup.ArgHasLike{Business: _businessLike, MessageIDs: []int64{aid}, Mid: mid, RealIP: ip}); e != nil {
			log.Error("ArcRelation s.thumbup.HasLike(%d,%d,%s) error %v", aid, mid, ip, e)
		} else if resp != nil {
			if v, ok := resp[aid]; ok {
				switch v {
				case thumbup.StateLike:
					data.Like = true
				case thumbup.StateDislike:
					data.Dislike = true
				}
			}
		}
		return nil
	})
	// coin
	group.Go(func() error {
		if resp, e := s.coinClient.ItemUserCoins(errCtx, &coinmdl.ItemUserCoinsReq{Mid: mid, Aid: aid, Business: model.CoinArcBusiness}); e != nil {
			log.Error("ArcRelation s.coinClient.ItemUserCoins(%d,%d,%s) error %v", mid, aid, ip, e)
		} else if resp != nil {
			data.Coin = resp.Number
		}
		return nil
	})
	group.Wait()
	return
}

func (s *Service) replyHot(c context.Context, aid int64) (res *model.ReplyHot, err error) {
	if res, err = s.dao.Hot(c, aid); err != nil {
		log.Error("s.dao.Hot(%d) error %+v", aid, err)
	}
	if res == nil {
		res = _emptyReplyHot
	}
	return
}

func (s *Service) arcTags(c context.Context, aid, mid int64) (res []*tagmdl.Tag, err error) {
	remoteIP := metadata.String(c, metadata.RemoteIP)
	var (
		arg = &tagmdl.ArgAid{
			Aid:    aid,
			Mid:    mid,
			RealIP: remoteIP,
		}
	)
	if res, err = s.tag.ArcTags(c, arg); err != nil {
		log.Error("s.tag.ArcTags(%v) error(%v)", arg, err)
	}
	if len(res) == 0 {
		res = _emptyTags
	}
	return
}

func (s *Service) ugcPayAsset(c context.Context, aid, mid int64) (data *ugcmdl.AssetQueryResp) {
	asset, err := s.ugcPayClient.AssetQuery(c, &ugcmdl.AssetQueryReq{Oid: aid, Otype: _ugcOtypeArc, Currency: _ugcCurrencyBp})
	if err != nil {
		log.Error("ugcPayAsset mid(%d) oid(%d) error(%v)", mid, aid, err)
		data = new(ugcmdl.AssetQueryResp)
		return
	}
	data = asset
	return
}
func (s *Service) loadManager() {
	for {
		time.Sleep(time.Duration(s.c.WEB.SpecailInterval))
		midsM, err := s.dao.Special(context.Background())
		if err != nil {
			log.Error("loadManager error(%+v)", err)
			continue
		}
		log.Info("load special mids(%+v)", midsM)
		s.specialMids = midsM
	}
}
