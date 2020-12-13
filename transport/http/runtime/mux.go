package runtime

import "net/http"

type ErrorHandler interface {
}

// ServerMux represents a HTTP RESTful multiplexer.
// It implements the http.Handler, so it can be handled by the
// http.ServerMux.
type ServerMux struct {
	handler Handler
}

func (m *ServerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := m.handler.ServeHTTP(w, r, nil)
	if err != nil {
		// error handling
	}
}
