package model

// ScoreChange DataBus spy score change.
type ScoreChange struct {
	Mid        int64  `json:"mid"`
	Score      int8   `json:"score"`
	BaseScore  int8   `json:"base_score"`
	EventScore int8   `json:"event_score"`
	TS         int64  `json:"ts"`
	Reason     string `json:"reason"`
	RiskLevel  int8   `json:"risk_level"`
}

const (
	//CoinReason coin reason ref update score
	CoinReason = "coin-service"

	//CoinHighRisk coin high risk
	CoinHighRisk = 7
)
