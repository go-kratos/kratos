package model

//Project struct
type Project struct {
	ID           int32           `json:"id"`
	Name         string          `json:"name"`
	StartTime    int64           `json:"start_time"`
	Type         int             `json:"type"`
	PV           map[int32]int64 `json:"pv"`
	UV           map[int32]int64 `json:"uv"`
	SaleInfo     map[int32]int64 `json:"sale_info"`
	WishInfo     map[int32]int64 `json:"wish_info"`
	CommentInfo  map[int32]int64 `json:"comment_info"`
	FavoriteInfo map[int32]int64 `json:"favorite_info"`
}

//OrderSkus struct
type OrderSkus struct {
	ProjectID   int    `json:"project_id"`
	ProjectName string `json:"project_name"`
	ScreenID    int    `json:"screen_id"`
	ScreenName  string `json:"screen_name"`
	SkuID       int    `json:"sku_id"`
	Count       int    `json:"count"`
	OriginPrice int    `json:"origin_price"`
	Price       int    `json:"price"`
	TicketType  string `json:"ticket_type"`
	Desc        string `json:"desc"`
}

//Order struct
type Order struct {
	OrderID      int         `json:"order_id"`
	Ctime        string      `json:"ctime"`
	Mtime        string      `json:"mtime"`
	OrderType    int         `json:"order_type"`
	Source       string      `json:"source"`
	PayMoney     int         `json:"pay_money"`
	TotalMoney   int         `json:"total_money"`
	ExpressFee   int         `json:"express_fee"`
	UID          string      `json:"uid"`
	PersonalID   string      `json:"personal_id"`
	Tel          string      `json:"tel"`
	Status       int         `json:"status"`
	SubStatus    int         `json:"sub_status"`
	RefundStatus int         `json:"refund_status"`
	PayTime      string      `json:"pay_time"`
	Skus         []OrderSkus `json:"skus"`
	ProjectID    int         `json:"project_id"`
	DistUserID   int         `json:"dist_user_id"`
	DistType     int         `json:"dist_type"`
	SendStatus   int         `json:"send_status"`
}

//Comment struct
type Comment struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    *struct {
		Page *struct {
			Count int `json:"count"`
			Num   int `json:"num"`
			Size  int `json:"size"`
		}
		Replies []*struct {
			Ctime  int64 `json:"ctime"`
			Member *struct {
				Mid   string `json:"mid"`
				Uname string `json:"uname"`
			} `json:"member"`
		} `json:"replies"`
	} `json:"data"`
}

// PUVResult get  responce
type PUVResult struct {
	PUV
	DaysBefore int32 `json:"days_before"`
}

// PUV puv info
type PUV struct {
	PV int64 `json:"pv"`
	UV int64 `json:"uv"`
}
