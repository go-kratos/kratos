package dao

import (
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/model"
	"go-common/library/cache/redis"
	"testing"
)

func TestDao_Lock(t *testing.T) {
	Convey("Lock", t, func() {
		once.Do(startService)
		key := "test_lock_key" + randomString(5)
		ttl := 2000
		retryDelay := 10000
		err, gotLock, lockValue := d.Lock(ctx, key, ttl, 0, retryDelay)
		So(err, ShouldBeNil)
		So(gotLock, ShouldEqual, true)
		So(lockValue, ShouldNotBeNil)
	})

	Convey("Lock Twice", t, func() {
		once.Do(startService)
		key := "test_lock_key" + randomString(5)
		ttl := 2000
		retryDelay := 10000
		err, gotLock, lockValue := d.Lock(ctx, key, ttl, 0, retryDelay)
		So(err, ShouldBeNil)
		So(gotLock, ShouldEqual, true)
		So(lockValue, ShouldNotBeNil)

		err, gotLock, lockValue = d.Lock(ctx, key, ttl, 3, retryDelay)
		So(err, ShouldBeNil)
		So(gotLock, ShouldEqual, false)

	})
}

func TestDao_UnLock(t *testing.T) {
	Convey("UnLock", t, func() {
		once.Do(startService)
		key := "test_lock_key" + randomString(5)
		value := "test"
		err := d.UnLock(ctx, key, value)
		So(err, ShouldEqual, redis.ErrNil)
	})

	Convey("UnLock wrong value", t, func() {
		once.Do(startService)
		key := "test_lock_key" + randomString(5)
		ttl := 2000
		retryDelay := 10000
		err, gotLock, lockValue := d.Lock(ctx, key, ttl, 3, retryDelay)
		So(err, ShouldBeNil)
		So(gotLock, ShouldEqual, true)
		So(lockValue, ShouldNotBeNil)
		value := "test"
		err = d.UnLock(ctx, key, value)
		So(err, ShouldEqual, ErrUnLockGet)

	})

	Convey("LockAndUnLock", t, func() {

		once.Do(startService)
		key := "test_lock_key" + randomString(5)
		ttl := 2000
		retryDelay := 10000
		err, gotLock, lockValue := d.Lock(ctx, key, ttl, 0, retryDelay)
		So(err, ShouldBeNil)
		So(gotLock, ShouldEqual, true)
		So(lockValue, ShouldNotBeNil)

		err = d.UnLock(ctx, key, lockValue)
		So(err, ShouldBeNil)

		err, gotLock, lockValue = d.Lock(ctx, key, ttl, 0, retryDelay)
		So(err, ShouldBeNil)
		So(gotLock, ShouldEqual, true)
		So(lockValue, ShouldNotBeNil)

	})
}

func TestDao_ForceUnLock(t *testing.T) {
	Convey("ForceUnLock", t, func() {
		once.Do(startService)
		key := "test_lock_key" + randomString(5)
		err := d.ForceUnLock(ctx, key)
		So(err, ShouldBeNil)
	})

	Convey("LockAndForceUnLock", t, func() {

		once.Do(startService)
		key := "test_lock_key" + randomString(5)
		ttl := 2000
		retryDelay := 10000
		err, gotLock, lockValue := d.Lock(ctx, key, ttl, 0, retryDelay)
		So(err, ShouldBeNil)
		So(gotLock, ShouldEqual, true)
		So(lockValue, ShouldNotBeNil)

		err = d.ForceUnLock(ctx, key)
		So(err, ShouldBeNil)

		err, gotLock, lockValue = d.Lock(ctx, key, ttl, 0, retryDelay)
		So(err, ShouldBeNil)
		So(gotLock, ShouldEqual, true)
		So(lockValue, ShouldNotBeNil)

	})

}

func TestDao_LockTransactionId(t *testing.T) {
	Convey("LockTid Twice", t, func() {
		tid := GetTid(model.PAYTYPE, "test")
		err := d.LockTransactionId(ctx, tid)
		So(err, ShouldBeNil)
		err = d.LockTransactionId(ctx, tid)
		So(err, ShouldEqual, ErrLockFailed)
	})
}

func TestDao_LockUser(t *testing.T) {
	Convey("LockUser", t, func() {
		uid := r.Int63n(100000000)
		err, gotLock, localValue := d.LockUser(ctx, uid)
		So(err, ShouldBeNil)
		So(gotLock, ShouldEqual, true)

		err = d.UnLockUser(ctx, uid, localValue)
		So(err, ShouldBeNil)

		err, gotLock, localValue = d.LockUser(ctx, uid)
		So(err, ShouldBeNil)
		So(gotLock, ShouldEqual, true)

	})

	Convey("LockUser Twice", t, func() {
		uid := r.Int63n(100000000)
		err, gotLock, _ := d.LockUser(ctx, uid)
		So(err, ShouldBeNil)
		So(gotLock, ShouldEqual, true)

		err, gotLock, _ = d.LockUser(ctx, uid)
		So(err, ShouldBeNil)
		So(gotLock, ShouldEqual, false)

	})
}
