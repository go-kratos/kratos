package service

import (
	"context"
	"encoding/json"
	"testing"
)

func TestService_RSAKey(t *testing.T) {
	once.Do(startService)
	originRSA, err := s.RSAKeyOrigin(context.TODO())
	if err != nil {
		t.Errorf("failed to get origin public key, service.RSAKeyOrigin() error(%v)", err)
		t.FailNow()
	} else {
		str, _ := json.Marshal(originRSA)
		t.Logf("cloud RSA: %s", str)
	}
	cloudRSA := s.RSAKey(context.TODO())
	if cloudRSA.Key == originRSA.Key {
		t.Errorf("cloud RSA public key cannot be equal to origin RSA public key")
		t.FailNow()
	} else {
		str, _ := json.Marshal(cloudRSA)
		t.Logf("cloud RSA: %s", str)
	}
}
