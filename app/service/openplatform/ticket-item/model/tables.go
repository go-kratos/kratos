package model

import (
	"go-common/library/time"
)

// Venue 场馆表
type Venue struct {
	ID            int64
	Name          string
	Status        int32
	Province      int64
	City          int64
	District      int64
	AddressDetail string
	Traffic       string
	Coordinate    string
	PlaceNum      int32
	Ctime         time.Time
	Mtime         time.Time
}

// Coor 项目图片结构
type Coor struct {
	Type string
	Coor string
}

// ItemDetail 项目详情表
type ItemDetail struct {
	ProjectID       int64  `json:"id"`
	PerformanceDesc string `json:"detail"`
}

// Item 项目表
type Item struct {
	ID               int64
	Name             string
	Status           int32
	IsSale           int32
	StartTime        int32
	EndTime          int32
	VenueID          int64
	PlaceID          int64
	CompID           int64
	ExpressFee       int32
	HasExpressFee    int32
	ExpressFreeFlag  int32
	PerformanceImage string
	TicketDesc       string
	BuyNumLimit      string
	Recommend        int32
	PromoTags        string
	VerID            uint64
	BuyerInfo        string
	Type             int32
	SponsorType      int32
	Label            string
	Img              *ItemImg
}

// ItemImg 项目图片结构
type ItemImg struct {
	First struct {
		URL  string
		Desc string
	}
	Banner struct {
		URL  string
		Desc string
	}
}

// Screen 场次表.
type Screen struct {
	ID           int64
	Name         string
	Status       int32
	Type         int32
	TicketType   int32
	ScreenType   int32
	DeliveryType int32
	PickSeat     int32
	StartTime    int32
	EndTime      int32
	ProjectID    int64
	SaleStart    int64
	SaleEnd      int64
}

// TicketPrice 票价表
type TicketPrice struct {
	ID            int64
	ParentID      int64
	Desc          string
	Type          int32
	SaleType      int32
	LinkSc        string
	LinkTicketID  int64
	Symbol        string
	Color         string
	BuyLimit      int32
	DescDetail    string
	ScreenID      int64
	IsSale        int32
	IsVisible     int32
	IsRefund      int32
	Price         int32
	OriginPrice   int32
	MarketPrice   int32
	ProjectID     int64
	PaymentMethod int32
	PaymentValue  int64
	SaleTime      string
	SaleStart     time.Time
	SaleEnd       time.Time
}

// TicketPriceExtra 票价额外表
type TicketPriceExtra struct {
	ID        int64
	ProjectID int64
	SkuID     int64
	Attrib    string
	Value     string
	IsDeleted int32
}

// Guest Build guest
type Guest struct {
	ID          int64  `json:"id"`
	GuestImg    string `json:"guest_img"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      int8   `json:"status"`
	GuestID     int64  `json:"guest_id"`
}

// ProjectGuest Build project_guest
type ProjectGuest struct {
	ID           int64  `json:"id"`
	ProjectID    int64  `json:"project_id"`
	GuestID      int64  `json:"guest_id"`
	Position     int64  `json:"position"`
	GuestImg     string `json:"guest_img"`
	DeleteStatus int32  `json:"delete_status"`
}

// Bulletin Build bulletin
type Bulletin struct {
	ID         int64     `json:"id"`
	Status     int8      `json:"status"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	ProjectID  int64     `json:"project_id"`
	VerID      uint64    `json:"ver_id"`
	BulletinID int64     `json:"bulletin_id"`
	Ctime      time.Time `json:"ctime"`
	Mtime      time.Time `json:"mtime"`
}

// BulletinExtra Build BulletinExtra
type BulletinExtra struct {
	ID         int64  `json:"id"`
	Detail     string `json:"detail"`
	BulletinID int64  `json:"bulletin_id"`
}

// Version Build Version
type Version struct {
	ID         int64     `json:"id"`
	VerID      uint64    `json:"ver_id"`
	Type       int32     `json:"type"`
	Status     int32     `json:"status"`
	ItemName   string    `json:"item_name"`
	Ver        string    `json:"ver"`
	TargetItem int64     `json:"target_item"`
	AutoPub    int32     `json:"auto_pub"`
	ParentID   int64     `json:"parent_id"`
	PubStart   time.Time `json:"pub_start"`
	PubEnd     time.Time `json:"pub_end"`
	For        int64     `json:"for"`
}

// VersionExt Build
type VersionExt struct {
	ID       int64     `json:"id"`
	VerID    uint64    `json:"ver_id"`
	Type     int32     `json:"type"`
	MainInfo string    `json:"main_info"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}

// VersionLog Build VersionLog
type VersionLog struct {
	ID     int64     `json:"id"`
	VerID  uint64    `json:"ver_id"`
	Type   int32     `json:"type"`
	Log    string    `json:"item_name"`
	IsPass int32     `json:"is_pass"`
	Uname  string    `json:"uname"`
	Ctime  time.Time `json:"ctime"`
	Mtime  time.Time `json:"mtime"`
}

// UserWish 想去表
type UserWish struct {
	ID     int64     `json:"id"`
	MID    int64     `json:"mid"`
	ItemID int64     `json:"item_id"`
	Face   string    `json:"face"`
	Ctime  time.Time `json:"ctime"`
	Mtime  time.Time `json:"mtime"`
}

// Place 场地表
type Place struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	BasePic string    `json:"base_pic"`
	Status  int32     `json:"status"`
	Venue   int64     `json:"venue"`
	DWidth  int32     `json:"d_width"`
	DHeight int32     `json:"d_height"`
	Ctime   time.Time `json:"ctime"`
	Mtime   time.Time `json:"mtime"`
}

// PlacePolygon 场地坐标表
type PlacePolygon struct {
	ID         int64  `json:"id"`
	Coordinate string `json:"coordinate"`
}

// Area 区域表
type Area struct {
	ID            int64     `json:"id"`
	AID           string    `json:"name"`
	Name          string    `json:"base_pic"`
	SeatsNum      int32     `json:"seats_num"`
	Width         int32     `json:"width"`
	Height        int32     `json:"height"`
	Place         int64     `json:"venue"`
	DeletedStatus int32     `json:"deleted_status"`
	ColStart      int32     `json:"col_start"`
	ColType       int32     `json:"col_type"`
	ColDirection  int32     `json:"col_direction"`
	RowList       string    `json:"row_list"`
	SeatStart     string    `json:"seat_start"`
	Ctime         time.Time `json:"ctime"`
	Mtime         time.Time `json:"mtime"`
}

// AreaSeats 区域座位表
type AreaSeats struct {
	ID      int64  `json:"id"`
	X       int32  `json:"x"`
	Y       int32  `json:"y"`
	Label   string `json:"label"`
	Bgcolor string `json:"bgcolor"`
	Area    int64  `json:"area"`
	Dstatus int32  `json:"dstatus"`
}

// AreaSeatmap 区域座位表
type AreaSeatmap struct {
	ID      int64  `json:"id"`
	SeatMap string `json:"seat_map"`
}

//Promotion 拼团表
type Promotion struct {
	ID        int64 `json:"id"`
	ItemID    int64 `json:"item_id"`
	SkuID     int64 `json:"sku_id"`
	Extra     int64 `json:"extra"`
	BeginTime int32 `json:"begin_time"`
	EndTime   int32 `json:"end_time"`
	Status    int32 `json:"status"`
}

//Stock 库存表
type Stock struct {
	SkuID       int64  `json:"sku_id"`
	ParentSkuID int64  `json:"parent_sku_id"`
	ItemID      int64  `json:"item_id"`
	Specs       string `json:"specs"`
	TotalStock  int64  `json:"total_stock"`
	Stock       int64  `json:"stock"`
}

//ProjectTags 项目标签表
type ProjectTags struct {
	ID        int64  `json:"id"`
	Status    int32  `json:"status"`
	ProjectID int64  `json:"project_id"`
	TagID     int64  `json:"tag_id"`
	TagName   string `json:"tag_name"`
}

// SeatOrder 座位
type SeatOrder struct {
	ID int64 `json:"id"`
	// 区域ID
	AreaID int64 `json:"area_id"`
	// 场次ID
	ScreenID int64 `json:"screen_id"`
	// 订单号
	OrderID int64 `json:"order_id"`
	// 行号
	Row int32 `json:"row"`
	// 列号
	Col int32 `json:"col"`
	// 状态 0-可售 1-已退票 2-已出票 3-已锁定 4-已预订
	Status int32 `json:"status"`
	// 价格ID
	PriceID int64 `json:"price_id"`
	// 价格
	Price int32 `json:"price"`
	// 操作ID，book.id或lock_rec.id
	OpID int64 `json:"op_id"`
	// 删除时间
	DeletedAt time.Time `json:"deleted_at"`
}

// SeatSet 单场座位及价格配置表
type SeatSet struct {
	ID int64 `json:"id"`
	// 区域ID
	AreaID int64 `json:"area_id"`
	// 场次ID
	ScreenID int64 `json:"screen_id"`
	// 票价设置图
	SeatChart string `json:"seat_chart"`
}

// Banner banner广告投放信息表
type Banner struct {
	ID         int64     `json:"id"`
	PubStart   time.Time `json:"pub_start"`
	PubEnd     time.Time `json:"pub_end"`
	Status     int32     `json:"status"`
	Name       string    `json:"name"`
	Pic        string    `json:"pic"`
	URL        string    `json:"url"`
	From       string    `json:"from"`
	TargetID   int64     `json:"target_id"`
	TargetUser int32     `json:"target_user"`
}

// BannerDistrict banner区域关系表
type BannerDistrict struct {
	ID          int64 `json:"id"`
	BannerID    int64 `json:"banner_id"`
	DistrictID  int64 `json:"district_id"`
	Position    int32 `json:"position"`
	SubPosition int32 `json:"sub_position"`
	Order       int32 `json:"order"`
}

// TableName project.
func (Item) TableName() string {
	return "project"
}

// TableName project_extra.
func (ItemDetail) TableName() string {
	return "project_extra"
}

// TableName project_guests
func (ProjectGuest) TableName() string {
	return "project_guests"
}

// TableName project_bulletin
func (Bulletin) TableName() string {
	return "project_bulletin"
}

// TableName bulletin_extra
func (BulletinExtra) TableName() string {
	return "bulletin_extra"
}

// TableName venue.
func (Venue) TableName() string {
	return "venue"
}

// TableName guest
func (Guest) TableName() string {
	return "guest"
}

// TableName version
func (Version) TableName() string {
	return "version"
}

// TableName version_ext
func (VersionExt) TableName() string {
	return "version_ext"
}

// TableName version_log
func (VersionLog) TableName() string {
	return "version_log"
}

// TableName screen
func (Screen) TableName() string {
	return "screen"
}

// TableName ticket_price
func (TicketPrice) TableName() string {
	return "ticket_price"
}

// TableName ticket_price_extra
func (TicketPriceExtra) TableName() string {
	return "ticket_price_extra"
}

// TableName place
func (Place) TableName() string {
	return "place"

}

// TableName place_polygon
func (PlacePolygon) TableName() string {
	return "place_polygon"
}

// TableName user_wish
func (UserWish) TableName() string {
	return "user_wish"
}

// TableName area
func (Area) TableName() string {
	return "area"
}

// TableName area_seats
func (AreaSeats) TableName() string {
	return "area_seats"
}

// TableName area_seatmap
func (AreaSeatmap) TableName() string {
	return "area_seatmap"
}

// TableName promotion
func (Promotion) TableName() string {
	return "promotion"
}

// TableName sku_stock
func (Stock) TableName() string {
	return "sku_stock"
}

// TableName project_tag
func (ProjectTags) TableName() string {
	return "project_tags"
}

// TableName seat_order
func (SeatOrder) TableName() string {
	return "seat_order"
}

// TableName seat_set
func (SeatSet) TableName() string {
	return "seat_set"
}

// TableName banner
func (Banner) TableName() string {
	return "banner"
}

// TableName banner_district
func (BannerDistrict) TableName() string {
	return "banner_district"
}
