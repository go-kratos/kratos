package biz

import (
	"context"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewUserUsecase)

type Transaction interface {
	Transaction(context.Context, func(ctx context.Context) error) error
}
