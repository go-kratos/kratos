package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go-common/library/log"
)

// HTTPServer provides an HTTP api for an agent.
type HTTPServer struct {
	*http.Server
	agent *Agent

	// proto is filled by the agent to "http" or "https".
	proto string
}

// NewHTTPServer http server provide simple query api
func NewHTTPServer(addr string, a *Agent) *HTTPServer {
	s := &HTTPServer{
		Server: &http.Server{Addr: addr},
		agent:  a,
	}
	s.Server.Handler = s.handler()
	return s
}

// handler is used to attach our handlers to the mux
func (s *HTTPServer) handler() http.Handler {
	mux := http.NewServeMux()
	// TODO: simple manage ui

	// API V1
	mux.HandleFunc("/v1/naming/", s.wrap(s.NSTranslation)) // naming path

	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

// wrap is used to wrap functions to make them more convenient
func (s *HTTPServer) wrap(handler func(resp http.ResponseWriter, req *http.Request) (interface{}, error)) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		logURL := req.URL.String()
		handleErr := func(err error) {
			log.Error("http: Request %s %v from %s, err: %v\n", req.Method, logURL, req.RemoteAddr, err)
			resp.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(resp, err.Error())
		}

		// Invoke the handler
		start := time.Now()
		defer func() {
			log.V(5).Info("http: Request %s %v from %s, Timing: %v\n", req.Method, logURL, req.RemoteAddr, time.Since(start))
		}()
		obj, err := handler(resp, req)
		if err != nil {
			handleErr(err)
			return
		}
		if obj == nil {
			return
		}

		buf, err := s.marshalJSON(req, obj)
		if err != nil {
			handleErr(err)
			return
		}
		resp.Header().Set("Content-Type", "application/json")
		resp.Write(buf)
	}
}

// marshalJSON marshals the object into JSON, respecting the user's pretty-ness
// configuration.
func (s *HTTPServer) marshalJSON(req *http.Request, obj interface{}) ([]byte, error) {
	if _, ok := req.URL.Query()["pretty"]; ok {
		buf, err := json.MarshalIndent(obj, "", "    ")
		if err != nil {
			return nil, err
		}
		buf = append(buf, "\n"...)
		return buf, nil
	}

	buf, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return buf, err
}
