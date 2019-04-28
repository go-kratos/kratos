package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	accmdl "go-common/app/service/main/account/model"
	"go-common/app/service/main/spy/model"
	"go-common/library/log"
)

// HandlerFunc block filter handler function.
type HandlerFunc func(af *Args)

// Args filter args.
type Args struct {
	block bool
	ui    *model.UserInfo
}

//BlockNo block no.
func (s *Service) BlockNo(seconds int64) int64 {
	return time.Now().Unix() / seconds
}

func (s *Service) scoreLessHandler(af *Args) {
	v, ok := s.Config(model.LessBlockScore)
	if !ok {
		log.Error("scoreLessHandler get config error(%s,%v)", model.LessBlockScore, s.spyConfig)
		return
	}
	if af.ui.Score <= v.(int8) {
		af.block = true
		return
	}
}

//reVerifyHandler double check user lv and vip info
func (s *Service) reVerifyHandler(c context.Context, ui *model.UserInfo) (b bool, err error) {
	b = true
	account, err := s.accInfoRetry(c, ui.Mid)
	if err != nil || account == nil {
		return
	}
	// vip or lv4
	if (account.Level < s.c.Property.DoubleCheckLevel) &&
		(account.Vip.Status != model.VipEnableStatus || account.Vip.Type == model.VipNonType) {
		return
	}
	udb, err := s.dao.UserInfo(c, ui.Mid)
	if err != nil {
		log.Error(" s.dao.UserInfo(%d) err(%v)", ui.Mid, err)
		return
	}
	if udb.ReliveTimes > model.ReliveCheckTimes {
		s.promBlockInfo.Incr("security_login_times_out_count")
		return
	}
	b = false
	l, err := s.dao.SetNXLockCache(c, ui.Mid)
	if !l || err != nil {
		log.Error("cycleblock had run (%v,%v)", l, err)
		return
	}
	if err = s.ReBuildPortrait(c, ui.Mid, model.DoubleCheckRemake); err != nil {
		log.Error(" s.ReBuildPortrait(%d) err(%v)", ui.Mid, err)
		return
	}
	reason, _ := s.blockReason(c, ui.Mid)
	if err = s.dao.SecurityLogin(c, ui.Mid, reason); err != nil {
		log.Error(" s.dao.SecurityLogin(%d) err(%v)", ui.Mid, err)
		return
	}
	s.promBlockInfo.Incr("security_login_count")
	return
}

// Verify block filter.
func (s *Service) Verify(c context.Context, af *Args, handlers ...HandlerFunc) (b bool) {
	v, ok := s.Config(model.AutoBlock)
	if !ok {
		log.Error("Verfiy get config error(%s,%v)", model.AutoBlock, s.spyConfig)
		return
	}
	if v.(int8) != model.AutoBlockOpen {
		log.Info("autoBlock Close(%v)", af)
		return
	}
	for _, h := range handlers {
		h(af)
		if b = af.block; b {
			break
		}
	}
	// double check
	if b {
		s.promBlockInfo.Incr("expect_block_count")
		b, _ = s.reVerifyHandler(c, af.ui)
	}
	return
}

// BlockFilter do user block.
func (s *Service) BlockFilter(c context.Context, ui *model.UserInfo) (state int8, err error) {
	var (
		args    = &Args{ui: ui}
		hs      = []HandlerFunc{s.scoreLessHandler}
		blockNo = s.BlockNo(s.c.Property.Block.CycleTimes)
	)
	state = ui.State
	if b := s.Verify(c, args, hs...); !b {
		return
	}
	state = model.StateBlock
	s.dao.AddBlockCache(c, ui.Mid, ui.Score, blockNo)
	if err = s.dao.AddPunishmentQueue(c, ui.Mid, blockNo); err != nil {
		log.Error("s.dao.AddPunishmentQueue(%d) error(%v)", ui.Mid, err)
		return
	}
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
	for _, v := range hs {
		if v.Reason == model.DoubleCheckRemake {
			continue
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

func (s *Service) accInfoRetry(c context.Context, mid int64) (account *accmdl.Card, err error) {
	for i := 1; i <= model.RetryTimes; i++ {
		account, err = s.accRPC.Card3(c, &accmdl.ArgMid{Mid: mid})
		if account == nil || err != nil {
			log.Error(" s.accRPC.Infos2(%d) err(%v)", mid, err)
			continue
		}
	}
	return
}
