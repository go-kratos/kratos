package dao

import (
	"context"
	notice "go-common/app/service/bbq/notice-service/api/v1"
	msg "go-common/app/service/bbq/sys-msg/api/v1"
	"go-common/library/log"
)

// 通知的业务类型
const (
	NoticeBizTypeSv      = 1
	NoticeBizTypeComment = 2
	NoticeBizTypeUser    = 3
	NoticeBizTypeSysMsg  = 4
)

// 通知类型
const (
	NoticeTypeLike    = 1
	NoticeTypeComment = 2
	NoticeTypeFan     = 3
	NoticeTypeSysMsg  = 4
)

const (
	_selectSQL = "select id, type, sender, receiver, jump_url, text, ctime, state from sys_msg where id > ? order by id asc"
)

// CreateNotice 创建通知
func (d *Dao) CreateNotice(ctx context.Context, notice *notice.NoticeBase) (err error) {
	_, err = d.noticeClient.CreateNotice(ctx, notice)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "create notice fail: notice="+notice.String()))
		return
	}

	log.V(10).Infov(ctx, log.KV("log", "create notice: notice="+notice.String()))
	return
}

// GetNewSysMsg 获取未被推送的系统消息
func (d *Dao) GetNewSysMsg(ctx context.Context, id int64) (list []*msg.SysMsg, err error) {

	rows, err := d.db.Query(ctx, _selectSQL, id)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "query mysql sys msg fail"))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var msg msg.SysMsg
		if err = rows.Scan(&msg.Id, &msg.Type, &msg.Sender, &msg.Receiver, &msg.JumpUrl, &msg.Text, &msg.Ctime, &msg.State); err != nil {
			log.Errorv(ctx, log.KV("log", "scan mysql sys msg fail"))
			return
		}
		list = append(list, &msg)
	}
	return
}
