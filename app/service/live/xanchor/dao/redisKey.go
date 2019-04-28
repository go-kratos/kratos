package dao

import (
	"fmt"
)

type resp struct {
	RedisKey string `json:"redis_key"`
	TimeOut  int    `json:"time_out"`
}

//DanmuTotalNum  开播弹幕数统计
func (d *Dao) DanmuNumKey(roomId int64) (resp *resp) {
	resp.RedisKey = fmt.Sprintf("danmu_num_key_%d", roomId)
	resp.TimeOut = 24 * 60 * 60
	return
}

//GiftTotalNum  开播礼物数统计
func (d *Dao) GiftNumKey(roomId int64) (resp *resp) {
	resp.RedisKey = fmt.Sprintf("gift_num_key_%d", roomId)
	resp.TimeOut = 24 * 60 * 60
	return
}

//GiftGoldTotalNum  开播金瓜子数量统计
func (d *Dao) GiftGoldNumKey(roomId int64) (resp *resp) {
	resp.RedisKey = fmt.Sprintf("gift_gold_num_key_%d", roomId)
	resp.TimeOut = 24 * 60 * 60
	return
}

//GiftGoldTotalAmount 开播金瓜子金额统计
func (d *Dao) GiftGoldAmountKey(roomId int64) (resp *resp) {
	resp.RedisKey = fmt.Sprintf("gift_gold_num_key_%d", roomId)
	resp.TimeOut = 24 * 60 * 60
	return
}
