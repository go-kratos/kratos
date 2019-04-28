package like

// ArgMatch arg match
type ArgMatch struct {
	Sid int64
}

// ArgSubjectUp .
type ArgSubjectUp struct {
	Sid int64
}

// ArgLikeUp .
type ArgLikeUp struct {
	Lid int64
}

// ArgLikeItem .
type ArgLikeItem struct {
	ID   int64
	Sid  int64
	Type int
}

// ArgActSubject .
type ArgActSubject struct {
	Sid int64
}

// ArgActProtocol .
type ArgActProtocol struct {
	Sid int64
}
