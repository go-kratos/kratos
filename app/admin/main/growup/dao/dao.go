package dao

import (
	"context"

	"go-common/app/admin/main/growup/conf"
	"go-common/library/database/orm"
	"go-common/library/database/sql"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

// Dao dao
type Dao struct {
	c         *conf.Config
	db        *gorm.DB
	rddb      *sql.DB
	client    *httpx.Client
	VideoURL  string
	ColumnURL string
}

// New fn
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:         c,
		db:        orm.NewMySQL(c.ORM.Growup),
		rddb:      sql.NewMySQL(c.ORM.Allowance),
		client:    httpx.NewClient(c.HTTPClient),
		VideoURL:  c.Host.VideoType + "/videoup/types",
		ColumnURL: c.Host.ColumnType + "/x/article/categories",
	}
	d.initORM()

	return
}

func (d *Dao) initORM() {
	d.db.LogMode(true)
}

// Ping check conn of db
func (d *Dao) Ping(c context.Context) (err error) {
	if d.db != nil {
		err = d.db.DB().PingContext(c)
	}
	return
}

// Close close conn of db
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

// BeginTran begin tran
func (d *Dao) BeginTran(c context.Context) (tx *sql.Tx, err error) {
	return d.rddb.Begin(c)
}

// Exec do exec
func (d *Dao) Exec(c context.Context, sql string) (rows int64, err error) {
	res, err := d.rddb.Exec(c, sql)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// DoInTx .
func (d *Dao) DoInTx(c context.Context, txFunc func(*sql.Tx) error) (err error) {
	tx, err := d.rddb.Begin(c)
	if err != nil {
		log.Error("d.rddb.Begin err(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback err(%v)", err1)
			}
			return
		}

	}()
	err = txFunc(tx)
	if err != nil {
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit err(%v)", err)
	}
	return
}
