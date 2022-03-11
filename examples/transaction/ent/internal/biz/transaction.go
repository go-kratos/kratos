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
	CreateUser(ctx context.Context, a *User) (int, error)
}

type CardRepo interface {
	CreateCard(ctx context.Context, id int) (int, error)
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
		id  int
	)
	if e := u.tm.ExecTx(ctx, func(ctx context.Context) error {
		id, err = u.userRepo.CreateUser(ctx, m)
		if err != nil {
			return err
		}
		if _, e := u.cardRepo.CreateCard(ctx, id); err != nil {
			return e
		}
		return nil
	}); e != nil {
		return 0, err
	}
	return id, nil
}
