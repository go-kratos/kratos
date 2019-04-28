package rpc

import (
	"context"

	"go-common/app/admin/main/aegis/conf"
	account "go-common/app/service/main/account/api"
	relmod "go-common/app/service/main/relation/model"
	relrpc "go-common/app/service/main/relation/rpc/client"
	uprpc "go-common/app/service/main/up/api/v1"

	"google.golang.org/grpc"
)

// Dao dao
type Dao struct {
	c *conf.Config

	//gorpc
	relRPC RelationRPC
	//grpc
	AccountClient AccRPC
	UpClient      UpRPC
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c: c,
	}
	if c.Debug != "local" {
		dao.relRPC = relrpc.New(c.RPC.Rel)
		var err error
		if dao.AccountClient, err = account.NewClient(c.GRPC.AccRPC); err != nil {
			panic(err)
		}
		if dao.UpClient, err = uprpc.NewClient(c.GRPC.UpRPC); err != nil {
			panic(err)
		}
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

//RelationRPC .
type RelationRPC interface {
	Stats(c context.Context, arg *relmod.ArgMids) (res map[int64]*relmod.Stat, err error)
}

//AccRPC .
type AccRPC interface {
	Info3(ctx context.Context, in *account.MidReq, opts ...grpc.CallOption) (*account.InfoReply, error)
	Cards3(ctx context.Context, in *account.MidsReq, opts ...grpc.CallOption) (*account.CardsReply, error)
	ProfileWithStat3(ctx context.Context, in *account.MidReq, opts ...grpc.CallOption) (*account.ProfileStatReply, error)
}

//UpRPC .
type UpRPC interface {
	UpSpecial(ctx context.Context, in *uprpc.UpSpecialReq, opts ...grpc.CallOption) (*uprpc.UpSpecialReply, error)
	UpsSpecial(ctx context.Context, in *uprpc.UpsSpecialReq, opts ...grpc.CallOption) (*uprpc.UpsSpecialReply, error)
	UpGroups(ctx context.Context, in *uprpc.NoArgReq, opts ...grpc.CallOption) (*uprpc.UpGroupsReply, error)
}
