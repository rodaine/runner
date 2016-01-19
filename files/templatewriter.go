package commands

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/rodaine/runner"
)

// WriteTemplate returns a FileWriterCommand that resolves a Template from data on the Context and writes the
// resulting output to the specified file.
func WriteTemplate(template Template, dataKey interface{}, destPath string) FileWriterCommand {
	renderedKey := fmt.Sprintf("RenderTemplateToFile-%d", rand.Int())
	rdr := RenderTemplate(template, dataKey, renderedKey)
	wrt := WriteFile(renderedKey, destPath)

	return &tplWriter{
		wrt: wrt,
		seq: runner.NewSequence(rdr, wrt),
	}
}

type tplWriter struct {
	wrt FileWriterCommand
	seq runner.Command
}

func (w *tplWriter) Run(ctx runner.Context, p runner.Printer) {
	w.seq.Run(ctx, p)
}

func (w *tplWriter) Rollback(ctx runner.Context, p runner.Printer) {
	w.seq.(runner.Rollbacker).Rollback(ctx, p)
}

func (w *tplWriter) DryRun(ctx runner.Context, p runner.Printer) {
	w.seq.(runner.DryRunner).DryRun(ctx, p)
}

func (w *tplWriter) SetAppend(append bool) FileWriterCommand {
	w.wrt.SetAppend(append)
	return w
}

func (w *tplWriter) SetRollback(rollback bool) FileWriterCommand {
	w.wrt.SetRollback(rollback)
	return w
}

func (w *tplWriter) SetFileMode(mode os.FileMode) FileWriterCommand {
	w.wrt.SetFileMode(mode)
	return w
}
