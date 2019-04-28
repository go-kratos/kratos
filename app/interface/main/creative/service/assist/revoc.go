package assist

import (
	"context"
	"strconv"

	"go-common/app/interface/main/creative/model/assist"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) revoc(c context.Context, assistLog *assist.AssistLog, ck, ip string) (err error) {
	switch {
	case assistLog.Type == 1 && assistLog.Action == 1:
		objectID, _ := strconv.ParseInt(assistLog.ObjectID, 10, 64)
		if objectID == 0 {
			err = ecode.RequestErr
			log.Error("strconv.ParseInt(%s) err(%v)", assistLog.ObjectID, err)
			return
		}
		if err = s.reply.ReplyRecover(c, assistLog.Mid, assistLog.SubjectID, objectID, ip); err != nil {
			log.Error("s.reply.ReplyRecover(%d,%d,%d,%s) err(%v)", assistLog.Mid, assistLog.SubjectID, assistLog.ObjectID, ip, err)
			return
		}
	case assistLog.Type == 2 && (assistLog.Action == 1 || assistLog.Action == 3):
		objectID, _ := strconv.ParseInt(assistLog.ObjectID, 10, 64)
		if err = s.dm.Edit(c, assistLog.Mid, assistLog.SubjectID, 0, []int64{objectID}, ip); err != nil {
			log.Error("s.dm.ResetDmStat(%d,%d,%d,%s) err(%v)", assistLog.Mid, assistLog.SubjectID, assistLog.ObjectID, ip, err)
			return
		}
	case assistLog.Type == 2 && assistLog.Action == 4:
		// 拉黑用户
		if err = s.dm.ResetUpBanned(c, assistLog.Mid, 0, assistLog.ObjectID, ip); err != nil {
			log.Error("s.dm.ResetUpBanned(%d,%d,%d,%s) err(%v)", assistLog.Mid, assistLog.SubjectID, assistLog.ObjectID, ip, err)
			return
		}
	case assistLog.Type == 2 && assistLog.Action == 5:
		// 移动弹幕到字幕池
		objectID, _ := strconv.ParseInt(assistLog.ObjectID, 10, 64)
		if err = s.dm.UpPool(c, assistLog.Mid, assistLog.SubjectID, []int64{objectID}, 0); err != nil {
			log.Error("s.dm.UpPool(%d,%d,%d,%s) err(%v)", assistLog.Mid, assistLog.SubjectID, assistLog.ObjectID, ip, err)
			return
		}
	case assistLog.Type == 2 && assistLog.Action == 6:
		// 忽略字幕池的弹幕
		objectID, _ := strconv.ParseInt(assistLog.ObjectID, 10, 64)
		if err = s.dm.UpPool(c, assistLog.Mid, assistLog.SubjectID, []int64{objectID}, 1); err != nil {
			log.Error("s.dm.UpPool(%d,%d,%d,%s) err(%v)", assistLog.Mid, assistLog.SubjectID, assistLog.ObjectID, ip, err)
			return
		}
	case assistLog.Type == 2 && assistLog.Action == 7:
		// 取消拉黑用户
		if err = s.dm.ResetUpBanned(c, assistLog.Mid, 1, assistLog.ObjectID, ip); err != nil {
			log.Error("s.dm.ResetUpBanned(%d,%d,%d,%s) err(%v)", assistLog.Mid, assistLog.SubjectID, assistLog.ObjectID, ip, err)
			return
		}
	case assistLog.Type == 3 && assistLog.Action == 8:
		// 取消拉黑用户
		if err = s.LiveRevocBanned(c, assistLog.Mid, assistLog.ObjectID, ck, ip); err != nil {
			log.Error("s.reply.LiveRevocBanned(%d,%d,%d,%s) err(%v)", assistLog.Mid, assistLog.SubjectID, assistLog.ObjectID, ip, err)
			return
		}
	}
	return
}
