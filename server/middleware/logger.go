package middleware

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

var (
	LogEntryCtxKey = &contextKey{"LogEntry"}
	DefaultLogger  func(next http.Handler) http.Handler
)

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return. When standard output is a TTY, Logger will
// print in color, otherwise it will print in black and white. Logger prints a
// request id if one is provided.
func Logger(next http.Handler) http.Handler {
	return DefaultLogger(next)
}

// RequestLogger returns a logger handler using a custom LogFormatter
func RequestLogger(lf LogFormatter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := lf.NewLogEntry(r)
			wrw := NewWrappedResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.Write(wrw.Status(), wrw.BytesWritten(), wrw.Header(), time.Since(t1), nil)
			}()

			next.ServeHTTP(wrw, WithLogEntry(r, entry))
		}
		return http.HandlerFunc(fn)
	}
}

// LogFormatter initiates the begining of a new log entry per request
type LogFormatter interface {
	NewLogEntry(r *http.Request) LogEntry
}

// LogEntry records the final log when the request completes
type LogEntry interface {
	Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{})
	Panic(v interface{}, stack []byte)
}

// GetLogEntry returns the in-context LogEntry for a request
func GetLogEntry(r *http.Request) LogEntry {
	entry, _ := r.Context().Value(LogEntryCtxKey).(LogEntry)
	return entry
}

// WithLogEntry sets the in-context LogEntry for a request
func WithLogEntry(r *http.Request, entry LogEntry) *http.Request {
	r = r.WithContext(context.WithValue(r.Context(), LogEntryCtxKey, entry))
	return r
}

// LoggerInterface accepts printing to stdlib logger or compatible logger
type LoggerInterface interface {
	Print(v ...interface{})
}

// DefaultLogFormatter is a simple logger that implements LogFormatter
type DefaultLogFormatter struct {
	Logger  LoggerInterface
	NoColor bool
}

func (dlf *DefaultLogFormatter) NewLogEntry(r *http.Request) LogEntry {
	useColor := !dlf.NoColor
	entry := &defaultLogEntry{
		DefaultLogFormatter: dlf,
		request:             r,
		buf:                 &bytes.Buffer{},
		useColor:            useColor,
	}

	// reqId := GetReqId(r.Context())

	cw(entry.buf, useColor, nCyan, "\"")
	cw(entry.buf, useColor, bMagenta, "%s ", r.Method)

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	cw(entry.buf, useColor, nCyan, "%s://%s%s %s\"", scheme, r.Host, r.RequestURI, r.Proto)

	entry.buf.WriteString(" from ")
	entry.buf.WriteString(r.RemoteAddr)
	entry.buf.WriteString(" - ")

	return entry
}

type defaultLogEntry struct {
	*DefaultLogFormatter
	request  *http.Request
	buf      *bytes.Buffer
	useColor bool
}

func (dle *defaultLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	switch {
	case status < 200:
		cw(dle.buf, dle.useColor, bBlue, "%03d", status)
	case status < 300:
		cw(dle.buf, dle.useColor, bGreen, "%03d", status)
	case status < 400:
		cw(dle.buf, dle.useColor, bCyan, "%03d", status)
	case status < 500:
		cw(dle.buf, dle.useColor, bYellow, "%03d", status)
	default:
		cw(dle.buf, dle.useColor, bRed, "%03d", status)
	}

	cw(dle.buf, dle.useColor, bBlue, " %dB", bytes)

	dle.buf.WriteString(" in ")
	if elapsed < 500*time.Millisecond {
		cw(dle.buf, dle.useColor, nGreen, "%s", elapsed)
	} else if elapsed < 5*time.Second {
		cw(dle.buf, dle.useColor, nYellow, "%s", elapsed)
	} else {
		cw(dle.buf, dle.useColor, nRed, "%s", elapsed)
	}

	dle.Logger.Print(dle.buf.String())
}

func (dle *defaultLogEntry) Panic(v interface{}, stack []byte) {
	PrintPrettyStack(v)
}

func init() {
	color := true
	if runtime.GOOS == "windows" {
		color = false
	}

	logger := &DefaultLogFormatter{
		Logger:  log.New(os.Stdout, "", log.LstdFlags),
		NoColor: !color,
	}

	DefaultLogger = RequestLogger(logger)
}
