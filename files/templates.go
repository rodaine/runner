package files

import (
	"bytes"
	"io"

	"github.com/rodaine/runner"
)

// Template describes the interface used to format data and output to the provided io.Writer. Both text.Template and
// html.Template satisfy this interface.
type Template interface {
	Execute(wr io.Writer, data interface{}) (err error)
}

// RenderTemplate returns a Command that applies arbitrary data at dataKey in the Context to the provided Template and
// makes the resulting io.Reader available on the Context at outputKey.
//
// This Command implements DryRunner and performs identical operations as Run; Rollbacker is not implemented.
func RenderTemplate(template Template, dataKey, outputKey interface{}) runner.Command {
	return &tpl{
		template:  template,
		dataKey:   dataKey,
		outputKey: outputKey,
	}
}

type tpl struct {
	template           Template
	dataKey, outputKey interface{}
}

func (t *tpl) Run(ctx runner.Context, p runner.Printer) {
	data, found := ctx.Get(t.dataKey)
	if !found {
		p.Warn("template data not found")
	} else {
		p.Debug("template data: %+v", data)
	}

	wr := &bytes.Buffer{}
	if err := t.template.Execute(wr, data); err != nil {
		p.Err("unable to execute template: %v", err)
		ctx.SetErr(err)
		return
	}

	ctx.Set(t.outputKey, wr)
}

func (t *tpl) DryRun(ctx runner.Context, p runner.Printer) {
	t.Run(ctx, p)
}
