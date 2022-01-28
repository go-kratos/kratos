package events

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport"
)

type Handler interface {
	Handle(ctx context.Context, msg Message) error
}

type Message struct {
	Topic    string
	Data     []byte
	Metadata map[string]interface{}
}

type PublishMetadata struct {
	Topic string
}

type SubRequest struct {
	Topic string
}

type Subscriber interface {
	transport.Server
	Subscribe(subReq SubRequest, handler Handler) error
}
