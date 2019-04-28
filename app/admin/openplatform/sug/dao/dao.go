package dao

import (
	"context"
	"net/http"
	"time"

	"go-common/app/admin/openplatform/sug/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"

	elastic "gopkg.in/olivere/elastic.v5"
)

// Dao dao
type Dao struct {
	c        *conf.Config
	redis    *redis.Pool
	es       *elastic.Client
	client   *http.Client
	xclient  *xhttp.Client
	dbMall   *sql.DB
	dbTicket *sql.DB
}

// New init redis,es
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:        c,
		redis:    redis.NewPool(c.Redis),
		xclient:  xhttp.NewClient(c.HTTPClient),
		client:   &http.Client{Timeout: time.Second * 5},
		dbMall:   sql.NewMySQL(c.DB.MallDB),
		dbTicket: sql.NewMySQL(c.DB.TicketDB),
	}
	es, err := elastic.NewClient(
		elastic.SetURL(c.Es.Addr...),
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
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.dbTicket.Ping(c); err != nil {
		log.Error("Ping ticket DB error(%v)", err)
		return
	}
	if err = d.dbMall.Ping(c); err != nil {
		log.Error("Ping mall DB error(%v)", err)
		return
	}
	if err = d.pingESCluster(c); err != nil {
		log.Error("es:ping", "ping %v", err)
		return
	}
	redisConn := d.redis.Get(c)
	defer redisConn.Close()
	if _, err = redisConn.Do("SET", "ping", "pong"); err != nil {
		redisConn.Close()
		log.Error("redis.SET error(%v)", err)
		return
	}
	return
}

// pingEsCluster ping es cluster
func (d *Dao) pingESCluster(ctx context.Context) (err error) {
	_, _, err = d.es.Ping(d.c.Es.Addr[0]).Do(ctx)
	if err != nil {
		log.Error("Es:Ping", "%s:Ping error(%v)", d.c.Es.Addr[0], err)
		return
	}
	return
}
