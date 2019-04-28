package service

import (
	"context"
	"fmt"

	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

func (s *Service) upBindTag(c context.Context, mid, aid int64, tags string, typeID int16) (err error) {
	typeName := s.typeName(typeID)
	topTypeName := s.topTypeName(typeID)
	if topTypeName != "" {
		typeName = fmt.Sprintf("%s,%s", typeName, topTypeName)
	}

	for i := 0; i < 3; i++ {
		if err = s.tagDao.UpBind(c, mid, aid, tags, typeName, ""); err == nil {
			return
		}
	}
	if err != nil {
		log.Error("upBindTag s.tagDao.UpBind(%d,%d,%s,%s) typeid(%d) error(%v)", mid, aid, tags, typeName, typeID, err)
	}
	return
}

func (s *Service) adminBindTag(c context.Context, mid, aid int64, tags string, typeID int16) (err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			log.Error("s.adminBindTag panic(%v)", pErr)
		}
	}()
	typeName := s.typeName(typeID)
	topTypeName := s.topTypeName(typeID)
	if topTypeName != "" {
		typeName = fmt.Sprintf("%s,%s", typeName, topTypeName)
	}

	for i := 0; i < 3; i++ {
		if err = s.tagDao.AdminBind(c, mid, aid, tags, typeName, ""); err == nil {
			return
		}
	}
	if err != nil {
		log.Error("adminBindTag s.tagDao.AdminBind(%d,%d,%s,%s) typeid(%d) error(%v)", mid, aid, tags, typeName, typeID, err)
	}
	return
}

func (s *Service) checkChannelReview(c context.Context, aid int64) (channelReview bool, channelIDs string, err error) {
	var (
		tmp bool
	)
	if tmp, channelIDs, err = s.tagDao.CheckChannelReview(c, aid); err != nil {
		log.Error("checkChannelReview s.tagDao.CheckChannelReview error(%v) aid(%d)", err, aid)
		return
	}
	if !tmp {
		log.Info("checkChannelReview s.tagDao.CheckChannelReview aid(%d) not in channelreview", aid)
		return
	}

	channelReview = true
	return
}

func (s *Service) txAddOrUpRecheckState(c context.Context, tx *sql.Tx, tp int, aid int64, state int8) (recheck *archive.Recheck, err error) {
	if recheck, err = s.arc.RecheckByAid(c, tp, aid); err != nil {
		log.Error("txAddOrUpRecheckState s.arc.RecheckByAid(%d,%d) error(%v)", tp, aid, err)
		return
	}
	if recheck != nil && recheck.State == state {
		return
	}

	if recheck == nil {
		recheck = &archive.Recheck{Type: tp, Aid: aid, State: archive.RecheckStateWait}
		if recheck.ID, err = s.arc.TxAddRecheckAID(tx, tp, aid); err != nil {
			log.Error("txAddOrUpRecheckState s.arc.TxAddRecheckAID(%d,%d) error(%v)", tp, aid, err)
			return
		}
	}
	if recheck.State != state {
		if _, err = s.arc.TxUpRecheckState(tx, tp, aid, state); err != nil {
			log.Error("txAddOrUpRecheckState s.arc.TxUpRecheckState(%d,%d,%d) error(%v)", tp, aid, state, err)
			return
		}
		recheck.State = state
	}
	return
}

func (s *Service) txAddChannelReview(c context.Context, tx *sql.Tx, aid int64) (operCont string, operRemark string, err error) {
	var (
		need       bool
		channelIDs string
	)
	if need, channelIDs, err = s.checkChannelReview(c, aid); err != nil || !need {
		return
	}
	log.Info("start to add channel review aid(%d)", aid)
	if _, err = s.txAddOrUpRecheckState(c, tx, archive.TypeChannelRecheck, aid, archive.RecheckStateWait); err != nil {
		log.Error("txAddChannelReview s.txAddOrUpRecheckState(%d) error(%v)", aid, err)
		return
	}

	if _, operCont, err = s.txAddOrUpdateFlowState(c, tx, aid, archive.FLowGroupIDChannel, 399, archive.PoolArcForbid, archive.FlowOpen, "投稿开启频道禁止"); err != nil {
		log.Error("txAddChannelReview s.txAddOrUpdateFlowState(%d) error(%v)", aid, err)
		return
	}
	operRemark = fmt.Sprintf("待频道回查,频道ID:%s", channelIDs)
	return
}
