package data

import (
	"context"

	"github.com/SeeMusic/kratos/examples/transaction/ent/internal/biz"
	"github.com/SeeMusic/kratos/examples/transaction/ent/internal/conf"
	"github.com/SeeMusic/kratos/examples/transaction/ent/internal/data/ent"

	"entgo.io/ent/dialect/sql"
	"github.com/SeeMusic/kratos/v2/log"
	"github.com/google/wire"

	// init mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewTransaction, NewUserRepo, NewCardRepo)

// Data .
type Data struct {
	db *ent.Client
}

type contextTxKey struct{}

func (d *Data) ExecTx(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := d.db.Tx(ctx)
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, contextTxKey{}, tx)
	if err := f(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (d *Data) User(ctx context.Context) *ent.UserClient {
	tx, ok := ctx.Value(contextTxKey{}).(*ent.Tx)
	if ok {
		return tx.User
	}
	return d.db.User
}

func (d *Data) Card(ctx context.Context) *ent.CardClient {
	tx, ok := ctx.Value(contextTxKey{}).(*ent.Tx)
	if ok {
		return tx.Card
	}
	return d.db.Card
}

func NewTransaction(d *Data) biz.Transaction {
	return d
}

// NewData .
func NewData(conf *conf.Data, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(logger)
	drv, err := sql.Open(
		conf.Database.Driver,
		conf.Database.Source,
	)
	client := ent.NewClient(ent.Driver(drv))
	if err != nil {
		log.Errorf("failed opening connection to sqlite: %v", err)
		return nil, nil, err
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Errorf("failed creating schema resources: %v", err)
		return nil, nil, err
	}

	d := &Data{
		db: client,
	}
	return d, func() {
		log.Info("message", "closing the data resources")
		if err := d.db.Close(); err != nil {
			log.Error(err)
		}
	}, nil
}
