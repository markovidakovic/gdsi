package middleware

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

// NewWrappedResponseWriter wraps an http.ResponseWriter, returning a proxy that allows
// hooking into various parts of the response process
func NewWrappedResponseWriter(w http.ResponseWriter, protoMajor int) WrappedResponseWriter {
	_, fl := w.(http.Flusher)

	bw := basicWriter{ResponseWriter: w}

	if protoMajor == 2 {
		_, ps := w.(http.Pusher)
		if fl && ps {
			return &http2FancyWriter{bw}
		}
	} else {
		_, hj := w.(http.Hijacker)
		_, rf := w.(io.ReaderFrom)
		if fl && hj && rf {
			return &httpFancyWriter{bw}
		}
		if fl && hj {
			return &flushHijackWriter{bw}
		}
		if hj {
			return &hijackWriter{bw}
		}
	}

	if fl {
		return &flushWriter{bw}
	}

	return &bw
}

// WrappedResponseWriter is a proxy around an http.ResponseWriter that allows
// hooking into various parts of the response process
type WrappedResponseWriter interface {
	http.ResponseWriter
	// Status returns the http status of the request, or 0 if one has not
	// yet been sent
	Status() int
	// BytesWritten returns the total number of bytes sent to the client
	BytesWritten() int
	// Tee causes the response body to be written to the given io.Writter in
	// addition to proxying the writes through. Only one io.Writer can be
	// tee'd to at once: setting a second one will overwrite the first.
	// Writes will be sent to the proxy before being written to this io.Writer.
	// It is illegal for the tee'd writer to be modified concurrently with writes
	Tee(io.Writer)
	// Unwrap returns the original proxied target
	Unwrap() http.ResponseWriter
	// Discard causes all writes to the original ResponseWriter be discarded,
	// instead writing only to the tee'd writer if it's set.
	// The caller is responsible for calling WriteHeader and Write on the
	// original ResponseWriter once the processing is done
	Discard()
}

// basicWriter wraps a http.ResponseWriter that implements the minimal
// http.ResponseWriter interface
type basicWriter struct {
	http.ResponseWriter
	wroteHeader bool
	code        int
	bytes       int
	tee         io.Writer
	discard     bool
}

func (bw *basicWriter) WriteHeader(code int) {
	if code >= 100 && code <= 199 && code != http.StatusSwitchingProtocols {
		if !bw.discard {
			bw.ResponseWriter.WriteHeader(code)
		}
	} else if !bw.wroteHeader {
		bw.code = code
		bw.wroteHeader = true
		if !bw.discard {
			bw.ResponseWriter.WriteHeader(code)
		}
	}
}

func (bw *basicWriter) Write(buf []byte) (n int, err error) {
	bw.maybeWriteHeader()

	if !bw.discard {
		n, err = bw.ResponseWriter.Write(buf)
		if bw.tee != nil {
			_, err2 := bw.tee.Write(buf[:n])
			if err == nil {
				err = err2
			}
		}
	} else if bw.tee != nil {
		n, err = bw.tee.Write(buf)
	} else {
		n, err = io.Discard.Write(buf)
	}

	bw.bytes += n
	return
}

func (bw *basicWriter) maybeWriteHeader() {
	if !bw.wroteHeader {
		bw.WriteHeader(http.StatusOK)
	}
}

func (bw *basicWriter) Status() int {
	return bw.code
}

func (bw *basicWriter) BytesWritten() int {
	return bw.bytes
}

func (bw *basicWriter) Tee(w io.Writer) {
	bw.tee = w
}

func (bw *basicWriter) Unwrap() http.ResponseWriter {
	return bw.ResponseWriter
}

func (bw *basicWriter) Discard() {
	bw.discard = true
}

// flushWriter...
type flushWriter struct {
	basicWriter
}

func (fw *flushWriter) Flush() {
	fw.wroteHeader = true
	fl := fw.basicWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}

var _ http.Flusher = &flushWriter{}

// hijackWriter...
type hijackWriter struct {
	basicWriter
}

func (hw *hijackWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := hw.basicWriter.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

var _ http.Hijacker = &hijackWriter{}

// flushHijackWriter...
type flushHijackWriter struct {
	basicWriter
}

func (fhw *flushHijackWriter) Flush() {
	fhw.wroteHeader = true
	fl := fhw.basicWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}

func (fhw *flushHijackWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := fhw.basicWriter.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

var _ http.Flusher = &flushHijackWriter{}
var _ http.Hijacker = &flushHijackWriter{}

// httpFancyWriter is a http writer that additionally satisfies
// http.Flusher, http.Hijacker, and io.ReadFrom. It exists for the common case
// of wrapping the http.ResponseWriter that package http gives, in order to
// make the proxied object support the full method set of the proxied object
type httpFancyWriter struct {
	basicWriter
}

func (hfw *httpFancyWriter) Flush() {
	hfw.wroteHeader = true
	fl := hfw.basicWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}

func (hfw *httpFancyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := hfw.basicWriter.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

func (hfw *httpFancyWriter) ReadFrom(r io.Reader) (int64, error) {
	if hfw.basicWriter.tee != nil {
		n, err := io.Copy(&hfw.basicWriter, r)
		hfw.basicWriter.bytes += int(n)
		return n, err
	}
	rf := hfw.basicWriter.ResponseWriter.(io.ReaderFrom)
	hfw.basicWriter.maybeWriteHeader()
	n, err := rf.ReadFrom(r)
	hfw.basicWriter.bytes += int(n)
	return n, err
}

var _ http.Flusher = &httpFancyWriter{}
var _ http.Hijacker = &httpFancyWriter{}
var _ io.ReaderFrom = &httpFancyWriter{}

// http2FancyWriter is a http2 writer that additionally satisfies
// http.Flusher, and io.ReaderFrom. It exists for the common case
// of wrapping the http.ResponseWriter that package http gives, in order to
// make the proxied object support the full method set of the proxied object
type http2FancyWriter struct {
	basicWriter
}

func (hfw *http2FancyWriter) Flush() {
	hfw.wroteHeader = true
	fl := hfw.basicWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}

func (hfw *http2FancyWriter) Push(target string, opts *http.PushOptions) error {
	return hfw.basicWriter.ResponseWriter.(http.Pusher).Push(target, opts)
}

var _ http.Flusher = &http2FancyWriter{}
var _ http.Pusher = &http2FancyWriter{}
