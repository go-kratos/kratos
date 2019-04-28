package like

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net"
	"strconv"
	"time"

	"go-common/app/interface/main/activity/conf"
	"go-common/app/interface/main/activity/dao/bnj"
	"go-common/app/interface/main/activity/dao/like"
	bnjmdl "go-common/app/interface/main/activity/model/bnj"
	l "go-common/app/interface/main/activity/model/like"
	tagrpc "go-common/app/interface/main/tag/rpc/client"
	accapi "go-common/app/service/main/account/api"
	arccli "go-common/app/service/main/archive/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	arcmdl "go-common/app/service/main/archive/model/archive"
	coinrpc "go-common/app/service/main/coin/api/gorpc"
	spymdl "go-common/app/service/main/spy/model"
	spyrpc "go-common/app/service/main/spy/rpc/client"
	thumbup "go-common/app/service/main/thumbup/rpc/client"
	suitrpc "go-common/app/service/main/usersuit/rpc/client"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
	"go-common/library/sync/pipeline/fanout"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

const (
	_yes           = 1
	_no            = 0
	_typeAll       = "all"
	_typeRegion    = "region"
	_like          = "like"
	_grade         = "grade"
	_vote          = "vote"
	_silenceForbid = 1
)

// Service struct
type Service struct {
	c              *conf.Config
	dao            *like.Dao
	bnjDao         *bnj.Dao
	arcRPC         *arcrpc.Service2
	arcClient      arccli.ArchiveClient
	coin           *coinrpc.Service
	suit           *suitrpc.Service2
	accClient      accapi.AccountClient
	spy            *spyrpc.Service
	tagRPC         *tagrpc.Service
	thumbup        thumbup.ThumbupRPC
	cache          *fanout.Fanout
	arcType        map[int16]*arcmdl.ArcType
	dialectTags    map[int64]struct{}
	dialectRegions map[int16]struct{}
	reward         map[int]*bnjmdl.Reward
	r              *rand.Rand
	newestSubTs    int64
}

// New Service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:       c,
		cache:   fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		dao:     like.New(c),
		bnjDao:  bnj.New(c),
		arcRPC:  arcrpc.New2(c.RPCClient2.Archive),
		coin:    coinrpc.New(c.RPCClient2.Coin),
		suit:    suitrpc.New(c.RPCClient2.Suit),
		spy:     spyrpc.New(c.RPCClient2.Spy),
		tagRPC:  tagrpc.New2(c.RPCClient2.Tag),
		thumbup: thumbup.New(c.RPCClient2.Thumbup),
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	var err error
	if s.arcClient, err = arccli.NewClient(c.ArcClient); err != nil {
		panic(err)
	}
	if s.accClient, err = accapi.NewClient(c.AccClient); err != nil {
		panic(err)
	}
	s.initDialect()
	s.initReward()
	go s.arcTypeproc()
	go s.actSourceproc()
	go s.newestSubTsproc()
	return
}

func checkIsLike(subType int64) (likesType bool) {
	switch subType {
	case l.PICTURELIKE, l.DRAWYOOLIKE, l.TEXTLIKE, l.VIDEOLIKE, l.VIDEO2, l.VIDEO, l.SMALLVIDEO, l.MUSIC, l.PHONEVIDEO, l.STORYKING:
		likesType = true
	default:
		likesType = false
	}
	return
}

// Subject service
func (s *Service) Subject(c context.Context, sid int64) (res *l.Subject, err error) {
	var (
		mc      = true
		subErr  error
		likeErr error
	)
	if res, err = s.dao.InfoCache(c, sid); err != nil {
		err = nil
		mc = false
	} else if res != nil {
		if res, err = s.LikeArc(c, res); err != nil {
			return
		}
	}
	eg, errCtx := errgroup.WithContext(c)
	var ls = make([]*l.Like, 0)
	eg.Go(func() error {
		res, subErr = s.dao.Subject(errCtx, sid)
		return subErr
	})
	eg.Go(func() error {
		ls, likeErr = s.dao.LikeTypeList(errCtx, sid)
		return likeErr
	})
	if err = eg.Wait(); err != nil {
		log.Error("eg.Wait error(%v)", err)
		return
	}
	if res != nil {
		res.List = ls
	}
	if mc {
		err = s.dao.SetInfoCache(c, res, sid)
		if err != nil {
			log.Error("SetInfoCache error(%v)", err)
		}
	}
	if res, err = s.LikeArc(c, res); err != nil {
		return
	}
	return
}

// LikeArc service
func (s *Service) LikeArc(c context.Context, sub *l.Subject) (res *l.Subject, err error) {
	if sub != nil {
		if sub.ID == 0 {
			res = nil
		} else {
			res = sub
			var (
				ok   bool
				arcs map[int64]*arccli.Arc
				aids []int64
			)
			for _, l := range res.List {
				aids = append(aids, l.Wid)
			}
			argAids := &arcmdl.ArgAids2{
				Aids: aids,
			}
			if arcs, err = s.arcRPC.Archives3(c, argAids); err != nil {
				log.Error("s.arcRPC.Archives(arcAids:(%v), arcs), err(%v)", aids, err)
				return
			}
			for _, l := range res.List {
				if l.Archive, ok = arcs[l.Wid]; !ok {
					log.Info("s.arcs.wid:(%d),ok(%v)", l.Wid, ok)
					continue
				}
			}
		}
	}
	return
}

// OnlineVote Service
func (s *Service) OnlineVote(c context.Context, mid, vote, stage, aid int64) (res bool, err error) {
	res = true
	if vote != _yes && vote != _no {
		err = nil
		res = false
		return
	}
	var incrKey string
	midStr := strconv.FormatInt(mid, 10)
	aidStr := strconv.FormatInt(aid, 10)
	stageStr := strconv.FormatInt(stage, 10)
	midKye := midStr + ":" + aidStr + ":" + stageStr
	if res, err = s.dao.RsSetNX(c, midKye); err != nil {
		log.Error("s.OnlineVote.reids(mid:(%v),,vote:(%v),stage:(%v)), err(%v)", mid, vote, stage, err)
		return
	}
	if !res {
		return
	}
	if vote == _yes {
		incrKey = aidStr + ":" + stageStr + ":yes"
	} else {
		incrKey = aidStr + ":" + stageStr + ":no"
	}
	if mid == 288239 || mid == 26366366 || mid == 20453897 {
		log.Info("288239,26366366,20453897")
		if res, err = s.dao.Incrby(c, incrKey); err != nil {
			log.Error("s.OnlineVote.Incrby(key:(%v)", incrKey)
			return
		}
	} else {
		if res, err = s.dao.Incr(c, incrKey); err != nil {
			log.Error("s.OnlineVote.Incr(key:(%v)", incrKey)
			return
		}
	}
	s.dao.CVoteLog(c, 0, aid, mid, stage, vote)
	return
}

// Ltime service
func (s *Service) Ltime(c context.Context, sid int64) (res map[string]interface{}, err error) {
	var key = "ltime:" + strconv.FormatInt(sid, 10)
	var b []byte
	if b, err = s.dao.Rb(c, key); err != nil {
		log.Error("s.dao.Rb((%v), err(%v)", key, err)
		return
	}
	if b == nil {
		res = nil
		return
	}
	if err = json.Unmarshal(b, &res); err != nil {
		log.Error("s.Ltime.Unmarshal((%v), err(%v)", b, err)
		return
	}
	if res["time"] != nil {
		if st, ok := res["time"].(float64); ok {
			res["currentTime"] = time.Now().Unix() - int64(st)
		}
	}
	return
}

// LikeAct service
func (s *Service) LikeAct(c context.Context, p *l.ParamAddLikeAct, mid int64) (res int64, err error) {
	var (
		subject   *l.SubjectItem
		likeItem  *l.Item
		memberRly *accapi.ProfileReply
		subErr    error
		likeErr   error
	)
	eg, errCtx := errgroup.WithContext(c)
	eg.Go(func() error {
		subject, subErr = s.dao.ActSubject(errCtx, p.Sid)
		return subErr
	})
	eg.Go(func() error {
		likeItem, likeErr = s.dao.Like(errCtx, p.Lid)
		return likeErr
	})
	if err = eg.Wait(); err != nil {
		log.Error("eg.Wait error(%v)", err)
		return
	}
	if subject.ID == 0 || subject.Type == l.STORYKING {
		err = ecode.ActivityHasOffLine
		return
	}
	if likeItem.ID == 0 || likeItem.Sid != p.Sid {
		err = ecode.ActivityLikeHasOffLine
		return
	}
	if memberRly, err = s.accClient.Profile3(c, &accapi.MidReq{Mid: mid}); err != nil {
		log.Error(" s.acc.Profile3(c,&accmdl.ArgMid{Mid:%d}) error(%v)", mid, err)
		return
	}
	if err = s.judgeUser(c, subject, memberRly.Profile); err != nil {
		return
	}
	nowTime := time.Now().Unix()
	if int64(memberRly.Profile.JoinTime) >= (nowTime - 86400*7) {
		err = ecode.ActivityLikeMemberLimit
		return
	}
	if int64(subject.Lstime) >= nowTime {
		err = ecode.ActivityLikeNotStart
		return
	}
	if int64(subject.Letime) <= nowTime {
		err = ecode.ActivityLikeHasEnd
		return
	}
	var (
		likeAct map[int64]int
		lids    = []int64{p.Lid}
	)
	if likeAct, err = s.dao.LikeActs(c, p.Sid, mid, lids); err != nil {
		log.Error("s.dao.LikeActMidList(%v) error(%+v)", p, err)
		return
	}
	if _, ok := likeAct[p.Lid]; !ok {
		log.Error("s.dao.LikeActMidList() get lid value error()")
		return
	}
	isLikeType := s.isLikeType(c, subject.Type)
	if likeAct[p.Lid] == like.HasLike {
		if isLikeType == _like {
			err = ecode.ActivityLikeHasLike
		} else if isLikeType == _vote {
			err = ecode.ActivityLikeHasVote
		} else {
			err = ecode.ActivityLikeHasGrade
		}
		return
	}
	var score int64
	if isLikeType == _like || isLikeType == _vote {
		score = l.LIKESCORE
	} else {
		score = p.Score
	}
	if err = s.dao.SetRedisCache(c, p.Sid, p.Lid, score, likeItem.Type); err != nil {
		log.Error("s.dao.SetRedisCache(%v) error(%+v)", p, err)
		return
	}
	likeActAdd := &l.Action{
		Lid:    p.Lid,
		Mid:    mid,
		Sid:    p.Sid,
		Action: score,
		IPv6:   make([]byte, 0),
	}
	if IPv6 := net.ParseIP(metadata.String(c, metadata.RemoteIP)); IPv6 != nil {
		likeActAdd.IPv6 = IPv6
	}
	if res, err = s.dao.LikeActAdd(c, likeActAdd); err != nil {
		log.Error("s.dao.LikeActAdd(%v) error(%+v)", p, err)
		return
	}
	s.dao.AddCacheLikeActs(c, p.Sid, mid, map[int64]int{p.Lid: like.HasLike})
	return
}

// StoryKingAct .
func (s *Service) StoryKingAct(c context.Context, p *l.ParamStoryKingAct, mid int64) (res map[string]int64, err error) {
	var (
		subject   *l.SubjectItem
		likeItem  *l.Item
		memberRly *accapi.ProfileReply
		subErr    error
		likeErr   error
		leftTime  int64
	)
	eg, errCtx := errgroup.WithContext(c)
	eg.Go(func() error {
		subject, subErr = s.dao.ActSubject(errCtx, p.Sid)
		return subErr
	})
	eg.Go(func() error {
		likeItem, likeErr = s.dao.Like(errCtx, p.Lid)
		return likeErr
	})
	if err = eg.Wait(); err != nil {
		log.Error("eg.Wait error(%v)", err)
		return
	}
	if subject.ID == 0 || subject.Type != l.STORYKING {
		err = ecode.ActivityHasOffLine
		return
	}
	if likeItem.ID == 0 || likeItem.Sid != p.Sid {
		err = ecode.ActivityLikeHasOffLine
		return
	}
	if memberRly, err = s.accClient.Profile3(c, &accapi.MidReq{Mid: mid}); err != nil {
		log.Error(" s.acc.Profile3(c,&accmdl.ArgMid{Mid:%d}) error(%v)", mid, err)
		return
	}
	if err = s.judgeUser(c, subject, memberRly.Profile); err != nil {
		return
	}
	nowTime := time.Now().Unix()
	if int64(subject.Lstime) >= nowTime {
		err = ecode.ActivityLikeNotStart
		return
	}
	if int64(subject.Letime) <= nowTime {
		err = ecode.ActivityLikeHasEnd
		return
	}
	if leftTime, err = s.storyLikeCheck(c, p.Sid, p.Lid, mid, subject.DailyLikeLimit, subject.DailySingleLikeLimit); err != nil {
		log.Error(" s.storyLikeCheck(%d,%d,%d) error(%+v)", p.Sid, p.Lid, mid, err)
		return
	}
	if leftTime < p.Score {
		if leftTime > 0 {
			p.Score = leftTime
		} else {
			err = ecode.ActivityOverDailyScore
			return
		}
	}
	if err = s.dao.SetRedisCache(c, p.Sid, p.Lid, p.Score, likeItem.Type); err != nil {
		log.Error("s.dao.SetRedisCache(%v) error(%+v)", p, err)
		return
	}
	likeActAdd := &l.Action{
		Lid:    p.Lid,
		Mid:    mid,
		Sid:    p.Sid,
		Action: int64(p.Score),
		IPv6:   make([]byte, 0),
	}
	if IPv6 := net.ParseIP(metadata.String(c, metadata.RemoteIP)); IPv6 != nil {
		likeActAdd.IPv6 = IPv6
	}
	res = make(map[string]int64, 2)
	if res["act_id"], err = s.dao.LikeActAdd(c, likeActAdd); err != nil {
		log.Error("s.dao.LikeActAdd(%v) error(%+v)", p, err)
		return
	}
	s.storyLikeActSet(c, p.Sid, p.Lid, mid, p.Score)
	res["score"] = p.Score
	return
}

// StoryKingLeftTime .
func (s *Service) StoryKingLeftTime(c context.Context, sid, mid int64) (res int64, err error) {
	var (
		subject   *l.SubjectItem
		memberRly *accapi.ProfileReply
	)
	if subject, err = s.dao.ActSubject(c, sid); err != nil {
		return
	}
	if subject.ID == 0 || subject.Type != l.STORYKING {
		err = ecode.ActivityHasOffLine
		return
	}
	if memberRly, err = s.accClient.Profile3(c, &accapi.MidReq{Mid: mid}); err != nil {
		log.Error(" s.acc.Profile3(c,&accmdl.ArgMid{Mid:%d}) error(%v)", mid, err)
		return
	}
	if err = s.simpleJudge(c, subject, memberRly.Profile); err != nil {
		err = nil
		res = 0
		return
	}
	nowTime := time.Now().Unix()
	if int64(subject.Lstime) >= nowTime || int64(subject.Letime) <= nowTime {
		res = 0
		return
	}
	if res, err = s.storySumUsed(c, sid, mid); err != nil {
		log.Error("s.storySumUsed(%d,%d) error(%+v)", sid, mid, err)
		return
	}
	res = subject.DailyLikeLimit - res
	if res < 0 {
		res = 0
	}
	return
}

// UpList .
func (s *Service) UpList(c context.Context, p *l.ParamList, mid int64) (res *l.ListInfo, err error) {
	switch p.Type {
	case like.EsOrderLikes, like.EsOrderCoin, like.EsOrderReply, like.EsOrderShare, like.EsOrderClick, like.EsOrderDm, like.EsOrderFav:
		res, err = s.EsList(c, p, mid)
	case like.ActOrderCtime, like.ActOrderLike, like.ActOrderRandom:
		res, err = s.StoryKingList(c, p, mid)

	default:
		err = errors.New("type error")
	}
	return
}

// EsList .
func (s *Service) EsList(c context.Context, p *l.ParamList, mid int64) (res *l.ListInfo, err error) {
	var (
		subject *l.SubjectItem
	)
	if subject, err = s.dao.ActSubject(c, p.Sid); err != nil {
		return
	}
	if subject.ID == 0 {
		err = ecode.ActivityHasOffLine
		return
	}
	if res, err = s.dao.ListFromES(c, p.Sid, p.Type, p.Ps, p.Pn, 0); err != nil {
		log.Error("s.dao.ListFromES(%d) error(%+v)", p.Sid, err)
		return
	}
	if res == nil || len(res.List) == 0 {
		return
	}
	if err = s.getContent(c, res.List, subject.Type, mid, p.Type); err != nil {
		log.Error("s.getContent(%d) error(%v)", p.Sid, err)
	}
	return
}

// StoryKingList .
func (s *Service) StoryKingList(c context.Context, p *l.ParamList, mid int64) (res *l.ListInfo, err error) {
	var (
		subject  *l.SubjectItem
		likeList []*l.List
		total    int64
	)
	if subject, err = s.dao.ActSubject(c, p.Sid); err != nil {
		return
	}
	if subject.ID == 0 {
		err = ecode.ActivityHasOffLine
		return
	}
	likesType := checkIsLike(subject.Type)
	switch p.Type {
	case like.ActOrderCtime:
		likeList, err = s.orderByCtime(c, p.Sid, p.Pn, p.Ps, likesType)
	case like.ActOrderRandom:
		likeList, err = s.orderByRandom(c, p.Sid, p.Pn, p.Ps, likesType)
	default:
		likeList, err = s.orderByLike(c, p.Sid, p.Pn, p.Ps)
	}
	if err != nil {
		log.Error("s.orderBy(%s)(%d) error(%v)", p.Type, p.Sid, err)
		return
	}
	if len(likeList) == 0 {
		return
	}
	if err = s.getContent(c, likeList, subject.Type, mid, p.Type); err != nil {
		log.Error("s.getContent(%d) error(%v)", p.Sid, err)
	}
	if p.Type == like.ActOrderRandom {
		total, _ = s.dao.LikeRandomCount(c, p.Sid)
	} else {
		total, _ = s.dao.LikeCount(c, p.Sid)
	}
	res = &l.ListInfo{List: likeList, Page: &l.Page{Size: p.Ps, Num: p.Pn, Total: total}}
	return
}

// getContent get likes extends 后期接入其他活动补充完善.
func (s *Service) getContent(c context.Context, list []*l.List, subType, mid int64, order string) (err error) {
	switch subType {
	case l.STORYKING:
		err = s.actContent(c, list, mid)
	case l.PICTURE, l.PICTURELIKE, l.DRAWYOO, l.DRAWYOOLIKE, l.TEXT, l.TEXTLIKE, l.QUESTION:
		err = s.contentAccount(c, list)
	case l.VIDEOLIKE, l.VIDEO, l.VIDEO2, l.SMALLVIDEO, l.PHONEVIDEO, l.ONLINEVOTE:
		err = s.arcTag(c, list, order, mid)
	default:
		err = ecode.RequestErr
	}
	return
}

// actContent get like_content and account info.
func (s *Service) contentAccount(c context.Context, list []*l.List) (err error) {
	var (
		lt     = len(list)
		lids   = make([]int64, 0, lt)
		mids   = make([]int64, 0, lt)
		cont   map[int64]*l.LikeContent
		actRly *accapi.InfosReply
	)
	for _, v := range list {
		lids = append(lids, v.ID)
		if v.Mid > 0 {
			mids = append(mids, v.Mid)
		}
	}
	eg, errCtx := errgroup.WithContext(c)
	eg.Go(func() (e error) {
		cont, e = s.dao.LikeContent(errCtx, lids)
		return
	})
	eg.Go(func() (e error) {
		actRly, e = s.accClient.Infos3(errCtx, &accapi.MidsReq{Mids: mids})
		return
	})
	if err = eg.Wait(); err != nil {
		log.Error("actContent:eg.Wait() error(%v)", err)
		return
	}
	for _, v := range list {
		obj := make(map[string]interface{}, 2)
		if _, ok := cont[v.ID]; ok {
			obj["cont"] = cont[v.ID]
		}
		if _, k := actRly.Infos[v.Mid]; k {
			obj["act"] = actRly.Infos[v.Mid]
		}
		v.Object = obj
	}
	return
}

// actContent get like_content and account info.
func (s *Service) actContent(c context.Context, list []*l.List, mid int64) (err error) {
	var (
		lt           = len(list)
		lids         = make([]int64, 0, lt)
		wids         = make([]int64, 0, lt)
		cont         map[int64]*l.LikeContent
		accRly       *accapi.InfosReply
		ip           = metadata.String(c, metadata.RemoteIP)
		followersRly *accapi.RelationsReply
	)
	for _, v := range list {
		lids = append(lids, v.ID)
		if v.Wid > 0 {
			wids = append(wids, v.Wid)
		}
	}
	eg, errCtx := errgroup.WithContext(c)
	eg.Go(func() (e error) {
		cont, e = s.dao.LikeContent(errCtx, lids)
		return
	})
	eg.Go(func() (e error) {
		accRly, e = s.accClient.Infos3(errCtx, &accapi.MidsReq{Mids: wids})
		return
	})
	if mid > 0 {
		eg.Go(func() (e error) {
			followersRly, e = s.accClient.Relations3(errCtx, &accapi.RelationsReq{Mid: mid, Owners: wids, RealIp: ip})
			return
		})
	}
	if err = eg.Wait(); err != nil {
		log.Error("actContent:eg.Wait() error(%v)", err)
		return
	}
	for _, v := range list {
		obj := make(map[string]interface{}, 2)
		if _, ok := cont[v.ID]; ok {
			obj["cont"] = cont[v.ID]
		}
		var t struct {
			*accapi.Info
			Following bool `json:"following"`
		}
		if _, k := accRly.Infos[v.Wid]; k {
			t.Info = accRly.Infos[v.Wid]
		}
		if mid > 0 {
			if _, f := followersRly.Relations[v.Wid]; f {
				t.Following = followersRly.Relations[v.Wid].Following
			}
		}
		obj["act"] = t
		v.Object = obj
	}
	return
}

// orderByCtime .
func (s *Service) orderByCtime(c context.Context, sid int64, pn, ps int, likesType bool) (res []*l.List, err error) {
	var (
		lids    []int64
		start   = (pn - 1) * ps
		end     = start + ps - 1
		items   map[int64]*l.Item
		likeAct map[int64]int64
	)
	if lids, err = s.dao.LikeCtime(c, sid, start, end); err != nil {
		log.Error("s.dao.LikeCtime(%d,%d,%d) error(%+v)", sid, start, end, err)
		return
	}
	lt := len(lids)
	if lt == 0 {
		return
	}
	eg, errCtx := errgroup.WithContext(c)
	eg.Go(func() (e error) {
		items, e = s.dao.Likes(errCtx, lids)
		return
	})
	if likesType {
		eg.Go(func() (e error) {
			likeAct, e = s.dao.LikeActLidCounts(c, lids)
			return
		})
	}
	if err = eg.Wait(); err != nil {
		log.Error("orderByCtime:eg.Wait() error(%+v)", err)
		return
	}
	res = make([]*l.List, 0, lt)
	for _, v := range lids {
		if _, ok := items[v]; ok && items[v].ID > 0 {
			t := &l.List{Item: items[v]}
			if likesType {
				if _, f := likeAct[v]; f {
					t.Like = likeAct[v]
				}
			}
			res = append(res, t)
		} else {
			log.Info("s.dao.CacheLikes(%d) not found", v)
		}
	}
	return
}

// orderByRandom order by random
func (s *Service) orderByRandom(c context.Context, sid int64, pn, ps int, likesType bool) (res []*l.List, err error) {
	var (
		lids      []int64
		start     = (pn - 1) * ps
		end       = start + ps - 1
		items     map[int64]*l.Item
		likeAct   map[int64]int64
		orderIDs  []int64
		orderList *l.ListInfo
	)
	if lids, err = s.dao.LikeRandom(c, sid, start, end); err != nil {
		log.Error("s.dao.LikeRandom(%d,%d,%d) error(%+v)", sid, start, end, err)
		return
	}
	lt := len(lids)
	if lt == 0 {
		if orderList, err = s.dao.ListFromES(c, sid, "", 500, 1, time.Now().Unix()); err != nil {
			log.Error("s.dao.ListFromES(%d) error(%+v)", sid, err)
			return
		}
		if orderList == nil || len(orderList.List) == 0 {
			return
		}
		orderLen := len(orderList.List)
		orderIDs = make([]int64, 0, orderLen)
		for _, v := range orderList.List {
			orderIDs = append(orderIDs, v.ID)
		}
		if err = s.dao.SetLikeRandom(c, sid, orderIDs); err != nil {
			log.Error("s.dao.SetLikeRandom(%d) error(%+v)", sid, err)
			return
		}
		if lids, err = s.dao.LikeRandom(c, sid, start, end); err != nil {
			log.Error("s.dao.LikeRandom(%d,%d,%d) error(%+v)", sid, start, end, err)
			return
		}
	}
	if len(lids) == 0 {
		return
	}
	eg, errCtx := errgroup.WithContext(c)
	eg.Go(func() (e error) {
		items, e = s.dao.Likes(errCtx, lids)
		return
	})
	if likesType {
		eg.Go(func() (e error) {
			likeAct, _ = s.dao.LikeActLidCounts(c, lids)
			return
		})
	}
	if err = eg.Wait(); err != nil {
		log.Error("orderByRandom:eg.Wait() error(%+v)", err)
		return
	}
	res = make([]*l.List, 0, lt)
	for _, v := range lids {
		if _, ok := items[v]; ok && items[v].ID > 0 {
			t := &l.List{Item: items[v]}
			if likesType {
				if _, f := likeAct[v]; f {
					t.Like = likeAct[v]
				}
			}
			res = append(res, t)
		} else {
			log.Info("s.dao.orderByRandom(%d) not found", v)
		}
	}
	return
}

// orderByLike only fo like .
func (s *Service) orderByLike(c context.Context, sid int64, pn, ps int) (res []*l.List, err error) {
	var (
		lids  []int64
		lt    int
		items map[int64]*l.Item
		start = (pn - 1) * ps
		end   = start + ps - 1
	)
	infos, err := s.dao.RedisCache(c, sid, start, end)
	if err != nil {
		log.Error("s.dao.RedisCache(%d,%d,%d) error(%+v)", sid, start, end, err)
		return
	}
	lt = len(infos)
	if lt == 0 {
		return
	}
	lids = make([]int64, 0, lt)
	for _, v := range infos {
		lids = append(lids, v.Lid)
	}
	if items, err = s.dao.Likes(c, lids); err != nil {
		log.Error("s.dao.CacheLikes(%v) error(%+v)", lids, err)
		return
	}
	res = make([]*l.List, 0, lt)
	for _, v := range infos {
		if _, ok := items[v.Lid]; ok && items[v.Lid].ID > 0 {
			t := &l.List{Item: items[v.Lid], Like: v.Score}
			res = append(res, t)
		} else {
			log.Info("s.dao.CacheLikes(%d) not found", v.Lid)
		}
	}
	return
}

// storyLikeCheck .
func (s *Service) storyLikeCheck(c context.Context, sid, lid, mid, storyMaxAct, storyEachMaxAct int64) (left int64, err error) {
	var (
		sumScore, lScore int64
		maxLeft, lLeft   int64
	)
	eg, errCtx := errgroup.WithContext(c)
	eg.Go(func() (e error) {
		sumScore, e = s.storySumUsed(errCtx, sid, mid)
		return
	})
	eg.Go(func() (e error) {
		lScore, e = s.storyEachUsed(errCtx, sid, mid, lid)
		return
	})
	if err = eg.Wait(); err != nil {
		err = errors.Wrap(err, "eg.Wait()")
		return
	}
	maxLeft = storyMaxAct - sumScore
	lLeft = storyEachMaxAct - lScore
	left = int64(math.Min(float64(maxLeft), float64(lLeft)))
	if left <= 0 {
		left = 0
	}
	return
}

// storyLikeActSet .
func (s *Service) storyLikeActSet(c context.Context, sid, lid, mid int64, score int64) (err error) {
	eg, errCtx := errgroup.WithContext(c)
	eg.Go(func() (e error) {
		_, e = s.dao.IncrStoryLikeSum(errCtx, sid, mid, score)
		return
	})
	eg.Go(func() (e error) {
		_, e = s.dao.IncrStoryEachLikeAct(errCtx, sid, mid, lid, score)
		return
	})
	if err = eg.Wait(); err != nil {
		log.Error("storyLikeActSet:eg.Wait() error(%+v)", err)
	}
	return
}

// storySumUsed .
func (s *Service) storySumUsed(c context.Context, sid, mid int64) (res int64, err error) {
	if res, err = s.dao.StoryLikeSum(c, sid, mid); err != nil {
		log.Error("s.dao.StoryLikeSum(%d,%d) error(%+v)", sid, mid, err)
		return
	}
	if res == -1 {
		today := time.Now().Format("2006-01-02")
		etime := fmt.Sprintf("%s 23:59:59", today)
		stime := fmt.Sprintf("%s 00:00:00", today)
		if res, err = s.dao.StoryLikeActSum(c, sid, mid, stime, etime); err != nil {
			log.Error("s.dao.StoryLikeActSum(%d,%d) error(%+v)", sid, mid, err)
			return
		}
		if err = s.dao.SetLikeSum(c, sid, mid, res); err != nil {
			log.Error("s.dao.SetLikeSum(%d,%d,%d) error(%+v)", sid, mid, res, err)
		}
	}
	return
}

func (s *Service) storyEachUsed(c context.Context, sid, mid, lid int64) (res int64, err error) {

	if res, err = s.dao.StoryEachLikeSum(c, sid, mid, lid); err != nil {
		log.Error("s.dao.StoryEachLikeSum(%d,%d) error(%+v)", sid, mid, err)
		return
	}
	if res == -1 {
		today := time.Now().Format("2006-01-02")
		etime := fmt.Sprintf("%s 23:59:59", today)
		stime := fmt.Sprintf("%s 00:00:00", today)
		if res, err = s.dao.StoryEachLikeAct(c, sid, mid, lid, stime, etime); err != nil {
			log.Error("s.dao.StoryLikeActSum(%d,%d) error(%+v)", sid, mid, err)
			return
		}
		if err = s.dao.SetEachLikeSum(c, sid, mid, lid, res); err != nil {
			log.Error("s.dao.SetEachLikeSum(%d,%d,%d) error(%+v)", sid, mid, res, err)
		}
	}
	return
}

// LikeActList get sid&lid likeact list .
func (s *Service) LikeActList(c context.Context, sid, mid int64, lids []int64) (res map[int64]interface{}, err error) {
	var (
		likeCounts map[int64]int64
		likeActs   map[int64]int
		likeCount  int64
		isLike     int
	)
	group, ctx := errgroup.WithContext(c)
	group.Go(func() (e error) {
		if likeCounts, e = s.dao.LikeActLidCounts(ctx, lids); e != nil {
			log.Error("s.dao.LikeActLidCounts(%v) error(%+v)", lids, e)
			return e
		}
		return nil
	})
	if mid > 0 {
		group.Go(func() (e error) {
			if likeActs, e = s.dao.LikeActs(ctx, sid, mid, lids); e != nil {
				log.Error("s.dao.LikeActMidList(%v,%d,%d) error(%+v)", lids, sid, mid, e)
				return e
			}
			return nil
		})
	}
	if err = group.Wait(); err != nil {
		log.Error("get likeactListerror(%v)", err)
		return
	}
	res = make(map[int64]interface{}, len(lids))
	for _, lid := range lids {
		if _, ok := likeCounts[lid]; ok {
			likeCount = likeCounts[lid]
		} else {
			likeCount = 0
		}
		if _, ok := likeActs[lid]; ok {
			isLike = likeActs[lid]
		} else {
			isLike = 0
		}
		res[lid] = map[string]interface{}{
			"likeCount": likeCount,
			"isLike":    isLike,
		}
	}
	return
}

// isLikeType range liketype find out real type .
func (s *Service) isLikeType(c context.Context, subType int64) (res string) {
	for _, ty := range l.LIKETYPE {
		if subType == ty {
			res = _like
			return
		}
	}
	if subType == l.MUSIC {
		res = _vote
	} else {
		res = _grade
	}
	return
}

// simpleJudge judge user could like or not .
func (s *Service) simpleJudge(c context.Context, subject *l.SubjectItem, member *accapi.Profile) (err error) {
	if member.Silence == _silenceForbid {
		err = ecode.ActivityMemberBlocked
		return
	}
	if subject.Flag == 0 {
		return
	}
	if (subject.Flag & l.FLAGSPY) == l.FLAGSPY {
		var userScore *spymdl.UserScore
		if userScore, err = s.spy.UserScore(c, &spymdl.ArgUserScore{Mid: member.Mid}); err != nil {
			log.Error("s.spy.UserScore(%d) error(%v)", member.Mid, err)
			return
		}
		if int64(userScore.Score) <= s.c.Rule.Spylike {
			err = ecode.ActivityLikeScoreLower
			return
		}
	}
	if (subject.Flag & l.FLAGUSTIME) == l.FLAGUSTIME {
		if subject.Ustime <= xtime.Time(member.JoinTime) {
			err = ecode.ActivityLikeRegisterLimit
			return
		}
	}
	if (subject.Flag & l.FLAGUETIME) == l.FLAGUETIME {
		if subject.Uetime >= xtime.Time(member.JoinTime) {
			err = ecode.ActivityLikeBeforeRegister
			return
		}
	}
	if (subject.Flag & l.FLAGPHONEBIND) == l.FLAGPHONEBIND {
		if member.TelStatus != 1 {
			err = ecode.ActivityTelValid
			return
		}
	}
	if (subject.Flag & l.FLAGLEVEL) == l.FLAGLEVEL {
		if subject.Level > int64(member.Level) {
			err = ecode.ActivityLikeLevelLimit
		}
	}
	return
}

// judgeUser judge user could like or not .
func (s *Service) judgeUser(c context.Context, subject *l.SubjectItem, member *accapi.Profile) (err error) {
	if member.Silence == _silenceForbid {
		err = ecode.ActivityMemberBlocked
		return
	}
	if subject.Flag == 0 {
		return
	}
	if (subject.Flag & l.FLAGIP) == l.FLAGIP {
		ip := metadata.String(c, metadata.RemoteIP)
		var used int
		if used, err = s.dao.IPReqquestCheck(c, ip); err != nil {
			log.Error("s.dao.IpReqquestCheck(%s) error(%+v)", ip, err)
			return
		}
		if used == 0 {
			if err = s.dao.SetIPRequest(c, ip); err != nil {
				log.Error("s.dao.SetIPRequest(%s) error(%+v)", ip, err)
				return
			}
		} else {
			err = ecode.ActivityLikeIPFrequence
			return
		}
	}
	if (subject.Flag & l.FLAGSPY) == l.FLAGSPY {
		var userScore *spymdl.UserScore
		if userScore, err = s.spy.UserScore(c, &spymdl.ArgUserScore{Mid: member.Mid}); err != nil {
			log.Error("s.spy.UserScore(%d) error(%v)", member.Mid, err)
			return
		}
		if int64(userScore.Score) <= s.c.Rule.Spylike {
			err = ecode.ActivityLikeScoreLower
			return
		}
	}
	if (subject.Flag & l.FLAGUSTIME) == l.FLAGUSTIME {
		if subject.Ustime <= xtime.Time(member.JoinTime) {
			err = ecode.ActivityLikeRegisterLimit
			return
		}
	}
	if (subject.Flag & l.FLAGUETIME) == l.FLAGUETIME {
		if subject.Uetime >= xtime.Time(member.JoinTime) {
			err = ecode.ActivityLikeBeforeRegister
			return
		}
	}
	if (subject.Flag & l.FLAGPHONEBIND) == l.FLAGPHONEBIND {
		if member.TelStatus != 1 {
			err = ecode.ActivityTelValid
			return
		}
	}
	if (subject.Flag & l.FLAGLEVEL) == l.FLAGLEVEL {
		if subject.Level > int64(member.Level) {
			err = ecode.ActivityLikeLevelLimit
		}
	}
	return
}

// Close service
func (s *Service) Close() {
	s.dao.Close()
}

// Ping service
func (s *Service) Ping(c context.Context) (err error) {
	err = s.dao.Ping(c)
	return
}

func (s *Service) arcTypeproc() {
	for {
		if types, err := s.arcRPC.Types2(context.Background()); err != nil {
			log.Error("s.arcRPC.Types2 error(%v)", err)
			time.Sleep(time.Second)
		} else {
			s.arcType = types
		}
		time.Sleep(time.Duration(s.c.Interval.PullArcTypeInterval))
	}
}

func (s *Service) initDialect() {
	tmpTag := make(map[int64]struct{}, len(s.c.Rule.DialectTags))
	for _, v := range s.c.Rule.DialectTags {
		tmpTag[v] = struct{}{}
	}
	tmpRegion := make(map[int16]struct{}, len(s.c.Rule.DialectRegions))
	for _, v := range s.c.Rule.DialectRegions {
		tmpRegion[v] = struct{}{}
	}
	s.dialectTags = tmpTag
	s.dialectRegions = tmpRegion
}

func (s *Service) actSourceproc() {
	for {
		if s.c.Rule.DialectSid != 0 {
			s.updateActSourceList(context.Background(), s.c.Rule.DialectSid, _typeAll)
		}
		if len(s.c.Rule.SpecialSids) > 0 {
			for _, sid := range s.c.Rule.SpecialSids {
				if sid > 0 {
					s.updateActSourceList(context.Background(), sid, _typeRegion)
				}
			}
		}
		time.Sleep(time.Duration(s.c.Interval.ActSourceInterval))
	}
}

func (s *Service) initReward() {
	tmp := make(map[int]*bnjmdl.Reward, len(s.c.Bnj2019.Reward))
	for _, v := range s.c.Bnj2019.Reward {
		tmp[v.Step] = v
	}
	s.reward = tmp
}
