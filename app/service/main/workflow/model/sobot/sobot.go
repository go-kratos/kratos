package sobot

const (
	// TicketFrom .
	TicketFrom = int8(12)
	// EcodeOK .
	EcodeOK = "000000"

	// ReplyTypePublic 评论可见状态
	ReplyTypePublic = int8(0)
	// ReplyTypeCSOnly .
	ReplyTypeCSOnly = int8(1)

	// CustomerSourcePC 客户来源
	CustomerSourcePC = int8(0)
	// CustomerSourceWX .
	CustomerSourceWX = int8(1)
	// CustomerSourceAPP .
	CustomerSourceAPP = int8(2)
	// CustomerSourceWB .
	CustomerSourceWB = int8(3)
	// CustomerSourceWAP .
	CustomerSourceWAP = int8(4)

	// TicketLevelLow 工单等级 .
	TicketLevelLow = int8(0)
	// TicketLevelMedium .
	TicketLevelMedium = int8(1)
	// TicketLevelHigh .
	TicketLevelHigh = int8(2)
	// TicketLevelurgency .
	TicketLevelurgency = int8(3)

	// TicketStatusPending 工单状态
	TicketStatusPending = int8(0)
	// TicketStatusHandling .
	TicketStatusHandling = int8(1)
	// TicketStatusReplying .
	TicketStatusReplying = int8(2)
	// TicketStatusSolved .
	TicketStatusSolved = int8(3)
	// TicketStatusClosed .
	TicketStatusClosed = int8(99)
	// TicketStatusDeleted .
	TicketStatusDeleted = int8(98)
)

// Ticket struct
type Ticket struct {
	TicketID string `json:"ticket_id"`
	Content  string `json:"ticket_content"`
	Level    int8   `json:"ticket_level"`
	State    int8   `json:"ticket_status"`
	Title    string `json:"ticket_title"`
	FileStr  string `json:"file_str"`
	CTime    int64  `json:"ctime"`
}

// Reply struct
type Reply struct {
	Face      string `json:"face_img"`
	FileStr   string `json:"file_str"`
	Content   string `json:"reply_content"`
	ReplyType int8   `json:"reply_type"`
	ShowName  string `json:"show_name"`
	StartType int8   `json:"start_type"`
	CTime     int64  `json:"reply_time"`
}

// ReplyParam reply param
type ReplyParam struct {
	TicketID      int32  `form:"ticket_id" validate:"required"`
	ReplyContent  string `form:"reply_content" validate:"required"`
	CustomerEmail string `form:"customer_email" validate:"required"`
	StartType     int8   `form:"start_type"`
	ReplyType     int8   `form:"reply_type"`
}

// TicketParam ticket param
type TicketParam struct {
	CustomerName   string `form:"customer_name"`
	CustomerQQ     string `form:"customer_qq"`
	CustomerNick   string `form:"customer_nick"`
	CustomerEmail  string `form:"customer_email" validate:"required"`
	CustomerSource int8   `form:"customer_source"`
	CustomerPhone  string `form:"customer_phone"`
	TicketID       int32  `form:"ticket_id" validate:"required"`
	TicketTitle    string `form:"ticket_title"`
	TicketContent  string `form:"ticket_content"`
	TicketLevel    int8   `form:"ticket_level"`
	TicketStatus   int8   `form:"ticket_status"`
	TicketFrom     int8   `form:"ticket_from"`
	StartType      int8   `form:"start_type"`
	FileStr        string `form:"file_str"`
}
