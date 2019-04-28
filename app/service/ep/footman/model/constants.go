package model

//tapd urls
const (
	IterationURL         = "https://api.tapd.cn/iterations?workspace_id=%s&name=%s&limit=%d&fields=id,name,startdate,enddate"
	AllIterationURL      = "https://api.tapd.cn/iterations?workspace_id=%s&limit=%d"
	StoryURL             = "https://api.tapd.cn/stories?workspace_id=%s&iteration_id=%s&limit=%d"
	StoryURLWithoutItera = "https://api.tapd.cn/stories?workspace_id=%s&limit=%d"
	SpecStoryURL         = "https://api.tapd.cn/stories?workspace_id=%s&id=%s"
	StoryChangeURL       = "https://api.tapd.cn/story_changes?workspace_id=%s&story_id=%s&limit=%d"
	NameMapURL           = "https://api.tapd.cn/workflows/status_map?system=story&workspace_id=%s"
	CategoryURL          = "https://api.tapd.cn/story_categories?workspace_id=%s&fields=id,name&limit=%d"
	BugURL               = "https://api.tapd.cn/bugs?workspace_id=%s&limit=%d"
	ReleaseURL           = "https://api.tapd.cn/releases?workspace_id=%s&id=%s"
	CategoryPreURL       = "https://api.tapd.cn/story_categories?workspace_id=%s&id=%s"
	BugPreURL            = "https://api.tapd.cn/bugs?workspace_id=%s&id=%s"
	CreateBugURL         = "https://api.tapd.cn/bugs"
)

//page size
const (
	IPS  = 30
	SPS  = 30
	SCPS = 50
	CPS  = 50
)

//ios and and android info
var (
	IOSRelease        = map[string]string{"1120055921001001566": "粉 - 5.22", "1120055921001001673": "粉 - 5.23", "1120055921001001766": "粉 - 5.24", "1120055921001001856": "粉 - 5.25", "1120055921001001930": "粉 - 5.26", "1120055921001002000": "粉 - 5.27", "1120055921001002073": "粉 - 5.28", "1120055921001002147": "粉 - 5.29", "1120055921001002213": "粉 - 5.30", "1120055921001002234": "粉 - 5.31"}
	AndroidRelease    = map[string]string{"1120060791001001567": "Android 5.22", "1120060791001001672": "Android 5.23", "1120060791001001731": "Android 5.24", "1120060791001001855": "Android 5.25", "1120060791001001931": "Android 5.26", "1120060791001001999": "Android 5.27", "1120060791001002072": "Android 5.28", "1120060791001002158": "Android 5.29", "1120060791001002212": "Android 5.30", "1120060791001002268": "Android 5.31"}
	IOSWorkflow       = []string{"规划中（研发不受理此状态！）", "需求评审", "需求评审未通过", "确定发布计划", "待开发", "开发中", "产品设计验收", "验收通过待测试", "测试中", "验收通过待测试（接入需求）", "测试通过待合入", "已合入总包", "需求完成", "需求挂起", "需求取消", "免测发布"}
	AndroidWorkflow   = []string{"规划中（研发不受理此状态！）", "需求评审", "需求评审未通过", "确定发布计划", "待开发", "开发中", "产品设计验收", "验收通过待测试", "测试中", "验收通过待测试（接入需求）", "测试通过", "待合入总包", "需求完成", "需求挂起", "需求取消", "免测发布"}
	BaseFields        = map[string]string{"id": "ID", "status": "状态", "priority": "优先级", "size": "规模", "category_id": "需求分类", "parent_id": "父需求", "release_id": "发布计划", "owner": "处理人", "developer": "开发人员", "creator": "创建人", "begin": "预计开始", "due": "预计结束", "created": "创建时间", "completed": "完成时间", "effort": "预估工时"}
	BaseFieldsList    = []string{"ID", "状态", "优先级", "规模", "需求分类", "父需求", "发布计划", "处理人", "开发人员", "创建人", "预计开始", "预计结束", "创建时间", "完成时间", "预估工时"}
	IOSFields         = map[string]string{"custom_field_99": "接口上线日", "custom_field_97": "双端都提得需求", "custom_field_93": "端范围（默认仅粉iPhone）", "custom_field_92": "是否可以单端上线"}
	IOSFieldsList     = []string{"接口上线日", "双端都提得需求", "端范围（默认仅粉iPhone）", "是否可以单端上线"}
	AndroidFields     = map[string]string{"custom_field_99": "接口上线日", "custom_field_97": "双端都提得需求", "custom_field_93": "是否可以单端上线"}
	AndroidFieldsList = []string{"接口上线日", "双端都提得需求", "是否可以单端上线"}
	LiveWorkflow      = []string{"需求规划中（未受理需求）", "需求文档评审", "待排期", "已排期", "开发中", "开发完成待体验", "验收通过待测试", "测试中", "测试通过", "产品设计验收通过", "需求挂起", "需求取消", "免测发布", "已实现"}
	BPlusWorkflow     = []string{"需求中", "已评审", "开发中", "产品/设计体验", "转测试", "测试中", "测试完成", "已实现", "已拒绝"}

	AndroidStoryWallFields   = []string{"auditing", "status_2", "status_11", "developing", "product_experience", "status_5", "testing", "status_6", "status_7", "status_1"}
	AndroidStoryWallColNames = map[string]string{"auditing": "需求评审", "status_2": "确定发布计划", "status_11": "待开发", "developing": "开发中", "product_experience": "产品设计验收", "status_5": "验收通过待测试", "testing": "测试中", "status_6": "测试通过", "status_7": "待合入总包", "status_1": "需求完成"}
	IosStoryWallFields       = []string{"auditing", "status_2", "status_11", "developing", "product_experience", "status_5", "testing", "status_7", "status_8", "status_1"}
	IosStoryWallColNames     = map[string]string{"auditing": "需求评审", "status_2": "确定发布计划", "status_11": "待开发", "developing": "开发中", "product_experience": "产品设计验收", "status_5": "验收通过待测试", "testing": "测试中", "status_7": "测试通过待合入", "status_8": "已合入总包", "status_1": "需求完成"}
)

//test and reject type
var (
	Test       = "test"
	Experience = "experience"
)

//workspace type
var (
	IOS     = "ios"
	Android = "android"
	Live    = "live"
	BPlus   = "bplus"
)
