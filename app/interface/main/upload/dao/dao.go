package dao

import (
	"context"

	"go-common/app/interface/main/upload/conf"
	"go-common/library/database/orm"

	"github.com/jinzhu/gorm"
)

// Dao dao struct
type Dao struct {
	orm *gorm.DB
}

// NewDao new a dao instance.
func NewDao(c *conf.Config) *Dao {
	return &Dao{
		orm: orm.NewMySQL(c.ORM),
	}
}

// Ping ping database.
func (d *Dao) Ping(c context.Context) error {
	return nil
}
