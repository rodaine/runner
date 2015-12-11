package runner

import "fmt"

type MockCommand struct {
	name string
	err  error

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
	c.maybeFail(ctx, p)
}

func (c *MockCommand) Rollback(ctx Context, p Printer) {
	p.Info("MOCK rolling back %s", c.name)
	c.rolledBack = true
}

func (c *MockCommand) DryRun(ctx Context, p Printer) {
	p.Info("MOCK dry run %s", c.name)
	c.dryRan = true
	c.maybeFail(ctx, p)
}

func (c *MockCommand) maybeFail(ctx Context, p Printer) {
	if c.err != nil {
		p.Err("MOCK error %s: %v", c.name, c.err)
		c.failed = true
		ctx.SetErr(c.err)
	}
}
