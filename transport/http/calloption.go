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

// EmptyCallOption does not alter the Call configuration.
// It can be embedded in another structure to carry satellite data for use
// by interceptors.
type EmptyCallOption struct{}

func (EmptyCallOption) before(*callInfo) error      { return nil }
func (EmptyCallOption) after(*callInfo, *csAttempt) {}

type callInfo struct {
	bodyPattern string
	method      string
}

type csAttempt struct{}

// BodyPattern is bodyPattern
func BodyPattern(bodyPattern string) CallOption {
	return BodyPatternCallOption{BodyPattern: bodyPattern}
}

// BodyPatternCallOption is BodyPattern
type BodyPatternCallOption struct {
	BodyPattern string
}

func (o BodyPatternCallOption) before(c *callInfo) error {
	c.bodyPattern = o.BodyPattern
	return nil
}
func (o BodyPatternCallOption) after(c *callInfo, attempt *csAttempt) {}

// Method is Method
func Method(method string) CallOption {
	return MethodCallOption{Method: method}
}

// MethodCallOption is BodyCallOption
type MethodCallOption struct {
	Method string
}

func (o MethodCallOption) before(c *callInfo) error {
	c.method = o.Method
	return nil
}
func (o MethodCallOption) after(c *callInfo, attempt *csAttempt) {}

func defaultCallInfo() callInfo {
	return callInfo{
		bodyPattern: "*",
		method:      "POST",
	}
}
