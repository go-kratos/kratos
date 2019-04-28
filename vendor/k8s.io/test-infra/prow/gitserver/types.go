package gitserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// GenericCommentEventAction coerces multiple actions into its generic equivalent.
type GenericCommentEventAction string

// Comments indicate values that are coerced to the specified value.
const (
	// GenericCommentActionCreated means something was created/opened/submitted
	GenericCommentActionCreated GenericCommentEventAction = "created" // "opened", "submitted"
	// GenericCommentActionEdited means something was edited.
	GenericCommentActionEdited = "edited"
	// GenericCommentActionDeleted means something was deleted/dismissed.
	GenericCommentActionDeleted = "deleted" // "dismissed"
)

// These are possible State entries for a Status.
const (
	StatusPending = "pending"
	StatusSuccess = "success"
	StatusError   = "error"
	StatusFailure = "failure"
)

type Provider string

var (
	// FoundingYear is the year GitHub was founded. This is just used so that
	// we can lower bound dates related to PRs and issues.
	FoundingYear, _ = time.Parse(SearchTimeFormat, "2007-01-01T00:00:00Z")
)

// Update action
var (
	UpdateClose  = String("close")
	UpdateReopen = String("reopen")
)

const (
	Gitlab Provider = "gitlab"
	Github          = "github"
)

const (
	// RoleAll lists both members and admins
	RoleAll = "all"
	// RoleAdmin specifies the user is an org admin, or lists only admins
	RoleAdmin = "admin"
	// RoleMaintainer specifies the user is a team maintainer, or lists only maintainers
	RoleMaintainer = "maintainer"
	// RoleMember specifies the user is a regular user, or only lists regular users
	RoleMember = "member"
	// StatePending specifies the user has an invitation to the org/team.
	StatePending = "pending"
	// StateActive specifies the user's membership is active.
	StateActive = "active"
)

const (
	// EventGUID is sent by Github in a header of every webhook request.
	// Used as a log field across prow.
	EventGUID = "event-GUID"
	// PrLogField is the number of a PR.
	// Used as a log field across prow.
	PrLogField = "pr"
	// OrgLogField is the organization of a PR.
	// Used as a log field across prow.
	OrgLogField = "org"
	// RepoLogField is the repository of a PR.
	// Used as a log field across prow.
	RepoLogField = "repo"

	// SearchTimeFormat is a time.Time format string for ISO8601 which is the
	// format that GitHub requires for times specified as part of a search query.
	SearchTimeFormat = "2006-01-02T15:04:05Z"
)

// NormLogin normalizes GitHub login strings
var NormLogin = strings.ToLower

type IssueEventAction string

const (
	// IssueActionAssigned means assignees were added.
	IssueActionAssigned IssueEventAction = "assigned"
	// IssueActionUnassigned means assignees were added.
	IssueActionUnassigned = "unassigned"
	// IssueActionLabeled means labels were added.
	IssueActionLabeled = "labeled"
	// IssueActionUnlabeled means labels were removed.
	IssueActionUnlabeled = "unlabeled"
	// IssueActionOpened means issue was opened/created.
	IssueActionOpened = "opened"
	// IssueActionEdited means issue body was edited.
	IssueActionEdited = "edited"
	// IssueActionMilestoned means the milestone was added/changed.
	IssueActionMilestoned = "milestoned"
	// IssueActionDemilestoned means a milestone was removed.
	IssueActionDemilestoned = "demilestoned"
	// IssueActionClosed means issue was closed.
	IssueActionClosed = "closed"
	// IssueActionReopened means issue was reopened.
	IssueActionReopened = "reopened"
)

// ReviewEventAction enumerates the triggers for this
// webhook payload type. See also:
// https://developer.github.com/v3/activity/events/types/#pullrequestreviewevent
type ReviewEventAction string

const (
	// ReviewActionSubmitted means the review was submitted.
	ReviewActionSubmitted ReviewEventAction = "submitted"
	// ReviewActionEdited means the review was edited.
	ReviewActionEdited = "edited"
	// ReviewActionDismissed means the review was dismissed.
	ReviewActionDismissed = "dismissed"
)

// ReviewCommentEventAction enumerates the triggers for this
// webhook payload type. See also:
// https://developer.github.com/v3/activity/events/types/#pullrequestreviewcommentevent
type ReviewCommentEventAction string

const (
	// ReviewCommentActionCreated means the comment was created.
	ReviewCommentActionCreated ReviewCommentEventAction = "created"
	// ReviewCommentActionEdited means the comment was edited.
	ReviewCommentActionEdited = "edited"
	// ReviewCommentActionDeleted means the comment was deleted.
	ReviewCommentActionDeleted = "deleted"
)

type Event string

const (
	IssueEvents                    Event = "issues"
	IssueCommentEvents                   = "issue_comment"
	PullRequestEvents                    = "pull_request"
	PullRequestReviewEvents              = "pull_request_review"
	PullRequestReviewCommentEvents       = "pull_request_review_comment"
	PushEvents                           = "push"
	StatusEvents                         = "status"
)

// PullRequestEventAction enumerates the triggers for this
// webhook payload type. See also:
// https://developer.github.com/v3/activity/events/types/#pullrequestevent
type PullRequestEventAction string

const (
	// PullRequestActionAssigned means assignees were added.
	PullRequestActionAssigned PullRequestEventAction = "assigned"
	// PullRequestActionUnassigned means assignees were removed.
	PullRequestActionUnassigned = "unassigned"
	// PullRequestActionReviewRequested means review requests were added.
	PullRequestActionReviewRequested = "review_requested"
	// PullRequestActionReviewRequestRemoved means review requests were removed.
	PullRequestActionReviewRequestRemoved = "review_request_removed"
	// PullRequestActionLabeled means means labels were added.
	PullRequestActionLabeled = "labeled"
	// PullRequestActionUnlabeled means labels were removed
	PullRequestActionUnlabeled = "unlabeled"
	// PullRequestActionOpened means the PR was created
	PullRequestActionOpened = "opened"
	// PullRequestActionEdited means means the PR body changed.
	PullRequestActionEdited = "edited"
	// PullRequestActionClosed means the PR was closed (or was merged).
	PullRequestActionClosed = "closed"
	// PullRequestActionReopened means the PR was reopened.
	PullRequestActionReopened = "reopened"
	// PullRequestActionSynchronize means the git state changed.
	PullRequestActionSynchronize = "synchronize"
)

func ToPullRequestEventAction(action string) (PullRequestEventAction, error) {
	switch action {
	case "assigned":
		return PullRequestActionAssigned, nil
	case "unassigned":
		return PullRequestActionUnassigned, nil
	case "review_requested":
		return PullRequestActionReviewRequested, nil
	case "review_request_removed":
		return PullRequestActionReviewRequestRemoved, nil
	case "labeled":
		return PullRequestActionLabeled, nil
	case "unlabeled":
		return PullRequestActionUnlabeled, nil
	case "opened":
		return PullRequestActionOpened, nil
	case "edited":
		return PullRequestActionEdited, nil
	case "closed":
		return PullRequestActionClosed, nil
	case "reopened":
		return PullRequestActionReopened, nil
	case "synchronize":
		return PullRequestActionSynchronize, nil
	}
	return PullRequestActionOpened, errors.New("invaild PullRequestEventAction")
}

// ReviewState is the state a review can be in.
type ReviewState string

// Possible review states.
const (
	ReviewStateApproved         ReviewState = "APPROVED"
	ReviewStateChangesRequested             = "CHANGES_REQUESTED"
	ReviewStateCommented                    = "COMMENTED"
	ReviewStateDismissed                    = "DISMISSED"
	ReviewStatePending                      = "PENDING"
)

// IssueEvent represents an issue event from a webhook payload (not from the events API).
type IssueEvent struct {
	GitType Provider
	//IssueEvent
	Action IssueEventAction `json:"action"`
	Issue  Issue            `json:"issue"`
	Repo   Repo             `json:"repository"`
	// Label is specified for IssueActionLabeled and IssueActionUnlabeled events.
	Label Label `json:"label"`

	// GUID is included in the header of the request received by Github.
	GUID string
}

// Issue represents general info about an issue.
type Issue struct {
	GitType Provider

	User      User      `json:"user"`
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	HTMLURL   string    `json:"html_url"`
	Labels    []Label   `json:"labels"`
	Assignees []User    `json:"assignees"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Milestone Milestone `json:"milestone"`

	// This will be non-nil if it is a pull request.
	PullRequest *struct{} `json:"pull_request,omitempty"`
}

// IsAssignee checks if a user is assigned to the issue.
func (i Issue) IsAssignee(login string) bool {
	for _, assignee := range i.Assignees {
		if NormLogin(login) == NormLogin(assignee.Login) {
			return true
		}
	}
	return false
}

// IsAuthor checks if a user is the author of the issue.
func (i Issue) IsAuthor(login string) bool {
	return NormLogin(i.User.Login) == NormLogin(login)
}

// IsPullRequest checks if an issue is a pull request.
func (i Issue) IsPullRequest() bool {
	return i.PullRequest != nil
}

// HasLabel checks if an issue has a given label.
func (i Issue) HasLabel(labelToFind string) bool {
	for _, label := range i.Labels {
		if strings.ToLower(label.Name) == strings.ToLower(labelToFind) {
			return true
		}
	}
	return false
}

// User is a GitHub user account.
type User struct {
	GitType   Provider
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
}

// Label describes a GitHub label.
type Label struct {
	GitType     Provider
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

// Milestone is a milestone defined on a github repository
type Milestone struct {
	GitType Provider
	Title   string `json:"title"`
	Number  int    `json:"number"`
}

// Repo contains general repository information.
type Repo struct {
	GitType       Provider
	Owner         User   `json:"owner"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	HTMLURL       string `json:"html_url"`
	Fork          bool   `json:"fork"`
	DefaultBranch string `json:"default_branch"`
	Archived      bool   `json:"archived"`
}

/*
func NewIssueEvent(payload interface{}) IssueEvent {
	switch payload.(type) {
	case github.IssueComment:
		return IssueEvent{
			Provider: Github,
			//IssueEvent: payload.(github.IssueEvent),
		}
	default:
		panic(fmt.Sprint("Don't support this type of payload(%s), please implement me", typeof(payload)))
	}
}
*/

// IssueCommentEventAction enumerates the triggers for this
// webhook payload type. See also:
// https://developer.github.com/v3/activity/events/types/#issuecommentevent
type IssueCommentEventAction string

const (
	// IssueCommentActionCreated means the comment was created.
	IssueCommentActionCreated IssueCommentEventAction = "created"
	// IssueCommentActionEdited means the comment was edited.
	IssueCommentActionEdited = "edited"
	// IssueCommentActionDeleted means the comment was deleted.
	IssueCommentActionDeleted = "deleted"
)

// IssueCommentEvent is what the provider sends us when an issue comment is changed.
type IssueCommentEvent struct {
	GitType Provider
	//github.IssueCommentEvent
	Action  IssueCommentEventAction `json:"action"`
	Issue   Issue                   `json:"issue"`
	Comment IssueComment            `json:"comment"`
	Repo    Repo                    `json:"repository"`

	// GUID is included in the header of the request received by Github.
	GUID string
}

func (ic *IssueCommentEvent) isPR() bool {
	return ic.Issue.PullRequest == nil
}

/*
func NewIssueCommentEvent(payload interface{}) IssueCommentEvent {
	switch payload.(type) {
	case github.IssueCommentEvent:
		return IssueCommentEvent{
			Provider: Github,
			//IssueCommentEvent: payload.(github.IssueCommentEvent),
		}
	default:
		panic(fmt.Sprint("Don't support this type of payload(%s), please implement me", typeof(payload)))
	}
}
*/

// PullRequestEvent is what the provider sends us when a PR is changed.
type PullRequestEvent struct {
	GitType Provider
	//github.PullRequestEvent
	Action      PullRequestEventAction `json:"action"`
	Number      int                    `json:"number"`
	PullRequest PullRequest            `json:"pull_request"`
	Repo        Repo                   `json:"repository"`
	Label       Label                  `json:"label"`
	Sender      User                   `json:"sender"`

	WorkInProgress bool `json:"work_in_progress"`
	// Changes holds raw change data, which we must inspect
	// and deserialize later as this is a polymorphic field
	Changes json.RawMessage `json:"changes"`

	// GUID is included in the header of the request received by Github.
	GUID string

	// gitlab changed label
	Labels []Label `json:"labels"`
}

/*
func NewPullRequestEvent(payload interface{}) PullRequestEvent {
	switch payload.(type) {
	case github.PullRequestEvent:
		return PullRequestEvent{
			Provider: Github,
			//PullRequestEvent: payload.(github.PullRequestEvent),
		}
	case gitlab.MergeRequestEvent:
		mr := payload.(gitlab.MergeRequestEvent)
		repo := mr.Repo.ToGithubRepo()
		action, sender := gitlab.MergeRequestAction(mr.ObjectAttributes.Action).GithubEventAction(mr)
		return PullRequestEvent{
			Provider: Gitlab,
			//PullRequestEvent: github.PullRequestEvent{
			//Action:      action,
			//Number:      mr.ObjectAttributes.IID,
			//PullRequest: mr.ObjectAttributes.ToGithubPR(mr.User, repo),
			//Repo:        repo,
			//Sender:      sender,
			//GUID:        strconv.FormatInt(time.Now().Unix(), 10),
			//},
		}
	default:
		panic(fmt.Sprint("Don't support this type of payload(%s), please implement me", typeof(payload)))
	}
}
*/

//ReviewEvent is what the provider sends us when a PR review is changed.
type ReviewEvent struct {
	GitType     Provider
	Action      ReviewEventAction `json:"action"`
	PullRequest PullRequest       `json:"pull_request"`
	Repo        Repo              `json:"repository"`
	Review      Review            `json:"review"`

	// GUID is included in the header of the request received by Github.
	GUID string
}

/*
func NewReviewEvent(payload interface{}) ReviewEvent {
	switch payload.(type) {
	case github.ReviewEvent:
		return ReviewEvent{
			Provider: Github,
			//ReviewEvent: payload.(github.ReviewEvent),
		}
	default:
		panic(fmt.Sprint("Don't support this type of payload(%s), please implement me", typeof(payload)))
	}
}
*/

// ReviewCommentEvent is what GitHub sends us when a PR review comment is changed.
type ReviewCommentEvent struct {
	GitType     Provider
	Action      ReviewCommentEventAction `json:"action"`
	PullRequest PullRequest              `json:"pull_request"`
	Repo        Repo                     `json:"repository"`
	Comment     ReviewComment            `json:"comment"`

	// GUID is included in the header of the request received by Github.
	GUID string
}

/*
func NewReviewCommentEvent(payload interface{}) ReviewCommentEvent {
	switch payload.(type) {
	case github.ReviewCommentEvent:
		return ReviewCommentEvent{
			GitType Provider: Github,
			//ReviewCommentEvent: payload.(github.ReviewCommentEvent),
		}
	default:
		panic(fmt.Sprint("Don't support this type of payload(%s), please implement me", typeof(payload)))
	}
}
*/

// PushEvent is what the provider sends us when a user pushes to a repo.
type PushEvent struct {
	GitType Provider
	//github.PushEvent
	Ref     string   `json:"ref"`
	Before  string   `json:"before"`
	After   string   `json:"after"`
	Compare string   `json:"compare"`
	Size    int64    `json:"size"`
	Commits []Commit `json:"commits"`
	// Pusher is the user that pushed the commit, valid in a webhook event.
	Pusher User `json:"pusher"`
	// Sender contains more information that Pusher about the user.
	Sender User `json:"sender"`
	Repo   Repo `json:"repository"`

	// GUID is included in the header of the request received by Github.
	GUID string
}

// Branch returns the name of the branch to which the user pushed.
func (pe PushEvent) Branch() string {
	refs := strings.Split(pe.Ref, "/")
	return refs[len(refs)-1]
}

/*
func NewPushEvent(payload interface{}) PushEvent {
	switch payload.(type) {
	case github.PushEvent:
		return PushEvent{
			GitType Provider: Github,
			//PushEvent: payload.(github.PushEvent),
		}
	case gitlab.PushEvent:
		p := payload.(gitlab.PushEvent)
		return PushEvent{
			GitType Provider: Gitlab,
			PushEvent: github.PushEvent{
				Ref:     p.Ref,
				Before:  p.Before,
				After:   p.After,
				Compare: p.Compare(),
				Repo:    p.Repo.ToGithubRepo(),
				Pusher: github.User{
					Login: p.UserName,
					Name:  p.Name,
					Email: p.UserEmail,
					ID:    p.UserID,
				},
				Sender: github.User{
					Login: p.UserName,
					Name:  p.Name,
					Email: p.UserEmail,
					ID:    p.UserID,
				},
				GUID: strconv.FormatInt(time.Now().Unix(), 10),
			},
		}
	default:
		panic(fmt.Sprint("Don't support this type of payload(%s), please implement me", typeof(payload)))
	}
}
*/

// StatusEvent fires whenever a git commit changes.
type StatusEvent struct {
	GitType Provider
	//github.StatusEvent
	SHA         string `json:"sha,omitempty"`
	State       string `json:"state,omitempty"`
	Description string `json:"description,omitempty"`
	TargetURL   string `json:"target_url,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Context     string `json:"context,omitempty"`
	Sender      User   `json:"sender,omitempty"`
	Repo        Repo   `json:"repository,omitempty"`

	// GUID is included in the header of the request received by Github.
	GUID string
}

/*
func NewStatusEvent(payload interface{}) StatusEvent {
	switch payload.(type) {
	case github.StatusEvent:
		return StatusEvent{
			GitType Provider: Github,
			//StatusEvent: payload.(github.StatusEvent),
		}
	default:
		panic(fmt.Sprint("Don't support this type of payload(%s), please implement me", typeof(payload)))
	}
}
*/

type Action struct {
	Event   Event
	Payload interface{}
	Header  http.Header
}

func (a Action) GUID() string {
	v := reflect.ValueOf(a)
	return v.FieldByName("GUID").String()
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

// PullRequest contains information about a PullRequest.
type PullRequest struct {
	GitType            Provider
	Number             int               `json:"number"`
	HTMLURL            string            `json:"html_url"`
	User               User              `json:"user"`
	Base               PullRequestBranch `json:"base"`
	Head               PullRequestBranch `json:"head"`
	Title              string            `json:"title"`
	Body               string            `json:"body"`
	RequestedReviewers []User            `json:"requested_reviewers"`
	Assignees          []User            `json:"assignees"`
	State              string            `json:"state"`
	Merged             bool              `json:"merged"`
	CreatedAt          time.Time         `json:"created_at,omitempty"`
	UpdatedAt          time.Time         `json:"updated_at,omitempty"`
	WorkInProgress     bool              `json:"work_in_progress"`
	// ref https://developer.github.com/v3/pulls/#get-a-single-pull-request
	// If Merged is true, MergeSHA is the SHA of the merge commit, or squashed commit
	// If Merged is false, MergeSHA is a commit SHA that github created to test if
	// the PR can be merged automatically.
	MergeSHA *string `json:"merge_commit_sha"`
	// ref https://developer.github.com/v3/pulls/#response-1
	// The value of the mergeable attribute can be true, false, or null. If the value
	// is null, this means that the mergeability hasn't been computed yet, and a
	// background job was started to compute it. When the job is complete, the response
	// will include a non-null value for the mergeable attribute.
	Mergable *bool `json:"mergeable,omitempty"`

	Labels []string
}

// PullRequestBranch contains information about a particular branch in a PR.
type PullRequestBranch struct {
	GitType Provider
	Ref     string `json:"ref"`
	SHA     string `json:"sha"`
	Repo    Repo   `json:"repo"`
}

// Review describes a Pull Request review.
type Review struct {
	GitType     Provider
	ID          int         `json:"id"`
	User        User        `json:"user"`
	Body        string      `json:"body"`
	State       ReviewState `json:"state"`
	HTMLURL     string      `json:"html_url"`
	SubmittedAt time.Time   `json:"submitted_at"`
}

// GenericCommentEvent is a fake event type that is instantiated for any github event that contains
// comment like content.
// The specific events that are also handled as GenericCommentEvents are:
// - issue_comment events
// - pull_request_review events
// - pull_request_review_comment events
// - pull_request events with action in ["opened", "edited"]
// - issue events with action in ["opened", "edited"]
//
// Issue and PR "closed" events are not coerced to the "deleted" Action and do not trigger
// a GenericCommentEvent because these events don't actually remove the comment content from GH.
type GenericCommentEvent struct {
	GitType      Provider
	IsPR         bool
	Action       GenericCommentEventAction
	Body         string
	HTMLURL      string
	Number       int
	Repo         Repo
	User         User
	IssueAuthor  User
	Assignees    []User
	IssueState   string
	IssueBody    string
	IssueHTMLURL string
	GUID         string
}

// ReviewComment describes a Pull Request review.
type ReviewComment struct {
	GitType   Provider
	ID        int       `json:"id"`
	ReviewID  int       `json:"pull_request_review_id"`
	User      User      `json:"user"`
	Body      string    `json:"body"`
	Path      string    `json:"path"`
	HTMLURL   string    `json:"html_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Position will be nil if the code has changed such that the comment is no
	// longer relevant.
	Position *int `json:"position"`
}

// Commit represents general info about a commit.
type Commit struct {
	GitType  Provider
	ID       string   `json:"id"`
	Message  string   `json:"message"`
	Added    []string `json:"added"`
	Removed  []string `json:"removed"`
	Modified []string `json:"modified"`
}

// IssueComment represents general info about an issue comment.
type IssueComment struct {
	GitType   Provider
	ID        int `json:"id,omitempty"`
	IssueID   int
	Body      string    `json:"body"`
	User      User      `json:"user,omitempty"`
	HTMLURL   string    `json:"html_url,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// ListedIssueEvent represents an issue event from the events API (not from a webhook payload).
// https://developer.github.com/v3/issues/events/
type ListedIssueEvent struct {
	GitType   Provider
	Event     Event     `json:"event"` // This is the same as IssueEvent.Action.
	Actor     User      `json:"actor"`
	Label     Label     `json:"label"`
	CreatedAt time.Time `json:"created_at"`
}

// IssuesSearchResult represents the result of an issues search.
type IssuesSearchResult struct {
	Total  int     `json:"total_count,omitempty"`
	Issues []Issue `json:"items,omitempty"`
}

// ReviewAction is the action that a review can be made with.
type ReviewAction string

// Possible review actions. Leave Action blank for a pending review.
const (
	Approve        ReviewAction = "APPROVE"
	RequestChanges              = "REQUEST_CHANGES"
	Comment                     = "COMMENT"
)

// DraftReviewComment is a comment in a draft review.
type DraftReviewComment struct {
	Path string `json:"path"`
	// Position in the patch, not the line number in the file.
	Position int    `json:"position"`
	Body     string `json:"body"`
}

type DraftReview struct {
	GitType Provider
	// If unspecified, defaults to the most recent commit in the PR.
	CommitSHA string `json:"commit_id,omitempty"`
	Body      string `json:"body"`
	// If unspecified, defaults to PENDING.
	Action   ReviewAction         `json:"event,omitempty"`
	Comments []DraftReviewComment `json:"comments,omitempty"`
}

// Team is a github organizational team
type Team struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	Privacy      string `json:"privacy,omitempty"`
	Parent       *Team  `json:"parent,omitempty"`         // Only present in responses
	ParentTeamID *int   `json:"parent_team_id,omitempty"` // Only valid in creates/edits
}

// TeamMember is a member of an organizational team
type TeamMember struct {
	Login string `json:"login"`
}

// Membership specifies the role and state details for an org and/or team.
type Membership struct {
	// admin or member
	Role string `json:"role"`
	// pending or active
	State string `json:"state,omitempty"`
}

// Content is some base64 encoded github file content
type Content struct {
	Content string `json:"content"`
	SHA     string `json:"sha"`
}

// Organization stores metadata information about an organization
type Organization struct {
	// BillingEmail holds private billing address
	BillingEmail string `json:"billing_email"`
	Company      string `json:"company"`
	// Email is publicly visible
	Email                        string `json:"email"`
	Location                     string `json:"location"`
	Name                         string `json:"name"`
	Description                  string `json:"description"`
	HasOrganizationProjects      bool   `json:"has_organization_projects"`
	HasRepositoryProjects        bool   `json:"has_repository_projects"`
	DefaultRepositoryPermission  string `json:"default_repository_permission"`
	MembersCanCreateRepositories bool   `json:"members_can_create_repositories"`
}

// OrgMembership contains Membership fields for user membership in an org.
type OrgMembership struct {
	Membership
}

// TeamMembership contains Membership fields for user membership on a team.
type TeamMembership struct {
	Membership
}

// OrgInvitation contains Login and other details about the invitation.
type OrgInvitation struct {
	TeamMember
	Inviter TeamMember `json:"login"`
}

// PullRequestChange contains information about what a PR changed.
type PullRequestChange struct {
	GitType          Provider
	SHA              string `json:"sha"`
	Filename         string `json:"filename"`
	Status           string `json:"status"`
	Additions        int    `json:"additions"`
	Deletions        int    `json:"deletions"`
	Changes          int    `json:"changes"`
	Patch            string `json:"patch"`
	BlobURL          string `json:"blob_url"`
	PreviousFilename string `json:"previous_filename"`
}

// Branch contains general branch information.
type Branch struct {
	GitType   Provider
	Name      string `json:"name"`
	Protected bool   `json:"protected"` // only included for ?protection=true requests
	// TODO(fejta): consider including undocumented protection key
}

// CombinedStatus is the latest statuses for a ref.
type CombinedStatus struct {
	GitType  Provider
	Statuses []Status `json:"statuses"`
}

// Status is used to set a commit status line.
type Status struct {
	GitType     Provider
	State       string `json:"state"`
	TargetURL   string `json:"target_url,omitempty"`
	Description string `json:"description,omitempty"`
	Context     string `json:"context,omitempty"`
}

// MissingUsers is an error specifying the users that could not be unassigned.
type MissingUsers struct {
	GitType Provider
	Users   []string
	action  string
}

func (m MissingUsers) Error() string {
	return fmt.Sprintf("could not %s the following user(s): %s.", m.action, strings.Join(m.Users, ", "))
}

// FileNotFound happens when github cannot find the file requested by GetFile().
type FileNotFound struct {
	Org, Repo, Path, Commit string
}

func (e *FileNotFound) Error() string {
	return fmt.Sprintf("%s/%s/%s @ %s not found", e.Org, e.Repo, e.Path, e.Commit)
}

type ListProjectMergeRequestsOptions struct {
	IIDs            []int      `url:"iids[],omitempty" json:"iids,omitempty"`
	State           *string    `url:"state,omitempty" json:"state,omitempty"`
	OrderBy         *string    `url:"order_by,omitempty" json:"order_by,omitempty"`
	Sort            *string    `url:"sort,omitempty" json:"sort,omitempty"`
	Milestone       *string    `url:"milestone,omitempty" json:"milestone,omitempty"`
	View            *string    `url:"view,omitempty" json:"view,omitempty"`
	Labels          []string   `url:"labels,omitempty" json:"labels,omitempty"`
	CreatedAfter    *time.Time `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore   *time.Time `url:"created_before,omitempty" json:"created_before,omitempty"`
	UpdatedAfter    *time.Time `url:"updated_after,omitempty" json:"updated_after,omitempty"`
	UpdatedBefore   *time.Time `url:"updated_before,omitempty" json:"updated_before,omitempty"`
	Scope           *string    `url:"scope,omitempty" json:"scope,omitempty"`
	AuthorID        *int       `url:"author_id,omitempty" json:"author_id,omitempty"`
	AssigneeID      *int       `url:"assignee_id,omitempty" json:"assignee_id,omitempty"`
	MyReactionEmoji *string    `url:"my_reaction_emoji,omitempty" json:"my_reaction_emoji,omitempty"`
	SourceBranch    *string    `url:"source_branch,omitempty" json:"source_branch,omitempty"`
	TargetBranch    *string    `url:"target_branch,omitempty" json:"target_branch,omitempty"`
	Search          *string    `url:"search,omitempty" json:"search,omitempty"`
}

// MergeDetails contains desired properties of the merge.
//
// See https://developer.github.com/v3/pulls/#merge-a-pull-request-merge-button
type MergeDetails struct {
	GitType Provider

	// CommitMessage defaults to the automatic message.
	CommitMessage string `json:"commit_message,omitempty"`
	// The PR HEAD must match this to prevent races.
	SHA string `json:"sha,omitempty"`

	// github
	// CommitTitle defaults to the automatic message.
	CommitTitle string `json:"commit_title,omitempty"`
	// Can be "merge", "squash", or "rebase". Defaults to merge.
	MergeMethod string `json:"merge_method,omitempty"`

	// gitlab
	ShouldRemoveSourceBranch  *bool `url:"should_remove_source_branch,omitempty" json:"should_remove_source_branch,omitempty"`
	MergeWhenPipelineSucceeds *bool `url:"merge_when_pipeline_succeeds,omitempty" json:"merge_when_pipeline_succeeds,omitempty"`
}

// BranchProtectionRequest represents
// protections in place for a branch.
// See also: https://developer.github.com/v3/repos/branches/#update-branch-protection
type BranchProtectionRequest struct {
	RequiredStatusChecks       *RequiredStatusChecks       `json:"required_status_checks"`
	EnforceAdmins              *bool                       `json:"enforce_admins"`
	RequiredPullRequestReviews *RequiredPullRequestReviews `json:"required_pull_request_reviews"`
	Restrictions               *Restrictions               `json:"restrictions"`
}

type Request struct {
	Method      string
	Path        string
	Accept      string
	RequestBody interface{}
	ExitCodes   []int
}

// SingleCommit is the commit part received when requesting a single commit
// https://developer.github.com/v3/repos/commits/#get-a-single-commit
type SingleCommit struct {
	Commit struct {
		Tree struct {
			SHA string `json:"sha"`
		} `json:"tree"`
	} `json:"commit"`
}

// RequiredStatusChecks specifies which contexts must pass to merge.
type RequiredStatusChecks struct {
	Strict   bool     `json:"strict"` // PR must be up to date (include latest base branch commit).
	Contexts []string `json:"contexts"`
}

// RequiredPullRequestReviews controls review rights.
type RequiredPullRequestReviews struct {
	DismissalRestrictions        Restrictions `json:"dismissal_restrictions"`
	DismissStaleReviews          bool         `json:"dismiss_stale_reviews"`
	RequireCodeOwnerReviews      bool         `json:"require_code_owner_reviews"`
	RequiredApprovingReviewCount int          `json:"required_approving_review_count"`
}

// Restrictions tells github to restrict an activity to people/teams.
//
// Use *[]string in order to distinguish unset and empty list.
// This is needed by dismissal_restrictions to distinguish
// do not restrict (empty object) and restrict everyone (nil user/teams list)
type Restrictions struct {
	Users *[]string `json:"users,omitempty"`
	Teams *[]string `json:"teams,omitempty"`
}

// LabelNotFound indicates that a label is not attached to an issue. For example, removing a
// label from an issue, when the issue does not have that label.
type LabelNotFound struct {
	Owner, Repo string
	Number      int
	Label       string
}

func (e *LabelNotFound) Error() string {
	return fmt.Sprintf("label %q does not exist on %s/%s/%d", e.Label, e.Owner, e.Repo, e.Number)
}

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}

type PullRequestMergeType string

// Possible types of merges for the GitHub merge API
const (
	MergeMerge  PullRequestMergeType = "merge"
	MergeRebase PullRequestMergeType = "rebase"
	MergeSquash PullRequestMergeType = "squash"
)

// ModifiedHeadError happens when github refuses to merge a PR because the PR changed.
type ModifiedHeadError string

func (e ModifiedHeadError) Error() string { return string(e) }

// UnmergablePRError happens when github refuses to merge a PR for other reasons (merge confclit).
type UnmergablePRError string

func (e UnmergablePRError) Error() string { return string(e) }

// UnmergablePRBaseChangedError happens when github refuses merging a PR because the base changed.
type UnmergablePRBaseChangedError string

func (e UnmergablePRBaseChangedError) Error() string { return string(e) }

// UnauthorizedToPushError happens when client is not allowed to push to github.
type UnauthorizedToPushError string

func (e UnauthorizedToPushError) Error() string { return string(e) }
