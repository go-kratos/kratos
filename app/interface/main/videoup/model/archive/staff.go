package archive

//联合投稿的分区配置
type StaffTypeConf struct {
	TypeID   int16 `json:"typeid"`
	MaxStaff int   `json:"max_staff"`
}
