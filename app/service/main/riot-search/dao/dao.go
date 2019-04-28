package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-common/app/service/main/riot-search/conf"
	"go-common/app/service/main/riot-search/model"
	"go-common/library/database/sql"
	"go-common/library/log"

	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
)

// 过审的增量数据
var _selIncrement = "SELECT id, title from archive where mtime>? and mtime<=?"

// Dao dao
type Dao struct {
	c        *conf.Config
	searcher *riot.Engine
	db       *sql.DB
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	validateConfig(c.Riot)
	if c.UT {
		dao = &Dao{
			c:        c,
			searcher: &riot.Engine{},
			db:       sql.NewMySQL(c.Mysql),
		}
		dao.searcher.Init(types.EngineOpts{})
		return
	}
	dao = &Dao{
		c:        c,
		searcher: &riot.Engine{},
		// db
		db: sql.NewMySQL(c.Mysql),
	}
	dao.searcher.Init(types.EngineOpts{
		GseDict:       c.Riot.Dict,
		StopTokenFile: c.Riot.StopToken,
		NumShards:     c.Riot.NumShards,
		IndexerOpts: &types.IndexerOpts{
			IndexType:    types.FrequenciesIndex,
			DocCacheSize: 5000,
		},
	})
	return
}

func validateConfig(conf *conf.RiotConfig) {
	if conf.Dict == "" || conf.StopToken == "" {
		panic("must provide a dict and stop_token file")
	}
	if conf.FlushTime <= 0 {
		panic("flush time must larger than 0")
	}
}

// Close close the resource.
func (d *Dao) Close() {
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	return d.db.Ping(c)
}

// IncrementBackup select mtime>now-24h data
func (d *Dao) IncrementBackup(c context.Context, stime, etime time.Time) (docs []*model.Document, err error) {
	var states []int
	for k := range model.PubStates.LegalStates {
		states = append(states, k)
	}
	query := _selIncrement + " and state in (" + strings.Trim(strings.Join(strings.Split(fmt.Sprint(states), " "), ","), "[]") + ")" + " order by id asc"
	rows, err := d.db.Query(c, query, stime, etime)
	log.Info("exec query(%s) args(stime:%v, etime:%v)", query, stime, etime)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		doc := &model.Document{}
		if err = rows.Scan(&doc.ID, &doc.Content); err != nil {
			return
		}
		docs = append(docs, doc)
	}
	err = rows.Err()
	return
}
