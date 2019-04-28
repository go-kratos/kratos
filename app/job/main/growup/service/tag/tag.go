package tag

import (
	"context"
	"fmt"

	task "go-common/app/job/main/growup/service"

	"go-common/library/database/sql"
	"go-common/library/log"
)

// TagRatioAll cal tag ratio all
func (s *Service) TagRatioAll(c context.Context, date string) (err error) {
	defer func() {
		task.GetTaskService().SetTaskStatus(c, task.TaskTagRatio, date, err)
	}()

	// ready
	err = task.GetTaskService().TaskReady(c, date, task.TaskAvCharge, task.TaskCmCharge, task.TaskBgmSync)
	if err != nil {
		return
	}

	err = s.TagRatio(c, date, _video)
	if err != nil {
		log.Error("s.TagRatio video error(%v)", err)
		return
	}
	err = s.TagRatio(c, date, _column)
	if err != nil {
		log.Error("s.TagRatio column error(%v)", err)
		return
	}
	err = s.TagRatio(c, date, _bgm)
	if err != nil {
		log.Error("s.TagRatio bgm error(%v)", err)
	}
	return
}

// TagRatio handle tag ratio
func (s *Service) TagRatio(c context.Context, date string, ctype int) (err error) {
	err = s.tagArchiveRatio(c, date, ctype)
	if err != nil {
		log.Error("s.tagArchiveRatio error(%v)", err)
		return
	}

	err = s.tagUpRatio(c, date, ctype)
	if err != nil {
		log.Error("s.tagUpRatio error(%v)", err)
	}
	return
}

// TagIncomeAll tag income all
func (s *Service) TagIncomeAll(c context.Context, date string) (err error) {
	defer func() {
		task.GetTaskService().SetTaskStatus(c, task.TaskTagIncome, date, err)
	}()

	// ready
	err = task.GetTaskService().TaskReady(c, date, task.TaskCreativeIncome)
	if err != nil {
		return
	}

	err = s.TagIncome(c, date, _video)
	if err != nil {
		log.Error("s.TagIncome video error(%v)", err)
		return
	}
	err = s.TagIncome(c, date, _column)
	if err != nil {
		log.Error("s.TagIncome column error(%v)", err)
		return
	}
	err = s.TagIncome(c, date, _bgm)
	if err != nil {
		log.Error("s.TagIncome bgm error(%v)", err)
	}
	return
}

// TagIncome handle tag
func (s *Service) TagIncome(c context.Context, date string, ctype int) (err error) {
	var getArchives IGetArchivesToTag
	switch ctype {
	case _video:
		getArchives = s.getAvsAfterCal
	case _column:
		getArchives = s.getCmsAfterCal
	case _bgm:
		getArchives = s.getBgmBeforeCal
	}

	tagAvs, err := s.getTagArchive(c, date, ctype, getArchives)
	if err != nil {
		log.Error("s.getTagArchive error(%v)", err)
		return
	}

	tagUps, err := s.getTagUps(c, date, ctype, getArchives)
	if err != nil {
		log.Error("s.getTagAvs error(%v)", err)
		return
	}

	err = s.updateTagIncome(c, tagAvs, tagUps, date, ctype)
	if err != nil {
		log.Error("s.getActivityAvs error(%v)", err)
		return
	}

	err = s.addActivityInfo(c, date, tagAvs, ctype)
	if err != nil {
		log.Error("s.addActivityInfo error(%v)", err)
	}
	return
}

// TxUpdateTagInfoIncome update tag_info's income
func (s *Service) TxUpdateTagInfoIncome(tx *sql.Tx, tagIncome map[int64]int64) (err error) {
	for id, income := range tagIncome {
		if income <= 0 {
			continue
		}
		var rows int64
		rows, err = s.dao.TxUpdateTagInfoIncome(tx, id, income)
		if err != nil {
			tx.Rollback()
			log.Error("s.dao.UpdateTagInfoIncome error(%v)", err)
			return
		}
		if rows == 0 {
			tx.Rollback()
			return fmt.Errorf("TxUpdateTagInfoIncome error rows = 0")
		}
	}
	return
}
