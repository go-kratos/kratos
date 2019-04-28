package coin

import (
	"context"

	"go-common/app/service/main/coin/model"
	"go-common/library/net/rpc"
)

const (
	_addCoins         = "RPC.AddCoins"
	_archiveUserCoins = "RPC.ArchiveUserCoins"
	_userCoins        = "RPC.UserCoins"
	_modifyCoin       = "RPC.ModifyCoin"
	_list             = "RPC.List"
	_userLog          = "RPC.UserLog"
	_addUserCoinExp   = "RPC.AddUserCoinExp"
	_updateAddCoin    = "RPC.UpdateAddCoin"
	_todayExp         = "RPC.TodayExp"
)

const (
	_appid = "community.service.coin"
)

var (
	_noRes = &struct{}{}
)

// Service rpc service.
type Service struct {
	client *rpc.Client2
}

//go:generate mockgen -source coin.go  -destination mock.go -package coin

// RPC coin rpc
type RPC interface {
	AddCoins(c context.Context, arg *model.ArgAddCoin) (err error)
	ArchiveUserCoins(c context.Context, arg *model.ArgCoinInfo) (res *model.ArchiveUserCoins, err error)
	UserCoins(c context.Context, arg *model.ArgCoinInfo) (count float64, err error)
}

// New new service.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// AddCoins coin to archive.
func (s *Service) AddCoins(c context.Context, arg *model.ArgAddCoin) (err error) {
	err = s.client.Call(c, _addCoins, arg, _noRes)
	return
}

// ArchiveUserCoins get archive User added coins.
func (s *Service) ArchiveUserCoins(c context.Context, arg *model.ArgCoinInfo) (res *model.ArchiveUserCoins, err error) {
	res = &model.ArchiveUserCoins{}
	err = s.client.Call(c, _archiveUserCoins, arg, res)
	return
}

// UserCoins get user coins.
func (s *Service) UserCoins(c context.Context, arg *model.ArgCoinInfo) (count float64, err error) {
	err = s.client.Call(c, _userCoins, arg, &count)
	return
}

// ModifyCoin modify user coin.
func (s *Service) ModifyCoin(c context.Context, arg *model.ArgModifyCoin) (count float64, err error) {
	err = s.client.Call(c, _modifyCoin, arg, &count)
	return
}

// List coin added list.
func (s *Service) List(c context.Context, arg *model.ArgList) (list []*model.List, err error) {
	err = s.client.Call(c, _list, arg, &list)
	return
}

// UserLog user log
func (s *Service) UserLog(c context.Context, arg *model.ArgLog) (logs []*model.Log, err error) {
	err = s.client.Call(c, _userLog, arg, &logs)
	return
}

// AddUserCoinExp add user coin exp
func (s *Service) AddUserCoinExp(c context.Context, arg *model.ArgAddUserCoinExp) (err error) {
	err = s.client.Call(c, _addUserCoinExp, arg, _noRes)
	return
}

// UpdateAddCoin update db after add coin for job.
func (s *Service) UpdateAddCoin(c context.Context, arg *model.Record) (err error) {
	err = s.client.Call(c, _updateAddCoin, arg, _noRes)
	return
}

// TodayExp get today exp
func (s *Service) TodayExp(c context.Context, arg *model.ArgMid) (number int64, err error) {
	err = s.client.Call(c, _todayExp, arg, &number)
	return
}
