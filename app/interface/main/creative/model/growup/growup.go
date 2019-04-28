package growup

import "go-common/library/time"

//UpInfo get up state info
type UpInfo struct {
	MID              int64     `json:"mid"`
	Fans             int64     `json:"fans"`                   //粉丝数量
	NickName         string    `json:"nickname"`               //用户昵称
	OriginalArcCount int       `json:"original_archive_count"` //UP主原创投稿数
	MainCategory     int       `json:"main_category"`          //UP主主要投稿区ID
	AccountState     int       `json:"account_state"`          //账号状态; 1: 未申请; 2: 待审核; 3: 已签约; 4.已驳回; 5.主动退出; 6:被动退出; 7:封禁
	SignType         int8      `json:"sign_type"`              //签约类型; 0: 基础, 1: 首发
	QuitType         int       `json:"quit_type"`              //退出类型: 0: 主动退出 1: 封禁; 2: 平台清退
	ApplyAt          time.Time `json:"apply_at"`               //申请时间
	CTime            time.Time `json:"ctime"`
	MTime            time.Time `json:"mtime"`
}

//UpStatus get up status info
type UpStatus struct {
	Blocked      bool   `json:"blocked"`
	AccountType  int    `json:"account_type"`  //账号类型 1-UGC 2- PGC
	AccountState int    `json:"account_state"` //账号状态; 1: 未申请; 2: 待审核; 3: 已签约; 4.已驳回; 5.主动退出; 6:被动退出; 7:封禁
	ExpiredIn    int64  `json:"expired_in"`    //冷却过期天数
	Reason       string `json:"reason"`        //封禁/驳回/清退（被动退出）理由
	InWhiteList  bool   `json:"in_white_list"` //是否在白名单中，blocked字段在第一期中被忽略，第二期会去掉该字段
	ArchiveType  []int  `json:"archive_type"`  //投稿类型，1：视频，2：音频，3：专栏
	ShowPanel    bool   `json:"show_panel"`
	ShowPanelMsg string `json:"show_panel_msg"`
}

//Summary get summary income.
type Summary struct {
	BreachMoney  float64 `json:"breachMoney"`  //违反金额
	Income       float64 `json:"income"`       //当月收入
	TotalIncome  float64 `json:"totalIncome"`  //累计收入
	WaitWithdraw float64 `json:"waitWithdraw"` //带提现
	Date         string  `json:"date"`
	DayIncome    float64 `json:"dayIncome"`
}

//Stat get statistic income.
type Stat struct {
	ProportionDraw map[string]float64 `json:"proportionDraw"` //比例图
	LineDraw       []*LineDraw        `json:"lineDraw"`
	Tops           []*TopArc          `json:"tops"`
	Desc           string             `json:"desc"`
}

//LineDraw for income data.
type LineDraw struct {
	DateKey int64   `json:"dateKey"`
	Income  float64 `json:"income"`
}

//TopArc get top archive.
type TopArc struct {
	AID         int64   `json:"aid"`
	Title       string  `json:"title"`
	TypeName    string  `json:"typeName"`    //type类型
	TotalIncome float64 `json:"totalIncome"` //累计收入
}

//IncomeList get income list.
type IncomeList struct {
	Page       int `json:"page"`
	TotalCount int `json:"total_count"`
	Data       []*struct {
		AID         int64   `json:"aid"`
		Title       string  `json:"title"`
		Income      float64 `json:"income"`      //当月收入
		TotalIncome float64 `json:"totalIncome"` //累计收入
	} `json:"data"`
}

//BreachList get reach list.
type BreachList struct {
	Page       int `json:"page"`
	TotalCount int `json:"total_count"`
	Data       []*struct {
		AID        int64   `json:"aid"`
		BreachTime int64   `json:"breachTime"` //时间戳
		Money      float64 `json:"money"`      //扣除金额
		Reason     string  `json:"reason"`     //原因
		Title      string  `json:"title"`      //稿件标题
	} `json:"data"`
}
