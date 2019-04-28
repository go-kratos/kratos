package service

import (
	"context"
	"encoding/json"
	"go-common/app/job/main/videoup/model/archive"
	"go-common/app/job/main/videoup/model/message"
	"go-common/app/job/main/videoup/model/redis"
	"go-common/library/log"
	"strconv"
	"time"
)

// sendBblog 发送粉丝动态databus
func (s *Service) sendBblog(a *archive.Result) {
	var (
		err    error
		c      = context.TODO()
		noPush = int64(a.AttrVal(archive.AttrNoPushBplus))
		//默认不展示
		show                           = int64(0)
		staffs                         []*archive.Staff
		dataPoi, dataVote, dataVoteFix []byte
	)
	if noPush == 0 && a.IsNormal() {
		show = int64(1)
	}
	if noPush == 1 {
		show = int64(2)
	}
	msg := &message.BlogCardMsg{
		Card: &archive.BlogCard{
			Type:    8,
			Rid:     a.Aid,
			OwnerID: a.Mid,
			Show:    show,
			Ts:      time.Now().Unix(),
			Dynamic: a.Dynamic,
		},
	}
	//lbs
	if dataPoi, err = s.arc.POI(context.TODO(), a.Aid); err != nil {
		log.Error("aid(%s)  s.videoupPub.SendBblog(%v) POI error(%v)", a.Aid, msg, err)
	}
	//vote
	if dataVote, err = s.arc.Vote(context.TODO(), a.Aid); err != nil {
		log.Error("aid(%s)  s.videoupPub.SendBblog(%v) Vote error(%v)", a.Aid, msg, err)
	}
	if dataPoi != nil || dataVote != nil {
		if dataVote != nil {
			var old *archive.VoteOld
			if err = json.Unmarshal(dataVote, &old); err != nil {
				log.Error("aid(%s)  s.videoupPub.SendBblog(%+v) Vote old  Unmarshal error(%v)", a.Aid, string(dataVote), err)
			}
			if dataVoteFix, err = json.Marshal(&archive.Vote{VoteID: old.VoteID, VoteTitle: old.VoteTitle}); err != nil {
				log.Error("aid(%s)  s.videoupPub.SendBblog(%+v) Vote new  Marshal error(%v)", a.Aid, old, err)
			}
		}
		var ext []byte
		if ext, err = json.Marshal(&archive.Ext{LBS: string(dataPoi), Vote: string(dataVoteFix)}); err != nil {
			log.Error("aid(%s)  s.videoupPub.SendBblog(%+v) Unmarshal error(%v)", a.Aid, msg, err)
		}
		msg.Card.Ext = string(ext)
	}
	//staffs
	if a.AttrVal(archive.AttrBitSTAFF) == archive.AttrYes {
		if staffs, err = s.arc.Staffs(c, a.Aid); err != nil {
			log.Error("aid(%s)  s.arc.Staffs (%+v)  error(%v)", a.Aid, msg, err)
		}
		var staffBox []*archive.StaffItem
		if staffs != nil && len(staffs) > 0 {
			for _, v := range staffs {
				item := &archive.StaffItem{Type: 1, UID: v.StaffMID}
				staffBox = append(staffBox, item)
			}
		}
		if staffBox != nil && len(staffBox) > 0 {
			msg.Card.Staffs = staffBox
			log.Info("aid(%d) SendBblog(%+v) staffs(%+v)", a.Aid, msg, msg.Card.Staffs)
		}
	}
	log.Info("aid(%d) start to send SendBblog(%+v) poi(%v) vote(%v) to databus", a.Aid, msg, string(dataPoi), string(dataVoteFix))
	k := strconv.FormatInt(a.Aid, 10)
	if err = s.blogPub.Send(c, k, msg); err != nil {
		s.syncRetry(c, a.Aid, a.Mid, redis.ActionForSendBblog, "", "")
		log.Error("aid(%s) %s s.videoupPub.SendBblog(%v) error(%v)", k, msg, err)
	}
}
