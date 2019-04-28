package service

import (
	"context"
	"errors"
	"hash/crc32"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/videoup/model/archive"
	mngmdl "go-common/app/admin/main/videoup/model/manager"
	accApi "go-common/app/service/main/account/api"
	upsrpc "go-common/app/service/main/up/api/v1"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"

	"go-common/library/xstr"
)

//ERROR
var (
	ErrRPCEmpty = errors.New("rpc reply empty")
)

func (s *Service) archiveRound(c context.Context, a *archive.Archive, aid, mid int64, typeID int16, nowRound, newState int8, cancelMission bool) (round int8) {
	//非特殊分区私单定时发布时，稿件进入四审
	//一个私单且定时发布的稿件 应该是先通过待审 传递到 私单四审 ，私单四审开放浏览后到达设置的定时时间最后再自动发布
	if s.isPorder(a) && newState == archive.StateForbidUserDelay && !s.isAuditType(typeID) {
		round = archive.RoundReviewFlow
		return
	}
	//定时发布的付费稿件 投上来就是round24  admin提交不需要再走付费审核

	//定时发布
	if newState == archive.StateForbidAdminDelay || newState == archive.StateForbidUserDelay {
		round = nowRound
		return
	}
	var addit *archive.Addit
	//活动稿件（非私单且非付费） | pgc 稿件 up_from(1,5,6) 直接round 99 简单流程业务
	if addit, _ = s.arc.Addit(c, aid); addit != nil && (addit.UpFrom == archive.UpFromPGC || addit.UpFrom == archive.UpFromSecretPGC || addit.UpFrom == archive.UpFromCoopera || (addit.MissionID > 0 && !s.isPorder(a) && !s.isUGCPay(a))) {
		if addit.UpFrom == archive.UpFromSecretPGC && newState == archive.StateForbidWait {
			//pgc 机密待审state=-1 不变更round
			round = nowRound
			return
		}
		if archive.NormalState(newState) {
			//pgc 生产组机密稿件(番剧、付费版权，严格要求时效和保密性) up_from=5  开放回查流程 90
			if addit.UpFrom == archive.UpFromSecretPGC && nowRound < archive.RoundTriggerClick {
				//三查
				round = archive.RoundTriggerClick
				return
			}
		}
		//pgc 生产组常规稿件(片包等其他内容) up_from=1  进三审
		if addit.UpFrom == archive.UpFromPGC && nowRound < archive.RoundAuditThird {
			round = archive.RoundAuditThird
			return
		}
		//pgc 合作方嵌套，终结
		if addit.UpFrom == archive.UpFromCoopera && nowRound < archive.RoundAuditThird {
			round = archive.RoundEnd
			return
		}
		round = archive.RoundEnd
		return
	}
	//其他稿件  （非定时，活动，pgc）
	if nowRound == archive.RoundAuditSecond {
		//二审阶段
		if newState == archive.StateForbidWait {
			//1、非特殊分区私单进入私单四审；2、私单活动稿件进入私单四审
			if (s.isPorder(a) && !s.isAuditType(typeID)) || (s.isPorder(a) && !cancelMission && addit.MissionID > 0) {
				//私单四审
				round = archive.RoundReviewFlow
			} else if s.isUGCPay(a) && (!s.isAuditType(typeID) || !s.isPorder(a) || (!cancelMission && addit.MissionID > 0)) {
				//付费稿件 非特殊分区 非私单
				round = archive.RoundAuditUGCPayFlow
			} else if s.isAuditType(typeID) || (addit != nil && addit.OrderID > 0) {
				//特殊分区 ,商单到三审 二审到三审
				round = archive.RoundAuditThird
			} else {
				//不变
				round = archive.RoundAuditSecond
			}
		} else if archive.NormalState(newState) {
			//回查控制
			round = s.normalRound(c, aid, mid, typeID, nowRound, newState)
		} else if newState == archive.StateForbidFixed {
			//1、非特殊分区私单进入私单四审修复待审；2、私单+活动稿件进入私单四审修复待审
			if (s.isPorder(a) && !s.isAuditType(typeID)) || (s.isPorder(a) && !cancelMission && addit.MissionID > 0) {
				round = archive.RoundReviewFlow
			} else if s.isAuditType(typeID) && !(!cancelMission && addit.MissionID > 0) {
				//特殊分区 二审到三审
				round = archive.RoundAuditThird
			} else if s.isUGCPay(a) && (!s.isAuditType(typeID) || !s.isPorder(a) || (!cancelMission && addit.MissionID > 0)) {
				//付费稿件 非特殊分区 非私单
				round = archive.RoundAuditUGCPayFlow
			} else {
				//不变
				round = archive.RoundAuditSecond
			}
		} else {
			round = archive.RoundEnd
		}
	} else if nowRound == archive.RoundReviewFlow {
		//私单四审21
		if archive.NormalState(newState) {
			if s.isAuditType(typeID) || addit.MissionID > 0 {
				round = archive.RoundEnd
			} else {
				round = s.normalRound(c, aid, mid, typeID, nowRound, newState)
			}
		} else {
			round = archive.RoundEnd
		}
	} else if nowRound == archive.RoundAuditUGCPayFlow {
		//付费待审 24
		if archive.NormalState(newState) {
			if s.isAuditType(typeID) || addit.MissionID > 0 {
				round = archive.RoundEnd
			} else {
				round = s.normalRound(c, aid, mid, typeID, nowRound, newState)
			}
		} else {
			round = archive.RoundEnd
		}
	} else if nowRound == archive.RoundReviewFirst {
		//一查
		if archive.NormalState(newState) {
			round = archive.RoundReviewFirstWaitTrigger
		} else {
			round = archive.RoundEnd
		}
	} else if nowRound == archive.RoundAuditThird {
		//三审
		if s.isPorder(a) {
			//私单四审
			round = archive.RoundReviewFlow
		} else if s.isUGCPay(a) {
			round = archive.RoundAuditUGCPayFlow
		} else {
			round = archive.RoundEnd
		}
	} else if nowRound == archive.RoundReviewSecond || nowRound == archive.RoundTriggerClick {
		//二查 三查
		round = archive.RoundEnd
	} else {
		round = nowRound
	}
	return
}

//normalRound 回查逻辑
func (s *Service) normalRound(c context.Context, aid, mid int64, typeID int16, nowRound, newState int8) (round int8) {
	if s.isWhite(mid) || s.isBlack(mid) {
		round = archive.RoundReviewSecond
	} else if plf, _ := s.profile(c, mid); plf != nil && plf.Follower >= s.fansCache {
		//社区回查 粉丝阈值
		round = archive.RoundReviewSecond
	} else if plf != nil && plf.Follower < s.fansCache && s.isRoundType(typeID) {
		//回查分区
		round = archive.RoundReviewFirst // NOTE: if audit type, state must not open!!! so cannot execute here...
	} else {
		//点击量回查
		round = archive.RoundReviewFirstWaitTrigger
		if plf == nil {
			log.Info("archive(%d) card(%d) is nil", aid, mid)
		} else {
			log.Info("archive(%d) card(%d) fans(%d) little than config(%d)", aid, mid, plf.Follower, s.fansCache)
		}
	}
	return
}

func (s *Service) hadPassed(c context.Context, aid int64) (had bool) {
	id, err := s.arc.GetFirstPassByAID(c, aid)
	if err != nil {
		log.Error("hadPassed s.arc.GetFirstPassByAID error(%v) aid(%d)", err, aid)
		return
	}

	had = id > 0
	return
}

func (s *Service) isPorder(a *archive.Archive) bool {
	if a == nil {
		return false
	}
	return a.AttrVal(archive.AttrBitIsPorder) == archive.AttrYes
}

//ugc pay only
func (s *Service) isUGCPay(a *archive.Archive) bool {
	if a == nil {
		return false
	}
	return a.AttrVal(archive.AttrBitUGCPay) == archive.AttrYes
}

func (s *Service) profile(c context.Context, mid int64) (card *accApi.ProfileStatReply, err error) {
	if card, err = s.accRPC.ProfileWithStat3(c, &accApi.MidReq{Mid: mid}); err != nil {
		log.Error("s.accRPC.ProfileWithStat3(%d) error(%v)", mid, err)
	}
	return
}

func (s *Service) upCards(c context.Context, mids []int64) (p map[int64]*accApi.Card, err error) {
	p = make(map[int64]*accApi.Card)
	if len(mids) == 0 {
		return
	}
	var res *accApi.CardsReply
	if res, err = s.accRPC.Cards3(c, &accApi.MidsReq{Mids: mids}); err != nil {
		p = nil
		log.Error("s.accRPC.Cards3(%v) error(%v)", mids, err)
		return
	}
	p = res.Cards
	return
}

func (s *Service) archivePtime(c context.Context, aid int64, newState int8, newPtime xtime.Time) (ptime xtime.Time) {
	if newState >= archive.StateOpen {
		if !s.hadPassed(c, aid) {
			ptime = xtime.Time(time.Now().Unix())
			return
		}
	}
	ptime = newPtime
	return
}

func (s *Service) archiveAttr(c context.Context, ap *archive.ArcParam, isExt bool) (attrs map[uint]int32, forbidAttrs map[string]map[uint]int32) {
	//批量和单个提交都需要补全对应属性值

	attrs = make(map[uint]int32)
	attrs[archive.AttrBitNoRank] = ap.Attrs.NoRank
	attrs[archive.AttrBitNoDynamic] = ap.Attrs.NoDynamic
	attrs[archive.AttrBitNoWeb] = ap.Attrs.NoWeb
	attrs[archive.AttrBitNoMobile] = ap.Attrs.NoMobile
	attrs[archive.AttrBitNoSearch] = ap.Attrs.NoSearch
	attrs[archive.AttrBitOverseaLock] = ap.Attrs.OverseaLock
	attrs[archive.AttrBitNoRecommend] = ap.Attrs.NoRecommend
	attrs[archive.AttrBitNoReprint] = ap.Attrs.NoReprint
	attrs[archive.AttrBitHasHD5] = ap.Attrs.HasHD5
	attrs[archive.AttrBitAllowBp] = ap.Attrs.AllowBp
	attrs[archive.AttrBitIsPorder] = ap.Attrs.IsPorder
	attrs[archive.AttrBitLimitArea] = ap.Attrs.LimitArea
	attrs[archive.AttrBitPushBlog] = ap.Attrs.PushBlog
	attrs[archive.AttrBitUGCPay] = ap.Attrs.UGCPay
	attrs[archive.AttrBitParentMode] = ap.Attrs.ParentMode
	// pgc
	attrs[archive.AttrBitIsMovie] = ap.Attrs.IsMovie
	attrs[archive.AttrBitBadgepay] = ap.Attrs.BadgePay
	attrs[archive.AttrBitIsBangumi] = ap.Attrs.IsBangumi
	attrs[archive.AttrBitIsPGC] = ap.Attrs.IsPGC
	if isExt {
		attrs[archive.AttrBitAllowTag] = ap.Attrs.AllowTag
		attrs[archive.AttrBitJumpURL] = ap.Attrs.JumpURL
	}
	forbidAttrs = make(map[string]map[uint]int32, 4)
	forbidAttrs[archive.ForbidRank] = map[uint]int32{
		archive.ForbidRankMain:      ap.Forbid.Rank.Main,
		archive.ForbidRankRecentArc: ap.Forbid.Rank.RecentArc,
		archive.ForbidRankAllArc:    ap.Forbid.Rank.AllArc,
	}
	forbidAttrs[archive.ForbidDynamic] = map[uint]int32{
		archive.ForbidDynamicMain: ap.Forbid.Dynamic.Main,
	}
	forbidAttrs[archive.ForbidRecommend] = map[uint]int32{
		archive.ForbidRecommendMain: ap.Forbid.Recommend.Main,
	}
	forbidAttrs[archive.ForbidShow] = map[uint]int32{
		archive.ForbidShowMain:    ap.Forbid.Show.Main,
		archive.ForbidShowMobile:  ap.Forbid.Show.Mobile,
		archive.ForbidShowWeb:     ap.Forbid.Show.Web,
		archive.ForbidShowOversea: ap.Forbid.Show.Oversea,
		archive.ForbidShowOnline:  ap.Forbid.Show.Online,
	}
	return
}

func (s *Service) isAccess(c context.Context, aid int64) (wm bool) {
	var vs []*archive.Video
	if vs, _ = s.arc.NewVideosByAid(c, aid); len(vs) <= 0 {
		return
	}
	for _, v := range vs {
		if v.Status == 10000 {
			wm = true
			return
		}
	}
	return
}

// CheckArchive check typeid
func (s *Service) CheckArchive(aps []*archive.ArcParam) bool {
	for _, ap := range aps {
		if ap.Aid == 0 || ap.UID == 0 {
			return false
		}
	}
	return true
}

// CheckVideo check video
func (s *Service) CheckVideo(vps []*archive.VideoParam) bool {
	for _, vp := range vps {
		if vp.ID == 0 || vp.Aid == 0 || vp.Filename == "" || vp.Cid == 0 || vp.UID == 0 {
			return false
		}
	}
	return true
}

// CheckStaff check
func (s *Service) CheckStaff(vps []*archive.StaffParam) bool {
	//允许为空 不允许为Nil
	if len(vps) == 0 {
		return true
	}
	for _, vp := range vps {
		if vp.MID == 0 || vp.Title == "" {
			return false
		}
	}
	return true
}

// TypeTopParent get archive type's first level type
func (s *Service) TypeTopParent(tid int16) (tp *archive.Type, err error) {
	if _, ok := s.typeCache[tid]; !ok {
		err = errors.New("archive type id not exist. id:" + strconv.Itoa(int(tid)))
		return
	}
	if s.typeCache[tid].PID == 0 {
		tp = s.typeCache[tid]
	} else {
		tp, err = s.TypeTopParent(s.typeCache[tid].PID)
		if err != nil {
			return
		}
	}
	return
}

//StringHandler handle two strings 以s1顺序为准
func StringHandler(s1 string, s2 string, delimiter string, subtraction bool) string {
	if strings.TrimSpace(s2) == "" {
		return s1
	}

	var (
		res       []string
		duplicate []int
	)
	s1Arr := strings.Split(s1, delimiter)
	s2Arr := strings.Split(s2, delimiter)
	for _, s1Item := range s1Arr {
		dupIndex := -1
		for k, s2Item := range s2Arr {
			if s1Item == s2Item {
				dupIndex = k
				break
			}
		}

		if dupIndex >= 0 && subtraction {
			continue
		}
		if dupIndex >= 0 {
			duplicate = append(duplicate, dupIndex)
		}
		res = append(res, s1Item)
	}

	if !subtraction {
		for k, s2Item := range s2Arr {
			s2Item = strings.TrimSpace(s2Item)
			if s2Item == "" {
				continue
			}
			add := true
			for _, dup := range duplicate {
				if k == dup {
					add = false
					break
				}
			}
			if add {
				res = append(res, s2Item)
			}
		}
	}
	return strings.Join(res, delimiter)
}

// SplitInts 去掉id字符串中的空白字符
func (s *Service) SplitInts(str string) ([]int64, error) {
	empties := []string{" ", "\n", "\t", "\r"}
	for _, v := range empties {
		str = strings.Replace(str, v, "", -1)
	}
	str = strings.Trim(str, ",")
	return xstr.SplitInts(str)
}

// coverURL convert cover url to full url.
func coverURL(uri string) (cover string) {
	if uri == "" {
		cover = "http://static.hdslb.com/images/transparent.gif"
		return
	}
	cover = uri
	if strings.Index(uri, "http://") == 0 {
		return
	}
	if len(uri) >= 10 && uri[:10] == "/templets/" {
		return
	}
	if strings.HasPrefix(uri, "group1") {
		cover = "http://i0.hdslb.com/" + uri
		return
	}
	if pos := strings.Index(uri, "/uploads/"); pos != -1 && (pos == 0 || pos == 3) {
		cover = uri[pos+8:]
	}
	cover = strings.Replace(cover, "{IMG}", "", -1)
	cover = "http://i" + strconv.FormatInt(int64(crc32.ChecksumIEEE([]byte(cover)))%3, 10) + ".hdslb.com" + cover
	return
}

func (s *Service) upGroupMids(c context.Context, gid int64) (mids []int64, err error) {
	var (
		total int
		maxps = 10000
		req   = &upsrpc.UpGroupMidsReq{
			Pn:      1,
			GroupID: gid,
			Ps:      maxps,
		}
		reply *upsrpc.UpGroupMidsReply
	)

	for {
		reply, err = s.upsRPC.UpGroupMids(c, req)
		if err == nil && (reply == nil || reply.Mids == nil) {
			err = ErrRPCEmpty
		}
		if err != nil {
			log.Error("UpGroupMids req(%+v) error(%v)", req, err)
			return
		}
		total = reply.Total
		mids = append(mids, reply.Mids...)
		if reply.Size() != maxps {
			break
		}
		req.Pn++
	}
	log.Info("upGroupMids(%d) reply total(%d) len(%d)", gid, total, len(mids))
	return
}

func (s *Service) upSpecial(c context.Context) (ups map[int8]map[int64]struct{}, err error) {
	var (
		g                                                                                                    errgroup.Group
		whitegroup, blackgroup, pgcgroup, ugcxgroup, policygroup, dangergroup, twoforbidgroup, pgcwhitegroup map[int64]struct{}
	)
	ups = make(map[int8]map[int64]struct{})

	f := func(gid int8) (map[int64]struct{}, error) {
		group := make(map[int64]struct{})
		mids, e := s.upGroupMids(c, int64(gid))
		if e != nil {
			return group, e
		}
		for _, mid := range mids {
			group[mid] = struct{}{}
		}
		return group, nil
	}
	g.Go(func() error {
		whitegroup, err = f(mngmdl.UpperTypeWhite)
		return err
	})
	g.Go(func() error {
		blackgroup, err = f(mngmdl.UpperTypeBlack)
		return err
	})
	g.Go(func() error {
		pgcgroup, err = f(mngmdl.UpperTypePGC)
		return err
	})
	g.Go(func() error {
		ugcxgroup, err = f(mngmdl.UpperTypeUGCX)
		return err
	})
	g.Go(func() error {
		policygroup, err = f(mngmdl.UpperTypePolity)
		return err
	})
	g.Go(func() error {
		dangergroup, err = f(mngmdl.UpperTypeDanger)
		return err
	})
	g.Go(func() error {
		twoforbidgroup, err = f(mngmdl.UpperTypeTwoForbid)
		return err
	})
	g.Go(func() error {
		pgcwhitegroup, err = f(mngmdl.UpperTypePGCWhite)
		return err
	})
	if err = g.Wait(); err != nil {
		return
	}

	ups[mngmdl.UpperTypeWhite] = whitegroup
	ups[mngmdl.UpperTypeBlack] = blackgroup
	ups[mngmdl.UpperTypePGC] = pgcgroup
	ups[mngmdl.UpperTypeUGCX] = ugcxgroup
	ups[mngmdl.UpperTypePolity] = policygroup
	ups[mngmdl.UpperTypeDanger] = dangergroup
	ups[mngmdl.UpperTypeTwoForbid] = twoforbidgroup
	ups[mngmdl.UpperTypePGCWhite] = pgcwhitegroup
	return
}
