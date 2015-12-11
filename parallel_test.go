package runner

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParallel_Interfaces(t *testing.T) {
	var (
		p *parallel
		_ Command      = p
		_ Rollbacker   = p
		_ DryRunner    = p
		_ fmt.Stringer = p
	)
}

func TestParallel_String(t *testing.T) {
	is := assert.New(t)
	p := MakeParallel()
	is.Contains(fmt.Sprint(p), "0")

	p = MakeParallel(&MockCommand{}, &MockCommand{})
	is.Contains(fmt.Sprint(p), "2")
}

func TestParallel_Run_Empty(t *testing.T) {
	is := assert.New(t)

	ctx := NewContext()
	p := MakeParallel()

	p.Run(ctx, DefaultPrinter)
	is.NoError(ctx.Err())
}

func TestParallel_Run_Success(t *testing.T) {
	is := assert.New(t)

	cmdA := &MockCommand{name: "A"}
	cmdB := &MockCommand{name: "B"}
	cmdC := &MockCommand{name: "C"}

	ctx := NewContext()
	p := MakeParallel(cmdA, cmdB, cmdC)
	p.Run(ctx, DefaultPrinter)

	is.NoError(ctx.Err())
	for _, cmd := range []*MockCommand{cmdA, cmdB, cmdC} {
		is.True(cmd.ran)
		is.False(cmd.failed)
		is.False(cmd.rolledBack)
	}
}

func TestParallel_Run_Rollback(t *testing.T) {
	is := assert.New(t)
	err := errors.New("foobar")

	cmdA := &MockCommand{name: "A"}
	cmdB := &MockCommand{name: "B"}
	cmdC := &MockCommand{name: "C", err: err}
	cmdD := &MockCommand{name: "D"}

	ctx := NewContext()
	p := MakeParallel(cmdA, cmdB, cmdC, cmdD)
	p.Run(ctx, DefaultPrinter)

	is.Equal(err, ctx.Err())

	is.True(cmdC.ran)
	is.True(cmdC.failed)
	is.False(cmdC.rolledBack)

	for _, cmd := range []*MockCommand{cmdA, cmdB, cmdD} {
		is.True(cmd.ran)
		is.False(cmd.failed, "%+v", cmd)
		is.True(cmd.rolledBack, "%+v", cmd)
	}
}

func TestParallel_Rollback_Panic(t *testing.T) {
	assert.Panics(t, func() {
		MakeParallel().(Rollbacker).Rollback(NewContext(), DefaultPrinter)
	})
}

func TestParallel_Rollback_Empty(t *testing.T) {
	is := assert.New(t)

	ctx := NewContext()
	p := MakeParallel().(*parallel)
	err := errors.New("foobar")

	p.Run(ctx, DefaultPrinter)
	ctx.SetErr(err)
	p.Rollback(ctx, DefaultPrinter)

	is.Equal(err, ctx.Err())
}

func TestParallel_Rollback_Success(t *testing.T) {
	is := assert.New(t)

	cmdA := &MockCommand{name: "A"}
	cmdB := &MockCommand{name: "B"}
	cmdC := &MockCommand{name: "C"}

	ctx := NewContext()
	p := MakeParallel(cmdA, cmdB, cmdC).(*parallel)
	p.Run(ctx, DefaultPrinter)
	p.Rollback(ctx, DefaultPrinter)

	for _, cmd := range []*MockCommand{cmdA, cmdB, cmdC} {
		is.True(cmd.ran)
		is.False(cmd.failed)
		is.True(cmd.rolledBack)
	}
}

func TestParallel_DryRun_Empty(t *testing.T) {
	ctx := NewContext()
	p := MakeParallel().(*parallel)

	p.DryRun(ctx, DefaultPrinter)
	assert.NoError(t, ctx.Err())
}

func TestParallel_DryRun_Success(t *testing.T) {
	is := assert.New(t)

	cmdA := &MockCommand{name: "A"}
	cmdB := &MockCommand{name: "B"}
	cmdC := &MockCommand{name: "C"}

	ctx := NewContext()
	p := MakeParallel(cmdA, cmdB, cmdC).(*parallel)
	p.DryRun(ctx, DefaultPrinter)

	is.NoError(ctx.Err())
	for _, cmd := range []*MockCommand{cmdA, cmdB, cmdC} {
		is.True(cmd.dryRan)
		is.False(cmd.ran)
		is.False(cmd.rolledBack)
		is.False(cmd.failed)
	}
}

func TestParallel_DryRun_Fail(t *testing.T) {
	is := assert.New(t)
	err := errors.New("foo")

	cmdA := &MockCommand{name: "A"}
	cmdB := &MockCommand{name: "B"}
	cmdC := &MockCommand{name: "C", err: err}

	ctx := NewContext()
	p := MakeParallel(cmdA, cmdB, cmdC).(*parallel)
	p.DryRun(ctx, DefaultPrinter)

	is.Equal(err, ctx.Err())
	is.True(cmdC.failed)

	for _, cmd := range []*MockCommand{cmdA, cmdB, cmdC} {
		is.True(cmd.dryRan)
		is.False(cmd.ran)
		is.False(cmd.rolledBack)
	}

}

// TODO: Test a downstream Command failing + rollback into Parallel Command
