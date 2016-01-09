package commands

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"io"

	"github.com/rodaine/runner"
	"github.com/stretchr/testify/assert"
)

func TestFileWriter_Run_Success(t *testing.T) {
	t.Parallel()
	key := "srcKey"
	expected := "foobar"
	sources := []interface{}{
		bytes.NewBufferString(expected), // io.Reader
		[]byte(expected),                // []byte
		expected,                        // string
	}

	is := assert.New(t)
	for _, src := range sources {
		fn := prepTempFile()
		defer cleanFile(fn)

		ctx := runner.NewContext()
		ctx.Set(key, src)

		cmd := WriteFile(key, fn)
		cmd.Run(ctx, runner.DefaultPrinter)

		b, err := ioutil.ReadFile(fn)
		is.NoError(err)
		is.Equal(expected, string(b))
	}
}

func TestFileWriter_Run_InvalidSource(t *testing.T) {
	t.Parallel()
	key := "foo"
	ctx := runner.NewContext()
	ctx.Set(key, 123)

	fn := prepTempFile()
	defer cleanFile(fn)

	WriteFile(key, fn).Run(ctx, runner.DefaultPrinter)
	assert.Error(t, ctx.Err())
}

func TestFileWriter_Run_BadReader(t *testing.T) {
	t.Parallel()
	key := "foo"
	ctx := runner.NewContext()
	ctx.Set(key, &BrokenReader{})

	fn := prepTempFile()
	defer cleanFile(fn)

	WriteFile(key, fn).Run(ctx, runner.DefaultPrinter)
	assert.Error(t, ctx.Err())
}

func TestFileWriter_Run_NoSource(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	key := "foo"

	fn := prepTempFile()
	defer cleanFile(fn)

	ctx := runner.NewContext()
	cmd := WriteFile(key, fn).(*fileWriter)
	cmd.Run(ctx, runner.DefaultPrinter)
	is.NoError(ctx.Err())

	is.NoError(ctx.Err())
	is.Equal(int64(0), cmd.bytesWritten(ctx))
}

func TestFileWriter_Run_InvalidPath(t *testing.T) {
	t.Parallel()
	key := "foo"

	fn := "/dev/null/bar"
	defer cleanFile(fn)

	ctx := runner.NewContext()
	WriteFile(key, fn).Run(ctx, runner.DefaultPrinter)
	assert.Error(t, ctx.Err())
}

func TestFileWriter_Rollback_Create(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	fn := prepTempFile()
	defer cleanFile(fn)
	ctx := runner.NewContext()

	cmd := WriteFile("", fn).SetRollback(true).(*fileWriter)
	cmd.Run(ctx, runner.DefaultPrinter)
	is.NoError(ctx.Err())

	cmd.Rollback(ctx, runner.DefaultPrinter)
	_, err := ioutil.ReadFile(fn)
	is.Error(err)

	fn = prepTempFile()
	defer cleanFile(fn)
	ctx = runner.NewContext()

	cmd = WriteFile("", fn).SetRollback(true).(*fileWriter)
	cmd.Rollback(ctx, runner.DefaultPrinter)

	_, err = ioutil.ReadFile(fn)
	is.Error(err)
}

func TestFileWriter_Rollback_Append(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	fn := prepTempFile()
	defer cleanFile(fn)

	ctx := runner.NewContext()
	ctx.Set("foo", "foo")
	ctx.Set("bar", "bar")

	cmd := runner.NewSequence(
		WriteFile("foo", fn).SetRollback(false),
		WriteFile("bar", fn).SetAppend(true).SetRollback(true),
	)
	cmd.Run(ctx, runner.DefaultPrinter)

	b, _ := ioutil.ReadFile(fn)
	is.Equal("foobar", string(b))

	rb := cmd.(runner.Rollbacker)
	rb.Rollback(ctx, runner.DefaultPrinter)

	b, _ = ioutil.ReadFile(fn)
	is.Equal("foo", string(b))
}

func TestFileWriter_Rollback_AppendNoBytes(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	fn := prepTempFile()
	defer cleanFile(fn)

	ctx := runner.NewContext()
	ctx.Set("foo", "foo")
	ctx.Set("bar", "")

	cmd := runner.NewSequence(
		WriteFile("foo", fn),
		WriteFile("bar", fn).SetAppend(true).SetRollback(true),
	)
	cmd.Run(ctx, runner.NewPrinter(os.Stdout, runner.LevelDebug))

	b, _ := ioutil.ReadFile(fn)
	is.Equal("foo", string(b))

	rb := cmd.(runner.Rollbacker)
	rb.Rollback(ctx, runner.NewPrinter(os.Stdout, runner.LevelDebug))

	b, _ = ioutil.ReadFile(fn)
	is.Equal("foo", string(b))
}

func TestFileWriter_Rollback_BadFile(t *testing.T) {
	t.Parallel()
	assert.NotPanics(t, func() {
		ctx := runner.NewContext()
		cmd := WriteFile("foo", "/dev/null/bar").SetAppend(true).SetRollback(true).(*fileWriter)
		cmd.setBytesWritten(ctx, 123)

		cmd.Rollback(ctx, runner.DefaultPrinter)
	})
}

func TestFileWriter_Rollback_ErroneousSize(t *testing.T) {
	t.Parallel()
	assert.NotPanics(t, func() {
		var f *os.File
		f, _ = ioutil.TempFile("", "")
		_ = f.Close()
		defer cleanFile(f.Name())

		ctx := runner.NewContext()

		cmd := WriteFile("foo", f.Name()).SetAppend(true).SetRollback(true).(*fileWriter)
		cmd.setBytesWritten(ctx, 123)

		cmd.Rollback(ctx, runner.DefaultPrinter)
	})
}

func TestFileWriter_DryRun_Success(t *testing.T) {
	t.Parallel()
	key := "srcKey"
	expected := "foobar"
	sources := []interface{}{
		bytes.NewBufferString(expected), // io.Reader
		[]byte(expected),                // []byte
		expected,                        // string
	}

	is := assert.New(t)
	for _, src := range sources {
		ctx := runner.NewContext()
		ctx.Set(key, src)

		cmd := WriteFile(key, "").(*fileWriter)
		cmd.DryRun(ctx, runner.DefaultPrinter)
		is.NoError(ctx.Err())
	}
}

func TestFileWriter_DryRun_InvalidSource(t *testing.T) {
	t.Parallel()
	key := "foo"
	ctx := runner.NewContext()
	ctx.Set(key, 123)

	WriteFile(key, "").(*fileWriter).DryRun(ctx, runner.DefaultPrinter)
	assert.Error(t, ctx.Err())
}

func TestFileWriter_DryRun_BadReader(t *testing.T) {
	t.Parallel()
	key := "foo"
	ctx := runner.NewContext()
	ctx.Set(key, &BrokenReader{})

	WriteFile(key, "").(*fileWriter).DryRun(ctx, runner.DefaultPrinter)
	assert.Error(t, ctx.Err())
}

func TestFileWriter_SetAppend(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := runner.NewContext()
	ctx.Set("foo", "foo")
	ctx.Set("bar", "bar")

	fn := prepTempFile()
	defer cleanFile(fn)

	runner.NewSequence(
		WriteFile("foo", fn),
		WriteFile("bar", fn).SetAppend(true),
	).Run(ctx, runner.DefaultPrinter)

	b, err := ioutil.ReadFile(fn)
	is.NoError(err)
	is.Equal("foobar", string(b))
}

func TestFileWriter_SetFileMode(t *testing.T) {
	is := assert.New(t)

	fn := prepTempFile()
	defer cleanFile(fn)

	WriteFile("", fn).Run(runner.NewContext(), runner.DefaultPrinter)

	info, _ := os.Stat(fn)
	is.Equal(DefaultFileWriterFileMode, info.Mode().Perm(), "%#o - %#o", DefaultFileWriterFileMode, info.Mode().Perm())

	fn = prepTempFile()
	defer cleanFile(fn)

	var mode os.FileMode = 0777
	WriteFile("", fn).SetFileMode(mode).Run(runner.NewContext(), runner.DefaultPrinter)

	info, _ = os.Stat(fn)
	is.Equal(mode, info.Mode().Perm(), "%#o - %#o", mode, info.Mode().Perm())
}

func TestFileWriter_bytesWritten_NeverSet(t *testing.T) {
	t.Parallel()
	cmd := WriteFile("foo", "/tmp/bar").(*fileWriter)
	assert.Equal(t, int64(0), cmd.bytesWritten(runner.NewContext()))
}

func TestFileWriter_fileSize(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	var f *os.File

	in := "foobar"
	f, _ = ioutil.TempFile("", "")
	fn := f.Name()
	_, _ = io.WriteString(f, in)

	cmd := WriteFile("", "").(*fileWriter)
	size := cmd.fileSize(f, runner.DefaultPrinter)
	is.Len([]byte(in), int(size))

	_ = f.Close()
	_ = os.Remove(fn)
	size = cmd.fileSize(f, runner.DefaultPrinter)
	is.Equal(int64(0), size)
}

func TestFileWriter_truncateFile(t *testing.T) {
	var f *os.File
	f, _ = ioutil.TempFile("", "")
	defer cleanFile(f.Name())
	n, _ := io.WriteString(f, "foobar")

	_ = f.Close()
	f, _ = os.Open(f.Name())

	assert.NotPanics(t, func() {
		WriteFile("", "").(*fileWriter).truncateFile(f, int64(n), runner.DefaultPrinter)
	})

	_ = f.Close()
}
