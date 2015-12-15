package runner

import "fmt"

type MockCommand struct {
	name string
	err  error

	set string
	see string

	setVal     bool
	seenVal    bool
	ran        bool
	failed     bool
	rolledBack bool
	dryRan     bool
}

func (c *MockCommand) String() string {
	return fmt.Sprintf("MOCK %s", c.name)
}

func (c *MockCommand) Run(ctx Context, p Printer) {
	p.Info("MOCK running %s", c.name)
	c.ran = true
	c.maybeSetVal(ctx, p)
	c.maybeSeeVal(ctx, p)
	c.maybeFail(ctx, p)
}

func (c *MockCommand) Rollback(ctx Context, p Printer) {
	p.Info("MOCK rolling back %s", c.name)
	c.rolledBack = true
	c.maybeSeeVal(ctx, p)
}

func (c *MockCommand) DryRun(ctx Context, p Printer) {
	p.Info("MOCK dry run %s", c.name)
	c.dryRan = true
	c.maybeSetVal(ctx, p)
	c.maybeSeeVal(ctx, p)
	c.maybeFail(ctx, p)
}

func (c *MockCommand) maybeFail(ctx Context, p Printer) {
	if c.err != nil {
		p.Err("MOCK error %s: %v", c.name, c.err)
		c.failed = true
		ctx.SetErr(c.err)
	}
}

func (c *MockCommand) maybeSetVal(ctx Context, p Printer) {
	if c.set != "" {
		ctx.Set(c.set, c.name)
		c.setVal = true
	}
}

func (c *MockCommand) maybeSeeVal(ctx Context, p Printer) {
	if c.see != "" {
		_, c.seenVal = ctx.Get(c.see)
		if !c.seenVal {
			p.Warn("Didn't see: %s", c.see)
		}
	}
}
