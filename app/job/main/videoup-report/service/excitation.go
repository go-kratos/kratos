package service

import (
	"context"
	"go-common/app/job/main/videoup-report/dao/data"
	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/app/job/main/videoup-report/model/manager"
	"go-common/library/log"
	"time"
)

// hdlExcitation 激励回查逻辑
func (s *Service) hdlExcitation(nw, old *archive.Archive) (err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			log.Error("s.hdlExcitation(%v,%v) panic(%v)", nw, old, pErr)
		}
	}()
	var (
		c        = context.TODO()
		state    int8
		archives []*archive.Archive
		addits   map[int64]*archive.Addit
		now      = time.Now()
		aids     []int64
		isIgn    bool
		aCount   = 20 //自制投稿阀值
	)
	log.Info("hdlExcitation() begin new archive(%v) old archive(%v)", nw, old)
	if old != nil { //只有新投稿才需要继续
		return
	}
	//检查UP主是否在白名单中
	if isIgn, err = s.isIgnMidExcitation(c, nw.Mid); err != nil {
		log.Error("s.hdlExcitation() new(%v) old(%v) s.isIgnMidExcitation() error(%v)", nw, old, err)
		go s.hdlExcitationRetry(nw, old)
		return
	} else if isIgn {
		log.Error("s.hdlExcitation() new(%v) old(%v) ignore white mid(%d)", nw, old, nw.Mid)
		return
	}
	if state, err = s.dataDao.UpProfitState(c, nw.Mid); err != nil {
		log.Error("hdlExcitation() s.dataDao.UpProfitState(%d) error(%v)", nw.Mid, err)
		go s.hdlExcitationRetry(nw, old)
		return
	}
	if state != data.UpProfitStateSigned {
		log.Info("hdlExcitation() ignore 非签约UP主(%d) state(%d)", nw.Mid, state)
		return
	}
	//查询7天内自制、过审、定时过审稿件
	st := now.Add(-168 * time.Hour)
	if archives, err = s.arc.ExcitationArchivesByTime(c, nw.Mid, st, now); err != nil {
		log.Error("hdlExcitation() s.arc.ExcitationAidsByTime(%d,%v,%v) error(%v)", nw.Mid, st, now, err)
		go s.hdlExcitationRetry(nw, old)
		return
	}
	if len(archives) == 0 {
		log.Info("hdlExcitation() ignore UP主(%d) 在 %v 到 %v 时间内没有自制、非特殊分区、非私单、非商单稿件。", nw.Mid, st, now)
		return
	}
	for _, a := range archives {
		if s.isAuditType(a.TypeID) || a.AttrVal(archive.AttrBitIsPOrder) == archive.AttrYes {
			log.Info("hdlExcitation() 忽略特殊分区、私单 archive(%v)", a)
			continue
		}
		aids = append(aids, a.ID)
	}
	if addits, err = s.arc.Addits(c, aids); err != nil {
		log.Error("hdlExcitation() aid(%d) s.arc.Addits(%v) error(%v)", nw.ID, aids, err)
		go s.hdlExcitationRetry(nw, old)
		return
	}

	for i := 0; i < len(aids); i++ {
		if _, ok := addits[aids[i]]; ok && addits[aids[i]].OrderID > 0 {
			log.Info("hdlExcitation() 忽略商单 addit(%v)", addits[aids[i]])
			aids = append(aids[:i], aids[i+1:]...)
			i--
			continue
		}
	}

	if len(aids) < aCount {
		log.Info("hdlExcitation() ignore UP主(%d) 在 %v 到 %v 时间内只有 %d 个自制、非特殊分区、非私单、非商单稿件没有达到 %d 个。aids:%v", nw.Mid, st, now, len(aids), aCount, aids)
		return
	}
	log.Info("hdlExcitation() 符合激励回查的UP主(%d) 正在插入激励回查aids(%v)", nw.Mid, aids)
	if err = s.arc.AddRecheckAids(c, archive.TypeExcitationRecheck, aids, true); err != nil {
		log.Error("s.hdlExcitation() s.arc.AddRecheckAids error(%v)", err)
		go s.hdlExcitationRetry(nw, old)
		return
	}
	return
}

// isIgnMidExcitation 检查UP主是否在白名单中
func (s *Service) isIgnMidExcitation(c context.Context, mid int64) (is bool, err error) {
	groups, err := s.dataDao.MidGroups(c, mid)
	if err != nil {
		log.Error("s.isIgnMidExcitation(%d) error(%v)", mid, err)
		return
	}
	if _, ok := groups[manager.UpTypeExcitationWhite]; ok {
		is = true
		return
	}
	return
}

// ignoreUpsExcitation 将UP主的未回查的稿件从激励回查去除
func (s *Service) ignoreUpsExcitation(c context.Context, mid int64) (err error) {
	log.Info("s.ignoreUpsExcitation(%d)", mid)
	if err = s.arc.UpdateMidRecheckState(c, archive.TypeExcitationRecheck, mid, archive.RecheckStateIgnore); err != nil {
		log.Error("s.ignoreUpsExcitation(%d) error(%v)", mid, err)
		return
	}
	return
}

// hdlExcitationRetry 激励回查重试
func (s *Service) hdlExcitationRetry(nw, old *archive.Archive) (err error) {
	mTime, _ := time.ParseInLocation("2006-01-02 15:04:05", nw.MTime, time.Local)
	if time.Now().Unix()-mTime.Unix() > 3600 {
		log.Error("s.hdlExcitationRetry() new(%v) old(%v) the retry handler finish. (over 60 mins)", nw, old)
		return
	}
	log.Error("s.hdlExcitationRetry() new(%v) old(%v) the handler will retry in 10 seconds", nw, old)
	time.Sleep(10 * time.Second)
	err = s.hdlExcitation(nw, old)
	return
}
