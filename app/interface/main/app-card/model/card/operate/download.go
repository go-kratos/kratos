package operate

import "go-common/app/interface/main/app-card/model"

type Download struct {
	ID          int64  `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Desc        string `json:"desc,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Cover       string `json:"cover,omitempty"`
	URLType     int    `json:"url_type,omitempty"`
	URLValue    string `json:"url_value,omitempty"`
	BtnTxt      int    `json:"btn_txt,omitempty"`
	ReType      int    `json:"re_type,omitempty"`
	ReValue     string `json:"re_value,omitempty"`
	DoubleCover string `json:"double_cover,omitempty"`
	Number      int32  `json:"number,omitempty"`
	// extra
	ButtonText string   `json:"button_text,omitempty"`
	Goto       model.Gt `json:"goto,omitempty"`
	Param      string   `json:"param,omitempty"`
}

func (c *Download) Change() {
	switch c.BtnTxt {
	case 0:
		c.ButtonText = "下载"
	case 1:
		c.ButtonText = "预约"
	case 2:
		c.ButtonText = "查看详情"
	}
	c.Goto = model.OperateType[c.URLType]
	c.Param = c.URLValue
}
