package exp

// Exp 经验DB层
type Exp struct {
	UID   int64
	Uexp  int64
	Rexp  int64
	CTime string
	MTime string
}

// ModelExpList exp list
type ModelExpList struct {
	Data []Exp
}

// ModelExpLog 行为日志上报结构
type ModelExpLog struct {
	Mid  int64
	Uexp int64
	Rexp int64
	Ts   int64
	// 业务来源
	ReqBizDesc string
	Buvid      string
	// 具体描述
	Content map[string]string
}

// LevelInfo 等级结构
type LevelInfo struct {
	UID         int64
	UserLevel   UserLevelInfo
	AnchorLevel AnchorLevelInfo
	CTime       string
	MTime       string
}

// UserLevelInfo 等级基础结构
type UserLevelInfo struct {
	Level            int64 `json:"level"`
	NextLevel        int64 `json:"nextLevel"`
	UserExpLeft      int64 `json:"userExpLeft"`
	UserExpRight     int64 `json:"userExpRight"`
	UserExp          int64 `json:"userExp"`
	UserExpNextLevel int64 `json:"userExpNextLevel"`
	Color            int64 `json:"color"`
	UserExpNextLeft  int64 `json:"userExpNextLeft"`
	UserExpNextRight int64 `json:"userExpNextRight"`
	IsLevelTop       int64 `json:"isLevelTop"`
}

// AnchorLevelInfo 等级基础结构
type AnchorLevelInfo struct {
	Level            int64 `json:"level"`
	NextLevel        int64 `json:"nextLevel"`
	UserExpLeft      int64 `json:"userExpLeft"`
	UserExpRight     int64 `json:"userExpRight"`
	UserExp          int64 `json:"userExp"`
	UserExpNextLevel int64 `json:"userExpNextLevel"`
	Color            int64 `json:"color"`
	UserExpNextLeft  int64 `json:"userExpNextLeft"`
	UserExpNextRight int64 `json:"userExpNextRight"`
	AnchorScore      int64 `json:"anchorScore"`
	IsLevelTop       int64 `json:"isLevelTop"`
}
