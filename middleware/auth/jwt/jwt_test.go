package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/stretchr/testify/assert"
)

type headerCarrier http.Header

func (hc headerCarrier) Get(key string) string { return http.Header(hc).Get(key) }

func (hc headerCarrier) Set(key string, value string) { http.Header(hc).Set(key, value) }

// Keys lists the keys stored in this carrier.
func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}

func newTokenHeader(headerKey string, token string) *headerCarrier {
	header := &headerCarrier{}
	header.Set(headerKey, token)
	return header
}

type Transport struct {
	kind      transport.Kind
	endpoint  string
	operation string
	reqHeader transport.Header
}

func (tr *Transport) Kind() transport.Kind {
	return tr.kind
}

func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

func (tr *Transport) Operation() string {
	return tr.operation
}

func (tr *Transport) RequestHeader() transport.Header {
	return tr.reqHeader
}

func (tr *Transport) ReplyHeader() transport.Header {
	return nil
}

func TestServer(t *testing.T) {
	testKey := "testKey"
	mapClaims := jwt.MapClaims{}
	mapClaims["name"] = "xiaoli"
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	token, err := claims.SignedString([]byte(testKey))
	if err != nil {
		panic(err)
	}
	token = fmt.Sprintf(bearerFormat, token)
	tests := []struct {
		name          string
		ctx           context.Context
		signingMethod jwt.SigningMethod
		exceptErr     error
		key           string
	}{
		{
			name:          "normal",
			ctx:           transport.NewServerContext(context.Background(), &Transport{reqHeader: newTokenHeader(authorizationKey, token)}),
			signingMethod: jwt.SigningMethodHS256,
			exceptErr:     nil,
			key:           testKey,
		},
		{
			name:          "miss token",
			ctx:           transport.NewServerContext(context.Background(), &Transport{reqHeader: headerCarrier{}}),
			signingMethod: jwt.SigningMethodHS256,
			exceptErr:     ErrMissingJwtToken,
			key:           testKey,
		},
		{
			name: "token invalid",
			ctx: transport.NewServerContext(context.Background(), &Transport{
				reqHeader: newTokenHeader(authorizationKey, fmt.Sprintf(bearerFormat, "12313123")),
			}),
			signingMethod: jwt.SigningMethodHS256,
			exceptErr:     ErrTokenInvalid,
			key:           testKey,
		},
		{
			name:          "method invalid",
			ctx:           transport.NewServerContext(context.Background(), &Transport{reqHeader: newTokenHeader(authorizationKey, token)}),
			signingMethod: jwt.SigningMethodES384,
			exceptErr:     ErrUnSupportSigningMethod,
			key:           testKey,
		},
		{
			name:          "miss signing method",
			ctx:           transport.NewServerContext(context.Background(), &Transport{reqHeader: newTokenHeader(authorizationKey, token)}),
			signingMethod: nil,
			exceptErr:     nil,
			key:           testKey,
		},
		{
			name:          "miss signing method",
			ctx:           transport.NewServerContext(context.Background(), &Transport{reqHeader: newTokenHeader(authorizationKey, token)}),
			signingMethod: nil,
			exceptErr:     nil,
			key:           testKey,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var testToken jwt.Claims
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				t.Log(req)
				testToken, _ = FromContext(ctx)
				return "reply", nil
			}
			var server middleware.Handler
			if test.signingMethod != nil {
				server = Server(func(token *jwt.Token) (interface{}, error) {
					return []byte(test.key), nil
				}, WithSigningMethod(test.signingMethod))(next)
			} else {
				server = Server(func(token *jwt.Token) (interface{}, error) {
					return []byte(test.key), nil
				})(next)
			}
			_, err2 := server(test.ctx, test.name)
			assert.Equal(t, test.exceptErr, err2)
			if test.exceptErr == nil {
				assert.NotNil(t, testToken)
				_, ok := testToken.(jwt.MapClaims)
				assert.True(t, ok)
			}
		})
	}
}

func TestClient(t *testing.T) {
	testKey := "testKey"
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{})
	token, err := claims.SignedString([]byte(testKey))
	if err != nil {
		panic(err)
	}
	tProvider := func(*jwt.Token) (interface{}, error) {
		return []byte(testKey), nil
	}
	tests := []struct {
		name          string
		expectError   error
		tokenProvider jwt.Keyfunc
	}{
		{
			name:          "normal",
			expectError:   nil,
			tokenProvider: tProvider,
		},
		{
			name:          "miss token provider",
			expectError:   ErrNeedTokenProvider,
			tokenProvider: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				return "reply", nil
			}
			handler := Client(test.tokenProvider)(next)
			header := &headerCarrier{}
			_, err2 := handler(transport.NewClientContext(context.Background(), &Transport{reqHeader: header}), "ok")
			assert.Equal(t, test.expectError, err2)
			if err2 == nil {
				assert.Equal(t, fmt.Sprintf(bearerFormat, token), header.Get(authorizationKey))
			}
		})
	}
}

func TestTokenExpire(t *testing.T) {
	testKey := "testKey"
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Millisecond).Unix(),
	})
	token, err := claims.SignedString([]byte(testKey))
	if err != nil {
		panic(err)
	}
	token = fmt.Sprintf(bearerFormat, token)
	time.Sleep(time.Second)
	next := func(ctx context.Context, req interface{}) (interface{}, error) {
		t.Log(req)
		return "reply", nil
	}
	ctx := transport.NewServerContext(context.Background(), &Transport{reqHeader: newTokenHeader(authorizationKey, token)})
	server := Server(func(token *jwt.Token) (interface{}, error) {
		return []byte(testKey), nil
	}, WithSigningMethod(jwt.SigningMethodHS256))(next)
	_, err2 := server(ctx, "test expire token")
	assert.Equal(t, ErrTokenExpired, err2)
}

func TestMissingKeyFunc(t *testing.T) {
	testKey := "testKey"
	mapClaims := jwt.MapClaims{}
	mapClaims["name"] = "xiaoli"
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	token, err := claims.SignedString([]byte(testKey))
	if err != nil {
		panic(err)
	}
	token = fmt.Sprintf(bearerFormat, token)
	test := struct {
		name          string
		ctx           context.Context
		signingMethod jwt.SigningMethod
		exceptErr     error
		key           string
	}{
		name:          "miss key",
		ctx:           transport.NewServerContext(context.Background(), &Transport{reqHeader: newTokenHeader(authorizationKey, token)}),
		signingMethod: jwt.SigningMethodHS256,
		exceptErr:     ErrMissingKeyFunc,
		key:           "",
	}

	var testToken jwt.Claims
	next := func(ctx context.Context, req interface{}) (interface{}, error) {
		t.Log(req)
		testToken, _ = FromContext(ctx)
		return "reply", nil
	}
	server := Server(nil)(next)
	_, err2 := server(test.ctx, test.name)
	assert.Equal(t, test.exceptErr, err2)
	if test.exceptErr == nil {
		assert.NotNil(t, testToken)
	}
}

func TestClientWithClaims(t *testing.T) {
	testKey := "testKey"
	mapClaims := jwt.MapClaims{}
	mapClaims["name"] = "xiaoli"
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	token, err := claims.SignedString([]byte(testKey))
	if err != nil {
		panic(err)
	}
	tProvider := func(*jwt.Token) (interface{}, error) {
		return []byte(testKey), nil
	}
	test := struct {
		name          string
		expectError   error
		tokenProvider jwt.Keyfunc
	}{
		name:          "normal",
		expectError:   nil,
		tokenProvider: tProvider,
	}

	t.Run(test.name, func(t *testing.T) {
		next := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "reply", nil
		}
		handler := Client(test.tokenProvider, WithClaims(mapClaims))(next)
		header := &headerCarrier{}
		_, err2 := handler(transport.NewClientContext(context.Background(), &Transport{reqHeader: header}), "ok")
		assert.Equal(t, test.expectError, err2)
		if err2 == nil {
			assert.Equal(t, fmt.Sprintf(bearerFormat, token), header.Get(authorizationKey))
		}
	})
}

func TestClientMissKey(t *testing.T) {
	testKey := "testKey"
	mapClaims := jwt.MapClaims{}
	mapClaims["name"] = "xiaoli"
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	token, err := claims.SignedString([]byte(testKey))
	if err != nil {
		panic(err)
	}
	tProvider := func(*jwt.Token) (interface{}, error) {
		return nil, errors.New("some error")
	}
	test := struct {
		name          string
		expectError   error
		tokenProvider jwt.Keyfunc
	}{
		name:          "normal",
		expectError:   ErrGetKey,
		tokenProvider: tProvider,
	}

	t.Run(test.name, func(t *testing.T) {
		next := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "reply", nil
		}
		handler := Client(test.tokenProvider, WithClaims(mapClaims))(next)
		header := &headerCarrier{}
		_, err2 := handler(transport.NewClientContext(context.Background(), &Transport{reqHeader: header}), "ok")
		assert.Equal(t, test.expectError, err2)
		if err2 == nil {
			assert.Equal(t, fmt.Sprintf(bearerFormat, token), header.Get(authorizationKey))
		}
	})
}
