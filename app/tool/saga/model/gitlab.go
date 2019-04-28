package model

// User def
type User struct {
	Name      string `json:"name"`
	UserName  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

// Project def
type Project struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	WebURL            string `json:"web_url"`
	AvatarURL         string `json:"avatar_url"`
	GitSSHURL         string `json:"git_ssh_url"`
	GitHTTPURL        string `json:"git_http_url"`
	Namespace         string `json:"namespace"`
	VisibilityLevel   int64  `json:"visibility_level"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
	Homepage          string `json:"homepage"`
	URL               string `json:"url"`
	SSHURL            string `json:"ssh_url"`
	HTTPURL           string `json:"http_url"`
}

// Repository def
type Repository struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	Description     string `json:"description"`
	Homepage        string `json:"homepage"`
	GitHTTPURL      string `json:"git_http_url"`
	GitSSHURL       string `json:"git_ssh_url"`
	VisibilityLevel int64  `json:"visibility_level"`
}

// Commit def
type Commit struct {
	ID        string   `json:"id"`
	Message   string   `json:"message"`
	Timestamp string   `json:"timestamp"`
	URL       string   `json:"url"`
	Author    *Author  `json:"author"`
	Added     []string `json:"added"`
	Modified  []string `json:"modified"`
	Removed   []string `json:"removed"`
}

// Author def
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// WebHook def
type WebHook struct {
	URL                      string `json:"url,omitempty"`
	PushEvents               bool   `json:"push_events,omitempty"`
	IssuesEvents             bool   `json:"issues_events,omitempty"`
	ConfidentialIssuesEvents bool   `json:"confidential_issues_events,omitempty"`
	MergeRequestsEvents      bool   `json:"merge_requests_events,omitempty"`
	TagPushEvents            bool   `json:"tag_push_events,omitempty"`
	NoteEvents               bool   `json:"note_events,omitempty"`
	JobEvents                bool   `json:"job_events,omitempty"`
	PipelineEvents           bool   `json:"pipeline_events,omitempty"`
	WikiPageEvents           bool   `json:"wiki_page_events,omitempty"`
}

// RepoInfo ...
type RepoInfo struct {
	Group  string `json:"group"`
	Name   string `json:"name"`
	Branch string `json:"branch"`
}
