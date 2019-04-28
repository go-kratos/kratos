package param

// ChallListParam describe challenge list search params of a business
type ChallListParam struct {
	Businesses []int8 `form:"businesses" validate:"required,min=1"`
	AssignNum  []int8 `form:"assign_num" validate:"required,min=0"`
	Order      string `form:"order" default:"id"`
	Sort       string `form:"sort" default:"desc"`
	PN         int    `form:"pn"`
	PS         int    `form:"ps"`
	R          int64  `form:"r" validate:"required"`
}

// ChallHandlingDoneListParam describe params challenge list handling of admin
type ChallHandlingDoneListParam struct {
	Businesses int8   `form:"businesses" validate:"required,min=1"`
	Order      string `form:"order"`
	Sort       string `form:"sort"`
	PN         int    `form:"pn"`
	PS         int    `form:"ps"`
}

// ChallCountParam describe challenge count in some states of a business
type ChallCountParam struct {
	Business int64   `form:"business" validate:"required,min=1"`
	States   []int64 `form:"states,split" validate:"dive,gt=-1"`
}

// ChallCreatedListParam return challenge list created by an admin
type ChallCreatedListParam struct {
	Businesses int8   `form:"businesses" validate:"required,min=1"`
	Order      string `form:"order"`
	Sort       string `form:"sort"`
	PN         int    `form:"pn"`
	PS         int    `form:"ps"`
}
