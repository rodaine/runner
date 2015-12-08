package runner

import "golang.org/x/net/context"

type Command interface {
	Run(context.Context, Printer) context.Context
}

type Rollbacker interface {
	Rollback(context.Context, Printer) context.Context
}

type DryRunner interface {
	DryRun(context.Context, Printer) context.Context
}
