package runner

import "fmt"

type MockCommand struct {
	name string
	err  error
}

func (c *MockCommand) String() string {
	return fmt.Sprintf("MOCK %s", c.name)
}

func (c *MockCommand) Run(ctx *Context, p Printer) {
	p.Info("MOCK running %s", c.name)

	if c.err != nil {
		p.Err("MOCK error %s: %v", c.name, c.err)
		ctx.SetErr(c.err)
	}
}

func (c *MockCommand) Rollback(ctx *Context, p Printer) {
	p.Info("MOCK rolling back %s", c.name)
}

func (c *MockCommand) DryRun(ctx *Context, p Printer) {
	p.Info("MOCK dry run %s", c.name)
}
