package tab

import (
	"context"

	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-feed/conf"
	"go-common/library/database/sql"
)

const (
	_followSQL = "SELECT `id`,`type`,`long_title`,`content` FROM `card_follow` WHERE `deleted`=0"
)

type Dao struct {
	db        *sql.DB
	followGet *sql.Stmt
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: sql.NewMySQL(c.MySQL.Show),
	}
	// prepare
	d.followGet = d.db.Prepared(_followSQL)
	return
}

func (d *Dao) Follow(c context.Context) (cm map[int64]*operate.Follow, err error) {
	var rows *sql.Rows
	if rows, err = d.followGet.Query(c); err != nil {
		return
	}
	defer rows.Close()
	cm = make(map[int64]*operate.Follow)
	for rows.Next() {
		c := &operate.Follow{}
		if err = rows.Scan(&c.ID, &c.Type, &c.Title, &c.Content); err != nil {
			return
		}
		c.Change()
		cm[c.ID] = c
	}
	return
}
