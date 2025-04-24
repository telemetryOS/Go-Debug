package trace_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/telemetrytv/trace"
)

func captureOutput(f func()) string {
	// Capture stdout
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function
	f()

	// Restore stdout
	w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestBindCreatesScope(t *testing.T) {
	scope := trace.Bind("test")
	assert.NotNil(t, scope, "Bind should return a non-nil Scope")
}

func TestTraceWithoutDebugEnv(t *testing.T) {
	// Ensure DEBUG is not set
	os.Unsetenv("DEBUG")

	output := captureOutput(func() {
		scope := trace.Bind("test")
		scope.Trace("This should not be printed")
	})

	assert.Empty(t, output, "Expected no output with DEBUG unset")
}

func TestTraceWithMatchingDebugEnv(t *testing.T) {
	// Set DEBUG to match our scope
	os.Setenv("DEBUG", "test")
	defer os.Unsetenv("DEBUG")

	output := captureOutput(func() {
		scope := trace.Bind("test")
		scope.Trace("This should be printed")
	})

	expected := "test: This should be printed\n"
	assert.Equal(t, expected, output, "Expected matching output with DEBUG set to scope name")
}

func TestTraceWithNonMatchingDebugEnv(t *testing.T) {
	// Set DEBUG to not match our scope
	os.Setenv("DEBUG", "other")
	defer os.Unsetenv("DEBUG")

	output := captureOutput(func() {
		scope := trace.Bind("test")
		scope.Trace("This should not be printed")
	})

	assert.Empty(t, output, "Expected no output with non-matching DEBUG scope")
}

func TestTraceWithWildcardDebugEnv(t *testing.T) {
	// Set DEBUG to wildcard
	os.Setenv("DEBUG", "*")
	defer os.Unsetenv("DEBUG")

	output := captureOutput(func() {
		scope := trace.Bind("test")
		scope.Trace("This should be printed")
	})

	expected := "test: This should be printed\n"
	assert.Equal(t, expected, output, "Expected matching output with DEBUG set to wildcard")
}

func TestTraceWithMultipleScopesInDebugEnv(t *testing.T) {
	// Set DEBUG with multiple scopes
	os.Setenv("DEBUG", "other,test,third")
	defer os.Unsetenv("DEBUG")

	output := captureOutput(func() {
		scope := trace.Bind("test")
		scope.Trace("This should be printed")
	})

	expected := "test: This should be printed\n"
	assert.Equal(t, expected, output, "Expected matching output with DEBUG set to multiple scopes")
}

func TestTracefWithMatchingDebugEnv(t *testing.T) {
	// Set DEBUG to match our scope
	os.Setenv("DEBUG", "test")
	defer os.Unsetenv("DEBUG")

	output := captureOutput(func() {
		scope := trace.Bind("test")
		scope.Tracef("Count: %d", 42)
	})

	expected := "test: Count: 42\n"
	assert.Equal(t, expected, output, "Expected formatted output with DEBUG set to scope name")
}

func TestTracefWithoutDebugEnv(t *testing.T) {
	// Ensure DEBUG is not set
	os.Unsetenv("DEBUG")

	output := captureOutput(func() {
		scope := trace.Bind("test")
		scope.Tracef("Count: %d", 42)
	})

	assert.Empty(t, output, "Expected no output with DEBUG unset")
}

func TestDebugWithWhitespace(t *testing.T) {
	// Set DEBUG with whitespace
	os.Setenv("DEBUG", " test , other ")
	defer os.Unsetenv("DEBUG")

	output := captureOutput(func() {
		scope := trace.Bind("test")
		scope.Trace("This should be printed")
	})

	expected := "test: This should be printed\n"
	assert.Equal(t, expected, output, "Expected matching output with DEBUG containing whitespace")
}
