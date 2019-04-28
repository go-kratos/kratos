package dao

import (
	"context"
	"encoding/json"

	"go-common/app/service/live/resource/api/grpc/v1"
	"go-common/library/log"
)

//TitansTeam 话题team值
const TitansTeam = 40

//TitansKeyword 话题标签值
const TitansKeyword = "topic"

// GetTopicList  获取话题列表
func (d *Dao) GetTopicList(c context.Context) (resp []string, err error) {
	reply, err := d.titansCli.GetConfigByKeyword(c, &v1.GetConfigReq{Team: TitansTeam, Keyword: TitansKeyword})
	if err != nil {
		log.Error("main_member_GetIdentityStatus_error:%v", err)
		return
	}
	log.Info("main_member_GetIdentityStatus:%v", reply)
	resp = make([]string, 0)
	e := json.Unmarshal([]byte(reply.Value), &resp)
	if e != nil {
		log.Error("GetTopicList_json_error:%v,res=%v,", e, reply.Value)
		return
	}
	return
}
