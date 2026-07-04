package debug

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestEnvEnablesDebug(t *testing.T) {
	tests := []struct {
		value string
		want  bool
	}{
		{"", false},
		{"false", false},
		{"0", false},
		{"no", false},
		{"1", true},
		{"true", true},
		{"TRUE", true},
		{"yes", true},
		{"Yes", true},
	}

	for _, tt := range tests {
		if got := EnvEnablesDebug(tt.value); got != tt.want {
			t.Errorf("EnvEnablesDebug(%q) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestEnabled(t *testing.T) {
	SetEnabled(false)
	if Enabled() {
		t.Fatal("expected debug disabled")
	}

	SetEnabled(true)
	if !Enabled() {
		t.Fatal("expected debug enabled")
	}
}

func TestLogRespectsEnabled(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	SetEnabled(false)
	Log("hidden")
	if buf.Len() != 0 {
		t.Fatalf("expected no output when disabled, got %q", buf.String())
	}

	SetEnabled(true)
	Log("visible")
	if !strings.Contains(buf.String(), "visible") {
		t.Fatalf("expected log output to contain message, got %q", buf.String())
	}
}

func TestLogfRespectsEnabled(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	SetEnabled(true)
	Logf("count=%d", 42)
	if !strings.Contains(buf.String(), "count=42") {
		t.Fatalf("expected formatted log output, got %q", buf.String())
	}
}

func TestPrintfRespectsEnabled(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = stdout })

	SetEnabled(false)
	Printf("hidden\n")
	w.Close()
	out, _ := io.ReadAll(r)
	if len(out) != 0 {
		t.Fatalf("expected no stdout when disabled, got %q", out)
	}

	r, w, err = os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	SetEnabled(true)
	Printf("visible=%s\n", "ok")
	w.Close()
	out, _ = io.ReadAll(r)
	if !strings.Contains(string(out), "visible=ok") {
		t.Fatalf("expected stdout output, got %q", out)
	}
}

func TestStack(t *testing.T) {
	stack := Stack()
	if len(stack) == 0 {
		t.Fatal("expected non-empty stack trace")
	}
	if !strings.Contains(string(stack), "TestStack") {
		t.Fatalf("expected stack to mention TestStack, got %q", stack)
	}
}

func TestPrintStack(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stderr := os.Stderr
	os.Stderr = w
	t.Cleanup(func() { os.Stderr = stderr })

	PrintStack()
	w.Close()
	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), "TestPrintStack") {
		t.Fatalf("expected stack trace on stderr, got %q", out)
	}
}

func TestStartPprof(t *testing.T) {
	SetEnabled(false)

	srv, err := StartPprof("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = srv.Close() })

	deadline := time.Now().Add(2 * time.Second)
	for {
		resp, err := http.Get("http://" + srv.Addr + "/debug/pprof/")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected status 200, got %d", resp.StatusCode)
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("pprof server did not become ready: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestStartPprofDefaultAddr(t *testing.T) {
	srv, err := StartPprof("")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = srv.Close() })

	_, port, err := net.SplitHostPort(srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	if port != "6060" {
		t.Fatalf("expected default port 6060, got %q", port)
	}
}
