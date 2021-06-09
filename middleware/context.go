package middleware

import "context"

type serviceMethodKey struct{}

// WithServiceMethod with service full method, i.e. /package.service/method.
func WithServiceMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, serviceMethodKey{}, method)
}

// ServiceMethod returns the service full method, i.e. /package.service/method.
func ServiceMethod(ctx context.Context) string {
	method, _ := ctx.Value(serviceMethodKey{}).(string)
	return method
}
