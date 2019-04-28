package intimacy

import (
	"fmt"
	"testing"
)

func TestLevelConf(t *testing.T) {
	var sumScore = 0

	level := 1
	for level < MaxLevel {
		conf := LevelConf[level]
		sumScore += conf.levelUpExp
		if sumScore != LevelConf[level+1].sumScore {
			t.Error("LevelConf")
		}
		level += 1
	}
}

type IncreaseCases struct {
	param Increaser
	ret   IncreaseResult
}

var cases = []IncreaseCases{
	{
		Increaser{0, 0, 10000000, true},
		IncreaseResult{
			false,
			[4]Result{
				{10000000, 10000000, 0, 20},
				{10000000, 10000000, 0, 20},
				{10000000, 10000000, 0, 20},
				{10000000, 10000000, 0, 20},
			},
		},
	},
	{
		Increaser{0, 0, 1, false},
		IncreaseResult{
			false,
			[4]Result{
				{1, 1, 1, 1},
				{1, 1, 1, 1},
				{1, 1, 1, 1},
				{1, 1, 1, 1},
			},
		},
	},
	{
		Increaser{0, 0, 201, false},
		IncreaseResult{
			false,
			[4]Result{
				{201, 201, 201, 2},
				{201, 201, 201, 2},
				{201, 201, 201, 2},
				{201, 201, 201, 2},
			},
		},
	},
	{
		Increaser{0, 0, 2700, false},
		IncreaseResult{
			true,
			[4]Result{
				{500, 500, 500, 2},
				{1500, 1500, 1500, 4},
				{1000, 1000, 1000, 3},
				{750, 750, 750, 3},
			},
		},
	},
	{
		Increaser{1700, 499, 3001, false},
		IncreaseResult{
			true,
			[4]Result{
				{501, 2201, 1000, 5},
				{2501, 4201, 3000, 7},
				{1501, 3201, 2000, 6},
				{1001, 2701, 1500, 6},
			},
		},
	},
	{
		Increaser{1700, 499, 2, false},
		IncreaseResult{
			true,
			[4]Result{
				{1, 1701, 500, 5},
				{2, 1702, 501, 5},
				{2, 1702, 501, 5},
				{2, 1702, 501, 5},
			},
		},
	},
	{
		Increaser{1700, 499, 2, false},
		IncreaseResult{
			true,
			[4]Result{
				{1, 1701, 500, 5},
				{2, 1702, 501, 5},
				{2, 1702, 501, 5},
				{2, 1702, 501, 5},
			},
		},
	},
}

// TestIncrease 测试亲密度增加逻辑
// 这是核心接口
func TestIncrease(t *testing.T) {
	testTimes := 2000
	for testTimes > 0 {
		testTimes -= 1
		for i, currCase := range cases {
			fmt.Println(currCase)
			ret := Increase(&currCase.param)

			if ret.passLimit != currCase.ret.passLimit {
				t.Error("case ", i, "passLimit 计算错误")
			}

			for n := 0; n < 4; n++ {
				ptr1 := &ret.retList[n]
				ptr2 := &currCase.ret.retList[n]

				if ptr1.score != ptr2.score {
					t.Error("case ", i, "retList ", n, " score 计算错误")
				}

				if ptr1.limit != ptr2.limit {
					t.Error("case ", i, "retList ", n, " limit 计算错误")
				}

				if ptr1.incr != ptr2.incr {
					t.Error("case ", i, "incr ", n, " score 计算错误")
				}

				if ptr1.level != ptr2.level {
					t.Error("case ", i, "level ", n, ptr1.level, " != ", ptr2.level)
				}
			}
		}
	}
}

func TestGetLevel(t *testing.T) {
	fmt.Println(GetLevel(201))
	for i, elem := range LevelConf {
		currLevel, currIntimacy := GetLevel(elem.sumScore)
		if currLevel != i || currIntimacy != 0 {
			t.Error("GetLevel 有错 ", currLevel, currIntimacy)
		}
	}
}
