# Go-Debug

A lightweight tracing library for Go packages. Go-Debug allows library authors to add diagnostic logging that end users can enable via a DEBUG environment variable.

## Installation

```bash
go get github.com/telemetryos/godebug/debug
```

## Usage

### For Library Authors

Library authors can add trace logs to their packages using the `Trace` function. This allows for detailed diagnostic information without adding noise for regular users.

```go
package mylib

import "github.com/telemetryos/godebug/debug"

// Create a scope for your library
var dbg = debug.Bind("mylib")

func DoSomething() {
    dbg.Trace("Starting DoSomething")

    // Complex operations...

    dbg.Tracef("Processing %d items", count)
}
```

#### Best Practices

1. Create scopes with meaningful names that identify the component or package
2. Consider using namespacing with colons (e.g., `app:http`, `app:db`) for better organization
3. Add trace calls for entry/exit points and key operations
4. Use consistent naming conventions for your scopes

### For End Users

End users can enable tracing by setting the `DEBUG` environment variable with a comma-separated list of scopes:

```bash
# Enable tracing for specific libraries
DEBUG=mylib,database,http go run main.go

# Enable tracing for namespaced components
DEBUG=app:http,app:db go run main.go

# Enable all tracing with wildcard
DEBUG=* go run main.go
```

#### Pattern Matching

The `DEBUG` environment variable supports several pattern matching techniques:

```bash
# Exact match: enable only the exact scope
DEBUG=app:http go run main.go

# Hierarchical match: enable a scope and all its children
# This enables app, app:http, app:http:get, etc.
DEBUG=app go run main.go

# Prefix match: enable all scopes that start with prefix
# This enables app:http, app:http:get, etc.
DEBUG=app: go run main.go

# Wildcard match: enable all scopes that match the pattern
# This enables app:get, app:post, but not app:http:get
DEBUG=app:* go run main.go

# Middle wildcards: enable all matching scopes
# This enables app:http:get, app:api:get, etc.
DEBUG=app:*:get go run main.go

# Multiple wildcards: match more complex patterns
# This enables app:http:get:v1, app:api:get:v2, etc.
DEBUG=app:*:get:* go run main.go
```

Multiple patterns can be combined with commas:
```bash
DEBUG=app:http,database:*,auth:oauth:* go run main.go
```

Trace logs will be printed to stdout with the scope name as a prefix:

```
mylib: Starting DoSomething
mylib: Processing 5 items
```

## API Reference

### `Bind(scopeName string) *Scope`

Creates a new named scope for tracing.

```go
dbg := debug.Bind("mylib")
```

### `(s *Scope) Trace(message string)`

Logs a simple message if the scope is enabled.

```go
dbg.Trace("Connection established")
```

### `(s *Scope) Tracef(format string, args ...any)`

Logs a formatted message if the scope is enabled. Uses the same formatting rules as `fmt.Printf`.

```go
dbg.Tracef("Processing item %d of %d", i, total)
```

## Example

See the [example](./example/main.go) for a complete demonstration of how to use the package.

To run the example with different trace levels:

```bash
# Show no traces
go run example/main.go

# Show only database traces
DEBUG=app:db go run example/main.go

# Show all traces
DEBUG=* go run example/main.go
```

## Features

- Zero overhead when tracing is disabled
- Simple API with string and formatted string support
- Scope-based filtering for selective debugging
- Advanced pattern matching for flexible debugging control
- No external dependencies
- Namespace support with colons for logical grouping
- Wildcard support with flexible patterns (*, a:*, a:*:c, etc.)

## License

This project is licensed under the MIT License with an Attribution Requirement - see the [LICENSE](LICENSE) file for details.

When using this software in your projects or integrating it into your products, please maintain the original copyright notice and include a link to the original repository (https://github.com/telemetryos/godebug).
