package view

import (
	"context"
	"strconv"

	"go-common/app/interface/main/app-card/model"
	"go-common/app/interface/main/app-card/model/card/operate"
)

func (s *Service) searchFollow(c context.Context, platform, mobiApp, device, buvid string, build int, mid, vmid int64) (follow *operate.Card, err error) {
	const _title = "关注TA的也关注了"
	ups, trackID, err := s.search.Follow(c, platform, mobiApp, device, buvid, build, mid, vmid)
	if err != nil {
		return
	}
	items := make([]*operate.Card, 0, len(ups))
	for _, up := range ups {
		if up.Mid != 0 {
			item := &operate.Card{ID: up.Mid, Goto: model.GotoMid, Param: strconv.FormatInt(up.Mid, 10), URI: strconv.FormatInt(up.Mid, 10), Desc: up.RecReason}
			items = append(items, item)
		}
	}
	if len(items) < 3 {
		return
	}
	id, _ := strconv.ParseInt(trackID, 10, 64)
	if id < 1 {
		return
	}
	follow = &operate.Card{ID: id, Param: trackID, Items: items, Title: _title, CardGoto: model.CardGotoSearchUpper}
	return
}
