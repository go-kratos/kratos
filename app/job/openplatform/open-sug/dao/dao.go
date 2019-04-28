package dao

import (
	"context"
	"net/http"
	"time"

	"go-common/app/job/openplatform/open-sug/conf"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"

	"gopkg.in/olivere/elastic.v5"
)

// Dao dao
type Dao struct {
	c              *conf.Config
	redis          *redis.Pool
	mallDB         *xsql.DB
	ugcDB          *xsql.DB
	ticketDB       *xsql.DB
	client         *http.Client
	ItemSalesMax   map[string]int
	ItemSalesMin   map[string]int
	ItemWishMax    map[string]int
	ItemWishMin    map[string]int
	ItemCommentMax map[string]int
	ItemCommentMin map[string]int
	es             *elastic.Client
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	var (
		err error
		es  *elastic.Client
	)
	dao = &Dao{
		c:              c,
		redis:          redis.NewPool(c.Redis),
		mallDB:         xsql.NewMySQL(c.MallMySQL),
		ugcDB:          xsql.NewMySQL(c.MallUgcMySQL),
		ticketDB:       xsql.NewMySQL(c.TicketMySQL),
		client:         &http.Client{Timeout: time.Second * 5},
		ItemSalesMax:   make(map[string]int),
		ItemSalesMin:   make(map[string]int),
		ItemWishMax:    make(map[string]int),
		ItemWishMin:    make(map[string]int),
		ItemCommentMax: make(map[string]int),
		ItemCommentMin: make(map[string]int),
	}
	es, err = elastic.NewClient(
		elastic.SetURL(c.ElasticSearch.Addr...),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(time.Duration(c.ElasticSearch.Check)),
		elastic.SetErrorLog(&elog{}),
		elastic.SetInfoLog(&ilog{}),
	)
	if err != nil {
		panic(err)
	}
	dao.es = es

	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
	d.mallDB.Close()
	d.ugcDB.Close()
	d.ticketDB.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	return d.mallDB.Ping(c)
}

type ilog struct{}
type elog struct{}

// Printf printf.
func (l *ilog) Printf(format string, v ...interface{}) {
	log.Info(format, v...)
}

// Printf printf.
func (l *elog) Printf(format string, v ...interface{}) {
	log.Error(format, v...)
}
