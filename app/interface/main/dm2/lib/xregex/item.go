package xregex

/*
golang version regex parser
refer to: https://github.com/aristotle9/as3cc/tree/master/java-template/src/org/lala/lex/utils/parser
*/

type rangeItem struct {
	from  int64
	to    int64
	value int64
}

type stateTransItem struct {
	isDead    bool
	toStates  []int64
	transEdge []*rangeItem
}

type productionItem struct {
	headerID   int64
	bodyLength int64
}
