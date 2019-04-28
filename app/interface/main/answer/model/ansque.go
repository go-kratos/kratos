package model

// AnsQueDetail .
type AnsQueDetail struct {
	ID            int64   `json:"qs_id"`
	AnsImg        string  `json:"ans_img"`
	QsHeight      float64 `json:"qs_h"`
	QsPositionY   float64 `json:"qs_y"`
	Ans1Hash      string  `json:"ans1_hash"`
	Ans2Hash      string  `json:"ans2_hash"`
	Ans3Hash      string  `json:"ans3_hash"`
	Ans4Hash      string  `json:"ans4_hash"`
	Ans0Height    float64 `json:"ans0_h"`
	Ans0PositionY float64 `json:"ans0_y"`
	Ans1Height    float64 `json:"ans1_h"`
	Ans1PositionY float64 `json:"ans1_y"`
	Ans2Height    float64 `json:"ans2_h"`
	Ans2PositionY float64 `json:"ans2_y"`
	Ans3Height    float64 `json:"ans3_h"`
	Ans3PositionY float64 `json:"ans3_y"`
}

// AnsQueDetailList .
type AnsQueDetailList struct {
	CurrentTime int64           `json:"current_time"`
	EndTime     int64           `json:"end_time"`
	QuesList    []*AnsQueDetail `json:"items"`
}

// AnsProType .
type AnsProType struct {
	List        []*AnsTypeList `json:"list"`
	CurrentTime int64          `json:"current_time"`
	EndTime     int64          `json:"end_time"`
	Repro       string         `json:"repro"`
}

// AnsTypeList .
type AnsTypeList struct {
	Name   string     `json:"name"`
	Fields []*AnsType `json:"fields"`
}

// AnsType info.
type AnsType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
