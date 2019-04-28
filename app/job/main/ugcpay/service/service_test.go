package service

import (
	"testing"
)

type c struct {
	balance          int64
	loss             int64
	expectUserRefund int64
	expectBizRefund  int64
}

func TestCalcRefundFee(t *testing.T) {
	var min int64 = -20000
	cases := []c{
		c{
			balance:          10000,
			loss:             20000,
			expectUserRefund: 20000,
			expectBizRefund:  0,
		}, c{
			balance:          0,
			loss:             10000,
			expectUserRefund: 10000,
			expectBizRefund:  0,
		}, c{
			balance:          -1,
			loss:             20000,
			expectUserRefund: 19999,
			expectBizRefund:  1,
		}, c{
			balance:          -19999,
			loss:             20000,
			expectUserRefund: 1,
			expectBizRefund:  19999,
		}, c{
			balance:          -20000,
			loss:             20000,
			expectUserRefund: 0,
			expectBizRefund:  20000,
		}, c{
			balance:          -30000,
			loss:             20000,
			expectUserRefund: 0,
			expectBizRefund:  20000,
		},
	}

	for _, c := range cases {
		bizRefund, userRefund := calcRefundFee(c.balance, c.loss, min)
		if userRefund != c.expectUserRefund {
			t.Fatalf("TestCalcRefundFee case: %+v expectUserRefund not right, actual: %d\n", c, userRefund)
		}
		if bizRefund != c.expectBizRefund {
			t.Fatalf("TestCalcRefundFee case: %+v expectBizRefund not right, actual: %d\n", c, bizRefund)
		}
	}
}
