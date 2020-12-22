package log

type filter struct {
	Logger
}

func newFilter() Logger {
	return &filter{}
}

func (n *filter) Print(kvpair ...interface{}) {
	n.Logger.Print(kvpair...)
}
