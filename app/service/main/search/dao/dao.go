package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/search/conf"
	"go-common/library/log"
	"go-common/library/stat/prom"

	elastic "gopkg.in/olivere/elastic.v5"
)

type Dao struct {
	// conf
	c *conf.Config

	// esPool
	esPool map[string]*elastic.Client
	// sms
	sms *sms
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
	}
	d.sms = newSMS(d)
	// cluster
	d.esPool = newEsPool(c, d)
	return
}

// BulkItem .
type BulkItem interface {
	IndexName() string
	IndexType() string
	IndexID() string
}

// BulkMapItem .
type BulkMapItem interface {
	IndexName() string
	IndexType() string
	IndexID() string
	PField() map[string]interface{}
}

// newEsCluster cluster action
func newEsPool(c *conf.Config, d *Dao) (esCluster map[string]*elastic.Client) {
	esCluster = make(map[string]*elastic.Client)
	for esName, e := range c.Es {
		if client, err := elastic.NewClient(elastic.SetURL(e.Addr...)); err == nil {
			esCluster[esName] = client
		} else {
			PromError("es:集群连接失败", "cluster: %s, %v", esName, err)
			if err := d.SendSMS(fmt.Sprintf("[search-job]%s集群连接失败", esName)); err != nil {
				PromError("es:集群连接短信失败", "cluster: %s, %v", esName, err)
			}
		}
	}
	return
}

// PromError prometheus error count.
func PromError(name, format string, args ...interface{}) {
	prom.BusinessErrCount.Incr(name)
	log.Error(format, args...)
}

// Ping health of db.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.pingESCluster(c); err != nil {
		PromError("es:ping", "Ping %v", err)
	}
	return
}

// pingESCluster ping es cluster
func (d *Dao) pingESCluster(ctx context.Context) (err error) {
	for name := range d.c.Es {
		client, ok := d.esPool[name]
		if !ok {
			continue
		}
		_, _, err = client.Ping(d.c.Es["replyExternal"].Addr[0]).Do(ctx)
		if err != nil {
			PromError("archiveESClient:Ping", "dao.pingESCluster error(%v) ", err)
			return
		}
	}
	return
}
