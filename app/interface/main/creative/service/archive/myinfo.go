package archive

import (
	"context"
	"go-common/library/log"
)

// AllowCommercial fn
func (s *Service) AllowCommercial(c context.Context, mid int64) (ok int) {
	isOrder := s.AllowOrderUps(mid)
	if isOrder {
		log.Info("s.AllowOrderUps mid(%d)", mid)
		ok = 1
	}
	return
}

// AppModuleShowMap fn
// 编辑的时候暂时不允许重新开启动态抽奖或者修改动态抽奖
func (s *Service) AppModuleShowMap(mid int64, lotteryCheck bool) (ret map[string]bool) {
	ret = map[string]bool{
		"cooperate":       true,
		"vote":            true,
		"subtitle":        true,
		"filter":          true,
		"audio_record":    true,
		"camera":          true,
		"sticker":         true,
		"videoup_sticker": true,
		"theme":           true,
		"lottery":         false,
	}
	// define by dymc check
	if lotteryCheck {
		ret["lottery"] = true
	}
	// isWhite := false
	// for _, m := range s.c.Whitelist.ArcMids {
	// 	if m == mid {
	// 		isWhite = true
	// 		break
	// 	}
	// }
	// define for app module
	// if isWhite {
	// 	ret["theme"] = true
	// 	return
	// }
	// mod := mid % 100
	// if mod <= 20 {
	// 	ret["theme"] = true
	// }
	return
}
