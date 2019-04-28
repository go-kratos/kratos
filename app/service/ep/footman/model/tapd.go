package model

//IterationResponse response for tapd iteration query
type IterationResponse struct {
	Status int                 `json:"status"`
	Data   []*IterationWrapper `json:"data"`
	Info   string              `json:"info"`
}

//IterationWrapper sub struct in IterationResponse
type IterationWrapper struct {
	Iteration *Iteration `json:"iteration"`
}

//Iteration tapd iteration
//type Iteration struct {
//	ID        string `json:"id"`
//	Name      string `json:"name"`
//	StartDate string `json:"startdate"`
//	EndDate   string `json:"enddate"`
//}

//Iteration tapd iteration
type Iteration struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	WorkspaceID   string `json:"workspace_id"`
	Startdate     string `json:"startdate"`
	Enddate       string `json:"enddate"`
	Status        string `json:"status"`
	ReleaseID     string `json:"release_id"`
	Description   string `json:"description"`
	Creator       string `json:"creator"`
	Created       string `json:"created"`
	Modified      string `json:"modified"`
	Completed     string `json:"completed"`
	CustomField1  string `json:"custom_field_1"`
	CustomField2  string `json:"custom_field_2"`
	CustomField3  string `json:"custom_field_3"`
	CustomField4  string `json:"custom_field_4"`
	CustomField5  string `json:"custom_field_5"`
	CustomField6  string `json:"custom_field_6"`
	CustomField7  string `json:"custom_field_7"`
	CustomField8  string `json:"custom_field_8"`
	CustomField9  string `json:"custom_field_9"`
	CustomField10 string `json:"custom_field_10"`
	CustomField11 string `json:"custom_field_11"`
	CustomField12 string `json:"custom_field_12"`
	CustomField13 string `json:"custom_field_13"`
	CustomField14 string `json:"custom_field_14"`
	CustomField15 string `json:"custom_field_15"`
	CustomField16 string `json:"custom_field_16"`
	CustomField17 string `json:"custom_field_17"`
	CustomField18 string `json:"custom_field_18"`
	CustomField19 string `json:"custom_field_19"`
	CustomField20 string `json:"custom_field_20"`
	CustomField21 string `json:"custom_field_21"`
	CustomField22 string `json:"custom_field_22"`
	CustomField23 string `json:"custom_field_23"`
	CustomField24 string `json:"custom_field_24"`
	CustomField25 string `json:"custom_field_25"`
	CustomField26 string `json:"custom_field_26"`
	CustomField27 string `json:"custom_field_27"`
	CustomField28 string `json:"custom_field_28"`
	CustomField29 string `json:"custom_field_29"`
	CustomField30 string `json:"custom_field_30"`
	CustomField31 string `json:"custom_field_31"`
	CustomField32 string `json:"custom_field_32"`
	CustomField33 string `json:"custom_field_33"`
	CustomField34 string `json:"custom_field_34"`
	CustomField35 string `json:"custom_field_35"`
	CustomField36 string `json:"custom_field_36"`
	CustomField37 string `json:"custom_field_37"`
	CustomField38 string `json:"custom_field_38"`
	CustomField39 string `json:"custom_field_39"`
	CustomField40 string `json:"custom_field_40"`
	CustomField41 string `json:"custom_field_41"`
	CustomField42 string `json:"custom_field_42"`
	CustomField43 string `json:"custom_field_43"`
	CustomField44 string `json:"custom_field_44"`
	CustomField45 string `json:"custom_field_45"`
	CustomField46 string `json:"custom_field_46"`
	CustomField47 string `json:"custom_field_47"`
	CustomField48 string `json:"custom_field_48"`
	CustomField49 string `json:"custom_field_49"`
	CustomField50 string `json:"custom_field_50"`
}

//StoryResponse response for tapd multiple stories query
type StoryResponse struct {
	Status int             `json:"status"`
	Data   []*StoryWrapper `json:"data"`
	Info   string          `json:"info"`
}

//SpecStoryResponse response for tapd specific story query
type SpecStoryResponse struct {
	Status int           `json:"status"`
	Data   *StoryWrapper `json:"data"`
	Info   string        `json:"info"`
}

//StoryWrapper sub struct in story response
type StoryWrapper struct {
	Story *Story `json:"story"`
}

//Story tapd story
type Story struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	WorkspaceID      string `json:"workspace_id"`
	Creator          string `json:"creator"`
	Created          string `json:"created"`
	Modified         string `json:"modified"`
	Status           string `json:"status"`
	Owner            string `json:"owner"`
	Cc               string `json:"cc"`
	Begin            string `json:"begin"`
	Due              string `json:"due"`
	Size             string `json:"size"`
	Priority         string `json:"priority"`
	Developer        string `json:"developer"`
	IterationID      string `json:"iteration_id"`
	TestFocus        string `json:"test_focus"`
	Type             string `json:"type"`
	Source           string `json:"source"`
	Module           string `json:"module"`
	Version          string `json:"version"`
	Completed        string `json:"completed"`
	CategoryID       string `json:"category_id"`
	ParentID         string `json:"parent_id"`
	ChildrenID       string `json:"children_id"`
	AncestorID       string `json:"ancestor_id"`
	BusinessValue    string `json:"business_value"`
	Effort           string `json:"effort"`
	EffortCompleted  string `json:"effort_completed"`
	Exceed           string `json:"exceed"`
	Remain           string `json:"remain"`
	ReleaseID        string `json:"release_id"`
	CustomFieldOne   string `json:"custom_field_one"`
	CustomFieldTwo   string `json:"custom_field_two"`
	CustomFieldThree string `json:"custom_field_three"`
	CustomFieldFour  string `json:"custom_field_four"`
	CustomFieldFive  string `json:"custom_field_five"`
	CustomFieldSix   string `json:"custom_field_six"`
	CustomFieldSeven string `json:"custom_field_seven"`
	CustomFieldEight string `json:"custom_field_eight"`
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
	CustomField51    string `json:"custom_field_51"`
	CustomField52    string `json:"custom_field_52"`
	CustomField53    string `json:"custom_field_53"`
	CustomField54    string `json:"custom_field_54"`
	CustomField55    string `json:"custom_field_55"`
	CustomField56    string `json:"custom_field_56"`
	CustomField57    string `json:"custom_field_57"`
	CustomField58    string `json:"custom_field_58"`
	CustomField59    string `json:"custom_field_59"`
	CustomField60    string `json:"custom_field_60"`
	CustomField61    string `json:"custom_field_61"`
	CustomField62    string `json:"custom_field_62"`
	CustomField63    string `json:"custom_field_63"`
	CustomField64    string `json:"custom_field_64"`
	CustomField65    string `json:"custom_field_65"`
	CustomField66    string `json:"custom_field_66"`
	CustomField67    string `json:"custom_field_67"`
	CustomField68    string `json:"custom_field_68"`
	CustomField69    string `json:"custom_field_69"`
	CustomField70    string `json:"custom_field_70"`
	CustomField71    string `json:"custom_field_71"`
	CustomField72    string `json:"custom_field_72"`
	CustomField73    string `json:"custom_field_73"`
	CustomField74    string `json:"custom_field_74"`
	CustomField75    string `json:"custom_field_75"`
	CustomField76    string `json:"custom_field_76"`
	CustomField77    string `json:"custom_field_77"`
	CustomField78    string `json:"custom_field_78"`
	CustomField79    string `json:"custom_field_79"`
	CustomField80    string `json:"custom_field_80"`
	CustomField81    string `json:"custom_field_81"`
	CustomField82    string `json:"custom_field_82"`
	CustomField83    string `json:"custom_field_83"`
	CustomField84    string `json:"custom_field_84"`
	CustomField85    string `json:"custom_field_85"`
	CustomField86    string `json:"custom_field_86"`
	CustomField87    string `json:"custom_field_87"`
	CustomField88    string `json:"custom_field_88"`
	CustomField89    string `json:"custom_field_89"`
	CustomField90    string `json:"custom_field_90"`
	CustomField91    string `json:"custom_field_91"`
	CustomField92    string `json:"custom_field_92"`
	CustomField93    string `json:"custom_field_93"`
	CustomField94    string `json:"custom_field_94"`
	CustomField95    string `json:"custom_field_95"`
	CustomField96    string `json:"custom_field_96"`
	CustomField97    string `json:"custom_field_97"`
	CustomField98    string `json:"custom_field_98"`
	CustomField99    string `json:"custom_field_99"`
	CustomField100   string `json:"custom_field_100"`
}

//IOSStory additional fields for ios story
type IOSStory struct {
	CustomField99 string `json:"custom_field_99"` //接口上线日
	CustomField97 string `json:"custom_field_97"` //双端都提得需求
	CustomField93 string `json:"custom_field_93"` //端范围（默认仅粉iPhone）
	CustomField92 string `json:"custom_field_92"` //是否可以单端上线
}

//AndroidStory additional fields for android story
type AndroidStory struct {
	CustomField99 string `json:"custom_field_99"` //接口上线日
	CustomField97 string `json:"custom_field_97"` //双端都提得需求
	CustomField93 string `json:"custom_field_93"` //是否可以单端上线
}

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

//StoryChangeResponse response for tapd story change query
type StoryChangeResponse struct {
	Status int                      `json:"status"`
	Data   []*WorkitemChangeWrapper `json:"data"`
	Info   string                   `json:"info"`
}

//WorkitemChangeWrapper sub struct in StoryChangeResponse
type WorkitemChangeWrapper struct {
	WorkitemChange *WorkitemChange `json:"WorkitemChange"`
}

//WorkitemChange sub struct in WorkitemChangeWrapper
type WorkitemChange struct {
	ID           string `json:"id"`
	WorkspaceID  string `json:"workspace_id"`
	Creator      string `json:"creator"`
	Created      string `json:"created"`
	ChangeSummay string `json:"change_summay"`
	Comment      string `json:"comment"`
	Changes      string `json:"changes"`
	EntityType   string `json:"entity_type"`
	StoryID      string `json:"story_id"`
}

//StoryChangeItem story change struct wrote to change file
type StoryChangeItem struct {
	ID           string
	WorkspaceID  string
	StoryID      string
	Number       string
	Field        string
	Creator      string
	Created      string
	ValueBefore  string
	ValueAfter   string
	ChangeSummay string
	Comment      string
	EntityType   string
}

//StoryChangeByIteration story changes organized by iteration
type StoryChangeByIteration struct {
	IterationName   string
	StoryCount      int
	StoryChangeList []*TargetStoryChange
}

//TargetStoryChange story and story changes
type TargetStoryChange struct {
	Story         *Story
	StatusChanges []*StatusChange
}

//StatusChange story change
type StatusChange struct {
	Creator     string
	Created     string
	ValueBefore string
	ValueAfter  string
}

//NameMapResponse story status name mapping
type NameMapResponse struct {
	Status int               `json:"status"`
	Data   map[string]string `json:"data"`
	Info   string
}

//RejectedStoryByIteration rejected stories organized by iteration
type RejectedStoryByIteration struct {
	IterationName      string
	RejectedStoryCount int
	RejectedStoryList  []string
}

//TestTimeByIteration stories' test time info organized by iteration
type TestTimeByIteration struct {
	IterationName string
	StoryCount    int
	TimeByStroy   []*TestTimeByStory
}

//TestTimeByStory story base info and test time
type TestTimeByStory struct {
	StoryName   string
	StorySize   string
	StoryEffort string
	TestTime    float64
}

//WaitTimeByIteration stories' wait time organized by iteration
type WaitTimeByIteration struct {
	IterationName string
	StoryCount    int
	TimeByStroy   []*WaitTimeByStory
}

//WaitTimeByStory story base info and wait time
type WaitTimeByStory struct {
	StoryName   string
	StorySize   string
	StoryEffort string
	WaitTime    float64
}

//CategoryResponse response for tapd category query
type CategoryResponse struct {
	Status int                `json:"status"`
	Data   []*CategoryWrapper `json:"data"`
	Info   string             `json:"info"`
}

//CategoryPreResponse response for tapd category query
type CategoryPreResponse struct {
	Status int              `json:"status"`
	Data   *CategoryWrapper `json:"data"`
	Info   string           `json:"info"`
}

//CategoryWrapper sub struct in CategoryResponse
type CategoryWrapper struct {
	Category *Category
}

//Category project category
type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
