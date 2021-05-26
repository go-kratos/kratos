package http

// CallOption configures a Call before it starts or extracts information from
// a Call after it completes.
type CallOption interface {
	// before is called before the call is sent to any server.  If before
	// returns a non-nil error, the RPC fails with that error.
	before(*callInfo) error

	// after is called after the call has completed.  after cannot return an
	// error, so any failures should be reported via output parameters.
	after(*callInfo, *csAttempt)
}

type callInfo struct {
	pathPattern string
	method      string
}

// EmptyCallOption does not alter the Call configuration.
// It can be embedded in another structure to carry satellite data for use
// by interceptors.
type EmptyCallOption struct{}

func (EmptyCallOption) before(*callInfo) error      { return nil }
func (EmptyCallOption) after(*callInfo, *csAttempt) {}

type csAttempt struct{}

// PathPattern is pathpattern
func PathPattern(pathPattern string) CallOption {
	return PathPatternCallOption{PathPattern: pathPattern}
}

// PathPatternCallOption is BodyPattern
type PathPatternCallOption struct {
	EmptyCallOption
	PathPattern string
}

func (o PathPatternCallOption) before(c *callInfo) error {
	c.pathPattern = o.PathPattern
	return nil
}

// Method is Method
func Method(method string) CallOption {
	return MethodCallOption{Method: method}
}

// MethodCallOption is BodyCallOption
type MethodCallOption struct {
	EmptyCallOption
	Method string
}

func (o MethodCallOption) before(c *callInfo) error {
	c.method = o.Method
	return nil
}

func defaultCallInfo() callInfo {
	return callInfo{
		method: "POST",
	}
}
