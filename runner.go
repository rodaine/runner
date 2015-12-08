package runner

import "golang.org/x/net/context"

func Run(ctx context.Context, cmds ...Command) error {
	return RunWithPrinter(ctx, DefaultPrinter, cmds...)
}

func RunWithPrinter(ctx context.Context, p Printer, cmds ...Command) error {
	seq := NewSequence(cmds...)
	return seq.Run(ctx, p).Err()
}
