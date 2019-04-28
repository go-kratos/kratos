package coin

import (
	"context"

	"go-common/app/interface/main/app-intl/conf"
	coinrpc "go-common/app/service/main/coin/api/gorpc"
	"go-common/app/service/main/coin/model"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Dao is coin dao
type Dao struct {
	coinRPC *coinrpc.Service
}

// New initial coin dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		coinRPC: coinrpc.New(c.CoinRPC),
	}
	return
}

// AddCoins add coin to upper.
func (d *Dao) AddCoins(c context.Context, aid, mid, upID, maxCoin, avtype, multiply int64, typeID int16, pubTime int64) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &model.ArgAddCoin{Mid: mid, UpMid: upID, MaxCoin: maxCoin, Aid: aid, AvType: avtype, Multiply: multiply, RealIP: ip, TypeID: typeID, PubTime: pubTime}
	return d.coinRPC.AddCoins(c, arg)
}

// ArchiveUserCoins .
func (d *Dao) ArchiveUserCoins(c context.Context, aid, mid, avType int64) (res *model.ArchiveUserCoins, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &model.ArgCoinInfo{Mid: mid, Aid: aid, AvType: avType, RealIP: ip}
	if res, err = d.coinRPC.ArchiveUserCoins(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

// UserCoins get user coins
func (d *Dao) UserCoins(c context.Context, mid int64) (count float64, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &model.ArgCoinInfo{Mid: mid, RealIP: ip}
	if count, err = d.coinRPC.UserCoins(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}
