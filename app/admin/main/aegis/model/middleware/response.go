package middleware

//IMiddleware handler
type IMiddleware interface {
	Process(data interface{})
}

//ResponseRender .
type ResponseRender func(data interface{}, err error)

//Response response handler
func Response(data interface{}, err error, r ResponseRender, i IMiddleware) {
	if data != nil && i != nil {
		i.Process(data)
	}
	if r != nil {
		r(data, err)
	}
}

//Request request handler
func Request(data interface{}, i IMiddleware) {
	if data != nil && i != nil {
		i.Process(data)
	}
}
