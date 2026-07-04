package debug

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
)

var enabled atomic.Bool

func init() {
	if EnvEnablesDebug(os.Getenv("DEBUG")) {
		enabled.Store(true)
	}
}

// EnvEnablesDebug reports whether the DEBUG env value turns debug mode on.
func EnvEnablesDebug(value string) bool {
	switch strings.ToLower(value) {
	case "1", "true", "yes":
		return true
	default:
		return false
	}
}

// Enabled reports whether debug mode is active.
func Enabled() bool {
	return enabled.Load()
}

// SetEnabled turns debug mode on or off at runtime.
func SetEnabled(on bool) {
	enabled.Store(on)
}

// Log writes a debug message when debug mode is enabled.
func Log(v ...any) {
	if enabled.Load() {
		log.Print(append([]any{"[debug]"}, v...)...)
	}
}

// Logf writes a formatted debug message when debug mode is enabled.
func Logf(format string, v ...any) {
	if enabled.Load() {
		log.Printf("[debug] "+format, v...)
	}
}

// Printf writes a formatted message to stdout when debug mode is enabled.
func Printf(format string, v ...any) {
	if enabled.Load() {
		fmt.Printf("[debug] "+format, v...)
	}
}

// Stack returns the current goroutine stack trace.
func Stack() []byte {
	buf := make([]byte, 4096)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}

// PrintStack writes the current goroutine stack trace to stderr.
func PrintStack() {
	fmt.Fprintln(os.Stderr, string(Stack()))
}

// StartPprof starts an HTTP server exposing Go pprof endpoints.
// Returns the server so callers can shut it down when done.
func StartPprof(addr string) (*http.Server, error) {
	if addr == "" {
		addr = ":6060"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	srv := &http.Server{Addr: ln.Addr().String()}
	go func() {
		Logf("pprof listening on http://localhost%s/debug/pprof/", srv.Addr)
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Printf("[debug] pprof server error: %v", err)
		}
	}()
	return srv, nil
}
