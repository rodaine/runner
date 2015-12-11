package runner

// A Command is an atomic operation to perform. Information necessary to execute the command should either be
// encapsulated by the implementation, or passed in via the Context. Any information pertaining to the execution of the
// Command should be written into the Printer at the appropriate LogLevel.
//
//Data can be read and written from the Context to share information with
// other Commands. If the Command will also implement Rollbacker, data necessary to perform the rollback should also
// be stored in the Context.
//
// A single instance of a Command should not have side-effects between separate executions with different Contexts.
// Encapsulated data should be treated as immutable (e.g., configuration), while variable data should be set within the
// Context.
type Command interface {
	Run(Context, Printer)
}

// Rollbacker can be implemented by Commands that are reversible in the event of a downstream failure. A Rollbacker
// will have access to the same Context when the Command's Run method was executed. The error that triggered the
// rollback may not be available in the Context, so its value should not be relied upon.
//
// Commands that don't implement Rollbacker will be skipped over during a rollback; they will not halt the execution.
type Rollbacker interface {
	Rollback(Context, Printer)
}

// DryRunner can be implemented by Commands to simulate their execution. A DryRunner can/should perform all read
// operations that would normally occur during a run, while write/destructive operations should be mocked. A DryRunner
// Command that produces an error on the Context will halt the execution; neither subsequent Commands will run nor
// rollback will occur.
//
// Commands that don't implement DryRunner will be skipped over during a dry run; they will not halt the execution.
type DryRunner interface {
	DryRun(Context, Printer)
}
