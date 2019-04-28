package reply

// eq
const (
	EQ = int8(0) // condition equal
	GT = int8(1) // greater
	LT = int8(2) // less
)

// Notice Notice
type Notice struct {
	ID         int64  `json:"id"`
	Plat       int8   `json:"-"`
	Condition  int8   `json:"-"`
	Build      int64  `json:"-"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Link       string `json:"link"`
	ClientType string `json:"-"`
}

// CheckBuild CheckBuild
func (n *Notice) CheckBuild(plat int8, build int64, clientType string) bool {
	if plat == PlatWeb {
		return true
	}
	//如果公告的clientType不为空说明是针对特定的客户端类型，需要特别检查一下
	if n.ClientType != "" && n.ClientType != clientType {
		return false
	}
	// PC OR GT OR LT OR EQUAL
	if (build >= n.Build && n.Condition == GT) || (build <= n.Build && n.Condition == LT) || (build == n.Build && n.Condition == EQ) {
		return true
	}
	return false
}
