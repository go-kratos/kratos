package email

const (
	//MailKey 实时邮件队列
	MailKey = "f_mail_list"
	//MailFastKey 超限名单的邮件队列
	MailFastKey = "f_mail_list_fast"
	//RetryListKey 重试列表
	RetryListKey = "f_retry_list"
	//RetryActionReply 评论重试
	RetryActionReply = "reply"
)

//Retry retry
type Retry struct {
	AID        int64  `json:"aid"`
	Action     string `json:"action"`
	Flag       int64  `json:"flag"`
	FlagA      int64  `json:"flag_a"`
	CreateTime int64  `json:"create_time"`
}
