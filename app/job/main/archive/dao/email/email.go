package email

import (
	"encoding/json"
	"fmt"
	"strconv"

	"go-common/app/job/main/archive/model/result"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"

	gomail "gopkg.in/gomail.v2"
)

var (
	_states = map[int]string{
		0:  "开放浏览",
		-1: "待审",
		-2: "打回稿件回收站",
		-3: "网警锁定删除",
		-4: "锁定稿件",
		// -5:   "锁定稿件开放浏览",
		-6: "修复待审",
		-7: "暂缓审核",
		// -8:   "补档待审",
		-9:  "等待转码",
		-10: "延迟发布",
		-11: "视频源待修",
		// -12:  "上传失败",
		-13: "允许评论待审",
		// -14:  "临时回收站",
		-15:  "分发中",
		-16:  "转码失败",
		-30:  "创建已提交",
		-40:  "UP主定时发布",
		-100: "UP主删除",
	}
)

func stateDescribe(state int) string {
	des, ok := _states[state]
	if ok {
		return des
	}
	return strconv.Itoa(state)
}

const (
	_bangumiMailSub = "番剧稿件《%s》状态变更"
	_movieMailSub   = "电影稿件《%s》状态变更"
	_mailBody       = `
稿件标题：%s
稿件状态：%s => %s
其他变更：%s => %s
稿件地址：http://www.bilibili.com/video/av%d
审核后台：http://manager.bilibili.co/#!/archive/modify/%d`
)

// PGCNotifyMail notify pgc mail
func (d *Dao) PGCNotifyMail(a *api.Arc, nw *result.Archive, old *result.Archive) {
	msg := gomail.NewMessage()
	msg.SetHeader("From", d.c.Mail.Username)
	switch {
	case a.AttrVal(archive.AttrBitIsBangumi) == archive.AttrYes:
		msg.SetHeader("To", d.c.Mail.Bangumi...)
		msg.SetHeader("Subject", fmt.Sprintf(_bangumiMailSub, a.Title))
	case a.AttrVal(archive.AttrBitIsMovie) == archive.AttrYes:
		msg.SetHeader("To", d.c.Mail.Movie...)
		msg.SetHeader("Subject", fmt.Sprintf(_movieMailSub, a.Title))
	default:
		return
	}
	obs, _ := json.Marshal(old)
	nbs, _ := json.Marshal(nw)

	if old.State != nw.State {
		oldState, newState := stateDescribe(old.State), stateDescribe(nw.State)
		msg.SetBody("text/plain", fmt.Sprintf(_mailBody, a.Title, oldState, newState, obs, nbs, a.Aid, a.Aid))
	} else {
		state := stateDescribe(int(a.State))
		msg.SetBody("text/plain", fmt.Sprintf(_mailBody, a.Title, state, state, obs, nbs, a.Aid, a.Aid))
	}
	if err := d.email.DialAndSend(msg); err != nil {
		log.Error("s.email.DialAndSend(aid: %d) error(%v)", a.Aid, err)
	}
}
