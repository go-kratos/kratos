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
	pathPattern string
	bodyPattern *string
}

type csAttempt struct{}

func Path(pathPattern string) CallOption {
	return PathCallOption{PathPattern: pathPattern}
}

type PathCallOption struct {
	PathPattern string
}

func (o PathCallOption) before(c *callInfo) error {
	c.pathPattern = o.PathPattern
	return nil
}
func (o PathCallOption) after(c *callInfo, attempt *csAttempt) {}

func Body(bodyPattern string) CallOption {
	return BodyCallOption{BodyPattern: &bodyPattern}
}

type BodyCallOption struct {
	BodyPattern *string
}

func (o BodyCallOption) before(c *callInfo) error {
	c.bodyPattern = o.BodyPattern
	return nil
}
func (o BodyCallOption) after(c *callInfo, attempt *csAttempt) {}
