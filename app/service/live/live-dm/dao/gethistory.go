package dao

import (
	"context"
	"fmt"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

//GetHistoryData 获取历史数据
func (d *Dao) GetHistoryData(ctx context.Context, roomid int64) (result map[string][]string, err error) {
	var conn = d.redis.Get(ctx)
	defer conn.Close()
	var admkey = fmt.Sprintf("%s%d", _adminMsgHistoryCache, roomid)
	var userkey = fmt.Sprintf("%s%d", _msgHistoryCache, roomid)

	if err = conn.Send("LRANGE", admkey, 0, 9); err != nil {
		log.Error("DM: LRANGE ADMIN KEY ROOMID %d ERR:%v", roomid, err)
	}
	if err = conn.Send("LRANGE", userkey, 0, 9); err != nil {
		log.Error("DM: LRANGE USER KEY ROOMID %d ERR:%v", roomid, err)
	}

	if err = conn.Flush(); err != nil {
		log.Error("DM: Flush KEY ROOMID %d ERR:%v", roomid, err)
		return nil, err
	}

	var admin, user [][]byte
	admin, err = redis.ByteSlices(conn.Receive())
	if err != nil {
		log.Error("DM: ByteSlices ADMIN KEY ROOMID %d ERR:%v", roomid, err)
		return nil, err
	}
	user, err = redis.ByteSlices(conn.Receive())
	if err != nil {
		log.Error("DM: ByteSlices USER KEY ROOMID %d ERR:%v", roomid, err)
		return nil, err
	}

	result = make(map[string][]string)
	result["admin"] = make([]string, 0, 10)
	result["room"] = make([]string, 0, 10)
	for i := len(admin) - 1; i >= 0; i-- {
		result["admin"] = append(result["admin"], string(admin[i]))
	}
	for i := len(user) - 1; i >= 0; i-- {
		result["room"] = append(result["room"], string(user[i]))
	}
	return result, nil
}
