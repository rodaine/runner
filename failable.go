package runner

import (
	"fmt"
	"math/rand"
)

type failable struct {
	id  string
	cmd Command
}

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

	err := ctx.Err()
	ctx.Set(f.id, err)

	if err != nil {
		p.Warn("failure supressed: %v", err)
		ctx.SetErr(nil)
	}
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
	}
}
