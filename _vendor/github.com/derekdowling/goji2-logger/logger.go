package gojilogger

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"goji.io/pattern"

	"goji.io"

	"golang.org/x/net/context"

	"github.com/derekdowling/go-stdlogger"
	"github.com/zenazn/goji/web/mutil"
)

const (
	// FastResponse is anything under this duration
	FastResponse = 500 * time.Millisecond
	// AcceptableResponse is anything under this duration
	AcceptableResponse = 5 * time.Second
)

// Logger contains instance state for a goji2logger to avoid configuration
// collisions if this middleware is used in multiple places
type Logger struct {
	// Logger is a https://github.com/derekdowling/go-stdlogger
	Logger std.Logger
	// Debug will increase verbosity of logging information, and causes Query params not
	// to be omitted. Do NOT use in production otherwise you risk logging sensitive
	// information.
	Debug bool
}

// New creates a new goji2logger instance
func New(logger std.Logger, debug bool) *Logger {
	l := &Logger{
		Debug:  debug,
		Logger: logger,
	}

	if l.Logger == nil {
		l.Logger = log.New(os.Stderr, "", log.LstdFlags)
	}

	return l
}

// Middleware logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return. When standard output is a TTY, Logger will
// print in color, otherwise it will print in black and white.
//
// Use like so with Goji2:
//	gLogger := gojilogger.New(nil, false)
//	yourGoji.UseC(gLogger.Middleware)
func (l *Logger) Middleware(next goji.Handler) goji.Handler {
	middleware := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		l.printRequest(ctx, r)

		// WrapWriter lets us peek at ResponseWriter outputs
		lw := mutil.WrapWriter(w)

		startTime := time.Now()
		next.ServeHTTPC(ctx, lw, r)

		if lw.Status() == 0 {
			lw.WriteHeader(http.StatusOK)
		}

		finishTime := time.Now()

		l.printResponse(lw, finishTime.Sub(startTime))
	}

	return goji.HandlerFunc(middleware)
}

func (l *Logger) printRequest(ctx context.Context, r *http.Request) {
	var buf bytes.Buffer

	if l.Debug {
		buf.WriteString("[DEBUG]")
	}

	buf.WriteString("Serving route: ")

	// Goji routing details
	colorWrite(&buf, bGreen, "%s", pattern.Path(ctx))

	// Server details
	buf.WriteString(fmt.Sprintf(" from %s ", r.RemoteAddr))

	// Request details
	buf.WriteString("for ")
	colorWrite(&buf, bMagenta, "%s ", r.Method)

	urlStr := r.URL.String()

	// if not in debug mode, remove Query params from logging as not to include any
	// sensitive information inadvertantly into user's logs
	if !l.Debug && r.URL.RawQuery != "" {
		tempURL := &url.URL{}
		*tempURL = *r.URL
		tempURL.RawQuery = "<omitted>"
		urlStr = tempURL.String()
	}

	colorWrite(&buf, bBlue, "%q", urlStr)
	log.Print(buf.String())
}

func (l *Logger) printResponse(w mutil.WriterProxy, delta time.Duration) {
	var buf bytes.Buffer

	buf.WriteString("Returning HTTP ")

	status := w.Status()
	if status < 200 {
		colorWrite(&buf, bBlue, "%03d", status)
	} else if status < 300 {
		colorWrite(&buf, bGreen, "%03d", status)
	} else if status < 400 {
		colorWrite(&buf, bCyan, "%03d", status)
	} else if status < 500 {
		colorWrite(&buf, bYellow, "%03d", status)
	} else {
		colorWrite(&buf, bRed, "%03d", status)
	}

	buf.WriteString(" in ")

	if delta < FastResponse {
		colorWrite(&buf, nGreen, "%s", delta.String())
	} else if delta < AcceptableResponse {
		colorWrite(&buf, nYellow, "%s", delta.String())
	} else {
		colorWrite(&buf, nRed, "%s", delta.String())
	}

	log.Print(buf.String())
}
