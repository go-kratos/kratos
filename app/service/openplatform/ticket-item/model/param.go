package model

// ParamID ID类型请求
type ParamID struct {
	ID int64 `form:"id" validate:"gt=0,required"`
}

// ParamCards 卡片类型请求
type ParamCards struct {
	IDs string `form:"ids" validate:"gt=0,required"`
}

// ParamBill 项目订单信息请求
type ParamBill struct {
	IDs  string `form:"ids" validate:"gt=0,required"`
	Sids string `form:"sids"`
	Tids string `form:"tids"`
}

// GuestParam 嘉宾参数
type GuestParam struct {
	ID          uint32 `form:"id"`
	Name        string `form:"name" validate:"required"`
	GuestImg    string `form:"guestimg" validate:"required"`
	Description string `form:"description"`
	GuestID     int64  `form:"guestid"`
}

// GuestStatusParam 嘉宾状态
type GuestStatusParam struct {
	ID     int64 `form:"id" validate:"required"`
	Status int8  `form:"status"`
}

// GuestSearchParam 嘉宾搜索
type GuestSearchParam struct {
	Keyword string `form:"keyword"`
	Ps      int    `form:"ps"`
	Pn      int    `form:"pn"`
}

// VenueSearchParam 场馆搜索参数
type VenueSearchParam struct {
	PageParam

	ProvinceID int    `form:"province_id"`
	CityID     int    `form:"city_id"`
	ID         int    `form:"id"`
	Name       string `form:"name"`
}

// VersionSearchParam 版本搜索参数
type VersionSearchParam struct {
	Type       int    `form:"type"`
	TargetItem int    `form:"target_item"`
	ItemName   string `form:"item_name"`
	Status     []int  `form:"status"`

	PageParam
}

// PageParam 分页参数
type PageParam struct {
	Pn int `form:"pn" validate:"min=1,gte=1"`
	Ps int `form:"ps" validate:"min=1,max=10000"`
}

// AreaInfoParam areaInfo接口请求
type AreaInfoParam struct {
	ID         int64  `form:"id" validate:"min=0"`          // 待修改区域的ID（为0表示创建）
	AID        string `form:"a_id" validate:"min=1,max=16"` // 区域自定义编号
	Name       string `form:"name" validate:"min=1,max=16"` // 区域名
	Place      int64  `form:"place" validate:"min=1"`       // 所属场地ID
	Coordinate string `form:"coordinate" validate:"min=1"`  // 区域坐标
}

// PlaceInfoParam placeInfo接口请求
type PlaceInfoParam struct {
	ID      int64  `form:"id" validate:"min=0"`           // 待修改场地的ID（为0表示创建）
	Status  int32  `form:"status" validate:"min=0,max=1"` // 状态
	Name    string `form:"name" validate:"max=16"`        // 场地名
	BasePic string `form:"base_pic" validate:"max=128"`   // 场地底图
	Venue   int64  `form:"venue" validate:"min=1"`        // 场馆ID
	DWidth  int32  `form:"d_width" validate:"min=0"`      // mis画框宽度
	DHeight int32  `form:"d_height" validate:"min=0"`     // mis画框高度
}

// SeatInfoParam seatInfo接口请求
type SeatInfoParam struct {
	Area      int64  `form:"area" validate:"min=1"`       // 区域ID
	SeatsNum  int32  `form:"seats_num" validate:"min=0"`  // 座位数
	SeatMap   string `form:"seat_map" validate:"min=1"`   // 座位图
	Seats     string `form:"seats"`                       // 座位数组JSON
	Width     int32  `form:"width" validate:"min=1"`      // 坐区宽度
	Height    int32  `form:"height" validate:"min=1"`     // 坐区高度
	RowList   string `form:"row_list" validate:"min=1"`   // 行号序列
	SeatStart string `form:"seat_start" validate:"min=1"` // 座位起始坐标
}

// SeatStockParam seatStock接口请求
type SeatStockParam struct {
	Screen   int64  `form:"screen" validate:"min=1"` // 场次ID
	Area     int64  `form:"area" validate:"min=1"`   // 区域ID
	SeatInfo string `form:"seat_info"`               // 座位票种定义数组JSON
}

// RemoveSeatOrdersParam removeSeatOrders接口请求
type RemoveSeatOrdersParam struct {
	Price int64 `form:"price" validate:"min=1"` // 票价ID
}

// VenueInfoParam venueInfo接口请求
type VenueInfoParam struct {
	ID            int64  `form:"id" validate:"min=0"`            // 待修改场馆的ID（为0表示创建）
	Name          string `form:"name" validate:"max=25"`         // 场馆名
	Status        int32  `form:"status" validate:"min=0,max=1"`  // 状态 1-启用 0-停用
	Province      int64  `form:"provid" validate:"min=0"`        // 省份ID
	City          int64  `form:"cityid" validate:"min=0"`        // 城市ID
	District      int64  `form:"distid" validate:"min=0"`        // 区县ID
	AddressDetail string `form:"addr" validate:"min=0,max=60"`   // 详细地址
	Coordinate    string `form:"coordinate" validate:"min=0"`    // 场馆地图坐标及类型字段
	Traffic       string `form:"traff" validate:"min=0,max=100"` // 交通信息
}
