package msgpack

import (
	"reflect"
	"testing"
)

type loginRequest struct {
	UserName string
	Password string
}

type testModel struct {
	ID   int32
	Name string
}

func TestName(t *testing.T) {
	c := new(codec)
	if !reflect.DeepEqual("msgpack", c.Name()) {
		t.Errorf("Name() should be msgpack, but got %s", c.Name())
	}
}

func TestCodec(t *testing.T) {
	c := new(codec)
	t2 := testModel{ID: 1, Name: "name"}
	m, err := c.Marshal(&t2)
	if err != nil {
		t.Errorf("Marshal() should be nil, but got %s", err)
	}
	var t3 testModel
	err = c.Unmarshal(m, &t3)
	if err != nil {
		t.Errorf("Unmarshal() should be nil, but got %s", err)
	}
	if !reflect.DeepEqual(t2.ID, t3.ID) {
		t.Errorf("ID should be %d, but got %d", t2.ID, t3.ID)
	}
	if !reflect.DeepEqual(t3.Name, t2.Name) {
		t.Errorf("Name should be %s, but got %s", t2.Name, t3.Name)
	}

	request := loginRequest{
		UserName: "username",
		Password: "password",
	}
	m, err = c.Marshal(&request)
	if err != nil {
		t.Errorf("Marshal() should be nil, but got %s", err)
	}
	var req loginRequest
	err = c.Unmarshal(m, &req)
	if err != nil {
		t.Errorf("Unmarshal() should be nil, but got %s", err)
	}
	if !reflect.DeepEqual(req.Password, request.Password) {
		t.Errorf("ID should be %s, but got %s", req.Password, request.Password)
	}
	if !reflect.DeepEqual(req.UserName, request.UserName) {
		t.Errorf("Name should be %s, but got %s", req.UserName, request.UserName)
	}
}
