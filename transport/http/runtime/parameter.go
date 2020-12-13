package runtime

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/hashicorp/vic/pkg/retry"
)

// ParameterView provides the access to the path parameters
// and query parameters.
// It could be returned by a ParameterParser.
type ParameterView interface {
	// Get returns the first value associated with the given key
	// It returns an empty value if there are no values associated
	// with the key.
	Get(key string) string

	// Get returns the first value associated with the given key
	// It returns ErrInvalidArgument if there are no values associated
	// with the key.
	Require(string) (string, error)
}

// QueryParameterParser is used to parse the query parameters to
// ParameterView.
type QueryParameterParser interface {
	Parse(context.Context, *http.Request) (ParameterView, error)
}

type queryParameterView url.Values

func (v queryParameterView) Get(key string) string {
	if v == nil {
		return ""
	}
	vs := v[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

func (v queryParameterView) Require(key string) (string, error) {
	if v == nil {
		return "", ErrInvalidArgument
	}

	vs := v[key]
	if len(vs) == 0 {
		return "", ErrInvalidArgument
	}

	return vs[0], nil
}

// CompositedParameterView is used to composite the multiple ParameterViews
// into one ParameterView.
// It will resolve the value from the first to the last.
type CompositedParameterView struct {
	views []ParameterView
}

func NewCompositedParameterView(vs ...ParameterView) *CompositedParameterView {
	return &CompositedParameterView{
		views: vs,
	}
}

func (v *CompositedParameterView) Get(key string) string {
	for _, view := range v.views {
		val := view.Get(key)
		if val != "" {
			return val
		}
	}
	return ""
}

func (v *CompositedParameterView) Require(key string) (string, error) {
	for _, view := range v.views {
		val, err := view.Require(key)
		if err == nil {
			return val, nil
		}
	}
	return "", ErrInvalidArgument
}

func QueryParameterHandle(next Handler, w http.ResponseWriter, r *http.Request, v ParameterView) error {
	if next != nil {
		qpv := queryParameterView(r.URL.Query())
		if v != nil {
			return next.ServeHTTP(w, r, NewCompositedParameterView(qpv, v))
		}
		return next.ServeHTTP(w, r, qpv)
	}
	return errors.New("No next handler")
}
