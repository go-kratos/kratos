package account

import (
	"context"
	"fmt"

	"go-common/app/interface/main/app-player/conf"
	accrpc "go-common/app/service/main/account/api"

	"github.com/pkg/errors"
)

// Dao is account dao.
type Dao struct {
	// rpc
	accRPC accrpc.AccountClient
}

// New account dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{}
	var err error
	d.accRPC, err = accrpc.NewClient(c.AccountClient)
	if err != nil {
		panic(fmt.Sprintf("account NewClient error(%v)", err))
	}
	return
}

// Card get card
func (d *Dao) Card(c context.Context, mid int64) (card *accrpc.CardReply, err error) {
	if card, err = d.accRPC.Card3(c, &accrpc.MidReq{Mid: mid}); err != nil {
		err = errors.Wrapf(err, "%v", mid)
	}
	return
}
