package global

import (
	"go-common/app/interface/main/mcn/conf"
	taggrpc "go-common/app/interface/main/tag/api"
	accgrpc "go-common/app/service/main/account/api"
	arcgrpc "go-common/app/service/main/archive/api"
	memgrpc "go-common/app/service/main/member/api"
	"go-common/library/cache/memcache"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

var (
	accGRPC accgrpc.AccountClient
	memGRPC memgrpc.MemberClient
	arcGRPC arcgrpc.ArchiveClient
	tagGRPC taggrpc.TagRPCClient

	mc *memcache.Pool

	bmClient *bm.Client
)

// GetAccGRPC .
func GetAccGRPC() accgrpc.AccountClient {
	return accGRPC
}

// GetMemGRPC .
func GetMemGRPC() memgrpc.MemberClient {
	return memGRPC
}

// GetArcGRPC .
func GetArcGRPC() arcgrpc.ArchiveClient {
	return arcGRPC
}

// GetTagGRPC .
func GetTagGRPC() taggrpc.TagRPCClient {
	return tagGRPC
}

// GetMc get mc
func GetMc() *memcache.Pool {
	return mc
}

// GetBMClient get http client
func GetBMClient() *bm.Client {
	return bmClient
}

//Init init global
func Init(c *conf.Config) {
	var err error
	if accGRPC, err = accgrpc.NewClient(c.GRPCClient.Account); err != nil {
		panic(errors.WithMessage(err, "Failed to dial account service"))
	}
	if memGRPC, err = memgrpc.NewClient(c.GRPCClient.Member); err != nil {
		panic(errors.WithMessage(err, "Failed to dial member service"))
	}
	if arcGRPC, err = arcgrpc.NewClient(c.GRPCClient.Archive); err != nil {
		panic(errors.WithMessage(err, "Failed to dial archive service"))
	}
	if tagGRPC, err = taggrpc.NewClient(c.GRPCClient.Tag); err != nil {
		panic(errors.WithMessage(err, "Failed to dial tag service"))
	}

	mc = memcache.NewPool(&c.Memcache.Config)
	bmClient = bm.NewClient(c.HTTPClient)
}
