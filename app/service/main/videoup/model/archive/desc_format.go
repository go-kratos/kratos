package archive

// DescFormat is desc model.
type DescFormat struct {
	ID         int64  `json:"id"`
	TypeID     int16  `json:"typeid"`
	Copyright  int8   `json:"copyright"`
	Components string `json:"components"`
	Lang       int8   `json:"lang"`
	Platform   int8   `json:"platform"`
}
