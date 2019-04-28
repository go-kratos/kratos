package dao

import (
	"context"

	"encoding/json"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"strconv"

	"github.com/jinzhu/gorm"
)

// ProjectMainInfo 项目版本内容
type ProjectMainInfo struct {
	Name            string         `json:"name"`
	Type            string         `json:"type"`
	Size            string         `json:"size"`
	ProvinceID      string         `json:"province_id"`
	CityID          string         `json:"city_id"`
	DistrictID      string         `json:"district_id"`
	VenueID         string         `json:"venue_id"`
	PlaceID         string         `json:"place_id"`
	StartTime       int32          `json:"start_time"`
	EndTime         int32          `json:"end_time"`
	Docs            []ImgInfo      `json:"docs"`
	Screens         []Screen       `json:"screens"`
	TicketsSingle   []TicketSingle `json:"tickets_single"`
	TicketsPass     []TicketPass   `json:"tickets_pass"`
	TicketsAllPass  []TicketPass   `json:"tickets_allpass"`
	PerformanceImg  PerformanceImg `json:"performance_image"`
	PerformanceDesc string         `json:"performance_desc"`
	SellingProp     string         `json:"selling_prop"`
	TagIDs          []string       `json:"tag_ids"`
	GuestIDs        []int64        `json:"guest_ids"`
	GuestImgs       []GuestImg     `json:"guest_imgs"`
	CompID          string         `json:"comp_id"`
	Label           string         `json:"label"`
	SponsorType     string         `json:"sponsor_type"`
}

// ImgInfo 图片信息
type ImgInfo struct {
	URL  string `json:"url"`
	Desc string `json:"desc"`
}

// GuestImg 嘉宾图片信息
type GuestImg struct {
	ID     int64  `json:"id"`
	ImgURL string `json:"img_url"`
}

// Screen 场次信息
type Screen struct {
	ScreenName      string `json:"screen_name"`
	ScreenStartTime int32  `json:"screen_start_time"`
	ScreenEndTime   int32  `json:"screen_end_time"`
	ScreenType      string `json:"screen_type"`
	PickSeat        int32  `json:"pick_seat"`
	TicketType      string `json:"ticket_type"`
	DeliveryType    string `json:"delivery_type"`
	ScreenID        string `json:"screen_id"`
}

// TicketSingle 单场票信息
type TicketSingle struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Color       string   `json:"color"`
	BuyLimit    string   `json:"buy_limit"`
	PayValue    int64    `json:"pay_value"`
	PayMethod   string   `json:"pay_method"`
	Desc        string   `json:"descrp"`
	TicketID    string   `json:"ticket_sc_id"`
	BuyLimitNum []string `json:"buy_limit_num"`
}

// TicketPass 通票信息
type TicketPass struct {
	Name        string   `json:"name"`
	LinkTicket  int32    `json:"link_ticket"`
	LinkScreens []int32  `json:"link_screens"`
	Color       string   `json:"color"`
	BuyLimit    string   `json:"buy_limit"`
	PayValue    int64    `json:"pay_value"`
	PayMethod   string   `json:"pay_method"`
	Desc        string   `json:"descrp"`
	TicketID    string   `json:"ticket_sc_id"`
	BuyLimitNum []string `json:"buy_limit_num"`
}

// PerformanceImg 项目图片信息
type PerformanceImg struct {
	First  ImgInfo `json:"first"`
	Banner ImgInfo `json:"banner"`
}

const (
	// BuyNumLimit 默认限购值
	BuyNumLimit = `{"per":8}`
)

// AddProject 项目初始化
func (d *Dao) AddProject(c context.Context, verID uint64) (pid int64, err error) {

	var verExtInfo model.VersionExt
	if dbErr := d.db.Where("ver_id = ?", verID).First(&verExtInfo).Error; dbErr != nil {
		log.Error("获取项目版本详情失败:%s", dbErr)
		return 0, dbErr
	}

	decodedMainInfo := d.GetDefaultMainInfo()
	err = json.Unmarshal([]byte(verExtInfo.MainInfo), &decodedMainInfo)
	if err != nil {
		return
	}

	// 开启事务
	tx := d.db.Begin()

	// 创建新project
	var project model.Item
	if project, err = d.CreateProject(c, tx, verID, decodedMainInfo); err != nil {
		return
	}
	pid = project.ID
	if pid == 0 {
		tx.Rollback()
		log.Error("pid not initialized")
		return 0, ecode.TicketPidIsEmpty
	}

	// 创建新项目详情在project_extra表
	if err = d.CreateProjectExtInfo(c, tx, pid, decodedMainInfo.PerformanceDesc); err != nil {
		return
	}

	// 创建场次screen
	scIDList := make(map[int32]int64)
	scStartTimes := make(map[int32]int32)
	scEndTimes := make(map[int32]int32)
	var screen model.Screen
	for k, v := range decodedMainInfo.Screens {
		screenType, _ := strconv.ParseInt(v.ScreenType, 10, 64)
		ticketType, _ := strconv.ParseInt(v.TicketType, 10, 64)
		deliveryType, _ := strconv.ParseInt(v.DeliveryType, 10, 64)
		screen = model.Screen{
			ProjectID:    pid,
			Name:         v.ScreenName,
			StartTime:    v.ScreenStartTime,
			EndTime:      v.ScreenEndTime,
			Type:         int32(screenType),
			TicketType:   int32(ticketType),
			DeliveryType: int32(deliveryType),
			PickSeat:     v.PickSeat,
			ScreenType:   1, // 单场票场次
		}
		if screen, err = d.CreateOrUpdateScreen(c, tx, screen); err != nil {
			return
		}
		// 获取的screen自增id复制到maininfo内
		decodedMainInfo.Screens[k].ScreenID = strconv.FormatInt(screen.ID, 10)
		scIDList[int32(k)] = screen.ID
		scStartTimes[int32(k)] = screen.StartTime
		scEndTimes[int32(k)] = screen.EndTime
	}

	// 创建票价
	TkSingleIDList := make(map[int32]int64)
	TkSingleTypeList := make(map[int32]int32)
	var tkPrice model.TicketPrice
	for k, v := range decodedMainInfo.TicketsSingle {
		saleType, _ := strconv.ParseInt(v.Type, 10, 64)
		buyLimit, _ := strconv.ParseInt(v.BuyLimit, 10, 64)
		payMethod, _ := strconv.ParseInt(v.PayMethod, 10, 64)
		ticketID, baseErr := model.GetTicketIDFromBase()
		if ticketID == 0 || baseErr != nil {
			tx.Rollback()
			log.Error("baseCenter获取ticketID失败 ticketid:%s,baseErr:%s", ticketID, baseErr)
			return 0, baseErr
		}
		tkPrice = model.TicketPrice{
			ID:            ticketID,
			ProjectID:     pid,
			Desc:          v.Name,
			Type:          1, // 单场票
			SaleType:      int32(saleType),
			Color:         v.Color,
			BuyLimit:      int32(buyLimit),
			PaymentMethod: int32(payMethod),
			PaymentValue:  v.PayValue,
			DescDetail:    v.Desc,
			IsSale:        1,   // 可售
			IsRefund:      -10, // 不可退
			OriginPrice:   -1,  // 未設置
			MarketPrice:   -1,
			SaleStart:     TimeNull, // 0000-00-00 00:00:00
			SaleEnd:       TimeNull,
		}
		if tkPrice, err = d.CreateOrUpdateTkPrice(c, tx, tkPrice, 0); err != nil {
			return
		}

		//票价限购
		limitData := d.FormatByPrefix(v.BuyLimitNum, "buy_limit_")
		if err = d.CreateOrUpdateTkPriceExtra(c, tx, limitData, ticketID, pid); err != nil {
			return
		}

		// 获取的ticketPrice自增id复制到mainInfo内
		decodedMainInfo.TicketsSingle[k].TicketID = strconv.FormatInt(tkPrice.ID, 10)
		TkSingleIDList[int32(k)] = ticketID
		TkSingleTypeList[int32(k)] = int32(saleType)

	}

	var passScID, allPassScID int64
	// 创建通票场次
	passScID, err = d.GetOrUpdatePassSc(c, tx, pid, decodedMainInfo.TicketsPass, scStartTimes, scEndTimes, scTypePass, 0)
	if err != nil {
		return
	}
	decodedMainInfo.TicketsPass, err = d.InsertOrUpdateTkPass(c, tx, pid, passScID, decodedMainInfo.TicketsPass, TkTypePass, scIDList, TkSingleIDList, TkSingleTypeList)
	if err != nil {
		return
	}
	// 创建联票场次
	allPassScID, err = d.GetOrUpdatePassSc(c, tx, pid, decodedMainInfo.TicketsAllPass, scStartTimes, scEndTimes, scTypeAllPass, 0)
	if err != nil {
		return
	}
	decodedMainInfo.TicketsAllPass, err = d.InsertOrUpdateTkPass(c, tx, pid, allPassScID, decodedMainInfo.TicketsAllPass, TkTypeAllPass, scIDList, TkSingleIDList, TkSingleTypeList)
	if err != nil {
		return
	}

	// 创建标签
	for _, tagID := range decodedMainInfo.TagIDs {
		err = d.CreateTag(c, tx, pid, tagID)
		if err != nil {
			return
		}
	}

	// 创建嘉宾
	guestImgMap := make(map[int64]string)
	// 组合嘉宾对应头像map
	if decodedMainInfo.GuestIDs != nil {
		if decodedMainInfo.GuestImgs == nil {
			var guestInfoList []model.Guest
			if err = d.db.Select("id, guest_img").Where("id IN (?)", decodedMainInfo.GuestIDs).Find(&guestInfoList).Error; err != nil {
				log.Error("获取嘉宾信息失败:%s", err)
				tx.Rollback()
				return
			}
			for _, v := range guestInfoList {
				guestImgMap[v.ID] = v.GuestImg
			}
		} else {
			for _, v := range decodedMainInfo.GuestImgs {
				guestImgMap[v.ID] = v.ImgURL
			}
		}
	}

	var position int64
	for _, guestID := range decodedMainInfo.GuestIDs {
		if err = tx.Create(&model.ProjectGuest{
			ProjectID: pid,
			GuestID:   guestID,
			Position:  position,
			GuestImg:  guestImgMap[guestID],
		}).Error; err != nil {
			log.Error("项目嘉宾添加失败:%s", err)
			tx.Rollback()
			return
		}
		position++
	}

	// 将有场次票价id的mainInfo更新到version_ext
	encodedMainInfo, jsonErr := json.Marshal(decodedMainInfo)
	if jsonErr != nil {
		log.Error("JSONEncode MainInfo失败")
		tx.Rollback()
		return pid, jsonErr
	}

	finalMainInfo := string(encodedMainInfo)
	if len(finalMainInfo) > 4000 {
		log.Error("信息量过大")
		tx.Rollback()
		return pid, ecode.TicketMainInfoTooLarge
	}

	// 编辑version_ext
	err = tx.Model(&model.VersionExt{}).Where("ver_id = ? and type = ?", verID, 1).Update("main_info", finalMainInfo).Error
	if err != nil {
		tx.Rollback()
		log.Error("更新versionext失败: %s", err)
	}

	// 提交事务
	tx.Commit()

	return pid, nil
}

// CreateProject 灌入项目表
func (d *Dao) CreateProject(c context.Context, tx *gorm.DB, verID uint64, decodedMainInfo ProjectMainInfo) (project model.Item, err error) {

	// 创建新project
	venueID, _ := strconv.ParseInt(decodedMainInfo.VenueID, 10, 64)
	placeID, _ := strconv.ParseInt(decodedMainInfo.PlaceID, 10, 64)
	compID, _ := strconv.ParseInt(decodedMainInfo.CompID, 10, 64)
	projectType, _ := strconv.ParseInt(decodedMainInfo.Type, 10, 64)
	projectSponsorType, _ := strconv.ParseInt(decodedMainInfo.SponsorType, 10, 64)
	performanceImg, _ := json.Marshal(decodedMainInfo.PerformanceImg)
	projectInfo := model.Item{
		Name:             decodedMainInfo.Name,
		StartTime:        decodedMainInfo.StartTime,
		EndTime:          decodedMainInfo.EndTime,
		VenueID:          venueID,
		PlaceID:          placeID,
		CompID:           compID,
		PerformanceImage: string(performanceImg),
		TicketDesc:       decodedMainInfo.SellingProp,
		Type:             int32(projectType),
		VerID:            verID,
		SponsorType:      int32(projectSponsorType),
		Label:            decodedMainInfo.Label,
		BuyNumLimit:      BuyNumLimit,
		IsSale:           1,  // 默认值可售
		ExpressFee:       -2, //默认值
	}

	if err = tx.Create(&projectInfo).Error; err != nil {
		log.Error("新建项目失败:%s", err)
		tx.Rollback()
		return model.Item{}, err
	}

	return projectInfo, nil
}

// CreateProjectExtInfo 灌入项目详情表
func (d *Dao) CreateProjectExtInfo(c context.Context, tx *gorm.DB, projectID int64, performanceDesc string) (err error) {

	projectExtInfo := model.ItemDetail{
		ProjectID:       projectID,
		PerformanceDesc: performanceDesc,
	}

	if err = tx.Create(&projectExtInfo).Error; err != nil {
		log.Error("新建项目详情失败:%s", err)
		tx.Rollback()
		return err
	}

	return nil
}

// GetDefaultMainInfo 获取初始化数组的mainInfo
func (d *Dao) GetDefaultMainInfo() ProjectMainInfo {
	return ProjectMainInfo{
		Docs:           []ImgInfo{},
		Screens:        []Screen{},
		TicketsSingle:  []TicketSingle{},
		TicketsPass:    []TicketPass{},
		TicketsAllPass: []TicketPass{},
		TagIDs:         []string{},
		GuestIDs:       []int64{},
		GuestImgs:      []GuestImg{},
	}
}
