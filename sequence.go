package runner

import (
	"fmt"

	"golang.org/x/net/context"
)

type sequence struct {
	cmds []Command
}

func NewSequence(cmds ...Command) Command {
	return &sequence{cmds: cmds}
}

func (s *sequence) String() string {
	return fmt.Sprintf("%d Command Sequence", len(s.cmds))
}

func (s *sequence) Run(ctx context.Context, p Printer) context.Context {

	var idx int
	for idx = 0; idx < len(s.cmds); idx++ {
		select {
		case <-ctx.Done():
			return s.rollbackFromIndex(ctx, p, idx-1)
		default:
			cmd := s.cmds[idx]
			ctx = cmd.Run(ctx, p)
		}
	}

	select {
	case <-ctx.Done():
		return s.Rollback(ctx, p)
	default:
		return ctx
	}
}

func (s *sequence) Rollback(ctx context.Context, p Printer) context.Context {
	return s.rollbackFromIndex(ctx, p, len(s.cmds)-1)
}

func (s *sequence) DryRun(ctx context.Context, p Printer) context.Context {
	// TODO
	return ctx
}

func (s *sequence) rollbackFromIndex(ctx context.Context, p Printer, idx int) context.Context {
	for ; idx >= 0; idx-- {
		if cmd, ok := s.cmds[idx].(Rollbacker); ok {
			ctx = cmd.Rollback(ctx, p)
		}
	}
	return ctx
}
