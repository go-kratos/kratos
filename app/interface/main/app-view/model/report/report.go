package report

const (
	ArchivePVR       = 0 // 稿件存在色情、暴力、反动
	ArchiveCopyWrite = 1 // 稿件存在版权内容
	ArchiveCrash     = 2 // 与站内其他视频撞车
	ArchiveNotOwn    = 3 // 不能参与充电计划：非自制作品
	ArchiveBusiness  = 4 // 不能参与充电计划：有商业推广内容
	Other            = 5 // 其他
	// oversea
	ArchivePVRG       = 6
	ArchiveCopyWriteG = 7
	ArchiveHarmG      = 8
	OtherG            = 9

	Reason1 = "存在色情、暴力、反动内容"
	Reason2 = "存在版权问题"
	Reason3 = "与站内其他视频撞车"
	Reason4 = "不能参与充电计划：非自制作品"
	Reason5 = "不能参与充电计划：有商业推广内容"
	Reason6 = "其他问题"
	//
	Reason7  = "存在色情、暴力血腥等內容"
	Reason8  = "存在版權問題"
	Reason9  = "存在有害或危險行為內容"
	Reason10 = "其他問題"
	Reason11 = "存在有害或危害行为内容"

	Desc1 = "为帮助审核人员更快处理, 请补充 违规内容出现位置等详细信息"
	Desc2 = "如果您认为这个稿件作品侵犯了您的相关权益， 请登录哔哩哔哩网页端 ，找到页面底部的「侵权投诉」入口，下载页面中的「侵权申诉表」，按照提示录入相关信息后以邮件的形式反馈我们。"
	Desc3 = "为帮助审核人员更快处理, 请补充重复稿件av号等详细信息"
	Desc4 = "为帮助审核人员更快处理, 请补充转载来源等详细信息"
	Desc5 = "为帮助审核人员更快处理, 请补充商业元素出现位置等详细信息"
	Desc6 = "为帮助审核人员更快处理, 请补充问题类型和出现位置等详细信息"
	//
	Desc7  = "為幫助稽核人員更快處理, 請補充違規內容出現位置等詳細資訊"
	Desc8  = "如果您認為這個稿件作品侵犯了您的相關權益， 請登入嗶哩嗶哩網頁端 ，找到頁面底部的「侵權投訴」入口，下載頁面中的「侵權申訴表」，按照提示錄入相關資訊後以郵件的形式反饋我們。"
	Desc9  = "為幫助稽核人員更快處理, 請補充違規內容出現位置等詳細資訊"
	Desc10 = "為幫助稽核人員更快處理, 請補充問題型別和出現位置等詳細資訊"
	Desc11 = "为帮助审核人员更快处理, 请补充违规内容出现位置等详细信息"
)

type CopyWriter struct {
	Typ      int    `json:"type"`
	Reason   string `json:"reason"`
	Desc     string `json:"desc"`
	AllowAdd bool   `json:"allow_add"`
}
