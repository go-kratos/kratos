package resource

import (
	"go-common/app/admin/main/growup/conf"
	accgrpc "go-common/app/service/main/account/api"
	vip "go-common/app/service/main/vip/rpc/client"
	httpx "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

var (
	vipRPC             *vip.Service
	client             *httpx.Client
	accCli             accgrpc.AccountClient
	videoCategoryURL   string
	articleCategoryURL string
)

// Init .
func Init(c *conf.Config) {
	var err error
	vipRPC = vip.New(c.VipRPC)
	client = httpx.NewClient(c.HTTPClient)
	videoCategoryURL = c.Host.VideoType + "/videoup/types"
	articleCategoryURL = c.Host.ColumnType + "/x/article/categories"
	if accCli, err = accgrpc.NewClient(c.Account); err != nil {
		panic(errors.WithMessage(err, "Failed to dial account service"))
	}
}
