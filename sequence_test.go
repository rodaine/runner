package runner

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSequence_Interfaces(t *testing.T) {
	seq := &sequence{}
	var _ Command = seq
	var _ Rollbacker = seq
	var _ DryRunner = seq
	var _ fmt.Stringer = seq
}

func TestSequence_String(t *testing.T) {
	is := assert.New(t)
	seq := NewSequence()
	is.Contains(fmt.Sprint(seq), "0")

	seq = NewSequence(&MockCommand{}, &MockCommand{})
	is.Contains(fmt.Sprint(seq), "2")
}

func TestSequence_Run_EmptySequence(t *testing.T) {
	is := assert.New(t)
	p, out := getTestPrinter()

	ctx := NewContext(0)
	seq := NewSequence()

	seq.Run(ctx, p)
	is.NoError(ctx.Err())
	is.Empty(out.String())
}

func TestSequence_Run_Success(t *testing.T) {
	is := assert.New(t)
	p, out := getTestPrinter()

	ctx := NewContext(3)
	seq := NewSequence(
		&MockCommand{name: "A"},
		&MockCommand{name: "B"},
		&MockCommand{name: "C"},
	)

	seq.Run(ctx, p)
	is.NoError(ctx.Err())
	is.Equal("MOCK running A\nMOCK running B\nMOCK running C\n", out.String())
}

func TestSequence_Run_Rollback(t *testing.T) {
	is := assert.New(t)
	p, out := getTestPrinter()
	err := errors.New("foobar")

	ctx := NewContext(3)
	seq := NewSequence(
		&MockCommand{name: "A"},
		&MockCommand{name: "B"},
		&MockCommand{name: "C", err: err},
		&MockCommand{name: "D"},
	)

	seq.Run(ctx, p)
	is.Equal(err, ctx.Err())
	is.Equal("MOCK running A\nMOCK running B\nMOCK running C\nMOCK error C: foobar\nMOCK rolling back B\nMOCK rolling back A\n", out.String())
}

func TestSequence_DryRun_EmptySequence(t *testing.T) {
	is := assert.New(t)
	p, out := getTestPrinter()

	ctx := NewContext(0)
	seq := NewSequence().(*sequence)

	seq.DryRun(ctx, p)
	is.NoError(ctx.Err())
	is.Empty(out.String())
}

func TestSequence_DryRun_Success(t *testing.T) {
	is := assert.New(t)
	p, out := getTestPrinter()

	ctx := NewContext(3)
	seq := NewSequence(
		&MockCommand{name: "A"},
		&MockCommand{name: "B"},
		&MockCommand{name: "C"},
	).(*sequence)

	seq.DryRun(ctx, p)
	is.NoError(ctx.Err())
	is.Equal("MOCK dry run A\nMOCK dry run B\nMOCK dry run C\n", out.String())
}

func TestSequence_Rollback_EmptySequence(t *testing.T) {
	is := assert.New(t)
	p, out := getTestPrinter()

	ctx := NewContext(0)
	seq := NewSequence().(*sequence)
	err := errors.New("foobar")

	seq.Run(ctx, p)
	out.Reset()
	ctx.SetErr(err)
	seq.Rollback(ctx, p)

	is.Equal(err, ctx.Err())
	is.Empty(out.String())
}

func TestSequence_Rollback_Success(t *testing.T) {
	is := assert.New(t)
	p, out := getTestPrinter()

	ctx := NewContext(3)
	seq := NewSequence(
		&MockCommand{name: "A"},
		&MockCommand{name: "B"},
		&MockCommand{name: "C"},
	).(*sequence)
	err := errors.New("foobar")

	seq.Run(ctx, p)
	out.Reset()
	ctx.SetErr(err)
	seq.Rollback(ctx, p)

	is.Equal(err, ctx.Err())
	is.Equal("MOCK rolling back C\nMOCK rolling back B\nMOCK rolling back A\n", out.String())
}
