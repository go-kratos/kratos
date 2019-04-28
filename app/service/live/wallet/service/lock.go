package service

type UserLock interface {
	lock(uid int64) error
	release()
}

type RedisUserLock struct {
	ws *WalletService
}

func (r *RedisUserLock) lock(uid int64) error {
	return r.ws.lockSpecificUser(uid)
}

func (r *RedisUserLock) release() {
	r.ws.unLockUser()
}

type NopUserLock struct {
}

func (r *NopUserLock) lock(uid int64) error {
	return nil
}

func (r *NopUserLock) release() {
}
