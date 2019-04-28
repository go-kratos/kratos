package archive

import (
	"context"
)

// StaffApplySubmit fn
func (s *Service) StaffApplySubmit(c context.Context, id, aid, mid, state, atype int64, flagAddBlack, flagRefuse int) (err error) {
	return s.arc.StaffApplySubmit(c, id, aid, mid, state, atype, flagAddBlack, flagRefuse)
}

// StaffValidate func
func (s *Service) StaffValidate(c context.Context, mid int64) (uv int, err error) {
	return s.arc.StaffMidValidate(c, mid)
}
