// Package trace provides a lightweight tracing system for Go libraries and applications.
// It allows library authors to add diagnostic logs that can be enabled by end users
// via the DEBUG environment variable.
package trace

import (
	"fmt"
	"os"
	"strings"
)

// Scope represents a named logging scope with trace capabilities.
// Each scope can be independently enabled by including its name in the DEBUG environment variable.
type Scope struct {
	scopeName string
}

// Bind creates a new named scope for tracing.
// The scopeName parameter is used to match against the DEBUG environment variable.
func Bind(scopeName string) *Scope {
	return &Scope{
		scopeName: scopeName,
	}
}

// Trace logs a simple message if the scope is enabled in the DEBUG environment variable.
func (s *Scope) Trace(message string) {
	s.maybeLog(message)
}

// Tracef logs a formatted message if the scope is enabled in the DEBUG environment variable.
// It accepts format specifiers and arguments similar to fmt.Printf.
func (s *Scope) Tracef(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	s.maybeLog(message)
}

// maybeLog checks if the current scope is enabled and outputs the message if it is.
// The DEBUG environment variable is checked each time, allowing for runtime updates.
func (s *Scope) maybeLog(message string) {
	debug := os.Getenv("DEBUG")
	if debug == "" {
		return
	}
	
	if debug == "*" {
		fmt.Printf("%s: %s\n", s.scopeName, message)
		return
	}
	
	scopes := strings.Split(debug, ",")
	for _, scope := range scopes {
		scope = strings.TrimSpace(scope)
		if scope == s.scopeName {
			fmt.Printf("%s: %s\n", s.scopeName, message)
			return
		}
	}
}
