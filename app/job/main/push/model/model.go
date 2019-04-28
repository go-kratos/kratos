package model

const (
	// CheckJobStatusOk 已完成
	CheckJobStatusOk = 1
	// CheckJobStatusErr 失败
	CheckJobStatusErr = 2
	// CheckJobStatusDoing 进行中
	CheckJobStatusDoing = 3
	// CheckJobStatusPending 等待执行
	CheckJobStatusPending = 4
)

// DpCheckJobResult .
type DpCheckJobResult struct {
	Code      int      `json:"code"`
	Msg       string   `json:"msg"`
	StatusID  int      `json:"statusId"`
	StatusMsg string   `json:"statusMsg"`
	Files     []string `json:"hdfsPath"`
}
