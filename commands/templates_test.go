package commands

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"text/template"

	"github.com/rodaine/runner"
	"github.com/stretchr/testify/assert"
)

var (
	txt = template.Must(template.New("txt").Parse(`{{.Field}}bar`))

	data   = struct{ Field string }{"foo"}
	key    = "tplData"
	outKey = "out"

	expected = "foobar"
)

func TestRenderTemplate_Run_Success(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := runner.NewContext()
	ctx.Set(key, data)

	RenderTemplate(txt, key, outKey).Run(ctx, runner.DefaultPrinter)

	is.NoError(ctx.Err())
	out, found := fetchRenderedTemplate(ctx)
	is.True(found)
	is.Equal(expected, renderedString(out))
}

func TestRenderTemplate_Run_Error(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := runner.NewContext()

	err := errors.New("foobar")
	RenderTemplate(&ErrorTemplate{err}, key, outKey).Run(ctx, runner.DefaultPrinter)

	is.Equal(err, ctx.Err())
	out, found := fetchRenderedTemplate(ctx)
	is.False(found)
	is.Nil(out)
}

func TestRenderTemplate_DryRun_Success(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := runner.NewContext()
	ctx.Set(key, data)

	RenderTemplate(txt, key, outKey).(*tpl).DryRun(ctx, runner.DefaultPrinter)

	is.NoError(ctx.Err())
	out, found := fetchRenderedTemplate(ctx)
	is.True(found)
	is.Equal(expected, renderedString(out))
}

func renderedString(r io.Reader) string {
	b := bytes.Buffer{}
	_, _ = b.ReadFrom(r)
	return b.String()
}

type ErrorTemplate struct {
	error
}

func (t *ErrorTemplate) Execute(wr io.Writer, data interface{}) (err error) {
	return t.error
}

func fetchRenderedTemplate(ctx runner.Context) (r io.Reader, found bool) {
	if val, ok := ctx.Get(outKey); ok {
		r, found = val.(io.Reader)
	}
	return
}
