package model

// Subtitle state
const (
	AttrSubtitleClose = uint(1) // 关闭稿件字幕
)

// SubtitleSubject .
type SubtitleSubject struct {
	Aid   int64 `json:"aid"`
	Allow bool  `json:"allow"`
	Attr  int32 `json:"attr"`
	Lan   uint8 `json:"lan"`
}

// AttrVal return val of subtitle subject'attr
func (s *SubtitleSubject) AttrVal(bit uint) int32 {
	return (s.Attr >> bit) & int32(1)
}

// AttrSet set val of subtitle subject'attr
func (s *SubtitleSubject) AttrSet(v int32, bit uint) {
	s.Attr = s.Attr&(^(1 << bit)) | (v << bit)
}
