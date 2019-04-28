package model

// ReleaseResponse Release Response
type ReleaseResponse struct {
	Status int             `json:"status"`
	Data   *ReleaseWrapper `json:"data"`
	Info   string          `json:"info"`
}

// ReleaseWrapper Release Wrapper
type ReleaseWrapper struct {
	Release *Release `json:"Release"`
}

// WorkspaceUserResponse Workspace User Response
type WorkspaceUserResponse struct {
	Status int                     `json:"status"`
	Data   []*WorkspaceUserWrapper `json:"data"`
	Info   string                  `json:"info"`
}

// WorkspaceUserWrapper Workspace User Wrapper
type WorkspaceUserWrapper struct {
	UserWrapper *UserWrapper `json:"UserWorkspace"`
}

// UserWrapper User Wrapper
type UserWrapper struct {
	User string `json:"user"`
}

// Release Release
type Release struct {
	ID          string `json:"id"`
	WorkSpaceID string `json:"workspace_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   string `json:"startdate"`
	EndDate     string `json:"enddate"`
	Creator     string `json:"creator"`
	Created     string `json:"created"`
	Modified    string `json:"modified"`
	Status      string `json:"status"`
}

// BugResponse Bug Response
type BugResponse struct {
	Status int           `json:"status"`
	Data   []*BugWrapper `json:"data"`
	Info   string        `json:"info"`
}

// BugSingleResponse Bug Response
type BugSingleResponse struct {
	Status int         `json:"status"`
	Data   *BugWrapper `json:"data"`
	Info   string      `json:"info"`
}

// BugWrapper Bug Wrapper
type BugWrapper struct {
	Bug *Bug `json:"Bug"`
}

// Bug Bug
type Bug struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Priority         string `json:"priority"`
	Severity         string `json:"severity"`
	Module           string `json:"module"`
	Status           string `json:"status"`
	Reporter         string `json:"reporter"`
	Deadline         string `json:"deadline"`
	Created          string `json:"created"`
	BugType          string `json:"bugtype"`
	Resolved         string `json:"resolved"`
	Closed           string `json:"closed"`
	Modified         string `json:"modified"`
	LastModify       string `json:"lastmodify"`
	Auditer          string `json:"auditer"`
	DE               string `json:"de"`
	VersionTest      string `json:"version_test"`
	VersionReport    string `json:"version_report"`
	VersionClose     string `json:"version_close"`
	VersionFix       string `json:"version_fix"`
	BaselineFind     string `json:"baseline_find"`
	BaselineJoin     string `json:"baseline_join"`
	BaselineClose    string `json:"baseline_close"`
	BaselineTest     string `json:"baseline_test"`
	SourcePhase      string `json:"sourcephase"`
	TE               string `json:"te"`
	CurrentOwner     string `json:"current_owner"`
	IterationID      string `json:"iteration_id"`
	Resolution       string `json:"resolution"`
	Source           string `json:"source"`
	OriginPhase      string `json:"originphase"`
	Confirmer        string `json:"confirmer"`
	Milestone        string `json:"milestone"`
	Participator     string `json:"participator"`
	Closer           string `json:"closer"`
	Platform         string `json:"platform"`
	OS               string `json:"os"`
	TestType         string `json:"testtype"`
	TestPhase        string `json:"testphase"`
	Frequency        string `json:"frequency"`
	CC               string `json:"cc"`
	RegressionNumber string `json:"regression_number"`
	Flows            string `json:"flows"`
	Feature          string `json:"feature"`
	TestMode         string `json:"testmode"`
	Estimate         string `json:"estimate"`
	IssueID          string `json:"issue_id"`
	CreatedFrom      string `json:"created_from"`
	InProgressTime   string `json:"in_progress_time"`
	VerifyTime       string `json:"verify_time"`
	RejectTime       string `json:"reject_time"`
	ReopenTime       string `json:"reopen_time"`
	AuditTime        string `json:"audit_time"`
	SuspendTime      string `json:"suspend_time"`
	Due              string `json:"due"`
	Begin            string `json:"begin"`
	ReleaseID        string `json:"release_id"`
	WorkspaceID      string `json:"workspace_id"`
	CustomFieldOne   string `json:"custom_field_one"`
	CustomFieldTwo   string `json:"custom_field_two"`
	CustomFieldThree string `json:"custom_field_three"`
	CustomFieldFour  string `json:"custom_field_four"`
	CustomFieldFive  string `json:"custom_field_five"`
	CustomField6     string `json:"custom_field_6"`
	CustomField7     string `json:"custom_field_7"`
	CustomField8     string `json:"custom_field_8"`
	CustomField9     string `json:"custom_field_9"`
	CustomField10    string `json:"custom_field_10"`
	CustomField11    string `json:"custom_field_11"`
	CustomField12    string `json:"custom_field_12"`
	CustomField13    string `json:"custom_field_13"`
	CustomField14    string `json:"custom_field_14"`
	CustomField15    string `json:"custom_field_15"`
	CustomField16    string `json:"custom_field_16"`
	CustomField17    string `json:"custom_field_17"`
	CustomField18    string `json:"custom_field_18"`
	CustomField19    string `json:"custom_field_19"`
	CustomField20    string `json:"custom_field_20"`
	CustomField21    string `json:"custom_field_21"`
	CustomField22    string `json:"custom_field_22"`
	CustomField23    string `json:"custom_field_23"`
	CustomField24    string `json:"custom_field_24"`
	CustomField25    string `json:"custom_field_25"`
	CustomField26    string `json:"custom_field_26"`
	CustomField27    string `json:"custom_field_27"`
	CustomField28    string `json:"custom_field_28"`
	CustomField29    string `json:"custom_field_29"`
	CustomField30    string `json:"custom_field_30"`
	CustomField31    string `json:"custom_field_31"`
	CustomField32    string `json:"custom_field_32"`
	CustomField33    string `json:"custom_field_33"`
	CustomField34    string `json:"custom_field_34"`
	CustomField35    string `json:"custom_field_35"`
	CustomField36    string `json:"custom_field_36"`
	CustomField37    string `json:"custom_field_37"`
	CustomField38    string `json:"custom_field_38"`
	CustomField39    string `json:"custom_field_39"`
	CustomField40    string `json:"custom_field_40"`
	CustomField41    string `json:"custom_field_41"`
	CustomField42    string `json:"custom_field_42"`
	CustomField43    string `json:"custom_field_43"`
	CustomField44    string `json:"custom_field_44"`
	CustomField45    string `json:"custom_field_45"`
	CustomField46    string `json:"custom_field_46"`
	CustomField47    string `json:"custom_field_47"`
	CustomField48    string `json:"custom_field_48"`
	CustomField49    string `json:"custom_field_49"`
	CustomField50    string `json:"custom_field_50"`
}

// UpdateBug Update Bug
type UpdateBug struct {
	*Bug
	CurrentUser string `json:"current_user"`
}
