package runner

import (
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
)

func TestContext_Pop_Root(t *testing.T) {
	assert.Panics(t, func() {
		NewContext().pop()
	}, "root context should panic if popped")
}

func TestContext_Errors(t *testing.T) {
	is := assert.New(t)
	err := errors.New("foobar")

	ctx := NewContext()
	is.NoError(ctx.Err(), "initially there should be no error")

	ctx.push()
	is.NoError(ctx.Err(), "child contexts should also have no error")

	ctx.SetErr(err)
	is.Equal(err, ctx.Err(), "child contexts should see set error")

	ctx.pop()
	is.Equal(err, ctx.Err(), "error should be propogated to parent")
}

func TestContext_GetSet(t *testing.T) {
	is := assert.New(t)
	key := "foo"
	unknown := "fizz"
	val := "bar"

	ctx := NewContext()
	out, found := ctx.Get(unknown)
	is.False(found, "unknown kv should not be found")

	ctx.Set(key, val)
	out, found = ctx.Get(key)
	is.True(found, "known kv should be found")
	is.Equal(val, out, "known kv should match")

	ctx.push()
	out, found = ctx.Get(unknown)
	is.False(found, "unknown kv should not be found on child")

	out, found = ctx.Get(key)
	is.True(found, "known kv set on parent should be available on child")
	is.Equal(val, out, "known kv should have same value as parent")

	ctx.Set(unknown, "buzz")
	out, found = ctx.Get(unknown)
	is.True(found, "now set unknown kv should be found on child")

	newVal := "baz"
	ctx.Set(key, newVal)
	out, found = ctx.Get(key)
	is.Equal(newVal, out, "new key value should be on the child")

	ctx.pop()
	out, found = ctx.Get(unknown)
	is.False(found, "parent should not know about kvs set by children")

	out, found = ctx.Get(key)
	is.Equal(val, out, "parent should not know new value for kv set by children")
}
