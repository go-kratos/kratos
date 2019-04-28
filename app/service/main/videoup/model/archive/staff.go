package archive

// staff state .
const (
	APPLYSTATEOPEN   = int8(0)
	APPLYSTATEACCEPT = int8(1)
	APPLYSTATEREFUSE = int8(2)
	//场景是 staff未上线前 up直接删除
	APPLYSTATEDEL    = int8(3)
	APPLYSTATEIGNORE = int8(4)

	DEALSTATEOPEN   = int8(1)
	DEALSTATEDONE   = int8(2)
	DEALSTATEIGNORE = int8(3)

	STATEON  = int8(1)
	STATEOFF = int8(2)

	TYPEUPADD    = int8(1)
	TYPEUPDEL    = int8(2)
	TYPEUPMODIFY = int8(3)
	TYPESTAFFDEL = int8(4)
	TYPEADMINDEL = int8(5)

	STAFFLogBizID = int(84)

	STAFFLogBizType1 = int(1)
	STAFFLogBizType2 = int(2)
	STAFFLogBizType3 = int(3)
	STAFFLogBizType4 = int(4)

	UPRELATIONBLACK = int64(128)
)

//Staff . 正式staff
type Staff struct {
	ID           int64  `json:"id"`
	AID          int64  `json:"aid"`
	MID          int64  `json:"mid"`
	StaffMID     int64  `json:"staff_mid"`
	StaffTitle   string `json:"staff_title"`
	StaffName    string `json:"staff_name"`
	StaffTitleID int64  `json:"staff_title_id"`
	State        int8   `json:"state"`
}

//StaffParam 提交的staff参数
type StaffParam struct {
	//apply_id 建议前端传 为后面预留
	ApplyID int64  `json:"apply_id"`
	Title   string `json:"title"`
	MID     int64  `json:"mid"`
	TitleID int64  `json:"title_id"`
}

//StaffBatchParam  批量提交的staff参数
type StaffBatchParam struct {
	AID      int64         `json:"aid"`
	SyncAttr bool          `json:"sync_attr"`
	Staffs   []*StaffParam `json:"staffs"`
}

//ApplyParam 提交申请单参数
type ApplyParam struct {
	ID            int64    `form:"id"`
	Type          int8     `form:"type"`
	ASID          int64    `form:"as_id"`
	ApplyAID      int64    `form:"apply_aid"`
	ApplyStaffMID int64    `form:"apply_staff_mid" validate:"required"`
	ApplyUpMID    int64    `form:"apply_up_mid"`
	ApplyTitle    string   `form:"apply_title"`
	OldTitle      string   `form:"old_title"`
	ApplyTitleID  int64    `form:"apply_title_id"`
	State         int8     `form:"state"`
	DealState     int8     `form:"deal_state"`
	RefuseMid     int64    `form:"refuse_mid"`
	FlagRefuse    bool     `form:"flag_refuse"`
	FlagAddBlack  bool     `form:"flag_add_black"`
	NoNotify      bool     `form:"no_notify"`
	SyncStaff     bool     `form:"sync_staff"`
	CleanCache    bool     `form:"clean_cache"`
	SyncDynamic   bool     `form:"sync_dynamic"`
	MsgId         int      `form:"msg_id"`
	StaffState    int8     `json:"staff_state"`
	StaffTitle    string   `json:"staff_title"`
	Archive       *Archive `json:"archive"`
	UpName        string   `json:"up_name"`
	StaffName     string   `json:"staff_name"`
	StaffsName    string   `json:"staffs_name"`
	NotifyUp      bool     `json:"notify_up"`
}

type SearchApplyIndex struct {
	Indexs []*Index `json:"creative_archive"`
}

type Index struct {
	ID   int64        `json:"id"`
	Item []*IndexItem `json:"apply_staff"`
}

type IndexItem struct {
	DealState     int8  `json:"deal_state"`
	ApplyStaffMID int64 `json:"apply_staff_mid"`
}

//StaffApply 申请单
type StaffApply struct {
	ID            int64  `json:"id"`
	Type          int8   `json:"apply_type"`
	ASID          int64  `json:"apply_as_id"`
	ApplyAID      int64  `json:"apply_aid"`
	ApplyUpMID    int64  `json:"apply_up_mid"`
	ApplyStaffMID int64  `json:"apply_staff_mid"`
	ApplyTitle    string `json:"apply_title"`
	ApplyTitleID  int64  `json:"apply_title_id"`
	State         int8   `json:"apply_state"`
	DealState     int8   `json:"deal_state"`
	StaffState    int8   `json:"staff_state"`
	StaffTitle    string `json:"staff_title"`
}

//Copy . apply转化成staff
func (s *Staff) Copy(v *ApplyParam) {
	s.AID = v.ApplyAID
	s.MID = v.ApplyUpMID
	s.StaffMID = v.ApplyStaffMID
	s.ID = v.ASID
	s.StaffTitle = v.ApplyTitle
	s.StaffTitleID = v.ApplyTitleID
	switch v.State {
	case APPLYSTATEACCEPT:
		switch v.Type {
		case TYPEUPADD:
			s.State = STATEON
		case TYPEUPMODIFY:
			s.State = STATEON
		case TYPEUPDEL, TYPEADMINDEL, TYPESTAFFDEL:
			s.State = STATEOFF
		}
	case APPLYSTATEREFUSE:
		switch v.Type {
		case TYPEUPADD:
			s.State = STATEOFF
		case TYPEUPMODIFY:
			s.State = STATEON
		case TYPEUPDEL, TYPEADMINDEL, TYPESTAFFDEL:
			s.State = STATEON
		}
	case APPLYSTATEDEL:
		switch v.Type {
		case TYPEADMINDEL, TYPESTAFFDEL:
			s.State = STATEOFF
		}
	default:
		s.State = STATEOFF
	}
}

//Copy . 稿件编辑时用
func (s *ApplyParam) Copy(v *StaffApply) {
	s.ApplyAID = v.ApplyAID
	s.ApplyStaffMID = v.ApplyStaffMID
	s.ASID = v.ASID
	s.ApplyTitle = v.ApplyTitle
	s.ApplyTitleID = v.ApplyTitleID
	s.State = v.State
	s.Type = v.Type
}
