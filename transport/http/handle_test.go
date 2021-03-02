package http

import (
	"context"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

type HelloRequest struct {
	Name string `json:"name"`
}
type HelloReply struct {
	Message string `json:"message"`
}
type GreeterService struct {
}

func (s *GreeterService) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "hello " + req.Name}, nil
}

func newGreeterHandler(srv *GreeterService, opts ...HandleOption) http.Handler {
	h := DefaultHandleOptions()
	for _, o := range opts {
		o(&h)
	}
	r := mux.NewRouter()
	r.HandleFunc("/helloworld", func(w http.ResponseWriter, req *http.Request) {
		var in HelloRequest
		if err := h.Decode(req, &in); err != nil {
			h.Error(w, err)
			return
		}
		next := func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SayHello(ctx, &in)
		}
		if h.Middleware != nil {
			next = h.Middleware(next)
		}
		out, err := next(req.Context(), &in)
		if err != nil {
			h.Error(w, err)
			return
		}
		if err := h.Encode(w, out); err != nil {
			h.Error(w, err)
		}
	}).Methods("POST")
	return r
}

func TestHandler(t *testing.T) {
	s := &GreeterService{}
	_ = newGreeterHandler(s)
}
