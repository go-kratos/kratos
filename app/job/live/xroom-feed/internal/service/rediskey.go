package service

import "fmt"

const (
	_recPoolTtl       = 60
	_recPoolKey       = "rec_pool_%d"
	_recInfoExpireTtl = 86400
	_recInfoKey       = "rec_info_%d"
	_recConfExpireTtl = 86400
	_recConfKey       = "rec_conf"
)

func (s *Service) getRecPoolKey(id int) (key string, expire int) {
	key = fmt.Sprintf(_recPoolKey, id)
	expire = _recPoolTtl
	return
}

func (s *Service) getRecInfoKey(roomId int) (key string, expire int) {
	key = fmt.Sprintf(_recInfoKey, roomId)
	expire = _recInfoExpireTtl
	return
}

func (s *Service) getRecConfKey() (key string, expire int) {
	return _recConfKey, _recConfExpireTtl
}
