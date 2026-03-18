package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "apialerts-functional-*")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}
	defer os.RemoveAll(tmp)

	binaryPath = filepath.Join(tmp, "apialerts")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic("failed to build binary: " + err.Error())
	}

	os.Exit(m.Run())
}

type result struct {
	stdout   string
	stderr   string
	exitCode int
}

func run(t *testing.T, homeDir string, args ...string) result {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	cmd.Env = []string{"HOME=" + homeDir}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	exitCode := 0
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}
	return result{
		stdout:   strings.TrimSpace(stdout.String()),
		stderr:   strings.TrimSpace(stderr.String()),
		exitCode: exitCode,
	}
}

func mockServer(t *testing.T, statusCode int, body map[string]any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if body != nil {
			json.NewEncoder(w).Encode(body)
		}
	}))
}

// --- Config tests ---

func TestConfigNoKey(t *testing.T) {
	home := t.TempDir()
	r := run(t, home, "config")
	if !strings.Contains(r.stdout, "No API key configured") {
		t.Errorf("expected 'No API key configured', got: %q", r.stdout)
	}
	if !strings.Contains(r.stdout, "apialerts init") {
		t.Errorf("expected hint to run 'apialerts init', got: %q", r.stdout)
	}
}

func TestConfigSetKey(t *testing.T) {
	home := t.TempDir()
	r := run(t, home, "config", "--key", "testapikey12345678")
	if !strings.Contains(r.stdout, "API key saved") {
		t.Errorf("expected 'API key saved', got: %q", r.stdout)
	}
}

func TestConfigViewMaskedKey(t *testing.T) {
	home := t.TempDir()
	run(t, home, "config", "--key", "testapikey12345678")
	r := run(t, home, "config")
	if !strings.Contains(r.stdout, "API Key:") {
		t.Errorf("expected masked key, got: %q", r.stdout)
	}
	if strings.Contains(r.stdout, "testapikey12345678") {
		t.Errorf("expected key to be masked, got full key: %q", r.stdout)
	}
}

func TestConfigUnsetKey(t *testing.T) {
	home := t.TempDir()
	run(t, home, "config", "--key", "testapikey12345678")
	r := run(t, home, "config", "--unset")
	if !strings.Contains(r.stdout, "API key removed") {
		t.Errorf("expected 'API key removed', got: %q", r.stdout)
	}
	r = run(t, home, "config")
	if !strings.Contains(r.stdout, "No API key configured") {
		t.Errorf("expected 'No API key configured' after unset, got: %q", r.stdout)
	}
}

// --- Init tests ---

func TestInitNoTTY(t *testing.T) {
	home := t.TempDir()
	r := run(t, home, "init")
	if r.exitCode == 0 {
		t.Error("expected non-zero exit code when no TTY")
	}
	if !strings.Contains(r.stderr, "no terminal detected") {
		t.Errorf("expected 'no terminal detected' error, got: %q", r.stderr)
	}
}

// --- Send validation tests ---

func TestSendNoMessage(t *testing.T) {
	home := t.TempDir()
	r := run(t, home, "send")
	if r.exitCode == 0 {
		t.Error("expected non-zero exit code when no message")
	}
	if !strings.Contains(r.stderr, "message is required") {
		t.Errorf("expected 'message is required', got: %q", r.stderr)
	}
}

func TestSendEmptyMessage(t *testing.T) {
	home := t.TempDir()
	r := run(t, home, "send", "-m", "")
	if r.exitCode == 0 {
		t.Error("expected non-zero exit code when message is empty")
	}
	if !strings.Contains(r.stderr, "message is required") {
		t.Errorf("expected 'message is required', got: %q", r.stderr)
	}
}

func TestSendNoKey(t *testing.T) {
	home := t.TempDir()
	r := run(t, home, "send", "-m", "hello")
	if r.exitCode == 0 {
		t.Error("expected non-zero exit code when no API key")
	}
	if !strings.Contains(r.stderr, "no API key configured") {
		t.Errorf("expected 'no API key configured', got: %q", r.stderr)
	}
}

// --- Send HTTP tests ---

func TestSendSuccess(t *testing.T) {
	server := mockServer(t, http.StatusOK, map[string]any{
		"workspace": "Acme Corp",
		"channel":   "general",
	})
	defer server.Close()

	home := t.TempDir()
	run(t, home, "config", "--server-url", server.URL)
	r := run(t, home, "send", "-m", "Deploy complete", "--key", "fake-api-key")
	if r.exitCode != 0 {
		t.Errorf("expected success, got exit code %d, stderr: %q", r.exitCode, r.stderr)
	}
	if !strings.Contains(r.stdout, "Acme Corp") {
		t.Errorf("expected workspace in output, got: %q", r.stdout)
	}
	if !strings.Contains(r.stdout, "general") {
		t.Errorf("expected channel in output, got: %q", r.stdout)
	}
}

func TestSendUnauthorized(t *testing.T) {
	server := mockServer(t, http.StatusUnauthorized, nil)
	defer server.Close()

	home := t.TempDir()
	run(t, home, "config", "--server-url", server.URL)
	r := run(t, home, "send", "-m", "Deploy complete", "--key", "bad-key")
	if r.exitCode == 0 {
		t.Error("expected non-zero exit code for unauthorized")
	}
	if !strings.Contains(r.stderr, "unauthorized") {
		t.Errorf("expected 'unauthorized' in error, got: %q", r.stderr)
	}
}

func TestSendRateLimit(t *testing.T) {
	server := mockServer(t, http.StatusTooManyRequests, nil)
	defer server.Close()

	home := t.TempDir()
	run(t, home, "config", "--server-url", server.URL)
	r := run(t, home, "send", "-m", "Deploy complete", "--key", "fake-api-key")
	if r.exitCode == 0 {
		t.Error("expected non-zero exit code for rate limit")
	}
	if !strings.Contains(r.stderr, "rate limit") {
		t.Errorf("expected 'rate limit' in error, got: %q", r.stderr)
	}
}

func TestSendWithData(t *testing.T) {
	server := mockServer(t, http.StatusOK, map[string]any{
		"workspace": "Acme Corp",
		"channel":   "general",
	})
	defer server.Close()

	home := t.TempDir()
	run(t, home, "config", "--server-url", server.URL)
	r := run(t, home, "send", "-m", "New signup", "-d", `{"plan":"pro","source":"organic"}`, "--key", "fake-api-key")
	if r.exitCode != 0 {
		t.Errorf("expected success, got exit code %d, stderr: %q", r.exitCode, r.stderr)
	}
	if !strings.Contains(r.stdout, "Acme Corp") {
		t.Errorf("expected workspace in output, got: %q", r.stdout)
	}
}

func TestSendWithInvalidData(t *testing.T) {
	home := t.TempDir()
	r := run(t, home, "send", "-m", "hello", "-d", `not valid json`, "--key", "fake-api-key")
	if r.exitCode == 0 {
		t.Error("expected non-zero exit code for invalid JSON")
	}
	if !strings.Contains(r.stderr, "invalid JSON") {
		t.Errorf("expected 'invalid JSON' error, got: %q", r.stderr)
	}
}

// --- Test command HTTP tests ---

func TestTestCommandSuccess(t *testing.T) {
	server := mockServer(t, http.StatusOK, map[string]any{
		"workspace": "Acme Corp",
		"channel":   "general",
	})
	defer server.Close()

	home := t.TempDir()
	run(t, home, "config", "--key", "fake-api-key")
	run(t, home, "config", "--server-url", server.URL)
	r := run(t, home, "test")
	if r.exitCode != 0 {
		t.Errorf("expected success, got exit code %d, stderr: %q", r.exitCode, r.stderr)
	}
	if !strings.Contains(r.stdout, "Acme Corp") {
		t.Errorf("expected workspace in output, got: %q", r.stdout)
	}
}

func TestTestCommandNoKey(t *testing.T) {
	home := t.TempDir()
	r := run(t, home, "test")
	if r.exitCode == 0 {
		t.Error("expected non-zero exit code when no API key")
	}
	if !strings.Contains(r.stderr, "no API key configured") {
		t.Errorf("expected 'no API key configured', got: %q", r.stderr)
	}
}
