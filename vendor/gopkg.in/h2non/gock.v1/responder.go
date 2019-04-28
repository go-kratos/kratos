package gock

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Responder builds a mock http.Response based on the given Response mock.
func Responder(req *http.Request, mock *Response, res *http.Response) (*http.Response, error) {
	// If error present, reply it
	err := mock.Error
	if err != nil {
		return nil, err
	}

	if res == nil {
		res = createResponse(req)
	}

	// Apply response filter
	for _, filter := range mock.Filters {
		if !filter(res) {
			return res, nil
		}
	}

	// Define mock status code
	if mock.StatusCode != 0 {
		res.Status = strconv.Itoa(mock.StatusCode) + " " + http.StatusText(mock.StatusCode)
		res.StatusCode = mock.StatusCode
	}

	// Define headers by merging fields
	res.Header = mergeHeaders(res, mock)

	// Define mock body, if present
	if len(mock.BodyBuffer) > 0 {
		res.ContentLength = int64(len(mock.BodyBuffer))
		res.Body = createReadCloser(mock.BodyBuffer)
	}

	// Apply response mappers
	for _, mapper := range mock.Mappers {
		if tres := mapper(res); tres != nil {
			res = tres
		}
	}

	// Sleep to simulate delay, if necessary
	if mock.ResponseDelay > 0 {
		time.Sleep(mock.ResponseDelay)
	}

	return res, err
}

// createResponse creates a new http.Response with default fields.
func createResponse(req *http.Request) *http.Response {
	return &http.Response{
		ProtoMajor: 1,
		ProtoMinor: 1,
		Proto:      "HTTP/1.1",
		Request:    req,
		Header:     make(http.Header),
		Body:       createReadCloser([]byte{}),
	}
}

// mergeHeaders copies the mock headers.
func mergeHeaders(res *http.Response, mres *Response) http.Header {
	for key := range mres.Header {
		res.Header.Set(key, mres.Header.Get(key))
	}
	return res.Header
}

// createReadCloser creates an io.ReadCloser from a byte slice that is suitable for use as an
// http response body.
func createReadCloser(body []byte) io.ReadCloser {
	return &dummyReadCloser{body: bytes.NewReader(body)}
}

// dummyReadCloser is used internally as io.ReadCloser capable interface for bodies.
type dummyReadCloser struct {
	body io.ReadSeeker
}

// Read implements the required method by io.ReadClose interface.
func (d *dummyReadCloser) Read(p []byte) (n int, err error) {
	n, err = d.body.Read(p)
	if err == io.EOF {
		d.body.Seek(0, 0)
	}
	return n, err
}

// Close implements a no-op required method by io.ReadClose interface.
func (d *dummyReadCloser) Close() error {
	return nil
}
