package http

import "net/http"

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
	contentType  string
	operation    string
	pathTemplate string
}

// EmptyCallOption does not alter the Call configuration.
// It can be embedded in another structure to carry satellite data for use
// by interceptors.
type EmptyCallOption struct{}

func (EmptyCallOption) before(*callInfo) error      { return nil }
func (EmptyCallOption) after(*callInfo, *csAttempt) {}

type csAttempt struct {
	res *http.Response
}

// ContentType with request content type.
func ContentType(contentType string) CallOption {
	return ContentTypeCallOption{ContentType: contentType}
}

// ContentTypeCallOption is BodyCallOption
type ContentTypeCallOption struct {
	EmptyCallOption
	ContentType string
}

func (o ContentTypeCallOption) before(c *callInfo) error {
	c.contentType = o.ContentType
	return nil
}

func defaultCallInfo(path string) callInfo {
	return callInfo{
		contentType:  "application/json",
		operation:    path,
		pathTemplate: path,
	}
}

// Operation is serviceMethod call option
func Operation(operation string) CallOption {
	return OperationCallOption{Operation: operation}
}

// OperationCallOption is set ServiceMethod for client call
type OperationCallOption struct {
	EmptyCallOption
	Operation string
}

func (o OperationCallOption) before(c *callInfo) error {
	c.operation = o.Operation
	return nil
}

// PathTemplate is http path template
func PathTemplate(pattern string) CallOption {
	return PathTemplateCallOption{Pattern: pattern}
}

// PathTemplateCallOption is set path template for client call
type PathTemplateCallOption struct {
	EmptyCallOption
	Pattern string
}

func (o PathTemplateCallOption) before(c *callInfo) error {
	c.pathTemplate = o.Pattern
	return nil
}

// Header returns a CallOptions that retrieves the http response header
// from server reply.
func Header(header *http.Header) CallOption {
	return HeaderCallOption{header: header}
}

// HeaderCallOption is retrive response header for client call
type HeaderCallOption struct {
	EmptyCallOption
	header *http.Header
}

func (o HeaderCallOption) after(c *callInfo, cs *csAttempt) {
	if cs.res != nil && cs.res.Header != nil {
		*o.header = cs.res.Header
	}
}
