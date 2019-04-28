package service

import (
	"context"
	"go-common/app/admin/main/videoup/model/archive"
)

// Staffs fn
func (s *Service) Staffs(c context.Context, aid int64) (data []*archive.Staff, err error) {
	return s.staff.Staffs(c, aid)
}

// StaffApplyBatchSubmit func
func (s *Service) StaffApplyBatchSubmit(c context.Context, ap *archive.StaffBatchParam) (err error) {
	return s.staff.StaffApplyBatchSubmit(c, ap)
}
