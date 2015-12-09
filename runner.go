package runner

func Run(cmds ...Command) error {
	return RunWithPrinter(DefaultPrinter, cmds...)
}

func RunWithPrinter(p Printer, cmds ...Command) error {
	ctx := NewContext(0)
	NewSequence(cmds...).Run(ctx, p)
	return ctx.Err()
}
