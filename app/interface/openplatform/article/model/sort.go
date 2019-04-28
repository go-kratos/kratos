package model

const (
	// FieldDefault with recommends
	FieldDefault = iota
	// FieldNew new list
	FieldNew
	// FieldLike like stat
	FieldLike
	// FieldReply reply stat
	FieldReply
	// FieldFav favorite stat
	FieldFav
	// FieldView view stat
	FieldView
)

// SortFields all field
var SortFields = []int{FieldNew, FieldLike, FieldReply, FieldFav, FieldView}
