package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestUIFunctions(t *testing.T) {
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	OutputWriter = &outBuf
	ErrorWriter = &errBuf

	Success("Task completed successfully")
	if !strings.Contains(outBuf.String(), "Task completed successfully") {
		t.Errorf("expected success message in output, got: %s", outBuf.String())
	}
	outBuf.Reset()

	Info("Processing files")
	if !strings.Contains(outBuf.String(), "Processing files") {
		t.Errorf("expected info message in output, got: %s", outBuf.String())
	}
	outBuf.Reset()

	Warn("Warning occurred")
	if !strings.Contains(outBuf.String(), "Warning occurred") {
		t.Errorf("expected warning message in output, got: %s", outBuf.String())
	}
	outBuf.Reset()

	Error("Error occurred")
	if !strings.Contains(errBuf.String(), "Error occurred") {
		t.Errorf("expected error message in stderr, got: %s", errBuf.String())
	}
	errBuf.Reset()

	Step(1, 5, "Initialize")
	if !strings.Contains(outBuf.String(), "[1/5] Initialize") {
		t.Errorf("expected step header in output, got: %s", outBuf.String())
	}
	outBuf.Reset()

	Banner()
	if !strings.Contains(outBuf.String(), "create-go-app CLI") {
		t.Errorf("expected banner in output, got: %s", outBuf.String())
	}
}
