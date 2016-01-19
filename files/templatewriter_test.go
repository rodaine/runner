package files

import (
	"testing"

	"os"

	"io/ioutil"

	"github.com/rodaine/runner"
	"github.com/stretchr/testify/assert"
)

func TestTemplateWriter_Run(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	fn := prepTempFile()
	defer cleanFile(fn)

	ctx := runner.NewContext()
	ctx.Set(TestTemplateDataKey, TestTemplateData)

	var mode os.FileMode = 0777
	cmd := WriteTemplate(TestTemplate, TestTemplateDataKey, fn).
		SetAppend(true).
		SetFileMode(mode)

	// run it twice to test the append
	cmd.Run(ctx, runner.DefaultPrinter)
	cmd.Run(ctx, runner.DefaultPrinter)
	expected := "foobarfoobar"

	b, err := ioutil.ReadFile(fn)
	is.NoError(ctx.Err())
	is.NoError(err)
	is.Equal(expected, string(b))

	info, _ := os.Stat(fn)
	is.Equal(mode, info.Mode().Perm())
}

func TestTemplateWriter_DryRun(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	fn := prepTempFile()
	defer cleanFile(fn)

	ctx := runner.NewContext()
	ctx.Set(TestTemplateDataKey, TestTemplateData)

	cmd := WriteTemplate(TestTemplate, TestTemplateDataKey, fn).(runner.DryRunner)
	cmd.DryRun(ctx, runner.DefaultPrinter)

	is.NoError(ctx.Err())
	_, err := os.Stat(fn)
	is.True(os.IsNotExist(err))
}

func TestTemplateWriter_Rollback(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	fn := prepTempFile()

	ctx := runner.NewContext()
	ctx.Set(TestTemplateDataKey, TestTemplateData)

	cmd := WriteTemplate(TestTemplate, TestTemplateDataKey, fn)
	cmd.Run(ctx, runner.DefaultPrinter)

	is.NoError(ctx.Err())
	_, err := os.Stat(fn)
	is.NoError(err)

	cmd.(runner.Rollbacker).Rollback(ctx, runner.DefaultPrinter)

	is.NoError(ctx.Err())
	_, err = os.Stat(fn)
	is.NoError(err)

	cleanFile(fn)
	fn = prepTempFile()
	defer cleanFile(fn)

	cmd = WriteTemplate(TestTemplate, TestTemplateDataKey, fn).SetRollback(true)
	cmd.Run(ctx, runner.DefaultPrinter)

	is.NoError(ctx.Err())
	_, err = os.Stat(fn)
	is.NoError(err)

	cmd.(runner.Rollbacker).Rollback(ctx, runner.DefaultPrinter)

	is.NoError(ctx.Err())
	_, err = os.Stat(fn)
	is.True(os.IsNotExist(err))
}
