package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/spy/conf"
	"go-common/app/job/main/spy/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

func (s *Service) lastBlockNo(seconds int64) int64 {
	return time.Now().Unix()/seconds - 1 //cron -1 get the last hours block mids
}

//BlockTask block task.
func (s *Service) BlockTask(c context.Context) {
	var (
		blockNo = s.lastBlockNo(s.c.Property.Block.CycleTimes)
	)
	v, ok := s.Config(model.AutoBlock)
	if !ok {
		log.Error("Verfiy get config error(%s,%v)", model.AutoBlock, s.spyConfig)
		return
	}
	if v.(int8) != model.AutoBlockOpen {
		log.Info("autoBlock Close(%s)", blockNo)
		return
	}
	mids, _ := s.blockUsers(c, blockNo)
	if len(mids) == 0 {
		log.Info("s.blockUsers len is zero")
		return
	}
	if err := s.dao.SetBlockCache(c, mids); err != nil {
		log.Error("s.dao.SetBlockCache(%v) error(%v)", mids, err)
		return
	}
}

func (s *Service) blockUsers(c context.Context, blockNo int64) (mids []int64, err error) {
	v, ok := s.Config(model.LimitBlockCount)
	if !ok {
		log.Error("blockUsers get config error(%v)", s.spyConfig)
		return
	}
	if mids, err = s.dao.BlockMidCache(c, blockNo, v.(int64)); err != nil {
		log.Error("s.dao.BlockMidCache(%s, %d) error(%v)", blockNo, v.(int64), err)
		return
	}
	return
}

func (s *Service) blockByMid(c context.Context, mid int64) (err error) {
	ui, ok := s.canBlock(c, mid)
	if !ok {
		log.Info("s.canBlock user had block(%d)", mid)
		return
	}
	reason, remake := s.blockReason(c, mid)
	if err = s.block(c, mid, ui, reason, remake); err != nil {
		log.Error("s.block(%d,%v,%s) err(%v)", mid, ui, reason, err)
	}
	return
}

func (s *Service) block(c context.Context, mid int64, ui *model.UserInfo, reason string, remake string) (err error) {
	var (
		tx *sql.Tx
	)
	ui.State = model.StateBlock
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("s.dao.BeginTran() err(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	if err = s.dao.BlockAccount(c, ui.Mid, reason); err != nil {
		log.Error("s.dao.BlockAccount(%d) error(%v)", ui.Mid, err)
		return
	}
	var ueh = &model.UserEventHistory{
		Mid:        mid,
		EventID:    conf.Conf.Property.BlockEvent,
		Score:      ui.Score,
		BaseScore:  ui.BaseScore,
		EventScore: ui.EventScore,
		Remark:     remake,
		Reason:     "自动封禁",
	}

	if err = s.dao.TxAddEventHistory(c, tx, ueh); err != nil {
		log.Error("s.dao.TxAddEventHistory(%+v) error(%v)", ueh, err)
		return
	}
	if err = s.dao.TxAddPunishment(c, tx, ui.Mid, model.PunishmentTypeBlock,
		reason, s.lastBlockNo(s.c.Property.Block.CycleTimes)); err != nil {
		log.Error("s.dao.TxAddPunishment(%d,%s) error(%v)", ui.Mid, reason, err)
		return
	}
	// update user state.
	if err = s.dao.TxUpdateUserState(c, tx, ui); err != nil {
		log.Error("s.dao.TxUpdateUserState(%v) error(%v)", ui, err)
		return
	}
	s.promBlockInfo.Incr("actual_block_count")
	return
}

func (s *Service) canBlock(c context.Context, mid int64) (ui *model.UserInfo, ok bool) {
	var (
		err error
	)
	ui, err = s.dao.UserInfo(c, mid)
	if err != nil {
		log.Error("s.UserInfo(%d) err(%v)", mid, err)
		return
	}
	v, b := s.Config(model.LessBlockScore)
	if !b {
		log.Error("scoreLessHandler get config error(%s,%v)", model.LessBlockScore, s.spyConfig)
		return
	}
	// if blocked already , return
	if ui.State == model.StateBlock || ui.Score > v.(int8) {
		log.Info("canBlock not block(%v)", ui)
		return
	}
	ok = true
	return
}

func (s *Service) blockReason(c context.Context, mid int64) (reason string, remake string) {
	var (
		err error
		hs  []*model.UserEventHistory
		buf bytes.Buffer
	)
	if hs, err = s.dao.HistoryList(c, mid, model.BlockReasonSize); err != nil || len(hs) == 0 {
		log.Error("s.dao.HistoryList(%d) err(%v)", mid, err)
		return
	}
	m := make(map[string]int)
	for i, v := range hs {
		if i == 0 {
			remake = v.Remark
		}
		if m[v.Reason] == 0 {
			m[v.Reason] = 1
		} else {
			m[v.Reason] = m[v.Reason] + 1
		}
	}
	for k, v := range m {
		buf.WriteString(k)
		buf.WriteString("x")
		buf.WriteString(fmt.Sprintf("%d ", v))
	}
	reason = buf.String()
	return
}
