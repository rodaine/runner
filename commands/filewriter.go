package commands

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"syscall"

	"github.com/rodaine/runner"
)

const (
	// DefaultFileWriterAppend describes the default append behavior for a FileWriterCommand.
	DefaultFileWriterAppend = false

	// DefaultFileWriterRollback specifies whether or not a rollback should be performed by a FileWriterCommand by default.
	DefaultFileWriterRollback = false

	// DefaultFileWriterFileMode specifies the default FileMode for opening/creating a file with FileWriterCommand.
	DefaultFileWriterFileMode os.FileMode = 0666 // -rw-rw-rw-

	dirFileMode os.FileMode = 0755
)

// A FileWriterCommand describes the optional configuration methods available on the Command returned by WriteFile.
// These methods mutate the underlying Command; the Command is passed through to provided for chaining convenience.
type FileWriterCommand interface {
	runner.Command

	// SetAppend specifies if the source data should be appended to the destination file (true) or if the file should be
	// truncated before writing (false). The default value is specified by DefaultFileWriterAppend.
	SetAppend(append bool) FileWriterCommand

	// SetRollback specifies if this command should perform a rollback operation if the write partially fails or if a
	// subsequent Command fails. If true, the written bytes will be truncated from the file, restoring it to its
	// previous state. If false, the file will be deleted from the file system. The default value is specified
	// by DefaultFileWriterRollback.
	SetRollback(rollback bool) FileWriterCommand

	// SetFileMode specifies the FileMode to create/open the destination file. The default value is specified by
	// DefaultFileWriterFileMode. When creating a file, this Command attempts to ignore the umask of the process.
	SetFileMode(mode os.FileMode) FileWriterCommand
}

// WriteFile returns a FileWriterCommand that takes source data and writes it to a destination file. Valid source data
// for this Command includes io.Reader, string, and []byte; other types will result in an error and rollback.
//
// This Command implements Rollbacker, and can be enabled/disabled via SetRollback. The rollback behavior depends on the
// value passed to SetAppend. If true, the written bytes will be truncated from the file, restoring it to its previous
// state. If false, the file will be deleted from the file system. The default value is specified by
// DefaultFileWriterRollback.
//
// This Command also implements DryRunner, however no file will be written to the file system.
func WriteFile(sourceKey interface{}, destPath string) FileWriterCommand {
	return &fileWriter{
		sourceKey: sourceKey,
		destPath:  destPath,
		append:    DefaultFileWriterAppend,
		rollback:  DefaultFileWriterRollback,
		mode:      DefaultFileWriterFileMode,
	}
}

type fileWriter struct {
	sourceKey        interface{}
	destPath         string
	append, rollback bool
	mode             os.FileMode
}

func (w *fileWriter) Run(ctx runner.Context, p runner.Printer) {
	src, err := w.getSource(ctx, p)
	if err != nil {
		p.Err("%v", err)
		ctx.SetErr(err)
		return
	}

	f, err := w.openFile()
	if err != nil {
		p.Err("unable to access file: %v", err)
		ctx.SetErr(err)
		return
	}

	n, err := io.Copy(f, src)
	p.Debug("bytes written: %d", n)
	w.setBytesWritten(ctx, n)

	_ = f.Close()
	if err != nil {
		p.Err("unable to write to file: %v", err)
		ctx.SetErr(err)
		w.Rollback(ctx, p)
	}
}

func (w *fileWriter) Rollback(ctx runner.Context, p runner.Printer) {
	if !w.rollback {
		p.Debug("rollback disabled for this command")
		return
	}

	if !w.append {
		if err := os.Remove(w.destPath); err != nil {
			p.Err("could not remove file: %v", err)
		}
		return
	}

	n := w.bytesWritten(ctx)
	if n == 0 {
		return
	}

	f, err := w.openFile()
	if err != nil {
		p.Err("could not open file: %v", err)
		return
	}
	defer func() { _ = f.Close() }()

	w.truncateFile(f, n, p)
}

func (w *fileWriter) DryRun(ctx runner.Context, p runner.Printer) {
	src, err := w.getSource(ctx, p)
	if err != nil {
		p.Err("%v", err)
		ctx.SetErr(err)
		return
	}

	n, err := io.Copy(ioutil.Discard, src)
	p.Debug("bytes written: %d", n)
	w.setBytesWritten(ctx, n)

	if err != nil {
		p.Err("unable to read file: %v", err)
		ctx.SetErr(err)
		return
	}
}

func (w *fileWriter) SetAppend(append bool) FileWriterCommand {
	w.append = append
	return w
}

func (w *fileWriter) SetRollback(rollback bool) FileWriterCommand {
	w.rollback = rollback
	return w
}

func (w *fileWriter) SetFileMode(mode os.FileMode) FileWriterCommand {
	w.mode = mode
	return w
}

func (w *fileWriter) getSource(ctx runner.Context, p runner.Printer) (r io.Reader, err error) {
	data, found := ctx.Get(w.sourceKey)
	if !found {
		p.Warn("source data not found")
		return &emptyReader{}, nil
	}

	switch data.(type) {
	case io.Reader:
		return data.(io.Reader), nil
	case string:
		return bytes.NewBufferString(data.(string)), nil
	case []byte:
		return bytes.NewBuffer(data.([]byte)), nil
	}

	err = fmt.Errorf("unsupported source data type: %T", data)
	return &emptyReader{}, err
}

func (w *fileWriter) openFile() (f *os.File, err error) {
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)

	if err = os.MkdirAll(filepath.Dir(w.destPath), dirFileMode); err != nil {
		return nil, err
	}

	flag := os.O_CREATE | os.O_WRONLY
	if w.append {
		flag = flag | os.O_APPEND
	} else {
		flag = flag | os.O_TRUNC
	}

	f, err = os.OpenFile(w.destPath, flag, w.mode)
	return
}

func (w *fileWriter) truncateFile(f *os.File, n int64, p runner.Printer) {
	size := w.fileSize(f, p)
	if size < n {
		p.Err("erroneous truncate size: %d", n)
		return
	}

	if err := f.Truncate(size - n); err != nil {
		p.Err("could not truncate file: %v", err)
		return
	}
}

func (w *fileWriter) fileSize(f *os.File, p runner.Printer) int64 {
	stats, err := f.Stat()
	if err != nil {
		p.Err("could not stat file: %v", err)
		return 0
	}
	return stats.Size()
}

func (w *fileWriter) setBytesWritten(ctx runner.Context, n int64) {
	ctx.Set(w.bytesWrittenKey(), n)
}

func (w *fileWriter) bytesWritten(ctx runner.Context) (n int64) {
	if val, ok := ctx.Get(w.bytesWrittenKey()); ok {
		if n, ok = val.(int64); ok {
			return
		}
	}
	return 0
}

func (w *fileWriter) bytesWrittenKey() interface{} {
	return fmt.Sprintf("bytes written - %s", w.destPath)
}

type emptyReader struct{}

func (r *emptyReader) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}
