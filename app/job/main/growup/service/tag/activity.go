package tag

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	model "go-common/app/job/main/growup/model/tag"
	"go-common/library/log"
	xtime "go-common/library/time"
)

func (s *Service) addActivityInfo(c context.Context, date string, tagAvs map[int64]*model.AvTagRatio, ctype int) (err error) {
	tagAvMap := make(map[int64]bool)
	for _, av := range tagAvs {
		tagAvMap[av.AvID] = true
	}

	tagInfo, err := s.dao.GetTagInfoByDate(c, 0, ctype, date, date)
	if err != nil {
		log.Error("s.dao.GetTagInfoByDate error(%v)", err)
		return
	}

	activityAvs, err := s.getActivityArchives(c, tagInfo, ctype)
	if err != nil {
		log.Error("s.getActivityArchives error(%v)", err)
		return
	}

	avsExists, err := s.dao.ListActivityInfo(c)
	if err != nil {
		log.Error("s.ListActivityInfo error(%v)", err)
		return
	}

	inList := make([]*model.ArchiveCharge, 0)
	for _, av := range activityAvs {
		if tagAvMap[av.AID] && !avsExists[av.AID] {
			inList = append(inList, av)
		}
	}
	log.Info("get %d avs, need to insert %d avs", len(activityAvs), len(inList))
	err = s.insertActivityInfo(c, inList)
	if err != nil {
		log.Error("s.insertActivityInfo error(%v)", err)
	}
	return
}

func (s *Service) insertActivityInfo(c context.Context, avs []*model.ArchiveCharge) (err error) {
	var buf bytes.Buffer
	for _, a := range avs {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(a.AID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.ActivityID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.CategoryID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(a.TagID, 10))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values := buf.String()
	buf.Reset()
	_, err = s.dao.InsertActivityInfo(c, values)
	return
}

func (s *Service) getActivityArchives(c context.Context, tagInfo []*model.TagInfo, ctype int) (archives []*model.ArchiveCharge, err error) {
	archives = make([]*model.ArchiveCharge, 0)
	if ctype == _bgm {
		return
	}
	activityMap := make(map[int64]struct{})
	activityIDs := make([]int64, 0)
	tagMap := make(map[int64]map[int64]int64) // map[activityID][categoryID]tagID
	for _, tag := range tagInfo {
		if tag.ActivityID == 0 {
			continue
		}
		if _, ok := activityMap[tag.ActivityID]; !ok {
			activityMap[tag.ActivityID] = struct{}{}
			activityIDs = append(activityIDs, tag.ActivityID)
			tagMap[tag.ActivityID] = make(map[int64]int64)
		}
		tagMap[tag.ActivityID][tag.CategoryID] = tag.ID
	}
	if len(activityIDs) == 0 {
		return
	}
	activityInfos, err := s.getActivityInfos(c, activityIDs, ctype)
	if err != nil {
		log.Error("s.GetActivityInfos error(%v)", err)
		return
	}

	// 获取一级分区和二级分区的对应关系 key:二级分区
	categoryMap := make(map[int64]int64)
	if ctype == _video {
		categoryMap, err = s.dao.GetVideoTypes(c)
		if err != nil {
			log.Error("s.dao.GetAllTypes error(%v)", err)
			return
		}
	} else if ctype == _column {
		categoryMap, err = s.dao.GetColumnTypes(c)
		if err != nil {
			log.Error("s.dao.GetColumnTypes error(%v)", err)
			return
		}
	}

	archives, err = convertArchiveCharge(activityInfos, categoryMap, tagMap)
	if err != nil {
		log.Error("convertArchiveCharge error(%v)", err)
		return
	}
	return
}

func convertArchiveCharge(activityInfo []*model.ActivityInfo, categoryMap map[int64]int64, tagMap map[int64]map[int64]int64) (archives []*model.ArchiveCharge, err error) {
	archives = make([]*model.ArchiveCharge, 0)
	for _, act := range activityInfo {
		var uploadTime time.Time
		uploadTime, err = time.Parse("2006-01-02 15:04:05", act.CDate)
		if err != nil {
			return
		}
		av := &model.ArchiveCharge{
			AID:        act.AvID,
			MID:        act.MID,
			CategoryID: categoryMap[act.TypeID],
			ActivityID: act.ActivityID,
			UploadTime: xtime.Time(uploadTime.Unix()),
		}
		if t, ok := tagMap[av.ActivityID]; ok {
			av.TagID = t[av.CategoryID]
		}
		archives = append(archives, av)
	}
	return
}

func (s *Service) getActivityInfos(c context.Context, activityID []int64, ctype int) ([]*model.ActivityInfo, error) {
	switch ctype {
	case _video:
		return s.getVideoActivity(c, activityID)
	case _column:
		return s.getCmActivity(c, activityID)
	}
	return nil, fmt.Errorf("getActivityInfos ctype error")
}

func (s *Service) getVideoActivity(c context.Context, activityID []int64) (activityInfos []*model.ActivityInfo, err error) {
	activityInfos = make([]*model.ActivityInfo, 0)
	start, offset := 0, 20
	if len(activityID) < offset {
		offset = len(activityID)
	}
	for start+offset <= len(activityID) {
		var act []*model.ActivityInfo
		act, err = s.getVideoActivityInfos(c, activityID[start:start+offset])
		if err != nil {
			return
		}
		if len(act) != 0 {
			activityInfos = append(activityInfos, act...)
		}
		start += offset
		if start < len(activityID) && start+offset > len(activityID) {
			offset = len(activityID) - start
		}
	}
	return
}

func (s *Service) getVideoActivityInfos(c context.Context, activityID []int64) (acts []*model.ActivityInfo, err error) {
	acts = make([]*model.ActivityInfo, 0)
	page, size := 1, 30
	for {
		var act []*model.ActivityInfo
		act, err = s.dao.GetVideoActivityInfo(c, activityID, page, size)
		if err != nil {
			return
		}
		if len(act) == 0 {
			break
		}
		acts = append(acts, act...)
		page++
		// qps控制
		if page%50 == 0 {
			time.Sleep(1 * time.Second)
		}
	}
	return
}

func (s *Service) getCmActivity(c context.Context, activityID []int64) (activityInfos []*model.ActivityInfo, err error) {
	activityInfos = make([]*model.ActivityInfo, 0)
	for _, id := range activityID {
		var act []*model.ActivityInfo
		act, err = s.getCmActivityInfo(c, id)
		if err != nil {
			return
		}
		if len(act) != 0 {
			activityInfos = append(activityInfos, act...)
		}
	}
	return
}

func (s *Service) getCmActivityInfo(c context.Context, id int64) (acts []*model.ActivityInfo, err error) {
	acts = make([]*model.ActivityInfo, 0)
	page, size := 1, 30
	for {
		var act []*model.ActivityInfo
		act, err = s.dao.GetCmActivityInfo(c, id, page, size)
		if err != nil {
			return
		}
		if len(act) == 0 {
			break
		}
		acts = append(acts, act...)
		page++
		// qps控制
		if page%50 == 0 {
			time.Sleep(1 * time.Second)
		}
	}
	return
}
