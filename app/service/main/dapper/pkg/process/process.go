package process

import (
	"context"

	"go-common/app/service/main/dapper/model"
)

// Processer .
type Processer interface {
	Process(ctx context.Context, protoSpan *model.ProtoSpan) error
}

// MockProcess MockProcess
type MockProcess func(ctx context.Context, protoSpan *model.ProtoSpan) error

// Process implement Processer
func (m MockProcess) Process(ctx context.Context, protoSpan *model.ProtoSpan) error {
	return m(ctx, protoSpan)
}
