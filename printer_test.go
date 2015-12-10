package runner

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestPrinter() (*stdPrinter, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	return (NewPrinter(buf, LevelAll)).(*stdPrinter), buf
}

func TestPrinter_Log(t *testing.T) {
	is := assert.New(t)
	p, out := getTestPrinter()

	p.Trace("foo")
	is.Equal("foo\n", out.String())
	out.Reset()
	p.level = LevelDebug
	p.Trace("bar")
	is.Empty(out.String())
	out.Reset()

	p.Debug("fizz")
	is.Equal("fizz\n", out.String())
	out.Reset()
	p.level = LevelInfo
	p.Debug("buzz")
	is.Empty(out.String())
	out.Reset()

	p.Info("alpha")
	is.Equal("alpha\n", out.String())
	out.Reset()
	p.level = LevelWarn
	p.Info("beta")
	is.Empty(out.String())
	out.Reset()

	p.Warn("foo")
	is.Equal("foo\n", out.String())
	out.Reset()
	p.level = LevelError
	p.Warn("bar")
	is.Empty(out.String())
	out.Reset()

	p.Err("fizz")
	is.Equal("fizz\n", out.String())
	out.Reset()
	p.level = LevelFatal
	p.Err("buzz")
	is.Empty(out.String())
	out.Reset()

	p.Fatal("alpha")
	is.Equal("alpha\n", out.String())
	out.Reset()
	p.level = LevelOff
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
