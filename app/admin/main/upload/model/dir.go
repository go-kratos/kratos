package model

import xtime "go-common/library/time"

// DirConfig dir config
type DirConfig struct {
	Pic  DirPicConfig  `json:"dir_pic_config"`
	Rate DirRateConfig `json:"dir_rate_config"`
}

// DirPicConfig pic config
type DirPicConfig struct {
	FileSize           uint    `json:"file_size"`             //文件大小上限 单位 Byte
	MaxPixelWidthSize  uint    `json:"max_pixel_width_size"`  //像素宽上限
	MinPixelWidthSize  uint    `json:"min_pixel_width_size"`  //像素高下限
	MaxPixelHeightSize uint    `json:"max_pixel_height_size"` //像素高上限
	MinPixelHeightSize uint    `json:"min_pixel_height_size"` //像素宽下限
	MaxAspectRatio     float64 `json:"max_aspect_ratio"`      //最大宽高比
	MinAspectRatio     float64 `json:"min_aspect_ratio"`      //最小宽高比
	AllowType          string  `json:"allow_type"`            //允许的MIME类型
}

// DirRateConfig rate config
type DirRateConfig struct {
	// SecondQPS 接受 CountQPS 个请求
	SecondQPS uint `json:"second_qps"`
	CountQPS  uint `json:"count_qps"`
}

// DirLimit table dir_limit ORM
type DirLimit struct {
	ID         int        `json:"id" gorm:"column:id"`
	BucketName string     `json:"bucket_name" gorm:"column:bucket_name"`
	Dir        string     `json:"dir" gorm:"column:dir"`
	ConfigPic  string     `json:"config_pic" gorm:"column:config_pic"`
	ConfigRate string     `json:"config_rate" gorm:"column:config_rate"`
	CTime      xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime      xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName dir_limit
func (dl DirLimit) TableName() string {
	return "dir_limit"
}
