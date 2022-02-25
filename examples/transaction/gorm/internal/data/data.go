package data

import (
	"context"

	"github.com/go-kratos/kratos/examples/transaction/gorm/internal/biz"
	"github.com/go-kratos/kratos/examples/transaction/gorm/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	// init mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewTransaction, NewUserRepo, NewCardRepo)

// Data .
type Data struct {
	db *gorm.DB
}

type contextTxKey struct{}

func (d *Data) ExecTx(ctx context.Context, f func(ctx context.Context) error) error {
	tx := d.db.Begin()
	ctx = context.WithValue(ctx, contextTxKey{}, tx)
	if err := f(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func NewTransaction(d *Data) biz.Transaction {
	return d
}

// NewData .
func NewData(conf *conf.Data, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(logger)

	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{})
	if err != nil {
		log.Errorf("failed opening connection to sqlite: %v", err)
		return nil, nil, err
	}

	d := &Data{
		db: db,
	}
	db.AutoMigrate(new(biz.User), new(biz.Card))

	return d, func() {
		log.Info("message", "closing the data resources")
		close, err := db.DB()
		if err != nil {
			log.Error(err)
		}
		if err := close.Close(); err != nil {
			log.Error(err)
		}
	}, nil
}
