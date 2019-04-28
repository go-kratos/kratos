package service

import (
	"context"
	"fmt"
	"go-common/app/admin/main/creative/model/logcli"
	"go-common/app/admin/main/creative/model/music"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

//SendNotify .
func (s *Service) SendNotify(c *bm.Context, sendIds []int64, data map[int64]*music.SidNotify) (err error) {
	var (
		//mid首次收录
		content = "您的音频稿件（au:%d）已被手机投稿BGM库收录，快加入素材激励计划(#{APP申请入口}{\"https://member.bilibili.com/studio/up-allowance-h5#\"},#{WEB申请入口}{\"https://member.bilibili.com/v/#/allowance\"})获取收益吧！被收录稿件名：《%s》"
		//sid首次收录
		content2 = "您的音频稿件【《%s》（au:%d）】已被手机投稿BGM库收录，期待您创作更多优秀的新作品哦"
		title    = "创作激励计划素材收录通知"
	)
	for _, sid := range sendIds {
		if _, ok := data[sid]; !ok {
			continue
		}
		sendConfig := data[sid]
		log.Info("svc.SendNotify param sendConfig(%+v) ", sendConfig)
		var (
			mids        []int64
			first, send bool
			sendContent string
		)
		//check exists
		exists := music.Music{}
		if err = s.DBArchive.Model(&music.Music{}).Where("sid=?", sid).First(&exists).Error; err != nil {
			continue
		}
		//每个mid 第一次收录 优先级最高
		if sendConfig.MidFirst {
			first = true
			send = true
		}
		if !sendConfig.MidFirst && sendConfig.SidFirst {
			first = false
			send = true
		}
		if !first {
			content = content2
			sendContent = fmt.Sprintf(content, exists.Name, exists.Sid)
		} else {
			sendContent = fmt.Sprintf(content, exists.Sid, exists.Name)
		}
		if !send {
			return
		}
		mids = []int64{exists.Mid}

		s.addAsyn(func() {
			if err = s.dao.MutliSendSysMsg(context.TODO(), mids, title, sendContent); err != nil {
				log.Error("s.d.MutliSendSysMsg(%s,%s,%s) error(%+v)", xstr.JoinInts(mids), title, sendContent, err)
				return
			}
		})
		s.SendMusicLog(c, logcli.LogClientArchiveMusicTypeCategoryRelation, &music.LogParam{ID: sid, UID: 0, UName: fmt.Sprintf("mid(%d)", exists.Mid), Action: "SendNotify", Name: sendContent})
	}
	return
}
