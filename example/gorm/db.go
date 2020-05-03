package dao

import (
	"github.com/go-kratos/kratos/databases/orm"
	"github.com/jinzhu/gorm"

	"github.com/go-kratos/kratos/pkg/conf/paladin"
)

func NewDB() (db *gorm.DB, cf func(), err error) {
	var (
		cfg orm.Config
		ct  paladin.TOML
	)
	if err = paladin.Get("db.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Client").UnmarshalTOML(&cfg); err != nil {
		return
	}
	db = orm.NewMySQL(&cfg)
	db.LogMode(true)
	cf = func() { db.Close() }
	return
}
