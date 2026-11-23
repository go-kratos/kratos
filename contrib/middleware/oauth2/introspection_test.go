package oauth2

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	kerrors "github.com/go-kratos/kratos/v3/errors"
)

func newIntrospectionServer(t *testing.T, calls *atomic.Int32) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if id, secret, ok := r.BasicAuth(); !ok || id != "rs" || secret != "s3cret" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err := r.ParseForm(); err != nil || r.PostForm.Get("token_type_hint") != tokenTypeHintAccessToken {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.PostForm.Get("token") {
		case "good":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"active":    true,
				"scope":     "read write",
				"client_id": "app",
				"username":  "alice",
				"sub":       "u-1",
				"aud":       "api",
				"iss":       "https://issuer",
				"exp":       time.Now().Add(time.Hour).Unix(),
				"tenant":    "t-1",
			})
		case "multi-aud":
			_, _ = w.Write([]byte(`{"active":true,"aud":["a","b"]}`))
		case "bad-aud":
			_, _ = w.Write([]byte(`{"active":true,"aud":123}`))
		case "broken":
			_, _ = w.Write([]byte(`{`))
		case "fail":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			_, _ = w.Write([]byte(`{"active":false}`))
		}
	}))
}

func TestIntrospectionValidator(t *testing.T) {
	var calls atomic.Int32
	srv := newIntrospectionServer(t, &calls)
	defer srv.Close()
	v := NewIntrospectionValidator(srv.URL, srv.Client(), WithBasicAuth("rs", "s3cret"))

	info, err := v.Validate(context.Background(), "good")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !info.Active || info.Subject != "u-1" || info.ClientID != "app" || info.Username != "alice" || info.Issuer != "https://issuer" {
		t.Fatalf("unexpected token info: %+v", info)
	}
	if !info.HasScopes("read", "write") || info.HasScopes("admin") {
		t.Fatalf("unexpected scopes: %v", info.Scopes)
	}
	if len(info.Audience) != 1 || info.Audience[0] != "api" {
		t.Fatalf("unexpected audience: %v", info.Audience)
	}
	if info.Expiry.IsZero() || info.Extra["tenant"] != "t-1" {
		t.Fatalf("unexpected expiry or extra: %v %v", info.Expiry, info.Extra)
	}

	info, err = v.Validate(context.Background(), "multi-aud")
	if err != nil || len(info.Audience) != 2 {
		t.Fatalf("multi-aud: info=%+v err=%v", info, err)
	}

	info, err = v.Validate(context.Background(), "unknown")
	if err != nil || info.Active {
		t.Fatalf("inactive: info=%+v err=%v", info, err)
	}

	for _, token := range []string{"bad-aud", "broken", "fail"} {
		if _, err := v.Validate(context.Background(), token); !kerrors.Is(err, ErrIntrospection) {
			t.Fatalf("%s: err = %v, want %v", token, err, ErrIntrospection)
		}
	}

	wrong := NewIntrospectionValidator(srv.URL, srv.Client(), WithBasicAuth("rs", "wrong"))
	if _, err := wrong.Validate(context.Background(), "good"); !kerrors.Is(err, ErrIntrospection) {
		t.Fatalf("wrong credentials: err = %v, want %v", err, ErrIntrospection)
	}
}

func TestIntrospectionValidatorCache(t *testing.T) {
	var calls atomic.Int32
	srv := newIntrospectionServer(t, &calls)
	defer srv.Close()

	now := time.Now()
	v := NewIntrospectionValidator(srv.URL, srv.Client(), WithBasicAuth("rs", "s3cret"), WithCacheTTL(time.Minute))
	v.now = func() time.Time { return now }

	for range 3 {
		if _, err := v.Validate(context.Background(), "good"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("introspection called %d times, want 1", got)
	}

	// Inactive results are not cached.
	for range 2 {
		if _, err := v.Validate(context.Background(), "unknown"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if got := calls.Load(); got != 3 {
		t.Fatalf("introspection called %d times, want 3", got)
	}

	// The cache entry expires after the TTL.
	now = now.Add(2 * time.Minute)
	if _, err := v.Validate(context.Background(), "good"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := calls.Load(); got != 4 {
		t.Fatalf("introspection called %d times, want 4", got)
	}
}

func TestIntrospectionValidatorCacheBounds(t *testing.T) {
	now := time.Now()
	v := NewIntrospectionValidator("", nil, WithCacheTTL(time.Hour), WithCacheSize(1))
	v.now = func() time.Time { return now }

	// The entry never outlives the token expiry.
	v.store([32]byte{1}, &TokenInfo{Active: true, Expiry: now.Add(time.Minute)})
	if e := v.cache[[32]byte{1}]; !e.expireAt.Equal(now.Add(time.Minute)) {
		t.Fatalf("expireAt = %v, want %v", e.expireAt, now.Add(time.Minute))
	}

	// A full cache drops new entries instead of growing.
	v.store([32]byte{2}, &TokenInfo{Active: true})
	if len(v.cache) != 1 {
		t.Fatalf("cache size = %d, want 1", len(v.cache))
	}

	// Expired entries are evicted to make room.
	now = now.Add(2 * time.Minute)
	v.store([32]byte{2}, &TokenInfo{Active: true})
	if _, ok := v.cache[[32]byte{2}]; !ok || len(v.cache) != 1 {
		t.Fatalf("expected the expired entry to be replaced, cache = %v", v.cache)
	}

	// Already expired tokens are not cached.
	v.store([32]byte{3}, &TokenInfo{Active: true, Expiry: now.Add(-time.Second)})
	if _, ok := v.cache[[32]byte{3}]; ok {
		t.Fatal("expired token should not be cached")
	}
}
