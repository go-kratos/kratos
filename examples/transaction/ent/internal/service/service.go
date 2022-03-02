package service

import (
	pb "github.com/go-kratos/kratos/examples/transaction/api/transaction/v1"
	"github.com/go-kratos/kratos/examples/transaction/ent/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewTransactionService)

type TransactionService struct {
	pb.UnimplementedTransactionServiceServer

	log *log.Helper

	user *biz.UserUsecase
}
