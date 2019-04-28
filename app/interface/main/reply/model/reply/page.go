package reply

const (
	// AttrNo attribute no
	AttrNo = uint32(0)
	// AttrYes attribute yes
	AttrYes = uint32(1)
)

// PageParams reply page param.
type PageParams struct {
	Mid  int64
	Oid  int64
	Type int8
	Sort int8

	PageNum  int
	PageSize int

	NeedHot    bool
	NeedSecond bool
	Escape     bool
}

// PageResult reply page result.
type PageResult struct {
	Subject  *Subject
	TopAdmin *Reply
	TopUpper *Reply
	Roots    []*Reply
	Hots     []*Reply
	Total    int
	AllCount int
}
