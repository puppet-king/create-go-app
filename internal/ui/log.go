package ui

import (
	"fmt"
	"io"
	"os"
)

var (
	// OutputWriter allows redirecting output in tests
	OutputWriter io.Writer = os.Stdout
	ErrorWriter  io.Writer = os.Stderr
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// Success prints a green success message
func Success(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(OutputWriter, "%s%s✓ %s%s\n", colorBold, colorGreen, msg, colorReset)
}

// Info prints a cyan informative message
func Info(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(OutputWriter, "%sℹ %s%s\n", colorCyan, msg, colorReset)
}

// Warn prints a yellow warning message
func Warn(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(OutputWriter, "%s%s⚠ %s%s\n", colorBold, colorYellow, msg, colorReset)
}

// Error prints a red error message to stderr
func Error(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(ErrorWriter, "%s%s✗ %s%s\n", colorBold, colorRed, msg, colorReset)
}

// Step prints a step header with total steps
func Step(current, total int, title string) {
	fmt.Fprintf(OutputWriter, "\n%s[%d/%d] %s%s\n", colorBold, current, total, title, colorReset)
}

// Banner prints a stylish CLI welcome banner
func Banner() {
	fmt.Fprintf(OutputWriter, `
%s%s╔══════════════════════════════════════════╗
║           create-go-app CLI              ║
║       Go Application Scaffolder          ║
╚══════════════════════════════════════════╝%s

`, colorBold, colorCyan, colorReset)
}
