package runner

import (
	"fmt"
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
)

func TestFailable_Interfaces(t *testing.T) {
	t.Parallel()

	f := &failable{}
	var _ Command = f
	var _ Rollbacker = f
	var _ DryRunner = f
	var _ fmt.Stringer = f
}

func TestFailable_String(t *testing.T) {
	t.Parallel()

	is := assert.New(t)
	cmd := &MockCommand{name: "foo"}
	f := MakeFailable(cmd)

	str := fmt.Sprint(f)
	is.Contains(str, fmt.Sprint(cmd))
	is.Contains(str, "failable")
}

func TestFailable_Run_Success(t *testing.T) {
	t.Parallel()

	is := assert.New(t)
	cmd := &MockCommand{name: "foo"}
	f := MakeFailable(cmd)

	err := Run(f)
	is.NoError(err)

	is.True(cmd.ran)
	is.False(cmd.failed)
	is.False(cmd.rolledBack)
}

func TestFailable_Run_Failure(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	cmdA := &MockCommand{name: "foo", err: errors.New("")}
	cmdB := &MockCommand{name: "bar"}

	err := Run(MakeFailable(cmdA), cmdB)
	is.NoError(err)

	is.True(cmdA.ran)
	is.True(cmdA.failed)
	is.False(cmdA.rolledBack)

	is.True(cmdB.ran)
	is.False(cmdB.failed)
	is.False(cmdB.rolledBack)
}

func TestFailable_Rollback_Success(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	cmdA := &MockCommand{name: "foo"}
	cmdB := &MockCommand{name: "bar", err: errors.New("")}

	res := Run(MakeFailable(cmdA), cmdB)
	is.Equal(cmdB.err, res)

	is.True(cmdA.ran)
	is.False(cmdA.failed)
	is.True(cmdA.rolledBack)

	is.True(cmdB.ran)
	is.True(cmdB.failed)
	is.False(cmdB.rolledBack)
}

func TestFailable_Rollback_Failure(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	cmdA := &MockCommand{name: "fizz"}
	cmdB := &MockCommand{name: "foo", err: errors.New("-")}
	cmdC := &MockCommand{name: "bar", err: errors.New("=")}

	err := Run(cmdA, MakeFailable(cmdB), cmdC)
	is.Equal(cmdC.err, err)

	is.True(cmdA.ran)
	is.False(cmdA.failed)
	is.True(cmdA.rolledBack)

	is.True(cmdB.ran)
	is.True(cmdB.failed)
	is.False(cmdB.rolledBack)

	is.True(cmdC.ran)
	is.True(cmdC.failed)
	is.False(cmdC.rolledBack)
}

func TestFailable_DryRun_Success(t *testing.T) {
	t.Parallel()

	is := assert.New(t)
	cmd := &MockCommand{name: "foo"}
	DryRun(MakeFailable(cmd))
	is.True(cmd.dryRan)
}

func TestFailable_DryRun_Failure(t *testing.T) {
	t.Parallel()

	is := assert.New(t)
	cmd := &MockCommand{name: "foo", err: errors.New("bar")}

	ctx := NewContext()
	MakeFailable(cmd).(*failable).DryRun(ctx, DefaultPrinter)

	is.NoError(ctx.Err())
	is.True(cmd.dryRan)
}
