package tag

import (
	"context"
	"fmt"

	model "go-common/app/job/main/growup/model/tag"
	"go-common/library/log"
)

// updateTagIncome update tag_income
func (s *Service) updateTagIncome(c context.Context, avTagRatio, upTagRatio map[int64]*model.AvTagRatio, date string, ctype int) (err error) {
	log.Info("GET %d av tag ratio", len(avTagRatio))
	log.Info("GET %d up tag ratio", len(upTagRatio))
	query := fmt.Sprintf("date = '%s'", date)
	archiveIncome, err := s.getArchiveIncome(c, query, ctype)
	if err != nil {
		log.Error("s.getArchiveIncome error(%v)", err)
		return
	}
	log.Info("GET %d archive_income", len(archiveIncome))

	upIncome, err := s.getUpIncome(c, query)
	if err != nil {
		log.Error("s.GetUpIncome error(%v)", err)
		return
	}
	log.Info("GET %d up_income", len(upIncome))

	// 获取之前有标签收入的稿件, 计算稿件活动期间总收入
	// preTagIncome, err := s.GetUpTagIncomeMap(c)
	// if err != nil {
	// 	log.Error("s.GetUpTagIncome error(%v)", err)
	// 	return
	// }

	// 获取有标签收入的稿件，每个标签的当日总收入
	tagAvs := make([]*model.AvTagRatio, 0)
	tagIncome := make(map[int64]int64)
	// 稿件维度
	for _, av := range archiveIncome {
		if val, ok := avTagRatio[av.AID]; ok {
			val.Income = av.Income
			val.BaseIncome = av.BaseIncome
			val.TotalIncome = av.TotalIncome
			val.TaxMoney = av.TaxMoney
			val.Date = date

			tagAvs = append(tagAvs, val)
			tagIncome[val.TagID] += av.Income
		}
	}

	// up主维度
	for _, up := range upIncome {
		if val, ok := upTagRatio[up.MID]; ok {
			switch ctype {
			case _video:
				up.Income, up.BaseIncome, up.TotalIncome, up.TaxMoney = up.AvIncome, up.AvBaseIncome, up.AvTotalIncome, up.AvTax
			case _column:
				up.Income, up.BaseIncome, up.TotalIncome, up.TaxMoney = up.ColumnIncome, up.ColumnBaseIncome, up.ColumnTotalIncome, up.ColumnTax
			case _bgm:
				up.Income, up.BaseIncome, up.TotalIncome, up.TaxMoney = up.BgmIncome, up.BgmBaseIncome, up.BgmTotalIncome, up.BgmTax
			}
			t := &model.AvTagRatio{
				Income:      up.Income,
				BaseIncome:  up.BaseIncome,
				TotalIncome: up.TotalIncome,
				TaxMoney:    up.TaxMoney,
				Date:        date,
				AvID:        0,
				MID:         up.MID,
				TagID:       val.TagID,
			}

			tagAvs = append(tagAvs, t)
			tagIncome[val.TagID] += up.Income
		}
	}

	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}

	log.Info("Insert(update) %d avs to up_tag_income", len(tagAvs))
	err = s.TxInsertUpTagIncome(tx, tagAvs)
	if err != nil {
		log.Error("s.InsertUpTagIncome error(%v)", err)
		return
	}

	err = s.TxUpdateTagInfoIncome(tx, tagIncome)
	if err != nil {
		log.Error("s.UpdateTagInfo error(%v)", err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error")
		return
	}
	return
}

func (s *Service) getArchiveIncome(c context.Context, query string, ctype int) (archives []*model.ArchiveIncome, err error) {
	var id int64
	limit := 2000
	for {
		var archive []*model.ArchiveIncome
		archive, err = s.dao.GetArchiveIncome(c, id, query, limit, ctype)
		if err != nil {
			return
		}
		archives = append(archives, archive...)
		if len(archive) < limit {
			break
		}
		id = archive[len(archive)-1].ID
	}

	if ctype == _bgm {
		bgms := make(map[int64]*model.ArchiveIncome)
		for _, archive := range archives {
			if b, ok := bgms[archive.AID]; ok {
				b.Income += archive.Income
				b.BaseIncome += archive.BaseIncome
				b.TaxMoney += archive.TaxMoney
			} else {
				bgms[archive.AID] = archive
			}
		}
		archives = make([]*model.ArchiveIncome, 0)
		for _, b := range bgms {
			archives = append(archives, b)
		}
	}
	return
}

// getUpIncome get up_income by query
func (s *Service) getUpIncome(c context.Context, query string) (ups []*model.UpIncome, err error) {
	var id int64
	limit := 2000
	for {
		var up []*model.UpIncome
		up, err = s.dao.GetUpIncome(c, id, query, limit)
		if err != nil {
			return
		}
		ups = append(ups, up...)
		if len(up) < limit {
			break
		}
		id = up[len(up)-1].ID
	}
	return
}

func (s *Service) getAvIncomeStatis(c context.Context) (avs []*model.ArchiveCharge, err error) {
	var id int64
	limit := 2000
	for {
		var av []*model.ArchiveCharge
		av, err = s.dao.GetAvIncomeStatis(c, id, limit)
		if err != nil {
			return
		}
		avs = append(avs, av...)
		if len(av) < limit {
			break
		}
		id = av[len(av)-1].ID
	}
	return
}

func (s *Service) getCmIncomeStatis(c context.Context) (cms []*model.ArchiveCharge, err error) {
	var id int64
	limit := 2000
	for {
		var cm []*model.ArchiveCharge
		cm, err = s.dao.GetCmIncomeStatis(c, id, limit)
		if err != nil {
			return
		}
		cms = append(cms, cm...)
		if len(cm) < limit {
			break
		}
		id = cm[len(cm)-1].ID
	}
	return
}
