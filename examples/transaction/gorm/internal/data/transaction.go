package data

import (
	"context"
	"time"

	"github.com/SeeMusic/kratos/examples/transaction/gorm/internal/biz"
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

type User struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Card struct {
	ID        int64
	UserID    int64
	Money     int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (u *userRepo) CreateUser(ctx context.Context, m *biz.User) (int64, error) {
	user := User{Name: m.Name, Email: m.Email}
	result := u.data.DB(ctx).Create(&user)
	return user.ID, result.Error
}

func NewCardRepo(data *Data, logger log.Logger) biz.CardRepo {
	return &cardRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (c *cardRepo) CreateCard(ctx context.Context, id int64) (int64, error) {
	var card Card
	card.UserID = id
	card.Money = 1000
	result := c.data.DB(ctx).Save(&card)
	return card.ID, result.Error
}
