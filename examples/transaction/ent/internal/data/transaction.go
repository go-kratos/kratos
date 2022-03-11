package data

import (
	"context"
	"strconv"

	"github.com/SeeMusic/kratos/examples/transaction/ent/internal/biz"
	"github.com/SeeMusic/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

type cardRepo struct {
	data *Data
	log  *log.Helper
}

func (u *userRepo) CreateUser(ctx context.Context, m *biz.User) (int, error) {
	user, err := u.data.User(ctx).
		Create().
		SetName(m.Name).
		SetEmail(m.Email).
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (c *cardRepo) CreateCard(ctx context.Context, id int) (int, error) {
	card, err := c.data.Card(ctx).
		Create().
		SetMoney("1000").
		SetUserID(strconv.Itoa(id)).
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return card.ID, nil
}

func NewCardRepo(data *Data, logger log.Logger) biz.CardRepo {
	return &cardRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
