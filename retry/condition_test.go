package retry

import "testing"

func TestByCode(t *testing.T) {
	testCases := []struct {
		cond     *ByCode
		status   int
		expected bool
	}{
		{cond: NewByCode(501), status: 501, expected: true},
		{cond: NewByCode(), status: 501, expected: false},
		{cond: NewByCode(501, 509), status: 501, expected: true},
		{cond: NewByCode(501, 509), status: 500, expected: false},
	}
	for _, testCase := range testCases {
		result := testCase.cond.Judge(Resp{Code: testCase.status})
		if result != testCase.expected {
			t.Errorf("cond:%v, status:%v, expected:%v, got:%v", testCase.cond.Codes, testCase.status, testCase.expected, result)
		}
	}
}

func TestByMetedata(t *testing.T) {
	testCases := []struct {
		cond     *ByMetedata
		metedata map[string][]string
		expected bool
	}{
		{cond: NewByMetedata("Grpc-Status", "5"), metedata: map[string][]string{"Grpc-Status": {"5"}}, expected: true},
		{cond: NewByMetedata("Grpc-Status", "5"), metedata: map[string][]string{"Grpc-Status": {"15", "5"}}, expected: true},
		{cond: NewByMetedata("Grpc-Status", "5"), metedata: map[string][]string{}, expected: false},
		{cond: NewByMetedata("Grpc-Status", "5"), metedata: map[string][]string{"Grpc-Status": {"10"}}, expected: false},
		{cond: NewByMetedata("Grpc-Status", "5", "10"), metedata: map[string][]string{"Grpc-Status": {"10"}}, expected: true},
	}
	for _, testCase := range testCases {
		result := testCase.cond.Judge(Resp{MD: testCase.metedata})
		if result != testCase.expected {
			t.Errorf("cond.Key:%v, cond.Vals:%v,metedata:%v, expected:%v, got:%v", testCase.cond.Key, testCase.cond.Vals, testCase.metedata, testCase.expected, result)
		}
	}
}
