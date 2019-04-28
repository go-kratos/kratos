package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"

	"go-common/app/admin/main/member/model"
	"go-common/app/admin/main/member/model/bom"
	coin "go-common/app/service/main/coin/model"
	member "go-common/app/service/main/member/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/queue/databus/report"

	"github.com/pkg/errors"
)

const (
	_logActionBaseAudit   = "base_audit"
	_logActionDeleteSign  = "delete_sign"
	_logActionRankUpdate  = "rank_set"
	_logActionCoinUpdate  = "coin_set"
	_logActionExpUpdate   = "exp_set"
	_logActionMoralUpdate = "moral_set"
)

func (s *Service) batchBase(ctx context.Context, mids []int64) (map[int64]*model.Base, error) {
	bs := make(map[int64]*model.Base, len(mids))
	for _, mid := range mids {
		b, err := s.dao.Base(ctx, mid)
		if err != nil {
			log.Error("Failed to retrive user base by mid: %d: %+v", mid, err)
			continue
		}
		bs[b.Mid] = b
	}
	return bs, nil
}

// BaseReview is.
func (s *Service) BaseReview(ctx context.Context, arg *model.ArgBaseReview) ([]*model.BaseReview, error) {
	mids := arg.Mids()
	if len(mids) == 0 {
		return nil, ecode.RequestErr
	}
	if len(mids) > 200 {
		//mids = mids[:200]
		return nil, ecode.SearchMidOverLimit
	}
	var mrs []*model.BaseReview
	for _, mid := range mids {
		base, err := s.dao.Base(ctx, mid)
		if err != nil {
			log.Error("Failed to retrive user base by mid: %d: %+v", mid, err)
			continue
		}
		logs, err := s.dao.SearchUserAuditLog(ctx, mid)
		if err != nil {
			log.Error("Failed to search user audit log, mid: %d error: %v", mid, err)
			continue
		}

		fixFaceLogs(logs.Result)

		mr := &model.BaseReview{
			Base: *base,
			Logs: logs.Result,
		}
		mrs = append(mrs, mr)
	}
	if err := s.reviewAddit(ctx, mids, mrs); err != nil {
		log.Error("Failed to fetch review violation count with mids: %+v: %+v", mids, err)
	}
	return mrs, nil
}

func fixFaceLogs(auditLogs []model.AuditLog) {
	for i, v := range auditLogs {
		if v.Type != model.BaseAuditTypeFace {
			continue
		}
		// 对头像进行重新签名
		ext := new(struct {
			Old string `json:"old"`
			New string `json:"new"`
		})
		if err := json.Unmarshal([]byte(v.Extra), &ext); err != nil {
			log.Error("Failed to unmarshal extra, additLog: %+v error: %v", v, err)
			continue
		}
		ext.New = model.BuildFaceURL(ext.New)
		extraDataBytes, err := json.Marshal(ext)
		if err != nil {
			log.Error("FaceExtra (%+v) json marshal err(%v)", ext, err)
			continue
		}
		auditLogs[i].Extra = string(extraDataBytes)
	}
}

// ClearFace is
func (s *Service) ClearFace(ctx context.Context, arg *model.ArgMids) error {
	for _, mid := range arg.Mid {
		b, err := s.dao.Base(ctx, mid)
		if err != nil {
			log.Error("Failed to retrive user base by mid: %d: %+v", mid, err)
			continue
		}
		if err = s.dao.UpFace(ctx, mid, ""); err != nil {
			log.Error("Failed to clear face mid %d error: %+v", mid, err)
			continue
		}
		privFace, err := s.mvToPrivate(ctx, urlPath(b.Face))
		if err != nil {
			log.Error("Failed to mv face To private bucket, mid: %d error: %+v", mid, err)
			err = nil
		}
		if err = s.dao.Message(ctx, "违规头像处理通知", "抱歉，由于你的头像涉嫌违规，已被修改。如有疑问请联系客服。", []int64{mid}); err != nil {
			log.Error("Failed to send message, mid: %d error: %+v", mid, err)
			err = nil
		}
		if err = s.dao.IncrViolationCount(ctx, mid); err != nil {
			log.Error("Failed to increase violation count mid: %d error: %+v", mid, err)
			err = nil
		}
		report.Manager(&report.ManagerInfo{
			Uname:    arg.Operator,
			UID:      arg.OperatorID,
			Business: model.ManagerLogID,
			Type:     model.BaseAuditTypeFace,
			Oid:      mid,
			Action:   _logActionBaseAudit,
			Ctime:    time.Now(),
			// extra
			Index: []interface{}{0, 0, 0, "", "", ""},
			Content: map[string]interface{}{
				"old": b.Face,
				"new": model.BuildFaceURL(privFace),
			},
		})
	}
	return nil
}

// ClearSign is
func (s *Service) ClearSign(ctx context.Context, arg *model.ArgMids) error {
	for _, mid := range arg.Mid {
		b, err := s.dao.Base(ctx, mid)
		if err != nil {
			log.Error("Failed to retrive user base by mid: %d: %+v", mid, err)
			continue
		}
		if err = s.dao.UpSign(ctx, mid, ""); err != nil {
			log.Error("Failed to clear sign mid %d error: %+v", mid, err)
			continue
		}
		if err = s.dao.Message(ctx, "违规签名处理通知", "抱歉，由于你的签名涉嫌违规，已被修改。如有疑问请联系客服。", []int64{mid}); err != nil {
			log.Error("Failed to send message, mid: %d error: %+v", mid, err)
			err = nil
		}
		if err = s.dao.IncrViolationCount(ctx, mid); err != nil {
			log.Error("Failed to increase violation count mid: %d error: %+v", mid, err)
			err = nil
		}
		report.Manager(&report.ManagerInfo{
			Uname:    arg.Operator,
			UID:      arg.OperatorID,
			Business: model.ManagerLogID,
			Type:     model.BaseAuditTypeSign,
			Oid:      mid,
			Action:   _logActionBaseAudit,
			Ctime:    time.Now(),
			// extra
			Index: []interface{}{0, 0, 0, "", "", ""},
			Content: map[string]interface{}{
				"old": b.Sign,
				"new": "",
			},
		})
	}
	return nil
}

// ClearName is
func (s *Service) ClearName(ctx context.Context, arg *model.ArgMids) error {
	for _, mid := range arg.Mid {
		b, err := s.dao.Base(ctx, mid)
		if err != nil {
			log.Error("Failed to retrive user base by mid: %d: %+v", mid, err)
			continue
		}
		defaultName := fmt.Sprintf("bili_%d", mid)
		if err = s.dao.UpdateUname(ctx, mid, defaultName); err != nil {
			log.Error("Failed to clear name mid %d error: %+v", mid, err)
			continue
		}
		if err = s.dao.Message(ctx, "违规昵称处理通知", "抱歉，由于你的昵称涉嫌违规，已被修改。如有疑问请联系客服。", []int64{mid}); err != nil {
			log.Error("Failed to send message, mid: %d error: %+v", mid, err)
			err = nil
		}
		if err = s.dao.IncrViolationCount(ctx, mid); err != nil {
			log.Error("Failed to increase violation count mid: %d error: %+v", mid, err)
			err = nil
		}
		report.Manager(&report.ManagerInfo{
			Uname:    arg.Operator,
			UID:      arg.OperatorID,
			Business: model.ManagerLogID,
			Type:     model.BaseAuditTypeName,
			Oid:      mid,
			Action:   _logActionBaseAudit,
			Ctime:    time.Now(),
			// extra
			Index: []interface{}{0, 0, 0, "", "", ""},
			Content: map[string]interface{}{
				"old": b.Name,
				"new": defaultName,
			},
		})
	}
	return nil
}

// Members is.
func (s *Service) Members(ctx context.Context, arg *model.ArgList) (*model.MemberPagination, error) {
	searched, err := s.dao.SearchMember(ctx, arg)
	if err != nil {
		return nil, err
	}
	mids := searched.Mids()

	bs, err := s.batchBase(ctx, mids)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Base, 0, len(mids))
	for _, mid := range mids {
		if b, ok := bs[mid]; ok {
			result = append(result, b)
		}
	}

	page := &model.MemberPagination{
		CommonPagination: searched.Pagination(),
		Members:          result,
	}
	return page, nil
}

// MemberProfile is.
func (s *Service) MemberProfile(ctx context.Context, mid int64) (*model.Profile, error) {
	p := model.NewProfile()

	// base
	b, err := s.dao.Base(ctx, mid)
	if err != nil {
		log.Error("Failed to retrive user base with mid: %d: %+v", mid, err)
		return nil, err
	}
	p.Base = *b

	// detail
	// remove later
	p.Detail = model.Detail{
		Mid:      b.Mid,
		Birthday: b.Birthday,
	}

	// exp
	e, err := s.dao.Exp(ctx, mid)
	if err != nil {
		log.Error("Failed to retrive user exp with mid: %d: %+v", mid, err)
	}
	if e != nil {
		p.Exp = *e
		p.Level.FromExp(e)
	}

	// moral
	mo, err := s.dao.Moral(ctx, mid)
	if err != nil {
		log.Error("Failed to retrive user moral with mid: %d: %+v", mid, err)
	}
	if mo != nil {
		p.Moral = *mo
	}

	// official
	of, err := s.dao.Official(ctx, mid)
	if err != nil {
		log.Error("Failed to retrive user official with mid: %d: %+v", mid, err)
	}
	if of != nil {
		p.Official = *of
	}

	// coin
	co, err := s.coinRPC.UserCoins(ctx, &coin.ArgCoinInfo{Mid: mid})
	if err != nil {
		log.Error("Failed to retrive user coins with mid: %d: %+v", mid, err)
	}
	p.Coin = model.Coin{Coins: co}

	// addit
	ad, err := s.dao.UserAddit(ctx, mid)
	if err != nil {
		log.Error("Failed to retrive user addit with mid: %d: %+v", mid, err)
	}
	if ad != nil {
		p.Addit = *ad
	}

	// realname
	dr, err := s.dao.RealnameInfo(ctx, mid)
	if err != nil {
		log.Error("Failed to retrive user realname with mid: %d: %+v", mid, err)
	}
	if dr != nil {
		p.Realanme.ParseInfo(dr)
	} else {
		p.Realanme.State = model.RealnameApplyStateNone
	}
	return p, nil
}

// DelSign is.
func (s *Service) DelSign(ctx context.Context, arg *model.ArgMids) error {
	for _, mid := range arg.Mid {
		err := s.dao.UpSign(ctx, mid, "")
		if err != nil {
			log.Error("Failed to delete sign mid: %d: %+v", mid, err)
			continue
		}
		if err := s.dao.Message(ctx, "违规签名处理通知", "抱歉，由于你的个性签名内容涉嫌违规，我们已将你的个性签名清空，如有问题请联系客服。", []int64{mid}); err != nil {
			log.Error("Failed to send message: mid: %d: %+v", mid, err)
		}
		report.Manager(&report.ManagerInfo{
			Uname:    arg.Operator,
			UID:      arg.OperatorID,
			Business: model.ManagerLogID,
			Type:     0,
			Oid:      mid,
			Action:   _logActionDeleteSign,
			Ctime:    time.Now(),
		})
	}
	return nil
}

// SetMoral is.
func (s *Service) SetMoral(ctx context.Context, arg *model.ArgMoralSet) error {
	moral, err := s.dao.Moral(ctx, arg.Mid)
	if err != nil {
		return errors.WithStack(err)
	}
	newMoral := int64(arg.Moral * 100)
	delta := newMoral - moral.Moral
	if err = s.memberRPC.AddMoral(ctx,
		&member.ArgUpdateMoral{
			Mid:      arg.Mid,
			Delta:    delta,
			Operator: arg.Operator,
			Reason:   arg.Reason,
			Remark:   "管理后台",
			Origin:   member.ManualChangeType,
			IP:       arg.IP}); err != nil {
		return errors.WithStack(err)
	}

	report.Manager(&report.ManagerInfo{
		Uname:    arg.Operator,
		UID:      arg.OperatorID,
		Business: model.ManagerLogID,
		Type:     0,
		Oid:      arg.Mid,
		Action:   _logActionMoralUpdate,
		Ctime:    time.Now(),
		// extra
		Index: []interface{}{0, 0, 0, "", "", ""},
		Content: map[string]interface{}{
			"new_moral": newMoral,
			"delta":     delta,
		},
	})
	return nil
}

// SetExp is.
func (s *Service) SetExp(ctx context.Context, arg *model.ArgExpSet) error {
	exp, err := func() (*model.Exp, error) {
		exp, err := s.dao.Exp(ctx, arg.Mid)
		if err != nil {
			if err == ecode.NothingFound {
				return &model.Exp{Mid: arg.Mid}, nil
			}
			return nil, errors.WithStack(err)
		}
		return exp, nil
	}()
	if err != nil {
		return err
	}

	delta := arg.Exp - float64(exp.Exp/100)
	if err = s.memberRPC.UpdateExp(ctx,
		&member.ArgAddExp{
			Mid:     arg.Mid,
			Count:   delta,
			Operate: arg.Operator,
			Reason:  arg.Reason,
			IP:      arg.IP}); err != nil {
		return errors.WithStack(err)
	}

	report.Manager(&report.ManagerInfo{
		Uname:    arg.Operator,
		UID:      arg.OperatorID,
		Business: model.ManagerLogID,
		Type:     0,
		Oid:      arg.Mid,
		Action:   _logActionExpUpdate,
		Ctime:    time.Now(),
		// extra
		Index: []interface{}{0, 0, 0, "", "", ""},
		Content: map[string]interface{}{
			"new_exp": arg.Exp,
			"delta":   delta,
		},
	})
	return nil
}

// SetRank is.
func (s *Service) SetRank(ctx context.Context, arg *model.ArgRankSet) error {
	if err := s.memberRPC.SetRank(ctx,
		&member.ArgUpdateRank{
			Mid:      arg.Mid,
			Rank:     arg.Rank,
			RemoteIP: arg.IP}); err != nil {
		return errors.WithStack(err)
	}

	report.Manager(&report.ManagerInfo{
		Uname:    arg.Operator,
		UID:      arg.OperatorID,
		Business: model.ManagerLogID,
		Type:     0,
		Oid:      arg.Mid,
		Action:   _logActionRankUpdate,
		Ctime:    time.Now(),
		// extra
		Index: []interface{}{0, 0, 0, "", "", ""},
		Content: map[string]interface{}{
			"new_rank": arg.Rank,
			"reason":   arg.Reason,
		},
	})
	return nil
}

// SetCoin is.
func (s *Service) SetCoin(ctx context.Context, arg *model.ArgCoinSet) error {
	coins, err := s.coinRPC.UserCoins(ctx, &coin.ArgCoinInfo{Mid: arg.Mid, RealIP: arg.IP})
	if err != nil {
		return errors.WithStack(err)
	}

	reason := "系统操作"
	if arg.Reason != "" {
		reason = fmt.Sprintf("系统操作：%s", arg.Reason)
	}
	delta := arg.Coins - coins
	if _, err = s.coinRPC.ModifyCoin(ctx,
		&coin.ArgModifyCoin{
			Mid:      arg.Mid,
			Operator: arg.Operator,
			Reason:   reason,
			Count:    delta,
			IP:       arg.IP}); err != nil {
		return errors.WithStack(err)
	}

	report.Manager(&report.ManagerInfo{
		Uname:    arg.Operator,
		UID:      arg.OperatorID,
		Business: model.ManagerLogID,
		Type:     0,
		Oid:      arg.Mid,
		Action:   _logActionCoinUpdate,
		Ctime:    time.Now(),
		// extra
		Index: []interface{}{0, 0, 0, "", "", ""},
		Content: map[string]interface{}{
			"new_coins": arg.Coins,
			"delta":     delta,
			"reason":    reason,
		},
	})
	return nil
}

// SetAdditRemark is.
func (s *Service) SetAdditRemark(ctx context.Context, arg *model.ArgAdditRemarkSet) error {
	return s.dao.UpAdditRemark(ctx, arg.Mid, arg.Remark)
}

// PubExpMsg is.
func (s *Service) PubExpMsg(ctx context.Context, arg *model.ArgPubExpMsg) (err error) {
	msg := &model.AddExpMsg{
		Event: arg.Event,
		Mid:   arg.Mid,
		IP:    arg.IP,
		Ts:    arg.Ts,
	}
	return s.dao.PubExpMsg(ctx, msg)
}

// ExpLog is.
func (s *Service) ExpLog(ctx context.Context, mid int64) ([]*model.UserLog, error) {
	return nil, ecode.MethodNotAllowed
}

func filterByStatus(status ...int8) func(*model.FaceRecord) bool {
	ss := make(map[int8]struct{}, len(status))
	for _, s := range status {
		ss[s] = struct{}{}
	}
	return func(fr *model.FaceRecord) bool {
		_, ok := ss[fr.Status]
		return ok
	}
}

func filterByMid(mid int64) func(*model.FaceRecord) bool {
	return func(fr *model.FaceRecord) bool {
		return fr.Mid == mid
	}
}

func filterByOP(operator string) func(*model.FaceRecord) bool {
	return func(fr *model.FaceRecord) bool {
		return fr.Operator == operator
	}
}

// FaceHistory is.
func (s *Service) FaceHistory(ctx context.Context, arg *model.ArgFaceHistory) (*model.FaceRecordPagination, error) {
	list, err := s.faceHistory(ctx, arg)
	if err != nil {
		return nil, err
	}
	plist := list.Paginate(arg.PS*(arg.PN-1), arg.PS)

	page := &model.FaceRecordPagination{
		Records: plist,
		CommonPagination: &model.CommonPagination{
			Page: model.Page{
				Num:   arg.PN,
				Size:  arg.PS,
				Total: len(list),
			},
		},
	}
	return page, nil
}

func (s *Service) faceHistory(ctx context.Context, arg *model.ArgFaceHistory) (res model.FaceRecordList, err error) {
	switch arg.Mode() {
	case "op":
		res, err = s.dao.FaceHistoryByOP(ctx, arg)
		if err != nil {
			return nil, err
		}
		if arg.Mid > 0 {
			res = res.Filter(filterByMid(arg.Mid))
		}
	case "mid":
		res, err = s.dao.FaceHistoryByMid(ctx, arg)
		if err != nil {
			return nil, err
		}
		if arg.Operator != "" {
			res = res.Filter(filterByOP(arg.Operator))
		}
	}
	res = res.Filter(filterByStatus(arg.Status...))
	sort.Slice(res, func(i, j int) bool {
		return res[i].ModifyTime > res[j].ModifyTime
	})
	for _, r := range res {
		r.BuildFaceURL()
	}
	return
}

// MoralLog is.
func (s *Service) MoralLog(ctx context.Context, mid int64) ([]*model.UserLog, error) {
	return nil, ecode.MethodNotAllowed
}

func (s *Service) reviewAddit(ctx context.Context, mids []int64, mrs []*model.BaseReview) error {
	uas, err := s.dao.BatchUserAddit(ctx, mids)
	if err != nil {
		return err
	}
	for _, mr := range mrs {
		if ua, ok := uas[mr.Mid]; ok {
			mr.Addit = *ua
		}
	}
	return nil
}

// BatchFormal is
func (s *Service) BatchFormal(ctx context.Context, arg *model.ArgBatchFormal) error {
	fp := bom.NewReader(bytes.NewReader(arg.FileData))
	reader := csv.NewReader(fp)
	ip := metadata.String(ctx, metadata.RemoteIP)

	columns, err := reader.Read()
	if err != nil {
		log.Error("Failed to read columns from csv: %+v", err)
		return ecode.RequestErr
	}

	findMidPositon := func() (int, error) {
		for i, col := range columns {
			if col == "mid" {
				return i, nil
			}
		}
		return 0, errors.New("No mid column")
	}

	midPosition, err := findMidPositon()
	if err != nil {
		log.Error("Failed to find mid column: %+v", err)
		return ecode.RequestErr
	}

	mids := make([]int64, 0)
	for {
		record, rerr := reader.Read()
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			log.Error("Failed to parse csv: %+v", errors.WithStack(rerr))
			return ecode.RequestErr
		}

		if len(record) < midPosition {
			log.Warn("Skip record due to no suitable position: %+v", record)
			continue
		}

		mid, perr := strconv.ParseInt(record[midPosition], 10, 64)
		if perr != nil {
			log.Warn("Failed to parse mid on data: %+v: %+v", record[midPosition], perr)
			continue
		}

		mids = append(mids, mid)
	}

	bases, err := s.batchBase(ctx, mids)
	if err != nil {
		log.Error("Failed to query bases with mids: %+v: %+v", mids, err)
		return ecode.RequestErr
	}

	for _, mid := range mids {
		base, ok := bases[mid]
		if !ok {
			log.Warn("No such user with mid: %d", mid)
			continue
		}
		if base.Rank >= 10000 {
			log.Warn("Rank already exceeded 10000 on mid: %d: %+v", mid, base)
			continue
		}

		rankArg := &model.ArgRankSet{
			Mid:        mid,
			Rank:       10000,
			Operator:   arg.Operator,
			OperatorID: arg.OperatorID,
			IP:         ip,
		}
		if err := s.SetRank(ctx, rankArg); err != nil {
			log.Warn("Failed to set rank with mid: %d: %+v", mid, err)
			continue
		}

		// 通过发放一次每日登录的经验奖励消息来使用户等级直升 lv1
		expArg := &model.ArgPubExpMsg{
			Mid:   mid,
			IP:    ip,
			Ts:    time.Now().Unix(),
			Event: "login",
		}
		if err := s.PubExpMsg(ctx, expArg); err != nil {
			log.Warn("Failed to pub exp message with mid: %d: %+v", mid, err)
			continue
		}
	}

	return nil
}
