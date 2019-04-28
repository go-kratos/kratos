package model

// UserFaceBFS .
type UserFaceBFS struct {
	URL      string  `json:"url,omitempty"`
	FileName string  `json:"file_name,omitempty"`
	Bucket   string  `json:"bucket,omitempty"`
	Sex      float32 `json:"sex,omitempty"`
	Violent  float32 `json:"violent,omitempty"`
	Blood    float32 `json:"blood,omitempty"`
	Politics float32 `json:"politics,omitempty"`
	IsYellow bool    `json:"is_yellow,omitempty"`
	ErrCode  int32   `json:"error_code,omitempty"`
	ErrMsg   string  `json:"error_msg,omitempty"`
}
