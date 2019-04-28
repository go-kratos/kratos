package model

// DirConfig directory config
type DirConfig struct {
	Pic  DirPicConfig  `json:"dir_pic_config"`
	Rate DirRateConfig `json:"dir_rate_config"`
}

// DirPicConfig directory picture config
type DirPicConfig struct {
	FileSize           int      `json:"file_size"`             //文件大小上限 单位 Byte
	MaxPixelWidthSize  int      `json:"max_pixel_width_size"`  //像素宽上限
	MinPixelWidthSize  int      `json:"min_pixel_width_size"`  //像素高下限
	MaxPixelHeightSize int      `json:"max_pixel_height_size"` //像素高上限
	MinPixelHeightSize int      `json:"min_pixel_height_size"` //像素宽下限
	MaxAspectRatio     float64  `json:"max_aspect_ratio"`      //最大宽高比
	MinAspectRatio     float64  `json:"min_aspect_ratio"`      //最小宽高比
	AllowType          string   `json:"allow_type"`            //允许的MIME类型
	AllowTypeSlice     []string // 允许的MIME类型列表,AllowTypeSlice = strings.Split(AllowType,",")
}

// DirRateConfig directory rate config
type DirRateConfig struct {
	// secondQPS 接受 countQPS 个请求
	SecondQPS int `json:"second_qps"`
	CountQPS  int `json:"count_qps"`
}

//{
//    file_size： 100                   文件大小上限 单位 Byte
//    max_pixel_width_size： 1024       像素宽上限
//    max_pixel_height_size：1024       像素高上限
//    min_pixel_width_size： 10         像素宽下限
//    min_pixel_height_size：10         像素高下限
//    max_aspect_ratio： 100            最大宽高比
//    min_aspect_ratio： 10             最小宽高比
//}

//{
//	max_user_qps                      最大用户qps
//	max_user_upload_number            每日最大用户上传数量
//}
