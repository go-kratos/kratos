package extern

import "context"

type MockExternHandler struct {
	ErrDeleteReply error
}

func (mf *MockExternHandler) DeleteReply(ctx context.Context, adminID int64, ks []*Reply) error {
	return mf.ErrDeleteReply
}
