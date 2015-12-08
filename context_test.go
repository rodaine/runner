package runner

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState_Error(t *testing.T) {
	is := assert.New(t)
	s := NewContext(0)
	is.NoError(s.Err())

	err := errors.New("test")
	s.SetErr(err)
	is.Equal(err, s.Err())
}

func TestState_GetSet(t *testing.T) {
	is := assert.New(t)
	s := NewContext(2)
	s.push()

	// should return nil if not set
	actual := s.Get("foo")
	is.Nil(actual)

	// should have no issue pulling from current context
	s.Set("foo", "bar")
	actual, ok := s.Get("foo").(string)
	is.True(ok)
	is.Equal("bar", actual)

	// should dig deeper if not on current context
	s.push()
	actual, ok = s.Get("foo").(string)
	is.True(ok)
	is.Equal("bar", actual)

	// should be overwritten by more-recent stack
	s.Set("foo", "baz")
	actual, ok = s.Get("foo").(string)
	is.True(ok)
	is.Equal("baz", actual)

	// previous value should still be available after pop
	s.pop()
	actual, ok = s.Get("foo").(string)
	is.True(ok)
	is.Equal("bar", actual)
}
