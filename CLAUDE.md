# Go-Debug

This file provides guidance to AI coding agents (Claude Code, Cursor, etc.) when working with code in this repository.

Lightweight tracing library for Go packages. Allows library authors to add diagnostic logging that end users enable via the `DEBUG` environment variable.

## Commands

```bash
# Run all tests
go test ./...

# Run single test
go test ./debug -run TestPatternMatching

# Run tests with verbose output
go test -v ./...
```

## Architecture

Single-package library in `debug/` with three core components:

- **Scope** - Named logging context created via `Bind("scopeName")`
- **Trace/Tracef** - Output methods that check DEBUG env var on each call
- **Pattern matching** - Hierarchical scope matching with wildcards

### Pattern Matching Logic

The `maybeLog` function checks DEBUG env var each invocation (allows runtime changes). Matching precedence:

1. Empty DEBUG = no output
2. `DEBUG=*` = match all
3. Exact match (`test` matches `test`)
4. Prefix match (`app:` matches `app:http`)
5. Hierarchical match (`app` matches `app:http:get`)
6. Wildcard patterns via `matchPattern()`

Wildcard matching supports:
- Trailing wildcard: `a:*` matches `a:b`, `a:b:c`
- Middle wildcard: `a:*:c` matches `a:b:c`
- Multiple wildcards: `a:*:*:d` matches `a:b:c:d`

## Patterns

**Scope naming**: Use colon-delimited hierarchies (`app:http:get`) for selective filtering.

**Zero overhead**: When DEBUG is unset, only a single `os.Getenv` call occurs per trace.

**Testing**: Tests use `captureOutput()` helper that redirects stdout via `os.Pipe()`. Always `defer os.Unsetenv("DEBUG")` after setting.
