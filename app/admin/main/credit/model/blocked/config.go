package blocked

// const config
const (
	ConfigCaseGiveHours      = "case_give_hours"       // 案件发放时长
	ConfigCaseCheckHours     = "case_check_hours"      // 单案审核时长
	ConfigJuryVoteRadio      = "jury_vote_radio"       // 投准率下限
	ConfigCaseJudgeRadio     = "case_judge_radio"      // 判决阙值
	ConfigCaseVoteMin        = "case_vote_min"         // 案件投票数下限
	ConfigCaseObtainMax      = "case_obtain_max"       // 每日获取案件数
	ConfigCaseVoteMax        = "case_vote_max"         // 结案投票数
	ConfigJuryApplyMax       = "jury_apply_max"        // 每日发放风纪委上限
	ConfigCaseLoadMax        = "case_load_max"         // 案件发放最大队列数
	ConfigCaseLoadSwitch     = "case_load_switch"      // 案件发放进入队列开关
	ConfigCaseVoteMaxPercent = "case_vote_max_percent" // 结案投票数的百分比
)

// Config blocked_config model
type Config struct {
	ID           int64
	ConfigKey    string
	Name         string
	Content      string
	Description  string
	OperID       int64 `json:"oper_id"`
	OperatorName string
}

// VoteNum .
type VoteNum struct {
	RateS int8 `json:"rate_s"`
	RateA int8 `json:"rate_a"`
	RateB int8 `json:"rate_b"`
	RateC int8 `json:"rate_c"`
	RateD int8 `json:"rate_d"`
}

// TableName case tablename
func (*Config) TableName() string {
	return "blocked_config"
}
