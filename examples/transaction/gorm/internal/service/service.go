package service

import (
	pb "github.com/SeeMusic/kratos/examples/transaction/api/transaction/v1"
	"github.com/SeeMusic/kratos/examples/transaction/gorm/internal/biz"

	"github.com/SeeMusic/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewTransactionService)

type TransactionService struct {
	pb.UnimplementedTransactionServiceServer

	log *log.Helper

	user *biz.UserUsecase
}
