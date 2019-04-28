package model

// AssociatePanelInfo associate panel info.
type AssociatePanelInfo struct {
	ID         int64   `json:"id"`
	Month      int32   `json:"month"`
	PdName     string  `json:"product_name"`
	PdID       string  `json:"product_id"`
	SubType    int32   `json:"sub_type"`
	SuitType   int32   `json:"suit_type"`
	OPrice     float64 `json:"original_price"`
	DPrice     float64 `json:"discount_price"`
	DRate      string  `json:"discount_rate"`
	Remark     string  `json:"remark"`
	Selected   int32   `json:"selected"`
	PayState   int8    `json:"pay_state"`
	PayMessage string  `json:"pay_message"`
}

// ArgAssociatePanel args.
type ArgAssociatePanel struct {
	Device    string `form:"device"`
	Build     int64  `form:"build"`
	MobiApp   string `form:"mobi_app"`
	Platform  string `form:"platform" default:"pc"`
	SortTP    int8   `form:"sort_type"`
	PanelType string `form:"panel_type" default:"normal"`
	Mid       int64
	IP        string
}
