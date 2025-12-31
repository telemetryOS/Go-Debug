package debug_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/telemetryos/godebug/debug"
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
	scope := debug.Bind("test")
	assert.NotNil(t, scope, "Bind should return a non-nil Scope")
}

func TestTraceWithoutDebugEnv(t *testing.T) {
	// Ensure DEBUG is not set
	os.Unsetenv("DEBUG")

	output := captureOutput(func() {
		scope := debug.Bind("test")
		scope.Trace("This should not be printed")
	})

	assert.Empty(t, output, "Expected no output with DEBUG unset")
}

func TestTraceWithMatchingDebugEnv(t *testing.T) {
	// Set DEBUG to match our scope
	os.Setenv("DEBUG", "test")
	defer os.Unsetenv("DEBUG")

	output := captureOutput(func() {
		scope := debug.Bind("test")
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
		scope := debug.Bind("test")
		scope.Trace("This should not be printed")
	})

	assert.Empty(t, output, "Expected no output with non-matching DEBUG scope")
}

func TestTraceWithWildcardDebugEnv(t *testing.T) {
	// Set DEBUG to wildcard
	os.Setenv("DEBUG", "*")
	defer os.Unsetenv("DEBUG")

	output := captureOutput(func() {
		scope := debug.Bind("test")
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
		scope := debug.Bind("test")
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
		scope := debug.Bind("test")
		scope.Tracef("Count: %d", 42)
	})

	expected := "test: Count: 42\n"
	assert.Equal(t, expected, output, "Expected formatted output with DEBUG set to scope name")
}

func TestTracefWithoutDebugEnv(t *testing.T) {
	// Ensure DEBUG is not set
	os.Unsetenv("DEBUG")

	output := captureOutput(func() {
		scope := debug.Bind("test")
		scope.Tracef("Count: %d", 42)
	})

	assert.Empty(t, output, "Expected no output with DEBUG unset")
}

func TestDebugWithWhitespace(t *testing.T) {
	// Set DEBUG with whitespace
	os.Setenv("DEBUG", " test , other ")
	defer os.Unsetenv("DEBUG")

	output := captureOutput(func() {
		scope := debug.Bind("test")
		scope.Trace("This should be printed")
	})

	expected := "test: This should be printed\n"
	assert.Equal(t, expected, output, "Expected matching output with DEBUG containing whitespace")
}

func TestPatternMatching(t *testing.T) {
	tests := []struct {
		debugValue  string
		scopeName   string
		shouldMatch bool
		message     string
	}{
		// Hierarchical matching - 'a' matches 'a', 'a:b', 'a:b:c', etc.
		{"a", "a", true, "Debug 'a' should match scope 'a'"},
		{"a", "a:b", true, "Debug 'a' should match scope 'a:b'"},
		{"a", "a:b:c", true, "Debug 'a' should match scope 'a:b:c'"},
		{"a", "b:c", false, "Debug 'a' should not match scope 'b:c'"},
		
		// Prefix matching - 'a:' matches 'a:b', 'a:c', 'a:b:c', etc.
		{"a:", "a:b", true, "Debug 'a:' should match scope 'a:b'"},
		{"a:", "a:b:c", true, "Debug 'a:' should match scope 'a:b:c'"},
		{"a:", "b:c", false, "Debug 'a:' should not match scope 'b:c'"},
		
		// Simple wildcard matching
		{"a:*", "a:b", true, "Debug 'a:*' should match scope 'a:b'"},
		{"a:*", "a:b:c", true, "Debug 'a:*' should match scope 'a:b:c'"},
		{"a:*", "b:c", false, "Debug 'a:*' should not match scope 'b:c'"},
		
		// Middle wildcard matching - 'a:*:c' matches 'a:b:c', 'a:x:c', etc.
		{"a:*:c", "a:b:c", true, "Debug 'a:*:c' should match scope 'a:b:c'"},
		{"a:*:c", "a:xyz:c", true, "Debug 'a:*:c' should match scope 'a:xyz:c'"},
		{"a:*:c", "a:b:c:d", false, "Debug 'a:*:c' should not match scope 'a:b:c:d'"},
		{"a:*:c", "a:b:d", false, "Debug 'a:*:c' should not match scope 'a:b:d'"},
		
		// Multiple wildcards - 'a:*:*:d' matches 'a:b:c:d', 'a:x:y:d', etc.
		{"a:*:*:d", "a:b:c:d", true, "Debug 'a:*:*:d' should match scope 'a:b:c:d'"},
		{"a:*:*:d", "a:x:y:d", true, "Debug 'a:*:*:d' should match scope 'a:x:y:d'"},
		{"a:*:*:d", "a:b:c:e", false, "Debug 'a:*:*:d' should not match scope 'a:b:c:e'"},
		
		// Leading wildcard - '*:c' matches 'a:c', 'b:c', etc.
		{"*:c", "a:c", true, "Debug '*:c' should match scope 'a:c'"},
		{"*:c", "b:c", true, "Debug '*:c' should match scope 'b:c'"},
		{"*:c", "a:b:c", false, "Debug '*:c' should not match scope 'a:b:c'"},
		
		// Wildcard at both ends - 'a:*:c:*' matches 'a:b:c:d', 'a:x:c:y', etc.
		{"a:*:c:*", "a:b:c:d", true, "Debug 'a:*:c:*' should match scope 'a:b:c:d'"},
		{"a:*:c:*", "a:x:c:y", true, "Debug 'a:*:c:*' should match scope 'a:x:c:y'"},
		{"a:*:c:*", "a:b:d:e", false, "Debug 'a:*:c:*' should not match scope 'a:b:d:e'"},
	}
	
	for _, test := range tests {
		os.Setenv("DEBUG", test.debugValue)
		defer os.Unsetenv("DEBUG")
		
		output := captureOutput(func() {
			scope := debug.Bind(test.scopeName)
			scope.Trace("test message")
		})
		
		if test.shouldMatch {
			expected := test.scopeName + ": test message\n"
			assert.Equal(t, expected, output, test.message)
		} else {
			assert.Empty(t, output, test.message)
		}
	}
}
