package oauth2

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	// maxResponseSize limits the introspection response body size.
	maxResponseSize = 1 << 20

	// tokenTypeHintAccessToken is the RFC 7662 token_type_hint for access tokens.
	tokenTypeHintAccessToken = "access_token"
)

// IntrospectionOption is token introspection validator option.
type IntrospectionOption func(*IntrospectionValidator)

// WithBasicAuth authenticates the resource server to the introspection
// endpoint with HTTP Basic authentication (client_secret_basic).
func WithBasicAuth(clientID, clientSecret string) IntrospectionOption {
	return func(v *IntrospectionValidator) {
		v.clientID = clientID
		v.clientSecret = clientSecret
	}
}

// WithCacheTTL caches active introspection results for at most ttl,
// and never beyond the token expiry. Zero disables the cache.
func WithCacheTTL(ttl time.Duration) IntrospectionOption {
	return func(v *IntrospectionValidator) {
		v.cacheTTL = ttl
	}
}

// WithCacheSize limits the number of cached tokens, default is 10000.
func WithCacheSize(size int) IntrospectionOption {
	return func(v *IntrospectionValidator) {
		v.cacheSize = size
	}
}

// IntrospectionValidator validates tokens with an RFC 7662 token
// introspection endpoint.
type IntrospectionValidator struct {
	endpoint     string
	client       *http.Client
	clientID     string
	clientSecret string
	cacheTTL     time.Duration
	cacheSize    int
	now          func() time.Time

	mu    sync.Mutex
	cache map[[sha256.Size]byte]cacheEntry
}

type cacheEntry struct {
	info     *TokenInfo
	expireAt time.Time
}

var _ Validator = (*IntrospectionValidator)(nil)

// NewIntrospectionValidator returns a validator that calls the introspection
// endpoint. If client is nil, http.DefaultClient is used.
//
// To authenticate with an OAuth2 access token instead of basic auth, pass a
// client created by golang.org/x/oauth2/clientcredentials.Config.Client.
func NewIntrospectionValidator(endpoint string, client *http.Client, opts ...IntrospectionOption) *IntrospectionValidator {
	if client == nil {
		client = http.DefaultClient
	}
	v := &IntrospectionValidator{
		endpoint:  endpoint,
		client:    client,
		cacheSize: 10000,
		now:       time.Now,
		cache:     make(map[[sha256.Size]byte]cacheEntry),
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

// Validate implements Validator.
func (v *IntrospectionValidator) Validate(ctx context.Context, token string) (*TokenInfo, error) {
	key := sha256.Sum256([]byte(token))
	if info, ok := v.load(key); ok {
		return info, nil
	}
	info, err := v.introspect(ctx, token)
	if err != nil {
		return nil, err
	}
	if info.Active {
		v.store(key, info)
	}
	return info, nil
}

func (v *IntrospectionValidator) introspect(ctx context.Context, token string) (*TokenInfo, error) {
	form := url.Values{"token": {token}, "token_type_hint": {tokenTypeHintAccessToken}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, v.endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, ErrIntrospection.WithCause(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	if v.clientID != "" {
		req.SetBasicAuth(url.QueryEscape(v.clientID), url.QueryEscape(v.clientSecret))
	}
	resp, err := v.client.Do(req)
	if err != nil {
		return nil, ErrIntrospection.WithCause(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxResponseSize))
		return nil, ErrIntrospection.WithCause(fmt.Errorf("introspection endpoint returned status %d", resp.StatusCode))
	}
	var r introspectionResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxResponseSize)).Decode(&r); err != nil {
		return nil, ErrIntrospection.WithCause(err)
	}
	return r.tokenInfo()
}

func (v *IntrospectionValidator) load(key [sha256.Size]byte) (*TokenInfo, bool) {
	if v.cacheTTL <= 0 {
		return nil, false
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	e, ok := v.cache[key]
	if !ok {
		return nil, false
	}
	if !v.now().Before(e.expireAt) {
		delete(v.cache, key)
		return nil, false
	}
	return e.info, true
}

func (v *IntrospectionValidator) store(key [sha256.Size]byte, info *TokenInfo) {
	if v.cacheTTL <= 0 || v.cacheSize <= 0 {
		return
	}
	now := v.now()
	expireAt := now.Add(v.cacheTTL)
	if !info.Expiry.IsZero() && info.Expiry.Before(expireAt) {
		expireAt = info.Expiry
	}
	if !now.Before(expireAt) {
		return
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	if len(v.cache) >= v.cacheSize {
		for k, e := range v.cache {
			if !now.Before(e.expireAt) {
				delete(v.cache, k)
			}
		}
		if len(v.cache) >= v.cacheSize {
			return
		}
	}
	v.cache[key] = cacheEntry{info: info, expireAt: expireAt}
}

// introspectionResponse is the RFC 7662 section 2.2 response.
type introspectionResponse struct {
	Active   bool            `json:"active"`
	Scope    string          `json:"scope"`
	ClientID string          `json:"client_id"`
	Username string          `json:"username"`
	Exp      int64           `json:"exp"`
	Sub      string          `json:"sub"`
	Aud      json.RawMessage `json:"aud"`
	Iss      string          `json:"iss"`

	raw map[string]any
}

func (r *introspectionResponse) UnmarshalJSON(data []byte) error {
	type alias introspectionResponse
	if err := json.Unmarshal(data, (*alias)(r)); err != nil {
		return err
	}
	return json.Unmarshal(data, &r.raw)
}

func (r *introspectionResponse) tokenInfo() (*TokenInfo, error) {
	info := &TokenInfo{
		Active:   r.Active,
		Subject:  r.Sub,
		ClientID: r.ClientID,
		Username: r.Username,
		Issuer:   r.Iss,
		Scopes:   strings.Fields(r.Scope),
		Extra:    r.raw,
	}
	if r.Exp > 0 {
		info.Expiry = time.Unix(r.Exp, 0)
	}
	// aud may be a single string or an array of strings.
	if len(r.Aud) > 0 && string(r.Aud) != "null" {
		var aud string
		if err := json.Unmarshal(r.Aud, &aud); err == nil {
			info.Audience = []string{aud}
		} else if err := json.Unmarshal(r.Aud, &info.Audience); err != nil {
			return nil, ErrIntrospection.WithCause(err)
		}
	}
	return info, nil
}
