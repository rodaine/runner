package runner

import (
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
)

func TestSubContext_Pop_Root(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		ctx := NewContext()
		ctx.push()
		newSubContext(ctx).pop()
	}, "root subcontext should panic if popped")
}

func TestSubContext_Errors(t *testing.T) {
	t.Parallel()

	is := assert.New(t)
	err := errors.New("foobar")

	ctx := NewContext()
	sctx := newSubContext(ctx)
	is.NoError(sctx.Err(), "initially subcontext should have no error")

	sctx.push()
	is.NoError(sctx.Err(), "child subcontexts should also have no error")

	sctx.SetErr(err)
	is.Equal(err, sctx.Err(), "child subcontexts should see set error")

	sctx.SetErr(nil)
	is.Equal(err, sctx.Err(), "subcontext errors cannot be unset via SetErr")

	sctx.pop()
	is.Equal(err, sctx.Err(), "error should propogate to subcontext parent")
	is.NoError(ctx.Err(), "parent context should not see errors in subcontext")

	sctx.unsetErr()
	is.NoError(sctx.Err(), "subcontext errors can only be unset via unsetErr")

	ctx.SetErr(err)
	sctx = newSubContext(ctx)
	is.NoError(sctx.Err(), "subcontext should not see parent errors")
}

func TestSubContext_GetSet(t *testing.T) {
	t.Parallel()

	is := assert.New(t)
	key := "foo"
	val := "bar"

	ctx := NewContext()
	ctx.Set(key, val)

	sctx1 := newSubContext(ctx)
	out, _ := sctx1.Get(key)
	is.Equal(val, out, "subcontext 1 should see parent's value")

	sctx2 := newSubContext(ctx)
	out, _ = sctx2.Get(key)
	is.Equal(val, out, "subcontext 2 should see parent's value")

	newVal := "baz"
	sctx1.Set(key, newVal)
	out, _ = sctx1.Get(key)
	is.Equal(newVal, out, "subcontext 1 should see new set value")

	out, _ = sctx2.Get(key)
	is.Equal(val, out, "subcontext 2 should still see parent's value")

	out, _ = ctx.Get(key)
	is.Equal(val, out, "parent context should see its original value")

	unknown := "fizz"
	sctx2.Set(unknown, "buzz")
	_, found := ctx.Get(unknown)
	is.False(found, "parent should not see new key from subcontext")
}
