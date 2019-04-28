package model

//Danmu ...
type Danmu struct {
	ID      int64
	OID     int64
	MID     int64
	Offset  int64
	Content string
	State   int8
}
