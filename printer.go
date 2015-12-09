package runner

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	startCommandFormat = "=== RUN %s"
	endCommandFormat   = "--- END %s (%s)"

	SubCommandPrefix = "   "
)

const (
	PRIORITY_ALL Priority = iota
	PRIORITY_TRACE
	PRIORITY_DEBUG
	PRIORITY_INFO
	PRIORITY_WARN
	PRIORITY_ERROR
	PRIORITY_FATAL
	PRIORITY_OFF
)

type Priority int8

type Printer interface {
	Log(lvl Priority, format string, values ...interface{})
	Trace(format string, values ...interface{})
	Debug(format string, values ...interface{})
	Info(format string, values ...interface{})
	Warn(format string, values ...interface{})
	Err(format string, values ...interface{})
	Fatal(format string, values ...interface{})

	WithPrefix(prefix string) Printer

	StartCommand(cmd Command) time.Time
	EndCommand(cmd Command, start time.Time)
}

type printer struct {
	parent   Printer
	w        io.Writer
	priority Priority
	prefix   string
}

func NewPrinter(w io.Writer, priority Priority) Printer {
	return &printer{
		w:        w,
		priority: priority,
	}
}

func (p *printer) Log(lvl Priority, format string, values ...interface{}) {
	if p.parent != nil {
		p.parent.Log(lvl, p.prefix+format, values...)
		return
	}

	if lvl < p.priority {
		return
	}

	if len(values) == 0 {
		fmt.Fprintln(p.w, format)
	} else {
		fmt.Fprintf(p.w, format+"\n", values...)
	}
}

func (p *printer) Trace(format string, values ...interface{}) {
	p.Log(PRIORITY_TRACE, format, values...)
}

func (p *printer) Debug(format string, values ...interface{}) {
	p.Log(PRIORITY_DEBUG, format, values...)
}

func (p *printer) Info(format string, values ...interface{}) {
	p.Log(PRIORITY_INFO, format, values...)
}

func (p *printer) Warn(format string, values ...interface{}) {
	p.Log(PRIORITY_WARN, format, values...)
}

func (p *printer) Err(format string, values ...interface{}) {
	p.Log(PRIORITY_ERROR, format, values...)
}

func (p *printer) Fatal(format string, values ...interface{}) {
	p.Log(PRIORITY_FATAL, format, values...)
}

func (p *printer) WithPrefix(prefix string) Printer {
	return &printer{
		parent: p,
		prefix: prefix,
	}
}

func (p *printer) StartCommand(cmd Command) time.Time {
	p.Info(startCommandFormat, cmd)
	return time.Now()
}

func (p *printer) EndCommand(cmd Command, start time.Time) {
	p.Info(endCommandFormat, cmd, time.Since(start))
}

var DefaultPrinter = NewPrinter(os.Stdout, PRIORITY_INFO)
