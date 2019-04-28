/*
Package databusutil provides a util for building databus based async job with
single partition message aggregation and parallel consumption features.

Group

The group is the primary struct for working with this util.

Applications create groups by calling the package NewGroup function with a
databusutil config and a databus message chan.

To start a initiated group, the application must call the group Start method.

The application must call the group Close method when the application is
done with the group.

Callbacks

After a new group is created, the following callbacks: New, Split and Do must
be assigned, otherwise the job will not works as your expectation.

The callback New represents how the consume proc of the group parsing the target
object from a new databus message that it received for merging, if the error
returned is not nil, the consume proc will omit this message and continue.

A example of the callback New is:

	func newTestMsg(msg *databus.Message) (res interface{}, err error) {
		res = new(testMsg)
		if err = json.Unmarshal(msg.Value, &res); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
		}
		return
	}

The callback Split represents how the consume proc of the group getting the
sharding dimension from a databus message or the object parsed from the databus
message, it will be used along with the configuration item Num to decide which
merge goroutine to use to merge the parsed object. In more detail, if we take
the result of callback Split as sr, then the sharding result will be sr % Num.

A example of the callback Split is:

	func split(msg *databus.Message, data interface{}) int {
		t, ok := data.(*testMsg)
		if !ok {
			return 0
		}
		return int(t.Mid)
	}

If your messages is already assigned to their partitions corresponding to the split you want,
you may want to directly use its partition as split, here is the example:

	func anotherSplit(msg *databus.Message, data interface{}) int {
		return int(msg.Partition)
	}

Do not forget to ensure the max value your callback Split returns, as maxSplit,
greater than or equal to the configuration item Num, otherwise the merge
goroutines will not be fully used, in more detail, the last (Num - maxSplit)
merge goroutines are initiated by will never be used.

The callback Do represents how the merge proc of the group processing the merged
objects, define your business in it.

A example of the callback Do is:

	func do(msgs []interface{}) {
		for _, m := range msgs {
			// process messages you merged here, the example type asserts and prints each
			if msg, ok := m.(*testMsg); ok {
				fmt.Printf("msg: %+v", msg)
			}
		}
	}

Usage Example

The typical usage for databusutil is:

	// new a databus to subscribe from
	dsSub := databus.New(dsSubConf)
	defer dsSub.Close()
	// new a group
	g := NewGroup(
		c,
		dsSub.Messages(),
	)
	// fill callbacks
	g.New = yourNewFunc
	g.Split = yourSplitFunc
	g.Do = yourDoFunc
	// start the group
	g.Start()
	// must close the group before the job exits
	defer g.Close()
	// signal handler
*/
package databusutil
