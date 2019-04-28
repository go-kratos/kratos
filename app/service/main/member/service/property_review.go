package service

import (
	"context"
	"net/url"

	"go-common/app/service/main/member/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AddUserMonitor is.
func (s *Service) AddUserMonitor(ctx context.Context, arg *model.ArgAddUserMonitor) error {
	if err := s.mbDao.AddUserMonitor(ctx, arg.Mid, arg.Operator, arg.Remark); err != nil {
		log.Error("Failed to add user monitor with arg: %+v: %+v", arg, err)
		return err
	}
	return nil
}

// IsInMonitor is.
func (s *Service) IsInMonitor(ctx context.Context, arg *model.ArgMid) (bool, error) {
	inMonitor, err := s.mbDao.IsInUserMonitor(ctx, arg.Mid)
	if err != nil {
		log.Error("Failed to check user is in monitor with arg: %+v: %+v", arg, err)
		return false, err
	}
	return inMonitor, nil
}

// AddPropertyReview is
func (s *Service) AddPropertyReview(ctx context.Context, arg *model.ArgAddPropertyReview) error {
	base, err := s.mbDao.BaseInfo(ctx, arg.Mid)
	if err != nil {
		return err
	}

	// prepase the old value
	old := ""
	switch arg.Property {
	case model.ReviewPropertyFace:
		old = path(base.Face)
	case model.ReviewPropertyName:
		old = base.Name
	case model.ReviewPropertySign:
		old = base.Sign
	default:
		return ecode.RequestErr
	}

	// check is in monitor
	isMonitor, err := s.mbDao.IsInUserMonitor(ctx, arg.Mid)
	if err != nil {
		log.Error("Failed to check user is in monitor with mid: %+v: %+v", arg.Mid, err)
	}
	if err = s.mbDao.ArchivePropertyReview(ctx, arg.Mid, arg.Property); err != nil {
		log.Error("Failed to archive property review by mid: %d, property: %d, error: %+v", arg.Mid, arg.Property, err)
	}
	r := &model.UserPropertyReview{
		Mid:       arg.Mid,
		Old:       old,
		New:       arg.New,
		State:     arg.State,
		Property:  arg.Property,
		IsMonitor: isMonitor,
		Extra:     arg.ExtraStr(),
	}
	return s.mbDao.AddPropertyReview(ctx, r)
}

func path(faceURL string) string {
	URL, err := url.Parse(faceURL)
	if err != nil {
		return ""
	}
	return URL.Path
}
