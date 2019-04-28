package dao

import "fmt"

func (d *Dao) getRecInfoKey(roomId int64) (key string, expire int){
	key = fmt.Sprintf(_recInfoKey, roomId)
	expire = _recInfoExpireTtl
	return
}