package model

//MoniTemp 监控邮件模板
type MoniTemp struct {
	From    string   `json:"from"`
	Members []string `json:"members"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}
