package biz

import (
	"context"

	"github.com/SeeMusic/kratos/v2/log"
)

type User struct {
	Name  string
	Email string
}

type UserRepo interface {
	CreateUser(ctx context.Context, a *User) (int64, error)
}

type CardRepo interface {
	CreateCard(ctx context.Context, id int64) (int64, error)
}

type UserUsecase struct {
	userRepo UserRepo
	cardRepo CardRepo
	tm       Transaction
}

func NewUserUsecase(user UserRepo, card CardRepo, tm Transaction, logger log.Logger) *UserUsecase {
	return &UserUsecase{userRepo: user, cardRepo: card, tm: tm}
}

func (u *UserUsecase) CreateUser(ctx context.Context, m *User) (int, error) {
	var (
		err error
		id  int64
	)
	err = u.tm.ExecTx(ctx, func(ctx context.Context) error {
		id, err = u.userRepo.CreateUser(ctx, m)
		if err != nil {
			return err
		}
		_, err = u.cardRepo.CreateCard(ctx, id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
