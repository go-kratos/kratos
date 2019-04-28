package dao

import (
	"context"

	"go-common/app/interface/main/ugcpay/conf"
	archive "go-common/app/service/main/archive/api"
	ugcpay "go-common/app/service/main/ugcpay/api/grpc/v1"
)

// Dao dao
type Dao struct {
	c          *conf.Config
	ugcpayAPI  ugcpay.UGCPayClient
	archiveAPI archive.ArchiveClient
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c: c,
	}
	var err error
	if dao.ugcpayAPI, err = ugcpay.NewClient(nil); err != nil {
		panic(err)
	}
	if dao.archiveAPI, err = archive.NewClient(nil); err != nil {
		panic(err)
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	return nil
}
