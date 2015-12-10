package runner

func Run(cmds ...Command) error {
	return RunWithPrinter(DefaultPrinter, cmds...)
}

func RunWithPrinter(p Printer, cmds ...Command) error {
	// TODO: estimate depth
	ctx := NewContext()
	(&sequence{cmds: cmds}).Run(ctx, p)
	return ctx.Err()
}

func DryRun(cmds ...Command) {
	DryRunWithPrinter(DefaultPrinter, cmds...)
}

func DryRunWithPrinter(p Printer, cmds ...Command) {
	// TODO: estimate depth
	ctx := NewContext()
	(&sequence{cmds: cmds}).DryRun(ctx, p)
}
