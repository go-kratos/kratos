package model

import "go-common/library/time"

//ListUpCardInfoArg arg
type ListUpCardInfoArg struct {
	Pn uint `form:"pn"` // page num
	Ps uint `form:"ps"` // query size
}

//UpCardInfoPage page result of card info
type UpCardInfoPage struct {
	Cards map[int64]*UpCard `json:"cards"`
	Page  *Pager            `json:"page"`
}

//GetCardByMidArg arg
type GetCardByMidArg struct {
	Mid int64 `form:"mid" validate:"required"`
}

//ListCardByMidsArg arg
type ListCardByMidsArg struct {
	Mids string `form:"mids" validate:"required"` // mids split by ","
}

//UpCard up card content
type UpCard struct {
	UpCardInfo *UpCardInfo      `json:"up_card_info"`
	Accounts   []*UpCardAccount `json:"accounts"`
	Videos     []*UpCardVideo   `json:"videos"`
	Images     []*UpCardImage   `json:"images"`
}

//UpCardInfo for up info in card info
type UpCardInfo struct {
	MID             int64     `json:"mid"`
	NameCN          string    `json:"name_cn"`
	NameEN          string    `json:"name_en"`
	NameAlias       string    `json:"name_alias"`
	Signature       string    `json:"signature"`
	Content         string    `json:"content"`
	Nationality     string    `json:"nationality"`
	Nation          string    `json:"nation"`
	Gender          string    `json:"gender"`
	BloodType       string    `json:"blood_type"`
	Constellation   string    `json:"constellation"`
	Height          int       `json:"height"`
	Weight          int       `json:"weight"`
	BirthPlace      string    `json:"birth_place"`
	BirthDate       time.Time `json:"birth_date"`
	Occupation      string    `json:"occupation"`
	Tags            string    `json:"tags"`
	Masterpieces    string    `json:"masterpieces"`
	School          string    `json:"school"`
	Location        string    `json:"location"`
	Interests       string    `json:"interests"`
	Platform        string    `json:"platform"`
	PlatformAccount string    `json:"platform_account"`
}

//UpCardAccount for accounts in card info
type UpCardAccount struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Picture string `json:"picture"`
}

//UpCardImage for images in card info
type UpCardImage struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

//UpCardVideo for videos in card info
type UpCardVideo struct {
	URL      string    `json:"url"`
	Title    string    `json:"title"`
	Picture  string    `json:"picture"`
	Duration int64     `json:"duration"`
	CTime    time.Time `json:"ctime"`
}
