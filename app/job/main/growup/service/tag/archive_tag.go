package tag

import (
	"context"
	"fmt"
	"strings"

	model "go-common/app/job/main/growup/model/tag"

	"go-common/library/log"
)

var (
	_video  = 1
	_column = 2
	_bgm    = 3
)

// IGetArchivesToTag get avs to tag
type IGetArchivesToTag func(c context.Context, date string, categoryIDs []int64) (archives []*model.ArchiveCharge, err error)

func (s *Service) tagArchiveRatio(c context.Context, date string, ctype int) (err error) {
	// delete
	err = s.delArchiveRatio(c, ctype)
	if err != nil {
		log.Error("s.delArchiveRatio error(%v)", err)
		return
	}

	var getArchives IGetArchivesToTag
	switch ctype {
	case _video:
		getArchives = s.getAvsBeforeCal
	case _column:
		getArchives = s.getCmsBeforeCal
	case _bgm:
		getArchives = s.getBgmBeforeCal
	}

	ratios, err := s.getTagArchive(c, date, ctype, getArchives)
	if err != nil {
		log.Error("s.getTagArchive error(%v)", err)
		return
	}

	// insert
	err = s.insertAvRatio(c, ratios, ctype)
	if err != nil {
		log.Error("s.InsertAvRatio error(%v)", err)
	}
	return

}

func (s *Service) getTagArchive(c context.Context, date string, ctype int, getArchives IGetArchivesToTag) (ratios map[int64]*model.AvTagRatio, err error) {
	tagInfo, err := s.dao.GetTagInfoByDate(c, 0, ctype, date, date)
	if err != nil {
		log.Error("s.dao.GetTagInfoByDate error(%v)", err)
		return
	}
	if len(tagInfo) == 0 {
		return
	}

	// get tag_id and category_id
	tagIDs := make([]int64, 0)
	categoryIDs := make([]int64, 0)
	categoryMap := make(map[int64]struct{})
	for _, tag := range tagInfo {
		tagIDs = append(tagIDs, tag.ID)
		if tag.ActivityID != 0 {
			continue
		}
		if _, ok := categoryMap[tag.CategoryID]; !ok {
			categoryMap[tag.CategoryID] = struct{}{}
			categoryIDs = append(categoryIDs, tag.CategoryID)
		}
	}

	archives, err := getArchives(c, date, categoryIDs)
	if err != nil {
		log.Error("s.getAvs error(%v)", err)
		return
	}
	log.Info("GET matched archives(%d)", len(archives))

	// get archives by activityID
	actArchives, err := s.getActivityArchives(c, tagInfo, ctype)
	if err != nil {
		log.Error("s.getActivityAvs error(%v)", err)
		return
	}
	log.Info("GET matched avtivity archive(%d)", len(actArchives))

	if len(actArchives) > 0 {
		archives = append(archives, actArchives...)
	}
	if len(archives) <= 0 {
		return
	}

	// 获取tag指定mid
	tagMID, err := s.GetTagUpInfoMID(c, tagIDs)
	if err != nil {
		log.Error("s.GetTagUpInfoMID error(%v)", err)
		return
	}

	ratios = getRatioByTagInfo(archives, tagInfo, tagMID)
	return
}

// get avs from av_daily_charge
func (s *Service) getAvsBeforeCal(c context.Context, date string, categoryIDs []int64) (archives []*model.ArchiveCharge, err error) {
	archives = make([]*model.ArchiveCharge, 0)
	dateS := strings.Split(date, "-")
	for _, id := range categoryIDs {
		var avc []*model.ArchiveCharge
		query := fmt.Sprintf("tag_id = %d AND date = '%s'", id, date)
		avc, err = s.getAvDailyCharge(c, dateS[1], query)
		if err != nil {
			log.Error("s.getAvDailyCharge error(%v)", err)
			return
		}
		archives = append(archives, avc...)
	}
	return
}

// get avs from av_income_statis
func (s *Service) getAvsAfterCal(c context.Context, date string, categoryIDs []int64) (archives []*model.ArchiveCharge, err error) {
	return s.getAvIncomeStatis(c)
}

// get columns from column_daily_charge
func (s *Service) getCmsBeforeCal(c context.Context, date string, categoryIDs []int64) (archives []*model.ArchiveCharge, err error) {
	archives = make([]*model.ArchiveCharge, 0)
	for _, id := range categoryIDs {
		var cm []*model.ArchiveCharge
		query := fmt.Sprintf("tag_id = %d AND date = '%s'", id, date)
		cm, err = s.getCmDailyCharge(c, query)
		if err != nil {
			log.Error("s.getAvDailyCharge error(%v)", err)
			return
		}
		archives = append(archives, cm...)
	}
	return
}

// get cms from column_income_statis
func (s *Service) getCmsAfterCal(c context.Context, date string, categoryIDs []int64) (archives []*model.ArchiveCharge, err error) {
	return s.getCmIncomeStatis(c)
}

// get bgm from backgroud_music
func (s *Service) getBgmBeforeCal(c context.Context, date string, categoryIDs []int64) (archives []*model.ArchiveCharge, err error) {
	return s.getBackgroundMusic(c, 2000)
}

func getRatioByTagInfo(archives []*model.ArchiveCharge, tagInfo []*model.TagInfo, tagMID map[int64][]int64) (avRatio map[int64]*model.AvTagRatio) {
	avRatio = make(map[int64]*model.AvTagRatio)
	midMap := make(map[int64]int64)

	categoryAvs := make(map[int64][]*model.ArchiveCharge)        // 非活动 category
	midCategoryAvs := make(map[string][]*model.ArchiveCharge)    // 非活动 category+mid
	actCategoryAvs := make(map[string][]*model.ArchiveCharge)    // 活动 activity+category
	actMidCategoryAvs := make(map[string][]*model.ArchiveCharge) // 活动 activity+category+mid

	for _, av := range archives {
		// mid和avid对应map
		midMap[av.AID] = av.MID

		// 活动id需要特殊处理
		if av.ActivityID != 0 {
			// 活动 不指定mid
			key := fmt.Sprintf("act%d+cate%d", av.ActivityID, av.CategoryID)
			if _, ok := actCategoryAvs[key]; ok {
				actCategoryAvs[key] = append(actCategoryAvs[key], av)
			} else {
				actCategoryAvs[key] = []*model.ArchiveCharge{av}
			}

			// 活动 指定mid
			mckey := fmt.Sprintf("act%d+cate%d+mid%d", av.ActivityID, av.CategoryID, av.MID)
			if _, ok := actMidCategoryAvs[mckey]; ok {
				actMidCategoryAvs[mckey] = append(actMidCategoryAvs[mckey], av)
			} else {
				actMidCategoryAvs[mckey] = []*model.ArchiveCharge{av}
			}

		} else {
			// 非活动 不指定mid
			if _, ok := categoryAvs[av.CategoryID]; ok {
				categoryAvs[av.CategoryID] = append(categoryAvs[av.CategoryID], av)
			} else {
				categoryAvs[av.CategoryID] = []*model.ArchiveCharge{av}
			}

			// 非活动 指定mid
			mcKey := fmt.Sprintf("mid%d+cate%d", av.MID, av.CategoryID)
			if _, ok := midCategoryAvs[mcKey]; ok {
				midCategoryAvs[mcKey] = append(midCategoryAvs[mcKey], av)
			} else {
				midCategoryAvs[mcKey] = []*model.ArchiveCharge{av}
			}
		}
	}

	for _, tag := range tagInfo {
		// 非活动 不指定mid
		if tag.IsCommon == 1 && tag.ActivityID == 0 {
			if av, ok := categoryAvs[tag.CategoryID]; ok {
				getAvTagRatio(avRatio, av, tag, midMap)
			}
		}

		// 活动 不指定mid
		if tag.IsCommon == 1 && tag.ActivityID != 0 {
			key := fmt.Sprintf("act%d+cate%d", tag.ActivityID, tag.CategoryID)
			if av, ok := actCategoryAvs[key]; ok {
				getAvTagRatio(avRatio, av, tag, midMap)
			}
		}

		// 指定mid
		if tag.IsCommon == 0 {
			if mids, ok := tagMID[tag.ID]; ok {
				for _, mid := range mids {
					// 非活动
					if tag.ActivityID == 0 {
						mcKey := fmt.Sprintf("mid%d+cate%d", mid, tag.CategoryID)
						if av, ok := midCategoryAvs[mcKey]; ok {
							getAvTagRatio(avRatio, av, tag, midMap)
						}
					} else { // 活动
						mcKey := fmt.Sprintf("act%d+cate%d+mid%d", tag.ActivityID, tag.CategoryID, mid)
						if av, ok := actMidCategoryAvs[mcKey]; ok {
							getAvTagRatio(avRatio, av, tag, midMap)
						}
					}
				}
			}
		}
	}
	return
}

func getAvTagRatio(avRatio map[int64]*model.AvTagRatio, avs []*model.ArchiveCharge, tag *model.TagInfo, midMap map[int64]int64) {
	for _, av := range avs {
		if av.UploadTime < tag.UploadStartTime || int64(av.UploadTime) >= tag.UploadEndTime.Time().AddDate(0, 0, 1).Unix() {
			continue
		}
		if val, ok := avRatio[av.AID]; !ok {
			avRatio[av.AID] = &model.AvTagRatio{
				AvID:       av.AID,
				TagID:      tag.ID,
				AdjustType: tag.AdjustType,
				Ratio:      tag.Ratio,
				MID:        midMap[av.AID],
			}
		} else if tag.AdjustType == val.AdjustType && tag.Ratio > val.Ratio { // 调节方式相同，取大
			avRatio[av.AID].Ratio = tag.Ratio
			avRatio[av.AID].TagID = tag.ID
		} else if tag.AdjustType != val.AdjustType && tag.AdjustType == 1 { // 调节方式不同，取固定调节
			avRatio[av.AID].Ratio = tag.Ratio
			avRatio[av.AID].TagID = tag.ID
			avRatio[av.AID].AdjustType = tag.AdjustType
		}
	}
}
