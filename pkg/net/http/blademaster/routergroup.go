package blademaster

import (
	"regexp"
)

// IRouter http router framework interface.
type IRouter interface {
	IRoutes
	Group(string, ...HandlerFunc) *RouterGroup
}

// IRoutes http router interface.
type IRoutes interface {
	UseFunc(...HandlerFunc) IRoutes
	Use(...Handler) IRoutes

	Handle(string, string, ...HandlerFunc) IRoutes
	HEAD(string, ...HandlerFunc) IRoutes
	GET(string, ...HandlerFunc) IRoutes
	POST(string, ...HandlerFunc) IRoutes
	PUT(string, ...HandlerFunc) IRoutes
	DELETE(string, ...HandlerFunc) IRoutes
}

// RouterGroup is used internally to configure router, a RouterGroup is associated with a prefix
// and an array of handlers (middleware).
type RouterGroup struct {
	Handlers   []HandlerFunc
	basePath   string
	engine     *Engine
	root       bool
	baseConfig *MethodConfig
}

var _ IRouter = &RouterGroup{}

// Use adds middleware to the group, see example code in doc.
func (group *RouterGroup) Use(middleware ...Handler) IRoutes {
	for _, m := range middleware {
		group.Handlers = append(group.Handlers, m.ServeHTTP)
	}
	return group.returnObj()
}

// UseFunc adds middleware to the group, see example code in doc.
func (group *RouterGroup) UseFunc(middleware ...HandlerFunc) IRoutes {
	group.Handlers = append(group.Handlers, middleware...)
	return group.returnObj()
}

// Group creates a new router group. You should add all the routes that have common middlwares or the same path prefix.
// For example, all the routes that use a common middlware for authorization could be grouped.
func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers: group.combineHandlers(handlers),
		basePath: group.calculateAbsolutePath(relativePath),
		engine:   group.engine,
		root:     false,
	}
}

// SetMethodConfig is used to set config on specified method
func (group *RouterGroup) SetMethodConfig(config *MethodConfig) *RouterGroup {
	group.baseConfig = config
	return group
}

// BasePath router group base path.
func (group *RouterGroup) BasePath() string {
	return group.basePath
}

func (group *RouterGroup) handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
	absolutePath := group.calculateAbsolutePath(relativePath)
	injections := group.injections(relativePath)
	handlers = group.combineHandlers(injections, handlers)
	group.engine.addRoute(httpMethod, absolutePath, handlers...)
	if group.baseConfig != nil {
		group.engine.SetMethodConfig(absolutePath, group.baseConfig)
	}
	return group.returnObj()
}

// Handle registers a new request handle and middleware with the given path and method.
// The last handler should be the real handler, the other ones should be middleware that can and should be shared among different routes.
// See the example code in doc.
//
// For HEAD, GET, POST, PUT, and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (group *RouterGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
	if matches, err := regexp.MatchString("^[A-Z]+$", httpMethod); !matches || err != nil {
		panic("http method " + httpMethod + " is not valid")
	}
	return group.handle(httpMethod, relativePath, handlers...)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func (group *RouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle("HEAD", relativePath, handlers...)
}

// GET is a shortcut for router.Handle("GET", path, handle).
func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle("GET", relativePath, handlers...)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (group *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle("POST", relativePath, handlers...)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func (group *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle("PUT", relativePath, handlers...)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func (group *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle("DELETE", relativePath, handlers...)
}

func (group *RouterGroup) combineHandlers(handlerGroups ...[]HandlerFunc) []HandlerFunc {
	finalSize := len(group.Handlers)
	for _, handlers := range handlerGroups {
		finalSize += len(handlers)
	}
	if finalSize >= int(_abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make([]HandlerFunc, finalSize)
	copy(mergedHandlers, group.Handlers)
	position := len(group.Handlers)
	for _, handlers := range handlerGroups {
		copy(mergedHandlers[position:], handlers)
		position += len(handlers)
	}
	return mergedHandlers
}

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(group.basePath, relativePath)
}

func (group *RouterGroup) returnObj() IRoutes {
	if group.root {
		return group.engine
	}
	return group
}

// injections is
func (group *RouterGroup) injections(relativePath string) []HandlerFunc {
	absPath := group.calculateAbsolutePath(relativePath)
	for _, injection := range group.engine.injections {
		if !injection.pattern.MatchString(absPath) {
			continue
		}
		return injection.handlers
	}
	return nil
}
