package runner

import (
	"fmt"
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
)

func TestFailable_Interfaces(t *testing.T) {
	f := &failable{}
	var _ Command = f
	var _ Rollbacker = f
	var _ DryRunner = f
	var _ fmt.Stringer = f
}

func TestFailable_String(t *testing.T) {
	is := assert.New(t)
	cmd := &MockCommand{name: "foo"}
	f := MakeFailable(cmd)

	str := fmt.Sprint(f)
	is.Contains(str, fmt.Sprint(cmd))
	is.Contains(str, "failable")
}

func TestFailable_Run_Success(t *testing.T) {
	is := assert.New(t)
	p, out := getTestPrinter()
	f := MakeFailable(&MockCommand{name: "foo"})

	err := RunWithPrinter(p, f)
	is.NoError(err)
	is.Equal("MOCK running foo\n", out.String())
}

func TestFailable_Run_Failure(t *testing.T) {
	is := assert.New(t)
	p, out := getTestPrinter()
	err := errors.New("fizzbuzz")
	f := MakeFailable(&MockCommand{name: "foo", err: err})

	res := RunWithPrinter(p, f, &MockCommand{name: "bar"})
	is.NoError(res)
	is.Equal("MOCK running foo\nMOCK error foo: fizzbuzz\nfailure supressed: fizzbuzz\nMOCK running bar\n", out.String())
}

func TestFailable_Rollback_Success(t *testing.T) {
	is := assert.New(t)
	p, out := getTestPrinter()
	err := errors.New("fizzbuzz")

	res := RunWithPrinter(
		p,
		MakeFailable(&MockCommand{name: "foo"}),
		&MockCommand{name: "bar", err: err},
	)
	is.Equal(err, res)
	is.Equal("MOCK running foo\nMOCK running bar\nMOCK error bar: fizzbuzz\nMOCK rolling back foo\n", out.String())
}

func TestFailable_Rollback_Failure(t *testing.T) {
	is := assert.New(t)
	p, out := getTestPrinter()
	err := errors.New("fizzbuzz")

	res := RunWithPrinter(
		p,
		MakeFailable(&MockCommand{name: "foo", err: errors.New("rawr")}),
		&MockCommand{name: "bar", err: err},
	)
	is.Equal(err, res)
	is.Equal("MOCK running foo\nMOCK error foo: rawr\nfailure supressed: rawr\nMOCK running bar\nMOCK error bar: fizzbuzz\nskipping rollback due to failure: rawr\n", out.String())
}

func TestFailable_DryRun(t *testing.T) {
	is := assert.New(t)
	p, out := getTestPrinter()

	DryRunWithPrinter(p, MakeFailable(&MockCommand{name: "foo"}))
	is.Equal("MOCK dry run foo\n", out.String())
}
