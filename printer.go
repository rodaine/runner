package runner

import (
	"fmt"
	"io"
	"os"
)

// A LogLevel value describes the minimum logging verbosity for a Printer to output messages passed to it. Messages with
// a lower LogLevel should be suppressed by the Printer implementation.
type LogLevel int8

// LogLevel constants are provided in descending order of verbosity. These constants (excluding LevelAll and LevelOff)
// correspond with the similarly named methods on the Printer.
const (
	LevelAll LogLevel = iota
	LevelTrace
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelOff
)

// A Printer is passed into every command and should be used exclusively for logging the behavior of a command. Direct
// use of the fmt or log packages is for maintaining the cleanliness of the output. All logging methods should emulate
// fmt.Printf interpolation.
type Printer interface {
	// Log writes an arbitrary message at the given LogLevel. Commands should not call this method directly, and instead
	// use the other exposed logging methods. This method is exposed for Printer implementations that compose with other
	// printers for more direct access to their logging logic.
	Log(level LogLevel, format string, values ...interface{})

	// Trace should be used for high-noise debug information.
	Trace(format string, values ...interface{})

	// Debug should be used for detailed information on the flow through the command.
	Debug(format string, values ...interface{})

	// Info should be used for interesting runtime events. This is the minimum level for the DefaultPrinter, so be
	// conservative and keep to a minimum
	Info(format string, values ...interface{})

	// Warn should be used for deprecations, incorrect use of a command, "almost" errors, other situations that are
	// undesirable or unexpected, but not necessarily "wrong". Presence of warn messages may not predicate a rollback.
	Warn(format string, values ...interface{})

	// Error should for runtime errors or unexpected conditions. Presence of error messages should precede a rollback.
	Err(format string, values ...interface{})

	// Fatal should be used for severe errors that should result in termination of all commands, foregoing even a
	// rollback. It is expected that the program will either panic or exit immediately after these logs.
	Fatal(format string, values ...interface{})

	// WithPrefix should return a new printer that prefixes all messages with the provided string.
	WithPrefix(prefix string) Printer
}

// NewPrinter returns a standard Printer which writes all logs to the provided io.Writer. LogLevels below the provided
// level are suppressed from output.
func NewPrinter(w io.Writer, level LogLevel) Printer {
	return &stdPrinter{
		w:     w,
		level: level,
	}
}

type stdPrinter struct {
	parent Printer
	w      io.Writer
	level  LogLevel
	prefix string
}

func (p *stdPrinter) Log(level LogLevel, format string, values ...interface{}) {
	if p.parent != nil {
		p.parent.Log(level, p.prefix+format, values...)
		return
	}

	if level < p.level {
		return
	}

	if len(values) == 0 {
		fmt.Fprintln(p.w, format)
	} else {
		fmt.Fprintf(p.w, format+"\n", values...)
	}
}

func (p *stdPrinter) Trace(format string, values ...interface{}) {
	p.Log(LevelTrace, format, values...)
}

func (p *stdPrinter) Debug(format string, values ...interface{}) {
	p.Log(LevelDebug, format, values...)
}

func (p *stdPrinter) Info(format string, values ...interface{}) {
	p.Log(LevelInfo, format, values...)
}

func (p *stdPrinter) Warn(format string, values ...interface{}) {
	p.Log(LevelWarn, format, values...)
}

func (p *stdPrinter) Err(format string, values ...interface{}) {
	p.Log(LevelError, format, values...)
}

func (p *stdPrinter) Fatal(format string, values ...interface{}) {
	p.Log(LevelFatal, format, values...)
}

func (p *stdPrinter) WithPrefix(prefix string) Printer {
	return &stdPrinter{
		parent: p,
		prefix: prefix,
	}
}

// DefaultPrinter writes to os.stdOut at the Info LogLevel. It is the printer used by Run and DryRun.
var DefaultPrinter = NewPrinter(os.Stdout, LevelInfo)
