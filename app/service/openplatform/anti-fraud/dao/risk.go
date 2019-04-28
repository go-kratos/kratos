package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"github.com/json-iterator/go"
	"go-common/app/service/openplatform/anti-fraud/conf"
	"go-common/app/service/openplatform/anti-fraud/model"
	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"
	"net/http"
	"time"
)

//CheckSalesTime 检查售卖时间
func (d *Dao) CheckSalesTime(c context.Context, mid, itemID, salesTime, saleTimeOut int64) (err error) {
	key := model.GetSalesLimitKey(mid)
	conn := d.redis.Get(c)
	defer conn.Close()

	var (
		flag int64
		data []byte
	)
	data, _ = redis.Bytes(conn.Do("GET", key))
	json.Unmarshal(data, &flag)
	if flag == 1 {
		return ecode.AntiSalesTimeErr
	}
	if salesTime > time.Now().Unix() {
		conn.Do("SET", key, 1, "EX", saleTimeOut)
		return ecode.AntiSalesTimeErr
	}
	return nil
}

//CheckIPChange 检查用户ip变更
func (d *Dao) CheckIPChange(c context.Context, mid int64, ip string, changeTime int64) (err error) {
	key := model.GetIPChangeKey(mid)
	conn := d.redis.Get(c)
	defer conn.Close()

	var (
		flag string
		data []byte
	)
	data, _ = redis.Bytes(conn.Do("GET", key))
	json.Unmarshal(data, &flag)

	go func(c context.Context, key string, ip string) {
		conn := d.redis.Get(c)
		defer conn.Close()
		conn.Do("SET", key, ip, "EX", changeTime)
	}(context.Background(), key, ip)
	if flag != "" && flag != ip {
		return ecode.AntiIPChangeLimit
	}
	return nil
}

//CheckLimitNum 检查限制次数
func (d *Dao) CheckLimitNum(c context.Context, key string, num int64, pastTime int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	pastTime = pastTime * 1e9
	currentTime := time.Now().UnixNano()
	conn.Do("ZADD", key, currentTime, currentTime)
	data, _ := redis.Int64(conn.Do("ZCOUNT", key, currentTime-pastTime, currentTime))
	go func(c context.Context, key string, pastTime, currentTime int64, flag bool) {
		conn := d.redis.Get(c)
		defer conn.Close()
		conn.Do("EXPIRE", key, pastTime/1e9)
		if flag {
			conn.Do("ZREMRANGEBYRANK", key, 0, 0)
		}
	}(context.Background(), key, pastTime, currentTime, data > num)
	if data > num {
		return ecode.AntiLimitNumUpper
	}
	return nil
}

//Voucher 凭证
func (d *Dao) Voucher(c context.Context, mid int64, ip string, itemID, customer, voucherType int64) (voucher string) {
	s := make([]byte, 0)
	buf := bytes.NewBuffer(s)
	binary.Write(buf, binary.BigEndian, mid)
	binary.Write(buf, binary.BigEndian, []byte(ip))
	binary.Write(buf, binary.BigEndian, itemID)
	binary.Write(buf, binary.BigEndian, customer)
	binary.Write(buf, binary.BigEndian, voucherType)
	binary.Write(buf, binary.BigEndian, time.Now().UnixNano())
	digest := md5.Sum(buf.Bytes())
	voucher = hex.EncodeToString(digest[:])
	conn := d.redis.Get(c)
	defer conn.Close()
	key := model.GetUserVoucherKey(mid, voucher, voucherType)
	conn.Do("SET", key, 1, "EX", model.RedisUserVoucherKeyTimeOut)
	return
}

//CheckVoucher 验证用户凭证,一次性
func (d *Dao) CheckVoucher(c context.Context, mid int64, voucher string, voucherType int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := model.GetUserVoucherKey(mid, voucher, voucherType)
	data, _ := redis.Int64(conn.Do("INCR", key))
	conn.Do("DEL", key)
	if data < 2 {
		return ecode.AntiCheckVoucherErr
	}
	return nil
}

//IncrGeetestCount 统计一小时内极验的请求数
func (d *Dao) IncrGeetestCount(c context.Context) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := model.GetGeetestCountKey()
	conn.Do("INCR", key)
	conn.Do("EXPIRE", key, model.RedisGeetestCountKeyTimeOut)
}

//CheckGeetestCount 检查极验总数是否达到上限
func (d *Dao) CheckGeetestCount(c context.Context) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := model.GetGeetestCountKey()
	var (
		data  []byte
		count int64
	)
	data, _ = redis.Bytes(conn.Do("GET", key))
	json.Unmarshal(data, &count)
	if count > d.c.Geetest.Count {
		log.Info("极验总数达到上限")
		return ecode.AntiGeetestCountUpper
	}
	return
}

// CheckBlack 检测黑名单
func (d *Dao) CheckBlack(c context.Context, customerId, mid int64, clientIP string) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	midKey := model.GetMIDBlackKey(customerId, mid)
	ipKey := model.GetIPBlackKey(customerId, clientIP)

	data, err := redis.Int64s(conn.Do("MGET", midKey, ipKey))

	if err != nil {
		log.Info("获取黑名单出错 %s %s", midKey, ipKey)
		return nil
	}

	for _, value := range data {
		if value == 1 {
			return ecode.AntiBlackErr
		}
	}

	return nil
}

// PayShield 支付风控
func (d *Dao) PayShield(c context.Context, data *model.ShieldData) {
	var res struct {
		errno int64
		msg   string
		data  interface{}
	}

	params, err := jsoniter.Marshal(data)
	if err != nil {
		log.Info("json marshal err %v", err)
		return
	}

	log.Info("req pay shield params %s", string(params))

	req, err := http.NewRequest("POST", conf.Conf.URL.Shield, bytes.NewBuffer(params))
	if err != nil {
		log.Warn("new request err  url %s, err %v", conf.Conf.URL.Shield, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	err = d.client.Do(c, req, &res)
	if err != nil || res.errno != 0 {
		log.Warn("client do err  url %s, err %v", conf.Conf.URL.Shield, err)
		return
	}

}

// SetexRedisKey 设置redis key
func (d *Dao) SetexRedisKey(c context.Context, key string, timeout int64) {
	conn := d.redis.Get(c)
	defer conn.Close()

	conn.Do("SET", key, 1, "EX", timeout)
}

// AddPayData .
func (d *Dao) AddPayData(data *model.ShieldData) {
	d.payData <- data
}

// SyncPayShield .
func (d *Dao) SyncPayShield(c context.Context) {
	for {
		data := <-d.payData
		d.PayShield(c, data)
	}
}
