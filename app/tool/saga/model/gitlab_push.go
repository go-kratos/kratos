package model

// HookPush def
type HookPush struct {
	ObjectKind        string      `json:"object_kind"`
	Before            string      `json:"before"`
	After             string      `json:"after"`
	Ref               string      `json:"ref"`
	CheckoutSHA       string      `json:"checkout_sha"`
	UserID            int64       `json:"user_id"`
	UserName          string      `json:"user_name"`
	UserUserName      string      `json:"user_username"`
	UserEmail         string      `json:"user_email"`
	UserAvatar        string      `json:"user_avatar"`
	ProjectID         int64       `json:"project_id"`
	Project           *Project    `json:"project"`
	Repository        *Repository `json:"repository"`
	Commits           []*Commit   `json:"commits"`
	TotalCommitsCount int64       `json:"total_commits_count"`
}
