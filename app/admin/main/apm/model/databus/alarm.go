package databus

// AlarmOpen ...
type AlarmOpen struct {
	ReqID   string  `json:"ReqId"`
	Action  string  `json:"Action"`
	RetCode int64   `json:"RetCode"`
	Data    []*Open `json:"Data"`
}

// Open ...
type Open struct {
	AdjustID string `json:"adjust_id"`
	PolicyID string `json:"policy_id"`
	Group    string `json:"group"`
}

//Scope ...
type Scope struct {
	Type int64
	Key  string
	Val  []string
}

//Owner ...
type Owner struct {
	Owner string `json:"owner"`
	App   string `json:"app"`
}

// Res ...
type Res struct {
	ReqID   string  `json:"ReqId"`
	Action  string  `json:"Action"`
	RetCode int64   `json:"RetCode"`
	Data    Opsmind `json:"Data"`
}

// Opsmind ...
type Opsmind struct {
	PolicyID         string `json:"policy_id"`
	AdjustID         string `json:"adjust_id"`
	Category         string `json:"category"`
	Scope            string `json:"scope"`
	TriggerID        string `json:"trigger_id"`
	TriggerLevel     string `json:"trigger_level"`
	TriggerFor       int64  `json:"trigger_for"`
	TriggerNotes     string `json:"trigger_notes"`
	TriggerOperator  string `json:"trigger_operator"`
	TriggerThreshold int64  `json:"trigger_threshold"`
	Silence          bool   `json:"silence"`
	Hashid           string `json:"hashid"`
	ExpiredAt        string `json:"expired_at"`
}

//Query ...
type Query struct {
	Key string
	Val []string
}

// ResQuery ...
type ResQuery struct {
	ReqID   string    `json:"ReqId"`
	Action  string    `json:"Action"`
	RetCode int64     `json:"RetCode"`
	Data    []*Querys `json:"Data"`
}

// Querys ...
type Querys struct {
	ID        string     `json:"id"`
	PolicyID  string     `json:"policy_id"`
	Creator   string     `json:"creator"`
	Ctime     int64      `json:"ctime"`
	Mtime     int64      `json:"mtime"`
	Scope     []*Scope   `json:"scope"`
	Triggers  []*Trigger `json:"triggers"`
	Notes     *Owner     `json:"notes"`
	Desc      string     `json:"desc"`
	Silence   bool       `json:"silence"`
	ExpiredAt int64      `json:"expired_at"`
}

//Trigger ...
type Trigger struct {
	ID         string  `json:"id"`
	Desc       string  `json:"desc"`
	Operator   string  `json:"operator"`
	For        int64   `json:"for"`
	Threshold  float64 `json:"threshold"`
	Level      string  `json:"level"`
	NodataType string  `json:"nodata_type"`
	NodataFor  int64   `json:"nodata_for"`
	Notes      *Owner  `json:"notes"`
}
