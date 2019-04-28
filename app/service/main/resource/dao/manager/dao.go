package manager

import (
	"go-common/app/service/main/resource/conf"
	"go-common/library/database/sql"
)

//Dao manager dao
type Dao struct {
	db *sql.DB
	c  *conf.Config
}

//New new manager dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:  c,
		db: sql.NewMySQL(c.DB.Manager),
	}
	return
}

// Close close db resource.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
