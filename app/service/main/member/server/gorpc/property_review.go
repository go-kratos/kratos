package gorpc

import (
	"go-common/app/service/main/member/model"
	"go-common/library/net/rpc/context"
)

// AddUserMonitor is add user into monitor
func (r *RPC) AddUserMonitor(ctx context.Context, arg *model.ArgAddUserMonitor, res *struct{}) error {
	return r.s.AddUserMonitor(ctx, arg)
}

// IsInMonitor check user is in monitor
func (r *RPC) IsInMonitor(ctx context.Context, arg *model.ArgMid, res *bool) error {
	isInMonitor, err := r.s.IsInMonitor(ctx, arg)
	if err != nil {
		return err
	}
	*res = isInMonitor
	return nil
}

// AddPropertyReview add user property update review.
func (r *RPC) AddPropertyReview(ctx context.Context, arg *model.ArgAddPropertyReview, res *struct{}) error {
	return r.s.AddPropertyReview(ctx, arg)
}
