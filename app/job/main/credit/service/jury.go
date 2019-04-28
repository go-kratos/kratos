package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-common/app/job/main/credit/model"
	"go-common/library/log"
)

const (
	_msgTitle      = "风纪委员任期结束"
	_msgContext    = "您的风纪委员资格已到期，任职期间的总结报告已生成。感谢您对社区工作的大力支持！#{点击查看}{" + "\"http://www.bilibili.com/judgement/\"" + "}"
	_appealTitle   = "账号违规处理通知"
	_appealContent = `抱歉，你的账号因"在%s中%s"，现已进行%s处理，账号解封需要满足以下两个条件:1.账号封禁时间已满。2.完成解封答题（ #{点击进入解封答题}{"http://www.bilibili.com/blackroom/releaseexame.html"} ）全部完成后解封。封禁期间将无法投稿、发送及回复消息，无法发布评论、弹幕，无法对他人评论进行回复、赞踩操作，无法进行投币、编辑标签、添加关注、添加收藏操作。如对处罚有异议，可在7日内进行申诉，
	#{点击申诉}{"http://www.bilibili.com/judgement/appeal?bid=%d"} 。请遵守社区规范，共同维护良好的社区氛围！`
)

// Judge judge case.
func (s *Service) Judge(c context.Context, nwMsg []byte, oldMsg []byte) (err error) {
	var (
		sum        int64
		mr         = &model.Case{}
		bc         model.Case
		judge      int64
		status     int64
		yratio     int64
		nratio     int64
		judgeRadio = s.c.Judge.JudgeRadio
		voteMin    = s.c.Judge.CaseVoteMin
	)
	if err = json.Unmarshal(nwMsg, mr); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
		return
	}
	if mr.CaseType == model.JudeCaseTypePublic {
		return
	}
	if mr.Status != model.CaseStatusDealing {
		return
	}
	if bc, err = s.dao.CaseByID(c, mr.ID); err != nil {
		log.Error("s.dao.CaseByID error(%v)", err)
		return
	}
	if bc.Status == model.CaseStatusDealed || bc.Status == model.CaseStatusUndealed {
		return
	}
	sum = mr.Against + mr.Agree + mr.VoteDelete
	if voteMin <= 0 {
		log.Error("CaseVoteMin(%d) error(%v)", voteMin, err)
		return
	}
	if sum < voteMin {
		status = model.CaseStatusUndealed
		judge = model.JudgeTypeUndeal
		s.dao.UpdateCase(c, status, judge, mr.ID)
		return
	}
	yratio = mr.Agree * 100 / sum
	nratio = (mr.Against + mr.VoteDelete) * 100 / sum
	if judgeRadio <= 50 {
		log.Error("CaseJudgeRadio(%d) error(%v)", judgeRadio, err)
		return
	}
	if yratio >= judgeRadio {
		status = model.CaseStatusDealed
		judge = model.JudgeTypeLegal
	} else if nratio >= judgeRadio {
		status = model.CaseStatusDealed
		judge = model.JudgeTypeViolate
	} else {
		status = model.CaseStatusUndealed
		judge = model.JudgeTypeUndeal
	}
	s.dao.UpdateCase(c, status, judge, mr.ID)
	mr.Status = status
	mr.JudgeType = judge
	if status == model.CaseStatusDealed {
		s.BlockUser(c, mr)
		s.UpdateVoteCount(c, mr)
		s.dao.DelGrantCase(c, []int64{mr.ID})
	} else {
		if mr.OriginType == int64(model.OriginDM) {
			if mr.RelationID != "" {
				args := strings.Split(mr.RelationID, "-")
				if len(args) != 4 {
					return
				}
				s.dao.ReportDM(c, args[2], args[1], model.DMNotifyNotDel)
			}
		}
	}
	return
}

// BlockUser add user block.
func (s *Service) BlockUser(c context.Context, mr *model.Case) (err error) {
	if mr.JudgeType == model.JudgeTypeViolate {
		s.DelOrigin(c, mr)
		if mr.Against <= mr.VoteDelete+mr.Agree {
			err = s.dealMoralCase(c, mr)
			return
		}
		var (
			ok         bool
			punishType int64
		)
		forever, days := mr.BlockDays()
		if forever != model.InBlockedForever {
			punishType = int64(model.PunishTypeBlock)
		} else {
			punishType = int64(model.PunishTypeForever)
		}
		r := &model.BlockedInfo{
			UID:            mr.Mid,
			PunishType:     punishType,
			BlockedType:    model.PunishJury,
			OperatorName:   mr.Operator,
			CaseID:         mr.ID,
			Origin:         mr.Origin,
			OPID:           mr.OPID,
			BlockedForever: forever,
			BlockedDays:    days,
		}
		r.OriginContentModify = r.OriginContent
		if ok, err = s.CheckBlock(c, mr.Mid); err != nil || !ok {
			return
		}
		if mr.BusinessTime != model.DefaultTime {
			ok, _, err = s.jugeBlockedUser(c, mr.Mid, mr.BusinessTime, model.DealTimeTypeNone)
			if err != nil {
				log.Error("s.jugeBlockedUser(%d,%s,%d) error(%v)", mr.Mid, mr.BusinessTime, model.DealTimeTypeNone, err)
				return
			}
		} else {
			ok, _, err = s.jugeBlockedUser(c, mr.Mid, mr.Ctime, model.DealTimeTypeDay)
			if err != nil {
				log.Error("s.jugeBlockedUser(%d,%s,%d) error(%v)", mr.Mid, mr.Ctime, model.DealTimeTypeDay, err)
				return
			}
		}
		if ok {
			var id int64
			id, err = s.dao.AddBlockInfo(c, r, time.Now())
			if err != nil {
				log.Error("s.dao.AddBlockInfo error(%v)", err)
				return
			}
			if err = s.dao.BlockAccount(c, r); err != nil {
				log.Error("s.dao.BlockAccount(%+v) error(%v)", r, err)
				return
			}
			s.dao.SendMsg(c, mr.Mid, _appealTitle, fmt.Sprintf(_appealContent, model.OriginTypeDesc(int8(mr.OriginType)), model.ReasonTypeDesc(int8(mr.ReasonType)), model.BlockedDayDesc(int8(mr.BlockedDay)), id))
		}
	} else if mr.JudgeType == model.JudgeTypeLegal {
		if mr.JudgeType == int64(model.OriginDM) {
			s.dao.UpdatePunishResult(c, mr.ID, model.BlockNone)
			if mr.RelationID != "" {
				args := strings.Split(mr.RelationID, "-")
				if len(args) != 4 {
					return
				}
				s.dao.ReportDM(c, args[2], args[1], model.DMNotifyNotDel)
			}
		}
	}
	return
}

func (s *Service) dealMoralCase(c context.Context, mr *model.Case) (err error) {
	if err = s.dao.UpdatePunishResult(c, mr.ID, model.BlockOnlyDel); err != nil {
		log.Error("UpdatePunishResult error(%v)", err)
		return
	}
	title, content := model.OriginMsgContent(mr.OriginTitle, mr.OriginURL, mr.OriginContent, int8(mr.OriginType))
	for i := 0; i <= 5; i++ {
		if err := s.dao.AddMoral(c, mr.Mid, model.DefealtMoralVal, model.OrginMoralType(int8(mr.OriginType)), model.BUSSINESS, model.ReasonTypeDesc(int8(mr.ReasonType)), model.MoralRemark, ""); err != nil {
			continue
		}
		break
	}
	for i := 0; i <= 5; i++ {
		if err := s.dao.SendMsg(c, mr.Mid, title, content); err != nil {
			continue
		}
		break
	}
	return
}
