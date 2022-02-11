package permission

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
)

// Builder . The best way to use this Builder is build a struct
// to implement this interface and use it through wire dependency inject
// so that you can do lots of things for check permission
type Builder interface {
	// build . Implement this function to filter path
	build() selector.MatchFunc
}

// Permission . New permission Middleware
func Permission(builder Builder) middleware.Middleware {
	return selector.Server(
		server(),
	).Match(
		builder.build(),
	).Build()
}
