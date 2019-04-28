package global

import (
	"go-common/app/admin/main/up-rating/conf"
	accrpc "go-common/app/service/main/account/rpc/client"
)

var (
	accRPC *accrpc.Service3
)

// Init resources
func Init(c *conf.Config) {
	accRPC = accrpc.New3(c.RPCClient.Account)
}

// GetAccRPC .
func GetAccRPC() *accrpc.Service3 {
	return accRPC
}
