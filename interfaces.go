package runner

type Command interface {
	Run(*Context, Printer)
}

type Rollbacker interface {
	Rollback(*Context, Printer)
}

type DryRunner interface {
	DryRun(*Context, Printer)
}
