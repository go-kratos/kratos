package service

import (
	"context"
	"strconv"

	pb "github.com/SeeMusic/kratos/examples/transaction/api/transaction/v1"
	"github.com/SeeMusic/kratos/examples/transaction/gorm/internal/biz"

	"github.com/SeeMusic/kratos/v2/log"
)

func NewTransactionService(user *biz.UserUsecase, logger log.Logger) *TransactionService {
	return &TransactionService{
		user: user,
		log:  log.NewHelper(logger),
	}
}

func (b *TransactionService) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	id, err := b.user.CreateUser(ctx, &biz.User{
		Name:  in.Name,
		Email: in.Email,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserReply{
		Id: strconv.Itoa(id),
	}, nil
}
