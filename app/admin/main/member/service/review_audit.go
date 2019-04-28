package service

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"path"

	"go-common/app/admin/main/member/model"
	comodel "go-common/app/service/main/coin/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

var upNameCostCoins = 6.0

type auditHandler func(ctx context.Context, origin *model.UserPropertyReview, arg *model.ArgReviewAudit) error

func (s *Service) initAuditHandler() {
	s.auditHandlers[auditKey(model.ReviewPropertySign, true)] = s.onSignAudit
	s.auditHandlers[auditKey(model.ReviewPropertyName, true)] = s.onNameAudit
	s.auditHandlers[auditKey(model.ReviewPropertyFace, true)] = s.onFaceMonitorAudit
	s.auditHandlers[auditKey(model.ReviewPropertyFace, false)] = s.onFaceAudit
}

func auditKey(property int8, isMonitor bool) string {
	return fmt.Sprintf("%d-%t", property, isMonitor)
}

func (s *Service) onSignAudit(ctx context.Context, origin *model.UserPropertyReview, arg *model.ArgReviewAudit) error {
	switch arg.State {
	case model.ReviewStatePass:
		return s.dao.UpSign(ctx, origin.Mid, origin.New)
	case model.ReviewStateNoPass:
		if err := s.dao.Message(ctx, "违规签名处理通知", "抱歉，由于你的签名涉嫌违规，系统已将你的签名回退。如有疑问请联系客服。", []int64{origin.Mid}); err != nil {
			log.Error("Failed to send message: mid: %d: %+v", origin.Mid, err)
		}
		return nil
	}
	return nil
}

func (s *Service) onNameAudit(ctx context.Context, origin *model.UserPropertyReview, arg *model.ArgReviewAudit) error {
	switch arg.State {
	case model.ReviewStatePass:
		if origin.NickFree() {
			return s.dao.UpdateUname(ctx, origin.Mid, origin.New)
		}
		coins, err := s.coinRPC.UserCoins(ctx, &comodel.ArgCoinInfo{Mid: origin.Mid})
		if err != nil {
			return err
		}
		if coins < upNameCostCoins {
			return ecode.UpdateUnameMoneyIsNot
		}
		if err := s.dao.UpdateUname(ctx, origin.Mid, origin.New); err != nil {
			log.Error("faild to update Name, mid: %d, name: %s, %+v", origin.Mid, origin.New, err)
			return err
		}
		arg := &comodel.ArgModifyCoin{
			Mid:    origin.Mid,
			Count:  -upNameCostCoins,
			Reason: fmt.Sprintf("UPDATE:NICK:%s=>%s", origin.Old, origin.New),
		}
		if _, err := s.coinRPC.ModifyCoin(ctx, arg); err != nil {
			log.Error("faild to modify coin, arg: %+v, %+v", arg, err)
			return err
		}

	case model.ReviewStateNoPass:
		if err := s.dao.Message(ctx, "违规昵称处理通知", "抱歉，由于你的昵称涉嫌违规，系统已将你的昵称回退。如有疑问请联系客服。", []int64{origin.Mid}); err != nil {
			log.Error("Failed to send message: mid: %d: %+v", origin.Mid, err)
		}
		return nil
	}
	return nil
}

func (s *Service) onFaceAudit(ctx context.Context, origin *model.UserPropertyReview, arg *model.ArgReviewAudit) error {
	switch arg.State {
	case model.ReviewStateNoPass:
		if err := s.dao.UpFace(ctx, origin.Mid, ""); err != nil {
			return err
		}
		return s.faceReject(ctx, origin, arg)
	}
	return nil
}

func (s *Service) onFaceMonitorAudit(ctx context.Context, origin *model.UserPropertyReview, arg *model.ArgReviewAudit) error {
	switch arg.State {
	case model.ReviewStatePass:
		return s.dao.UpFace(ctx, origin.Mid, origin.New)
	case model.ReviewStateNoPass:
		return s.faceReject(ctx, origin, arg)
	}
	return nil
}

func (s *Service) mvToPrivate(ctx context.Context, face string) (string, error) {
	file, err := s.dao.Image(buildURL(face))
	if err != nil {
		return "", err
	}
	ftype := http.DetectContentType(file)
	privURL, err := s.dao.UploadImage(ctx, ftype, file, s.c.FacePriBFS)
	if err != nil {
		return "", err
	}
	if err := s.dao.DelImage(ctx, path.Base(privURL), s.c.FaceBFS); err != nil {
		log.Error("s.dao.DelImage(%v) error(%+v)", privURL, err)
	}
	return urlPath(privURL), nil
}

func buildURL(path string) string {
	return fmt.Sprintf("http://i%d.hdslb.com%s", rand.Int63n(3), path)
}

func urlPath(in string) string {
	URL, err := url.Parse(in)
	if err != nil {
		return ""
	}
	return URL.Path
}

func (s *Service) faceReject(ctx context.Context, origin *model.UserPropertyReview, arg *model.ArgReviewAudit) error {
	privFace, err := s.mvToPrivate(ctx, origin.New)
	if err != nil {
		log.Error("s.mvToPrivate(%d) error(%+v)", origin.Mid, err)
		return err
	}
	if err := s.dao.UpdateReviewFace(ctx, origin.ID, privFace); err != nil {
		log.Error("Failed to update review face: id: %d face: %s: %+v", origin.ID, privFace, err)
		return err
	}
	if err := s.dao.MvArchivedFaceToPriv(ctx, origin.New, privFace, arg.Operator, arg.Remark); err != nil {
		log.Error("mv archived face to private bucket mid(%d) error(%v)", origin.Mid, err)
		return err
	}
	if err := s.dao.Message(ctx, "违规头像处理通知", "抱歉，由于你的头像涉嫌违规，系统已将你的头像回退。如有疑问请联系客服。", []int64{origin.Mid}); err != nil {
		log.Error("Failed to send message: mid: %d: %+v", origin.Mid, err)
	}
	if err := s.dao.IncrFaceReject(ctx, origin.Mid); err != nil {
		log.Error("IncrFaceReject faild: mid: %d: %+v", origin.Mid, err)
	}
	return nil
}
