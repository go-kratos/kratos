package model

// Mail def.
type Mail struct {
	ToAddress []*MailAddress
	Subject   string
	Body      string
}

// MailAddress def.
type MailAddress struct {
	Address string
	Name    string
}

// MailData def.
type MailData struct {
	UserName       string
	SourceBranch   string
	TargetBranch   string
	Title          string
	Description    string
	URL            string
	Info           string
	PipelineStatus string
	PipeStatus     string
}
