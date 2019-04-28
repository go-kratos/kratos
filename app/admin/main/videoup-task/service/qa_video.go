package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/admin/main/videoup-task/dao"
	"go-common/app/admin/main/videoup-task/model"
	accmdl "go-common/app/service/main/account/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

//GetVideoList list qa video tasks
func (s *Service) GetVideoList(ctx context.Context, pm *model.ListParams) (list *model.QAVideoList, err error) {
	var (
		listLen     int
		detailMap   map[int64]map[string]int64
		midList     []int64
		users       map[int64]*model.UserRole
		infos       map[int64]*accmdl.Info
		upGroupList map[int64][]*model.UPGroup
	)

	if list, err = s.searchQAVideo(ctx, pm); err != nil || list == nil || len(list.Result) <= 0 {
		return
	}

	ids := make([]int64, listLen)
	uids := make([]int64, listLen)
	for _, item := range list.Result {
		ids = append(ids, item.ID)
		uids = append(uids, item.UID)
	}
	if detailMap, midList, err = s.dao.QAVideoDetail(ctx, ids); err != nil {
		return
	}
	if users, err = s.dao.GetUsernameAndRole(ctx, uids); err != nil {
		return
	}
	if infos, err = s.dao.AccountInfos(ctx, midList); err != nil {
		return
	}
	//获取列表页获取
	if upGroupList, err = s.dao.UPGroups(ctx, midList); err != nil {
		return
	}

	for _, item := range list.Result {
		item.User = users[item.UID]
		item.StateName = model.QAStates[item.State]
		item.UPName = ""
		dt, exist := detailMap[item.ID]
		if exist {
			item.DetailID = dt["detail_id"]
			item.TaskUTime = dt["task_utime"]
			item.MID = dt["mid"]
			item.UPGroupList = upGroupList[item.MID]
			if infos[item.MID] != nil {
				item.UPName = infos[item.MID].Name
			}
		}
	}

	return
}

//AddQATaskVideo add a qa video task
func (s *Service) AddQATaskVideo(ctx context.Context, detail *model.AddVideoParams) (taskID int64, err error) {
	var vid int64
	if vid, err = s.dao.GetVID(ctx, detail.AID, detail.CID); err != nil {
		log.Error("AddQATaskVideo s.dao.GetVID(aid(%d), cid(%d)) error(%v)", detail.AID, detail.CID, err)
		return
	}
	if vid <= 0 {
		log.Error("AddQATaskVideo non-deleted video(aid(%d), cid(%d)) not exist", detail.AID, detail.CID)
		return
	}

	taskID, err = s.insertVideoTask(ctx, detail)
	return
}

func (s *Service) insertVideoTask(ctx context.Context, detail *model.AddVideoParams) (taskID int64, err error) {
	var (
		tx       *sql.Tx
		detailID int64
	)
	defer func() {
		if msg := recover(); msg != nil {
			if tx != nil {
				tx.Rollback()
			}
			log.Error("insertVideoTask panic recover, msg(%s)", msg)
		}
	}()

	if tx, err = s.dao.BeginTran(ctx); err != nil {
		return
	}

	if detailID, err = s.dao.InsertQAVideo(tx, &detail.VideoDetail); err != nil {
		tx.Rollback()
		return
	}

	if taskID, err = s.dao.InTaskQA(tx, detail.OUID, detailID, model.QATypeVideo); err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		dao.PromeErr("arcdb: commit", "insertVideoTask commit error(%v) aid(%d) cid(%d)", err, detail.AID, detail.CID)
	}
	return
}

func (s *Service) getQATaskVideo(ctx context.Context, id int64, simple bool) (task *model.QATaskVideo, err error) {
	if simple {
		task, err = s.dao.QATaskVideoSimpleByID(ctx, id)
	} else {
		task, err = s.dao.QATaskVideoByID(ctx, id)
	}

	if err != nil {
		return
	}
	if task == nil {
		err = ecode.NothingFound
		return
	}

	task.GetAttributeList()
	return
}

//GetDetail get qa video task detail
func (s *Service) GetDetail(ctx context.Context, id int64) (dt *model.TaskVideoDetail, err error) {
	var (
		taskVideo *model.QATaskVideo
		video     *model.Video
		history   []*model.VideoOperInfo
	)
	if taskVideo, err = s.getQATaskVideo(ctx, id, false); err != nil {
		return
	}
	info := &model.VideoTaskInfo{
		QATaskVideo: *taskVideo,
	}
	groups, _ := s.dao.UPGroups(ctx, []int64{taskVideo.MID})
	info.UPGroupList = groups[taskVideo.MID]
	info.GetWarnings()

	if video, err = s.getVideo(ctx, taskVideo.AID, taskVideo.CID); err != nil {
		return
	}
	if history, err = s.getVideoOperInfo(ctx, video.ID); err != nil {
		return
	}

	dt = &model.TaskVideoDetail{
		Task:         info,
		Video:        video,
		VideoHistory: history,
	}
	return
}

//QAVideoSubmit submit qa video task
func (s *Service) QAVideoSubmit(ctx context.Context, username string, uid int64, vp *model.QASubmitParams) (err error) {
	var (
		task  *model.QATaskVideo
		video *model.Video
	)
	if task, err = s.getQATaskVideo(ctx, vp.ID, true); err != nil {
		log.Error("QAVideoSubmit s.arc.QATaskVideoByID error(%v), id(%d) task(%+v)", err, vp.ID, task)
		return
	}
	if video, err = s.getVideo(ctx, task.AID, task.CID); err != nil {
		log.Error("sendLog s.getVideo error(%v) qa.id(%d) aid(%d) cid(%d)", err, vp.ID, task.AID, task.CID)
		return
	}
	//不重复质检
	if task.State == model.QAStateFinish {
		return
	}

	//更新task
	task.State = model.QAStateFinish
	task.FTime = time.Now()
	if _, err = s.dao.UpTask(ctx, vp.ID, task.State, task.FTime); err != nil {
		return
	}

	s.dao.AddVideoOper(ctx, task.AID, uid, video.ID, video.Attribute, video.Status, 0, fmt.Sprintf("一审任务质检TAG: [%s]", vp.QATag), vp.QaNote)
	s.sendLog(ctx, username, uid, video, task, vp)
	return
}

func (s *Service) sendLog(ctx context.Context, username string, uid int64, video *model.Video, task *model.QATaskVideo, vp *model.QASubmitParams) (err error) {
	var (
		note        string
		taskUIDName string
	)

	if task == nil || len(task.AttributeList) == 0 {
		log.Error("sendLog task/task.attributelist not exist, task(%+v) params(%+v)", task, vp)
		return
	}
	if note, err = task.GetNote(); err != nil {
		log.Error("sendLog task.GetNote() error(%v), params(%+v)", err, vp)
		return
	}
	if vp.Norank == 1 {
		video.AttributeList["norank"] = 1
	}
	if vp.Nodynamic == 1 {
		video.AttributeList["nodynamic"] = 1
	}
	if vp.Norecommend == 1 {
		video.AttributeList["norecommend"] = 1
	}
	if vp.Nosearch == 1 {
		video.AttributeList["nosearch"] = 1
	}
	if vp.PushBlog == 1 {
		video.AttributeList["push_blog"] = 1
	}
	if vp.OverseaBlock == 1 {
		video.AttributeList["oversea_block"] = 1
	}
	video.TagID = vp.TagID
	video.Status = vp.AuditStatus
	video.Note = vp.Note
	video.Reason = vp.Reason
	video.Encoding = vp.Encoding
	if taskUIDNames, err := s.dao.GetUsername(ctx, []int64{task.UID}); err != nil {
		taskUIDName = ""
		err = nil
	} else {
		taskUIDName = taskUIDNames[task.UID]
	}

	content := map[string]interface{}{
		"audit_status": task.AuditStatus,
		"audit_attr":   task.AttributeList,
		"audit_tag_id": task.TagID,
		"audit_note":   note,
		"qa_status":    video.Status,
		"qa_attr":      video.AttributeList,
		"qa_tag_id":    video.TagID,
		"qa_note":      video.Note,
	}

	data := &report.ManagerInfo{
		Uname:    username,
		UID:      uid,
		Business: model.LogQATask,
		Type:     model.LogQATaskVideo,
		Oid:      vp.ID,
		Action:   strconv.Itoa(int(video.Status)),
		Ctime:    task.FTime,
		Index:    []interface{}{task.TaskID, task.MID, task.CTime.Unix(), strconv.FormatInt(vp.QaTagID, 10), strconv.FormatInt(task.UID, 10), taskUIDName},
		Content:  content,
	}
	report.Manager(data)
	log.Info(" sendLog data(%+v)", data)
	return
}

//UpVideoUTime update qa video task utime
func (s *Service) UpVideoUTime(ctx context.Context, aid, cid, taskID, utime int64) (err error) {
	var id int64
	if id, err = s.dao.GetQAVideoID(ctx, aid, cid, taskID); err != nil {
		log.Error("UpVideoUTime s.dao.GetQAVideoID error(%v) aid(%d) cid(%d) taskid(%d) utime(%d)", err, aid, cid, taskID, utime)
		return
	}
	if id <= 0 {
		log.Error("UpVideoUTime s.dao.GetQAVideoID not found aid(%d) cid(%d) taskid(%d) utime(%d)", aid, cid, taskID, utime)
		err = ecode.RequestErr
		return
	}

	return s.dao.UpdateQAVideoUTime(ctx, aid, cid, taskID, utime)
}

func (s *Service) delProc() {
	var (
		err                     error
		qaVideoRows, qaTaskRows int64 = 1, 1
		limit                         = 1000
	)

	for {
		deadLine := time.Now().AddDate(0, -1, 0)
		for {
			if qaVideoRows > 0 {
				if qaVideoRows, err = s.dao.DelQAVideo(context.TODO(), deadLine, limit); err != nil {
					log.Error("delProc s.dao.DelQAVideo(%v,%d) error(%v)", deadLine, limit, err)
				}
			}
			if qaTaskRows > 0 {
				if qaTaskRows, err = s.dao.DelQATask(context.TODO(), deadLine, limit); err != nil {
					log.Error("delProc s.dao.DelQATask(%v,%d) error(%v)", deadLine, limit, err)
				}
			}

			if qaVideoRows+qaTaskRows == 0 {
				break
			}
			time.Sleep(time.Minute)
		}

		time.Sleep(time.Hour * 24)
	}
}
