package model

//RecPoolConf  投放配置
type RecPoolConf struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Type        int64   `json:"type"`
	Rule        string  `json:"rules"`
	Priority    int64   `json:"priority"`
	Percent     float64 `json:"percent"`
	TruePercent float64 `json:"true_percent"`
	ModuleType  int64   `json:"module_type"`
	Position    int64   `json:"position"`
}

//RecRoomInfo 房间信息
type RecRoomInfo struct {
	Uid             int64  `json:"uid"`
	Title           string `json:"title"`
	PopularityCount int64  `json:"popularity_count"`
	KeyFrame        string `josn:"Keyframe"`
	Cover           string `josn:"cover"`
	ParentAreaID    int64  `json:"parent_area_id"`
	ParentAreaName  string `json:"parent_area_name"`
	AreaID          int64  `json:"area_id"`
	AreaName        string `josn:"area_name"`
}

//NewRecPoolConf 创建
func NewRecPoolConf() *RecPoolConf {
	return &RecPoolConf{}
}

//NewRecRoomInfo 创建
func NewRecRoomInfo() *RecRoomInfo {
	return &RecRoomInfo{}
}

//RecPoolSlice 配置
type RecPoolSlice []*RecPoolConf

//Len 返回长度
func (R RecPoolSlice) Len() int {
	return len(R)
}

//Less 根据优先级降序排序
func (R RecPoolSlice) Less(i, j int) bool {
	return R[i].Priority > R[j].Priority
}

//Swap 交换数据
func (R RecPoolSlice) Swap(i, j int) {
	R[i], R[j] = R[j], R[i]
}
