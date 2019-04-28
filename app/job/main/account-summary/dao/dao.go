package dao

import (
	"context"

	"go-common/app/job/main/account-summary/conf"
	member "go-common/app/service/main/member/api/gorpc"
	relation "go-common/app/service/main/relation/rpc/client"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"

	"go-common/library/database/hbase.v2"
)

// Dao dao
type Dao struct {
	c               *conf.Config
	AccountSumHBase *hbase.Client
	MemberService   *member.Service
	RelationService *relation.Service
	httpClient      *bm.Client
	MemberDB        *sql.DB
	RelationDB      *sql.DB
	PassportDB      *sql.DB
}

// New init mysql db
func New(c *conf.Config) *Dao {
	dao := &Dao{
		c:               c,
		AccountSumHBase: hbase.NewClient(c.AccountSummaryHBase),
		MemberService:   member.New(c.MemberService),
		RelationService: relation.New(c.RelationService),
		httpClient:      bm.NewClient(c.HTTPClient),
		MemberDB:        sql.NewMySQL(c.MemberDB),
		RelationDB:      sql.NewMySQL(c.RelationDB),
		PassportDB:      sql.NewMySQL(c.PassportDB),
	}
	return dao
}

// Close close the resource.
func (d *Dao) Close() {
}

// Ping dao ping
func (d *Dao) Ping(ctx context.Context) error {
	return nil
}
