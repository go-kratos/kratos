package middleware

import "context"

type methodKey struct{}

// WithMethod with service full method, i.e. /package.service/method.
func WithMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, methodKey{}, method)
}

// Method returns the method string for the server context.
// The returned string is in the format of "/package.service/method".
func Method(ctx context.Context) (string, bool) {
	method, ok := ctx.Value(methodKey{}).(string)
	return method, ok
}
