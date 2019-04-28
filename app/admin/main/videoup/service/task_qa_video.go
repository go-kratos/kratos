package service

import (
	"context"
	"encoding/json"
	"strconv"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/log"
)

func (s *Service) fetchQAVideo(c context.Context, vp *archive.VideoParam) (task *archive.QAVideo, err error) {
	auditDetails, err := s.fetchVideoAuditDetail(c, vp.ID, vp.Aid, vp.Cid, vp.TaskID)
	if err != nil || auditDetails == nil {
		return nil, err
	}

	fans := int64(0)
	if auditDetails.UserInfo != nil && auditDetails.UserInfo["fans"] != nil {
		fanss := strconv.FormatFloat(auditDetails.UserInfo["fans"].(float64), 'f', 0, 64)
		if fans, err = strconv.ParseInt(fanss, 10, 64); err != nil {
			log.Error("fetchQAVideo strconv.ParseInt(%v) error(%v)", auditDetails.UserInfo["fans"], err)
			return
		}
	}

	video := auditDetails.Video
	details, err := json.Marshal(auditDetails)
	if err != nil {
		log.Error("fetchQAVideo json.Marshal(auditdetails) error(%v) aid(%d) cid(%d) details(%v)", err, vp.Aid, vp.Cid, auditDetails)
		return
	}

	auditSubmit := &archive.AuditSubmit{
		Encoding: strconv.Itoa(int(vp.Encoding)),
		Reason:   vp.Reason,
		ReasonID: strconv.Itoa(int(vp.RegionID)),
		Note:     vp.Note,
	}
	submit, err := json.Marshal(auditSubmit)
	if err != nil {
		log.Error("fetchQAVideo json.Marshal(auditsubmit) error(%v) aid(%d) cid(%d) submit(%s)", err, vp.Aid, vp.Cid, auditSubmit)
		return
	}

	task = &archive.QAVideo{
		UID:          vp.UID,
		Oname:        vp.Oname,
		AID:          vp.Aid,
		CID:          vp.Cid,
		TaskID:       vp.TaskID,
		TagID:        vp.TagID,
		ArcTitle:     video.Title,
		ArcTypeid:    video.Typeid,
		AuditStatus:  vp.Status,
		AuditSubmit:  string(submit),
		AuditDetails: string(details),
		MID:          video.MID,
		UPGroups:     s.getAllUPGroups(video.MID),
		Fans:         fans,
	}

	return
}

func (s *Service) fetchVideoAuditDetail(c context.Context, vid, aid, cid, taskID int64) (dt *archive.AuditDetails, err error) {
	video, err := s.arc.VideoInfo(c, aid, cid)
	if err != nil {
		return nil, err
	}
	if video == nil {
		log.Error("fetchVideoAuditDetail video not exist, aid(%d) cid(%d)", aid, cid)
		return
	}
	video.XcodeStateName = archive.XcodeStateNames[video.XcodeState]
	if tp, exist := s.typeCache[int16(video.Typeid)]; exist {
		video.Typename = tp.Name
	}
	video.Cover = coverURL(video.Cover)

	relationVideo, err := s.arc.VideoRelated(c, aid)
	if err != nil {
		return nil, err
	}

	task, err := s.arc.TaskDispatchByID(c, taskID)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.arc.GetUserCard(c, video.MID)
	if err != nil || len(userInfo) <= 0 {
		return nil, err
	}

	mosaic, err := s.arc.Mosaic(c, cid)
	if err != nil {
		return nil, err
	}

	watermark, err := s.arc.Watermark(c, video.MID)
	if err != nil {
		return nil, err
	}

	dt = &archive.AuditDetails{
		UserInfo:       userInfo,
		RelationVideos: relationVideo,
		Task:           []*archive.Task{task},
		Video:          video,
		Watermark:      watermark,
		Mosaic:         mosaic,
	}

	return
}

func (s *Service) addQAVideo(c context.Context, task *archive.QAVideo) (err error) {
	if task == nil {
		return
	}

	var bs []byte
	if bs, err = json.Marshal(task); err != nil {
		log.Error("addQAVideo json.Marshal error(%v) task(%+v)", err, task)
		return
	}
	err = s.arc.SendQAVideoAdd(c, bs)
	return
}
