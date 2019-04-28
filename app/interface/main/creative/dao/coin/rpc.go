package coin

import (
	"context"
	"fmt"
	coinclient "go-common/app/service/main/coin/api"
	"go-common/library/log"
	"time"
)

// AddCoin ModifyCoins with grpc client
func (d *Dao) AddCoin(c context.Context, mid, aid int64, coin float64, ip string) (err error) {
	arg := &coinclient.ModifyCoinsReq{
		Mid:       mid,
		Count:     coin,
		Reason:    fmt.Sprintf("删除稿件av%d,扣硬币", aid),
		IP:        ip,
		Operator:  "main.archive.creative",
		CheckZero: 1,
		Ts:        time.Now().Unix(),
	}
	_, err = d.coinClient.ModifyCoins(c, arg)
	if err != nil {
		log.Error("ModifyCoins arg(%+v), err(%+v)", arg, err)
	}
	return
}
