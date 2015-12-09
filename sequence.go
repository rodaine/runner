package runner

import "fmt"

type sequence struct {
	id   string
	cmds []Command
}

func NewSequence(cmds ...Command) Command {
	return &sequence{
		cmds: cmds,
	}
}

func (s *sequence) String() string {
	return fmt.Sprintf("%d Command Sequence", len(s.cmds))
}

func (s *sequence) Run(ctx *Context, p Printer) {
	ctx.push()
	s.runSubCommands(ctx, p, s.cmds)
}

func (s *sequence) Rollback(ctx *Context, p Printer) {
	s.rollbackSubCommands(ctx, p, s.cmds)
	ctx.pop()
}

func (s *sequence) DryRun(ctx *Context, p Printer) {
	s.dryRunSubCommands(ctx, p, s.cmds)
}

func (s *sequence) runSubCommands(ctx *Context, p Printer, cmds []Command) {
	// no Commands remain, or there is an existing error.
	if len(cmds) == 0 || ctx.Err() != nil {
		return
	}

	// run the next Command
	ctx.push()
	cmds[0].Run(ctx, p)

	// if it was the last or there was an error, exit now
	if len(cmds) == 1 || ctx.Err() != nil {
		return
	}

	// run subsequent Commands
	s.runSubCommands(ctx, p, cmds[1:])

	// if there was an error from subsequent Commands, rollback if possible
	if ctx.Err() != nil {
		ctx.pop()
		if cmd, ok := cmds[0].(Rollbacker); ok {
			cmd.Rollback(ctx, p)
		}
	}
}

func (s *sequence) rollbackSubCommands(ctx *Context, p Printer, cmds []Command) {
	if len(cmds) == 0 {
		return
	}

	ctx.pop()
	if cmd, ok := cmds[len(cmds)-1].(Rollbacker); ok {
		cmd.Rollback(ctx, p)
	}

	if len(cmds) == 1 {
		return
	}

	s.rollbackSubCommands(ctx, p, cmds[:len(cmds)-1])
}

func (s *sequence) dryRunSubCommands(ctx *Context, p Printer, cmds []Command) {
	if len(cmds) == 0 {
		return
	}

	ctx.push()
	if cmd, ok := cmds[0].(DryRunner); ok {
		cmd.DryRun(ctx, p)
	}

	if len(cmds) == 1 {
		return
	}

	s.dryRunSubCommands(ctx, p, cmds[1:])
}
