package data

import (
	"context"

	"github.com/go-kratos/kratos/examples/transaction/gorm/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
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
	txx := ctx.Value(contextTxKey{})
	if tx, ok := txx.(*gorm.DB); ok {
		err := tx.WithContext(ctx).Create(m).Error
		return m.ID, err
	}
	err := u.data.db.WithContext(ctx).Create(m).Error
	return m.ID, err
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (c *cardRepo) CreateCard(ctx context.Context, m *biz.Card) (int, error) {
	txx := ctx.Value(contextTxKey{})
	if tx, ok := txx.(*gorm.DB); ok {
		err := tx.WithContext(ctx).Create(m).Error
		return m.ID, err
	}
	err := c.data.db.WithContext(ctx).Create(m).Error
	return m.ID, err
}

func NewCardRepo(data *Data, logger log.Logger) biz.CardRepo {
	return &cardRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
