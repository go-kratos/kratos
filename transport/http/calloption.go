package http

import "github.com/go-kratos/kratos/v2/metadata"

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
	contentType string
	operation   string
	metatada    metadata.Metadata
}

// EmptyCallOption does not alter the Call configuration.
// It can be embedded in another structure to carry satellite data for use
// by interceptors.
type EmptyCallOption struct{}

func (EmptyCallOption) before(*callInfo) error      { return nil }
func (EmptyCallOption) after(*callInfo, *csAttempt) {}

type csAttempt struct{}

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
		contentType: "application/json",
		operation:   path,
		metatada:    metadata.Metadata{},
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

// Metadata is Metadata call option
func Metadata(metatada metadata.Metadata) CallOption {
	return MetadataCallOption{Metatada: metatada}
}

// MetadataCallOption is set Metadata for client call
type MetadataCallOption struct {
	EmptyCallOption
	Metatada metadata.Metadata
}

func (o MetadataCallOption) before(c *callInfo) error {
	c.metatada = o.Metatada
	return nil
}
