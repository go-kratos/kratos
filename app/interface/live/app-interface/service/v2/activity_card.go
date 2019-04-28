package v2

import (
	"context"

	v2pb "go-common/app/interface/live/app-interface/api/http/v2"
)

// getActivityCard 活动模块
func (s *IndexService) getActivityCard(ctx context.Context) (resp []*v2pb.MActivityCard) {
	resp = []*v2pb.MActivityCard{}
	ids := s.getIdsFromModuleMap(ctx, []int64{_activityType})
	if len(ids) <= 0 {
		return
	}
	err, activityCardMap := s.roomDao.GetActivityCard(ctx, ids, "GetAllList")
	if err != nil {
		return
	}
	listMap := make(map[int64][]*v2pb.ActivityCardItem)
	for i, ac := range activityCardMap {
		respAc := &v2pb.ActivityCardItem{Room: []*v2pb.RoomCardItem{}, Av: []*v2pb.AvCardItem{}}
		respAc.Card = &v2pb.BannerCardItem{
			Aid:        ac.Card.Aid,
			Pic:        ac.Card.Pic,
			Title:      ac.Card.Title,
			Text:       ac.Card.Text,
			PicLink:    ac.Card.PicLink,
			GoLink:     ac.Card.GoLink,
			ButtonText: ac.Card.ButtonText,
			Status:     ac.Card.Status,
			Sort:       ac.Card.Sort,
		}
		if len(ac.Room) > 0 {
			for _, room := range ac.Room {
				roomCard := &v2pb.RoomCardItem{
					IsLive:         room.IsLive,
					RoomId:         room.Roomid,
					Title:          room.Title,
					UName:          room.Uname,
					Online:         room.Online,
					Cover:          room.Cover,
					AreaV2ParentId: room.AreaV2ParentId,
					AreaV2Id:       room.AreaV2Id,
					Sort:           room.Sort,
				}
				respAc.Room = append(respAc.Room, roomCard)
			}
		}
		if len(ac.Av) > 0 {
			for _, av := range ac.Av {
				avCard := &v2pb.AvCardItem{
					Avid:      av.Avid,
					Title:     av.Title,
					ViewCount: av.ViewCount,
					DanMaKu:   av.Danmaku,
					Duration:  av.Duration,
					Cover:     av.Cover,
					Sort:      av.Sort,
				}
				respAc.Av = append(respAc.Av, avCard)
			}
		}
		listMap[i] = append(listMap[i], respAc)
	}
	moduleInfoMap := s.getAllModuleInfoMap(ctx)
	for _, m := range moduleInfoMap[_activityType] {
		if l, ok := listMap[m.Id]; ok {
			resp = append(resp, &v2pb.MActivityCard{
				ModuleInfo: m,
				List:       l,
			})
		}
	}
	return
}
