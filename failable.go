package runner

import (
	"fmt"
	"math/rand"
)

type failable struct {
	id  string
	cmd Command
}

// MakeFailable returns a Command that wraps another Command and suppresses any errors raised within. This Command will
// never trigger a rollback. Sequence Commands wrapped by MakeFailable will still rollback internally if a sub-Command
// fails, however the raised error is suppressed.
//
// If a rollback occurs, Commands wrapped by MakeFailable will only be rolled back if they did not fail internally.
//
// This command implements the Rollbacker and DryRunner interfaces.
func MakeFailable(cmd Command) Command {
	return &failable{
		id:  fmt.Sprintf("failable%d", rand.Int()),
		cmd: cmd,
	}
}

func (f *failable) String() string {
	return fmt.Sprintf("%s [failable]", f.cmd)
}

func (f *failable) Run(ctx Context, p Printer) {
	f.cmd.Run(ctx, p)
	f.suppressError(ctx, p)
}

func (f *failable) Rollback(ctx Context, p Printer) {
	if val, found := ctx.Get(f.id); found {
		if err, ok := val.(error); ok && err != nil {
			p.Warn("skipping rollback due to failure: %v", err)
			return
		}
	}

	if cmd, ok := f.cmd.(Rollbacker); ok {
		cmd.Rollback(ctx, p)
	}
}

func (f *failable) DryRun(ctx Context, p Printer) {
	if cmd, ok := f.cmd.(DryRunner); ok {
		cmd.DryRun(ctx, p)
		f.suppressError(ctx, p)
	}
}

func (f *failable) suppressError(ctx Context, p Printer) {
	err := ctx.Err()
	ctx.Set(f.id, err)

	if err != nil {
		p.Warn("failure supressed: %v", err)
		ctx.unsetErr()
	}
}
