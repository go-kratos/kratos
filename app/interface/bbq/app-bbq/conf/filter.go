package conf

import "github.com/BurntSushi/toml"

// UploadFilter 创作中心上传过滤
type UploadFilter struct {
	MidFilter *UploadMidFilter
}

// UploadMidFilter 创作中心上传用户mid过滤
type UploadMidFilter struct {
	White []int64
	Black []int64
}

// Set .
func (uf *UploadFilter) Set(text string) error {
	if _, err := toml.Decode(text, uf); err != nil {
		panic(err)
	}
	return nil
}

// Set .
func (umf *UploadMidFilter) Set(text string) error {
	if _, err := toml.Decode(text, umf); err != nil {
		panic(err)
	}
	return nil
}
