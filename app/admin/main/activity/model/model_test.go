package model

import "testing"

func Test_Funcs(t *testing.T) {
	act := new(ActMatchsObject)
	s := act.TableName()
	t.Logf("tablename  %s", s)
}
