package income

import (
	"context"
	"fmt"
	"time"

	model "go-common/app/admin/main/growup/model/income"

	"go-common/library/log"
	"go-common/library/xstr"
)

func (s *Service) lotteryStatis(c context.Context, categoryID []int64, from, to time.Time, groupType int) (data interface{}, err error) {
	query := ""
	if len(categoryID) != 0 {
		query = fmt.Sprintf("tag_id in (%s)", xstr.JoinInts(categoryID))
	}
	avs, err := s.GetArchiveIncome(c, _lottery, query, from.Format(_layout), to.Format(_layout))
	if err != nil {
		log.Error("s.GetArchiveIncome error(%v)")
		return
	}
	avsMap := make(map[string]*model.ArchiveStatis)
	for _, av := range avs {
		date := formatDateByGroup(av.Date.Time(), groupType)
		if val, ok := avsMap[date]; ok {
			val.Income += av.Income
			val.Avs++
		} else {
			avsMap[date] = &model.ArchiveStatis{
				Income: av.Income,
				Avs:    1,
			}
		}
	}
	data = parseArchiveStatis(avsMap, from, to, groupType)
	return
}
