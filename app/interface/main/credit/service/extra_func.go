package service

import (
	"context"
	"math/rand"
	"sync"
	"time"

	model "go-common/app/interface/main/credit/model"
	acmdl "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

func (s *Service) caseVoteID(c context.Context, mid int64, pubCid int64) (cid int64, err error) {
	// 获取发放中cids的列表
	mcases, err := s.dao.GrantCases(c)
	if err != nil {
		log.Error("s.dao.GrantCases error(%v)", err)
		return
	}
	if len(mcases) == 0 {
		log.Warn("no grant cases(%+v)!", mcases)
		return
	}
	// 7天内已投cids
	weekcases, err := s._caseObtainMID(c, mid, 7, model.CaseObtainNoToday)
	if err != nil {
		log.Error("s._caseObtainMID(%d, 7, %t) error(%v)", mid, model.CaseObtainNoToday, err)
		return
	}
	// 今天内已投cids
	todaycases, err := s._caseObtainMID(c, mid, 0, model.CaseObtainToday)
	if err != nil {
		log.Error("s._caseObtainMID(%d, 0, %t) error(%v)", mid, model.CaseObtainToday, err)
		return
	}
	todayPubVote := 0
	tpmcids := make(map[int64]struct{})
	// 今天已投cid的map
	for _, tcase := range todaycases {
		if tcase.CaseType == model.JudeCaseTypePublic {
			tpmcids[tcase.ID] = struct{}{}
			todayPubVote++
		}
	}
	vcids := []int64{}
	for wcid := range weekcases {
		vcids = append(vcids, wcid)
	}
	for cid, m := range mcases {
		if _, ok := tpmcids[cid]; !ok && m.CaseType == model.JudeCaseTypePublic {
			todayPubVote++
		}
		// 把被举报人是风纪委用户自己和结案时间后10分钟的稿件  加入用户已投稿列表
		if m.Mid == mid || m.Etime.Time().Add(-time.Duration(s.c.Judge.ReservedTime)).Before(time.Now()) {
			vcids = append(vcids, cid)
			continue
		}
	}
	tLen := len(todaycases)
	// 获取案件最大数判断
	if int64(tLen-todayPubVote) >= s.c.Judge.CaseObtainMax {
		err = ecode.CreditCaseMax
		return
	}
	// 制作非可投cids的map
	vmcids := make(map[int64]struct{})
	for _, uncid := range vcids {
		vmcids[uncid] = struct{}{}
	}
	// 小众
	pteCids := make([]int64, 0)
	// 大众
	pubCids := make([]int64, 0)
	// 大众cid的map
	mpCids := make(map[int64]struct{})
	// 取出可投cids
	for kcid, m := range mcases {
		if _, ok := vmcids[kcid]; ok {
			continue
		}
		if m.CaseType == model.JudeCaseTypePublic {
			pubCids = append(pubCids, kcid)
			mpCids[kcid] = struct{}{}
		} else {
			pteCids = append(pteCids, kcid)
		}
	}
	pubLen := len(pubCids)
	pteLen := len(pteCids)
	// 没有可投的案件
	if pubLen+pteLen == 0 {
		log.Warn("mid(%d) no case can vote!", mid)
		return
	}
	var caseType int8
	_, ok := mpCids[pubCid]
	if pubCid != 0 && !ok {
		return
	}
	if pubCid != 0 && ok {
		cid = pubCid
		caseType = model.JudeCaseTypePublic
	} else {
		radio := rand.New(rand.NewSource(time.Now().UnixNano()))
		if pubLen > 0 {
			cid = s._randCid(pubCids, radio)
			caseType = model.JudeCaseTypePublic
		} else if pteLen > 0 {
			cid = s._randCid(pteCids, radio)
			caseType = model.JudeCaseTypePrivate
		}
	}
	// db插入用户投票数据
	if err = s.dao.InsVote(c, mid, cid, s.c.Judge.CaseCheckTime); err != nil {
		log.Error("s.dao.InsVote( mid(%d), cid(%d), s.c.Judge.CaseCheckTime(%d)) error(%v)", mid, cid, s.c.Judge.CaseCheckTime, err)
		return
	}
	mcid := &model.SimCase{ID: cid, CaseType: caseType}
	// 从redis的set中设置用户已投cids
	if err = s.dao.SetVoteCaseMID(c, mid, mcid); err != nil {
		log.Error("s.dao.SetVoteCaseMID(%d,%+v) error(%v)", mid, mcid, err)
		return
	}
	log.Info("CaseObtain mid:%d total:%d CaseObtainMax:%d cid:%d", mid, int64(tLen+todayPubVote), s.c.Judge.CaseObtainMax, cid)
	// db插入case投放总数
	if err = s.dao.AddCaseVoteTotal(c, "put_total", cid, 1); err != nil {
		log.Error("s.dao.InsVote( mid(%d), cid(%d), s.c.Judge.CaseCheckTime(%d)) error(%v)", mid, cid, s.c.Judge.CaseCheckTime, err)
	}
	return
}

// 获取用户N天内已投列表
func (s *Service) _caseObtainMID(c context.Context, mid int64, day int, isToday bool) (cases map[int64]*model.SimCase, err error) {
	isExpired, err := s.dao.IsExpiredObtainMID(c, mid, isToday)
	if err != nil {
		log.Error("s.dao.IsExpiredObtainMID(%d,%t) error(%v)", mid, isToday, err)
		return
	}
	if isExpired {
		cases, err = s.dao.CaseObtainMID(c, mid, isToday)
		if err != nil {
			log.Error("s.dao.CaseObtainByMID(%d,%t) error(%v)", mid, isToday, err)
			return
		}
	} else {
		if cases, err = s.dao.LoadVoteIDsMid(c, mid, day); err != nil {
			log.Error("s.dao.LoadVoteIDsMid(%d,%d) error(%v)", mid, day, err)
			return
		}
		if len(cases) == 0 {
			return
		}
		s.addCache(func() {
			if err = s.dao.LoadVoteCaseMID(context.TODO(), mid, cases, isToday); err != nil {
				log.Error("s.dao.LoadVoteCaseMID(%d,%v,%t) error(%v)", mid, cases, isToday, err)
				return
			}
		})
	}
	return
}

// cid 取值的随机算法
func (s *Service) _randCid(cids []int64, radio *rand.Rand) (cid int64) {
	// 随机取出数组的游标
	rand := int64(radio.Intn(len(cids)))
	cid = cids[rand]
	return
}

// 列表批量异步获取用户信息
func (s *Service) infoMap(c context.Context, uids []int64) (infoMap map[int64]*acmdl.Info, err error) {
	total := len(uids)
	pageNum := total / model.JuryMultiJuryerInfoMax
	if total%model.JuryMultiJuryerInfoMax != 0 {
		pageNum++
	}
	var (
		g  errgroup.Group
		lk sync.RWMutex
	)
	infoMap = make(map[int64]*acmdl.Info, total)
	for i := 0; i < pageNum; i++ {
		start := i * model.JuryMultiJuryerInfoMax
		end := (i + 1) * model.JuryMultiJuryerInfoMax
		if end > total {
			end = total
		}
		g.Go(func() (err error) {
			var (
				arg = &acmdl.MidsReq{Mids: uids[start:end]}
				res *acmdl.InfosReply
			)
			if res, err = s.accountClient.Infos3(c, arg); err != nil {
				log.Error("s.accountClient.Infos3(%v) error(%v)", arg, err)
				err = nil
			} else {
				for uid, info := range res.Infos {
					lk.Lock()
					infoMap[uid] = info
					lk.Unlock()
				}
			}
			return
		})
	}
	if err = g.Wait(); err != nil {
		log.Error("g.Wait() error(%v)", err)
	}
	return
}

// userInfo one user card.
func (s *Service) userInfo(c context.Context, mid int64) (card *acmdl.CardReply, err error) {
	if mid == 0 {
		return
	}
	arg := &acmdl.MidReq{
		Mid: mid,
	}
	if card, err = s.accountClient.Card3(c, arg); err != nil {
		log.Error("s.accountClient.Card3(%+v) error(%+v)", arg, err)
	}
	return
}
