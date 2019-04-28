package common

const (
	//NotDeleted db not deleted
	NotDeleted = 0
	//Deleted db deleted
	Deleted = 1
	//Verify 待审核
	Verify = 1
	//Pass 已通过
	Pass = 2
	//Rejecte 已拒绝
	Rejecte = 3
	//Valid 已生效
	Valid = 4
	//InValid 已失效
	InValid = 5
	//StatusOnline status online
	StatusOnline = 1
	//StatusDownline status downline
	StatusDownline = 0
	//OptionOnline option online
	OptionOnline = "online"
	//OptionHidden option downline
	OptionHidden = "hidden"
	//OptionPass option pass
	OptionPass = "pass"
	//OptionReject option reject
	OptionReject = "reject"
)

//Page pager
type Page struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

//CardPreview card preview
type CardPreview struct {
	Title string `json:"title"`
}
