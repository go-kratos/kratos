package group

import "fmt"

type Counter struct {
	Value int
}

func (c *Counter) Incr() {
	c.Value++
}

func ExampleGroup_Get() {
	group := NewGroup(func() any {
		fmt.Println("Only Once")
		return &Counter{}
	})

	// Create a new Counter
	group.Get("pass").(*Counter).Incr()

	// Get the created Counter again.
	group.Get("pass").(*Counter).Incr()
	// Output:
	// Only Once
}

func ExampleGroup_Reset() {
	group := NewGroup(func() any {
		return &Counter{}
	})

	// Reset the new function and clear all created objects.
	group.Reset(func() any {
		fmt.Println("reset")
		return &Counter{}
	})

	// Create a new Counter
	group.Get("pass").(*Counter).Incr()
	// Output:reset
}
