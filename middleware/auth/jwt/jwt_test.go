package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
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

type CustomerClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func TestJWTServerParse(t *testing.T) {
	testKey := "testKey"
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomerClaims{})
	token, err := claims.SignedString([]byte(testKey))
	if err != nil {
		panic(err)
	}
	ctx := transport.NewServerContext(
		context.Background(),
		&Transport{
			reqHeader: newTokenHeader(authorizationKey, fmt.Sprintf(bearerFormat, token)),
		},
	)

	next := func(ctx context.Context, req interface{}) (interface{}, error) {
		testToken, _ := FromContext(ctx)
		if _, ok := testToken.(*CustomerClaims); ok {
			t.Log("good result, token claims is CustomerClaims") // 期望打印
		} else {
			if _, ok := testToken.(jwt.MapClaims); ok {
				t.Log("bad result, token claims is MapClaims") // 实际打印
			}
			t.Fatal("fail")
		}
		return "reply", nil
	}

	server := Server(
		func(token *jwt.Token) (interface{}, error) { return []byte(testKey), nil },
		WithServerClaims(func() jwt.Claims { return &CustomerClaims{} }),
	)(next)

	_, err2 := server(ctx, "customer claim")
	if err2 != nil {
		t.Fatal("fail", err2)
	}
}

func TestJWTServerConcurrentWrite(t *testing.T) {
	testKey := "testKey"
	claims1 := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomerClaims{
		Name: "1",
	})
	claims2 := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomerClaims{
		Name: "2",
	})
	token1, err := claims1.SignedString([]byte(testKey))
	if err != nil {
		panic(err)
	}
	token2, err := claims2.SignedString([]byte(testKey))
	if err != nil {
		panic(err)
	}
	ctx1 := transport.NewServerContext(
		context.Background(),
		&Transport{
			reqHeader: newTokenHeader(authorizationKey, fmt.Sprintf(bearerFormat, token1)),
		},
	)
	ctx2 := transport.NewServerContext(
		context.Background(),
		&Transport{
			reqHeader: newTokenHeader(authorizationKey, fmt.Sprintf(bearerFormat, token2)),
		},
	)

	counter := 1
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	next := func(ctx context.Context, req interface{}) (interface{}, error) {
		if counter == 1 {
			var name string
			testToken1, _ := FromContext(ctx)
			if customerClaims, ok := testToken1.(*CustomerClaims); !ok {
				t.Fatal("claims is not *CustomerClaims")
			} else {
				name = customerClaims.Name
			}

			counter++
			ch1 <- struct{}{}
			<-ch2
			testToken2, _ := FromContext(ctx)
			if customerClaims, ok := testToken2.(*CustomerClaims); !ok {
				t.Fatal("claims is not *CustomerClaims")
			} else {
				if customerClaims.Name != name {
					t.Fatal("claims were modified concurrently")
				}
			}
			return nil, nil
		}
		ch2 <- struct{}{}
		return "reply", nil
	}

	server := Server(
		func(token *jwt.Token) (interface{}, error) { return []byte(testKey), nil },
		WithServerClaims(func() jwt.Claims { return &CustomerClaims{} }),
	)(next)

	go func() {
		_, err2 := server(ctx1, "first request")
		if err2 != nil {
			t.Error("fail", err)
		}
	}()
	<-ch1
	_, err = server(ctx2, "second request")
	if err != nil {
		t.Fatal("fail", err)
	}
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
			if !errors.Is(test.exceptErr, err2) {
				t.Errorf("except error %v, but got %v", test.exceptErr, err2)
			}
			if test.exceptErr == nil {
				if testToken == nil {
					t.Errorf("except testToken not nil, but got nil")
				}
				_, ok := testToken.(jwt.MapClaims)
				if !ok {
					t.Errorf("except testToken is jwt.MapClaims, but got %T", testToken)
				}
			}
		})
	}
}

func TestClient(t *testing.T) {
	testKey := "testKey"

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{})
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
			if !errors.Is(test.expectError, err2) {
				t.Errorf("except error %v, but got %v", test.expectError, err2)
			}
			if err2 == nil {
				if !reflect.DeepEqual(header.Get(authorizationKey), fmt.Sprintf(bearerFormat, token)) {
					t.Errorf("except header %s, but got %s", fmt.Sprintf(bearerFormat, token), header.Get(authorizationKey))
				}
			}
		})
	}
}

func TestTokenExpire(t *testing.T) {
	testKey := "testKey"
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Millisecond)),
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
	if !errors.Is(ErrTokenExpired, err2) {
		t.Errorf("except error %v, but got %v", ErrTokenExpired, err2)
	}
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
	if !errors.Is(test.exceptErr, err2) {
		t.Errorf("except error %v, but got %v", test.exceptErr, err2)
	}
	if test.exceptErr == nil {
		if testToken == nil {
			t.Errorf("except testToken not nil, but got nil")
		}
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
		if !errors.Is(test.expectError, err2) {
			t.Errorf("except error %v, but got %v", test.expectError, err2)
		}
		if err2 == nil {
			if !reflect.DeepEqual(header.Get(authorizationKey), fmt.Sprintf(bearerFormat, token)) {
				t.Errorf("except header %s, but got %s", fmt.Sprintf(bearerFormat, token), header.Get(authorizationKey))
			}
		}
	})
}

func TestClientWithHeader(t *testing.T) {
	testKey := "testKey"
	mapClaims := jwt.MapClaims{}
	mapClaims["name"] = "xiaoli"
	tokenHeader := map[string]interface{}{
		"test": "test",
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	for k, v := range tokenHeader {
		claims.Header[k] = v
	}
	token, err := claims.SignedString([]byte(testKey))
	if err != nil {
		panic(err)
	}
	tProvider := func(*jwt.Token) (interface{}, error) {
		return []byte(testKey), nil
	}
	next := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "reply", nil
	}
	handler := Client(tProvider, WithClaims(mapClaims), WithTokenHeader(tokenHeader))(next)
	header := &headerCarrier{}
	_, err2 := handler(transport.NewClientContext(context.Background(), &Transport{reqHeader: header}), "ok")
	if err2 != nil {
		t.Errorf("except error nil, but got %v", err2)
	}
	if !reflect.DeepEqual(header.Get(authorizationKey), fmt.Sprintf(bearerFormat, token)) {
		t.Errorf("except header %s, but got %s", fmt.Sprintf(bearerFormat, token), header.Get(authorizationKey))
	}
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
		if !errors.Is(test.expectError, err2) {
			t.Errorf("except error %v, but got %v", test.expectError, err2)
		}
		if err2 == nil {
			if !reflect.DeepEqual(header.Get(authorizationKey), fmt.Sprintf(bearerFormat, token)) {
				t.Errorf("except header %s, but got %s", fmt.Sprintf(bearerFormat, token), header.Get(authorizationKey))
			}
		}
	})
}
