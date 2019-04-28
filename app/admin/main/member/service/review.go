package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/member/model"
	"go-common/app/admin/main/member/model/block"
	relation "go-common/app/service/main/relation/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/queue/databus/report"
	xtime "go-common/library/time"
)

const (
	_logActionUserPropertyAudit = "user_property_review_audit"
)

// Reviews is.
func (s *Service) Reviews(ctx context.Context, arg *model.ArgReviewList) ([]*model.UserPropertyReview, int, error) {
	bySearch := func() ([]*model.UserPropertyReview, int, error) {
		stime := arg.STime.Time().Format("2006-01-02 15:04:05")
		if arg.ETime == 0 {
			arg.ETime = xtime.Time(time.Now().Unix())
		}
		etime := arg.ETime.Time().Format("2006-01-02 15:04:05")
		property := int8ToInt(arg.Property)
		state := int8ToInt(arg.State)
		result, err := s.dao.SearchUserPropertyReview(ctx, arg.Mid, property, state, arg.IsMonitor, arg.IsDesc, arg.Operator, stime, etime, arg.Pn, arg.Ps)
		if err != nil {
			return nil, 0, err
		}
		ids := result.IDs()

		rws, err := s.dao.ReviewByIDs(ctx, ids, arg.State)
		if err != nil {
			return nil, 0, err
		}
		rws = arrange(rws, ids)

		return rws, result.Total(), nil
	}
	byDB := func() ([]*model.UserPropertyReview, int, error) {
		return s.dao.Reviews(ctx, arg.Mid, arg.Property, arg.State, arg.IsMonitor, arg.IsDesc, arg.Operator, arg.STime, arg.ETime, arg.Pn, arg.Ps)
	}

	rws, total, err := bySearch()
	if arg.ForceDB {
		log.Info("Force user property review query to db")
		rws, total, err = byDB()
	}
	if err != nil {
		return nil, 0, err
	}

	for _, rw := range rws {
		if rw.Property == model.ReviewPropertyFace {
			rw.BuildFaceURL()
		}
	}
	s.reviewsName(ctx, rws)
	s.reviewsFaceReject(ctx, rws)
	s.reviewsRelationStat(ctx, rws)
	return rws, total, err
}

func (s *Service) onReviewSuccess(ctx context.Context, waitRws []*model.UserPropertyReview, arg *model.ArgReviewAudit) error {
	if !arg.BlockUser {
		return nil
	}
	blockArg := &block.ParamBatchBlock{
		MIDs:      waitRwMids(waitRws),
		AdminName: arg.Operator,
		AdminID:   arg.OperatorID,
		Source:    arg.Source,
		Area:      arg.Area,
		Reason:    arg.Reason,
		Comment:   arg.Comment,
		Action:    arg.Action,
		Duration:  arg.Duration,
		Notify:    arg.Notify,
	}
	if !blockArg.Validate() {
		log.Error("Failed to validate block parama, arg: %v", blockArg)
		return ecode.RequestErr
	}
	if err := s.block.BatchBlock(ctx, blockArg); err != nil {
		log.Error("Failed to batch block, error: %v, arg: %v", err, blockArg)
		return err
	}
	return nil
}

// ReviewAudit is.
func (s *Service) ReviewAudit(ctx context.Context, arg *model.ArgReviewAudit) error {
	waitRws, err := s.dao.ReviewByIDs(ctx, arg.ID, []int8{model.ReviewStateWait})
	if err != nil {
		return err
	}
	if err := s.dao.ReviewAudit(ctx, arg.ID, arg.State, arg.Remark, arg.Operator); err != nil {
		return err
	}
	for _, r := range waitRws {
		ak := auditKey(r.Property, r.IsMonitor)
		handler, ok := s.auditHandlers[ak]
		if !ok {
			log.Warn("Unable to handle property update: review: %+v audit: %+v", r, arg)
			continue
		}
		if err := handler(ctx, r, arg); err != nil {
			log.Error("Failed to trigger review audit event: review: %+v error: %+v", r, err)
			remark := fmt.Sprintf("操作异常：%s, 备注: %s", ecode.Cause(err).Message(), arg.Remark)
			if err = s.dao.UpdateRemark(ctx, r.ID, remark); err != nil {
				log.Error("Failed to update remark error: %v", err)
			}
		}
		report.Manager(&report.ManagerInfo{
			Uname:    arg.Operator,
			UID:      arg.OperatorID,
			Business: model.ManagerLogID,
			Type:     0,
			Oid:      r.Mid,
			Action:   _logActionUserPropertyAudit,
			Ctime:    time.Now(),
			// extra
			Index: []interface{}{r.ID, arg.State, 0, arg.Remark, "", ""},
			Content: map[string]interface{}{
				"remark": arg.Remark,
				"state":  arg.State,
				"id":     r.ID,
				"mid":    r.Mid,
			},
		})
	}
	s.onReviewSuccess(ctx, waitRws, arg)
	return nil
}

// Review is.
func (s *Service) Review(ctx context.Context, arg *model.ArgReview) (*model.UserPropertyReview, error) {
	r, err := s.dao.Review(ctx, arg.ID)
	if err != nil {
		return nil, err
	}
	r.Block, err = s.block.History(ctx, r.Mid, 2, 1, true)
	if err != nil {
		log.Error("Failed to get block review info, error: %v, mid: %v", err, r.Mid)
		err = nil
	}
	return r, nil
}

func (s *Service) reviewsName(ctx context.Context, rws []*model.UserPropertyReview) {
	mids := make([]int64, 0, len(rws))
	for _, rw := range rws {
		mids = append(mids, rw.Mid)
	}
	bs, err := s.dao.Bases(ctx, mids)
	if err != nil {
		log.Error("Failed to fetch bases with mids: %+v: %+v", mids, err)
		return
	}
	for _, rw := range rws {
		b, ok := bs[rw.Mid]
		if !ok {
			continue
		}
		rw.Name = b.Name
	}
}

func (s *Service) reviewsFaceReject(ctx context.Context, rws []*model.UserPropertyReview) {
	mids := make([]int64, 0, len(rws))
	for _, rw := range rws {
		mids = append(mids, rw.Mid)
	}
	frs, err := s.dao.BatchUserAddit(ctx, mids)
	if err != nil {
		log.Error("Failed to fetch FaceRejects with mids: %+v: %+v", mids, err)
		return
	}
	for _, rw := range rws {
		if fr, ok := frs[rw.Mid]; ok {
			rw.FaceReject = fr.FaceReject
		}
	}
}

func (s *Service) reviewsRelationStat(ctx context.Context, rws []*model.UserPropertyReview) {
	mids := make([]int64, 0, len(rws))
	for _, rw := range rws {
		mids = append(mids, rw.Mid)
	}
	stats, err := s.relationRPC.Stats(ctx, &relation.ArgMids{
		Mids:   mids,
		RealIP: metadata.String(ctx, metadata.RemoteIP),
	})
	if err != nil {
		log.Error("Failed to fetch relation stat with mids: %+v: %+v", mids, err)
		return
	}
	for _, rw := range rws {
		stat, ok := stats[rw.Mid]
		if !ok {
			continue
		}
		rw.Follower = stat.Follower
	}
}

func int8ToInt(in []int8) []int {
	res := []int{}
	for _, i := range in {
		res = append(res, int(i))
	}
	return res
}

func arrange(rws []*model.UserPropertyReview, ids []int64) []*model.UserPropertyReview {
	res := []*model.UserPropertyReview{}
	tmp := make(map[int64]*model.UserPropertyReview, len(ids))
	for _, rw := range rws {
		tmp[rw.ID] = rw
	}
	for _, id := range ids {
		if rw, ok := tmp[id]; ok {
			res = append(res, rw)
		}
	}
	return res
}

func waitRwMids(waitRws []*model.UserPropertyReview) []int64 {
	mids := make([]int64, 0, len(waitRws))
	for _, w := range waitRws {
		mids = append(mids, w.Mid)
	}
	return mids
}
