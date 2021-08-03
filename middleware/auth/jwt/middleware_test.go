package jwt

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
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

func newTokenHeader(token string) *headerCarrier {
	header := &headerCarrier{}
	header.Set(JWTHeaderKey, token)
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
	var testKey = "testKey"
	mapClaims := jwt.MapClaims{}
	mapClaims["name"] = "xiaoli"
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	token, err := claims.SignedString([]byte(testKey))
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name          string
		ctx           context.Context
		signingMethod jwt.SigningMethod
		exceptErr     error
		key           string
	}{
		{
			name:          "normal",
			ctx:           transport.NewServerContext(context.Background(), &Transport{reqHeader: newTokenHeader(token)}),
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
			name:          "token invalid",
			ctx:           transport.NewServerContext(context.Background(), &Transport{reqHeader: newTokenHeader("123123123")}),
			signingMethod: jwt.SigningMethodHS256,
			exceptErr:     ErrTokenInvalid,
			key:           testKey,
		},
		{
			name:          "method invalid",
			ctx:           transport.NewServerContext(context.Background(), &Transport{reqHeader: newTokenHeader(token)}),
			signingMethod: jwt.SigningMethodES384,
			exceptErr:     ErrUnSupportSigningMethod,
			key:           testKey,
		},
		{
			name:          "miss key",
			ctx:           transport.NewServerContext(context.Background(), &Transport{reqHeader: newTokenHeader(token)}),
			signingMethod: jwt.SigningMethodHS256,
			exceptErr:     ErrMissingAccessSecret,
			key:           "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var testToken interface{}
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				t.Log(req)
				testToken = ctx.Value(JWTClaimsContextKey)
				return "reply", nil
			}
			server := Server(test.key, WithSigningMethod(test.signingMethod))(next)
			_, err2 := server(test.ctx, test.name)
			assert.Equal(t, test.exceptErr, err2)
			if test.exceptErr == nil {
				assert.NotNil(t, testToken)
			}
		})
	}
}

type tokeBuilder struct {
	token string
}

func (t tokeBuilder) GetToken() string {
	return t.token
}

func TestClient(t *testing.T) {
	var testKey = "testKey"
	mapClaims := jwt.MapClaims{}
	mapClaims["name"] = "xiaoli"
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	token, err := claims.SignedString([]byte(testKey))
	if err != nil {
		panic(err)
	}
	tProvider := tokeBuilder{
		token: token,
	}

	tests := []struct {
		name          string
		expectError   error
		tokenProvider TokenProvider
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
				assert.Equal(t, test.tokenProvider.GetToken(), header.Get(JWTHeaderKey))
			}
		})
	}
}
