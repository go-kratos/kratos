package global

import (
	"go-common/app/admin/main/up/conf"
	accgrpc "go-common/app/service/main/account/api"

	"github.com/pkg/errors"
)

var (
	accCli accgrpc.AccountClient
)

// GetAccClient .
func GetAccClient() accgrpc.AccountClient {
	return accCli
}

//Init init global
func Init(c *conf.Config) {
	var err error
	if accCli, err = accgrpc.NewClient(c.GRPCClient.Account); err != nil {
		panic(errors.WithMessage(err, "Failed to dial account service"))
	}
}
