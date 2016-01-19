package files

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/rodaine/runner"
	"github.com/stretchr/testify/assert"
)

var (
	outKey = "out"
)

func TestRenderTemplate_Run_Success(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := runner.NewContext()
	ctx.Set(TestTemplateDataKey, TestTemplateData)

	RenderTemplate(TestTemplate, TestTemplateDataKey, outKey).Run(ctx, runner.DefaultPrinter)

	is.NoError(ctx.Err())
	out, found := fetchRenderedTemplate(ctx)
	is.True(found)
	is.Equal(TestTemplateExpected, renderedString(out))
}

func TestRenderTemplate_Run_Error(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := runner.NewContext()

	err := errors.New("foobar")
	RenderTemplate(&ErrorTemplate{err}, TestTemplateDataKey, outKey).Run(ctx, runner.DefaultPrinter)

	is.Equal(err, ctx.Err())
	out, found := fetchRenderedTemplate(ctx)
	is.False(found)
	is.Nil(out)
}

func TestRenderTemplate_DryRun_Success(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := runner.NewContext()
	ctx.Set(TestTemplateDataKey, TestTemplateData)

	RenderTemplate(TestTemplate, TestTemplateDataKey, outKey).(*tpl).DryRun(ctx, runner.DefaultPrinter)

	is.NoError(ctx.Err())
	out, found := fetchRenderedTemplate(ctx)
	is.True(found)
	is.Equal(TestTemplateExpected, renderedString(out))
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
