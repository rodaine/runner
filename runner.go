package runner

// Run executes the passed in Commands in sequence, returning an error if the execution failed and a rollback occurred.
// The DefaultPrinter is passed to all commands for logging.
func Run(cmds ...Command) error {
	return RunWithPrinter(DefaultPrinter, cmds...)
}

// RunWithPrinter executes the passed in Commands in sequence, returning an error if the execution failed and a rollback
// occurred. The provided Printer is passed to all commands for logging.
func RunWithPrinter(p Printer, cmds ...Command) error {
	// TODO: estimate depth
	ctx := NewContext()
	(&sequence{cmds: cmds}).Run(ctx, p)
	return ctx.Err()
}

// DryRun simulates a Run of the passed in Commands, without write/destructive actions. The DefaultPrinter is passed to
// all commands for logging.
func DryRun(cmds ...Command) {
	DryRunWithPrinter(DefaultPrinter, cmds...)
}

// DryRunWithPrinter simulates a Run of the passed in Commands, without write/destructive actions. The provided Printer
// is passed to all commands for logging.
func DryRunWithPrinter(p Printer, cmds ...Command) {
	// TODO: estimate depth
	ctx := NewContext()
	(&sequence{cmds: cmds}).DryRun(ctx, p)
}
