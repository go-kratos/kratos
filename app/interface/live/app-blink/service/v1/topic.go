package v1

import (
	"context"

	v1pb "go-common/app/interface/live/app-blink/api/http/v1"
	"go-common/app/interface/live/app-blink/conf"
	"go-common/app/interface/live/app-blink/dao"
	"go-common/library/ecode"
)

// TopicService struct
type TopicService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *dao.Dao
}

//NewTopicService init
func NewTopicService(c *conf.Config) (s *TopicService) {
	s = &TopicService{
		conf: c,
		dao:  dao.New(c),
	}
	return s
}

// GetTopicList implementation
// 获取话题列表
// `method:"GET" midware:"auth"`
func (s *TopicService) GetTopicList(ctx context.Context, req *v1pb.GetTopicListReq) (resp *v1pb.GetTopicListResp, err error) {
	reply, err := s.dao.GetTopicList(ctx)
	if err != nil {
		err = ecode.CallResourceError
		return
	}
	resp = &v1pb.GetTopicListResp{}
	resp.TopicList = reply
	return
}

// CheckTopic implementation
// 检验话题是否有效
// `method:"GET" midware:"auth"`
func (s *TopicService) CheckTopic(ctx context.Context, req *v1pb.CheckTopicReq) (resp *v1pb.CheckTopicResp, err error) {
	msg := req.Topic
	area := "live_biaoti"
	reply, err := s.dao.CheckMsgIsLegal(ctx, msg, area, 0)
	if err != nil {
		err = ecode.Error(ecode.CallMainFilterError, "出错啦，再试试吧")
		return
	}
	resp = &v1pb.CheckTopicResp{}
	if reply {
		err = ecode.Error(ecode.FILTERNOTPASS, "当前输入的话题里面包含敏感内容~请修改~")
		return
	}
	return
}
