package model

const (
	//UpperTypeWhite 优质
	UpperTypeWhite int8 = 1
	//UpperTypeBlack 高危
	UpperTypeBlack int8 = 2
	//UpperTypePGC 生产组
	UpperTypePGC int8 = 3
	//UpperTypeUGCX don't know
	UpperTypeUGCX int8 = 3
	//UpperTypePolitices 时政
	UpperTypePolitices int8 = 5
	//UpperTypeEnterprise 企业
	UpperTypeEnterprise int8 = 7
	//UpperTypeSigned 签约
	UpperTypeSigned int8 = 15
)

//UPGroup up主所属的所有特殊用户组
type UPGroup struct {
	ID        int64  `json:"id"`
	Tag       string `json:"tag"`
	ShortTag  string `json:"short_tag"`
	FontColor string `json:"font_color"` //字体颜色
	BgColor   string `json:"bg_color"`   //背景颜色
	Note      string `json:"note"`
}
