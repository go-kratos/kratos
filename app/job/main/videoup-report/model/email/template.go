package email

//邮件的发件人、收件列表、抄送列表、主题
const (
	FROM    = "From"
	TO      = "To"
	CC      = "Cc"
	SUBJECT = "Subject"
)

//Template 邮件模板
type Template struct {
	Headers     map[string][]string
	Body        string
	ContentType string
	Type        string
	AID         int64
	UID         int64
	Username    string
	Department  string
}
