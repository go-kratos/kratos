package account

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-wall/conf"
	account "go-common/app/service/main/account/model"
	accrpc "go-common/app/service/main/account/rpc/client"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_addVIPURL = "/x/internal/big/useBatchInfo"
)

// Dao account dao
type Dao struct {
	client *httpx.Client
	// account rpc
	accRPC *accrpc.Service3
	// url
	addVIPURL string
}

// New account dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client: httpx.NewClient(c.HTTPClient),
		// account rpc
		accRPC: accrpc.New3(c.AccountRPC),
		// url
		addVIPURL: c.Host.APICo + _addVIPURL,
	}
	return
}

// AddVIP add user vip
func (d *Dao) AddVIP(c context.Context, mid, orderNo int64, batchID int, remark string) (msg string, err error) {
	params := url.Values{}
	params.Set("orderNo", strconv.FormatInt(orderNo, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("remark", remark)
	params.Set("batchId", strconv.Itoa(batchID))
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
	}
	if err = d.client.Post(c, d.addVIPURL, "", params, &res); err != nil {
		log.Error("account add vip url(%v) error(%v)", d.addVIPURL+"?"+params.Encode(), err)
		return
	}
	msg = res.Msg
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("account add vip url(%v) res code(%d)", d.addVIPURL+"?"+params.Encode(), res.Code)
		return
	}
	return
}

// Info user info
func (d *Dao) Info(ctx context.Context, mid int64) (res *account.Info, err error) {
	arg := &account.ArgMid{
		Mid: mid,
	}
	if res, err = d.accRPC.Info3(ctx, arg); err != nil {
		log.Error("d.accRPC.Info3(%v) error(%v)", arg, err)
		return
	}
	return
}
