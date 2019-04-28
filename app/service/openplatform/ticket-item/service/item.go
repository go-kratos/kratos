package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"fmt"
	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/time"
)

// Info get item info
func (s *ItemService) Info(ctx context.Context, in *item.InfoRequest) (reply *item.InfoReply, err error) {
	var items map[int64]*model.Item
	reply = new(item.InfoReply)
	if err = v.Struct(in); err != nil {
		err = ecode.RequestErr
		return
	}
	ids := []int64{in.ID}
	items, err = s.dao.Items(ctx, ids)
	if err != nil {
		log.Warn("service.Item Items() error(%v) or empty", err)
		return
	} else if len(items) == 0 {
		err = ecode.NothingFound
		return
	}
	i := items[in.ID]
	reply = &item.InfoReply{
		ID:         i.ID,
		Name:       i.Name,
		Status:     i.Status,
		Type:       i.Type,
		Rec:        i.Recommend,
		IsSale:     i.IsSale,
		TicketDesc: i.TicketDesc,
		PromTag:    i.PromoTags,
	}

	scs, err := s.dao.ScListByItem(ctx, ids)
	if err != nil {
		log.Error("service.Item ScListByItem() error(%v) or empty", err)
		return
	}
	tks, err := s.dao.TkListByItem(ctx, ids)
	if err != nil {
		log.Error("service.Item TkListByItem() error(%v) or empty", err)
		return
	}
	sctks := make(map[int64][]*model.TicketInfo)
	for _, tk := range tks[in.ID] {
		sctks[tk.ScreenID] = append(sctks[tk.ScreenID], tk)
	}

	detail, err := s.dao.ItemDetails(ctx, ids)
	if err != nil {
		log.Error("service.Item ItemDetails() error(%v) or empty", err)
		return
	}

	v, err := s.dao.Venues(ctx, []int64{i.VenueID})
	if err != nil {
		log.Error("service.Item Venues() error(%v) or empty", err)
		return
	} else if len(v) == 0 {
		err = ecode.NothingFound
		log.Error("service.Item Venues(%v) not found", i.VenueID)
		return
	}
	venue := v[i.VenueID]

	place, err := s.dao.Place(ctx, i.PlaceID)
	if err != nil {
		log.Error("service.Item Place() error(%v)", err)
		return
	} else if place == nil {
		err = ecode.NothingFound
		log.Error("service.Item Place(%v) not found", i.PlaceID)
		return
	}

	if sclist, ok := scs[in.ID]; ok {
		for _, sc := range sclist {
			reply.Screen = make(map[int64]*item.ScreenInfo)
			tklist := make(map[int64]*item.TicketInfo)
			if tkslist, ok := sctks[sc.ID]; ok {
				for _, tk := range tkslist {
					info := &item.TicketInfo{
						ID:           tk.ID,
						Desc:         tk.Desc,
						Type:         tk.Type,
						SaleType:     tk.SaleType,
						LinkSc:       tk.LinkSc,
						LinkTicketID: strconv.FormatInt(tk.LinkTicketID, 10),
						Symbol:       tk.Symbol,
						Color:        tk.Color,
						BuyLimit:     tk.BuyLimit,
						DescDetail:   tk.DescDetail,
						PriceList: &item.TicketPriceList{
							Price:    tk.Price,
							OriPrice: tk.OriginPrice,
							MktPrice: tk.MarketPrice,
						},
						StatusList: &item.TicketStatus{
							IsSale:    tk.IsSale,
							IsVisible: tk.IsVisible,
							IsRefund:  tk.IsRefund,
						},
						Time: &item.TicketTime{
							SaleStime: int64(tk.SaleStart),
							SaleEtime: int64(tk.SaleEnd),
						},
					}
					info.BuyNumLimit = new(item.TicketBuyNumLimit)
					tk.FormatTicketBuyLimit(info.BuyNumLimit)
					info.SaleFlag = tk.CalTkSaleFlag()
					tklist[tk.ID] = info
				}
			}
			reply.Screen[sc.ID] = &item.ScreenInfo{
				ID:     sc.ID,
				Name:   sc.Name,
				Status: sc.Status,
				ScTime: &item.ScreenTime{
					Stime: sc.StartTime,
					Etime: sc.EndTime,
				},
				Type:         sc.Type,
				TicketType:   sc.TicketType,
				ScreenType:   sc.ScreenType,
				DeliveryType: sc.DeliveryType,
				PickSeat:     sc.PickSeat,
				Ticket:       tklist,
			}
		}
	}
	reply.Img = &item.ImgList{
		First:  i.Img.First.URL,
		Banner: i.Img.Banner.URL,
	}
	reply.Time = &item.ItemTime{
		Stime: i.StartTime,
		Etime: i.EndTime,
	}
	var coor model.Coor
	json.Unmarshal([]byte(venue.Coordinate), &coor)
	reply.Ext = &item.ItemExt{
		Label:  i.Label,
		SpType: i.SponsorType,
		VerID:  i.VerID,
		Detail: detail[i.ID].PerformanceDesc,
		Venue: &item.VenueInfo{
			ID:     venue.ID,
			Name:   venue.Name,
			Status: venue.Status,
			AddrInfo: &item.VenueAddrInfo{
				Province:      venue.Province,
				City:          venue.City,
				District:      venue.District,
				AddressDetail: venue.AddressDetail,
				Traffic:       venue.Traffic,
				LonLatType:    coor.Type,
				LonLat:        coor.Coor,
			},
			PlaceInfo: &item.PlaceInfo{
				ID:      place.ID,
				Name:    place.Name,
				Status:  place.Status,
				BasePic: place.BasePic,
				DWidth:  place.DWidth,
				DHeight: place.DHeight,
			},
		},
	}
	reply.BillOpt = &item.BillOpt{
		BuyerInfo:  i.BuyerInfo,
		ExpTip:     i.ExpressFee,
		ExpFree:    i.HasExpressFee,
		VipExpFree: i.ExpressFreeFlag,
	}
	return reply, nil
}

// Cards get item cardlist
func (s *ItemService) Cards(ctx context.Context, in *item.CardsRequest) (reply *item.CardsReply, err error) {
	var vids []int64
	reply = new(item.CardsReply)
	items, err := s.dao.Items(ctx, in.IDs)
	if err != nil {
		log.Warn("service.Item Items() error(%v)", err)
		return
	} else if len(items) == 0 {
		err = ecode.NothingFound
		return
	}
	cm := make(map[int64]*item.CardReply)
	for id, i := range items {
		cm[id] = &item.CardReply{
			ID:         i.ID,
			Name:       i.Name,
			Status:     i.Status,
			Type:       i.Type,
			Rec:        i.Recommend,
			IsSale:     i.IsSale,
			TicketDesc: i.TicketDesc,
			PromTag:    i.PromoTags,
			Img: &item.ImgList{
				First:  i.Img.First.URL,
				Banner: i.Img.Banner.URL,
			},
			Time: &item.ItemTime{
				Stime: i.StartTime,
				Etime: i.EndTime,
			},
		}
		vids = append(vids, i.VenueID)
	}
	vs, err := s.dao.Venues(ctx, model.UniqueInt64(vids))
	if err != nil {
		log.Error("service.Item Venues() error(%v)", err)
		return
	}
	for id, i := range items {
		venue, ok := vs[i.VenueID]
		if !ok {
			log.Error("venue %v not found", i.VenueID)
			continue
		}
		var coor model.Coor
		json.Unmarshal([]byte(venue.Coordinate), &coor)
		cm[id].Venue = &item.VenueInfo{
			ID:     venue.ID,
			Name:   venue.Name,
			Status: venue.Status,
			AddrInfo: &item.VenueAddrInfo{
				Province:      venue.Province,
				City:          venue.City,
				District:      venue.District,
				AddressDetail: venue.AddressDetail,
				Traffic:       venue.Traffic,
				LonLatType:    coor.Type,
				LonLat:        coor.Coor,
			},
		}
	}
	reply.Cards = cm
	return
}

// BillInfo get item cardlist
func (s *ItemService) BillInfo(ctx context.Context, in *item.BillRequest) (reply *item.BillReply, err error) {
	reply = new(item.BillReply)
	reply.BaseInfo = make(map[int64]*item.ItemBase)
	reply.BillOpt = make(map[int64]*item.BillOpt)
	if err = v.Struct(in); err != nil {
		err = ecode.RequestErr
		return
	}
	if len(in.IDs) > 10 || len(in.ScIDs) > 20 || len(in.TkIDs) > 20 {
		err = ecode.RequestErr
		log.Error("len(in.IDs) > 10 or len(in.ScIDs) > 20 or len(in.TkIDs) > 20 ", err)
		return
	}
	items, err := s.dao.Items(ctx, in.IDs)
	if err != nil {
		log.Error("service.Item Items() error(%v)", err)
		return
	} else if len(items) == 0 {
		err = ecode.NothingFound
		return
	}
	scs, err := s.dao.ScList(ctx, in.ScIDs)
	if err != nil {
		log.Error("service.Item ScList() error(%v)", err)
		return
	}
	tks, err := s.dao.TkList(ctx, in.TkIDs)
	if err != nil {
		log.Error("service.Item TkList() error(%v)", err)
		return
	}
	for id, i := range items {
		reply.BaseInfo[id] = &item.ItemBase{
			ID:      i.ID,
			Name:    i.Name,
			Status:  i.Status,
			Type:    i.Type,
			IsSale:  i.IsSale,
			PromTag: i.PromoTags,
			VerID:   i.VerID,
			Time: &item.ItemTime{
				Stime: i.StartTime,
				Etime: i.EndTime,
			},
			Img: &item.ImgList{
				First:  i.Img.First.URL,
				Banner: i.Img.Banner.URL,
			},
		}
		reply.BillOpt[id] = &item.BillOpt{
			BuyerInfo:  i.BuyerInfo,
			ExpTip:     i.ExpressFee,
			ExpFree:    i.HasExpressFee,
			VipExpFree: i.ExpressFreeFlag,
		}
	}
	for id, sc := range scs {
		if _, ok := reply.BaseInfo[sc.ProjectID]; !ok {
			log.Error("screen %v not belong to any of item", id)
			continue
		}
		if reply.BaseInfo[sc.ProjectID].Screen == nil {
			reply.BaseInfo[sc.ProjectID].Screen = make(map[int64]*item.ScreenInfo)
		}
		reply.BaseInfo[sc.ProjectID].Screen[id] = &item.ScreenInfo{
			ID:     sc.ID,
			Name:   sc.Name,
			Status: sc.Status,
			ScTime: &item.ScreenTime{
				Stime: sc.StartTime,
				Etime: sc.EndTime,
			},
			Type:         sc.Type,
			TicketType:   sc.TicketType,
			ScreenType:   sc.ScreenType,
			DeliveryType: sc.DeliveryType,
			PickSeat:     sc.PickSeat,
		}
	}

	for id, tk := range tks {
		if _, ok := reply.BaseInfo[tk.ProjectID]; !ok {
			log.Error("ticket %v not belong to any of item", id)
			continue
		}
		if _, ok := reply.BaseInfo[tk.ProjectID].Screen[tk.ScreenID]; !ok {
			log.Error("ticket %v not belong to any of screen", id)
			continue
		}
		if reply.BaseInfo[tk.ProjectID].Screen[tk.ScreenID].Ticket == nil {
			reply.BaseInfo[tk.ProjectID].Screen[tk.ScreenID].Ticket = make(map[int64]*item.TicketInfo)
		}
		reply.BaseInfo[tk.ProjectID].Screen[tk.ScreenID].Ticket[id] = &item.TicketInfo{
			ID:           tk.ID,
			Desc:         tk.Desc,
			Type:         tk.Type,
			SaleType:     tk.SaleType,
			LinkSc:       tk.LinkSc,
			LinkTicketID: strconv.FormatInt(tk.LinkTicketID, 10),
			Symbol:       tk.Symbol,
			Color:        tk.Color,
			BuyLimit:     tk.BuyLimit,
			DescDetail:   tk.DescDetail,
			PriceList: &item.TicketPriceList{
				Price:    tk.Price,
				OriPrice: tk.OriginPrice,
				MktPrice: tk.MarketPrice,
			},
			StatusList: &item.TicketStatus{
				IsSale:    tk.IsSale,
				IsVisible: tk.IsVisible,
				IsRefund:  tk.IsRefund,
			},
			Time: &item.TicketTime{
				SaleStime: int64(tk.SaleStart),
				SaleEtime: int64(tk.SaleEnd),
			},
		}
		reply.BaseInfo[tk.ProjectID].Screen[tk.ScreenID].Ticket[id].BuyNumLimit = new(item.TicketBuyNumLimit)
		tk.FormatTicketBuyLimit(reply.BaseInfo[tk.ProjectID].Screen[tk.ScreenID].Ticket[id].BuyNumLimit)
	}
	return
}

// Wish 想去
func (s *ItemService) Wish(ctx context.Context, in *item.WishRequest) (reply *item.WishReply, err error) {
	if in.Face == "" {
		in.Face = s.c.URL.DefaultHead
	}

	reply = &item.WishReply{
		MID:    in.MID,
		ItemID: in.ItemID,
	}

	// TODO 改进点 - 先查有没有再写入
	wish := &model.UserWish{
		MID:    in.MID,
		Face:   in.Face,
		ItemID: in.ItemID,
	}
	if err = s.dao.AddWish(ctx, wish); err != nil {
		// 重复插入
		if strings.Contains(err.Error(), "1062") {
			log.Warn("d.AddWish(%+v) error(%v) duplicate entry", in, err)
			// TODO 返回 nil 还是 ecode
			err = nil
		} else {
			log.Error("s.Wish(%+v) error(%v)", in, err)
		}
		return
	}

	err = s.dao.WishCacheUpdate(ctx, wish)
	return
}

// Fav 收藏
func (s *ItemService) Fav(ctx context.Context, in *item.FavRequest) (reply *item.FavReply, err error) {
	reply = &item.FavReply{
		ItemID: in.ItemID,
		MID:    in.MID,
		Type:   in.Type,
		Status: in.Status,
	}
	if err = s.dao.FavUpdate(ctx, in.ItemID, in.MID, in.Type, in.Status); err != nil {
		log.Error("s.Fav() dao.FavUpdate(%d, %d, %d) error(%v)", in.ItemID, in.MID, in.Type, err)
		return
	}

	if err = s.dao.UserFavStateCache(ctx, in.ItemID, in.MID, in.Type, in.Status); err != nil {
		log.Error("s.Fav(%+v) d.UserFavState() error(%v)", in, err)
		return
	}

	if err = s.dao.UpdateFavListCache(ctx, in.ItemID, in.MID, in.Status); err != nil {
		return
	}

	return
}

// Version Add or Update Project with Version
func (s *ItemService) Version(c context.Context, info *item.VersionRequest) (ret *item.VersionReply, err error) {
	if paramErr := v.Struct(info); paramErr != nil {
		return ret, ecode.TicketMainInfoTooLarge
	}

	var pid int64
	var daoErr error
	if info.OpType == 0 {
		pid, daoErr = s.dao.AddProject(c, info.VerId)
	}

	return &item.VersionReply{ProjectId: pid}, daoErr
}

// BannerEdit Add or Update A Banner
func (s *ItemService) BannerEdit(c context.Context, info *item.BannerEditRequest) (ret *item.BannerEditReply, err error) {
	var bannerID int64
	var verID uint64
	if info.VerId == 0 {
		bannerID, verID, err = s.dao.AddBanner(c, info)
	} else {
		verID = info.VerId
		verInfo, _, getErr := s.dao.GetVersion(c, verID, false)
		if getErr != nil {
			return nil, getErr
		}
		verStatus := verInfo.Status
		bannerID = verInfo.TargetItem
		if verStatus == model.VerStatusNotReviewed {
			//草稿状态 随意编辑版本信息
			err = s.dao.EditBanner(c, info)
		} else if verStatus == model.VerStatusReadyForSale || verStatus == model.VerStatusOnShelf {
			//审核通过或进行中 仅新增版本
			var mainInfo []byte
			mainInfo, err = json.Marshal(info)
			if err != nil {
				log.Error("jsonMarshal版本详情失败:%s", err)
				return nil, ecode.TicketAddVersionFailed
			}
			str := fmt.Sprintf("%d%d%.2d", info.Position, info.SubPosition, info.Order)
			forInt, _ := strconv.ParseInt(str, 10, 64)
			err = s.dao.AddVersion(c, nil, &model.Version{
				Type:       model.VerTypeBanner,
				Status:     info.OpType,
				ItemName:   info.Name,
				TargetItem: bannerID,
				AutoPub:    1, // 自动上架
				PubStart:   time.Time(info.PubStart),
				PubEnd:     time.Time(info.PubEnd),
				For:        forInt,
			}, &model.VersionExt{
				Type:     model.VerTypeBanner,
				MainInfo: string(mainInfo),
			})
			if err != nil {
				log.Error("编辑创建banner版本失败: %s", err)
				return nil, ecode.TicketAddVersionFailed
			}
		} else {
			//其他状态 不可编辑
			return nil, ecode.TicketVerCannotEdit
		}

	}
	if info.OpType == 1 {
		//编辑操作为提交审核时 记入versionLog
		err = s.dao.AddVersionLog(c, &model.VersionLog{
			VerID: verID,
			Type:  2, //用户操作记录
			Log:   "提交审核",
			Uname: info.Uname,
		})
		if err != nil {
			return nil, err
		}
	}

	return &item.BannerEditReply{BannerId: bannerID}, err
}

// VersionReview Pass or Reject A Version
func (s *ItemService) VersionReview(c context.Context, info *item.VersionReviewRequest) (ret *item.VersionReviewReply, err error) {
	if info.OpType == model.VerReviewPass {
		//通过
		switch info.VerType {
		case model.VerTypeBanner:
			err = s.dao.PassOrPublishBanner(c, info.VerId)
			if err != nil {
				return nil, err
			}
		default:
			return nil, ecode.NothingFound
		}
	} else {
		//驳回
		_, err = s.dao.RejectVersion(c, info.VerId, info.VerType)
		if err != nil {
			return nil, err
		}
	}
	//审核操作记入versionLog
	err = s.dao.AddVersionLog(c, &model.VersionLog{
		VerID:  info.VerId,
		Type:   1, //用户操作记录
		Log:    info.Msg,
		Uname:  info.Uname,
		IsPass: info.OpType,
	})
	return &item.VersionReviewReply{VerId: info.VerId}, err
}

// VersionStatus Change A Version's Status
func (s *ItemService) VersionStatus(c context.Context, info *item.VersionStatusRequest) (ret *item.VersionStatusReply, err error) {

	switch info.VerType {
	case model.VerTypeBanner:
		err = s.dao.CgBannerStatus(c, info)
		return &item.VersionStatusReply{VerId: info.VerId}, err
	default:
		return nil, ecode.NothingFound
	}
}
