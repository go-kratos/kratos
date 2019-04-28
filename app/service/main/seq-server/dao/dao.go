package dao

import (
	"context"

	"go-common/app/service/main/seq-server/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// Dao is seq-server dao.
type Dao struct {
	c *conf.Config
	// db
	db *sql.DB
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// db
		db: sql.NewMySQL(c.DB.Number),
	}
	return
}

// Ping ping db
func (d *Dao) Ping() (err error) {
	if err = d.db.Ping(context.TODO()); err != nil {
		log.Error("d.db.Ping error(%v)", err)
	}
	return
}

// Close close resource.
func (d *Dao) Close() {
	d.db.Close()
}
