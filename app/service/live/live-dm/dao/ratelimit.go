package dao

import (
	"context"
	"crypto/md5"
	"fmt"
	"strconv"
	"time"

	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"
)

// LimitCheckInfo  频率限制检查参数
type LimitCheckInfo struct {
	UID     int64
	RoomID  int64
	Msg     string
	MsgType int64
	Dao     *Dao
	Conf    *LimitConf
}

const (
	_MaxGiftMsgNum = "r:m:stormmsgcount:"
	_MAXMsgNum     = "r:m:rmx:"
)

// LimitPerSec 每秒发言限制
func (l *LimitCheckInfo) LimitPerSec(ctx context.Context) error {
	key := fmt.Sprintf("%d.%d", l.UID, time.Now().Unix())
	var conn = l.Dao.redis.Get(ctx)
	defer conn.Close()
	ret, err := conn.Do("SET", key, 1, "EX", 1, "NX")
	if err != nil {
		log.Error("limitPerSec:conn.Do(SET EX NX) %s %d %v", key, 1, err)
		return nil
	}
	if ret != nil {
		return nil
	}
	return ecode.Error(0, "msg in 1s")
}

//LimitSameMsg 同一个用户同一房间5s 只能发送一条相同弹幕
func (l *LimitCheckInfo) LimitSameMsg(ctx context.Context) error {
	key := fmt.Sprintf("%d.%s", l.RoomID, md5.Sum([]byte(strconv.FormatInt(l.UID, 10)+l.Msg)))
	var conn = l.Dao.redis.Get(ctx)
	defer conn.Close()
	ret, err := conn.Do("SET", key, 1, "EX", 5, "NX")
	if err != nil {
		log.Error("DM LimitSameMsg conn.Do(SET, %s, 1, EX, 5, NX) error(%v)", key, err)
		return nil
	}
	if ret != nil {
		return nil
	}
	return ecode.Error(0, "msg repeat")
}

//LimitRoomPerSecond 单房间每秒只能发送制定条数弹幕
func (l *LimitCheckInfo) LimitRoomPerSecond(ctx context.Context) error {
	maxNum := l.Conf.DmNum
	percent := l.Conf.DMPercent
	danNum := maxNum * percent / 100.0
	giftNum := maxNum - danNum

	msgKey := fmt.Sprintf("%s.%d.%d", _MAXMsgNum, l.RoomID, time.Now().Unix())
	giftKey := fmt.Sprintf("%s.%d.%d", _MaxGiftMsgNum, l.RoomID, time.Now().Unix())
	var conn = l.Dao.redis.Get(ctx)
	defer conn.Close()

	if l.MsgType != 0 {
		//礼物弹幕
		if count, err := redis.Int64(conn.Do("INCRBY", giftKey, 1)); err != nil {
			log.Error("DMRateLimit: LimitRoomPerSecond INCRBY err: %v", err)
		} else {
			if count > giftNum {
				return ecode.Error(0, "max limit")
			}
		}
		if _, err := conn.Do("EXPIRE", giftKey, 2); err != nil {
			log.Error("DMRateLimit: LimitRoomPerSecond EXPIRE err: %v", err)
		}
		return nil
	}
	//普通弹幕
	var max = maxNum
	if exit, err := redis.Bool(conn.Do("EXISTS", giftKey)); err != nil {
		log.Error("DMRateLimit: LimitRoomPerSecond EXISTS gift err: %v", err)
	} else {
		if exit {
			max = danNum
		} else {
			max = maxNum
		}
	}

	if count, err := redis.Int64(conn.Do("INCRBY", msgKey, 1)); err != nil {
		log.Error("DMRateLimit: LimitRoomPerSecond INCR err: %v", err)
	} else {
		if count > max {
			return ecode.Error(0, "max limit")
		}
	}
	if _, err := conn.Do("EXPIRE", msgKey, 2); err != nil {
		log.Error("DMRateLimit: LimitRoomPerSecond EXPIRE err: %v", err)
	}
	return nil
}
