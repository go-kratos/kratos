package service

import (
	"fmt"

	"go-common/app/admin/main/tv/model"
	"go-common/library/log"
)

// UPlayurl get ugc play url
func (s *Service) UPlayurl(aid int) (playurl string, err error) {
	w := map[string]interface{}{
		"deleted": 0,
		"result":  1,
		"cid":     aid,
	}
	video := model.Video{}
	if err = s.dao.DB.Model(&model.Video{}).Where(w).First(&video).Error; err != nil {
		err = fmt.Errorf("找不到aid为%d的数据", aid)
		return
	}
	if playurl, err = s.dao.UPlayurl(ctx, aid); err != nil {
		log.Error("UPlayurl API Error(%d) (%v)", aid, err)
		return
	}
	log.Info("UPlayurl aid = %d, playurl = %s", aid, playurl)
	return
}
