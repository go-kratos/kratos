package model

// 推送消息类型
const (
	NoticeTypeLike    = 1
	NoticeTypeComment = 2
	NoticeTypeFan     = 3
	NoticeTypeSysMsg  = 4
)

// 推送业务类型
const (
	NoticeBizTypeSv        = 1
	NoticeBizTypeComment   = 2
	NoticeBizTypeUser      = 3
	NoticeBizTypeCmsReview = 4
)

// 推送文案
const (
	PushMsgVideoLike    = "%s%s赞了你的作品"
	PushMsgVideoComment = "%s%s评论了你的作品"
	PushMsgCommentLike  = "%s%s赞了你的评论"
	PushMsgCommentReply = "%s%s回复了你的评论"
	PushMsgFollow       = "%s%s关注了你"
)

// 推送跳转schema
const (
	PushSchemaVideo   = "qing://videoplayer?svid=%d"
	PushSchemaComment = "qing://commentdetail?svid=%d&rootid=%s"
	PushSchemaUser    = "qing://profile?mid=%d"
	PushSchemaNotice  = "qing://notification?type=%d"
)
