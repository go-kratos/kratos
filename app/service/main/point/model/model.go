package model

// point change type
const (
	ExchangeVip            = iota + 1
	Charge                 //充电
	Contract               //承包
	PointSystem            //系统发放
	FYMReward              //分院帽奖励
	ExchangePendant        //兑换挂件
	MJActive               //萌节活动
	ReAcquirePointDedution //重复领取
)

// system.
const (
	ActivityGiveRemark = "承包额外赠送"
	ActivitySendTimes1 = 1
	ActivitySendTimes2 = 2
	ActivityMixBuyBp   = 1
	ActivityOutOfBuyBp = 100
	ActivityGivePoint  = 1000
	SUCCESS            = 1
)
