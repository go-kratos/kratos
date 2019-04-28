package model

import (
	"fmt"
	"time"
)

//缓存key
const (
	//RedisSalesLimitKey 售卖时间限制,mid:SALESLIMIT:ANTI
	RedisSalesLimitKey = "%d:SALESLIMIT:ANTI"

	//RedisIPChangeKey 用户ip变更，mid:IPCHANGELIMIT:ANTI
	RedisIPChangeKey = "%d:IPCHANGELIMIT:ANTI"

	//RedisCreateMIDLimitKey 同一mid下单次数限制,mid:itemId:MID:CREATELIMIT:ANTI
	RedisCreateMIDLimitKey = "%d:%d:MID:CREATELIMIT:ANTI"

	//RedisCreateIPLimitKey 同一mid下单次数限制,ip:itemId:IP:CREATELIMIT:ANTI
	RedisCreateIPLimitKey = "%s|%d:IP:CREATELIMIT:ANTI"

	//RedisUserVoucherKey 用户凭证key，mid:voucher:voucherType:VOUCHER:ANTI
	RedisUserVoucherKey        = "%d:%s:%d:VOUCHER:ANTI"
	RedisUserVoucherKeyTimeOut = 600

	//RedisGeetestCountKey 拉起极验的总数
	RedisGeetestCountKey        = "%d:ANTI:GEETEST:COUNT"
	RedisGeetestCountKeyTimeOut = 3600

	//RedisMIDBlackKey mid黑名单key
	RedisMIDBlackKey = "ANTI:MID:BLACK:%d:%d"

	//RedisIPBlackKey ip黑名单key
	RedisIPBlackKey = "ANTI:IP:BLACK:%d:%s"
)

//GetSalesLimitKey 获取售卖时间限制key
func GetSalesLimitKey(mid int64) (key string) {
	return fmt.Sprintf(RedisSalesLimitKey, mid)
}

//GetIPChangeKey 获取用户ip变更key
func GetIPChangeKey(mid int64) (key string) {
	return fmt.Sprintf(RedisIPChangeKey, mid)
}

//GetCreateMIDLimitKey 获取mid创单限制key
func GetCreateMIDLimitKey(mid int64, itemID int64) (key string) {
	return fmt.Sprintf(RedisCreateMIDLimitKey, mid, itemID)
}

//GetCreateIPLimitKey 获取ip创单限制key
func GetCreateIPLimitKey(ip string, itemID int64) (key string) {
	return fmt.Sprintf(RedisCreateIPLimitKey, ip, itemID)
}

//GetUserVoucherKey 获取用户凭证key
func GetUserVoucherKey(mid int64, voucher string, voucherType int64) (key string) {
	return fmt.Sprintf(RedisUserVoucherKey, mid, voucher, voucherType)
}

//GetGeetestCountKey 获取极验总数key
func GetGeetestCountKey() (key string) {
	current := time.Now().Unix()
	return fmt.Sprintf(RedisGeetestCountKey, current/RedisGeetestCountKeyTimeOut)
}

//GetMIDBlackKey 获取mid黑名单key
func GetMIDBlackKey(customerId int64, mid int64) (key string) {
	return fmt.Sprintf(RedisMIDBlackKey, customerId, mid)
}

//GetIPBlackKey 获取mid黑名单key
func GetIPBlackKey(customerId int64, clientIP string) (key string) {
	return fmt.Sprintf(RedisIPBlackKey, customerId, clientIP)
}
