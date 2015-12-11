package runner

import (
	"fmt"
	"math/rand"
	"sync"
)

// MakeParallel returns a Command that executes the passed in cmds in parallel, threading the parent context into each
// independently. Parallel Commands only share Context before forking; neither errors or key-value pairs are shared
// between parallel Commands.
//
// After all parallel Commands complete execution, if any Command failed, all other successful parallel Commands are
// rolled back. A rollback initiated from a downstream command will also trigger rollbacks on each successful parallel
// command. The rolled back parallel Commands will have access to their individual forked contexts.
//
// Commands that don't satisfy the Rollbacker or DryRunner interfaces are noted and skipped during a rollback or dry
// run, respectively.
//
// This command implements the Rollbacker and DryRunner interfaces.
func MakeParallel(cmds ...Command) Command {
	return &parallel{
		id:   fmt.Sprintf("parallel%d", rand.Int()),
		cmds: cmds,
	}
}

type parallel struct {
	id   string
	cmds []Command
}

func (c *parallel) String() string {
	return fmt.Sprintf("%d Parallel Commands", len(c.cmds))
}

func (c *parallel) Run(ctx Context, p Printer) {
	sctx := c.makeSubContexts(ctx)
	ctx.Set(c.id, sctx)

	wg := sync.WaitGroup{}
	wg.Add(len(c.cmds))

	for i := range sctx {
		go c.runParallelCommand(c.cmds[i], sctx[i], p, &wg)
	}

	wg.Wait()

	var err error
	for i := range sctx {
		if err = sctx[i].Err(); err != nil {
			break
		}
	}

	if err == nil {
		return
	}

	c.Rollback(ctx, p)
	ctx.SetErr(err)
}

func (c *parallel) Rollback(ctx Context, p Printer) {
	var sctx []Context
	val, ok := ctx.Get(c.id)
	if sctx, ok = val.([]Context); !ok {
		panic("contexts for parallel tasks missing")
	}

	wg := sync.WaitGroup{}
	wg.Add(len(c.cmds))

	for i := range sctx {
		go c.rollbackParallelCommand(c.cmds[i], sctx[i], p, &wg)
	}

	wg.Wait()
}

func (c *parallel) DryRun(ctx Context, p Printer) {
	sctx := c.makeSubContexts(ctx)
	ctx.Set(c.id, sctx)

	wg := sync.WaitGroup{}
	wg.Add(len(c.cmds))

	for i := range sctx {
		go c.dryRunParallelCommand(c.cmds[i], sctx[i], p, &wg)
	}

	wg.Wait()

	var err error
	for i := range sctx {
		if err = sctx[i].Err(); err != nil {
			break
		}
	}

	ctx.SetErr(err)
}

func (c *parallel) runParallelCommand(cmd Command, ctx Context, p Printer, wg *sync.WaitGroup) {
	cmd.Run(ctx, p)
	wg.Done()
}

func (c *parallel) rollbackParallelCommand(cmd Command, ctx Context, p Printer, wg *sync.WaitGroup) {
	if rb, ok := cmd.(Rollbacker); ok && ctx.Err() == nil {
		rb.Rollback(ctx, p)
	}
	wg.Done()
}

func (c *parallel) dryRunParallelCommand(cmd Command, ctx Context, p Printer, wg *sync.WaitGroup) {
	if dr, ok := cmd.(DryRunner); ok {
		dr.DryRun(ctx, p)
	}
	wg.Done()
}

func (c *parallel) makeSubContexts(ctx Context) (sctx []Context) {
	sctx = make([]Context, len(c.cmds))
	for i := range sctx {
		sctx[i] = newSubContext(ctx)
	}
	return
}
