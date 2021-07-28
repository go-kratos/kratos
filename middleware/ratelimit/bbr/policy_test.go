package bbr

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func GetRollingPolicy() *RollingPolicy {
	w := NewWindow(WindowOpts{Size: 10})
	return NewRollingPolicy(w, RollingPolicyOpts{BucketDuration: 300 * time.Millisecond})
}

func Handler(t *testing.T, table []map[string][]int) {
	for _, hm := range table {
		var totalTs, lastOffset int
		offsetAndPoints := hm["offsetAndPoints"]
		timeSleep := hm["timeSleep"]
		policy := GetRollingPolicy()
		for i, n := range timeSleep {
			totalTs += n
			time.Sleep(time.Duration(n) * time.Millisecond)
			policy.Add(1)
			offset, points := offsetAndPoints[2*i], offsetAndPoints[2*i+1]

			if int(policy.window.window[offset].Points[0]) != points {
				t.Errorf("error, time since last append: %vms, last offset: %v", totalTs, lastOffset)
			}
			lastOffset = offset
		}
	}
}

func TestRollingPolicy_Add(t *testing.T) {
	rand.Seed(time.Now().Unix())

	// test add after 400ms and 601ms relative to the policy created time
	policy := GetRollingPolicy()
	time.Sleep(400 * time.Millisecond)
	policy.Add(1)
	time.Sleep(201 * time.Millisecond)
	policy.Add(1)
	for _, b := range policy.window.window {
		fmt.Println(b.Points)
	}
	if int(policy.window.window[1].Points[0]) != 1 {
		t.Errorf("error, time since last append: %vms, last offset: %v", 300, 0)
	}
	if int(policy.window.window[2].Points[0]) != 1 {
		t.Errorf("error, time since last append: %vms, last offset: %v", 301, 0)
	}

	// test func timespan return real span
	table := []map[string][]int{
		{
			"timeSleep":       []int{294, 3200},
			"offsetAndPoints": []int{0, 1, 0, 1},
		},
		{
			"timeSleep":       []int{305, 3200, 6400},
			"offsetAndPoints": []int{1, 1, 1, 1, 1, 1},
		},
	}

	Handler(t, table)
}
