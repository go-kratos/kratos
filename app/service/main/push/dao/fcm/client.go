package fcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// PriorityHigh used for high notification priority
	PriorityHigh = "high"

	// PriorityNormal used for normal notification priority
	PriorityNormal = "normal"

	// HeaderRetryAfter HTTP header constant
	HeaderRetryAfter = "Retry-After"

	// ErrorKey readable error caching
	ErrorKey = "error"

	// MethodPOST indicates http post method
	MethodPOST = "POST"

	// ServerURL push server url
	ServerURL = "https://fcm.googleapis.com/fcm/send"
)

// retryableErrors whether the error is a retryable
var retryableErrors = map[string]bool{
	"Unavailable":         true,
	"InternalServerError": true,
}

// Client stores client with api key to firebase
type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new client
func NewClient(apiKey string, timeout time.Duration) *Client {
	return &Client{
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: timeout},
	}
}

func (f *Client) authorization() string {
	return fmt.Sprintf("key=%v", f.APIKey)
}

// Send sends message to FCM
func (f *Client) Send(message *Message) (*Response, error) {
	data, err := json.Marshal(message)
	if err != nil {
		return &Response{}, err
	}
	req, err := http.NewRequest(MethodPOST, ServerURL, bytes.NewBuffer(data))
	if err != nil {
		return &Response{}, err
	}
	req.Header.Set("Authorization", f.authorization())
	req.Header.Set("Content-Type", "application/json")
	resp, err := f.HTTPClient.Do(req)
	if err != nil {
		return &Response{}, err
	}
	defer resp.Body.Close()
	response := &Response{StatusCode: resp.StatusCode}
	if resp.StatusCode >= 500 {
		response.RetryAfter = resp.Header.Get(HeaderRetryAfter)
	}
	if resp.StatusCode != 200 {
		return response, fmt.Errorf("fcm status code(%d)", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return response, err
	}
	if err := f.Failed(response); err != nil {
		return response, err
	}
	response.Ok = true
	return response, nil
}

// Failed method indicates if the server couldn't process
// the request in time.
func (f *Client) Failed(response *Response) error {
	for _, response := range response.Results {
		if retryableErrors[response.Error] {
			return fmt.Errorf("fcm push error(%s)", response.Error)
		}
	}
	return nil
}
