package memcache

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"
)

var testClient *Memcache

func Test_client_Set(t *testing.T) {
	type args struct {
		c    context.Context
		item *Item
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "set value", args: args{c: context.Background(), item: &Item{Key: "Test_client_Set", Value: []byte("abc")}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testClient.Set(tt.args.c, tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("client.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_Add(t *testing.T) {
	type args struct {
		c    context.Context
		item *Item
	}
	key := fmt.Sprintf("Test_client_Add_%d", time.Now().Unix())
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "add not exist value", args: args{c: context.Background(), item: &Item{Key: key, Value: []byte("abc")}}, wantErr: false},
		{name: "add exist value", args: args{c: context.Background(), item: &Item{Key: key, Value: []byte("abc")}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testClient.Add(tt.args.c, tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("client.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_Replace(t *testing.T) {
	key := fmt.Sprintf("Test_client_Replace_%d", time.Now().Unix())
	ekey := "Test_client_Replace_exist"
	testClient.Set(context.Background(), &Item{Key: ekey, Value: []byte("ok")})
	type args struct {
		c    context.Context
		item *Item
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "not exist value", args: args{c: context.Background(), item: &Item{Key: key, Value: []byte("abc")}}, wantErr: true},
		{name: "exist value", args: args{c: context.Background(), item: &Item{Key: ekey, Value: []byte("abc")}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testClient.Replace(tt.args.c, tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("client.Replace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_CompareAndSwap(t *testing.T) {
	key := fmt.Sprintf("Test_client_CompareAndSwap_%d", time.Now().Unix())
	ekey := "Test_client_CompareAndSwap_k"
	testClient.Set(context.Background(), &Item{Key: ekey, Value: []byte("old")})
	cas := testClient.Get(context.Background(), ekey).Item().cas
	type args struct {
		c    context.Context
		item *Item
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "not exist value", args: args{c: context.Background(), item: &Item{Key: key, Value: []byte("abc")}}, wantErr: true},
		{name: "exist value", args: args{c: context.Background(), item: &Item{Key: ekey, cas: cas, Value: []byte("abc")}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testClient.CompareAndSwap(tt.args.c, tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("client.CompareAndSwap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_Get(t *testing.T) {
	key := fmt.Sprintf("Test_client_Get_%d", time.Now().Unix())
	ekey := "Test_client_Get_k"
	testClient.Set(context.Background(), &Item{Key: ekey, Value: []byte("old")})
	type args struct {
		c   context.Context
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "not exist value", args: args{c: context.Background(), key: key}, wantErr: true},
		{name: "exist value", args: args{c: context.Background(), key: ekey}, wantErr: false, want: "old"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var res string
			if err := testClient.Get(tt.args.c, tt.args.key).Scan(&res); (err != nil) != tt.wantErr || res != tt.want {
				t.Errorf("client.Get() = %v, want %v, got err: %v, want err: %v", err, tt.want, err, tt.wantErr)
			}
		})
	}
}

func Test_client_Touch(t *testing.T) {
	key := fmt.Sprintf("Test_client_Touch_%d", time.Now().Unix())
	ekey := "Test_client_Touch_k"
	testClient.Set(context.Background(), &Item{Key: ekey, Value: []byte("old")})
	type args struct {
		c       context.Context
		key     string
		timeout int32
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "not exist value", args: args{c: context.Background(), key: key, timeout: 100000}, wantErr: true},
		{name: "exist value", args: args{c: context.Background(), key: ekey, timeout: 100000}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testClient.Touch(tt.args.c, tt.args.key, tt.args.timeout); (err != nil) != tt.wantErr {
				t.Errorf("client.Touch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_Delete(t *testing.T) {
	key := fmt.Sprintf("Test_client_Delete_%d", time.Now().Unix())
	ekey := "Test_client_Delete_k"
	testClient.Set(context.Background(), &Item{Key: ekey, Value: []byte("old")})
	type args struct {
		c   context.Context
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "not exist value", args: args{c: context.Background(), key: key}, wantErr: true},
		{name: "exist value", args: args{c: context.Background(), key: ekey}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testClient.Delete(tt.args.c, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("client.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_Increment(t *testing.T) {
	key := fmt.Sprintf("Test_client_Increment_%d", time.Now().Unix())
	ekey := "Test_client_Increment_k"
	testClient.Set(context.Background(), &Item{Key: ekey, Value: []byte("1")})
	type args struct {
		c     context.Context
		key   string
		delta uint64
	}
	tests := []struct {
		name         string
		args         args
		wantNewValue uint64
		wantErr      bool
	}{
		{name: "not exist value", args: args{c: context.Background(), key: key, delta: 10}, wantErr: true},
		{name: "exist value", args: args{c: context.Background(), key: ekey, delta: 10}, wantErr: false, wantNewValue: 11},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewValue, err := testClient.Increment(tt.args.c, tt.args.key, tt.args.delta)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.Increment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNewValue != tt.wantNewValue {
				t.Errorf("client.Increment() = %v, want %v", gotNewValue, tt.wantNewValue)
			}
		})
	}
}

func Test_client_Decrement(t *testing.T) {
	key := fmt.Sprintf("Test_client_Decrement_%d", time.Now().Unix())
	ekey := "Test_client_Decrement_k"
	testClient.Set(context.Background(), &Item{Key: ekey, Value: []byte("100")})
	type args struct {
		c     context.Context
		key   string
		delta uint64
	}
	tests := []struct {
		name         string
		args         args
		wantNewValue uint64
		wantErr      bool
	}{
		{name: "not exist value", args: args{c: context.Background(), key: key, delta: 10}, wantErr: true},
		{name: "exist value", args: args{c: context.Background(), key: ekey, delta: 10}, wantErr: false, wantNewValue: 90},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewValue, err := testClient.Decrement(tt.args.c, tt.args.key, tt.args.delta)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.Decrement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNewValue != tt.wantNewValue {
				t.Errorf("client.Decrement() = %v, want %v", gotNewValue, tt.wantNewValue)
			}
		})
	}
}

func Test_client_GetMulti(t *testing.T) {
	key := fmt.Sprintf("Test_client_GetMulti_%d", time.Now().Unix())
	ekey1 := "Test_client_GetMulti_k1"
	ekey2 := "Test_client_GetMulti_k2"
	testClient.Set(context.Background(), &Item{Key: ekey1, Value: []byte("1")})
	testClient.Set(context.Background(), &Item{Key: ekey2, Value: []byte("2")})
	keys := []string{key, ekey1, ekey2}
	rows, err := testClient.GetMulti(context.Background(), keys)
	if err != nil {
		t.Errorf("client.GetMulti() error = %v, wantErr %v", err, nil)
	}
	tests := []struct {
		key          string
		wantNewValue string
		wantErr      bool
		nilItem      bool
	}{
		{key: key, wantErr: true, nilItem: true},
		{key: ekey1, wantErr: false, wantNewValue: "1", nilItem: false},
		{key: ekey2, wantErr: false, wantNewValue: "2", nilItem: false},
	}
	if reflect.DeepEqual(keys, rows.Keys()) {
		t.Errorf("got %v, expect: %v", rows.Keys(), keys)
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			var gotNewValue string
			err = rows.Scan(tt.key, &gotNewValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("rows.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNewValue != tt.wantNewValue {
				t.Errorf("rows.Value() = %v, want %v", gotNewValue, tt.wantNewValue)
			}
			if (rows.Item(tt.key) == nil) != tt.nilItem {
				t.Errorf("rows.Item() = %v, want %v", rows.Item(tt.key) == nil, tt.nilItem)
			}
		})
	}
	err = rows.Close()
	if err != nil {
		t.Errorf("client.Replies.Close() error = %v, wantErr %v", err, nil)
	}
}

func Test_client_Conn(t *testing.T) {
	conn := testClient.Conn(context.Background())
	defer conn.Close()
	if conn == nil {
		t.Errorf("expect get conn, get nil")
	}
}
