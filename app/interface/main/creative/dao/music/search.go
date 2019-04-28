package music

import (
	"context"
	"go-common/app/interface/main/creative/dao/tool"
	"go-common/app/interface/main/creative/model/search"
	"go-common/library/database/elastic"
	"go-common/library/ecode"
	"go-common/library/log"
)

// SearchBgmSIDs fn
func (d *Dao) SearchBgmSIDs(c context.Context, keyword string, pn, ps int) (ret []int64, page *search.Pager, err error) {
	retSIDs := make([]int64, 0)
	sres := &search.BgmResult{}
	r := d.es.NewRequest("archive_music").Fields(
		"sid",
	)
	r.Index("archive_music").Pn(1).Ps(5000).OrderScoreFirst(true)
	if len(keyword) > 0 { // "frontname", "name", "uname"
		r.WhereLike([]string{"music_frontname", "music_name", "uname"}, []string{keyword}, true, elastic.LikeLevelLow)
	}
	r.WhereEq("music_category_state", 0).WhereEq("music_state", 0).WhereEq("state", 0)
	r.Order("music_ctime", "desc") // 默认按入库时间倒序排
	log.Info("SearchBgmSIDs params(%s)", r.Params())
	if err = r.Scan(c, sres); err != nil {
		log.Error("SearchBgmSIDs r.Scan error(%v)", err)
		err = ecode.CreativeSearchErr
		return
	}
	if len(sres.Result) > 0 {
		for _, v := range sres.Result {
			retSIDs = append(retSIDs, v.SID)
		}
	}
	retSIDs = tool.DeDuplicationSlice(retSIDs)
	total := len(retSIDs)
	start := (pn - 1) * ps
	end := pn * ps
	page = &search.Pager{
		Num:   pn,
		Size:  ps,
		Total: total,
	}
	if total <= start {
		ret = make([]int64, 0)
	} else if total <= end {
		ret = retSIDs[start:total]
	} else {
		ret = retSIDs[start:end]
	}
	return
}

// ExtAidsWithSameBgm fn  获取这个sid最多的使用的aids列表,
func (d *Dao) ExtAidsWithSameBgm(c context.Context, sid int64, pn int) (retAIDs []int64, total int, err error) {
	retAIDs = make([]int64, 0)
	sres := &search.BgmExtResult{}
	r := d.es.NewRequest("archive_material").Fields("aid")
	r.Index("archive_material").Pn(1).Ps(pn).OrderScoreFirst(true)
	r.WhereEq("type", 3).WhereIn("data", sid).GroupBy(elastic.EnhancedModeDistinct, "aid", nil)
	r.Order("click", "desc")
	log.Info("ExtAidsWithSameBgm params(%s)", r.Params())
	if err = r.Scan(c, sres); err != nil {
		log.Error("ExtAidsWithSameBgm r.Scan error(%v)", err)
		err = ecode.CreativeSearchErr
		return
	}
	if len(sres.Result) > 0 {
		for _, v := range sres.Result {
			retAIDs = append(retAIDs, v.AID)
		}
	}
	if sres.Page != nil {
		total = sres.Page.Total
	}
	return
}
