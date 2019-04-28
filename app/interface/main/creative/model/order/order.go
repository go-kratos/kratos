package order

// Order str
type Order struct {
	ExeOdID    int64  `json:"execute_order_id"`
	BzOdName   string `json:"business_order_name"`
	IDCode     int64  `json:"id_code"`
	GameBaseID int64  `json:"game_base_id"`
	GameName   string `json:"game_name"`
}

// Oasis 绿洲计划
type Oasis struct {
	State        int   `json:"state"`
	RealeseOrder int64 `json:"running_execute_order_count"` //投放商单数
	TotalOrder   int64 `json:"total_execute_order_count"`   //总商单数
}

// Growth for order.
type Growth struct {
	State           int     `json:"state"`
	YesterdayIncome float64 `json:"yesterday_income"`     //昨日收入
	MonthIncome     float64 `json:"present_month_income"` //本月收入
}

// OasisEarnings for order.
type OasisEarnings struct {
	State   int   `json:"state"`
	Realese int64 `json:"release"` //投放商单数
	Total   int64 `json:"total"`   //总商单数
}

// GrowthEarnings for order.
type GrowthEarnings struct {
	State     int     `json:"state"`
	Yesterday float64 `json:"yesterday"` //昨日收入
	Month     float64 `json:"month"`     //本月收入
}

// UpValidate for up validate.
type UpValidate struct {
	MID   int64 `json:"mid"`
	State int   `json:"state"` //0:禁止，1:允许
}
