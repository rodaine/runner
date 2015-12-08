package runner

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestPrinter() (*printer, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	return (NewPrinter(buf, PRIORITY_ALL)).(*printer), buf
}

func TestPrinter_Log(t *testing.T) {
	is := assert.New(t)
	p, out := getTestPrinter()

	p.Trace("foo")
	is.Equal("foo\n", out.String())
	out.Reset()
	p.priority = PRIORITY_DEBUG
	p.Trace("bar")
	is.Empty(out.String())
	out.Reset()

	p.Debug("fizz")
	is.Equal("fizz\n", out.String())
	out.Reset()
	p.priority = PRIORITY_INFO
	p.Debug("buzz")
	is.Empty(out.String())
	out.Reset()

	p.Info("alpha")
	is.Equal("alpha\n", out.String())
	out.Reset()
	p.priority = PRIORITY_WARN
	p.Info("beta")
	is.Empty(out.String())
	out.Reset()

	p.Warn("foo")
	is.Equal("foo\n", out.String())
	out.Reset()
	p.priority = PRIORITY_ERROR
	p.Warn("bar")
	is.Empty(out.String())
	out.Reset()

	p.Error("fizz")
	is.Equal("fizz\n", out.String())
	out.Reset()
	p.priority = PRIORITY_FATAL
	p.Error("buzz")
	is.Empty(out.String())
	out.Reset()

	p.Fatal("alpha")
	is.Equal("alpha\n", out.String())
	out.Reset()
	p.priority = PRIORITY_OFF
	p.Fatal("beta")
	is.Empty(out.String())
	out.Reset()
}

func TestPrinter_Format(t *testing.T) {
	p, out := getTestPrinter()
	p.Info("foo%s", "bar")
	assert.Equal(t, "foobar\n", out.String())
}

func TestPrinter_WithPrefix(t *testing.T) {
	p, out := getTestPrinter()
	prefixed := p.WithPrefix("foo")
	prefixed.Info("bar")
	assert.Equal(t, "foobar\n", out.String())
}
