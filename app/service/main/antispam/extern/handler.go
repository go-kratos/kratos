package extern

import "context"

type Handler interface {
	ReplyHandler
}

type ReplyHandler interface {
	DeleteReply(ctx context.Context, adminId int64, rs []*Reply) error
}
