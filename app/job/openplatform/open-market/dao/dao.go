package dao

import (
	"context"
	"time"

	"go-common/app/job/openplatform/open-market/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	elastic "gopkg.in/olivere/elastic.v5"
)

//Dao struct
type Dao struct {
	c *conf.Config
	// http client
	client *bm.Client
	// db
	ticketDB *sql.DB
	//es
	es    *elastic.Client
	esUgc *elastic.Client
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	var (
		err   error
		es    *elastic.Client
		esUgc *elastic.Client
	)
	d = &Dao{
		c:        c,
		client:   bm.NewClient(c.HTTPClient),
		ticketDB: sql.NewMySQL(c.DB.TicketDB),
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
	esUgc, err = elastic.NewClient(
		elastic.SetURL(c.ElasticSearchUgc.Addr...),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(time.Duration(c.ElasticSearch.Check)),
		elastic.SetErrorLog(&elog{}),
		elastic.SetInfoLog(&ilog{}),
	)
	if err != nil {
		panic(err)
	}
	d.es = es
	d.esUgc = esUgc
	return d
}

// Ping ping health.
func (d *Dao) Ping(c context.Context) (err error) {
	return d.ticketDB.Ping(c)
}

// Close close.
func (d *Dao) Close() (err error) {
	if err = d.ticketDB.Close(); err != nil {
		log.Error("dao.ticketDB.Close() error(%v)", err)
	}
	return
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
