package dao

import (
	"encoding/json"
	"go-common/app/service/live/live-dm/conf"
	titansSdk "go-common/app/service/live/resource/sdk"
	"go-common/library/log"
	"go-common/library/log/infoc"
	"go-common/library/queue/databus"
)

var (
	//拜年祭投递databus
	bndatabus *databus.Databus
	//InfocDMSend 弹幕发送成功上报
	InfocDMSend *infoc.Infoc
	//InfocDMErr 弹幕发送失败上报
	InfocDMErr *infoc.Infoc
)

//InitDatabus 初始化databus
func InitDatabus(c *conf.Config) {
	bndatabus = databus.New(c.BNDatabus)
}

//InitLancer 初始化lancer上报
func InitLancer(c *conf.Config) {
	InfocDMErr = infoc.New(c.Lancer.DMErr)
	InfocDMSend = infoc.New(c.Lancer.DMSend)
}

//InitTitan 初始化kv配置
func InitTitan() {
	conf := &titansSdk.Config{
		TreeId: 72389,
		Expire: 1,
	}
	titansSdk.Init(conf)
}

//LimitConf 弹幕限制配置
type LimitConf struct {
	AreaLimit        bool   `json:"areaLimit"`
	AllUserLimit     bool   `json:"AllUserLimit"`
	LevelLimitStatus bool   `json:"LevelLimitStatus"`
	LevelLimit       int64  `json:"LevelLimit"`
	RealName         bool   `json:"RealName"`
	PhoneLimit       bool   `json:"PhoneLimit"`
	MsgLength        int    `json:"MsgLength"`
	DmNum            int64  `json:"DmNum"`
	DMPercent        int64  `json:"DMPercent"`
	DMwhitelist      bool   `json:"DMwhitelist"`
	DMwhitelistID    string `json:"DMwhitelistID"`
}

//GetDMCheckConf 获取弹幕配置参数
func (l *LimitConf) GetDMCheckConf() {

	cf, err := titansSdk.Get("dmLimit")
	if err != nil {
		log.Error("DM: get conf err:%+v", err)
		return
	}

	conf := &LimitConf{}
	err = json.Unmarshal([]byte(cf), conf)
	if err != nil {
		log.Error("DM: decode conf jsons err:%+v conf:%s", err, cf)
		return
	}

	l.AreaLimit = conf.AreaLimit
	l.AllUserLimit = conf.AllUserLimit
	l.LevelLimitStatus = conf.LevelLimitStatus
	l.LevelLimit = conf.LevelLimit
	l.RealName = conf.RealName
	l.PhoneLimit = conf.PhoneLimit
	l.DmNum = conf.DmNum
	l.DMwhitelist = conf.DMwhitelist
	l.DMPercent = conf.DMPercent
	l.MsgLength = conf.MsgLength
	l.DMwhitelistID = conf.DMwhitelistID

}
