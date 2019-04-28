package tag

import (
	"context"

	model "go-common/app/job/main/growup/model/tag"

	"go-common/library/log"
)

// tagUpRatio update tag up_charge_ratio
func (s *Service) tagUpRatio(c context.Context, date string, ctype int) (err error) {
	// delete
	err = s.delUpChargeRatio(c, ctype)
	if err != nil {
		log.Error("s.delUpChargeRatio error(%v)", err)
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

	ups, err := s.getTagUps(c, date, ctype, getArchives)
	if err != nil {
		log.Error("s.getTagUps error(%v)", err)
		return
	}
	log.Info("s.getTagUps insert ups(%d)", len(ups))

	// insert
	err = s.insertUpRatio(c, ups, ctype)
	if err != nil {
		log.Error("s.insertUpRatio error(%v)", err)
	}
	return
}

func (s *Service) getTagUps(c context.Context, date string, ctype int, getArchives IGetArchivesToTag) (ups map[int64]*model.AvTagRatio, err error) {
	ratioForAll := false
	tagInfo, err := s.dao.GetTagInfoByDate(c, 1, ctype, date, date)
	if err != nil {
		log.Error("s.dao.GetTagInfoByDate error(%v)", err)
		return
	}
	if len(tagInfo) == 0 {
		return
	}

	tagIDs := make([]int64, 0)
	for _, tag := range tagInfo {
		if tag.IsCommon == 1 {
			ratioForAll = true
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	categoryIDs, err := s.getArchiveCategory(c, ctype)
	if err != nil {
		log.Error("s.getArchiveCategory error(%v)", err)
		return
	}

	archives := make([]*model.ArchiveCharge, 0)
	if ratioForAll {
		archives, err = getArchives(c, date, categoryIDs)
		if err != nil {
			log.Error("s.getAvserror(%v)", err)
			return
		}
	}

	// 获取tag指定mid
	tagMID, err := s.GetTagUpInfoMID(c, tagIDs)
	if err != nil {
		log.Error("s.GetTagUpInfoMID error(%v)", err)
		return
	}

	ups = getUpsByTagInfo(archives, tagInfo, tagMID)
	return
}

func getUpsByTagInfo(archives []*model.ArchiveCharge, tagInfo []*model.TagInfo, tagMID map[int64][]int64) (upRatio map[int64]*model.AvTagRatio) {
	upRatio = make(map[int64]*model.AvTagRatio)

	chargeMIDs := make([]int64, 0)
	for _, archive := range archives {
		chargeMIDs = append(chargeMIDs, archive.MID)
	}
	for _, tag := range tagInfo {
		if tag.IsCommon == 0 {
			if mids, ok := tagMID[tag.ID]; ok {
				getUpTagRatio(upRatio, mids, tag)
			}
		} else {
			getUpTagRatio(upRatio, chargeMIDs, tag)
		}
	}
	return
}

func getUpTagRatio(upRatio map[int64]*model.AvTagRatio, mids []int64, tag *model.TagInfo) {
	for _, mid := range mids {
		if val, ok := upRatio[mid]; !ok {
			upRatio[mid] = &model.AvTagRatio{
				TagID:      tag.ID,
				AdjustType: tag.AdjustType,
				Ratio:      tag.Ratio,
				MID:        mid,
			}
		} else if tag.AdjustType == val.AdjustType && tag.Ratio > val.Ratio { // 调节方式相同，取大
			upRatio[mid].Ratio = tag.Ratio
			upRatio[mid].TagID = tag.ID
		} else if tag.AdjustType != val.AdjustType && tag.AdjustType == 1 { // 调节方式不同，取固定调节
			upRatio[mid].Ratio = tag.Ratio
			upRatio[mid].TagID = tag.ID
			upRatio[mid].AdjustType = tag.AdjustType
		}
	}
}
